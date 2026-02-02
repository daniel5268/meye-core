package user

import "errors"

var (
	ErrUsernameAlreadyExists = errors.New("USERNAME_ALREADY_EXISTS")
	ErrInvalidCredentials    = errors.New("INVALID_CREDENTIALS")
)
