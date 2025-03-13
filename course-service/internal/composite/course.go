package composite

import (
	v1 "github.com/JojoWeyn/duo-proj/course-service/internal/controller/http/v1"
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/kafka"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/usecase"
	"github.com/JojoWeyn/duo-proj/course-service/internal/repository/cache"
	"github.com/JojoWeyn/duo-proj/course-service/internal/repository/db/postgres"
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
		&entity.QuestionOption{},
		&entity.QuestionType{},
		&entity.MatchingPair{},
		&entity.Attempt{}); err != nil {
		return nil, err
	}

	courseRepo := postgres.NewCourseRepository(db)
	exerciseRepo := postgres.NewExerciseRepository(db)
	questionRepo := postgres.NewQuestionRepository(db)
	lessonRepo := postgres.NewLessonRepository(db)
	attemptRepo := postgres.NewAttemptRepository(db)

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

	courseCache := cache.NewCourseCache(redisClient)

	handler := gin.Default()

	courseUseCase := usecase.NewCourseUseCase(courseRepo, courseCache)
	exerciseUseCase := usecase.NewExerciseUseCase(exerciseRepo)
	questionUseCase := usecase.NewQuestionUseCase(questionRepo, progressProducer)
	lessonUseCase := usecase.NewLessonUseCase(lessonRepo)
	attemptUseCase := usecase.NewAttemptUseCase(attemptRepo)

	v1.NewRouter(handler, courseUseCase, lessonUseCase, exerciseUseCase, questionUseCase, attemptUseCase, cfg.GatewayURL)

	return &CourseComposite{
		handler: handler,
	}, nil
}

func (c *CourseComposite) Handler() *gin.Engine {
	return c.handler
}
