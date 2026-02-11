package getpjs

import (
	"context"
	applicationcampaign "meye-core/internal/application/campaign"
	domaincampaign "meye-core/internal/domain/campaign"
)

var _ applicationcampaign.GetPjsUseCase = (*UseCase)(nil)

type UseCase struct {
	pjQueryService domaincampaign.PjQueryService
}

func New(pjQueryService domaincampaign.PjQueryService) *UseCase {
	return &UseCase{
		pjQueryService: pjQueryService,
	}
}

func (uc *UseCase) Execute(ctx context.Context, userID string) ([]applicationcampaign.PjBasicInfoOutput, error) {
	pjs, err := uc.pjQueryService.GetPjsBasicInfo(ctx, userID)
	if err != nil {
		return []applicationcampaign.PjBasicInfoOutput{}, err
	}

	output := make([]applicationcampaign.PjBasicInfoOutput, 0, len(pjs))
	for i := range pjs {
		output = append(output, applicationcampaign.MapPjBasicInfo(pjs[i]))
	}

	return output, nil
}
