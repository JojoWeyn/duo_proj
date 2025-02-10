package v1

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, uc IdentityUsecase, tr TokenRepository, gatewayUrl string) {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{gatewayUrl}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true

	handler.Use(cors.New(config))

	v1 := handler.Group("/v1")
	{
		newIdentityRoutes(v1, uc, tr)
	}
}
