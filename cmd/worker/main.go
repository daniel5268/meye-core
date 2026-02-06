package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Starting worker...")

	// Initialize the application container with all dependencies
	container, err := NewDependencyContainer()
	if err != nil {
		logrus.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consumer in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := container.Consumer.Start(ctx); err != nil && err != context.Canceled {
			errChan <- err
		}
	}()

	// Wait for interrupt signal or error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		logrus.Info("Received shutdown signal")
	case err := <-errChan:
		logrus.Errorf("Worker error: %v", err)
	}

	// Cancel context to stop consumer
	cancel()

	logrus.Info("Worker shutdown complete")
}
