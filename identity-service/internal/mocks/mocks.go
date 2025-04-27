package mocks

import (
	"context"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/usecase"
	"github.com/stretchr/testify/mock"
)

// IdentityUseCaseMock мокирует интерфейс IdentityUseCase
type IdentityUseCaseMock struct {
	mock.Mock
}

func (m *IdentityUseCaseMock) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	args := m.Called(ctx, token)
	return args.Bool(0), args.Error(1)
}

func (m *IdentityUseCaseMock) GetByUserUUID(ctx context.Context, userUUID string) (*entity.Identity, error) {
	args := m.Called(ctx, userUUID)
	return args.Get(0).(*entity.Identity), args.Error(1)
}

func (m *IdentityUseCaseMock) Register(ctx context.Context, email, password string) error {
	args := m.Called(ctx, email, password)
	return args.Error(0)
}

func (m *IdentityUseCaseMock) Login(ctx context.Context, email, password string) (*usecase.Tokens, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(*usecase.Tokens), args.Error(1)
}

func (m *IdentityUseCaseMock) RefreshToken(ctx context.Context, refreshToken string) (*usecase.Tokens, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(*usecase.Tokens), args.Error(1)
}

func (m *IdentityUseCaseMock) Logout(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *IdentityUseCaseMock) ResetPassword(ctx context.Context, email, newPassword string) error {
	args := m.Called(ctx, email, newPassword)
	return args.Error(0)
}

func (m *IdentityUseCaseMock) AddVerificationCode(ctx context.Context, email, code string) error {
	args := m.Called(ctx, email, code)
	return args.Error(0)
}

func (m *IdentityUseCaseMock) VerifyCode(ctx context.Context, email, code string) (bool, error) {
	args := m.Called(ctx, email, code)
	return args.Bool(0), args.Error(1)
}

func (m *IdentityUseCaseMock) ValidateToken(ctx context.Context, token string, isRefreshToken bool) (string, error) {
	args := m.Called(ctx, token, isRefreshToken)
	return args.String(0), args.Error(1)
}

func (m *IdentityUseCaseMock) ConfirmEmail(ctx context.Context, email, code string) error {
	args := m.Called(ctx, email, code)
	return args.Error(0)
}

// VerificationServiceMock мокирует интерфейс VerificationService
type VerificationServiceMock struct {
	mock.Mock
}

func (m *VerificationServiceMock) GenerateVerificationCode() string {
	args := m.Called()
	return args.String(0)
}

func (m *VerificationServiceMock) SendVerificationCode(email, code string) error {
	args := m.Called(email, code)
	return args.Error(0)
}

// TokenRepositoryMock мокирует интерфейс TokenRepository
type TokenRepositoryMock struct {
	mock.Mock
}

func (m *TokenRepositoryMock) BlacklistToken(ctx context.Context, token *entity.BlacklistedToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *TokenRepositoryMock) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	args := m.Called(ctx, token)
	return args.Bool(0), args.Error(1)
}

func (m *TokenRepositoryMock) CleanupExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
