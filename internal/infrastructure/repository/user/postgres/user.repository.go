package postgres

import (
	"context"
	"errors"
	"meye-core/internal/domain/event"
	"meye-core/internal/domain/user"
	"meye-core/internal/infrastructure/repository/shared"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ user.Repository = (*Repository)(nil)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Save performs an upsert operation into DB by ID.
func (r *Repository) Save(ctx context.Context, us *user.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		userModel := GetModelFromDomainUser(us)

		result := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"username":        userModel.Username,
				"hashed_password": userModel.HashedPassword,
				"updated_at":      time.Now(),
			}),
		}).Create(userModel)

		if result.Error != nil {
			return result.Error
		}

		domainEvents := getUncommittedEvents(us)

		return tx.Create(&domainEvents).Error
	})
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	var userModel User
	result := r.db.Where("username = ?", username).First(&userModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}

	domainUser := userModel.ToDomainUser()
	return domainUser, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*user.User, error) {
	var userModel User
	result := r.db.Where("id = ?", id).First(&userModel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}

	domainUser := userModel.ToDomainUser()
	return domainUser, nil
}

func getUncommittedEvents(user *user.User) []shared.DomainEvent {
	events := user.UncommittedEvents()
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

func extractEventData(event event.DomainEvent) shared.EventData {
	data := make(shared.EventData)
	switch e := event.(type) {
	case user.UserCreatedEvent:
		data["role"] = e.Role()
	}
	return data
}
