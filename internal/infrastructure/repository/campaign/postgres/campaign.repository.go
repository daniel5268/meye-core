package postgres

import (
	"context"
	"errors"
	"meye-core/internal/domain/campaign"
	"meye-core/internal/domain/event"
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

	return campaignModel.ToDomain(invitationModels, pjModels), nil
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

		// Get current invitations from the domain
		domainInvitations := c.Invitations()
		newInvitationIDs := make(map[string]bool)

		// Insert or update invitations
		for _, domainInvitation := range domainInvitations {
			invitationModel := GetModelFromDomainInvitation(domainInvitation)
			newInvitationIDs[invitationModel.ID] = true

			result := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"state":      invitationModel.State,
					"updated_at": gorm.Expr("CURRENT_TIMESTAMP"),
				}),
			}).Create(invitationModel)

			if result.Error != nil {
				return result.Error
			}
		}

		// Delete invitations that are no longer in the domain array
		if len(newInvitationIDs) > 0 {
			invitationIDsToKeep := make([]string, 0, len(newInvitationIDs))
			for id := range newInvitationIDs {
				invitationIDsToKeep = append(invitationIDsToKeep, id)
			}

			result := tx.Where("campaign_id = ? AND id NOT IN ?", c.ID(), invitationIDsToKeep).
				Delete(&CampaignInvitation{})

			if result.Error != nil {
				return result.Error
			}
		} else {
			// If no invitations in the domain, delete all invitations for this campaign
			result := tx.Where("campaign_id = ?", c.ID()).Delete(&CampaignInvitation{})
			if result.Error != nil {
				return result.Error
			}
		}

		// Get current PJs from the domain
		domainPJs := c.PJs()
		newPjIDs := make(map[string]bool)

		// Insert or update PJs
		for _, domainPJ := range domainPJs {
			pjModel := GetModelFromDomainPJ(domainPJ)
			newPjIDs[pjModel.ID] = true

			result := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"name":                        pjModel.Name,
					"weight":                      pjModel.Weight,
					"height":                      pjModel.Height,
					"age":                         pjModel.Age,
					"look":                        pjModel.Look,
					"charisma":                    pjModel.Charisma,
					"villainy":                    pjModel.Villainy,
					"heroism":                     pjModel.Heroism,
					"pj_type":                     pjModel.PjType,
					"strength":                    pjModel.Strength,
					"agility":                     pjModel.Agility,
					"speed":                       pjModel.Speed,
					"resistance":                  pjModel.Resistance,
					"is_physical_talented":        pjModel.IsPhysicalTalented,
					"inteligence":                 pjModel.Inteligence,
					"wisdom":                      pjModel.Wisdom,
					"concentration":               pjModel.Concentration,
					"will":                        pjModel.Will,
					"is_mental_talented":          pjModel.IsMentalTalented,
					"precision":                   pjModel.Precision,
					"calculation":                 pjModel.Calculation,
					"range":                       pjModel.Range,
					"reflexes":                    pjModel.Reflexes,
					"is_coordination_talented":    pjModel.IsCoordinationTalented,
					"life":                        pjModel.Life,
					"empowerment":                 pjModel.Empowerment,
					"vital_control":               pjModel.VitalControl,
					"is_physical_skills_talented": pjModel.IsPhysicalSkillsTalented,
					"ilusion":                     pjModel.Ilusion,
					"mental_control":              pjModel.MentalControl,
					"is_mental_skills_talented":   pjModel.IsMentalSkillsTalented,
					"object_handling":             pjModel.ObjectHandling,
					"energy_handling":             pjModel.EnergyHandling,
					"energy_tank":                 pjModel.EnergyTank,
					"is_energy_skills_talented":   pjModel.IsEnergySkillsTalented,
					"is_energy_talented":          pjModel.IsEnergyTalented,
					"supernatural_stats":          pjModel.SupernaturalStats,
					"updated_at":                  gorm.Expr("CURRENT_TIMESTAMP"),
				}),
			}).Create(pjModel)

			if result.Error != nil {
				return result.Error
			}
		}

		// Delete PJs that are no longer in the domain array
		if len(newPjIDs) > 0 {
			pjIDsToKeep := make([]string, 0, len(newPjIDs))
			for id := range newPjIDs {
				pjIDsToKeep = append(pjIDsToKeep, id)
			}

			result := tx.Where("campaign_id = ? AND id NOT IN ?", c.ID(), pjIDsToKeep).
				Delete(&PJ{})

			if result.Error != nil {
				return result.Error
			}
		} else {
			// If no PJs in the domain, delete all PJs for this campaign
			result := tx.Where("campaign_id = ?", c.ID()).Delete(&PJ{})
			if result.Error != nil {
				return result.Error
			}
		}

		events := getUncommittedEvents(c)

		return tx.Create(&events).Error

	})
}

func extractEventData(evt event.DomainEvent) shared.EventData {
	data := make(shared.EventData)

	switch e := evt.(type) {
	case campaign.UserInvitedEvent:
		data["campaign_id"] = e.CampaignID()
	case campaign.PjAddedEvent:
		data["campaign_id"] = e.CampaignID()
	}

	return data
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
			Data:          extractEventData(evt),
			CreatedAt:     evt.CreatedAt(),
			OccurredAt:    evt.OccurredAt(),
		}

		domainEvents = append(domainEvents, eventModel)
	}

	return domainEvents
}
