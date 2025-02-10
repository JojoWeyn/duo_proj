package main

import (
	"log"
	"os"
	"time"

	v1 "github.com/JojoWeyn/duo-proj/gateway/internal/controller/http/v1"
	"github.com/JojoWeyn/duo-proj/gateway/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	identityServiceURL := getEnv("IDENTITY_SERVICE_URL", "http://localhost:8080")
	proxyHandler, err := v1.NewProxyHandler(identityServiceURL)
	if err != nil {
		log.Fatalf("Failed to initialize proxy handler: %s", err.Error())
	}

	router := gin.Default()

	refreshLimiter := rate.NewLimiter(rate.Every(1*time.Second), 1)

	public := router.Group("/v1/auth")
	{
		public.POST("/register", proxyHandler.ProxyIdentityService())
		public.POST("/login", proxyHandler.ProxyIdentityService())
		public.POST("/refresh", middleware.RateLimitMiddleware(refreshLimiter), proxyHandler.ProxyIdentityService())
	}

	protected := router.Group("/v1/auth")
	protected.Use(middleware.AuthMiddleware(identityServiceURL))
	{
		protected.POST("/logout", proxyHandler.ProxyIdentityService())
		protected.GET("/token/status", proxyHandler.ProxyIdentityService())
	}

	port := getEnv("PORT", "3211")

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
