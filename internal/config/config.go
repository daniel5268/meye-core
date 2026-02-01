package config

import (
	"fmt"
	"os"
)

type Api struct {
	Port   string
	ApiKey string
}

type Database struct {
	DSN string
}

type Config struct {
	Api      Api
	Database Database
}

func getInvalidVarErr(varName string) error {
	return fmt.Errorf("%s environment variable is invalid", varName)
}

func (cfg *Config) loadApi() error {
	cfg.Api.Port = os.Getenv("API_PORT")
	if cfg.Api.Port == "" {
		return getInvalidVarErr("API_PORT")
	}

	cfg.Api.ApiKey = os.Getenv("API_KEY")
	if cfg.Api.ApiKey == "" {
		return getInvalidVarErr("API_KEY")
	}

	return nil
}

func (cfg *Config) loadDB() error {
	cfg.Database.DSN = os.Getenv("DATABASE_DSN")
	if cfg.Database.DSN == "" {
		return getInvalidVarErr("DATABASE_DSN")
	}

	return nil
}

// New loads configuration from environment and returns the structure.
func New() (*Config, error) {
	cfg := &Config{}

	if err := cfg.loadApi(); err != nil {
		return nil, err
	}

	if err := cfg.loadDB(); err != nil {
		return nil, err
	}

	return cfg, nil
}
