package createuser_test

import (
	"context"
	"errors"
	applicationuser "meye-core/internal/application/user"
	"meye-core/internal/application/user/createuser"
	"meye-core/tests/mocks"
	"meye-core/tests/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateUserUseCase_Execute(t *testing.T) {
	var idServiceMock *mocks.MockIdentificationService
	var userRepoMock *mocks.MockRepository
	var hashServiceMock *mocks.MockHashService

	type want struct {
		output applicationuser.UserOutput
		err    error
	}

	ctx := context.Background()

	s := createuser.NewUseCase(userRepoMock, idServiceMock, hashServiceMock)

	errTest := errors.New("mock_err")

	defaultInput := createuser.Input{
		Username: testdata.Username,
		Password: testdata.Password,
		Role:     testdata.Role,
	}

	tests := []struct {
		name       string
		input      createuser.Input
		want       want
		setupMocks func()
	}{
		{
			name:  "successful user creation",
			input: defaultInput,
			want: want{
				output: applicationuser.UserOutput{
					ID:       testdata.UserID,
					Username: testdata.Username,
					Role:     testdata.Role,
				},
				err: nil,
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

				userRepoMock.EXPECT().
					Save(ctx, gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:  "error on repository save",
			input: defaultInput,
			want: want{
				output: applicationuser.UserOutput{},
				err:    errTest,
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

				userRepoMock.EXPECT().
					Save(ctx, gomock.Any()).
					Return(errTest).
					Times(1)
			},
		},
		{
			name:  "error on domain user creation",
			input: defaultInput,
			want: want{
				output: applicationuser.UserOutput{},
				err:    errTest,
			},
			setupMocks: func() {
				idServiceMock.EXPECT().
					GenerateID().
					Return(testdata.UserID).
					Times(1)

				hashServiceMock.EXPECT().
					Hash(testdata.Password).
					Return("", errTest).
					Times(1)

				userRepoMock.EXPECT().
					Save(ctx, gomock.Any()).
					Times(0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			idServiceMock = mocks.NewMockIdentificationService(ctrl)
			hashServiceMock = mocks.NewMockHashService(ctrl)
			userRepoMock = mocks.NewMockRepository(ctrl)

			tt.setupMocks()

			s = createuser.NewUseCase(userRepoMock, idServiceMock, hashServiceMock)

			output, err := s.Execute(ctx, tt.input)

			assert.Equal(t, tt.want.output, output)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
