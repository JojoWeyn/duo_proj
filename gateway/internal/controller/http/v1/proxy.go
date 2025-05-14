package v1

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"errors"
	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	timeout      time.Duration
	jwtPublicKey *rsa.PublicKey
}

type Claims struct {
	Sub  string `json:"sub"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewProxyHandler(jwtPublicKey *rsa.PublicKey) *ProxyHandler {
	return &ProxyHandler{
		timeout:      10 * time.Second,
		jwtPublicKey: jwtPublicKey,
	}
}

func (h *ProxyHandler) ProxyService(serviceURL string, addUUID bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedURL, err := url.Parse(serviceURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid service URL"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(parsedURL)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			c.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable"})
		}

		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			if auth := c.GetHeader("Authorization"); auth != "" {
				req.Header.Set("Authorization", auth)
				if addUUID {
					if uuid, role, err := h.extractUUIDFromJWT(auth); err == nil {
						req.Header.Set("X-User-UUID", uuid)
						req.Header.Set("X-User-Role", role)
					}
				}
			}
			req.Header.Set("Origin", "http://37.18.102.166:3211")
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func (p *ProxyHandler) ProxySwagger(targetService, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		target := fmt.Sprintf("%s/swagger%s", targetService, path)

		req, err := http.NewRequest(c.Request.Method, target, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
			return
		}

		req.Header = c.Request.Header

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to forward request"})
			return
		}
		defer resp.Body.Close()

		for k, v := range resp.Header {
			for _, vv := range v {
				c.Writer.Header().Add(k, vv)
			}
		}

		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	}
}

func (h *ProxyHandler) extractUUIDFromJWT(tokenString string) (string, string, error) {
	if tokenString == "" {
		return "", "", errors.New("empty token")
	}

	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", "", errors.New("invalid token format")
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return h.jwtPublicKey, nil
	})
	if err != nil {
		return "", "", err
	}

	if !token.Valid {
		return "", "", errors.New("invalid token")
	}

	if claims.Sub == "" {
		return "", "", errors.New("uuid is not found")
	}

	return claims.Sub, claims.Role, nil
}
