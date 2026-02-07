package handler

import (
	"meye-core/internal/domain/campaign"
	"meye-core/internal/domain/user"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	apiKey             string
	jwtService         user.JWTService
	userRepository     user.Repository
	campaignRepository campaign.Repository
	pjRepository       campaign.PjRepository
}

type responseError struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

var unauthorizedError = responseError{Error: "Unauthorized"}
var forbiddenError = responseError{Error: "Forbidden"}

const AuthKey = "auth"

func NewAuthHandler(
	apiKey string,
	jwtService user.JWTService,
	userRepo user.Repository,
	campaignRepo campaign.Repository,
	pjRepo campaign.PjRepository,
) *AuthHandler {
	return &AuthHandler{
		apiKey:             apiKey,
		jwtService:         jwtService,
		userRepository:     userRepo,
		campaignRepository: campaignRepo,
		pjRepository:       pjRepo,
	}
}

func (h *AuthHandler) InternalAPIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("Api-Key")

		if apiKey != h.apiKey {
			c.JSON(http.StatusUnauthorized, unauthorizedError)
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedError)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedError)
			return
		}

		tokenString := parts[1]
		userID, err := h.jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedError)
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedError)
			return
		}

		auth, ok := authValue.(AuthContext)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedError)
			return
		}

		// Fetch user from repository
		usr, err := h.userRepository.FindByID(c.Request.Context(), auth.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedError)
			return
		}

		if usr == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedError)
			return
		}

		// Check if user has any of the allowed roles
		userRole := usr.Role()
		if slices.Contains(allowedRoles, userRole) {
			c.Next()
			return
		}

		// User doesn't have any of the required roles
		c.AbortWithStatusJSON(http.StatusForbidden, forbiddenError)
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

// RequireCampaignMaster is a middleware that checks if the authenticated user is the master of the campaign.
// This middleware should be used after AuthMiddleware, as it depends on the AuthContext being set.
// It expects a campaignID parameter in the URI path.
func (h *AuthHandler) RequireCampaignMaster() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get auth context set by AuthMiddleware
		authValue, exists := c.Get(AuthKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedError)
			return
		}

		auth, ok := authValue.(AuthContext)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedError)
			return
		}

		// Get campaign ID from URI parameter
		campaignID := c.Param("campaignID")
		if campaignID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, responseError{Error: "Parameter campaignID is required", Code: "MISSING_CAMPAIGN_ID"})
			return
		}

		// Fetch campaign from repository
		cmp, err := h.campaignRepository.FindByID(c.Request.Context(), campaignID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, responseError{Error: "Failed to retrieve campaign", Code: "FAILED_TO_RETRIEVE_CAMPAIGN"})
			return
		}

		if cmp == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, responseError{Error: "Campaign not found", Code: "CAMPAIGN_NOT_FOUND"})
			return
		}

		// Check if authenticated user is the campaign master
		if cmp.MasterID() != auth.UserID {
			c.AbortWithStatusJSON(http.StatusForbidden, forbiddenError)
			return
		}

		c.Next()
	}
}

func (h *AuthHandler) RequirePjUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get auth context set by AuthMiddleware
		authValue, exists := c.Get(AuthKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedError)
			return
		}

		auth, ok := authValue.(AuthContext)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedError)
			return
		}

		pjID := c.Param("pjID")
		if pjID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, responseError{Error: "Parameter pjID is required", Code: "MISSING_PJ_ID"})
			return
		}

		pj, err := h.pjRepository.FindByID(c.Request.Context(), pjID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, responseError{Error: "Failed to retrieve pj", Code: "FAILED_TO_RETRIEVE_CAMPAIGN"})
			return
		}

		if pj == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, responseError{Error: "PJ not found", Code: "PJ_NOT_FOUND"})
			return
		}

		if pj.UserID() != auth.UserID {
			c.AbortWithStatusJSON(http.StatusForbidden, forbiddenError)
			return
		}

		c.Next()
	}
}
