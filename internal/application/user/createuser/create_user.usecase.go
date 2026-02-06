package createuser

import (
	"context"
	applicationuser "meye-core/internal/application/user"
	"meye-core/internal/domain/event"
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
	eventPublisher        event.Publisher
}

func New(
	usRepo domainuser.Repository,
	idServ shared.IdentificationService,
	hashServ domainuser.HashService,
	eventPub event.Publisher,
) *UseCase {
	return &UseCase{
		userRepository:        usRepo,
		identificationService: idServ,
		hashService:           hashServ,
		eventPublisher:        eventPub,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input applicationuser.CreateUserInput) (applicationuser.UserOutput, error) {
	existingUser, err := uc.userRepository.FindByUsername(ctx, input.Username)
	if err != nil {
		return applicationuser.UserOutput{}, err
	}

	if existingUser != nil {
		return applicationuser.UserOutput{}, applicationuser.ErrUsernameAlreadyExists
	}

	newUser, err := user.NewUser(input.Username, input.Password, input.Role, uc.identificationService, uc.hashService)
	if err != nil {
		return applicationuser.UserOutput{}, err
	}

	if err := uc.userRepository.Save(ctx, newUser); err != nil {
		return applicationuser.UserOutput{}, err
	}

	if err := uc.eventPublisher.Publish(ctx, newUser.UncommittedEvents()); err != nil {
		return applicationuser.UserOutput{}, err
	}

	return applicationuser.MapUserOutput(newUser), nil
}
