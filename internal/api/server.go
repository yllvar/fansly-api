package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"fansly-api/internal/logger"
)

// Server represents the HTTP server
type Server struct {
	server  *http.Server
	router *chi.Mux
	log    logger.Logger
}

// NewServer creates a new HTTP server
func NewServer(log logger.Logger) *Server {
	s := &Server{
		router: chi.NewRouter(),
		log:    log,
	}

	// Initialize the router and middleware
	s.setupMiddleware()
	s.setupRoutes()

	s.server = &http.Server{
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

func (s *Server) setupMiddleware() {
	// Basic middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration - simple for development
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.Get("/health", s.handleHealthCheck)

	// API v1 routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Get("/health", s.handleHealthCheck)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(s.requireAuth)
			r.Get("/creators", s.handleListCreators)
		})
	})
}

// Start starts the HTTP server on the specified address
func (s *Server) Start(addr string) error {
	s.server.Addr = addr
	s.log.Infof("Server starting on %s", addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Health check handler
func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	s.respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Helper function to send JSON responses
func (s *Server) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		s.log.Errorf("Error encoding JSON response: %v", err)
	}
}

// Helper function to send error responses
func (s *Server) respondWithError(w http.ResponseWriter, code int, message string) {
	s.respondWithJSON(w, code, map[string]string{"error": message})
}
