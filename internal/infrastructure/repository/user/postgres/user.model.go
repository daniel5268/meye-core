package postgres

import (
	"meye-core/internal/domain/user"
	"time"
)

type User struct {
	Id             string `gorm:"primaryKey"`
	Username       string
	HashedPassword string
	Role           user.UserRole
	CreatedAt      time.Time `gorm:"default:current_timestamp"`
	UpdatedAt      time.Time `gorm:"default:current_timestamp"`
}

func GetModelFromDomainUser(u *user.User) *User {
	return &User{
		Id:             u.ID(),
		Username:       u.Username(),
		HashedPassword: u.HashedPassword(),
		Role:           u.Role(),
	}
}

func (u *User) ToDomainUser() *user.User {
	return user.CreateUserWithoutValidation(
		u.Id,
		u.Username,
		u.HashedPassword,
		u.Role,
	)
}
