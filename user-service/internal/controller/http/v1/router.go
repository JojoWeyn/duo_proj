package v1

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, userUseCase UserUseCase, achievementUseCase AchievementUseCase, ts TokenService, gatewayUrl string) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{gatewayUrl}
	config.AllowMethods = []string{"GET", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true

	handler.Use(cors.New(config))

	v1 := handler.Group("/v1")
	{
		newUserRoutes(v1, userUseCase, ts)
		newAchievementRoutes(v1, achievementUseCase)
	}
}
