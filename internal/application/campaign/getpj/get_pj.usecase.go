package getpj

import (
	"context"
	applicationcampaign "meye-core/internal/application/campaign"
	domaincampaign "meye-core/internal/domain/campaign"
)

var _ applicationcampaign.GetPjUseCase = (*UseCase)(nil)

type UseCase struct {
	pjRepository domaincampaign.PjRepository
}

func New(pjRepo domaincampaign.PjRepository) *UseCase {
	return &UseCase{
		pjRepository: pjRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, pjID string) (applicationcampaign.PJOutput, error) {
	pj, err := uc.pjRepository.FindByID(ctx, pjID)
	if err != nil {
		return applicationcampaign.PJOutput{}, err
	}

	if pj == nil {
		return applicationcampaign.PJOutput{}, applicationcampaign.ErrCampaignNotFound
	}

	return applicationcampaign.MapPJOutput(pj), nil
}
