package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

type TokenService struct {
	accessSecret string
}

func NewTokenService(accessSecret string) *TokenService {
	return &TokenService{
		accessSecret: accessSecret,
	}
}

func (s *TokenService) Validate(token string) (string, error) {
	secret := s.accessSecret

	// Объявляем claims для стандартных зарегистрированных полей
	claims := &jwt.RegisteredClaims{}

	// Парсим токен с проверкой подписи
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", fmt.Errorf("invalid token signature")
		}
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	// Проверка валидности токена
	if !parsedToken.Valid {
		return "", jwt.ErrSignatureInvalid
	}

	return string(claims.Subject), nil
}
