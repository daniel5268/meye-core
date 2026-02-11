package postgres

import (
	"context"
	domaincampaign "meye-core/internal/domain/campaign"

	"gorm.io/gorm"
)

var _ domaincampaign.CampaignQueryService = (*CampaignQueryService)(nil)

type CampaignQueryService struct {
	db *gorm.DB
}

func NewQueryService(db *gorm.DB) *CampaignQueryService {
	return &CampaignQueryService{
		db: db,
	}
}

func (qd *CampaignQueryService) GetCampaignsBasicInfo(ctx context.Context, masterID string) ([]*domaincampaign.CampaignBasicInfo, error) {
	var campaigns []Campaign

	err := qd.db.WithContext(ctx).
		Select("id", "name", "master_id").
		Where("master_id = ?", masterID).
		Order("created_at DESC").
		Find(&campaigns).Error

	if err != nil {
		return nil, err
	}

	result := make([]*domaincampaign.CampaignBasicInfo, len(campaigns))
	for i, c := range campaigns {
		result[i] = domaincampaign.CreateCampaignBasicInfo(c.ID, c.Name, c.MasterID)
	}

	return result, nil
}
