// internal/domain/user/user.entity_test.go
package user_test

import (
	"errors"
	"testing"

	"meye-core/internal/domain/user"
	"meye-core/tests/mocks"
	"meye-core/tests/testdata"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewUser(t *testing.T) {
	type parameters struct {
		username string
		password string
		role     user.UserRole
	}

	type want struct {
		user *user.User
		err  error
	}

	var idServiceMock *mocks.MockIdentificationService
	var hashServiceMock *mocks.MockHashService

	testDefaultParams := parameters{
		username: testdata.Username,
		password: testdata.Password,
		role:     testdata.Role,
	}

	errTest := errors.New("mock_err")

	tests := []struct {
		name       string
		params     parameters
		want       want
		setupMocks func()
	}{
		{
			name:   "successful user creation",
			params: testDefaultParams,
			want: want{
				user: testdata.User(t),
				err:  nil,
			},
			setupMocks: func() {
				idServiceMock.EXPECT().
					GenerateID().
					Return(testdata.UserID).
					Times(1)

				hashServiceMock.EXPECT().
					Hash(testdata.Password).
					Return(testdata.HashedPassword, nil).
					Times(1)
			},
		},
		{
			name:   "hash service failure",
			params: testDefaultParams,
			want: want{
				user: nil,
				err:  errTest,
			},
			setupMocks: func() {
				idServiceMock.EXPECT().
					GenerateID().
					Return("").
					Times(1)

				hashServiceMock.EXPECT().
					Hash(gomock.Any()).
					Return("", errTest).
					Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			idServiceMock = mocks.NewMockIdentificationService(ctrl)
			hashServiceMock = mocks.NewMockHashService(ctrl)

			tt.setupMocks()

			newUser, err := user.NewUser(
				tt.params.username,
				tt.params.password,
				tt.params.role,
				idServiceMock,
				hashServiceMock,
			)

			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.user, newUser)
		})
	}
}
