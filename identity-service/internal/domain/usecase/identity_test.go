package usecase_test

import (
	"context"
	"errors"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/usecase"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/usecase/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegister_Success(t *testing.T) {
	ctx := context.TODO()

	identityRepo := new(mocks.IdentityRepositoryMock)
	tokenService := new(mocks.TokenServiceMock)
	tokenRepo := new(mocks.TokenRepositoryMock)
	producer := new(mocks.ProducerMock)

	identityRepo.On("FindByEmail", ctx, "test@example.com").Return(nil, errors.New("not found"))
	identityRepo.On("Create", ctx, mock.AnythingOfType("*entity.Identity")).Return(nil)

	uc := usecase.NewIdentityUseCase(identityRepo, tokenService, tokenRepo, producer)

	err := uc.Register(ctx, "test@example.com", "StrongP@ssw0rd")

	require.NoError(t, err)
	identityRepo.AssertExpectations(t)
}

func TestLogin_InvalidEmail(t *testing.T) {
	ctx := context.TODO()

	identityRepo := new(mocks.IdentityRepositoryMock)
	tokenService := new(mocks.TokenServiceMock)
	tokenRepo := new(mocks.TokenRepositoryMock)
	producer := new(mocks.ProducerMock)

	identityRepo.On("FindByEmail", ctx, "test@example.com").Return(nil, errors.New("not found"))

	uc := usecase.NewIdentityUseCase(identityRepo, tokenService, tokenRepo, producer)

	tokens, err := uc.Login(ctx, "test@example.com", "wrongpass")

	require.Nil(t, tokens)
	require.Error(t, err)
	require.EqualError(t, err, "invalid email or password")
}

func TestLogin_EmailNotConfirmed(t *testing.T) {
	ctx := context.TODO()

	identityRepo := new(mocks.IdentityRepositoryMock)
	tokenService := new(mocks.TokenServiceMock)
	tokenRepo := new(mocks.TokenRepositoryMock)
	producer := new(mocks.ProducerMock)

	identity := &entity.Identity{
		Email:          "test@example.com",
		PasswordHash:   hashPassword("password123"),
		IsConfirmEmail: false,
	}

	identityRepo.On("FindByEmail", ctx, "test@example.com").Return(identity, nil)

	uc := usecase.NewIdentityUseCase(identityRepo, tokenService, tokenRepo, producer)

	tokens, err := uc.Login(ctx, "test@example.com", "password123")

	require.Nil(t, tokens)
	require.Error(t, err)
	require.EqualError(t, err, "emails is not confirmed")
}

func TestConfirmEmail_Success(t *testing.T) {
	ctx := context.TODO()

	identityRepo := new(mocks.IdentityRepositoryMock)
	tokenService := new(mocks.TokenServiceMock)
	tokenRepo := new(mocks.TokenRepositoryMock)
	producer := new(mocks.ProducerMock)

	identity := &entity.Identity{
		Email:            "test@example.com",
		UserUUID:         uuid.New(),
		VerificationCode: "123456",
	}

	identityRepo.On("FindByEmail", ctx, "test@example.com").Return(identity, nil)
	identityRepo.On("Update", ctx, mock.AnythingOfType("*entity.Identity")).Return(nil)
	producer.On("SendUserCreated", identity.UserUUID.String(), "test@example.com").Return(nil)

	uc := usecase.NewIdentityUseCase(identityRepo, tokenService, tokenRepo, producer)

	err := uc.ConfirmEmail(ctx, "test@example.com", "123456")

	require.NoError(t, err)
	identityRepo.AssertExpectations(t)
	producer.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	ctx := context.TODO()

	identityRepo := new(mocks.IdentityRepositoryMock)
	tokenService := new(mocks.TokenServiceMock)
	tokenRepo := new(mocks.TokenRepositoryMock)
	producer := new(mocks.ProducerMock)

	pass := hashPassword("password123")
	identity, err := entity.NewIdentity("test@example.com", pass)
	identity.IsConfirmEmail = true

	require.NoError(t, err)

	identityRepo.On("FindByEmail", ctx, "test@example.com").Return(identity, nil)
	tokenService.On("GenerateTokenPair", identity.UserUUID.String(), "user").Return("access", "refresh", nil)
	producer.On("SendUserLogin", identity.UserUUID.String(), "test@example.com").Return(nil)

	uc := usecase.NewIdentityUseCase(identityRepo, tokenService, tokenRepo, producer)

	tokens, err := uc.Login(ctx, "test@example.com", "password123")

	require.NoError(t, err)
	require.Equal(t, "access", tokens.AccessToken)
	require.Equal(t, "refresh", tokens.RefreshToken)
}

func TestRefreshToken_Success(t *testing.T) {
	ctx := context.TODO()

	identityRepo := new(mocks.IdentityRepositoryMock)
	tokenService := new(mocks.TokenServiceMock)
	tokenRepo := new(mocks.TokenRepositoryMock)
	producer := new(mocks.ProducerMock)

	tokenService.On("ValidateToken", "refresh_token", true).Return("user-id", "user", nil)
	tokenService.On("BlacklistToken", ctx, "refresh_token").Return(nil)
	tokenService.On("GenerateTokenPair", "user-id", "user").Return("new_access", "new_refresh", nil)

	uc := usecase.NewIdentityUseCase(identityRepo, tokenService, tokenRepo, producer)

	tokens, err := uc.RefreshToken(ctx, "refresh_token")

	require.NoError(t, err)
	require.Equal(t, "new_access", tokens.AccessToken)
	require.Equal(t, "new_refresh", tokens.RefreshToken)
}

