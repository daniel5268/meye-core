package session

import (
	"meye-core/internal/domain/event"
	"meye-core/internal/domain/shared"
	"time"
)

type SessionCreatedEvent struct {
	id         string
	sessionID  string
	userID     string
	createdAt  time.Time
	occurredAt time.Time

	campaignID string
}

// Compile-time check to ensure UseCase implements the port interface
var _ event.DomainEvent = (*SessionCreatedEvent)(nil)

func (e SessionCreatedEvent) ID() string            { return e.id }
func (e SessionCreatedEvent) UserID() string        { return e.userID }
func (e SessionCreatedEvent) AggregateID() string   { return e.sessionID }
func (e SessionCreatedEvent) CreatedAt() time.Time  { return e.createdAt }
func (e SessionCreatedEvent) OccurredAt() time.Time { return e.createdAt }

func (e SessionCreatedEvent) Type() event.EventType              { return event.EventTypeSessionCreated }
func (e SessionCreatedEvent) AggregateType() event.AggregateType { return event.AggregateTypeSession }

func (e SessionCreatedEvent) CampaignID() string { return e.campaignID }

func newSessionCreatedEvent(s *Session, masterID string, idServ shared.IdentificationService) SessionCreatedEvent {
	return SessionCreatedEvent{
		id:         idServ.GenerateID(),
		userID:     masterID,
		sessionID:  s.id,
		createdAt:  time.Now(),
		occurredAt: s.createdAt,
		campaignID: s.campaignID,
	}
}
