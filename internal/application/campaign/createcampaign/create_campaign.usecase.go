package createcampaign

import (
	"context"
	aplicationcampaign "meye-core/internal/application/campaign"
	domaincampaign "meye-core/internal/domain/campaign"
	"meye-core/internal/domain/shared"
)

// Compile-time check to ensure UseCase implements the port interface
var _ aplicationcampaign.CreateCampaignUseCase = (*UseCase)(nil)

type UseCase struct {
	campaignRepository    domaincampaign.Repository
	identificationService shared.IdentificationService
}

func NewUseCase(campaignRepository domaincampaign.Repository, identificationService shared.IdentificationService) *UseCase {
	return &UseCase{
		campaignRepository:    campaignRepository,
		identificationService: identificationService,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input *aplicationcampaign.CreateCampaignInput) (aplicationcampaign.CampaignOutput, error) {
	campaign := domaincampaign.NewCampaign(input.MasterID, input.Name, uc.identificationService)

	if err := uc.campaignRepository.Save(ctx, campaign); err != nil {
		return aplicationcampaign.CampaignOutput{}, err
	}

	return aplicationcampaign.MapCampaignOutput(campaign), nil
}
