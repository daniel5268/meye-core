package campaign

import "context"

type CreateCampaignInput struct {
	Name     string
	MasterID string
}

type CreateCampaignUseCase interface {
	Execute(ctx context.Context, input *CreateCampaignInput) (CampaignOutput, error)
}
