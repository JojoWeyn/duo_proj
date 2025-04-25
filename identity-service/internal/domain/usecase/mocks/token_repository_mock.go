package mocks

import (
	"context"

	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

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
