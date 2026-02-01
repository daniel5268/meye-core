package user

import "context"

//go:generate mockgen -destination=../../../tests/mocks/user_repository_mock.go -package=mocks meye-core/internal/domain/user Repository
type Repository interface {
	Save(ctx context.Context, user *User) error
}
