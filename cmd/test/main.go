package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"fansly-api/internal/api"
	"fansly-api/internal/logger"
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
	log.Info("Initializing Fansly API client...")

	// Create a new Fansly client
	client := api.NewFanslyClient(authToken, log)

	// Test account info
	log.Info("Fetching account info...")
	accountInfo, err := client.GetAccountInfo(context.Background())
	if err != nil {
		log.Fatalf("Error getting account info: %v", err)
	}

	// Print account info
	log.Info("Successfully retrieved account info:")
	printJSON(accountInfo)

	// Test followed users
	log.Info("Fetching followed users...")
	followed, err := client.GetFollowedUsers(context.Background(), 10, 0)
	if err != nil {
		log.Fatalf("Error getting followed users: %v", err)
	}

	log.Infof("Found %d followed users:", len(followed))
	for i, user := range followed {
		if username, ok := user["username"].(string); ok {
			log.Infof("%d. %s", i+1, username)
		}
	}

	log.Info("Test completed successfully")
}

// printJSON pretty prints the given data as JSON
func printJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}
	fmt.Println(string(jsonData))
}
