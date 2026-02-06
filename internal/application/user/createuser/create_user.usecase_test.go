package createuser_test

import (
	"context"
	"errors"
	applicationuser "meye-core/internal/application/user"
	"meye-core/internal/application/user/createuser"
	"meye-core/internal/domain/user"
	"meye-core/tests/data"
	"meye-core/tests/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateUserUseCase_Execute(t *testing.T) {
	var idServiceMock *mocks.MockIdentificationService
	var userRepoMock *mocks.MockUserRepository
	var hashServiceMock *mocks.MockHashService
	var publisherMock *mocks.MockPublisher

	type want struct {
		output applicationuser.UserOutput
		err    error
	}

	ctx := context.Background()

	s := createuser.New(userRepoMock, idServiceMock, hashServiceMock, publisherMock)

	errTest := errors.New("mock_err")

	defaultInput := applicationuser.CreateUserInput{
		Username: data.Username,
		Password: data.Password,
		Role:     data.Role,
	}

	tests := []struct {
		name       string
		input      applicationuser.CreateUserInput
		want       want
		setupMocks func()
	}{
		{
			name:  "successful user creation",
			input: defaultInput,
			want: want{
				output: applicationuser.UserOutput{
					ID:       data.UserID,
					Username: data.Username,
					Role:     data.Role,
				},
				err: nil,
			},
			setupMocks: func() {
				idServiceMock.EXPECT().
					GenerateID().
					Return(data.UserID).
					Times(2)

				hashServiceMock.EXPECT().
					Hash(data.Password).
					Return(data.HashedPassword, nil).
					Times(1)

				userRepoMock.EXPECT().
					FindByUsername(ctx, data.Username).
					Return(nil, nil).
					Times(1)

				userRepoMock.EXPECT().
					Save(ctx, gomock.Any()).
					Return(nil).
					Times(1)

				publisherMock.EXPECT().
					Publish(ctx, gomock.Any()).
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
					Return(data.UserID).
					Times(2)

				hashServiceMock.EXPECT().
					Hash(data.Password).
					Return(data.HashedPassword, nil).
					Times(1)

				userRepoMock.EXPECT().
					FindByUsername(ctx, data.Username).
					Return(nil, nil).
					Times(1)

				userRepoMock.EXPECT().
					Save(ctx, gomock.Any()).
					Return(errTest).
					Times(1)
				publisherMock.EXPECT().
					Publish(ctx, gomock.Any()).
					Times(0)
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
					Return(data.UserID).
					Times(1)

				hashServiceMock.EXPECT().
					Hash(data.Password).
					Return("", errTest).
					Times(1)

				userRepoMock.EXPECT().
					FindByUsername(ctx, data.Username).
					Times(1)

				userRepoMock.EXPECT().
					Save(ctx, gomock.Any()).
					Times(0)

				publisherMock.EXPECT().
					Publish(ctx, gomock.Any()).
					Times(0)
			},
		},
		{
			name:  "username already exists",
			input: defaultInput,
			want: want{
				output: applicationuser.UserOutput{},
				err:    applicationuser.ErrUsernameAlreadyExists,
			},
			setupMocks: func() {
				existingUser := user.CreateUserWithoutValidation(
					"existing-id",
					data.Username,
					data.HashedPassword,
					data.Role,
				)

				idServiceMock.EXPECT().
					GenerateID().
					Return(data.UserID).
					Times(0)

				hashServiceMock.EXPECT().
					Hash(data.Password).
					Return(data.HashedPassword, nil).
					Times(0)

				userRepoMock.EXPECT().
					FindByUsername(ctx, data.Username).
					Return(existingUser, nil).
					Times(1)

				userRepoMock.EXPECT().
					Save(ctx, gomock.Any()).
					Times(0)

				publisherMock.EXPECT().
					Publish(ctx, gomock.Any()).
					Times(0)
			},
		},
		{
			name:  "error on find by username",
			input: defaultInput,
			want: want{
				output: applicationuser.UserOutput{},
				err:    errTest,
			},
			setupMocks: func() {
				idServiceMock.EXPECT().
					GenerateID().
					Return(data.UserID).
					Times(0)

				hashServiceMock.EXPECT().
					Hash(data.Password).
					Return(data.HashedPassword, nil).
					Times(0)

				userRepoMock.EXPECT().
					FindByUsername(ctx, data.Username).
					Return(nil, errTest).
					Times(1)

				userRepoMock.EXPECT().
					Save(ctx, gomock.Any()).
					Times(0)

				publisherMock.EXPECT().
					Publish(ctx, gomock.Any()).
					Times(0)
			},
		},
		{
			name:  "successful user creation",
			input: defaultInput,
			want: want{
				output: applicationuser.UserOutput{},
				err:    errTest,
			},
			setupMocks: func() {
				idServiceMock.EXPECT().
					GenerateID().
					Return(data.UserID).
					Times(2)

				hashServiceMock.EXPECT().
					Hash(data.Password).
					Return(data.HashedPassword, nil).
					Times(1)

				userRepoMock.EXPECT().
					FindByUsername(ctx, data.Username).
					Return(nil, nil).
					Times(1)

				userRepoMock.EXPECT().
					Save(ctx, gomock.Any()).
					Return(nil).
					Times(1)

				publisherMock.EXPECT().
					Publish(ctx, gomock.Any()).
					Return(errTest).
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
			userRepoMock = mocks.NewMockUserRepository(ctrl)
			publisherMock = mocks.NewMockPublisher(ctrl)

			tt.setupMocks()

			s = createuser.New(userRepoMock, idServiceMock, hashServiceMock, publisherMock)

			output, err := s.Execute(ctx, tt.input)

			assert.Equal(t, tt.want.output, output)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
