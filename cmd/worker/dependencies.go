package main

import (
	"fmt"
	"os"

	"meye-core/internal/application/campaign/consumexp"
	"meye-core/internal/config"
	"meye-core/internal/infrastructure/messaging/rabbitmq"
	postgresCampaignRepo "meye-core/internal/infrastructure/repository/campaign/postgres"
	"meye-core/internal/infrastructure/worker"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UseCases struct {
	ConsumeXp *consumexp.UseCase
}

type Repositories struct {
	PJ *postgresCampaignRepo.PjRepository
}

type Services struct {
	EventPublisher *rabbitmq.Publisher
}

type DependencyContainer struct {
	Config       *config.Config
	Database     *gorm.DB
	Services     *Services
	Repositories *Repositories
	UseCases     *UseCases
	EventHandler *worker.EventHandler
	Consumer     *rabbitmq.Consumer
}

func (c *DependencyContainer) loadEnvironment() error {
	if os.Getenv("GO_ENV") == "development" {
		if err := godotenv.Load(); err != nil {
			logrus.Warn("No .env file found, using system environment variables")
		}
	}

	return nil
}

func (c *DependencyContainer) loadConfig() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	c.Config = cfg

	return nil
}

func (c *DependencyContainer) connectDatabase() error {
	db, err := gorm.Open(postgres.Open(c.Config.Database.DSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	c.Database = db

	logrus.Info("Database connection established")

	return nil
}


func (c *DependencyContainer) initializeServices() error {
	publisher, err := rabbitmq.New(c.Config.RabbitMQ.URL, c.Config.RabbitMQ.EventsQueue)
	if err != nil {
		return fmt.Errorf("failed to create RabbitMQ publisher: %w", err)
	}

	c.Services = &Services{
		EventPublisher: publisher,
	}

	logrus.Info("Services initialized")
	return nil
}

func (c *DependencyContainer) initializeConsumer() error {
	consumer, err := rabbitmq.NewConsumer(
		c.Config.RabbitMQ.URL,
		c.Config.RabbitMQ.EventsQueue,
		c.EventHandler,
	)
	if err != nil {
		return fmt.Errorf("failed to create RabbitMQ consumer: %w", err)
	}

	c.Consumer = consumer

	return nil
}

func NewDependencyContainer() (*DependencyContainer, error) {
	container := &DependencyContainer{}

	if err := container.loadEnvironment(); err != nil {
		return nil, fmt.Errorf("failed to load environment: %w", err)
	}

	if err := container.loadConfig(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if err := container.connectDatabase(); err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	if err := container.initializeServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	container.initializeRepositories()
	container.initializeUseCases()
	container.initializeEventHandler()

	if err := container.initializeConsumer(); err != nil {
		return nil, fmt.Errorf("failed to initialize consumer: %w", err)
	}

	return container, nil
}

func (c *DependencyContainer) initializeRepositories() {
	c.Repositories = &Repositories{
		PJ: postgresCampaignRepo.NewPjRepository(c.Database),
	}
}

func (c *DependencyContainer) initializeUseCases() {
	c.UseCases = &UseCases{
		ConsumeXp: consumexp.New(
			c.Repositories.PJ,
			c.Services.EventPublisher,
		),
	}
}

func (c *DependencyContainer) initializeEventHandler() {
	c.EventHandler = worker.NewEventHandler(c.UseCases.ConsumeXp)
}

// Close gracefully closes all resources
func (c *DependencyContainer) Close() error {
	if c.Consumer != nil {
		if err := c.Consumer.Close(); err != nil {
			logrus.Errorf("Failed to close RabbitMQ consumer: %v", err)
			return err
		}
	}

	if c.Services != nil && c.Services.EventPublisher != nil {
		if err := c.Services.EventPublisher.Close(); err != nil {
			logrus.Errorf("Failed to close RabbitMQ connection: %v", err)
			return err
		}
	}

	if c.Database != nil {
		if sqlDB, err := c.Database.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				logrus.Errorf("Failed to close database connection: %v", err)
				return err
			}
			logrus.Info("Database connection closed")
		}
	}

	return nil
}
