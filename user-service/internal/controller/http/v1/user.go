package v1

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/JojoWeyn/duo-proj/user-service/internal/controller/http/dto"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserUseCase interface {
	GetLeaderboard(ctx context.Context, limit, offset int) ([]entity.Leaderboard, error)
	GetUser(ctx context.Context, uuid uuid.UUID) (*entity.User, error)
	UpdateUser(ctx context.Context, uuid uuid.UUID, user *entity.User) error
	DeleteUser(ctx context.Context, uuid uuid.UUID) error
	GetAllUsers(ctx context.Context, limit, offset int) ([]*entity.User, error)
	UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarFile multipart.File, fileSize int64) (string, error)
}

type ProgressUseCase interface {
	GetProgress(ctx context.Context, userID uuid.UUID) ([]*entity.Progress, error)
	GetStreak(ctx context.Context, userID uuid.UUID) (int, error)
}

type userRoutes struct {
	userUseCase     UserUseCase
	progressUseCase ProgressUseCase
}

func newUserRoutes(handler *gin.RouterGroup, uc UserUseCase, puc ProgressUseCase) {
	r := &userRoutes{
		userUseCase:     uc,
		progressUseCase: puc,
	}

	users := handler.Group("/users")
	{
		users.GET("/:uuid", r.getUser)
		users.PATCH("/me", r.updateUser)
		users.GET("/all", r.getAllUsers)
		users.GET("/me", r.getMe)
		users.POST("/me/avatar", r.updateAvatar)
		users.GET("/me/progress", r.getProgress)
		users.GET("/leaderboard", r.getLeaderboard)
		users.GET("/me/streak", r.getStreak)
	}
}

// @Summary Получить streak пользователя
// @Description Возвращает количество последовательных дней активности пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} dto.StreakResponseDTO "Серия активности пользователя (streak)"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/streak [get]
func (r *userRoutes) getStreak(c *gin.Context) {
	userUUID, err := uuid.Parse(c.GetHeader("X-User-UUID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	streak, err := r.progressUseCase.GetStreak(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "streak not found"})
		return
	}

	c.JSON(http.StatusOK, dto.StreakResponseDTO{Days: streak})
}

// @Summary Get leaderboard
// @Description Получить таблицу лидеров
// @Tags users
// @Produce json
// @Param limit query int false "Лимит"
// @Param offset query int false "Смещение"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/leaderboard [get]
func (r *userRoutes) getLeaderboard(c *gin.Context) {
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

	leaderboard, err := r.userUseCase.GetLeaderboard(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get leaderboard"})
		return
	}

	dtoList := make([]dto.LeaderboardDTO, 0, len(leaderboard))
	for _, l := range leaderboard {
		dtoList = append(dtoList, dto.ToLeaderboardDTO(l))
	}

	c.JSON(http.StatusOK, gin.H{
		"leaderboard": dtoList,
		"limit":       limit,
		"offset":      offset,
	})
}

// @Summary Получить прогресс пользователя
// @Description Возвращает прогресс пользователя, сгруппированный по упражнениям, урокам и курсам
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} dto.ProgressResponseDTO "Прогресс пользователя"
// @Failure 404 {object} map[string]string
// @Router /users/progress [get]
func (r *userRoutes) getProgress(c *gin.Context) {
	sub := c.GetHeader("X-User-UUID")

	progresses, err := r.progressUseCase.GetProgress(c.Request.Context(), uuid.MustParse(sub))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "progress not found"})
		return
	}

	exerciseProgress := []dto.ExerciseProgressDTO{}
	lessonProgress := []dto.LessonProgressDTO{}
	courseProgress := []dto.CourseProgressDTO{}

	for _, p := range progresses {
		switch p.EntityType {
		case "exercise":
			exerciseProgress = append(exerciseProgress, dto.ExerciseProgressDTO{
				UUID:         p.UUID,
				ExerciseUUID: p.EntityUUID,
				TotalPoints:  p.Points,
				CompletedAt:  p.CompletedAt,
			})
		case "lesson":
			lessonProgress = append(lessonProgress, dto.LessonProgressDTO{
				UUID:        p.UUID,
				LessonUUID:  p.EntityUUID,
				TotalPoints: p.Points,
				CompletedAt: p.CompletedAt,
			})
		case "course":
			courseProgress = append(courseProgress, dto.CourseProgressDTO{
				UUID:        p.UUID,
				CourseUUID:  p.EntityUUID,
				TotalPoints: p.Points,
				CompletedAt: p.CompletedAt,
			})
		}
	}

	c.JSON(http.StatusOK, dto.ProgressResponseDTO{
		Exercises: exerciseProgress,
		Lessons:   lessonProgress,
		Courses:   courseProgress,
	})
}

