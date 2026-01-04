package controlplane

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Vaibhav2154/ShadowNet/internal/controlplane/api"
	"github.com/Vaibhav2154/ShadowNet/internal/controlplane/service"
	"github.com/Vaibhav2154/ShadowNet/internal/controlplane/store"
)

// Config holds server configuration
type Config struct {
	ListenAddr    string
	DBPath        string
	ActiveTimeout time.Duration
	APIKey        string
}

// Server represents the control plane HTTP server
type Server struct {
	config      *Config
	httpServer  *http.Server
	repo        store.PeerRepository
	peerService *service.PeerService
	authService *service.AuthService
}

// NewServer creates a new control plane server
func NewServer(config *Config) (*Server, error) {
	// Initialize repository
	repo, err := store.NewSQLiteRepository(config.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	// Initialize services
	peerService := service.NewPeerService(repo, config.ActiveTimeout)
	authService := service.NewAuthService(config.APIKey)

	// Create server
	server := &Server{
		config:      config,
		repo:        repo,
		peerService: peerService,
		authService: authService,
	}

	// Setup HTTP server
	mux := http.NewServeMux()
	server.setupRoutes(mux)

	server.httpServer = &http.Server{
		Addr:         config.ListenAddr,
		Handler:      server.corsMiddleware(server.loggingMiddleware(mux)),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server, nil
}

// setupRoutes registers all API routes
func (s *Server) setupRoutes(mux *http.ServeMux) {
	// API handlers
	registerHandler := api.NewRegisterHandler(s.peerService)
	peersHandler := api.NewPeersHandler(s.peerService)
	heartbeatHandler := api.NewHeartbeatHandler(s.peerService)
	metricsHandler := api.NewMetricsHandler(s.peerService)

	// Register routes
	mux.Handle("/register", registerHandler)
	mux.Handle("/peers", peersHandler)
	mux.Handle("/heartbeat", heartbeatHandler)
	mux.Handle("/metrics", metricsHandler)
	
	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

// loggingMiddleware logs all HTTP requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// corsMiddleware adds CORS headers for dashboard
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("Starting control plane server on %s", s.config.ListenAddr)
	log.Printf("Database: %s", s.config.DBPath)
	log.Printf("Active timeout: %s", s.config.ActiveTimeout)
	
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	
	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down control plane server...")
	
	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}
	
	// Close repository
	if err := s.repo.Close(); err != nil {
		return fmt.Errorf("failed to close repository: %w", err)
	}
	
	log.Println("Server shutdown complete")
	return nil
}
