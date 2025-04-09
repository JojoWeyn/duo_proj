package admin

import (
	"context"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type UserUseCase interface {
	GetAllUsers(ctx context.Context, limit, offset int) ([]*entity.User, error)
	DeleteUser(ctx context.Context, uuid uuid.UUID) error
}
type adminRoutes struct {
	userUseCase UserUseCase
}

func newAdminRoutes(handler *gin.RouterGroup, uc UserUseCase) {
	r := &adminRoutes{
		userUseCase: uc,
	}
	h := handler.Group("/admin")
	{
		h.GET("/users", r.getAllUsers)
		h.DELETE("/users/:id", r.deleteUser)
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

	c.JSON(http.StatusOK, gin.H{
		"users":  users,
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
