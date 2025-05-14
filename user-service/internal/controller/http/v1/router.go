package v1

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(handler *gin.Engine, userUseCase UserUseCase, achievementUseCase AchievementUseCase, progressUseCase ProgressUseCase, gatewayUrl string) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{gatewayUrl}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true

	handler.Use(cors.New(config))

	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	v1 := handler.Group("/v1")
	{
		newUserRoutes(v1, userUseCase, progressUseCase)
		newAchievementRoutes(v1, achievementUseCase)
	}
}
