package v1

import (
	"context"
	"github.com/JojoWeyn/duo-proj/user-service/internal/controller/http/dto"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type AchievementUseCase interface {
	GetAllAchievements(ctx context.Context) ([]entity.Achievement, error)
	GetUserAchievements(ctx context.Context, userID uuid.UUID) ([]dto.UserAchievementsDTO, error)
}

type achievementRoutes struct {
	uc AchievementUseCase
}

func newAchievementRoutes(handler *gin.RouterGroup, uc AchievementUseCase) {
	r := &achievementRoutes{
		uc: uc,
	}

	achievements := handler.Group("/achievements")
	{
		achievements.GET("/list", r.getAllAchievements)
	}

	ac := handler.Group("/users/achievements")
	{
		ac.GET("/:uuid", r.getUserAchievements)
	}
}

func (r *achievementRoutes) getAllAchievements(c *gin.Context) {
	achievements, err := r.uc.GetAllAchievements(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get achievements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"achievements": achievements,
	})
}

func (r *achievementRoutes) getUserAchievements(c *gin.Context) {
	userUUID := c.Param("uuid")

	achievements, err := r.uc.GetUserAchievements(c.Request.Context(), uuid.MustParse(userUUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get achievements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"achievements": achievements,
	})
}