// @Description Получить данные текущего пользователя
// @Tags users
// @Produce json
// @Success 200 {object} dto.UserDTO
// @Failure 404 {object} map[string]string
// @Router /users/me [get]
func (r *userRoutes) getMe(c *gin.Context) {
	sub := c.GetHeader("X-User-UUID")

	user, err := r.userUseCase.GetUser(c.Request.Context(), uuid.MustParse(sub))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	userDTO := dto.ToUserDTO(user)
	c.JSON(http.StatusOK, userDTO)
}

// @Summary Получить список пользователей
// @Description Возвращает список пользователей с поддержкой limit и offset
// @Tags Users
// @Accept json
// @Produce json
// @Param limit query int false "Максимальное количество пользователей (по умолчанию 50)"
// @Param offset query int false "Смещение для пагинации (по умолчанию 0)"
//
//	@Success 200 {object} map[string]string
//
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/all [get]
func (r *userRoutes) getAllUsers(c *gin.Context) {
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

// getUser обрабатывает GET-запрос для получения информации о пользователе по UUID.
//
// @Summary Получить пользователя
// @Description Возвращает информацию о пользователе по его UUID
// @Tags Users
// @Accept json
// @Produce json
// @Param uuid path string true "UUID пользователя"
// @Success 200 {object} dto.UserDTO
// @Failure 400 {object} map[string]string "Неверный формат UUID"
// @Failure 404 {object} map[string]string "Пользователь не найден"
// @Router /users/{uuid} [get]
func (r *userRoutes) getUser(c *gin.Context) {
	userUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	user, err := r.userUseCase.GetUser(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, dto.ToUserDTO(user))
}

// @Summary Update current user
// @Description Обновить данные текущего пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.UserUpdateDTO true "Данные для обновления"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me [patch]
func (r *userRoutes) updateUser(c *gin.Context) {
	sub := c.GetHeader("X-User-UUID")

	var updateDTO dto.UserUpdateDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entityUpdate := &entity.User{}
	if updateDTO.Name != nil {
		entityUpdate.Name = *updateDTO.Name
	}
	if updateDTO.SecondName != nil {
		entityUpdate.SecondName = *updateDTO.SecondName
	}
	if updateDTO.LastName != nil {
		entityUpdate.LastName = *updateDTO.LastName
	}
	if updateDTO.Login != nil {
		entityUpdate.Login = *updateDTO.Login
	}

	if err := r.userUseCase.UpdateUser(c.Request.Context(), uuid.MustParse(sub), entityUpdate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Update avatar
// @Description Загрузить и обновить аватар текущего пользователя
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл аватара"
// @Success 200 {object} dto.UserAvatarResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me/avatar [post]
func (r *userRoutes) updateAvatar(c *gin.Context) {
	sub := c.GetHeader("X-User-UUID")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileSize := file.Size
	avatarFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open avatar file"})
		return
	}
	defer avatarFile.Close()

	avatarURL, err := r.userUseCase.UpdateAvatar(c.Request.Context(), uuid.MustParse(sub), avatarFile, fileSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update avatar on storage"})
		return
	}

	c.JSON(http.StatusOK, dto.UserAvatarResponseDTO{AvatarURL: avatarURL})
}

func extractUserUUID(c *gin.Context) (uuid.UUID, error) {
	sub := c.GetHeader("X-User-UUID")
	if sub == "" {
		return uuid.Nil, errors.New("X-User-UUID header missing")
	}
	return uuid.Parse(sub)
}
