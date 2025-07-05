package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/api"
	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

func main() {
	// Load configuration with defaults
	cfg := sharedconfig.GetDefaults()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	log.Printf("Starting fr0g-ai-io service with config: HTTP=%s, gRPC=%s",
		cfg.HTTP.Address, cfg.GRPC.Address)

	// Create server
	server, err := api.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := server.Start(ctx); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down fr0g-ai-io service...")
	if err := server.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("fr0g-ai-io service stopped")
}
