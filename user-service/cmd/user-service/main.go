package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/JojoWeyn/duo-proj/user-service/internal/controller/kafka"

	"github.com/JojoWeyn/duo-proj/user-service/internal/composite"
	"github.com/JojoWeyn/duo-proj/user-service/pkg/client/postgresql"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: loading .env file")
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

	cfg := composite.Config{
		KafkaBrokers: getEnv("KAFKA_BROKERS", "kafka:29092"),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "user_create"),
		GatewayURL:   getEnv("GATEWAY_URL", "176.109.108.209:3211"),
		Secret:       getEnv("JWT_SIGNING_KEY", "your-signing-key"),
		S3Endpoint:   getEnv("S3_ENDPOINT", "minio:9000"),
		S3AccessKey:  getEnv("S3_ACCESS_KEY", "minio"),
		S3SecretKey:  getEnv("S3_SECRET_KEY", "minio123"),
		S3Bucket:     getEnv("S3_BUCKET", "user-avatar"),
	}

	app, err := composite.NewUserComposite(db, cfg)

	if err != nil {
		log.Fatal("Failed to create application:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := kafka.NewSaramaConsumerGroup(app.UserUseCase, app.AchievementUseCase, []string{cfg.KafkaBrokers}, cfg.KafkaTopic, "user-service-group")
	go consumer.Start(ctx)

	port := getEnv("USER_PORT", "8082")
	if err := app.Handler().Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
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
