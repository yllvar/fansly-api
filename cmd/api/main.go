package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fansly-api/internal/api"
	"fansly-api/internal/logger"
)

func main() {
	// Initialize logger
	log := logger.New()

	// Load environment variables from .env file
	if err := loadEnv(); err != nil {
		log.Errorf("Error loading .env file: %v", err)
		os.Exit(1)
	}

	// Get auth token from environment
	authToken := os.Getenv("FANSLY_AUTH_TOKEN")
	if authToken == "" {
		log.Errorf("FANSLY_AUTH_TOKEN is required")
		os.Exit(1)
	}

	// Create and start the server
	server := api.NewServer(log)
	log.Info("Starting server on :8080")
	
	// Start the server in a goroutine
	go func() {
		if err := server.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Errorf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited properly")
}

// loadEnv loads environment variables from .env file
func loadEnv() error {
	// Check if .env file exists
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		// .env file doesn't exist, create it with default values
		return createDefaultEnvFile()
	}

	// Load .env file
	file, err := os.Open(".env")
	if err != nil {
		return fmt.Errorf("error opening .env file: %w", err)
	}
	defer file.Close()

	return nil
}

// createDefaultEnvFile creates a default .env file with required variables
func createDefaultEnvFile() error {
	defaultEnv := `# Fansly API Configuration
FANSLY_AUTH_TOKEN=your_auth_token_here
`
	
	if err := os.WriteFile(".env", []byte(defaultEnv), 0644); err != nil {
		return fmt.Errorf("error creating default .env file: %w", err)
	}
	
	return nil
}
