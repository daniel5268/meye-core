package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

// Server wraps the HTTP server with graceful shutdown capabilities
type Server struct {
	httpServer   *http.Server
	dependencies *DependencyContainer
}

// NewServer creates a new HTTP server with the provided container
func NewServer(dependencies *DependencyContainer) *Server {
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", dependencies.Config.Api.Port),
		Handler:      dependencies.APIRouter.Engine(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer:   httpServer,
		dependencies: dependencies,
	}
}

// Start begins listening for HTTP requests
func (s *Server) Start() error {
	migrationsPath := "file://migrations"
	m, err := migrate.New(migrationsPath, s.dependencies.Config.Database.DSN)
	if err != nil {
		logrus.Fatalf("could not create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("failed to run migrations: %v", err)
	}

	logrus.Info("Migrations applied successfully")

	logrus.Infof("Starting HTTP server on %s", s.httpServer.Addr)

	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	logrus.Info("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}
