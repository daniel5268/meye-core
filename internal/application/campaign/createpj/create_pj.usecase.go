package createpj

import (
	"context"
	applicationcampaign "meye-core/internal/application/campaign"
	applicationuser "meye-core/internal/application/user"
	domaincampaign "meye-core/internal/domain/campaign"
	"meye-core/internal/domain/event"
	"meye-core/internal/domain/shared"
	domainuser "meye-core/internal/domain/user"
)

// Compile-time check to ensure UseCase implements the port interface
var _ applicationcampaign.CreatePJUseCase = (*UseCase)(nil)

type UseCase struct {
	campaignRepository    domaincampaign.Repository
	userRepository        domainuser.Repository
	identificationService shared.IdentificationService
	eventPublisher        event.Publisher
}

func New(
	campRepo domaincampaign.Repository,
	userRepo domainuser.Repository,
	idServ shared.IdentificationService,
	evtPub event.Publisher,
) *UseCase {
	return &UseCase{
		campaignRepository:    campRepo,
		userRepository:        userRepo,
		identificationService: idServ,
		eventPublisher:        evtPub,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input applicationcampaign.CreatePJInput) (applicationcampaign.PJOutput, error) {
	user, err := uc.userRepository.FindByID(ctx, input.IDs.UserID)
	if err != nil {
		return applicationcampaign.PJOutput{}, err
	}

	if user == nil {
		return applicationcampaign.PJOutput{}, applicationuser.ErrUserNotFound
	}

	camp, err := uc.campaignRepository.FindByID(ctx, input.IDs.CampaignID)
	if err != nil {
		return applicationcampaign.PJOutput{}, err
	}

	if camp == nil {
		return applicationcampaign.PJOutput{}, applicationcampaign.ErrCampaignNotFound
	}

	params := domaincampaign.PJCreateParameters{
		Name:                     input.PJInfo.Name,
		Weight:                   input.PJInfo.Weight,
		Height:                   input.PJInfo.Height,
		Age:                      input.PJInfo.Age,
		Look:                     input.PJInfo.Look,
		Charisma:                 input.PJInfo.Charisma,
		Villainy:                 input.PJInfo.Villainy,
		Heroism:                  input.PJInfo.Heroism,
		PjType:                   input.PJInfo.PjType,
		IsPhysicalTalented:       input.PJInfo.IsPhysicalTalented,
		IsMentalTalented:         input.PJInfo.IsMentalTalented,
		IsCoordinationTalented:   input.PJInfo.IsCoordinationTalented,
		IsPhysicalSkillsTalented: input.PJInfo.IsPhysicalSkillsTalented,
		IsMentalSkillsTalented:   input.PJInfo.IsMentalSkillsTalented,
		IsEnergySkillsTalented:   input.PJInfo.IsEnergySkillsTalented,
		IsEnergyTalented:         input.PJInfo.IsEnergyTalented,
	}

	pj, err := camp.AddPJ(input.IDs.UserID, params, uc.identificationService)
	if err != nil {
		return applicationcampaign.PJOutput{}, err
	}

	if err = uc.campaignRepository.Save(ctx, camp); err != nil {
		return applicationcampaign.PJOutput{}, err
	}

	if err = uc.eventPublisher.Publish(ctx, camp.UncommittedEvents()); err != nil {
		return applicationcampaign.PJOutput{}, err
	}

	return applicationcampaign.MapPJOutput(pj), nil
}
