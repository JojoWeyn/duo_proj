package composite

import (
	"crypto/rsa"
	"github.com/JojoWeyn/duo-proj/identity-service/pkg/client/smtp"
	"time"

	v1 "github.com/JojoWeyn/duo-proj/identity-service/internal/controller/http/v1"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/controller/kafka"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/usecase"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/repository/db/postgres"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IdentityComposite struct {
	handler *gin.Engine
}

type Config struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	SigningKey      *rsa.PrivateKey
	RefreshKey      *rsa.PrivateKey
	SigningPublic   *rsa.PublicKey
	RefreshPublic   *rsa.PublicKey
	GatewayURL      string
	KafkaBrokers    string
	SmtpServer      string
	SmtpPort        string
	SmtpSender      string
	SmtpPassword    string
}

func NewIdentityComposite(db *gorm.DB, cfg Config) (*IdentityComposite, error) {
	if err := db.AutoMigrate(&entity.Identity{}, &entity.BlacklistedToken{}); err != nil {
		return nil, err
	}

	identityRepo := postgres.NewIdentityRepository(db)
	tokenRepo := postgres.NewTokenRepository(db)

	tokenService := service.NewTokenService(
		cfg.SigningKey,
		cfg.SigningPublic,
		cfg.RefreshKey,
		cfg.RefreshPublic,
		cfg.AccessTokenTTL,
		cfg.RefreshTokenTTL,
		tokenRepo,
	)

	smtpClient := smtp.NewSMTPClient(smtp.SMTPConfig{
		Server:   cfg.SmtpServer,
		Port:     cfg.SmtpPort,
		Sender:   cfg.SmtpSender,
		Password: cfg.SmtpPassword,
	})

	verificationService := service.NewVerificationService(smtpClient)

	producer, err := kafka.NewProducer(cfg.KafkaBrokers, "user_create")
	if err != nil {
		return nil, err
	}

	identityUseCase := usecase.NewIdentityUseCase(
		identityRepo,
		tokenService,
		tokenRepo,
		producer,
	)

	handler := gin.Default()

	v1.NewRouter(handler, verificationService, identityUseCase, cfg.GatewayURL)

	return &IdentityComposite{
		handler: handler,
	}, nil
}

func (c *IdentityComposite) Handler() *gin.Engine {
	return c.handler
}
