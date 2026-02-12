package handler

import (
	"meye-core/internal/application/campaign"
	"meye-core/internal/application/session"
	dto "meye-core/internal/infrastructure/api/handler/dto/campaign"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CampaignHandler struct {
	createCampaignUseCase campaign.CreateCampaignUseCase
	inviteUserUseCase     campaign.InviteUserUseCase
	createPJUseCase       campaign.CreatePJUseCase
	createSessionUseCase  session.CreateSessionUseCase
	updatePJStatsUseCase  campaign.UpdateStatsUseCase
	getCampaignUseCase    campaign.GetCampaignUseCase
	getPjUseCase          campaign.GetPjUseCase
	getCampaignsUseCase   campaign.GetCampaignsUseCase
	getPjsUseCase         campaign.GetPjsUseCase
	getInvitations        campaign.GetInvitationsUseCase
}

func NewCampaignHandler(
	createCampaignUseCase campaign.CreateCampaignUseCase,
	inviteUserUseCase campaign.InviteUserUseCase,
	createPJUseCase campaign.CreatePJUseCase,
	createSessionUseCase session.CreateSessionUseCase,
	updatePJStatsUseCase campaign.UpdateStatsUseCase,
	getcampaignUseCase campaign.GetCampaignUseCase,
	getPjUseCase campaign.GetPjUseCase,
	getCampaignsUseCase campaign.GetCampaignsUseCase,
	getPjsUseCase campaign.GetPjsUseCase,
	getInvitations campaign.GetInvitationsUseCase,
) *CampaignHandler {
	return &CampaignHandler{
		createCampaignUseCase: createCampaignUseCase,
		inviteUserUseCase:     inviteUserUseCase,
		createPJUseCase:       createPJUseCase,
		createSessionUseCase:  createSessionUseCase,
		updatePJStatsUseCase:  updatePJStatsUseCase,
		getCampaignUseCase:    getcampaignUseCase,
		getPjUseCase:          getPjUseCase,
		getCampaignsUseCase:   getCampaignsUseCase,
		getPjsUseCase:         getPjsUseCase,
		getInvitations:        getInvitations,
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

func (h *CampaignHandler) CreateSession(c *gin.Context) {
	var pathParams dto.CampaignPathParams

	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reqBody dto.CreateSessionInputBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	xpAss := make([]session.XPAssignation, 0, len(reqBody.XPAssignations))
	for _, xpA := range reqBody.XPAssignations {
		xpAss = append(xpAss, session.XPAssignation{
			PjID: xpA.PjID,
			Amounts: session.XPAmounts{
				Basic:        xpA.Amounts.Basic,
				Special:      xpA.Amounts.Special,
				SuperNatural: xpA.Amounts.SuperNatural,
			},
			Reason: xpA.Reason,
		})
	}

	input := session.CreateSessionInput{
		CampaignID:     pathParams.CampaignID,
		Summary:        reqBody.Summary,
		XPAssignations: xpAss,
	}

	output, err := h.createSessionUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		respondMappedError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.MapSessionOutput(output))
}

func (h *CampaignHandler) UpdatePJStats(c *gin.Context) {
	var pathParams dto.PJPathParams

	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reqBody dto.UpdatePJStatsInputBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := dto.MapUpdatePJStatsInput(pathParams, reqBody)

	output, err := h.updatePJStatsUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		respondMappedError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MapPJOutputBody(output))
}

func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	var pathParams dto.CampaignPathParams

	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.getCampaignUseCase.Execute(c.Request.Context(), pathParams.CampaignID)
	if err != nil {
		respondMappedError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MapCampaignOutputBody(output))
}

func (h *CampaignHandler) GetPj(c *gin.Context) {
	var pathParams dto.PJPathParams

	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.getPjUseCase.Execute(c.Request.Context(), pathParams.PJID)
	if err != nil {
		respondMappedError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MapPJOutputBody(output))
}

func (h *CampaignHandler) GetCampaignsBasicInfo(c *gin.Context) {
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

	output, err := h.getCampaignsUseCase.Execute(c.Request.Context(), auth.UserID)
	if err != nil {
		respondMappedError(c, err)
		return
	}

	outputBody := make([]dto.CampaignBasicInfoOutputBody, 0, len(output))
	for _, o := range output {
		outputBody = append(outputBody, dto.MapCampaignBasicInfoOutputBody(o))
	}

	c.JSON(http.StatusOK, outputBody)
}

func (h *CampaignHandler) GetPjs(c *gin.Context) {
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

	output, err := h.getPjsUseCase.Execute(c.Request.Context(), auth.UserID)
	if err != nil {
		respondMappedError(c, err)
		return
	}

	outputBody := make([]dto.PjBasicInfoOutputBody, 0, len(output))
	for _, o := range output {
		outputBody = append(outputBody, dto.MapPjBasicInfoOutputBody(o))
	}

	c.JSON(http.StatusOK, outputBody)
}

func (h *CampaignHandler) GetUserInvitations(c *gin.Context) {
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

	result, err := h.getInvitations.Execute(c.Request.Context(), auth.UserID)
	if err != nil {
		respondMappedError(c, err)
		return
	}

	output := make([]dto.InvitationOutputBody, 0, len(result))
	for _, o := range result {
		output = append(output, dto.MapInvitationOutputBody(o))
	}

	c.JSON(http.StatusOK, output)
}
