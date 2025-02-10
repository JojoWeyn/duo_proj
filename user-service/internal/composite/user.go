package composite

import (
	v1 "github.com/JojoWeyn/duo-proj/user-service/internal/controller/http/v1"
	"github.com/JojoWeyn/duo-proj/user-service/internal/controller/kafka"
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
}

type UserComposite struct {
	handler *gin.Engine
}

func NewUserComposite(db *gorm.DB, cfg Config) (*UserComposite, error) {
	if err := db.AutoMigrate(&entity.User{}, &entity.Rank{}); err != nil {
		return nil, err
	}

	userRepo := postgres.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)

	consumer := kafka.NewSaramaConsumer([]string{cfg.KafkaBrokers}, cfg.KafkaTopic)
	go consumer.Start()

	handler := gin.Default()
	v1.NewRouter(handler, userUseCase, cfg.GatewayURL)

	return &UserComposite{
		handler: handler,
	}, nil
}

func (c *UserComposite) Handler() *gin.Engine {
	return c.handler
}
