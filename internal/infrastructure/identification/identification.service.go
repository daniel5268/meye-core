package identification

import (
	"meye-core/internal/domain/shared"

	"github.com/google/uuid"
)

var _ shared.IdentificationService = (*Service)(nil)

type Service struct{}

// New creates a new Identification Service.
func New() *Service {
	return &Service{}
}

// GenerateID generates a new unique identifier.
func (s *Service) GenerateID() string {
	id := uuid.New().String()
	return id
}
