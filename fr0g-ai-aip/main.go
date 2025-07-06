package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/api"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/config"
	grpcserver "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/grpc"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/persona"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/registry"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/storage"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Initialize storage
	var store storage.Storage
	var err2 error
	switch cfg.Storage.Type {
	case "memory":
		store = storage.NewMemoryStorage()
	case "file":
		store, err2 = storage.NewFileStorage(cfg.Storage.DataDir)
		if err2 != nil {
			log.Fatalf("Failed to initialize file storage: %v", err2)
		}
	default:
		log.Fatalf("Unsupported storage type: %s", cfg.Storage.Type)
	}

	// Initialize persona service
	personaService := persona.NewService(store)

	// Initialize registry client
	var registryClient *registry.RegistryClient
	if os.Getenv("ENABLE_REGISTRY") != "false" {
		registryClient, err = registry.NewRegistryClientSimple("http://localhost:8500")
		if err != nil {
			log.Printf("Warning: Failed to create registry client: %v", err)
		} else {
			// Register service
			if err := registryClient.RegisterServiceSimple("fr0g-ai-aip", cfg.HTTP.Port, cfg.GRPC.Port); err != nil {
				log.Printf("Warning: Failed to register service: %v", err)
			}
		}
	}

	// Start gRPC server
	go func() {
		log.Printf("Starting gRPC server on port %s", cfg.GRPC.Port)
		if err := grpcserver.StartGRPCServerWithConfig(cfg, personaService); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Start REST API server
	restServer := api.NewServer(cfg, personaService, registryClient)
	
	go func() {
		log.Printf("Starting REST API server on port %s", cfg.HTTP.Port)
		if err := restServer.Start(); err != nil {
			log.Printf("REST server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down servers...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown REST server
	if err := restServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("REST server shutdown error: %v", err)
	}

	log.Println("Servers stopped")
}
