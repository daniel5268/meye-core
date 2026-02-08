package user

import "meye-core/internal/domain/user"

type CreateUserInput struct {
	Username string
	Password string
	Role     user.UserRole
}

type LoginInput struct {
	Username string
	Password string
}

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

type PaginationInput struct {
	Page int
	Size int
}
