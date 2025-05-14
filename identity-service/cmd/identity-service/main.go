package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"

	_ "github.com/JojoWeyn/duo-proj/identity-service/docs"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/composite"
	"github.com/JojoWeyn/duo-proj/identity-service/pkg/client/postgresql"
	"github.com/joho/godotenv"
)

// @title Identity Service API
// @version 1.0
// @description Сервис для управления пользователями: регистрация, аутентификация, токены и подтверждение email.
// @host localhost:8081
// @BasePath /v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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

	privateKeySigned, err := loadPrivateKey("/app/private.pem")
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	publicKeySigned, err := loadPublicKey("/app/public.pem")
	if err != nil {
		log.Fatalf("Failed to load public key: %v", err)
	}

	privateKeyRef, err := loadPrivateKey("/app/privateRef.pem")
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	publicKeyRef, err := loadPublicKey("/app/publicRef.pem")
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
		SmtpServer:      getEnv("SMTP_SERVER", ""),
		SmtpPort:        getEnv("SMTP_PORT", "443"),
		SmtpSender:      getEnv("SMTP_SENDER", ""),
		SmtpPassword:    getEnv("SMTP_PASSWORD", ""),
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
	privData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %v", err)
	}

	block, _ := pem.Decode(privData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block with private key")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %v", err)
		}
	}

	return privKey.(*rsa.PrivateKey), nil
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
