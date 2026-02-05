package shared

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type DomainEvent struct {
	ID            string
	UserID        string
	Type          string
	AggregateType string
	AggregateID   string
	Data          EventData
	CreatedAt     time.Time
	OccurredAt    time.Time
}

type EventData map[string]interface{}

func (e EventData) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *EventData) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan EventData")
	}
	return json.Unmarshal(bytes, e)
}
