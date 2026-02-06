package postgres

import (
	"context"
	"errors"
	"meye-core/internal/domain/campaign"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ campaign.PjRepository = (*PjRepository)(nil)

type PjRepository struct {
	db *gorm.DB
}

func NewPjRepository(db *gorm.DB) *PjRepository {
	return &PjRepository{db: db}
}

func (r *PjRepository) Save(ctx context.Context, pj *campaign.PJ) error {
	model := GetModelFromDomainPJ(pj)

	result := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(model)

	return result.Error
}

func (r *PjRepository) FindByID(ctx context.Context, id string) (*campaign.PJ, error) {
	var pjModel PJ
	result := r.db.Where("id = ?", id).First(&pjModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return pjModel.ToDomain(), nil
}
