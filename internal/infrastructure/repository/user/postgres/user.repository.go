package postgres

import (
	"context"
	"errors"
	"meye-core/internal/domain/user"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Save performs an upsert operation into DB by ID.
func (r *Repository) Save(ctx context.Context, user *user.User) error {
	userModel := GetModelFromDomainUser(user)

	result := r.db.Clauses(clause.OnConflict{
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

	return nil
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
