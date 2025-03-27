package v1

import (
	"context"
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	timeout      time.Duration
	jwtPublicKey *rsa.PublicKey
}

type Claims struct {
	Sub string `json:"sub"`
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
					if uuid, err := h.extractUUIDFromJWT(auth); err == nil {
						req.Header.Set("X-User-UUID", uuid)
					}
				}
			}

			req.Header.Set("Origin", "http://37.18.102.166:3211")
			log.Println(req.Header.Get("Origin"))
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
func (h *ProxyHandler) extractUUIDFromJWT(tokenString string) (string, error) {
	if tokenString == "" {
		return "", errors.New("empty token")
	}

	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid token format")
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return h.jwtPublicKey, nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	if claims.Sub == "" {
		return "", errors.New("uuid not found in token")
	}

	return claims.Sub, nil
}
