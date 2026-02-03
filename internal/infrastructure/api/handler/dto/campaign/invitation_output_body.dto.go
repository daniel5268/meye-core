package campaign

import "meye-core/internal/application/campaign"

type InvitationOutputBody struct {
	ID         string `json:"id"`
	CampaignID string `json:"campaign_id"`
	UserID     string `json:"user_id"`
	State      string `json:"state"`
}

func MapInvitationOutputBody(i campaign.InvitationOutput) InvitationOutputBody {
	return InvitationOutputBody{
		ID:         i.ID,
		CampaignID: i.CampaignID,
		UserID:     i.UserID,
		State:      string(i.State),
	}
}
