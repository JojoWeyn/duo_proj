package mocks

import (
	"context"

	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type IdentityRepositoryMock struct {
	mock.Mock
}

func (m *IdentityRepositoryMock) Create(ctx context.Context, identity *entity.Identity) error {
	args := m.Called(ctx, identity)
	return args.Error(0)
}

func (m *IdentityRepositoryMock) FindByUUID(ctx context.Context, userID string) (*entity.Identity, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*entity.Identity), args.Error(1)
}

func (m *IdentityRepositoryMock) FindByLogin(ctx context.Context, login string) (*entity.Identity, error) {
	args := m.Called(ctx, login)
	return args.Get(0).(*entity.Identity), args.Error(1)
}

func (m *IdentityRepositoryMock) FindByEmail(ctx context.Context, email string) (*entity.Identity, error) {
	args := m.Called(ctx, email)
	identity := args.Get(0)
	if identity == nil {
		return nil, args.Error(1)
	}
	return identity.(*entity.Identity), args.Error(1)
}

func (m *IdentityRepositoryMock) Update(ctx context.Context, identity *entity.Identity) error {
	args := m.Called(ctx, identity)
	return args.Error(0)
}
