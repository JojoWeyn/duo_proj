package v1

import (
	"context"
	"net/http"
	"strings"

	"github.com/JojoWeyn/duo-proj/identity-service/internal/controller/http/dto"
	"github.com/JojoWeyn/duo-proj/identity-service/internal/domain/usecase"
	"github.com/gin-gonic/gin"
)

type TokenRepository interface {
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

type IdentityUsecase interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (*usecase.Tokens, error)
	RefreshToken(ctx context.Context, refreshToken string) (*usecase.Tokens, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string, isRefreshToken bool) (string, error)
}

type identityRoutes struct {
	identityUsecase IdentityUsecase
	tokenRepo       TokenRepository
}

func newIdentityRoutes(handler *gin.RouterGroup, identityUsecase IdentityUsecase, tokenRepository TokenRepository) {
	r := &identityRoutes{
		identityUsecase: identityUsecase,
		tokenRepo:       tokenRepository,
	}

	h := handler.Group("/auth")
	{
		h.POST("/register", r.register)
		h.POST("/login", r.login)
		h.POST("/refresh", r.refresh)
		h.POST("/logout", r.logout)
		h.GET("/token/status", r.checkToken)
	}
}

func (r *identityRoutes) checkToken(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no token provided"})
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")

	validateToken := func(token string, checkRefresh bool) (string, error) {
		isBlacklisted, err := r.identityUsecase.ValidateToken(c.Request.Context(), token, checkRefresh)
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

	c.JSON(http.StatusOK, gin.H{"token_blacklisted": isBlacklisted})
}

func (r *identityRoutes) register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.identityUsecase.Register(c.Request.Context(), req.Email, req.Password)
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

	tokens, err := r.identityUsecase.Login(c.Request.Context(), req.Email, req.Password)
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

	tokens, err := r.identityUsecase.RefreshToken(c.Request.Context(), req.RefreshToken)
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

	if err := r.identityUsecase.Logout(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
