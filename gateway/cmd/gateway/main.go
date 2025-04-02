package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
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

	jwtPublicKey, err := loadPublicKey("/app/public.pem")
	if err != nil {
		log.Fatal(err)
	}

	proxy := v1.NewProxyHandler(jwtPublicKey)

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
		public.POST("/auth/verification/email", proxy.ProxyService(serviceURLs["identity"], false))
	}

	protected := router.Group("/v1", middleware.AuthMiddleware(serviceURLs["identity"]))
	{
		protected.POST("/auth/logout", proxy.ProxyService(serviceURLs["identity"], false))
		protected.GET("/auth/token/status", proxy.ProxyService(serviceURLs["identity"], false))
		protected.GET("/auth/me", proxy.ProxyService(serviceURLs["identity"], false))

		userEndpointsGET := []string{
			"/users/:uuid",
			"/users/all",
			"/users/me",
			"/users/achievements/:uuid",
			"/achievements/list",
			"/users/me/progress",
			"/users/leaderboard",
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

		protected.GET("/admin/course/list", proxy.ProxyService(serviceURLs["course"], true))
		protected.GET("/admin/course/:course_id/lesson", proxy.ProxyService(serviceURLs["course"], true))
		protected.GET("/admin/lesson/:lesson_id/exercise", proxy.ProxyService(serviceURLs["course"], true))
		protected.GET("/admin/exercise/:exercise_id/question", proxy.ProxyService(serviceURLs["course"], true))
		protected.GET("/admin/question/:question_id/question-option", proxy.ProxyService(serviceURLs["course"], true))
		protected.GET("/admin/question/:question_id/matching-pair", proxy.ProxyService(serviceURLs["course"], true))

		protected.GET("/admin/course/:course_id", proxy.ProxyService(serviceURLs["course"], true))
		protected.GET("/admin/lesson/:lesson_id", proxy.ProxyService(serviceURLs["course"], true))
		protected.GET("/admin/exercise/:exercise_id", proxy.ProxyService(serviceURLs["course"], true))
		protected.GET("/admin/question/:question_id", proxy.ProxyService(serviceURLs["course"], true))
		protected.GET("/admin/question-option/:uuid", proxy.ProxyService(serviceURLs["course"], true))
		protected.GET("/admin/matching-pair/:uuid", proxy.ProxyService(serviceURLs["course"], true))

		protected.POST("/admin/course", proxy.ProxyService(serviceURLs["course"], true))
		protected.POST("/admin/lesson", proxy.ProxyService(serviceURLs["course"], true))
		protected.POST("/admin/exercise", proxy.ProxyService(serviceURLs["course"], true))
		protected.POST("/admin/question", proxy.ProxyService(serviceURLs["course"], true))
		protected.POST("/admin/question-option", proxy.ProxyService(serviceURLs["course"], true))
		protected.POST("/admin/matching-pair", proxy.ProxyService(serviceURLs["course"], true))

		protected.PATCH("/admin/course/:uuid", proxy.ProxyService(serviceURLs["course"], true))
		protected.PATCH("/admin/lesson/:uuid", proxy.ProxyService(serviceURLs["course"], true))
		protected.PATCH("/admin/exercise/:uuid", proxy.ProxyService(serviceURLs["course"], true))
		protected.PATCH("/admin/question/:uuid", proxy.ProxyService(serviceURLs["course"], true))
		protected.PATCH("/admin/question-option/:uuid", proxy.ProxyService(serviceURLs["course"], true))
		protected.PATCH("/admin/matching-pair/:uuid", proxy.ProxyService(serviceURLs["course"], true))

		protected.DELETE("/admin/course/:uuid", proxy.ProxyService(serviceURLs["course"], true))
		protected.DELETE("/admin/lesson/:uuid", proxy.ProxyService(serviceURLs["course"], true))
		protected.DELETE("/admin/exercise/:uuid", proxy.ProxyService(serviceURLs["course"], true))
		protected.DELETE("/admin/question/:uuid", proxy.ProxyService(serviceURLs["course"], true))
		protected.DELETE("/admin/question-option/:uuid", proxy.ProxyService(serviceURLs["course"], true))
		protected.DELETE("/admin/matching-pair/:uuid", proxy.ProxyService(serviceURLs["course"], true))

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

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block with public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}

	return publicKey, nil
}
