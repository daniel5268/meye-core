package createuser

import (
	"context"
	applicationuser "meye-core/internal/application/user"
	"meye-core/internal/domain/shared"
	"meye-core/internal/domain/user"
	domainuser "meye-core/internal/domain/user"
)

// Compile-time check to ensure UseCase implements the port interface
var _ applicationuser.CreateUserUseCase = (*UseCase)(nil)

type UseCase struct {
	userRepository        domainuser.Repository
	identificationService shared.IdentificationService
	hashService           domainuser.HashService
}

func NewUseCase(usRepo domainuser.Repository, idServ shared.IdentificationService, hashServ domainuser.HashService) *UseCase {
	return &UseCase{
		userRepository:        usRepo,
		identificationService: idServ,
		hashService:           hashServ,
	}
}

func (c *UseCase) Execute(ctx context.Context, input applicationuser.CreateUserInput) (applicationuser.UserOutput, error) {
	newUser, err := user.NewUser(input.Username, input.Password, input.Role, c.identificationService, c.hashService)
	if err != nil {
		return applicationuser.UserOutput{}, err
	}

	existingUser, err := c.userRepository.FindByUsername(ctx, input.Username)
	if err != nil {
		return applicationuser.UserOutput{}, err
	}

	if existingUser != nil {
		return applicationuser.UserOutput{}, applicationuser.ErrUsernameAlreadyExists
	}

	if err := c.userRepository.Save(ctx, newUser); err != nil {
		return applicationuser.UserOutput{}, err
	}

	return applicationuser.MapUserOutput(newUser), nil
}
