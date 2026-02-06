package worker

import (
	"context"
	"fmt"
	"meye-core/internal/application/campaign"
	"meye-core/internal/domain/event"
	"meye-core/internal/infrastructure/messaging/rabbitmq"

	"github.com/sirupsen/logrus"
)

// EventHandler routes events to their appropriate use case handlers
type EventHandler struct {
	consumeXpUseCase campaign.ConsumeXpUseCase
}

// NewEventHandler creates a new event handler
func NewEventHandler(consumeXpUseCase campaign.ConsumeXpUseCase) *EventHandler {
	return &EventHandler{
		consumeXpUseCase: consumeXpUseCase,
	}
}

// Handle processes an event message
func (h *EventHandler) Handle(ctx context.Context, message rabbitmq.EventMessage) error {
	switch event.EventType(message.Type) {
	case event.EventTypeXPAssigned:
		return h.handleXPAssigned(ctx, message)
	default:
		// Ignore unknown event types
		logrus.WithFields(logrus.Fields{
			"event_type": message.Type,
			"event_id":   message.ID,
		}).Debug("Ignoring unknown event type")
		return nil
	}
}

// handleXPAssigned processes the XPAssigned event
func (h *EventHandler) handleXPAssigned(ctx context.Context, message rabbitmq.EventMessage) error {
	logrus.WithFields(logrus.Fields{
		"event_id":     message.ID,
		"aggregate_id": message.AggregateID, // This is the PJ ID
	}).Info("Processing XPAssigned event")

	assignedXPData, ok := message.Data["assigned_xp"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("missing or invalid assigned_xp data in event %s", message.ID)
	}

	basic, err := convertToUint(assignedXPData["basic"])
	if err != nil {
		return fmt.Errorf("invalid basic xp value: %w", err)
	}

	special, err := convertToUint(assignedXPData["special"])
	if err != nil {
		return fmt.Errorf("invalid special xp value: %w", err)
	}

	supernatural, err := convertToUint(assignedXPData["supernatural"])
	if err != nil {
		return fmt.Errorf("invalid supernatural xp value: %w", err)
	}

	input := campaign.ConsumeXpInput{
		PjID: message.AggregateID, // PJ ID from aggregate_id
		Xp: campaign.XpAmounts{
			Basic:        basic,
			Special:      special,
			Supernatural: supernatural,
		},
	}

	if err := h.consumeXpUseCase.Execute(ctx, input); err != nil {
		return fmt.Errorf("failed to consume XP for PJ %s: %w", message.AggregateID, err)
	}

	return nil
}

// convertToUint converts an interface{} value to uint
func convertToUint(value interface{}) (uint, error) {
	switch v := value.(type) {
	case float64:
		return uint(v), nil
	case int:
		return uint(v), nil
	case uint:
		return v, nil
	case int64:
		return uint(v), nil
	case uint64:
		return uint(v), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to uint", value)
	}
}
