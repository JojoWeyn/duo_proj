package v1

import (
	"context"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/entity"
	"log"
	"net/http"
	"strings"

	"github.com/JojoWeyn/duo-proj/identity-service/internal/controller/http/dto"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/usecase"
	"github.com/gin-gonic/gin"
)

type IdentityUseCase interface {
	IsBlacklisted(ctx context.Context, token string) (bool, error)
	GetByUserUUID(ctx context.Context, userUUID string) (*entity.Identity, error)
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (*usecase.Tokens, error)
	RefreshToken(ctx context.Context, refreshToken string) (*usecase.Tokens, error)
	Logout(ctx context.Context, token string) error
	ResetPassword(ctx context.Context, email, newPassword string) error
	AddVerificationCode(ctx context.Context, email, code string) error
	VerifyCode(ctx context.Context, email, code string) (bool, error)
	ValidateToken(ctx context.Context, token string, isRefreshToken bool) (string, error)
	ConfirmEmail(ctx context.Context, email, code string) error
}

type VerificationService interface {
	GenerateVerificationCode() string
	SendVerificationCode(email, code string) error
}

type identityRoutes struct {
	identityUseCase     IdentityUseCase
	verificationService VerificationService
}

func newIdentityRoutes(handler *gin.RouterGroup, verificationService VerificationService, identityUseCase IdentityUseCase) {
	r := &identityRoutes{
		identityUseCase:     identityUseCase,
		verificationService: verificationService,
	}

	h := handler.Group("/auth")
	{
		h.POST("/register", r.register)
		h.POST("/login", r.login)
		h.POST("/refresh", r.refresh)
		h.POST("/logout", r.logout)
		h.GET("/token/status", r.checkToken)
		h.POST("/password/reset", r.resetPassword)
		h.POST("/verification/code", r.sendVerificationCode)
		h.POST("/verification/email", r.confirmEmail)
		h.GET("/me", r.getIdentity)
	}
}

func (r *identityRoutes) getIdentity(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no token provided"})
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")

	userUUID, err := r.identityUseCase.ValidateToken(c.Request.Context(), token, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	identity, err := r.identityUseCase.GetByUserUUID(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email": identity.Email,
	})

}

func (r *identityRoutes) sendVerificationCode(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	code := r.verificationService.GenerateVerificationCode()
	if err := r.identityUseCase.AddVerificationCode(c.Request.Context(), email, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func() {
		if err := r.verificationService.SendVerificationCode(email, code); err != nil {
			log.Println(err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "verification code generated",
		"code":    code,
	})
}

func (r *identityRoutes) confirmEmail(c *gin.Context) {
	var req dto.ConfirmEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.identityUseCase.ConfirmEmail(c.Request.Context(), req.Email, req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email confirmed successfully"})
}

func (r *identityRoutes) resetPassword(c *gin.Context) {
	var req dto.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, err := r.identityUseCase.VerifyCode(c.Request.Context(), req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !code {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid verification code"})
		return
	}

	err = r.identityUseCase.ResetPassword(c.Request.Context(), req.Email, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

func (r *identityRoutes) checkToken(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no token provided"})
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")

	validateToken := func(token string, checkRefresh bool) (string, error) {
		isBlacklisted, err := r.identityUseCase.ValidateToken(c.Request.Context(), token, checkRefresh)
		if err != nil {
			return "", err
		}
		return isBlacklisted, nil
	}

	isBlacklisted, err := validateToken(token, false)
	if err != nil {
		isBlacklisted, err = validateToken(token, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	if isBlacklisted == "" {
		c.JSON(http.StatusOK, gin.H{"is_blacklisted": "true"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"is_blacklisted": "false"})
}

func (r *identityRoutes) register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.identityUseCase.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})
}

func (r *identityRoutes) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := r.identityUseCase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (r *identityRoutes) refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := r.identityUseCase.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (r *identityRoutes) logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no token provided"})
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")

	if err := r.identityUseCase.Logout(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
