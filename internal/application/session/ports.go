package session

import "context"

type CreateSessionUseCase interface {
	Execute(ctx context.Context, input CreateSessionInput) (SessionOutput, error)
}
