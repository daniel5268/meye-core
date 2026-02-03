package handler

import (
	"meye-core/internal/application/campaign"
	dto "meye-core/internal/infrastructure/api/handler/dto/campaign"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CampaignHandler struct {
	createCampaignUseCase campaign.CreateCampaignUseCase
	inviteUserUseCase     campaign.InviteUserUseCase
}

func NewCampaignHandler(createCampaignUseCase campaign.CreateCampaignUseCase, inviteUserUseCase campaign.InviteUserUseCase) *CampaignHandler {
	return &CampaignHandler{
		createCampaignUseCase: createCampaignUseCase,
		inviteUserUseCase:     inviteUserUseCase,
	}
}

func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var reqBody dto.CreateCampaignInputBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get auth context set by AuthMiddleware
	authValue, exists := c.Get(AuthKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	auth, ok := authValue.(AuthContext)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid authentication context"})
		return
	}

	input := &campaign.CreateCampaignInput{
		Name:     reqBody.Name,
		MasterID: auth.UserID,
	}

	output, err := h.createCampaignUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		respondMappedError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.MapCampaignOutputBody(output))
}

func (h *CampaignHandler) InviteUser(c *gin.Context) {
	var pathParams dto.InviteUserPathParams

	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reqBody dto.InviteUserInputBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := &campaign.InviteUserInput{
		CampaignID: pathParams.CampaignID,
		UserID:     reqBody.UserID,
	}

	output, err := h.inviteUserUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		respondMappedError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.MapInvitationOutputBody(output))
}
