package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type TokenServiceMock struct {
	mock.Mock
}

func (m *TokenServiceMock) GenerateTokenPair(userID, userRole string) (string, string, error) {
	args := m.Called(userID, userRole)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *TokenServiceMock) ValidateToken(token string, isRefreshToken bool) (string, string, error) {
	args := m.Called(token, isRefreshToken)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *TokenServiceMock) BlacklistToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
