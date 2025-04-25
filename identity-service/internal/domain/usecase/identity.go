package usecase

import (
	"context"
	"errors"
	"log"

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
	GenerateTokenPair(userID, userRole string) (accessToken string, refreshToken string, err error)
	ValidateToken(token string, isRefreshToken bool) (userID, userRole string, err error)
	BlacklistToken(ctx context.Context, token string) error
}

type EventProducer interface {
	SendUserCreated(userUUID string, email string) error
	SendUserLogin(userUUID string, email string) error
}

type IdentityUseCase struct {
	identityRepo IdentityRepository
	tokenService TokenService
	tokenRepo    TokenRepository
	producer     EventProducer
}

func NewIdentityUseCase(identityRepo IdentityRepository, tokenService TokenService, tokenRepo TokenRepository, producer EventProducer) *IdentityUseCase {
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

	userID, _, err := uc.tokenService.ValidateToken(token, isRefreshToken)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (uc *IdentityUseCase) RefreshToken(ctx context.Context, refreshToken string) (*Tokens, error) {
	userID, userRole, err := uc.tokenService.ValidateToken(refreshToken, true)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if err := uc.tokenService.BlacklistToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	tokens, err := uc.generateTokens(userID, userRole)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (uc *IdentityUseCase) Register(ctx context.Context, email, password string) error {
	if _, err := uc.identityRepo.FindByEmail(ctx, email); err == nil {
		return errors.New("email already exists")
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
		email,
		string(hashedPassword),
	)
	if err != nil {
		return err
	}

	return uc.identityRepo.Create(ctx, identity)
}

func (uc *IdentityUseCase) ConfirmEmail(ctx context.Context, email, code string) error {
	identity, err := uc.identityRepo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	if identity.VerificationCode != code {
		return errors.New("verification code is incorrect")
	}

	identity.ConfirmEmail()

	identity.RemoveVerificationCode()

	if err := uc.identityRepo.Update(ctx, identity); err != nil {
		return err
	}

	if err := uc.producer.SendUserCreated(identity.UserUUID.String(), email); err != nil {
		log.Printf("Failed to send user created event: %v", err)
	}

	return nil

}

func (uc *IdentityUseCase) VerifyCode(ctx context.Context, email, code string) (bool, error) {
	identity, err := uc.identityRepo.FindByEmail(ctx, email)
	if err != nil {
		return false, errors.New("user not found")
	}
	if identity.VerificationCode != code {
		return false, errors.New("invalid verification code")
	}

	identity.RemoveVerificationCode()
	if err := uc.identityRepo.Update(ctx, identity); err != nil {
		return false, err
	}
	return true, nil

}

func (uc *IdentityUseCase) AddVerificationCode(ctx context.Context, email, code string) error {
	identity, err := uc.identityRepo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}
	identity.AddVerificationCode(code)

	return uc.identityRepo.Update(ctx, identity)
}

func (uc *IdentityUseCase) Login(ctx context.Context, email, password string) (*Tokens, error) {
	identity, err := uc.identityRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(identity.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !identity.IsConfirmEmail {
		return nil, errors.New("emails is not confirmed")
	}

	tokens, err := uc.generateTokens(identity.UserUUID.String(), identity.Role)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := uc.producer.SendUserLogin(identity.UserUUID.String(), email); err != nil {
			log.Printf("Failed to send user login event: %v", err)
		}
	}()

	return tokens, nil

}

func (uc *IdentityUseCase) ResetPassword(ctx context.Context, email, newPassword string) error {
	identity, err := uc.identityRepo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	if err := entity.ValidatePassword(newPassword); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	identity.UpdatePassword(string(hashedPassword))
	return uc.identityRepo.Update(ctx, identity)
}

func (uc *IdentityUseCase) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	return uc.tokenRepo.IsBlacklisted(ctx, token)
}

func (uc *IdentityUseCase) GetByUserUUID(ctx context.Context, userUUID string) (*entity.Identity, error) {
	identity, err := uc.identityRepo.FindByUUID(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	return identity, nil
}

func (uc *IdentityUseCase) generateTokens(userID, userRole string) (*Tokens, error) {
	accessToken, refreshToken, err := uc.tokenService.GenerateTokenPair(userID, userRole)
	if err != nil {
		return nil, err
	}
	return &Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
