package getcampaign

import (
	"context"
	applicationcampaign "meye-core/internal/application/campaign"
	domaincampaign "meye-core/internal/domain/campaign"
)

var _ applicationcampaign.GetCampaignUseCase = (*UseCase)(nil)

type UseCase struct {
	campaignRepository domaincampaign.Repository
}

func New(campRepo domaincampaign.Repository) *UseCase {
	return &UseCase{
		campaignRepository: campRepo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, campID string) (applicationcampaign.CampaignOutput, error) {
	camp, err := uc.campaignRepository.FindByID(ctx, campID)
	if err != nil {
		return applicationcampaign.CampaignOutput{}, err
	}

	if camp == nil {
		return applicationcampaign.CampaignOutput{}, applicationcampaign.ErrCampaignNotFound
	}

	return applicationcampaign.MapCampaignOutput(camp), nil
}
