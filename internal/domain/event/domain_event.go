package event

import "time"

type EventType string

// Session Events
const (
	EventTypeSessionCreated EventType = "session_created"
	EventTypeXPAssigned     EventType = "xp_assigned"
)

type AggregateType string

const (
	AggregateTypeSession AggregateType = "session"
	AggregateTypePJ      AggregateType = "pj"
)

type DomainEvent interface {
	ID() string
	UserID() string
	Type() EventType
	AggregateID() string
	AggregateType() AggregateType
	CreatedAt() time.Time
	OccurredAt() time.Time
}
