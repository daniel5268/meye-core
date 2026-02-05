package login

import (
	"context"
	applicationuser "meye-core/internal/application/user"
	domainuser "meye-core/internal/domain/user"
)

// Compile-time check to ensure UseCase implements the port interface
var _ applicationuser.LoginUseCase = (*UseCase)(nil)

type UseCase struct {
	userRepository domainuser.Repository
	hashService    domainuser.HashService
	jwtService     domainuser.JWTService
}

func New(userRepo domainuser.Repository, hashServ domainuser.HashService, jwtServ domainuser.JWTService) *UseCase {
	return &UseCase{
		userRepository: userRepo,
		hashService:    hashServ,
		jwtService:     jwtServ,
	}
}

func (uc *UseCase) Execute(ctx context.Context, input applicationuser.LoginInput) (string, error) {
	user, err := uc.userRepository.FindByUsername(ctx, input.Username)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", applicationuser.ErrInvalidCredentials
	}

	if err := uc.hashService.Compare(input.Password, user.HashedPassword()); err != nil {
		return "", applicationuser.ErrInvalidCredentials
	}

	token, err := uc.jwtService.GenerateSignedToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
