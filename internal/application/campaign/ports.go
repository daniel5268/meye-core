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
