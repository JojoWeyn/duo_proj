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

	identityServiceURL := getEnv("IDENTITY_SERVICE_URL", "http://localhost:8081")
	userServiceURL := getEnv("USER_SERVICE_URL", "http://176.109.108.209:8082")

	proxy, err := v1.NewProxyHandler(identityServiceURL, userServiceURL)
	if err != nil {
		log.Fatalf("Failed to initialize identity proxy handler: %s", err.Error())
	}

	router := gin.Default()

	refreshLimiter := rate.NewLimiter(rate.Every(1*time.Second), 1)

	public := router.Group("/v1")
	{
		public.POST("/auth/register", proxy.ProxyIdentityService())
		public.POST("/auth/login", proxy.ProxyIdentityService())
		public.POST("/auth/refresh", middleware.RateLimitMiddleware(refreshLimiter), proxy.ProxyIdentityService())
		public.POST("/auth/password/reset", proxy.ProxyIdentityService())
		public.POST("/auth/verification/code", proxy.ProxyIdentityService())
	}

	protected := router.Group("/v1")
	protected.Use(middleware.AuthMiddleware(identityServiceURL))
	{
		protected.POST("/auth/logout", proxy.ProxyIdentityService())
		protected.GET("/auth/token/status", proxy.ProxyIdentityService())

		protected.GET("/users/:uuid", proxy.ProxyUserService())
		protected.GET("/users/all", proxy.ProxyUserService())
		protected.GET("/users/me", proxy.ProxyUserService())
		protected.PATCH("/users/me", proxy.ProxyUserService())
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
