package getuser

import (
	"context"
	applicationuser "meye-core/internal/application/user"
	domainuser "meye-core/internal/domain/user"
)

var _ applicationuser.GetUserUseCase = (*UseCase)(nil)

type UseCase struct {
	userRepository domainuser.Repository
}

func New(userRepo domainuser.Repository) *UseCase {
	return &UseCase{
		userRepository: userRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, id string) (applicationuser.UserOutput, error) {
	u, err := uc.userRepository.FindByID(ctx, id)
	if err != nil {
		return applicationuser.UserOutput{}, err
	}

	return applicationuser.MapUserOutput(u), nil
}
