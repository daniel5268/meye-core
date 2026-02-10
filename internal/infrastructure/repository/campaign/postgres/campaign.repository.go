package postgres

import (
	"context"
	"errors"
	"meye-core/internal/domain/campaign"
	"meye-core/internal/infrastructure/repository/shared"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ campaign.Repository = (*Repository)(nil)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByID(ctx context.Context, id string) (*campaign.Campaign, error) {
	var campaignModel Campaign
	result := r.db.Where("id = ?", id).First(&campaignModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}

	var invitationModels []CampaignInvitation
	result = r.db.Where("campaign_id = ?", id).Find(&invitationModels)
	if result.Error != nil {
		return nil, result.Error
	}

	var pjModels []PJ
	result = r.db.Where("campaign_id = ?", id).Find(&pjModels)
	if result.Error != nil {
		return nil, result.Error
	}

	var sessionModels []Session
	result = r.db.Where("campaign_id = ?", id).
		Order("created_at DESC").
		Find(&sessionModels)
	if result.Error != nil {
		return nil, result.Error
	}

	return campaignModel.ToDomain(invitationModels, pjModels, sessionModels), nil
}

func (r *Repository) Save(ctx context.Context, c *campaign.Campaign) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Save or update campaign
		campaignModel := GetModelFromDomainCampaign(c)

		result := tx.Clauses(clause.OnConflict{
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

		events := getUncommittedEvents(c)

		return tx.Create(&events).Error

	})
}

func getUncommittedEvents(c *campaign.Campaign) []shared.DomainEvent {
	events := c.UncommittedEvents()
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
