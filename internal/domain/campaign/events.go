package campaign

import (
	"meye-core/internal/domain/event"
	"meye-core/internal/domain/shared"
	"time"
)

var _ event.DomainEvent = (*CampaignCreatedEvent)(nil)

type CampaignCreatedEvent struct {
	id         string
	campaignID string
	createdAt  time.Time
	occurredAt time.Time
}

func (e CampaignCreatedEvent) ID() string                         { return e.id }
func (e CampaignCreatedEvent) Type() event.EventType              { return event.EventTypeCampaignCreated }
func (e CampaignCreatedEvent) AggregateID() string                { return e.campaignID }
func (e CampaignCreatedEvent) AggregateType() event.AggregateType { return event.AggregateTypeCampaign }
func (e CampaignCreatedEvent) CreatedAt() time.Time               { return e.createdAt }
func (e CampaignCreatedEvent) OccurredAt() time.Time              { return e.occurredAt }

func newCampaignCreatedEvent(c *Campaign, idServ shared.IdentificationService) CampaignCreatedEvent {
	return CampaignCreatedEvent{
		id:         idServ.GenerateID(),
		campaignID: c.id,
		createdAt:  time.Now(),
		occurredAt: time.Now(),
	}
}

var _ event.DomainEvent = (*UserInvitedEvent)(nil)

type UserInvitedEvent struct {
	id         string
	campaignID string
	userID     string
	createdAt  time.Time
	occurredAt time.Time
}

func (e UserInvitedEvent) ID() string                         { return e.id }
func (e UserInvitedEvent) Type() event.EventType              { return event.EventTypeUserInvited }
func (e UserInvitedEvent) AggregateID() string                { return e.userID }
func (e UserInvitedEvent) AggregateType() event.AggregateType { return event.AggregateTypeUser }
func (e UserInvitedEvent) CreatedAt() time.Time               { return e.createdAt }
func (e UserInvitedEvent) OccurredAt() time.Time              { return e.occurredAt }

func (e UserInvitedEvent) CampaignID() string { return e.campaignID }

func newUserInvitedEvent(userID, campaignID string, idServ shared.IdentificationService) UserInvitedEvent {
	return UserInvitedEvent{
		id:         idServ.GenerateID(),
		campaignID: campaignID,
		userID:     userID,
		createdAt:  time.Now(),
		occurredAt: time.Now(),
	}
}

var _ event.DomainEvent = (*PjAddedEvent)(nil)

type PjAddedEvent struct {
	id         string
	campaignID string
	pjID       string
	createdAt  time.Time
	occurredAt time.Time
}

func (e PjAddedEvent) ID() string                         { return e.id }
func (e PjAddedEvent) Type() event.EventType              { return event.EventTypePjAdded }
func (e PjAddedEvent) AggregateID() string                { return e.pjID }
func (e PjAddedEvent) AggregateType() event.AggregateType { return event.AggregateTypePJ }
func (e PjAddedEvent) CreatedAt() time.Time               { return e.createdAt }
func (e PjAddedEvent) OccurredAt() time.Time              { return e.occurredAt }

func (e PjAddedEvent) CampaignID() string { return e.campaignID }

func newPjAddedEvent(pjID, campaignID string, idServ shared.IdentificationService) PjAddedEvent {
	return PjAddedEvent{
		id:         idServ.GenerateID(),
		campaignID: campaignID,
		pjID:       pjID,
		createdAt:  time.Now(),
		occurredAt: time.Now(),
	}
}
