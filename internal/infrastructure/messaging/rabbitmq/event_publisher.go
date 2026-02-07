package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"meye-core/internal/domain/event"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

// Compile-time check to ensure Publisher implements the port interface
var _ event.Publisher = (*Publisher)(nil)

// Publisher implements the event.Publisher interface using RabbitMQ
type Publisher struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
}

// New creates a new RabbitMQ event publisher
func New(url, queueName string) (*Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare queue (idempotent operation)
	// durable: true - queue survives broker restart
	// autoDelete: false - queue is not deleted when last consumer unsubscribes
	// exclusive: false - queue can be accessed by other connections
	// noWait: false - wait for server confirmation
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"queue": queueName,
	}).Info("RabbitMQ event publisher initialized")

	return &Publisher{
		connection: conn,
		channel:    ch,
		queueName:  queueName,
	}, nil
}

// EventMessage represents the structure of the message sent to RabbitMQ
type EventMessage struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	AggregateID   string         `json:"aggregate_id"`
	AggregateType string         `json:"aggregate_type"`
	Data          map[string]any `json:"data"`
	CreatedAt     string         `json:"created_at"`
	OccurredAt    string         `json:"occurred_at"`
}

// Publish publishes a batch of domain events to RabbitMQ
func (p *Publisher) Publish(ctx context.Context, events []event.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	for _, evt := range events {
		message := EventMessage{
			ID:            evt.ID(),
			Type:          string(evt.Type()),
			AggregateID:   evt.AggregateID(),
			AggregateType: string(evt.AggregateType()),
			Data:          evt.GetSerializedData(),
			CreatedAt:     evt.CreatedAt().Format("2006-01-02T15:04:05.999Z07:00"),
			OccurredAt:    evt.OccurredAt().Format("2006-01-02T15:04:05.999Z07:00"),
		}

		body, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("failed to marshal event %s: %w", evt.ID(), err)
		}

		err = p.channel.PublishWithContext(
			ctx,
			"",          // exchange (default)
			p.queueName, // routing key (queue name)
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent, // make messages persistent
				ContentType:  "application/json",
				Body:         body,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to publish event %s: %w", evt.ID(), err)
		}

		logrus.WithFields(logrus.Fields{
			"event_id":       evt.ID(),
			"event_type":     evt.Type(),
			"aggregate_id":   evt.AggregateID(),
			"aggregate_type": evt.AggregateType(),
		}).Debug("Event published to RabbitMQ")
	}

	return nil
}

// Close closes the RabbitMQ channel and connection
func (p *Publisher) Close() error {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			logrus.Errorf("Failed to close RabbitMQ channel: %v", err)
			return err
		}
	}

	if p.connection != nil {
		if err := p.connection.Close(); err != nil {
			logrus.Errorf("Failed to close RabbitMQ connection: %v", err)
			return err
		}
	}

	logrus.Info("RabbitMQ event publisher closed")
	return nil
}
