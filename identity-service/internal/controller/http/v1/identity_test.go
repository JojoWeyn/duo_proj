package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/JojoWeyn/duo-proj/identity-service/internal/controller/http/v1"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/usecase"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Тест для POST /auth/register - Успешная регистрация
func TestRegister_Success(t *testing.T) {
	// Initialize mocks
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockVerification := new(mocks.VerificationServiceMock)

	// Set up the mock expectation: Register returns nil error on success
	mockUseCase.On("Register", context.Background(), "test@example.com", "password123.").Return(nil)

	// Set up Gin router
	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	// Create request body
	reqBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password123.",
	})
	req, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert expectations
	require.Equal(t, http.StatusOK, w.Code)
	expectedResponse := `{"message":"user registered successfully"}`
	require.JSONEq(t, expectedResponse, w.Body.String())

	// Verify mock was called as expected
	mockUseCase.AssertExpectations(t)
}

// Тест для POST /auth/register - Ошибка (существующий пользователь)
func TestRegister_ExistingUser(t *testing.T) {
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockUseCase.On("Register", context.Background(), "test@example.com", "password123.").Return(errors.New("user already exists"))

	mockVerification := new(mocks.VerificationServiceMock)

	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	reqBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password123.",
	})
	req, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"error":"user already exists"}`, w.Body.String())

	mockUseCase.AssertExpectations(t)
}

// Тест для POST /auth/login - Успешный вход
func TestLogin_Success(t *testing.T) {
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockUseCase.On("Login", context.Background(), "test@example.com", "password123").Return(&usecase.Tokens{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
	}, nil)

	mockVerification := new(mocks.VerificationServiceMock)

	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	reqBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	})
	req, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	expectedResponse := `{"access_token":"access_token","refresh_token":"refresh_token"}`
	require.JSONEq(t, expectedResponse, w.Body.String())

	mockUseCase.AssertExpectations(t)
}

// Тест для POST /auth/login - Неверные данные
func TestLogin_InvalidCredentials(t *testing.T) {
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockUseCase.On("Login", context.Background(), "test@example.com", "wrongpassword").Return((*usecase.Tokens)(nil), errors.New("invalid email or password"))

	mockVerification := new(mocks.VerificationServiceMock)

	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	reqBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "wrongpassword",
	})
	req, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.JSONEq(t, `{"error":"invalid email or password"}`, w.Body.String())

	mockUseCase.AssertExpectations(t)
}

// Тест для POST /auth/refresh - Успешное обновление
func TestRefresh_Success(t *testing.T) {
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockUseCase.On("RefreshToken", context.Background(), "valid_refresh_token").Return(&usecase.Tokens{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
	}, nil)

	mockVerification := new(mocks.VerificationServiceMock)

	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	reqBody, _ := json.Marshal(map[string]string{
		"refresh_token": "valid_refresh_token",
	})
	req, _ := http.NewRequest("POST", "/v1/auth/refresh", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	expectedResponse := `{"access_token":"new_access_token","refresh_token":"new_refresh_token"}`
	require.JSONEq(t, expectedResponse, w.Body.String())

	mockUseCase.AssertExpectations(t)
}

// Тест для POST /auth/refresh - Неверный токен
func TestRefresh_InvalidToken(t *testing.T) {
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockUseCase.On("RefreshToken", context.Background(), "invalid_refresh_token").Return((*usecase.Tokens)(nil), errors.New("invalid token"))

	mockVerification := new(mocks.VerificationServiceMock)

	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	reqBody, _ := json.Marshal(map[string]string{
		"refresh_token": "invalid_refresh_token",
	})
	req, _ := http.NewRequest("POST", "/v1/auth/refresh", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.JSONEq(t, `{"error":"invalid token"}`, w.Body.String())

	mockUseCase.AssertExpectations(t)
}

// Тест для POST /auth/logout - Успешный выход
func TestLogout_Success(t *testing.T) {
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockUseCase.On("Logout", context.Background(), "valid_token").Return(nil)

	mockVerification := new(mocks.VerificationServiceMock)

	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	req, _ := http.NewRequest("POST", "/v1/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer valid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	mockUseCase.AssertExpectations(t)
}

// Тест для POST /auth/logout - Неверный токен
func TestLogout_InvalidToken(t *testing.T) {
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockUseCase.On("Logout", context.Background(), "invalid_token").Return(errors.New("invalid token"))

	mockVerification := new(mocks.VerificationServiceMock)

	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	req, _ := http.NewRequest("POST", "/v1/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.JSONEq(t, `{"error":"invalid token"}`, w.Body.String())

	mockUseCase.AssertExpectations(t)
}

// Тест для GET /auth/token/status - Действительный токен
func TestCheckToken_Valid(t *testing.T) {
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockUseCase.On("ValidateToken", context.Background(), "valid_token", false).Return("user123", nil)

	mockVerification := new(mocks.VerificationServiceMock)

	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	req, _ := http.NewRequest("GET", "/v1/auth/token/status", nil)
	req.Header.Set("Authorization", "Bearer valid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"is_blacklisted":false}`, w.Body.String())

	mockUseCase.AssertExpectations(t)
}

