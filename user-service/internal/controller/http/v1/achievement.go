package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/JojoWeyn/duo-proj/user-service/internal/controller/http/dto"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// @Summary Получить список достижений
// @Description Возвращает список всех публичных достижений
// @Tags Achievements
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Список достижений"
// @Failure 500 {object} map[string]string "Ошибка при получении или парсинге достижений"
// @Router /achievements [get]
func (r *achievementRoutes) getAllAchievements(c *gin.Context) {
	achievements, err := r.uc.GetAllAchievements(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get achievements"})
		return
	}

	var newAchievements []dto.AchievementsDTO
	for _, achievement := range achievements {
		if achievement.Secret {
			continue
		}

		var condition json.RawMessage
		if err := json.Unmarshal([]byte(achievement.Condition), &condition); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse condition"})
			return
		}

		newAchievements = append(newAchievements, dto.AchievementsDTO{
			ID:          achievement.ID,
			Title:       achievement.Title,
			Description: achievement.Description,
			Condition:   condition,
			Secret:      achievement.Secret,
			CreatedAt:   achievement.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"achievements": newAchievements,
	})
}

// @Summary Get user achievements
// @Description Получить достижения пользователя по UUID
// @Tags achievements
// @Produce json
// @Param uuid path string true "UUID пользователя"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /users/achievements/{uuid} [get]
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
