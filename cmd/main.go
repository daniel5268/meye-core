package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("starting...")
	// Initialize the application container with all dependencies
	container, err := NewDependencyContainer()
	if err != nil {
		logrus.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	// Create and configure the server
	server := NewServer(container)

	// Start server in a goroutine
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	logrus.Println("Server exited gracefully")
}
