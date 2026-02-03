package postgres

import (
	"meye-core/internal/domain/campaign"
	"time"
)

type CampaignInvitation struct {
	ID         string `gorm:"primaryKey"`
	CampaignID string
	UserID     string
	State      string
	CreatedAt  time.Time `gorm:"default:current_timestamp"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp"`
}

func GetModelFromDomainInvitation(i *campaign.Invitation) *CampaignInvitation {
	return &CampaignInvitation{
		ID:         i.ID(),
		CampaignID: i.CampaignID(),
		UserID:     i.UserID(),
		State:      string(i.State()),
	}
}

func (i *CampaignInvitation) ToDomain() *campaign.Invitation {
	return campaign.CreateInvitationWithoutValidation(
		i.ID,
		i.CampaignID,
		i.UserID,
		campaign.InvitationState(i.State),
	)
}
