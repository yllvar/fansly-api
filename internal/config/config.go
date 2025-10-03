package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`  // e.g., ":8080"
	LogLevel      string `mapstructure:"LOG_LEVEL"`       // "debug", "info", "warn", "error"
	Environment   string `mapstructure:"ENV"`            // "development" or "production"
	JWTSecret     string `mapstructure:"JWT_SECRET"`     // Secret for signing JWT tokens
	
	// OAuth2 Configuration
	AuthURL       string `mapstructure:"AUTH_URL"`       // OAuth2 authorization URL
	TokenURL      string `mapstructure:"TOKEN_URL"`      // OAuth2 token URL
	CallbackURL   string `mapstructure:"CALLBACK_URL"`   // OAuth2 callback URL
	ClientID      string `mapstructure:"CLIENT_ID"`      // OAuth2 client ID
	ClientSecret  string `mapstructure:"CLIENT_SECRET"`  // OAuth2 client secret
	Scopes        string `mapstructure:"SCOPES"`         // OAuth2 scopes (comma-separated)
	
	// Session configuration
	SessionSecret string `mapstructure:"SESSION_SECRET"` // Secret for encrypting sessions
	SessionMaxAge int    `mapstructure:"SESSION_MAX_AGE"` // Session max age in seconds
}

func Load() (*Config, error) {
	// Set default values
	viper.SetDefault("SERVER_ADDRESS", ":8080")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("ENV", "development")
	viper.SetDefault("AUTH_URL", "https://fansly.com/oauth2/authorize")
	viper.SetDefault("TOKEN_URL", "https://fansly.com/oauth2/token")
	viper.SetDefault("CALLBACK_URL", "http://localhost:8080/api/v1/auth/callback")

	// Read from environment variables
	viper.AutomaticEnv()

	// Read from .env file if it exists
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Try to read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.JWTSecret == "" && c.Environment == "production" {
		return fmt.Errorf("JWT_SECRET is required in production")
	}
	return nil
}

// IsValidAPIKey checks if the provided API key is valid
// In a real application, you would validate against a database or configuration
// For this example, we'll just check if it's not empty
func (c *Config) IsValidAPIKey(apiKey string) bool {
	// In a real application, you would validate the API key against a database or configuration
	// For now, we'll just check if it's not empty
	return apiKey != ""
}

// GetDataDir returns the directory where data should be stored
func (c *Config) GetDataDir() (string, error) {
	// In production, this should be a proper data directory
	if c.Environment == "production" {
		return "/var/lib/fansly-api", nil
	}

	// In development, use a local directory
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	dataDir := filepath.Join(home, ".fansly-api")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create data directory: %w", err)
	}

	return dataDir, nil
}
