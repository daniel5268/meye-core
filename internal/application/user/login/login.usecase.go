package login

import (
	"context"
	applicationuser "meye-core/internal/application/user"
	domainuser "meye-core/internal/domain/user"
)

type UseCase struct {
	userRepository domainuser.Repository
	hashService    domainuser.HashService
	jwtService     domainuser.JWTService
}

func NewUseCase(userRepo domainuser.Repository, hashServ domainuser.HashService, jwtServ domainuser.JWTService) *UseCase {
	return &UseCase{
		userRepository: userRepo,
		hashService:    hashServ,
		jwtService:     jwtServ,
	}
}

func (c *UseCase) Execute(ctx context.Context, input Input) (string, error) {
	user, err := c.userRepository.FindByUsername(ctx, input.Username)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", applicationuser.ErrInvalidCredentials
	}

	if err := c.hashService.Compare(input.Password, user.HashedPassword()); err != nil {
		return "", applicationuser.ErrInvalidCredentials
	}

	token, err := c.jwtService.GenerateSignedToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
