package v1

import (
	"context"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type UserUseCase interface {
	GetUser(ctx context.Context, uuid uuid.UUID) (*entity.User, error)
	UpdateUser(ctx context.Context, uuid uuid.UUID, user *entity.User) error
	DeleteUser(ctx context.Context, uuid uuid.UUID) error
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
}

type TokenService interface {
	Validate(token string) (string, error)
}

type userRoutes struct {
	uc UserUseCase
	ts TokenService
}

func newUserRoutes(handler *gin.RouterGroup, uc UserUseCase, ts TokenService) {
	r := &userRoutes{
		uc: uc,
		ts: ts,
	}

	users := handler.Group("/users")
	{
		users.GET("/:uuid", r.getUser)
		users.PUT("/update", r.updateUser)
		users.DELETE("/", r.deleteUser)
		users.GET("/", r.getAllUsers)
	}
}

func (r *userRoutes) getAllUsers(c *gin.Context) {
	users, err := r.uc.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
		return
	}

	c.JSON(http.StatusOK, users)
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
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no token provided"})
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")

	sub, err := r.ts.Validate(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
	}

	userUUID, err := uuid.Parse(sub)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": userUUID})
		return
	}

	var updateData entity.User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.uc.UpdateUser(c.Request.Context(), userUUID, &updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.Status(http.StatusOK)
}

func (r *userRoutes) deleteUser(c *gin.Context) {
	userUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
		return
	}

	if err := r.uc.DeleteUser(c.Request.Context(), userUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.Status(http.StatusOK)
}
