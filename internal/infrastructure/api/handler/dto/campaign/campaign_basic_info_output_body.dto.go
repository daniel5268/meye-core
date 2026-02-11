package campaign

import "meye-core/internal/application/campaign"

type CampaignBasicInfoOutputBody struct {
	ID       string `json:"id"`
	MasterID string `json:"master_id"`
	Name     string `json:"name"`
}

func MapCampaignBasicInfoOutputBody(c campaign.CampaignBasicInfoOutput) CampaignBasicInfoOutputBody {
	return CampaignBasicInfoOutputBody{
		ID:       c.ID,
		MasterID: c.MasterID,
		Name:     c.Name,
	}
}
