package campaign

import "meye-core/internal/application/campaign"

type CampaignOutputBody struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	MasterID string `json:"master_id"`
}

func MapCampaignOutputBody(c campaign.CampaignOutput) CampaignOutputBody {
	return CampaignOutputBody{
		ID:       c.ID,
		Name:     c.Name,
		MasterID: c.MasterID,
	}
}