func TestLogout_Success(t *testing.T) {
	ctx := context.TODO()

	tokenService := new(mocks.TokenServiceMock)

	tokenService.On("BlacklistToken", ctx, "some_token").Return(nil)

	uc := usecase.NewIdentityUseCase(nil, tokenService, nil, nil)

	err := uc.Logout(ctx, "some_token")
	require.NoError(t, err)
}

func TestValidateToken_Blacklisted(t *testing.T) {
	ctx := context.TODO()

	tokenService := new(mocks.TokenServiceMock)
	tokenRepo := new(mocks.TokenRepositoryMock)

	tokenRepo.On("IsBlacklisted", ctx, "some_token").Return(true, nil)

	uc := usecase.NewIdentityUseCase(nil, tokenService, tokenRepo, nil)

	uid, err := uc.ValidateToken(ctx, "some_token", false)
	require.Error(t, err)
	require.Empty(t, uid)
	require.EqualError(t, err, "token is blacklisted")
}

func TestConfirmEmail_UserNotFound(t *testing.T) {
	ctx := context.TODO()

	identityRepo := new(mocks.IdentityRepositoryMock)
	producer := new(mocks.ProducerMock)

	identityRepo.On("FindByEmail", ctx, "test@example.com").Return(nil, errors.New("not found"))

	uc := usecase.NewIdentityUseCase(identityRepo, nil, nil, producer)

	err := uc.ConfirmEmail(ctx, "test@example.com", "123456")

	require.Error(t, err)
	require.EqualError(t, err, "user not found")
}

func TestConfirmEmail_WrongCode(t *testing.T) {
	ctx := context.TODO()

	identityRepo := new(mocks.IdentityRepositoryMock)
	producer := new(mocks.ProducerMock)

	identity := &entity.Identity{
		Email:            "test@example.com",
		UserUUID:         uuid.New(),
		VerificationCode: "123456",
	}

	identityRepo.On("FindByEmail", ctx, "test@example.com").Return(identity, nil)

	uc := usecase.NewIdentityUseCase(identityRepo, nil, nil, producer)

	err := uc.ConfirmEmail(ctx, "test@example.com", "000000")

	require.Error(t, err)
	require.EqualError(t, err, "verification code is incorrect")
}

func TestValidateToken_Success(t *testing.T) {
	ctx := context.TODO()

	tokenService := new(mocks.TokenServiceMock)
	tokenRepo := new(mocks.TokenRepositoryMock)

	tokenRepo.On("IsBlacklisted", ctx, "valid_token").Return(false, nil)
	tokenService.On("ValidateToken", "valid_token", false).Return("user-id", "user", nil)

	uc := usecase.NewIdentityUseCase(nil, tokenService, tokenRepo, nil)

	uid, err := uc.ValidateToken(ctx, "valid_token", false)
	require.NoError(t, err)
	require.Equal(t, "user-id", uid)
}

func TestVerifyCode_Success(t *testing.T) {
	ctx := context.TODO()

	identityRepo := new(mocks.IdentityRepositoryMock)

	identity := &entity.Identity{
		Email:            "test@example.com",
		VerificationCode: "123456",
	}

	identityRepo.On("FindByEmail", ctx, "test@example.com").Return(identity, nil)
	identityRepo.On("Update", ctx, mock.AnythingOfType("*entity.Identity")).Return(nil)

	uc := usecase.NewIdentityUseCase(identityRepo, nil, nil, nil)

	ok, err := uc.VerifyCode(ctx, "test@example.com", "123456")

	require.True(t, ok)
	require.NoError(t, err)
}

func TestAddVerificationCode_Success(t *testing.T) {
	ctx := context.TODO()

	identityRepo := new(mocks.IdentityRepositoryMock)
	identity := &entity.Identity{Email: "test@example.com"}

	identityRepo.On("FindByEmail", ctx, "test@example.com").Return(identity, nil)
	identityRepo.On("Update", ctx, identity).Return(nil)

	uc := usecase.NewIdentityUseCase(identityRepo, nil, nil, nil)

	err := uc.AddVerificationCode(ctx, "test@example.com", "654321")
	require.NoError(t, err)
}

func TestResetPassword_Success(t *testing.T) {
	ctx := context.TODO()

	identityRepo := new(mocks.IdentityRepositoryMock)
	identity := &entity.Identity{
		Email: "test@example.com",
	}

	identityRepo.On("FindByEmail", ctx, "test@example.com").Return(identity, nil)
	identityRepo.On("Update", ctx, identity).Return(nil)

	uc := usecase.NewIdentityUseCase(identityRepo, nil, nil, nil)

	err := uc.ResetPassword(ctx, "test@example.com", "NewP@ssw0rd")
	require.NoError(t, err)
}

func TestIsBlacklisted_Success(t *testing.T) {
	ctx := context.TODO()

	tokenRepo := new(mocks.TokenRepositoryMock)
	tokenRepo.On("IsBlacklisted", ctx, "some_token").Return(true, nil)

	uc := usecase.NewIdentityUseCase(nil, nil, tokenRepo, nil)

	result, err := uc.IsBlacklisted(ctx, "some_token")
	require.NoError(t, err)
	require.True(t, result)
}

func TestGetByUserUUID_Success(t *testing.T) {
	ctx := context.TODO()

	identityRepo := new(mocks.IdentityRepositoryMock)
	identity := &entity.Identity{
		UserUUID: uuid.New(),
		Email:    "test@example.com",
	}

	identityRepo.On("FindByUUID", ctx, identity.UserUUID.String()).Return(identity, nil)

	uc := usecase.NewIdentityUseCase(identityRepo, nil, nil, nil)

	result, err := uc.GetByUserUUID(ctx, identity.UserUUID.String())
	require.NoError(t, err)
	require.Equal(t, identity, result)
}

func hashPassword(pwd string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash)
}
