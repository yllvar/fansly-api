package api

import (
	"net/http"
	"strings"
)

// APIKeyAuth is a middleware that checks for a valid API key in the request
func (s *Server) APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health check and auth endpoints
		if r.URL.Path == "/" || 
		   r.URL.Path == "/api/v1/health" ||
		   strings.HasPrefix(r.URL.Path, "/api/v1/auth/") {
			next.ServeHTTP(w, r)
			return
		}

		// Get API key from header or query parameter
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			apiKey = r.URL.Query().Get("api_key")
		}

		// Validate API key
		if !s.config.IsValidAPIKey(apiKey) {
			s.log.Warnf("Invalid or missing API key from %s", r.RemoteAddr)
			respondWithError(w, http.StatusUnauthorized, "Invalid or missing API key")
			return
		}

		// Add user context if needed
		ctx := r.Context()
		// ctx = context.WithValue(ctx, "userID", userID) // You can add user context here
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth is a middleware that ensures the user is authenticated
func (s *Server) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for session token
		sessionID, err := r.Cookie("session_id")
		if err != nil || sessionID.Value == "" {
			s.log.Warn("Unauthorized access attempt - no session")
			respondWithError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		// TODO: Validate session
		// if !isValidSession(sessionID.Value) {
		// 	respondWithError(w, http.StatusUnauthorized, "Invalid or expired session")
		// 	return
		// }

		next.ServeHTTP(w, r)
	})
}
