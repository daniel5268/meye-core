package campaign

type CreateCampaignInputBody struct {
	Name string `json:"name" binding:"required"`
}
