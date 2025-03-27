package composite

import (
	v1 "github.com/JojoWeyn/duo-proj/course-service/internal/controller/http/v1"
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/http/v1/admin"
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/kafka"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/usecase"
	"github.com/JojoWeyn/duo-proj/course-service/internal/repository/cache"
	"github.com/JojoWeyn/duo-proj/course-service/internal/repository/db/postgres"
	"github.com/JojoWeyn/duo-proj/course-service/internal/service"
	"github.com/JojoWeyn/duo-proj/course-service/pkg/client/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Config struct {
	GatewayURL   string
	RedisURL     string
	RedisDB      int
	KafkaBrokers string
}

type CourseComposite struct {
	handler *gin.Engine
}

func NewCourseComposite(db *gorm.DB, cfg Config) (*CourseComposite, error) {

	if err := db.AutoMigrate(
		&entity.Course{},
		&entity.Lesson{},
		&entity.Exercise{},
		&entity.Question{},
		&entity.QuestionImage{},
		&entity.QuestionOption{},
		&entity.QuestionType{},
		&entity.MatchingPair{},
		&entity.Attempt{},
		&entity.AttemptSession{},
		&postgres.Completion{},
	); err != nil {
		return nil, err
	}

	courseRepo := postgres.NewCourseRepository(db)
	exerciseRepo := postgres.NewExerciseRepository(db)
	questionRepo := postgres.NewQuestionRepository(db)
	lessonRepo := postgres.NewLessonRepository(db)
	attemptRepo := postgres.NewAttemptRepository(db)
	completionRepo := postgres.NewCompletionRepo(db)
	matchingPairRepo := postgres.NewMatchingPairRepository(db)
	questionOptionsRepo := postgres.NewQuestionOptionRepository(db)

	redisClient, err := redis.NewRedisClient(redis.Config{
		Addr: cfg.RedisURL,
		DB:   cfg.RedisDB,
	})
	if err != nil {
		return nil, err
	}

	progressProducer, err := kafka.NewProducer(cfg.KafkaBrokers, "user_progress")
	if err != nil {
		return nil, err
	}

	courseCache := cache.NewRedisCache(redisClient)
	attemptService := service.NewAttemptService(questionRepo, exerciseRepo, attemptRepo, lessonRepo, completionRepo)

	handler := gin.Default()

	courseUseCase := usecase.NewCourseUseCase(courseRepo, courseCache)
	exerciseUseCase := usecase.NewExerciseUseCase(exerciseRepo)
	questionUseCase := usecase.NewQuestionUseCase(questionRepo, progressProducer, attemptService)
	lessonUseCase := usecase.NewLessonUseCase(lessonRepo)
	attemptUseCase := usecase.NewAttemptUseCase(attemptRepo)

	matchingPairUseCase := usecase.NewMatchingPairUseCase(matchingPairRepo)
	questionOptionUseCase := usecase.NewQuestionOptionUseCase(questionOptionsRepo)

	v1.NewRouter(handler, courseUseCase, lessonUseCase, exerciseUseCase, questionUseCase, attemptUseCase, cfg.GatewayURL)
	admin.NewRouter(handler, courseUseCase, lessonUseCase, exerciseUseCase, questionUseCase, matchingPairUseCase, questionOptionUseCase, "*")
	return &CourseComposite{
		handler: handler,
	}, nil
}

func (c *CourseComposite) Handler() *gin.Engine {
	return c.handler
}
