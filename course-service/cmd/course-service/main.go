package main

import (
	"github.com/JojoWeyn/duo-proj/course-service/internal/composite"
	"github.com/JojoWeyn/duo-proj/course-service/pkg/client/postgresql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	db, err := postgresql.NewPostgresDB(postgresql.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "admin"),
		DBName:   getEnv("DB_NAME", "postgres"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		SslMode:  getEnv("DB_SSL_MODE", "disable"),
	})
	if err != nil {
		log.Fatalf("Failed to initialize db: %s", err.Error())
	}

	courseComposite, err := composite.NewCourseComposite(db, composite.Config{
		GatewayURL:   getEnv("GATEWAY_URL", "176.109.108.209:3211"),
		RedisURL:     getEnv("REDIS_URL", "redis:6379"),
		RedisDB:      getEnvAsInt("REDIS_DB", 0),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost"),
		S3Endpoint:   getEnv("S3_ENDPOINT", "minio:9000"),
		S3AccessKey:  getEnv("S3_ACCESS_KEY", "minio"),
		S3SecretKey:  getEnv("S3_SECRET_KEY", "minio123"),
		S3Bucket:     getEnv("S3_BUCKET", "duo-bucket"),
	})
	if err != nil {
		panic(err)
	}

	port := getEnv("COURSE_PORT", "8083")

	if err := courseComposite.Handler().Run(":" + port); err != nil {
		panic(err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := time.ParseDuration(value); err == nil {
			return int(intValue.Minutes())
		}
	}
	return defaultValue
}
