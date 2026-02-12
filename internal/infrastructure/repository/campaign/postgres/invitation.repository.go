package postgres

import (
	"context"
	domaincampaign "meye-core/internal/domain/campaign"

	"gorm.io/gorm"
)

var _ domaincampaign.InvitationRepository = (*InvitationRepository)(nil)

type InvitationRepository struct {
	db *gorm.DB
}

func NewInvitationRepository(db *gorm.DB) *InvitationRepository {
	return &InvitationRepository{
		db: db,
	}
}

func (r *InvitationRepository) FindByUserID(ctx context.Context, userID string) ([]*domaincampaign.Invitation, error) {
	var invitationModels []CampaignInvitation

	result := r.db.Where("user_id = ?", userID).Find(&invitationModels)
	if result.Error != nil {
		return nil, result.Error
	}

	invs := make([]*domaincampaign.Invitation, 0, len(invitationModels))
	for _, invM := range invitationModels {
		invs = append(invs, invM.ToDomain())
	}

	return invs, nil
}
