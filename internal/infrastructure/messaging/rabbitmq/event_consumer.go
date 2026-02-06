package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

// EventHandler handles incoming events from RabbitMQ
type EventHandler interface {
	Handle(ctx context.Context, message EventMessage) error
}

// Consumer consumes events from RabbitMQ
type Consumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
	handler    EventHandler
}

// NewConsumer creates a new RabbitMQ event consumer
func NewConsumer(url, queueName string, handler EventHandler) (*Consumer, error) {
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

	// Set prefetch count to control how many messages are processed concurrently
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"queue": queueName,
	}).Info("RabbitMQ event consumer initialized")

	return &Consumer{
		connection: conn,
		channel:    ch,
		queueName:  queueName,
		handler:    handler,
	}, nil
}

// Start begins consuming messages from the queue
func (c *Consumer) Start(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.queueName, // queue
		"",          // consumer
		false,       // auto-ack (manual ack for reliability)
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	logrus.Info("Worker started. Waiting for messages...")

	for {
		select {
		case <-ctx.Done():
			logrus.Info("Stopping consumer due to context cancellation")
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				logrus.Warn("Message channel closed")
				return fmt.Errorf("message channel closed")
			}

			if err := c.processMessage(ctx, msg); err != nil {
				logrus.WithError(err).Error("Failed to process message")
				// Reject the message and requeue it
				if nackErr := msg.Nack(false, true); nackErr != nil {
					logrus.WithError(nackErr).Error("Failed to nack message")
				}
			} else {
				// Acknowledge the message
				if ackErr := msg.Ack(false); ackErr != nil {
					logrus.WithError(ackErr).Error("Failed to ack message")
				}
			}
		}
	}
}

// processMessage unmarshals and handles a single message
func (c *Consumer) processMessage(ctx context.Context, msg amqp.Delivery) error {
	var eventMessage EventMessage
	if err := json.Unmarshal(msg.Body, &eventMessage); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	if err := c.handler.Handle(ctx, eventMessage); err != nil {
		return fmt.Errorf("handler error for event %s: %w", eventMessage.ID, err)
	}

	return nil
}

// Close closes the RabbitMQ channel and connection
func (c *Consumer) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			logrus.Errorf("Failed to close RabbitMQ channel: %v", err)
			return err
		}
	}

	if c.connection != nil {
		if err := c.connection.Close(); err != nil {
			logrus.Errorf("Failed to close RabbitMQ connection: %v", err)
			return err
		}
	}

	logrus.Info("RabbitMQ event consumer closed")
	return nil
}
