package postgres

import (
	"context"
	"meye-core/internal/domain/campaign"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, c *campaign.Campaign) error {
	campaignModel := GetModelFromDomainCampaign(c)

	result := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"name":       campaignModel.Name,
			"master_id":  campaignModel.MasterID,
			"updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}),
	}).Create(campaignModel)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
