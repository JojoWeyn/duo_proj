package usecase

import (
	"context"
	"errors"
	"log"

	"github.com/JojoWeyn/duo-proj/identity-service/internal/controller/kafka"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type IdentityRepository interface {
	Create(ctx context.Context, identity *entity.Identity) error
	FindByUUID(ctx context.Context, userID string) (*entity.Identity, error)
	FindByLogin(ctx context.Context, login string) (*entity.Identity, error)
	FindByEmail(ctx context.Context, email string) (*entity.Identity, error)
	Update(ctx context.Context, identity *entity.Identity) error
}

type TokenRepository interface {
	BlacklistToken(ctx context.Context, token *entity.BlacklistedToken) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
	CleanupExpired(ctx context.Context) error
}

type TokenService interface {
	GenerateTokenPair(userID string) (accessToken string, refreshToken string, err error)
	ValidateToken(token string, isRefreshToken bool) (userID string, err error)
	BlacklistToken(ctx context.Context, token string) error
}

type IdentityUseCase struct {
	identityRepo IdentityRepository
	tokenService TokenService
	tokenRepo    TokenRepository
	producer     *kafka.Producer
}

func NewIdentityUseCase(identityRepo IdentityRepository, tokenService TokenService, tokenRepo TokenRepository, producer *kafka.Producer) *IdentityUseCase {
	return &IdentityUseCase{
		identityRepo: identityRepo,
		tokenService: tokenService,
		tokenRepo:    tokenRepo,
		producer:     producer,
	}
}

func (uc *IdentityUseCase) Logout(ctx context.Context, token string) error {
	return uc.tokenService.BlacklistToken(ctx, token)
}

func (uc *IdentityUseCase) ValidateToken(ctx context.Context, token string, isRefreshToken bool) (string, error) {
	isBlacklisted, err := uc.tokenRepo.IsBlacklisted(ctx, token)
	if err != nil {
		return "", err
	}
	if isBlacklisted {
		return "", errors.New("token is blacklisted")
	}

	return uc.tokenService.ValidateToken(token, isRefreshToken)
}

func (uc *IdentityUseCase) RefreshToken(ctx context.Context, refreshToken string) (*Tokens, error) {
	userID, err := uc.tokenService.ValidateToken(refreshToken, true)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if err := uc.tokenService.BlacklistToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	accessToken, newRefreshToken, err := uc.tokenService.GenerateTokenPair(userID)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (uc *IdentityUseCase) Register(ctx context.Context, login, email, password string) error {
	if _, err := uc.identityRepo.FindByEmail(ctx, email); err == nil {
		return errors.New("email already exists")
	}

	if _, err := uc.identityRepo.FindByLogin(ctx, login); err == nil {
		return errors.New("login already exists")
	}

	err := entity.ValidatePassword(password)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	identity, err := entity.NewIdentity(
		login,
		email,
		string(hashedPassword),
	)
	if err != nil {
		return err
	}

	if err := uc.producer.SendUserCreated(identity.UserUUID.String()); err != nil {

		log.Printf("Failed to send user created event: %v", err)
	}

	return uc.identityRepo.Create(ctx, identity)

}

func (uc *IdentityUseCase) Login(ctx context.Context, login, password string) (*Tokens, error) {
	identity, err := uc.identityRepo.FindByLogin(ctx, login)
	if err != nil {
		return nil, errors.New("invalid login or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(identity.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid login or password")
	}
	accessToken, refreshToken, err := uc.tokenService.GenerateTokenPair(identity.UserUUID.String())
	if err != nil {
		return nil, err
	}
	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
