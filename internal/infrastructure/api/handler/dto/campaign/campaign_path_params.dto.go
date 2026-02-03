package campaign

type CampaignPathParams struct {
	CampaignID string `uri:"campaignID" binding:"required"`
}
