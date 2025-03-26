package main

import (
	"github.com/gin-contrib/cors"
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
		log.Println("Warning: .env file not found")
	}

	serviceURLs := map[string]string{
		"identity": getEnv("IDENTITY_SERVICE_URL", "http://localhost:8081"),
		"user":     getEnv("USER_SERVICE_URL", "http://176.109.108.209:8082"),
		"course":   getEnv("COURSE_SERVICE_URL", "http://176.109.108.209:8083"),
	}

	jwtSecret := getEnv("JWT_SIGNING_KEY", "your-signing-key")

	proxy := v1.NewProxyHandler(jwtSecret)

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Authorization"}

	router := gin.Default()
	router.Use(cors.New(config))

	refreshLimiter := rate.NewLimiter(rate.Every(1*time.Second), 1)

	public := router.Group("/v1")
	{
		public.POST("/auth/register", proxy.ProxyService(serviceURLs["identity"], false))
		public.POST("/auth/login", proxy.ProxyService(serviceURLs["identity"], false))
		public.POST("/auth/refresh", middleware.RateLimitMiddleware(refreshLimiter), proxy.ProxyService(serviceURLs["identity"], false))
		public.POST("/auth/password/reset", middleware.RateLimitMiddleware(refreshLimiter), proxy.ProxyService(serviceURLs["identity"], false))
		public.POST("/auth/verification/code", proxy.ProxyService(serviceURLs["identity"], false))
	}

	protected := router.Group("/v1", middleware.AuthMiddleware(serviceURLs["identity"]))
	{
		protected.POST("/auth/logout", proxy.ProxyService(serviceURLs["identity"], false))
		protected.GET("/auth/token/status", proxy.ProxyService(serviceURLs["identity"], false))

		userEndpointsGET := []string{
			"/users/:uuid",
			"/users/all",
			"/users/me",
			"/users/achievements/:uuid",
			"/achievements/list",
			"/users/me/progress",
		}
		for _, endpoint := range userEndpointsGET {
			protected.GET(endpoint, proxy.ProxyService(serviceURLs["user"], true))
		}
		protected.PATCH("/users/me", proxy.ProxyService(serviceURLs["user"], true))
		protected.POST("/users/me/avatar", proxy.ProxyService(serviceURLs["user"], true))

		courseEndpointsGET := []string{
			"/course/list",
			"/course/:uuid/info",
			"/course/:uuid/content",
			"/lesson/:uuid/info",
			"/lesson/:uuid/content",
			"/exercise/:uuid/info",
			"/question/:uuid/info",
			"/exercise/:uuid/question",
		}
		for _, endpoint := range courseEndpointsGET {
			protected.GET(endpoint, proxy.ProxyService(serviceURLs["course"], true))
		}
		protected.POST("/question/:uuid/check", proxy.ProxyService(serviceURLs["course"], true))
		protected.POST("/attempts/start/:exercise_id", proxy.ProxyService(serviceURLs["course"], true))
		protected.POST("/attempts/:session_id/answer", proxy.ProxyService(serviceURLs["course"], true))
		protected.POST("/attempts/:session_id/finish", proxy.ProxyService(serviceURLs["course"], true))
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
