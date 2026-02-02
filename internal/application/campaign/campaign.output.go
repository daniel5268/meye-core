package campaign

import "meye-core/internal/domain/campaign"

type CampaignOutput struct {
	ID       string
	Name     string
	MasterID string
}

func MapCampaignOutput(c *campaign.Campaign) CampaignOutput {
	return CampaignOutput{
		ID:       c.ID(),
		Name:     c.Name(),
		MasterID: c.MasterID(),
	}
}
