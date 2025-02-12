package composite

import (
	v1 "github.com/JojoWeyn/duo-proj/user-service/internal/controller/http/v1"
	"github.com/JojoWeyn/duo-proj/user-service/internal/service"
	"log"

	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/usecase"
	"github.com/JojoWeyn/duo-proj/user-service/internal/repository/db/postgres"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Config struct {
	KafkaBrokers string
	KafkaTopic   string
	GatewayURL   string
	Secret       string
}

type UserComposite struct {
	handler     *gin.Engine
	UserUseCase *usecase.UserUseCase
}

func NewUserComposite(db *gorm.DB, cfg Config) (*UserComposite, error) {
	if err := db.AutoMigrate(&entity.User{}, &entity.Rank{}); err != nil {
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

	userRepo := postgres.NewUserRepository(db)
	UserUseCase := usecase.NewUserUseCase(userRepo)
	TokenService := service.NewTokenService(cfg.Secret)

	handler := gin.Default()
	v1.NewRouter(handler, UserUseCase, TokenService, cfg.GatewayURL)

	return &UserComposite{
		handler:     handler,
		UserUseCase: UserUseCase,
	}, nil
}

func (c *UserComposite) Handler() *gin.Engine {
	return c.handler
}
