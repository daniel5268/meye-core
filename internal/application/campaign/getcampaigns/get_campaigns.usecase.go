package getcampaigns

import (
	"context"
	applicationcampaign "meye-core/internal/application/campaign"
	domaincampaign "meye-core/internal/domain/campaign"
)

var _ applicationcampaign.GetCampaignsUseCase = (*UseCase)(nil)

type UseCase struct {
	queryService domaincampaign.CampaignQueryService
}

func New(queryServ domaincampaign.CampaignQueryService) *UseCase {
	return &UseCase{
		queryService: queryServ,
	}
}

func (uc *UseCase) Execute(ctx context.Context, masterID string) ([]applicationcampaign.CampaignBasicInfoOutput, error) {
	c, err := uc.queryService.GetCampaignsBasicInfo(ctx, masterID)
	if err != nil {
		return []applicationcampaign.CampaignBasicInfoOutput{}, err
	}

	output := make([]applicationcampaign.CampaignBasicInfoOutput, 0, len(c))
	for _, info := range c {
		output = append(output, applicationcampaign.MapCampaignBasicInfoOutput(info))
	}

	return output, nil
}
