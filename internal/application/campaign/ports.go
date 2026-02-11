package campaign

import (
	"context"
)

type CreateCampaignUseCase interface {
	Execute(ctx context.Context, input CreateCampaignInput) (CampaignOutput, error)
}

type InviteUserUseCase interface {
	Execute(ctx context.Context, input InviteUserInput) (InvitationOutput, error)
}

type CreatePJUseCase interface {
	Execute(ctx context.Context, input CreatePJInput) (PJOutput, error)
}

type ConsumeXpUseCase interface {
	Execute(ctx context.Context, input ConsumeXpInput) error
}

type UpdateStatsUseCase interface {
	Execute(ctx context.Context, input UpdatePjStatsInput) (PJOutput, error)
}

type GetCampaignUseCase interface {
	Execute(ctx context.Context, campID string) (CampaignOutput, error)
}

type GetPjUseCase interface {
	Execute(ctx context.Context, pjID string) (PJOutput, error)
}

type GetCampaignsUseCase interface {
	Execute(ctx context.Context, masterID string) ([]CampaignBasicInfoOutput, error)
}

type GetPjsUseCase interface {
	Execute(ctx context.Context, userID string) ([]PjBasicInfoOutput, error)
}
