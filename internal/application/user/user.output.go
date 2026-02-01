package user

import "meye-core/internal/domain/user"

type UserOutput struct {
	ID       string
	Username string
	Role     user.UserRole
}

func MapUserOutput(u *user.User) UserOutput {
	return UserOutput{
		ID:       u.ID(),
		Username: u.Username(),
		Role:     u.Role(),
	}
}
