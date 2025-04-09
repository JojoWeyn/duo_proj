package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetHeader("X-User-Role")

		if userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
