package campaign

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, campaign *Campaign) error
	FindByID(ctx context.Context, id string) (*Campaign, error)
}
