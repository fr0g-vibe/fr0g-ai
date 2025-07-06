package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/api"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/config"
	grpcserver "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/grpc"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/persona"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/registry"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/storage"
	pb "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/grpc/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Load configuration using AIP-specific config
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if validationErrors := cfg.Validate(); validationErrors != nil {
		log.Fatalf("Configuration validation failed: %v", validationErrors)
	}

	log.Println("STARTING: fr0g.ai AIP servers...")

	// Initialize storage
	var store storage.Storage
	switch cfg.Storage.Type {
	case "file":
		store, err = storage.NewFileStorage(cfg.Storage.DataDir)
		if err != nil {
			log.Fatalf("Failed to initialize file storage: %v", err)
		}
	case "memory":
		store = storage.NewMemoryStorage()
	default:
		log.Fatalf("Unsupported storage type: %s", cfg.Storage.Type)
	}

	// Initialize persona service
	personaService := persona.NewService(store)

	// Initialize service registry client if configured
	var registryClient *registry.RegistryClient
	if registryURL := os.Getenv("REGISTRY_URL"); registryURL != "" {
		registryClient = registry.NewRegistryClient(registryURL, logger)
		
		// Register service
		serviceInfo := &registry.ServiceInfo{
			ID:      getEnvOrDefault("SERVICE_ID", "aip-001"),
			Name:    getEnvOrDefault("SERVICE_NAME", "fr0g-ai-aip"),
			Address: "localhost",
			Port:    getPortFromEnv("HTTP_PORT", 8080),
			Tags:    []string{"ai", "persona", "identity"},
			Check: &registry.HealthCheck{
				HTTP:     fmt.Sprintf("http://localhost:%s/health", cfg.HTTP.Port),
				Interval: "30s",
				Timeout:  "10s",
			},
		}
		
		if err := registryClient.RegisterService(serviceInfo); err != nil {
			logger.WithError(err).Warn("Failed to register with service registry")
		} else {
			logger.Info("Successfully registered with service registry")
		}
	}

	// Start gRPC server
	grpcServer := grpc.NewServer()
	grpcPersonaServer := grpcserver.NewServer()
	pb.RegisterPersonaServiceServer(grpcServer, grpcPersonaServer)

	grpcPort := cfg.GRPC.Port
	grpcListener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port %s: %v", grpcPort, err)
	}

	go func() {
		log.Printf("SUCCESS: gRPC server starting on port %s", grpcPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Start REST API server
	restServer := api.NewServer(cfg, personaService, registryClient)

	go func() {
		log.Printf("SUCCESS: REST API server starting on port %s", cfg.HTTP.Port)
		if err := restServer.Start(); err != nil && err != context.Canceled {
			log.Printf("REST server error: %v", err)
		}
	}()

	log.Println("RUNNING: fr0g.ai AIP is running...")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("SHUTTING DOWN: Servers...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown REST server
	if err := restServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("REST server shutdown error: %v", err)
	}

	// Shutdown gRPC server
	grpcServer.GracefulStop()

	log.Println("COMPLETE: fr0g.ai AIP shutdown complete")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getPortFromEnv(key string, defaultPort int) int {
	if portStr := os.Getenv(key); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			return port
		}
	}
	return defaultPort
}
