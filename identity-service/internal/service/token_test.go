package service_test

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/errors"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/mocks"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/service"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func generateTestKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PrivateKey) {
	accessPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "failed to generate access private key")

	refreshPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err, "failed to generate refresh private key")

	return accessPrivateKey, refreshPrivateKey
}

func TestGenerateTokenPair_Success(t *testing.T) {
	accessPrivateKey, refreshPrivateKey := generateTestKeys(t)
	tokenRepo := new(mocks.TokenRepositoryMock)
	service := service.NewTokenService(accessPrivateKey, &accessPrivateKey.PublicKey,
		refreshPrivateKey, &refreshPrivateKey.PublicKey, 15*time.Minute, 7*24*time.Hour, tokenRepo)

	accessToken, refreshToken, err := service.GenerateTokenPair("user123", "admin")
	require.NoError(t, err, "должна быть успешная генерация токенов")
	require.NotEmpty(t, accessToken, "access-токен не должен быть пустым")
	require.NotEmpty(t, refreshToken, "refresh-токен не должен быть пустым")
}

func TestGenerateTokenPair_EmptyUserID(t *testing.T) {
	accessPrivateKey, refreshPrivateKey := generateTestKeys(t)
	tokenRepo := new(mocks.TokenRepositoryMock)
	service := service.NewTokenService(accessPrivateKey, &accessPrivateKey.PublicKey,
		refreshPrivateKey, &refreshPrivateKey.PublicKey, 15*time.Minute, 7*24*time.Hour, tokenRepo)

	_, _, err := service.GenerateTokenPair("", "admin")
	require.Error(t, err, "должна вернуться ошибка для пустого userID")
	require.Equal(t, errors.ErrEmptyUserID, err, "ошибка должна быть ErrEmptyUserID")
}

func TestGenerateTokenPair_EmptyUserRole(t *testing.T) {
	accessPrivateKey, refreshPrivateKey := generateTestKeys(t)
	tokenRepo := new(mocks.TokenRepositoryMock)
	service := service.NewTokenService(accessPrivateKey, &accessPrivateKey.PublicKey,
		refreshPrivateKey, &refreshPrivateKey.PublicKey, 15*time.Minute, 7*24*time.Hour, tokenRepo)

	accessToken, _, err := service.GenerateTokenPair("user123", "")
	require.NoError(t, err, "должна быть успешная генерация токенов с ролью по умолчанию")

	claims, err := service.ParseToken(accessToken, &accessPrivateKey.PublicKey)
	require.NoError(t, err, "токен должен успешно парситься")
	require.Equal(t, "user", claims.Role, "роль должна быть по умолчанию 'user'")
}

func TestValidateToken_AccessToken(t *testing.T) {
	accessPrivateKey, refreshPrivateKey := generateTestKeys(t)
	tokenRepo := new(mocks.TokenRepositoryMock)
	service := service.NewTokenService(accessPrivateKey, &accessPrivateKey.PublicKey,
		refreshPrivateKey, &refreshPrivateKey.PublicKey, 15*time.Minute, 7*24*time.Hour, tokenRepo)

	userID := "user123"
	userRole := "admin"
	accessToken, err := service.GenerateToken(userID, userRole, accessPrivateKey, 15*time.Minute)
	require.NoError(t, err, "должен успешно сгенерироваться access-токен")

	id, role, err := service.ValidateToken(accessToken, false)
	require.NoError(t, err, "валидация access-токена должна пройти успешно")
	require.Equal(t, userID, id, "userID должен совпадать")
	require.Equal(t, userRole, role, "роль должна совпадать")
}
