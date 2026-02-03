package user

import (
	"context"
)

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
