package main

import (
	"log"
	"net/http"
	"os"

	"github.com/agnosto/fansly-scraper/auth"
	"github.com/agnosto/fansly-scraper/headers"
	"github.com/agnosto/fansly-scraper/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get auth token from environment
	authToken := os.Getenv("FANSLY_AUTH_TOKEN")
	if authToken == "" {
		log.Fatal("FANSLY_AUTH_TOKEN is not set in .env file")
	}

	// Initialize logger
	log := logger.New()
	log.Info("Successfully loaded environment variables and initialized logger")

	// Truncate the token for logging purposes to avoid logging the full token
	tokenStart := ""
	tokenEnd := ""
	if len(authToken) > 10 {
		tokenStart = authToken[:5]
		tokenEnd = authToken[len(authToken)-5:]
	}
	log.Infof("Using auth token: %s...%s", tokenStart, tokenEnd)

	// Initialize headers
	headers := headers.New()
	headers.Authorization = authToken

	// Test authentication
	log.Info("Testing authentication...")
	user, err := auth.Login(headers)
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	log.Info("Successfully authenticated with Fansly")
	log.Infof("User ID: %s", user.ID)
	log.Infof("Username: %s", user.Username)
	log.Infof("Display Name: %s", user.DisplayName)
}
