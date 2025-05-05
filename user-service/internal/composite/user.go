package composite

import (
	"context"
	"log"

	v1 "github.com/JojoWeyn/duo-proj/user-service/internal/controller/http/v1"
	"github.com/JojoWeyn/duo-proj/user-service/internal/controller/http/v1/admin"
	"github.com/JojoWeyn/duo-proj/user-service/internal/controller/kafka"
	"github.com/JojoWeyn/duo-proj/user-service/internal/repository/cache"
	"github.com/JojoWeyn/duo-proj/user-service/pkg/client/redis"
	"github.com/JojoWeyn/duo-proj/user-service/pkg/client/s3"

	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/usecase"
	"github.com/JojoWeyn/duo-proj/user-service/internal/repository/db/postgres"
	s3Repo "github.com/JojoWeyn/duo-proj/user-service/internal/repository/db/s3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Config struct {
	KafkaBrokers string
	KafkaTopic   string
	GatewayURL   string
	Secret       string
	S3Endpoint   string
	S3AccessKey  string
	S3SecretKey  string
	S3Bucket     string
	RedisURL     string
	RedisDB      int
}

type UserComposite struct {
	handler            *gin.Engine
	UserUseCase        *usecase.UserUseCase
	AchievementUseCase *usecase.AchievementUseCase
	ProgressUseCase    *usecase.ProgressUseCase
}

func NewUserComposite(ctx context.Context, db *gorm.DB, cfg Config) (*UserComposite, error) {
	if err := db.AutoMigrate(&entity.User{}, &entity.Rank{}, &entity.Progress{}); err != nil {
		return nil, err
	}

	var count int64
	db.Model(&entity.Rank{}).Count(&count)
	if count == 0 {
		ranks := []entity.Rank{
			{ID: 1, Name: "Новичек"},
		}

		if err := db.Create(&ranks).Error; err != nil {
			log.Printf("Failed to add default ranks: %v", err)
		}
		log.Println("Default ranks added")
	}

	if err := db.AutoMigrate(&entity.Achievement{}, &entity.UserAchievementProgress{}); err != nil {
		return nil, err
	}
	db.Model(&entity.Achievement{}).Count(&count)
	if count == 0 {
		achievements := []entity.Achievement{
			{ID: 1, Title: "Обновки", Description: "Обновить данные профиля", Condition: `{"action": "update"}`},
			{ID: 2, Title: "Вхождение", Description: "3 раза войти в аккаунт", Condition: `{"action": "login", "count": 3}`},
		}

		if err := db.Create(&achievements).Error; err != nil {
			log.Printf("Failed to add default achievements: %v", err)
		}
		log.Println("Default achievements added")
	}
	s3Client, err := s3.NewS3Client(cfg.S3Endpoint, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3Bucket)
	if err != nil {
		log.Printf("Failed to create S3 client: %v", err)
	}

	redisClient, err := redis.NewRedisClient(ctx, redis.Config{
		Addr: cfg.RedisURL,
		DB:   cfg.RedisDB,
	})
	if err != nil {
		return nil, err
	}

	userRepo := postgres.NewUserRepository(db)
	userS3Repo := s3Repo.NewUserS3Repository(s3Client)
	achievementRepo := postgres.NewAchievementRepository(db)
	progressRepo := postgres.NewProgressRepository(db)
	cacher := cache.NewRedisCache(redisClient)
	producer, err := kafka.NewProducer(cfg.KafkaBrokers, "user_create")
	if err != nil {
		return nil, err
	}

	UserUseCase := usecase.NewUserUseCase(userRepo, cacher, userS3Repo, producer)
	AchievementUseCase := usecase.NewAchievementUseCase(achievementRepo)
	progressUseCase := usecase.NewProgressUseCase(progressRepo)

	handler := gin.Default()
	v1.NewRouter(handler, UserUseCase, AchievementUseCase, progressUseCase, cfg.GatewayURL)
	admin.NewAdminRouter(handler, UserUseCase, AchievementUseCase)

	return &UserComposite{
		handler:            handler,
		UserUseCase:        UserUseCase,
		AchievementUseCase: AchievementUseCase,
		ProgressUseCase:    progressUseCase,
	}, nil
}

func (c *UserComposite) Handler() *gin.Engine {
	return c.handler
}
