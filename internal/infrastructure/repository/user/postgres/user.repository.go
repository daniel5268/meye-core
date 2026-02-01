package postgres

import (
	"context"
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
			"role":            userModel.Role,
			"updated_at":      time.Now(),
		}),
	}).Create(userModel)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
