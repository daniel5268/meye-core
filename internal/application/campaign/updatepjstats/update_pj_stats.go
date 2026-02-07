package updatepjstats

import (
	"context"
	applicationcampaign "meye-core/internal/application/campaign"
	domaincampaign "meye-core/internal/domain/campaign"
)

var _ applicationcampaign.UpdateStatsUseCase = (*UseCase)(nil)

type UseCase struct {
	pjRepository domaincampaign.PjRepository
}

func New(pjRepository domaincampaign.PjRepository) *UseCase {
	return &UseCase{
		pjRepository: pjRepository,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input applicationcampaign.UpdatePjStatsInput) (applicationcampaign.PJOutput, error) {
	pj, err := uc.pjRepository.FindByID(ctx, input.PjID)
	if err != nil {
		return applicationcampaign.PJOutput{}, err
	}

	if pj == nil {
		return applicationcampaign.PJOutput{}, domaincampaign.ErrPjNotFound
	}

	updateParams := applicationcampaign.MapToUpdatePjStatsParameters(input)

	err = pj.UpdateStats(updateParams)
	if err != nil {
		return applicationcampaign.PJOutput{}, err
	}

	err = uc.pjRepository.Save(ctx, pj)
	if err != nil {
		return applicationcampaign.PJOutput{}, err
	}

	output := applicationcampaign.MapPJOutput(pj)

	return output, nil
}
