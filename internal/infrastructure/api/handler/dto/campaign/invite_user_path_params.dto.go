package campaign

type InviteUserPathParams struct {
	CampaignID string `uri:"campaignID" binding:"required"`
}