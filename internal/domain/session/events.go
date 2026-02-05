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
func (e SessionCreatedEvent) OccurredAt() time.Time { return e.occurredAt }

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

type AssignedXP struct {
	basic        uint
	special      uint
	supernatural uint
}

func (a AssignedXP) Basic() uint        { return a.basic }
func (a AssignedXP) Special() uint      { return a.special }
func (a AssignedXP) Supernatural() uint { return a.supernatural }

type XPAssignedEvent struct {
	id         string
	sessionID  string
	userID     string
	pjID       string
	createdAt  time.Time
	occurredAt time.Time
	assignedXP AssignedXP
}

// Compile-time check to ensure XPAssignedEvent implements the DomainEvent interface
var _ event.DomainEvent = (*XPAssignedEvent)(nil)

func (e XPAssignedEvent) ID() string                         { return e.id }
func (e XPAssignedEvent) UserID() string                     { return e.userID }
func (e XPAssignedEvent) AggregateID() string                { return e.pjID }
func (e XPAssignedEvent) CreatedAt() time.Time               { return e.createdAt }
func (e XPAssignedEvent) OccurredAt() time.Time              { return e.occurredAt }
func (e XPAssignedEvent) Type() event.EventType              { return event.EventTypeXPAssigned }
func (e XPAssignedEvent) AggregateType() event.AggregateType { return event.AggregateTypePJ }

func (e XPAssignedEvent) SessionID() string { return e.sessionID }
func (e XPAssignedEvent) AssignedXP() AssignedXP { return e.assignedXP }

func newXPAssignedEvent(xpAssignation XPAssignation, sessionID string, sessionCreatedAt time.Time, masterID string, idServ shared.IdentificationService) XPAssignedEvent {
	return XPAssignedEvent{
		id:         idServ.GenerateID(),
		sessionID:  sessionID,
		userID:     masterID,
		pjID:       xpAssignation.PjID(),
		createdAt:  time.Now(),
		occurredAt: sessionCreatedAt,
		assignedXP: AssignedXP{
			basic:        xpAssignation.Basic(),
			special:      xpAssignation.Special(),
			supernatural: xpAssignation.SuperNatural(),
		},
	}
}
