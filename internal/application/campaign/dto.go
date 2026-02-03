package campaign

import "meye-core/internal/domain/campaign"

type CreateCampaignInput struct {
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

type CampaignOutput struct {
	ID       string
	Name     string
	MasterID string
}

type InviteUserInput struct {
	CampaignID string
	UserID     string
}

func MapInvitationOutput(i *campaign.Invitation) InvitationOutput {
	return InvitationOutput{
		ID:         i.ID(),
		CampaignID: i.CampaignID(),
		UserID:     i.UserID(),
		State:      i.State(),
	}
}

type InvitationOutput struct {
	ID         string
	CampaignID string
	UserID     string
	State      campaign.InvitationState
}
