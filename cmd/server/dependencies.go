package main

import (
	"fmt"
	"os"

	"meye-core/internal/application/campaign"
	"meye-core/internal/application/campaign/createcampaign"
	"meye-core/internal/application/campaign/inviteuser"
	"meye-core/internal/application/user"
	"meye-core/internal/application/user/createuser"
	"meye-core/internal/application/user/login"
	"meye-core/internal/config"
	"meye-core/internal/infrastructure/api"
	"meye-core/internal/infrastructure/api/handler"
	"meye-core/internal/infrastructure/hash"
	"meye-core/internal/infrastructure/identification"
	"meye-core/internal/infrastructure/jwt"
	postgresCampaignRepo "meye-core/internal/infrastructure/repository/campaign/postgres"
	postgresUserRepo "meye-core/internal/infrastructure/repository/user/postgres"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserUseCases struct {
	CreateUser user.CreateUserUseCase
	Login      user.LoginUseCase
}

type CampaignUseCases struct {
	CreateCampaign campaign.CreateCampaignUseCase
	InviteUser     campaign.InviteUserUseCase
}

type UseCases struct {
	User     *UserUseCases
	Campaign *CampaignUseCases
}

type Repositories struct {
	User     *postgresUserRepo.Repository
	Campaign *postgresCampaignRepo.Repository
}

type Services struct {
	Hash           *hash.Service
	Identification *identification.Service
	JWT            *jwt.Service
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
	container.initializeRepositories()
	container.initializeUseCases()
	container.initializeHandlers()
	container.initializeRouter()

	return container, nil
}

func (c *DependencyContainer) initializeServices() {
	c.Services = &Services{
		Hash:           hash.NewService(),
		Identification: identification.NewService(),
		JWT:            jwt.NewService(c.Config.JWT.Secret, c.Config.JWT.Issuer, c.Config.JWT.ExpirationTime),
	}
}

func (c *DependencyContainer) initializeRepositories() {
	c.Repositories = &Repositories{
		User:     postgresUserRepo.New(c.Database),
		Campaign: postgresCampaignRepo.New(c.Database),
	}
}

func (c *DependencyContainer) initializeUseCases() {
	c.UseCases = &UseCases{
		User: &UserUseCases{
			CreateUser: createuser.NewUseCase(
				c.Repositories.User,
				c.Services.Identification,
				c.Services.Hash,
			),
			Login: login.NewUseCase(
				c.Repositories.User,
				c.Services.Hash,
				c.Services.JWT,
			),
		},
		Campaign: &CampaignUseCases{
			CreateCampaign: createcampaign.NewUseCase(
				c.Repositories.Campaign,
				c.Services.Identification,
			),
			InviteUser: inviteuser.NewInviteUserUseCase(
				c.Repositories.Campaign,
				c.Repositories.User,
				c.Services.Identification,
			),
		},
	}
}

func (c *DependencyContainer) initializeHandlers() {
	c.Handlers = &Handlers{
		User:     handler.NewUserHandler(c.UseCases.User.CreateUser, c.UseCases.User.Login),
		Auth:     handler.NewAuthHandler(c.Config.Api.ApiKey, *c.Services.JWT, c.Repositories.User, c.Repositories.Campaign),
		Campaign: handler.NewCampaignHandler(c.UseCases.Campaign.CreateCampaign, c.UseCases.Campaign.InviteUser),
	}
}

func (c *DependencyContainer) initializeRouter() {
	c.APIRouter = api.NewRouter(&api.Handlers{
		UserHandler:     c.Handlers.User,
		AuthHandler:     c.Handlers.Auth,
		CampaignHandler: c.Handlers.Campaign,
	})
	logrus.Debug("Router initialized")
}

// Close gracefully closes all resources
func (c *DependencyContainer) Close() error {
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
