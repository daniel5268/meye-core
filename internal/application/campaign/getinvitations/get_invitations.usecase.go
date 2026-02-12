package getinvitations

import (
	"context"
	applicationcampaign "meye-core/internal/application/campaign"
	domaincampaign "meye-core/internal/domain/campaign"
)

var _ applicationcampaign.GetInvitationsUseCase = (*UseCase)(nil)

type UseCase struct {
	invitationRepository domaincampaign.InvitationRepository
}

func New(invRepo domaincampaign.InvitationRepository) *UseCase {
	return &UseCase{
		invitationRepository: invRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, userID string) ([]applicationcampaign.InvitationOutput, error) {
	invs, err := uc.invitationRepository.FindByUserID(ctx, userID)
	if err != nil {
		return []applicationcampaign.InvitationOutput{}, err
	}

	output := make([]applicationcampaign.InvitationOutput, 0, len(invs))
	for i := range invs {
		output = append(output, applicationcampaign.MapInvitationOutput(invs[i]))
	}

	return output, nil
}
