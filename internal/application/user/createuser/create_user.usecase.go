package createuser

import (
	"context"
	applicationuser "meye-core/internal/application/user"
	"meye-core/internal/domain/shared"
	"meye-core/internal/domain/user"
	domainuser "meye-core/internal/domain/user"
)

type UseCase struct {
	identificationService shared.IdentificationService
	userRepository        domainuser.Repository
	hashService           domainuser.HashService
}

func NewUseCase(usRepo domainuser.Repository, idServ shared.IdentificationService, hashServ domainuser.HashService) *UseCase {
	return &UseCase{
		identificationService: idServ,
		userRepository:        usRepo,
		hashService:           hashServ,
	}
}

func (c *UseCase) Execute(ctx context.Context, input Input) (applicationuser.UserOutput, error) {
	newUser, err := user.NewUser(input.Username, input.Password, input.Role, c.identificationService, c.hashService)
	if err != nil {
		return applicationuser.UserOutput{}, err
	}

	if err := c.userRepository.Save(ctx, newUser); err != nil {
		return applicationuser.UserOutput{}, err
	}

	return applicationuser.MapUserOutput(newUser), nil
}
