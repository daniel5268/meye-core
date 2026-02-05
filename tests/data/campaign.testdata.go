package data

import (
	"meye-core/internal/domain/campaign"
	"meye-core/tests/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

const (
	CampaignID       = "campaign-id"
	CampaignMasterID = "master-id"
	CampaignName     = "Test Campaign"
	InvitationID     = "invitation-id"
)

func Campaign(t *testing.T) *campaign.Campaign {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idServ := mocks.NewMockIdentificationService(ctrl)
	idServ.EXPECT().GenerateID().Return(InvitationID)

	c := campaign.CreateCampaignWithoutValidation(
		CampaignID,
		CampaignMasterID,
		CampaignName,
		[]campaign.Invitation{},
		[]campaign.PJ{},
	)

	_, _ = c.InviteUser(UserID, idServ)

	return c
}
