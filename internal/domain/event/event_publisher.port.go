package event

import "context"

// Publisher defines the port for publishing domain events
type Publisher interface {
	// Publish publishes a batch of domain events
	Publish(ctx context.Context, events []DomainEvent) error
}