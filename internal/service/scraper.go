package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/agnosto/fansly-scraper/auth"
	"github.com/agnosto/fansly-scraper/headers"
	"github.com/agnosto/fansly-scraper/logger"
)

// ScraperService handles all interactions with the fansly-scraper
type ScraperService struct {
	logger          logger.Logger
	authConfig      *auth.Config
	headers         *headers.FanslyHeaders
	isAuthenticated bool
}

// NewScraperService creates a new ScraperService instance
func NewScraperService(logger logger.Logger) *ScraperService {
	headers := headers.New()
	authConfig := &auth.Config{
		Client:    *http.DefaultClient,
		UserAgent: headers.GetUserAgent(),
	}

	return &ScraperService{
		logger:       logger,
		headers:      headers,
		authConfig:   authConfig,
	}
}

// GetCreators retrieves a list of creators from Fansly
func (s *ScraperService) GetCreators(ctx context.Context, limit, offset int) ([]Creator, error) {
	s.logger.Infof("Fetching creators with limit=%d, offset=%d", limit, offset)

	if !s.isAuthenticated {
		return nil, errors.New("not authenticated with Fansly")
	}

	// TODO: Implement actual creator fetching logic using the fansly-scraper package
	// For now, return mock data
	creators := []Creator{
		{
			ID:          "1",
			Name:        "Example Creator 1",
			Username:    "creator1",
			IsVerified:  true,
			IsFollowing: true,
			LastUpdated: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:          "2",
			Name:        "Example Creator 2",
			Username:    "creator2",
			IsVerified:  false,
			IsFollowing: true,
			LastUpdated: time.Now().Add(-1 * time.Hour),
		},
	}

	// Apply pagination
	if offset >= len(creators) {
		return []Creator{}, nil
	}

	end := offset + limit
	if end > len(creators) {
		end = len(creators)
	}

	return creators[offset:end], nil
}

// Creator represents a Fansly creator
// TODO: Move this to a shared types package
type Creator struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	IsVerified  bool      `json:"is_verified"`
	IsFollowing bool      `json:"is_following"`
	LastUpdated time.Time `json:"last_updated"`
}

// Authenticate handles Fansly authentication
func (s *ScraperService) Authenticate(ctx context.Context, authToken string) error {
	if authToken == "" {
		return errors.New("authentication token is required")
	}

	s.logger.Info("Authenticating with Fansly")
	
	// Set the auth token in headers
	s.headers.Authorization = authToken
	s.authConfig.Authorization = authToken

	// Test authentication by getting account info
	_, err := auth.Login(s.headers)
	if err != nil {
		s.logger.Errorf("Authentication failed: %v", err)
		return fmt.Errorf("authentication failed: %w", err)
	}

	s.isAuthenticated = true
	s.logger.Info("Successfully authenticated with Fansly")
	return nil
}
