package admin

import (
	"github.com/JojoWeyn/duo-proj/user-service/internal/controller/http/middleware"
	"github.com/gin-gonic/gin"
)

func NewAdminRouter(handler *gin.Engine, uuc UserUseCase, auc AchievementUseCase) {
	v1 := handler.Group("/v1", middleware.RoleMiddleware("admin"))
	{
		newAdminRoutes(v1, uuc, auc)
	}
}
