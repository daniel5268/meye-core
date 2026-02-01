package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	apiKey string
}

func NewAuthHandler(apiKey string) *AuthHandler {
	return &AuthHandler{
		apiKey: apiKey,
	}
}

func (h *AuthHandler) InternalAPIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("Api-Key")

		if apiKey != h.apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
