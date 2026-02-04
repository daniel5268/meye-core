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
	createPJUseCase       campaign.CreatePJUseCase
}

func NewCampaignHandler(createCampaignUseCase campaign.CreateCampaignUseCase, inviteUserUseCase campaign.InviteUserUseCase, createPJUseCase campaign.CreatePJUseCase) *CampaignHandler {
	return &CampaignHandler{
		createCampaignUseCase: createCampaignUseCase,
		inviteUserUseCase:     inviteUserUseCase,
		createPJUseCase:       createPJUseCase,
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

	input := campaign.CreateCampaignInput{
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
	var pathParams dto.CampaignPathParams

	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reqBody dto.InviteUserInputBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := campaign.InviteUserInput{
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

func (h *CampaignHandler) CreatePJ(c *gin.Context) {
	var pathParams dto.CampaignPathParams

	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reqBody dto.CreatePJInputBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	input := campaign.CreatePJInput{
		IDs: campaign.UserCampaignIDs{
			UserID:     auth.UserID,
			CampaignID: pathParams.CampaignID,
		},
		PJInfo: campaign.CreatePJInfo{
			Name:                     reqBody.Name,
			Weight:                   reqBody.Weight,
			Height:                   reqBody.Height,
			Age:                      reqBody.Age,
			Look:                     reqBody.Look,
			Charisma:                 reqBody.Charisma,
			Villainy:                 reqBody.Villainy,
			Heroism:                  reqBody.Heroism,
			PjType:                   reqBody.PjType,
			IsPhysicalTalented:       reqBody.IsPhysicalTalented,
			IsMentalTalented:         reqBody.IsMentalTalented,
			IsCoordinationTalented:   reqBody.IsCoordinationTalented,
			IsPhysicalSkillsTalented: reqBody.IsPhysicalSkillsTalented,
			IsMentalSkillsTalented:   reqBody.IsMentalSkillsTalented,
			IsEnergySkillsTalented:   reqBody.IsEnergySkillsTalented,
			IsEnergyTalented:         reqBody.IsEnergyTalented,
		},
	}

	output, err := h.createPJUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		respondMappedError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.MapPJOutputBody(output))
}
