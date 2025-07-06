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
	
	// Override with environment variables for Docker deployment
	if httpPort := os.Getenv("HTTP_PORT"); httpPort != "" {
		cfg.HTTP.Port = httpPort
	} else {
		cfg.HTTP.Port = "8083"
	}
	
	if grpcPort := os.Getenv("GRPC_PORT"); grpcPort != "" {
		cfg.GRPC.Port = grpcPort
	} else {
		cfg.GRPC.Port = "9092"
	}
	
	// Set host to listen on all interfaces in container
	cfg.HTTP.Host = "0.0.0.0"
	cfg.GRPC.Host = "0.0.0.0"

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	log.Printf("Starting fr0g-ai-io service with config: HTTP=%s:%s, gRPC=%s:%s",
		cfg.HTTP.Host, cfg.HTTP.Port, cfg.GRPC.Host, cfg.GRPC.Port)

	// Create HTTP API server
	server, err := api.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// TODO: Add gRPC server once protobuf generation is fixed
	log.Println("gRPC server temporarily disabled - protobuf generation needed")

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
