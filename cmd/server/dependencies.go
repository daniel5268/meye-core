package main

import (
	"fmt"
	"os"

	"meye-core/internal/application/campaign/createcampaign"
	"meye-core/internal/application/campaign/createpj"
	"meye-core/internal/application/campaign/getcampaign"
	"meye-core/internal/application/campaign/getcampaigns"
	"meye-core/internal/application/campaign/getpj"
	"meye-core/internal/application/campaign/inviteuser"
	"meye-core/internal/application/campaign/updatepjstats"
	"meye-core/internal/application/session/createsession"
	"meye-core/internal/application/user/createuser"
	"meye-core/internal/application/user/getplayers"
	"meye-core/internal/application/user/login"
	"meye-core/internal/config"
	"meye-core/internal/infrastructure/api"
	"meye-core/internal/infrastructure/api/handler"
	"meye-core/internal/infrastructure/hash"
	"meye-core/internal/infrastructure/identification"
	"meye-core/internal/infrastructure/jwt"
	"meye-core/internal/infrastructure/messaging/rabbitmq"
	postgresCampaignRepo "meye-core/internal/infrastructure/repository/campaign/postgres"
	postgresSessionRepo "meye-core/internal/infrastructure/repository/session/postgres"
	postgresUserRepo "meye-core/internal/infrastructure/repository/user/postgres"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserUseCases struct {
	CreateUser *createuser.UseCase
	Login      *login.UseCase
	GetPlayers *getplayers.UseCase
}

type CampaignUseCases struct {
	CreateCampaign *createcampaign.UseCase
	InviteUser     *inviteuser.UseCase
	CreatePJ       *createpj.UseCase
	UpdatePjStats  *updatepjstats.UseCase
	GetCampaign    *getcampaign.UseCase
	GetPj          *getpj.UseCase
	GetCampaigns   *getcampaigns.UseCase
}

type SessionUseCases struct {
	CreateSession *createsession.UseCase
}

type UseCases struct {
	User     *UserUseCases
	Campaign *CampaignUseCases
	Session  *SessionUseCases
}

type Repositories struct {
	User                 *postgresUserRepo.Repository
	Campaign             *postgresCampaignRepo.Repository
	Session              *postgresSessionRepo.Repository
	PJ                   *postgresCampaignRepo.PjRepository
	CampaignQueryService *postgresCampaignRepo.CampaignQueryService
}

type Services struct {
	Hash           *hash.Service
	Identification *identification.Service
	JWT            *jwt.Service
	EventPublisher *rabbitmq.Publisher
}

type Handlers struct {
	User     *handler.UserHandler
	Auth     *handler.AuthHandler
	Campaign *handler.CampaignHandler
}

type DependencyContainer struct {
	Config       *config.Config
	Database     *gorm.DB
	Services     *Services
	Repositories *Repositories
	UseCases     *UseCases
	Handlers     *Handlers
	APIRouter    *api.Router
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

func (c *DependencyContainer) connectRabbitMQ() error {
	publisher, err := rabbitmq.New(c.Config.RabbitMQ.URL, c.Config.RabbitMQ.EventsQueue)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	c.Services.EventPublisher = publisher

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

	container.initializeServices()

	if err := container.connectRabbitMQ(); err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	container.initializeRepositories()
	container.initializeUseCases()
	container.initializeHandlers()
	container.initializeRouter()

	return container, nil
}

func (c *DependencyContainer) initializeServices() {
	c.Services = &Services{
		Hash:           hash.New(),
		Identification: identification.New(),
		JWT:            jwt.New(c.Config.JWT.Secret, c.Config.JWT.Issuer, c.Config.JWT.ExpirationTime),
	}
}

func (c *DependencyContainer) initializeRepositories() {
	c.Repositories = &Repositories{
		User:                 postgresUserRepo.New(c.Database),
		Campaign:             postgresCampaignRepo.New(c.Database),
		Session:              postgresSessionRepo.New(c.Database),
		PJ:                   postgresCampaignRepo.NewPjRepository(c.Database),
		CampaignQueryService: postgresCampaignRepo.NewQueryService(c.Database),
	}
}

func (c *DependencyContainer) initializeUseCases() {
	c.UseCases = &UseCases{
		User: &UserUseCases{
			CreateUser: createuser.New(
				c.Repositories.User,
				c.Services.Identification,
				c.Services.Hash,
				c.Services.EventPublisher,
			),
			Login: login.New(
				c.Repositories.User,
				c.Services.Hash,
				c.Services.JWT,
			),
			GetPlayers: getplayers.New(
				c.Repositories.User,
			),
		},
		Campaign: &CampaignUseCases{
			CreateCampaign: createcampaign.New(
				c.Repositories.Campaign,
				c.Services.Identification,
				c.Services.EventPublisher,
			),
			InviteUser: inviteuser.New(
				c.Repositories.Campaign,
				c.Repositories.User,
				c.Services.Identification,
				c.Services.EventPublisher,
			),
			CreatePJ: createpj.New(
				c.Repositories.Campaign,
				c.Repositories.User,
				c.Services.Identification,
				c.Services.EventPublisher,
			),
			UpdatePjStats: updatepjstats.New(
				c.Repositories.PJ,
			),
			GetCampaign: getcampaign.New(
				c.Repositories.Campaign,
			),
			GetPj: getpj.New(
				c.Repositories.PJ,
			),
			GetCampaigns: getcampaigns.New(
				c.Repositories.CampaignQueryService,
			),
		},
		Session: &SessionUseCases{
			CreateSession: createsession.New(
				c.Repositories.Session,
				c.Repositories.Campaign,
				c.Services.Identification,
				c.Services.EventPublisher,
			),
		},
	}
}

func (c *DependencyContainer) initializeHandlers() {
	c.Handlers = &Handlers{
		User: handler.NewUserHandler(
			c.UseCases.User.CreateUser,
			c.UseCases.User.Login,
			c.UseCases.User.GetPlayers,
		),
		Auth: handler.NewAuthHandler(
			c.Config.Api.ApiKey,
			c.Services.JWT,
			c.Repositories.User,
			c.Repositories.Campaign,
			c.Repositories.PJ,
		),
		Campaign: handler.NewCampaignHandler(
			c.UseCases.Campaign.CreateCampaign,
			c.UseCases.Campaign.InviteUser,
			c.UseCases.Campaign.CreatePJ,
			c.UseCases.Session.CreateSession,
			c.UseCases.Campaign.UpdatePjStats,
			c.UseCases.Campaign.GetCampaign,
			c.UseCases.Campaign.GetPj,
			c.UseCases.Campaign.GetCampaigns,
		),
	}
}

func (c *DependencyContainer) initializeRouter() {
	c.APIRouter = api.NewRouter(&api.Handlers{
		UserHandler:     c.Handlers.User,
		AuthHandler:     c.Handlers.Auth,
		CampaignHandler: c.Handlers.Campaign,
	}, c.Config.Api.AllowedOrigins)
	logrus.Debug("Router initialized")
}

// Close gracefully closes all resources
func (c *DependencyContainer) Close() error {
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
