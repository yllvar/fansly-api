package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// authResponse represents the response for authentication endpoints
type authResponse struct {
	URL       string `json:"url,omitempty"`
	Token     string `json:"token,omitempty"`
	ExpiresIn int    `json:"expires_in,omitempty"`
}

// authRequest represents the request for completing authentication
type authRequest struct {
	AuthToken string `json:"auth_token"`
	UserAgent string `json:"user_agent"`
}

type claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// handleAuthInitiate starts the authentication process
func (s *Server) handleAuthInitiate(w http.ResponseWriter, r *http.Request) {
	// Generate a random token for this authentication attempt
	authToken, err := generateRandomString(32)
	if err != nil {
		s.log.Errorf("Failed to generate auth token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to start authentication")
		return
	}

	// Store the auth token in memory (in a real app, use a proper session store)
	s.authTokens[authToken] = time.Now().Add(10 * time.Minute) // Token valid for 10 minutes

	// Return the authentication URL and token
	respondWithJSON(w, http.StatusOK, authResponse{
		URL:   "https://fansly.com/account/security", // URL where user can find their auth token
		Token: authToken,
	})
}

// handleAuthComplete completes the authentication process
func (s *Server) handleAuthComplete(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.log.Warnf("Invalid request: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Verify the auth token is valid and not expired
	expiresAt, exists := s.authTokens[req.AuthToken]
	if !exists {
		s.log.Warnf("Invalid auth token: %s", req.AuthToken)
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired authentication token")
		return
	}

	if time.Now().After(expiresAt) {
		delete(s.authTokens, req.AuthToken)
		s.log.Warnf("Expired auth token: %s", req.AuthToken)
		respondWithError(w, http.StatusUnauthorized, "Authentication token expired")
		return
	}

	// Clean up the used token
	delete(s.authTokens, req.AuthToken)

	// In a real implementation, you would validate the auth token
	// and exchange it for a JWT

	// For now, just return a placeholder response
	token, err := s.generateJWT("user123") // Replace with actual user ID
	if err != nil {
		s.log.Errorf("Failed to generate JWT: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	respondWithJSON(w, http.StatusOK, authResponse{
		Token:     token,
		ExpiresIn: 3600, // 1 hour
	})
}

// requireAuth is a middleware that ensures the request is authenticated
func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		// Extract the token from the header (format: "Bearer <token>")
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		if tokenString == "" {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			s.log.Warnf("Invalid token: %v", err)
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Token is valid, continue with the request
		next.ServeHTTP(w, r)
	})
}

// generateRandomString generates a random string of the specified length
func generateRandomString(length int) (string, error) {
	b := make([]byte, (length+1)/2) // Can be simplified to length/2, but this accounts for odd lengths
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b)[:length], nil
}

// generateJWT creates a new JWT token for the given user ID
func (s *Server) generateJWT(userID string) (string, error) {
	// Create the claims
	claims := &claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "fansly-api",
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
