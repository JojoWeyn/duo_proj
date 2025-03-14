package v1

import (
	"context"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mime/multipart"
	"net/http"
	"strconv"
)

type UserUseCase interface {
	GetUser(ctx context.Context, uuid uuid.UUID) (*entity.User, error)
	UpdateUser(ctx context.Context, uuid uuid.UUID, user *entity.User) error
	DeleteUser(ctx context.Context, uuid uuid.UUID) error
	GetAllUsers(ctx context.Context, limit, offset int) ([]*entity.User, error)
	UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarFile multipart.File, fileSize int64) (string, error)
}

type userRoutes struct {
	uc UserUseCase
}

func newUserRoutes(handler *gin.RouterGroup, uc UserUseCase) {
	r := &userRoutes{
		uc: uc,
	}

	users := handler.Group("/users")
	{
		users.GET("/:uuid", r.getUser)
		users.PATCH("/me", r.updateUser)
		users.GET("/all", r.getAllUsers)
		users.GET("/me", r.getMe)
		users.POST("/me/avatar", r.updateAvatar)
	}
}

func (r *userRoutes) getMe(c *gin.Context) {
	sub := c.GetHeader("X-User-UUID")

	user, err := r.uc.GetUser(c.Request.Context(), uuid.MustParse(sub))
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

	users, err := r.uc.GetAllUsers(c.Request.Context(), limit, offset)
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

	user, err := r.uc.GetUser(c.Request.Context(), userUUID)
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

	if err := r.uc.UpdateUser(c.Request.Context(), uuid.MustParse(sub), &updateData); err != nil {
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

	avatarURL, err := r.uc.UpdateAvatar(c.Request.Context(), uuid.MustParse(sub), avatarFile, fileSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update avatar on storage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"avatar_url": avatarURL})
}
