package getplayers

import (
	"context"
	applicationuser "meye-core/internal/application/user"
	domainuser "meye-core/internal/domain/user"
)

var _ applicationuser.GetPlayersUseCase = (*UseCase)(nil)

type UseCase struct {
	usersRepository domainuser.Repository
}

func New(userRepo domainuser.Repository) *UseCase {
	return &UseCase{
		usersRepository: userRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input applicationuser.PaginationInput) ([]applicationuser.UserOutput, error) {
	players, err := uc.usersRepository.FindByRole(ctx, domainuser.UserRolePlayer, input.Page, input.Size)
	if err != nil {
		return []applicationuser.UserOutput{}, err
	}

	output := make([]applicationuser.UserOutput, 0, len(players))
	for i := range players {
		output = append(output, applicationuser.MapUserOutput(players[i]))
	}

	return output, nil
}
