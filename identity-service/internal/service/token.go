package service

import (
	"context"
	"time"

	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"
	"github.com/golang-jwt/jwt/v4"
)

type TokenRepository interface {
	BlacklistToken(ctx context.Context, token *entity.BlacklistedToken) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
	CleanupExpired(ctx context.Context) error
}

type TokenService struct {
	accessSecret   string
	refreshSecret  string
	accessTimeout  time.Duration
	refreshTimeout time.Duration
	tokenRepo      TokenRepository
}

func NewTokenService(
	accessSecret string,
	refreshSecret string,
	accessTimeout time.Duration,
	refreshTimeout time.Duration,
	tokenRepo TokenRepository,
) *TokenService {
	return &TokenService{
		accessSecret:   accessSecret,
		refreshSecret:  refreshSecret,
		accessTimeout:  accessTimeout,
		refreshTimeout: refreshTimeout,
		tokenRepo:      tokenRepo,
	}
}

func (s *TokenService) BlacklistToken(ctx context.Context, token string) error {
	claims, err := s.parseToken(token, s.accessSecret)
	if err != nil {
		claims, err = s.parseToken(token, s.refreshSecret)
		if err != nil {
			return err
		}
	}

	blacklistedToken := entity.NewBlacklistedToken(token, time.Unix(claims.ExpiresAt.Unix(), 0))
	return s.tokenRepo.BlacklistToken(ctx, blacklistedToken)
}

func (s *TokenService) parseToken(tokenString, secret string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func (s *TokenService) GenerateTokenPair(userID string) (string, string, error) {
	accessToken, err := s.generateToken(userID, s.accessSecret, s.accessTimeout)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.generateToken(userID, s.refreshSecret, s.refreshTimeout)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *TokenService) ValidateToken(token string, isRefreshToken bool) (string, error) {
	secret := s.accessSecret
	if isRefreshToken {
		secret = s.refreshSecret
	}

	claims := &jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if !parsedToken.Valid {
		return "", jwt.ErrSignatureInvalid
	}

	return claims.Subject, nil
}

func (s *TokenService) generateToken(userID string, secret string, expiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	return token.SignedString([]byte(secret))
}
