package campaign

import "meye-core/internal/application/campaign"

type CampaignOutputBody struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	MasterID    string                 `json:"master_id"`
	Invitations []InvitationOutputBody `json:"invitations"`
	PJs         []PJOutputBody         `json:"pjs"`
}

func MapCampaignOutputBody(c campaign.CampaignOutput) CampaignOutputBody {
	invitations := make([]InvitationOutputBody, len(c.Invitations))
	for i, inv := range c.Invitations {
		invitations[i] = MapInvitationOutputBody(inv)
	}

	pjs := make([]PJOutputBody, len(c.PJs))
	for i, pj := range c.PJs {
		pjs[i] = MapPJOutputBody(pj)
	}

	return CampaignOutputBody{
		ID:          c.ID,
		Name:        c.Name,
		MasterID:    c.MasterID,
		Invitations: invitations,
		PJs:         pjs,
	}
}
