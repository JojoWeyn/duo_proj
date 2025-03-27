package service

import (
	"context"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type TokenRepository interface {
	BlacklistToken(ctx context.Context, token *entity.BlacklistedToken) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
	CleanupExpired(ctx context.Context) error
}

type Claims struct {
	Sub  string `json:"sub"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type TokenService struct {
	accessPrivateKey  *rsa.PrivateKey
	refreshPrivateKey *rsa.PrivateKey
	accessPublicKey   *rsa.PublicKey
	refreshPublicKey  *rsa.PublicKey
	accessTimeout     time.Duration
	refreshTimeout    time.Duration
	tokenRepo         TokenRepository
}

func NewTokenService(accessPrivateKey *rsa.PrivateKey, accessPublicKey *rsa.PublicKey,
	refreshPrivateKey *rsa.PrivateKey, refreshPublicKey *rsa.PublicKey,
	accessTimeout time.Duration, refreshTimeout time.Duration, tokenRepo TokenRepository) *TokenService {
	return &TokenService{
		accessPrivateKey:  accessPrivateKey,
		accessPublicKey:   accessPublicKey,
		refreshPrivateKey: refreshPrivateKey,
		refreshPublicKey:  refreshPublicKey,
		accessTimeout:     accessTimeout,
		refreshTimeout:    refreshTimeout,
		tokenRepo:         tokenRepo,
	}
}

func (s *TokenService) GenerateTokenPair(userID, userRole string) (string, string, error) {
	accessToken, err := s.generateToken(userID, userRole, s.accessPrivateKey, s.accessTimeout)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.generateToken(userID, userRole, s.refreshPrivateKey, s.refreshTimeout)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *TokenService) ValidateToken(token string, isRefreshToken bool) (string, string, error) {
	publicKey := s.accessPublicKey
	if isRefreshToken {
		publicKey = s.refreshPublicKey
	}

	claims, err := s.parseToken(token, publicKey)
	if err != nil {
		return "", "", err
	}

	return claims.Sub, claims.Role, nil
}

func (s *TokenService) parseToken(tokenString string, publicKey *rsa.PublicKey) (*Claims, error) {
	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}

func (s *TokenService) generateToken(userID, userRole string, privateKey *rsa.PrivateKey, expiration time.Duration) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("userID must not be empty")
	}
	if userRole == "" {
		userRole = "user"
	}

	claims := Claims{
		Sub:  userID,
		Role: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "identity-service",
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func (s *TokenService) BlacklistToken(ctx context.Context, token string) error {
	claims, err := s.parseToken(token, s.accessPublicKey)
	if err != nil {
		claims, err = s.parseToken(token, s.refreshPublicKey)
		if err != nil {
			return err
		}
	}

	blacklistedToken := entity.NewBlacklistedToken(token, time.Unix(claims.ExpiresAt.Unix(), 0))
	return s.tokenRepo.BlacklistToken(ctx, blacklistedToken)
}

func (s *TokenService) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	return s.tokenRepo.IsBlacklisted(ctx, token)
}
