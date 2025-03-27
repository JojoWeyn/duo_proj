package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
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

	privateKeySigned, err := loadPrivateKey("private.pem")
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	publicKeySigned, err := loadPublicKey("public.pem")
	if err != nil {
		log.Fatalf("Failed to load public key: %v", err)
	}

	privateKeyRef, err := loadPrivateKey("privateRef.pem")
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	publicKeyRef, err := loadPublicKey("publicRef.pem")
	if err != nil {
		log.Fatalf("Failed to load public key: %v", err)
	}

	identityComposite, err := composite.NewIdentityComposite(db, composite.Config{
		AccessTokenTTL:  time.Duration(getEnvAsInt("ACCESS_TOKEN_TTL", 15)) * time.Minute,
		RefreshTokenTTL: time.Duration(getEnvAsInt("REFRESH_TOKEN_TTL", 24)) * time.Minute,
		SigningKey:      privateKeySigned,
		SigningPublic:   publicKeySigned,
		RefreshKey:      privateKeyRef,
		RefreshPublic:   publicKeyRef,
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

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block with private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
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
