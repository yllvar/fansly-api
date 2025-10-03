package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"fansly-api/internal/logger"
)

type FanslyClient struct {
	baseURL    string
	authToken  string
	httpClient *http.Client
	logger     logger.Logger
}

func NewFanslyClient(authToken string, logger logger.Logger) *FanslyClient {
	return &FanslyClient{
		baseURL:    "https://apiv3.fansly.com",
		authToken:  authToken,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		logger:     logger,
	}
}

// GetAccountInfo retrieves the current user's account information
func (c *FanslyClient) GetAccountInfo(ctx context.Context) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/account/me", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", c.authToken)
	req.Header.Set("User-Agent", "fansly-api/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return result, nil
}

// GetFollowedUsers retrieves a list of users that the authenticated user is following
func (c *FanslyClient) GetFollowedUsers(ctx context.Context, limit, offset int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/account/me/following?limit=%d&offset=%d", c.baseURL, limit, offset)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", c.authToken)
	req.Header.Set("User-Agent", "fansly-api/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Response struct {
			Accounts []map[string]interface{} `json:"accounts"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return result.Response.Accounts, nil
}
