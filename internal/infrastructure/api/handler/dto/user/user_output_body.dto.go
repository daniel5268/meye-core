package dto

import (
	applicationuser "meye-core/internal/application/user"
	domainuser "meye-core/internal/domain/user"
)

type UserOutputBody struct {
	ID       string              `json:"id"`
	Username string              `json:"username"`
	Role     domainuser.UserRole `json:"role"`
}

func MapUserOutput(u applicationuser.UserOutput) UserOutputBody {
	return UserOutputBody{
		ID:       u.ID,
		Username: u.Username,
		Role:     u.Role,
	}
}
