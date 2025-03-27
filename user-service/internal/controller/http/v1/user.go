package v1

import (
	"context"
	"github.com/JojoWeyn/duo-proj/user-service/internal/controller/http/dto"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mime/multipart"
	"net/http"
	"strconv"
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
	}
}

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

	c.JSON(http.StatusOK, gin.H{
		"leaderboard": leaderboard,
		"limit":       limit,
		"offset":      offset,
	})
}

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

	c.JSON(http.StatusOK, gin.H{
		"exercises": exerciseProgress,
		"lessons":   lessonProgress,
		"courses":   courseProgress,
	})
}

func (r *userRoutes) getMe(c *gin.Context) {
	sub := c.GetHeader("X-User-UUID")

	user, err := r.userUseCase.GetUser(c.Request.Context(), uuid.MustParse(sub))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
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

	c.JSON(http.StatusOK, gin.H{
		"users":  users,
		"limit":  limit,
		"offset": offset,
	})
}
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

	c.JSON(http.StatusOK, user)
}

func (r *userRoutes) updateUser(c *gin.Context) {
	sub := c.GetHeader("X-User-UUID")

	var updateData entity.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.userUseCase.UpdateUser(c.Request.Context(), uuid.MustParse(sub), &updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

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

	c.JSON(http.StatusOK, gin.H{"avatar_url": avatarURL})
}
