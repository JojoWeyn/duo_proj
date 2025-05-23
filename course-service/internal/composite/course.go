package composite

import (
	"context"

	v1 "github.com/JojoWeyn/duo-proj/course-service/internal/controller/http/v1"
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/http/v1/admin"
	"github.com/JojoWeyn/duo-proj/course-service/internal/controller/kafka"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/entity"
	"github.com/JojoWeyn/duo-proj/course-service/internal/domain/usecase"
	"github.com/JojoWeyn/duo-proj/course-service/internal/repository/cache"
	"github.com/JojoWeyn/duo-proj/course-service/internal/repository/db/postgres"
	s3Repo "github.com/JojoWeyn/duo-proj/course-service/internal/repository/db/s3"
	"github.com/JojoWeyn/duo-proj/course-service/internal/service"
	"github.com/JojoWeyn/duo-proj/course-service/pkg/client/redis"
	"github.com/JojoWeyn/duo-proj/course-service/pkg/client/s3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Config struct {
	GatewayURL   string
	RedisURL     string
	RedisDB      int
	KafkaBrokers string

	S3Endpoint  string
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
}

type CourseComposite struct {
	handler *gin.Engine
}

func NewCourseComposite(ctx context.Context, db *gorm.DB, cfg Config) (*CourseComposite, error) {

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
		&entity.CourseFile{},
		&entity.LessonFile{},
		&entity.ExerciseFile{},
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

	s3Client, err := s3.NewS3Client(
		cfg.S3Endpoint,
		cfg.S3AccessKey,
		cfg.S3SecretKey,
		cfg.S3Bucket)
	if err != nil {
		return nil, err
	}

	fileS3Repo := s3Repo.NewFileS3Repository(s3Client)
	fileS3UseCase := usecase.NewFileS3UseCase(fileS3Repo)

	redisClient, err := redis.NewRedisClient(ctx, redis.Config{
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

	// Инициализируем сервис импорта Excel
	excelImportUseCase := usecase.NewExcelImportUseCase(
		courseRepo,
		lessonRepo,
		exerciseRepo,
		questionRepo,
		matchingPairRepo,
		questionOptionsRepo,
	)

	// Инициализируем маршрутизаторы
	v1.NewRouter(handler, courseUseCase, lessonUseCase, exerciseUseCase, questionUseCase, attemptUseCase, cfg.GatewayURL)
	admin.NewRouter(handler, courseUseCase, lessonUseCase, exerciseUseCase, questionUseCase, matchingPairUseCase, questionOptionUseCase, excelImportUseCase, fileS3UseCase)
	return &CourseComposite{
		handler: handler,
	}, nil
}

func (c *CourseComposite) Handler() *gin.Engine {
	return c.handler
}
