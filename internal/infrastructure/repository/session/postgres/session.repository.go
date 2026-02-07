// infrastructure/repository/session/postgres/session.repository.go
package postgres

import (
	"context"
	"meye-core/internal/domain/session"
	"meye-core/internal/infrastructure/repository/shared"

	"gorm.io/gorm"
)

var _ session.Repository = (*Repository)(nil)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, s *session.Session) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		sessionModel := GetModelFromDomainSession(s)

		if err := tx.Create(sessionModel).Error; err != nil {
			return err
		}

		events := getUncommittedEvents(s)

		return tx.Create(&events).Error
	})
}

func getUncommittedEvents(s *session.Session) []shared.DomainEvent {
	events := s.UncommittedEvents()
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