// Тест для GET /auth/token/status - Недействительный токен
func TestCheckToken_Invalid(t *testing.T) {
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockUseCase.On("ValidateToken", context.Background(), "invalid_token", false).Return("", errors.New("invalid token"))
	mockUseCase.On("ValidateToken", context.Background(), "invalid_token", true).Return("", errors.New("invalid token"))

	mockVerification := new(mocks.VerificationServiceMock)

	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	req, _ := http.NewRequest("GET", "/v1/auth/token/status", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.JSONEq(t, `{"is_blacklisted":true}`, w.Body.String())

	mockUseCase.AssertExpectations(t)
}

// Тест для POST /auth/password/reset - Успешный сброс
func TestResetPassword_Success(t *testing.T) {
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockUseCase.On("VerifyCode", context.Background(), "test@example.com", "123456").Return(true, nil)
	mockUseCase.On("ResetPassword", context.Background(), "test@example.com", "newpassword").Return(nil)

	mockVerification := new(mocks.VerificationServiceMock)

	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	reqBody, _ := json.Marshal(map[string]string{
		"email":        "test@example.com",
		"code":         "123456",
		"new_password": "newpassword",
	})
	req, _ := http.NewRequest("POST", "/v1/auth/password/reset", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"message":"password reset successfully"}`, w.Body.String())

	mockUseCase.AssertExpectations(t)
}

// Тест для POST /auth/password/reset - Неверный код
func TestResetPassword_InvalidCode(t *testing.T) {
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockUseCase.On("VerifyCode", context.Background(), "test@example.com", "wrongcode").Return(false, nil)

	mockVerification := new(mocks.VerificationServiceMock)

	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	reqBody, _ := json.Marshal(map[string]string{
		"email":        "test@example.com",
		"code":         "wrongcode",
		"new_password": "newpassword",
	})
	req, _ := http.NewRequest("POST", "/v1/auth/password/reset", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"error":"invalid verification code"}`, w.Body.String())

	mockUseCase.AssertExpectations(t)
}

// Тест для POST /auth/verification/code - Неверный email
func TestSendVerificationCode_InvalidEmail(t *testing.T) {
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockUseCase.On("AddVerificationCode", context.Background(), "invalid@example.com", mock.Anything).Return(errors.New("email not found"))

	mockVerification := new(mocks.VerificationServiceMock)
	mockVerification.On("GenerateVerificationCode").Return("123456")

	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	req, _ := http.NewRequest("POST", "/v1/auth/verification/code?email=invalid@example.com", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.JSONEq(t, `{"error":"email not found"}`, w.Body.String())

	mockUseCase.AssertExpectations(t)
	mockVerification.AssertExpectations(t)
}

// Тест для POST /auth/verification/email - Успешное подтверждение
func TestConfirmEmail_Success(t *testing.T) {
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockUseCase.On("ConfirmEmail", context.Background(), "test@example.com", "123456").Return(nil)

	mockVerification := new(mocks.VerificationServiceMock)

	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	reqBody, _ := json.Marshal(map[string]string{
		"email": "test@example.com",
		"code":  "123456",
	})
	req, _ := http.NewRequest("POST", "/v1/auth/verification/email", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"message":"email confirmed successfully"}`, w.Body.String())

	mockUseCase.AssertExpectations(t)
}

// Тест для POST /auth/verification/email - Неверный код
func TestConfirmEmail_InvalidCode(t *testing.T) {
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockUseCase.On("ConfirmEmail", context.Background(), "test@example.com", "wrongcode").Return(errors.New("invalid code"))

	mockVerification := new(mocks.VerificationServiceMock)

	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	reqBody, _ := json.Marshal(map[string]string{
		"email": "test@example.com",
		"code":  "wrongcode",
	})
	req, _ := http.NewRequest("POST", "/v1/auth/verification/email", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"error":"invalid code"}`, w.Body.String())

	mockUseCase.AssertExpectations(t)
}

// Тест для GET /auth/me - Успешное получение
func TestGetIdentity_Success(t *testing.T) {
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockUseCase.On("ValidateToken", context.Background(), "valid_token", false).Return("user123", nil)
	mockUseCase.On("GetByUserUUID", context.Background(), "user123").Return(&entity.Identity{Email: "test@example.com"}, nil)

	mockVerification := new(mocks.VerificationServiceMock)

	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	req, _ := http.NewRequest("GET", "/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer valid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"email":"test@example.com"}`, w.Body.String())

	mockUseCase.AssertExpectations(t)
}

// Тест для GET /auth/me - Неверный токен
func TestGetIdentity_InvalidToken(t *testing.T) {
	mockUseCase := new(mocks.IdentityUseCaseMock)
	mockUseCase.On("ValidateToken", context.Background(), "invalid_token", false).Return("", errors.New("invalid token"))

	mockVerification := new(mocks.VerificationServiceMock)

	router := gin.Default()
	v1.NewIdentityRoutes(router.Group("/v1"), mockVerification, mockUseCase)

	req, _ := http.NewRequest("GET", "/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.JSONEq(t, `{"error":"invalid token"}`, w.Body.String())

	mockUseCase.AssertExpectations(t)
}
