package hash

import (
	"golang.org/x/crypto/bcrypt"
)

const cost = 12

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// Hash recieves a string and returns an encrypted hash.
func (s *Service) Hash(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), cost)

	return string(bytes), err
}

// Compare compares a hash with a secret string.
func (s *Service) Compare(secret, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
}
