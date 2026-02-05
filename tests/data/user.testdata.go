package data

import (
	"meye-core/internal/domain/user"
	"meye-core/tests/mocks"
	"testing"

	"go.uber.org/mock/gomock"
)

const (
	UserID         = "user-123"
	Username       = "testuser"
	Password       = "password123"
	HashedPassword = "hashed:password123"
	Role           = user.UserRolePlayer
)

func User(t *testing.T) *user.User {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idServ := mocks.NewMockIdentificationService(ctrl)
	hashServ := mocks.NewMockHashService(ctrl)

	idServ.EXPECT().GenerateID().Return(UserID)
	hashServ.EXPECT().Hash(Password).Return(HashedPassword, nil)

	u, _ := user.NewUser(Username, Password, Role, idServ, hashServ)

	return u
}
