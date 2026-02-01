package identification

import "github.com/google/uuid"

type Service struct{}

// NewService creates a new Identification Service.
func NewService() *Service {
	return &Service{}
}

// GenerateID generates a new unique identifier.
func (s *Service) GenerateID() string {
	id := uuid.New().String()
	return id
}
