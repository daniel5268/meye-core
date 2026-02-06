package createcampaign

import (
	"context"
	applicationcampaign "meye-core/internal/application/campaign"
	domaincampaign "meye-core/internal/domain/campaign"
	"meye-core/internal/domain/event"
	"meye-core/internal/domain/shared"
)

// Compile-time check to ensure UseCase implements the port interface
var _ applicationcampaign.CreateCampaignUseCase = (*UseCase)(nil)

type UseCase struct {
	campaignRepository    domaincampaign.Repository
	identificationService shared.IdentificationService
	eventPublisher        event.Publisher
}

func New(campaignRepository domaincampaign.Repository, identificationService shared.IdentificationService, evtPub event.Publisher) *UseCase {
	return &UseCase{
		campaignRepository:    campaignRepository,
		identificationService: identificationService,
		eventPublisher:        evtPub,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input applicationcampaign.CreateCampaignInput) (applicationcampaign.CampaignOutput, error) {
	campaign := domaincampaign.NewCampaign(input.MasterID, input.Name, uc.identificationService)

	if err := uc.campaignRepository.Save(ctx, campaign); err != nil {
		return applicationcampaign.CampaignOutput{}, err
	}

	if err := uc.eventPublisher.Publish(ctx, campaign.UncommittedEvents()); err != nil {
		return applicationcampaign.CampaignOutput{}, err
	}

	return applicationcampaign.MapCampaignOutput(campaign), nil
}
