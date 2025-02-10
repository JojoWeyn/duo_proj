package v1

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	userServiceURL     *url.URL
	identityServiceURL *url.URL
	timeout            time.Duration
}

func NewProxyHandler(identityServiceURL, userServiceURL string) (*ProxyHandler, error) {
	parsedIdentityServiceURL, err := url.Parse(identityServiceURL)
	if err != nil {
		return nil, err
	}

	parsedUserServiceURL, err := url.Parse(userServiceURL)
	if err != nil {
		return nil, err
	}

	return &ProxyHandler{
		identityServiceURL: parsedIdentityServiceURL,
		userServiceURL:     parsedUserServiceURL,
		timeout:            10 * time.Second,
	}, nil
}

func (h *ProxyHandler) ProxyUserService() gin.HandlerFunc {
	return func(c *gin.Context) {
		proxy := httputil.NewSingleHostReverseProxy(h.userServiceURL)

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			c.JSON(http.StatusBadGateway, gin.H{"error": "user service unavailable"})
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func (h *ProxyHandler) ProxyIdentityService() gin.HandlerFunc {
	return func(c *gin.Context) {
		proxy := httputil.NewSingleHostReverseProxy(h.identityServiceURL)

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			c.JSON(http.StatusBadGateway, gin.H{"error": "identity service unavailable"})
		}

		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			if auth := c.GetHeader("Authorization"); auth != "" {
				req.Header.Set("Authorization", auth)
			}
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), h.timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
