package postgres

import (
	"meye-core/internal/domain/campaign"
	"time"
)

type Campaign struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	MasterID  string
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
}

func GetModelFromDomainCampaign(c *campaign.Campaign) *Campaign {
	return &Campaign{
		ID:       c.ID(),
		Name:     c.Name(),
		MasterID: c.MasterID(),
	}
}

func (c *Campaign) ToDomain(invitations []CampaignInvitation, pjs []PJ) *campaign.Campaign {
	domainInvitations := make([]campaign.Invitation, 0, len(invitations))
	for _, inv := range invitations {
		domainInvitations = append(domainInvitations, *inv.ToDomain())
	}

	domainPJs := make([]campaign.PJ, 0, len(pjs))
	for _, pj := range pjs {
		domainPJs = append(domainPJs, *pj.ToDomain())
	}

	return campaign.CreateCampaignWithoutValidation(
		c.ID,
		c.MasterID,
		c.Name,
		domainInvitations,
		domainPJs,
	)
}
