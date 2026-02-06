package campaign

import (
	"context"
)

//go:generate mockgen -destination=../../../tests/mocks/campaign_repository_mock.go -package=mocks -mock_names=Repository=MockCampaignRepository meye-core/internal/domain/campaign Repository
type Repository interface {
	Save(ctx context.Context, campaign *Campaign) error
	FindByID(ctx context.Context, id string) (*Campaign, error)
}

//go:generate mockgen -destination=../../../tests/mocks/pj_repository_mock.go -package=mocks meye-core/internal/domain/campaign PjRepository
type PjRepository interface {
	Save(ctx context.Context, pj *PJ) error
	FindByID(ctx context.Context, id string) (*PJ, error)
}
