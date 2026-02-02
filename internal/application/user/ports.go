package user

import (
	"context"
	domainuser "meye-core/internal/domain/user"
)

// CreateUserInput represents the input data required to create a new user
type CreateUserInput struct {
	Username string
	Password string
	Role     domainuser.UserRole
}

// LoginInput represents the input data required for user authentication
type LoginInput struct {
	Username string
	Password string
}

// CreateUserUseCase defines the port for creating a new user.
// This interface allows the infrastructure layer to depend on an abstraction
// rather than concrete use case implementations.
type CreateUserUseCase interface {
	Execute(ctx context.Context, input CreateUserInput) (UserOutput, error)
}

// LoginUseCase defines the port for user authentication.
// This interface allows the infrastructure layer to depend on an abstraction
// rather than concrete use case implementations.
type LoginUseCase interface {
	Execute(ctx context.Context, input LoginInput) (string, error)
}