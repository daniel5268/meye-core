package inviteuser

import (
	"context"
	applicationcampaign "meye-core/internal/application/campaign"
	applicationuser "meye-core/internal/application/user"
	domaincampaign "meye-core/internal/domain/campaign"
	"meye-core/internal/domain/shared"
	domainuser "meye-core/internal/domain/user"
)

var _ applicationcampaign.InviteUserUseCase = (*UseCase)(nil)

type UseCase struct {
	campainRepository     domaincampaign.Repository
	userRepository        domainuser.Repository
	identificationService shared.IdentificationService
}

func New(campainRepository domaincampaign.Repository, userRepository domainuser.Repository, identificationService shared.IdentificationService) *UseCase {
	return &UseCase{
		campainRepository:     campainRepository,
		userRepository:        userRepository,
		identificationService: identificationService,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input applicationcampaign.InviteUserInput) (applicationcampaign.InvitationOutput, error) {
	cmp, err := uc.campainRepository.FindByID(ctx, input.CampaignID)
	if err != nil {
		return applicationcampaign.InvitationOutput{}, err
	}

	if cmp == nil {
		return applicationcampaign.InvitationOutput{}, applicationcampaign.ErrCampaignNotFound
	}

	user, err := uc.userRepository.FindByID(ctx, input.UserID)
	if err != nil {
		return applicationcampaign.InvitationOutput{}, err
	}

	if user == nil {
		return applicationcampaign.InvitationOutput{}, applicationuser.ErrUserNotFound
	}

	err = user.MustBePlayer()
	if err != nil {
		return applicationcampaign.InvitationOutput{}, err
	}

	inv, err := cmp.InviteUser(user.ID(), uc.identificationService)
	if err != nil {
		return applicationcampaign.InvitationOutput{}, err
	}

	if err := uc.campainRepository.Save(ctx, cmp); err != nil {
		return applicationcampaign.InvitationOutput{}, err
	}

	return applicationcampaign.MapInvitationOutput(inv), nil
}
