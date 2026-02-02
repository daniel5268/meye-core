package handler

import (
	"meye-core/internal/domain/user"
	"meye-core/internal/infrastructure/jwt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	apiKey         string
	jwtService     jwt.Service
	userRepository user.Repository
}

const AuthKey = "auth"

func NewAuthHandler(apiKey string, jwtService jwt.Service, userRepo user.Repository) *AuthHandler {
	return &AuthHandler{
		apiKey:         apiKey,
		jwtService:     jwtService,
		userRepository: userRepo,
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

type AuthContext struct {
	UserID string
}

// AuthMiddleware is a Gin middleware that validates a JWT and sets the claims in the context.
func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		tokenString := parts[1]
		userID, err := h.jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization failed"})
			return
		}

		auth := AuthContext{
			UserID: userID,
		}

		c.Set(AuthKey, auth)
		c.Next()
	}
}

// RequireRole is a generic Gin middleware that checks if the authenticated user has any of the specified roles.
// This middleware should be used after AuthMiddleware, as it depends on the AuthContext being set.
// It accepts one or more roles and grants access if the user has any of them (OR logic).
func (h *AuthHandler) RequireRole(allowedRoles ...user.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get auth context set by AuthMiddleware
		authValue, exists := c.Get(AuthKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		auth, ok := authValue.(AuthContext)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid authentication context"})
			return
		}

		// Fetch user from repository
		usr, err := h.userRepository.FindByID(c.Request.Context(), auth.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		if usr == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// Check if user has any of the allowed roles
		userRole := usr.Role()
		if slices.Contains(allowedRoles, userRole) {
			c.Next()
			return
		}

		// User doesn't have any of the required roles
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
	}
}

// RequireMasterRole is a convenience middleware that checks if the authenticated user has the master role.
func (h *AuthHandler) RequireMasterRole() gin.HandlerFunc {
	return h.RequireRole(user.UserRoleMaster)
}

// RequireAdminRole is a convenience middleware that checks if the authenticated user has the admin role.
func (h *AuthHandler) RequireAdminRole() gin.HandlerFunc {
	return h.RequireRole(user.UserRoleAdmin)
}

// RequirePlayerRole is a convenience middleware that checks if the authenticated user has the player role.
func (h *AuthHandler) RequirePlayerRole() gin.HandlerFunc {
	return h.RequireRole(user.UserRolePlayer)
}
