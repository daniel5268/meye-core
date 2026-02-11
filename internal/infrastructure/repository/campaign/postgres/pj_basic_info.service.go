package postgres

import (
	"context"
	domaincampaign "meye-core/internal/domain/campaign"

	"gorm.io/gorm"
)

var _ domaincampaign.PjQueryService = (*PjQueryService)(nil)

type PjQueryService struct {
	db *gorm.DB
}

func NewPjQueryService(db *gorm.DB) *PjQueryService {
	return &PjQueryService{
		db: db,
	}
}

func (qs *PjQueryService) GetPjsBasicInfo(ctx context.Context, userID string) ([]*domaincampaign.PjBasicInfo, error) {
	var pjs []*PJ

	err := qs.db.WithContext(ctx).
		Select("id", "name").
		Where("user_id = ?", userID).
		Find(&pjs).Error

	if err != nil {
		return nil, err
	}

	result := make([]*domaincampaign.PjBasicInfo, 0, len(pjs))
	for i := range pjs {
		result = append(result, domaincampaign.CreatePjBasicInfo(pjs[i].ID, pjs[i].Name))
	}

	return result, nil
}
