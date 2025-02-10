package main

import (
	"log"
	"os"
	"time"

	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"

	"github.com/JojoWeyn/duo-proj/identity-service/internal/composite"
	"github.com/JojoWeyn/duo-proj/identity-service/pkg/client/postgresql"
	"github.com/joho/godotenv"
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

	if err := db.AutoMigrate(&entity.Identity{}, entity.BlacklistedToken{}); err != nil {
		log.Fatalf("Failed to migrate db: %s", err.Error())
	}

	identityComposite, err := composite.NewIdentityComposite(db, composite.Config{
		AccessTokenTTL:  time.Duration(getEnvAsInt("ACCESS_TOKEN_TTL", 15)) * time.Minute,
		RefreshTokenTTL: time.Duration(getEnvAsInt("REFRESH_TOKEN_TTL", 24)) * time.Hour,
		SigningKey:      getEnv("JWT_SIGNING_KEY", "your-signing-key"),
		RefreshKey:      getEnv("JWT_REFRESH_KEY", "your-refresh-key"),
		GatewayURL:      getEnv("GATEWAY_URL", "176.109.108.209:3211"),
		KafkaBrokers:    getEnv("KAFKA_BROKERS", "kafka:29092"),
	})
	if err != nil {
		log.Fatalf("Failed to initialize composite: %s", err.Error())
	}

	port := getEnv("IDENTITY_PORT", "8081")
	log.Printf("Starting server on port %s", port)
	if err := identityComposite.Handler().Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
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
