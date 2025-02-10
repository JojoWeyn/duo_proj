package main

import (
	"log"
	"os"
	"time"

	"github.com/JojoWeyn/duo-proj/user-service/internal/composite"
	"github.com/JojoWeyn/duo-proj/user-service/internal/domain/entity"
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

	if err := db.AutoMigrate(&entity.User{}); err != nil {
		log.Fatalf("Failed to migrate db: %s", err.Error())
	}

	cfg := composite.Config{
		KafkaBrokers: getEnv("KAFKA_BROKERS", "176.109.108.209:9092"),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "user_created"),
		GatewayURL:   getEnv("GATEWAY_URL", "176.109.108.209:3211"),
	}

	app, err := composite.NewUserComposite(db, cfg)
	if err != nil {
		log.Fatal("Failed to create application:", err)
	}

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
