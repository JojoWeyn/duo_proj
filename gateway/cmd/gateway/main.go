package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"

	v1 "github.com/JojoWeyn/duo-proj/gateway/internal/controller/http/v1"
	"github.com/JojoWeyn/duo-proj/gateway/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

type route struct {
	Method    string
	Path      string
	Service   string
	Protected bool
}

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

	public := router.Group("/v1", middleware.RateLimitMiddleware(refreshLimiter))
	protected := router.Group("/v1", middleware.AuthMiddleware(serviceURLs["identity"]))

	publicRoutes := []route{
		{"POST", "/auth/register", "identity", false},
		{"POST", "/auth/login", "identity", false},
		{"POST", "/auth/refresh", "identity", false},
		{"POST", "/auth/password/reset", "identity", false},
		{"POST", "/auth/verification/code", "identity", false},
		{"POST", "/auth/verification/email", "identity", false},
	}

	protectedRoutes := []route{
		{"POST", "/auth/logout", "identity", false},
		{"GET", "/auth/token/status", "identity", false},
		{"GET", "/auth/me", "identity", false},

		// User
		{"GET", "/users/:uuid", "user", true},
		{"GET", "/users/all", "user", true},
		{"GET", "/users/me", "user", true},
		{"GET", "/users/achievements/:uuid", "user", true},
		{"GET", "/achievements/list", "user", true},
		{"GET", "/users/me/progress", "user", true},
		{"GET", "/users/leaderboard", "user", true},
		{"GET", "/users/me/streak", "user", true},
		{"PATCH", "/users/me", "user", true},
		{"POST", "/users/me/avatar", "user", true},

		// Course
		{"GET", "/course/list", "course", true},
		{"GET", "/course/:uuid/info", "course", true},
		{"GET", "/course/:uuid/content", "course", true},
		{"GET", "/lesson/:uuid/info", "course", true},
		{"GET", "/lesson/:uuid/content", "course", true},
		{"GET", "/exercise/:uuid/info", "course", true},
		{"GET", "/question/:uuid/info", "course", true},
		{"GET", "/exercise/:uuid/question", "course", true},
		{"POST", "/question/:uuid/check", "course", true},
		{"POST", "/attempts/start/:exercise_id", "course", true},
		{"POST", "/attempts/answer", "course", true},
		{"POST", "/attempts/finish", "course", true},

		// Admin
		{"GET", "/admin/users", "user", true},
		{"DELETE", "/admin/users/:uuid", "user", true},

		{"POST", "/admin/achievements/create", "user", true},
		{"GET", "/admin/achievements/:uuid", "user", true},
		{"PATCH", "/admin/achievements/:uuid", "user", true},
		{"DELETE", "/admin/achievements/:uuid", "user", true},
		{"GET", "/admin/achievements/list", "user", true},

		{"GET", "/admin/course/list", "course", true},
		{"POST", "/admin/course/import-excel", "course", true},
		{"GET", "/admin/course/:course_id/lesson", "course", true},
		{"GET", "/admin/lesson/:lesson_id/exercise", "course", true},
		{"GET", "/admin/exercise/:exercise_id/question", "course", true},
		{"GET", "/admin/question/:question_id/question-option", "course", true},
		{"GET", "/admin/question/:question_id/matching-pair", "course", true},

		{"GET", "/admin/course/:course_id", "course", true},
		{"GET", "/admin/lesson/:lesson_id", "course", true},
		{"GET", "/admin/exercise/:exercise_id", "course", true},
		{"GET", "/admin/question/:question_id", "course", true},
		{"GET", "/admin/question-option/:uuid", "course", true},
		{"GET", "/admin/matching-pair/:uuid", "course", true},

		{"POST", "/admin/course", "course", true},
		{"POST", "/admin/lesson", "course", true},
		{"POST", "/admin/exercise", "course", true},
		{"POST", "/admin/question", "course", true},
		{"POST", "/admin/question-option", "course", true},
		{"POST", "/admin/matching-pair", "course", true},

		{"PATCH", "/admin/course/:uuid", "course", true},
		{"PATCH", "/admin/lesson/:uuid", "course", true},
		{"PATCH", "/admin/exercise/:uuid", "course", true},
		{"PATCH", "/admin/question/:uuid", "course", true},
		{"PATCH", "/admin/question-option/:uuid", "course", true},
		{"PATCH", "/admin/matching-pair/:uuid", "course", true},

		{"DELETE", "/admin/course/:uuid", "course", true},
		{"DELETE", "/admin/lesson/:uuid", "course", true},
		{"DELETE", "/admin/exercise/:uuid", "course", true},
		{"DELETE", "/admin/question/:uuid", "course", true},
		{"DELETE", "/admin/question-option/:uuid", "course", true},
		{"DELETE", "/admin/matching-pair/:uuid", "course", true},

		// File Admin
		{"POST", "/admin/file/upload", "course", true},
		{"GET", "/admin/file/list", "course", true},
		{"POST", "/admin/file/add", "course", true},
		{"POST", "/admin/file/unpin", "course", true},
		{"DELETE", "/admin/file/delete", "course", true},
	}

	registerRoutes(public, proxy, publicRoutes, serviceURLs)
	registerRoutes(protected, proxy, protectedRoutes, serviceURLs)

	router.GET("/swagger/:service/*path", func(c *gin.Context) {
		service := c.Param("service")
		path := c.Param("path")

		targetURL, ok := serviceURLs[service]
		if !ok {
			c.JSON(404, gin.H{"error": "Service not found"})
			return
		}

		// Проксирование Swagger-запроса
		proxy.ProxySwagger(targetURL, path)(c)
	})

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

func registerRoutes(group *gin.RouterGroup, proxy *v1.ProxyHandler, routes []route, services map[string]string) {
	for _, r := range routes {
		handler := proxy.ProxyService(services[r.Service], r.Protected)
		switch r.Method {
		case "GET":
			group.GET(r.Path, handler)
		case "POST":
			group.POST(r.Path, handler)
		case "PATCH":
			group.PATCH(r.Path, handler)
		case "DELETE":
			group.DELETE(r.Path, handler)
		}
	}
}
