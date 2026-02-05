package createsession

import (
	"context"
	applicationsession "meye-core/internal/application/session"
	"meye-core/internal/domain/campaign"
	domainsession "meye-core/internal/domain/session"
	"meye-core/internal/domain/shared"
)

var _ applicationsession.CreateSessionUseCase = (*UseCase)(nil)

type UseCase struct {
	sessionRepository     domainsession.Repository
	campaignRepository    campaign.Repository
	identificationService shared.IdentificationService
}

func New(sessRepo domainsession.Repository, campRepo campaign.Repository, idServ shared.IdentificationService) *UseCase {
	return &UseCase{
		sessionRepository:     sessRepo,
		campaignRepository:    campRepo,
		identificationService: idServ,
	}
}

// Compile-time check to ensure UseCase implements the port interface
var _ applicationsession.CreateSessionUseCase = (*UseCase)(nil)

func (uc *UseCase) Execute(ctx context.Context, input applicationsession.CreateSessionInput) (applicationsession.SessionOutput, error) {
	camp, err := uc.campaignRepository.FindByID(ctx, input.CampaignID)
	if err != nil {
		return applicationsession.SessionOutput{}, err
	}

	xpALength := len(input.XPAssignations)

	pjIDs := make([]string, 0, xpALength)
	for _, xpA := range input.XPAssignations {
		pjIDs = append(pjIDs, xpA.PjID)
	}

	if err = camp.MustContainPjs(pjIDs); err != nil {
		return applicationsession.SessionOutput{}, err
	}

	xpAssignations := make([]domainsession.XPAssignation, 0, xpALength)
	for _, xpInput := range input.XPAssignations {
		xpAssignations = append(xpAssignations, domainsession.NewXPAssignation(xpInput.PjID, xpInput.Amounts.Basic, xpInput.Amounts.Special, xpInput.Amounts.SuperNatural, xpInput.Reason))
	}

	session, err := domainsession.NewSession(camp.MasterID(), input.CampaignID, input.Summary, xpAssignations, uc.identificationService)
	if err != nil {
		return applicationsession.SessionOutput{}, err
	}

	err = uc.sessionRepository.Save(ctx, session)
	if err != nil {
		return applicationsession.SessionOutput{}, err
	}

	return applicationsession.MapSessionOutput(session), nil
}
