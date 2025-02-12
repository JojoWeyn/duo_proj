package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type TokenStatus struct {
	IsBlacklisted string `json:"is_blacklisted"`
}

func AuthMiddleware(identityServiceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
			c.Abort()
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/auth/token/status", identityServiceURL), nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			c.Abort()
			return
		}

		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate token"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		var tokenStatus TokenStatus
		if err := json.NewDecoder(resp.Body).Decode(&tokenStatus); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse token status"})
			c.Abort()
			return
		}

		c.Next()
	}
}
