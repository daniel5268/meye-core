package config

import (
	"fmt"
	"os"
	"time"
)

type Api struct {
	Port   string
	ApiKey string
}

type Database struct {
	DSN string
}

type JWT struct {
	Secret         string
	Issuer         string
	ExpirationTime time.Duration
}

type Config struct {
	Api      Api
	Database Database
	JWT      JWT
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

func (cfg *Config) loadJWT() error {
	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	if cfg.JWT.Secret == "" {
		return getInvalidVarErr("JWT_SECRET")
	}

	cfg.JWT.Issuer = os.Getenv("JWT_ISSUER")
	if cfg.JWT.Issuer == "" {
		return getInvalidVarErr("JWT_ISSUER")
	}

	jwtExpTimeStr := os.Getenv("JWT_EXPIRATION_TIME")
	jwtExpTime, err := time.ParseDuration(jwtExpTimeStr)
	if err != nil {
		return getInvalidVarErr("JWT_EXPIRATION_TIME")
	}
	cfg.JWT.ExpirationTime = jwtExpTime
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

	if err := cfg.loadJWT(); err != nil {
		return nil, err
	}

	return cfg, nil
}
