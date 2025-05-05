package admin

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JojoWeyn/duo-proj/user-service/internal/controller/http/dto"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserUseCase interface {
	GetAllUsers(ctx context.Context, limit, offset int) ([]*entity.User, error)
	DeleteUser(ctx context.Context, uuid uuid.UUID) error
}

type AchievementUseCase interface {
	CreateAchievement(ctx context.Context, achievement entity.Achievement) (entity.Achievement, error)
	GetAchievementByID(ctx context.Context, id int) (entity.Achievement, error)
	UpdateAchievement(ctx context.Context, id int, achievement entity.Achievement) (entity.Achievement, error)
	DeleteAchievement(ctx context.Context, id int) error
	GetAllAchievements(ctx context.Context) ([]entity.Achievement, error)
}

type adminRoutes struct {
	userUseCase        UserUseCase
	achievementUseCase AchievementUseCase
}

func newAdminRoutes(handler *gin.RouterGroup, uc UserUseCase, ac AchievementUseCase) {
	r := &adminRoutes{
		userUseCase:        uc,
		achievementUseCase: ac,
	}
	h := handler.Group("/admin")
	{
		h.GET("/users", r.getAllUsers)
		h.DELETE("/users/:id", r.deleteUser)

		h.GET("/achievements/list", r.getAllAchievements)
		h.POST("/achievements/create", r.createAchievement)
		h.GET("/achievements/:id", r.getAchievementByID)
		h.PATCH("/achievements/:id", r.updateAchievement)
		h.DELETE("/achievements/:id", r.deleteAchievement)
	}
}

func (r *adminRoutes) getAllUsers(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	users, err := r.userUseCase.GetAllUsers(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
		return
	}

	userDTOs := make([]dto.UserDTO, 0, len(users))
	for _, u := range users {
		userDTOs = append(userDTOs, dto.ToUserDTO(u))
	}

	c.JSON(http.StatusOK, gin.H{
		"users":  userDTOs,
		"limit":  limit,
		"offset": offset,
	})
}

func (r *adminRoutes) deleteUser(c *gin.Context) {
	id := c.Param("id")
	err := r.userUseCase.DeleteUser(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.Status(http.StatusOK)
}

func (r *adminRoutes) createAchievement(c *gin.Context) {
	var createDTO dto.CreateAchievementDTO

	if err := c.ShouldBindJSON(&createDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}

	var conditionObj interface{}
	if err := json.Unmarshal(createDTO.Condition, &conditionObj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат условия достижения"})
		return
	}

	achievement := entity.Achievement{
		Title:       createDTO.Title,
		Description: createDTO.Description,
		Condition:   string(createDTO.Condition),
		Secret:      createDTO.Secret,
	}

	createdAchievement, err := r.achievementUseCase.CreateAchievement(c.Request.Context(), achievement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать достижение"})
		return
	}

	var condition json.RawMessage
	if err := json.Unmarshal([]byte(createdAchievement.Condition), &condition); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при обработке условия"})
		return
	}

	response := dto.AchievementsDTO{
		ID:          createdAchievement.ID,
		Title:       createdAchievement.Title,
		Description: createdAchievement.Description,
		Condition:   condition,
		Secret:      createdAchievement.Secret,
		CreatedAt:   createdAchievement.CreatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{"achievement": response})
}

func (r *adminRoutes) getAchievementByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	achievement, err := r.achievementUseCase.GetAchievementByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "достижение не найдено"})
		return
	}

	var condition json.RawMessage
	if err := json.Unmarshal([]byte(achievement.Condition), &condition); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при обработке условия"})
		return
	}

	response := dto.AchievementsDTO{
		ID:          achievement.ID,
		Title:       achievement.Title,
		Description: achievement.Description,
		Condition:   condition,
		Secret:      achievement.Secret,
		CreatedAt:   achievement.CreatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"achievement": response})
}

func (r *adminRoutes) updateAchievement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	var updateDTO dto.CreateAchievementDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}

	var conditionObj interface{}
	if err := json.Unmarshal(updateDTO.Condition, &conditionObj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат условия достижения"})
		return
	}

	achievement := entity.Achievement{
		Title:       updateDTO.Title,
		Description: updateDTO.Description,
		Condition:   string(updateDTO.Condition),
		Secret:      updateDTO.Secret,
	}

	updatedAchievement, err := r.achievementUseCase.UpdateAchievement(c.Request.Context(), id, achievement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось обновить достижение"})
		return
	}

	var condition json.RawMessage
	if err := json.Unmarshal([]byte(updatedAchievement.Condition), &condition); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при обработке условия"})
		return
	}

	response := dto.AchievementsDTO{
		ID:          updatedAchievement.ID,
		Title:       updatedAchievement.Title,
		Description: updatedAchievement.Description,
		Condition:   condition,
		Secret:      updatedAchievement.Secret,
		CreatedAt:   updatedAchievement.CreatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"achievement": response})
}

func (r *adminRoutes) deleteAchievement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат ID"})
		return
	}

	err = r.achievementUseCase.DeleteAchievement(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось удалить достижение"})
		return
	}

	c.Status(http.StatusOK)
}

func (r *adminRoutes) getAllAchievements(c *gin.Context) {
	achievements, err := r.achievementUseCase.GetAllAchievements(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get achievements"})
		return
	}

	var newAchievements []dto.AchievementsDTO
	for _, achievement := range achievements {
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
