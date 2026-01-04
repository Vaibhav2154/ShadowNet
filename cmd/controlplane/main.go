package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vaibhav2154/ShadowNet/internal/controlplane"
)

func main() {
	// Parse command-line flags
	listenAddr := flag.String("listen", getEnv("LISTEN_ADDR", ":8080"), "HTTP server listen address")
	dbPath := flag.String("db", getEnv("DB_PATH", "./data/controlplane.db"), "SQLite database path")
	activeTimeout := flag.Duration("active-timeout", 5*time.Minute, "Peer active timeout duration")
	apiKey := flag.String("api-key", getEnv("API_KEY", ""), "Optional API key for authentication")
	
	flag.Parse()

	// Create server configuration
	config := &controlplane.Config{
		ListenAddr:    *listenAddr,
		DBPath:        *dbPath,
		ActiveTimeout: *activeTimeout,
		APIKey:        *apiKey,
	}

	// Create server
	server, err := controlplane.NewServer(config)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	// Graceful shutdown
	log.Println("Received shutdown signal")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	log.Println("Control plane stopped")
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}