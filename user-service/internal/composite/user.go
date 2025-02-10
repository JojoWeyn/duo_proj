package composite

import (
	v1 "github.com/JojoWeyn/duo-proj/user-service/internal/controller/http/v1"

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
	handler     *gin.Engine
	UserUseCase *usecase.UserUseCase
}

func NewUserComposite(db *gorm.DB, cfg Config) (*UserComposite, error) {
	if err := db.AutoMigrate(&entity.User{}, &entity.Rank{}); err != nil {
		return nil, err
	}

	userRepo := postgres.NewUserRepository(db)
	UserUseCase := usecase.NewUserUseCase(userRepo)

	handler := gin.Default()
	v1.NewRouter(handler, UserUseCase, cfg.GatewayURL)

	return &UserComposite{
		handler:     handler,
		UserUseCase: UserUseCase,
	}, nil
}

func (c *UserComposite) Handler() *gin.Engine {
	return c.handler
}
