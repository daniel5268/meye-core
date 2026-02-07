package postgres

import (
	"context"
	"errors"
	"meye-core/internal/domain/campaign"
	"meye-core/internal/infrastructure/repository/shared"

	"gorm.io/gorm"
)

var _ campaign.PjRepository = (*PjRepository)(nil)

type PjRepository struct {
	db *gorm.DB
}

func NewPjRepository(db *gorm.DB) *PjRepository {
	return &PjRepository{db: db}
}

func (r *PjRepository) Save(ctx context.Context, pj *campaign.PJ) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		model := GetModelFromDomainPJ(pj)

		if err := tx.WithContext(ctx).Save(model).Error; err != nil {
			return err
		}

		events := getPjUncommittedEvents(pj)

		return tx.Create(&events).Error
	})
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

func getPjUncommittedEvents(pj *campaign.PJ) []shared.DomainEvent {
	events := pj.UncommittedEvents()
	domainEvents := make([]shared.DomainEvent, 0, len(events))
	for _, evt := range events {
		eventModel := shared.DomainEvent{
			ID:            evt.ID(),
			Type:          string(evt.Type()),
			AggregateType: string(evt.AggregateType()),
			AggregateID:   evt.AggregateID(),
			Data:          evt.GetSerializedData(),
			CreatedAt:     evt.CreatedAt(),
			OccurredAt:    evt.OccurredAt(),
		}

		domainEvents = append(domainEvents, eventModel)
	}

	return domainEvents
}
