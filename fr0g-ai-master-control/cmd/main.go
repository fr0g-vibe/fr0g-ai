package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

func main() {
	log.Println("fr0g.ai Master Control Program")
	log.Println("===============================")

	// Load configuration using shared config system
	loader := sharedconfig.NewLoader(sharedconfig.LoaderOptions{
		ConfigPath: "",
		EnvPrefix:  "FR0G_MCP",
	})

	// Load environment files
	if err := loader.LoadEnvFiles(); err != nil {
		log.Printf("Warning: failed to load env files: %v", err)
	}

	// Create default config
	cfg := &sharedconfig.Config{
		HTTP: sharedconfig.HTTPConfig{
			Host: "0.0.0.0",
			Port: "8081",
		},
		GRPC: sharedconfig.GRPCConfig{
			Host:               "0.0.0.0", 
			Port:               "9091",
			MaxRecvMsgSize:     4 * 1024 * 1024, // 4MB
			MaxSendMsgSize:     4 * 1024 * 1024, // 4MB
			ConnectionTimeout:  30 * time.Second,
			EnableReflection:   false,
		},
		Storage: sharedconfig.StorageConfig{
			Type:    "file",
			DataDir: "/app/data",
		},
	}
	
	// Override with environment variables
	if httpPort := os.Getenv("HTTP_PORT"); httpPort != "" {
		cfg.HTTP.Port = httpPort
	}
	if storageType := os.Getenv("STORAGE_TYPE"); storageType != "" {
		cfg.Storage.Type = storageType
	}
	if dataDir := os.Getenv("DATA_DIR"); dataDir != "" {
		cfg.Storage.DataDir = dataDir
	}

	// Load from file
	if err := loader.LoadFromFile(cfg); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	log.Println("Configuration loaded and validated successfully")

	// Create MCP server
	log.Println("Initializing Master Control Program...")
	mcpServer := NewMCPServer(cfg)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the server
	log.Println("Starting Master Control Program...")
	if err := mcpServer.Start(ctx); err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
	}

	log.Println("Master Control Program is now operational!")
	log.Printf("   - HTTP Server: http://%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	log.Printf("   - Health Check: http://%s:%s/health", cfg.HTTP.Host, cfg.HTTP.Port)

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Master Control Program is running...")
	log.Println("   Press Ctrl+C to shutdown gracefully")

	// Wait for shutdown signal
	<-sigChan

	log.Println("Shutdown signal received...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Stop the server
	if err := mcpServer.Stop(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("Master Control Program shutdown complete")
}

// MCPServer represents the main MCP server
type MCPServer struct {
	config *sharedconfig.Config
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(cfg *sharedconfig.Config) *MCPServer {
	return &MCPServer{
		config: cfg,
	}
}

// Start starts the MCP server
func (s *MCPServer) Start(ctx context.Context) error {
	log.Println("Starting MCP server components...")

	// TODO: Implement actual HTTP server startup
	log.Printf("HTTP server starting on %s:%s", s.config.HTTP.Host, s.config.HTTP.Port)
	log.Printf("gRPC server would start on port %s", s.config.GRPC.Port)
	log.Println("Webhook input system would be initialized")

	// Log configuration
	log.Println("Master Control Program configuration:")
	log.Printf("   - HTTP: %s:%s", s.config.HTTP.Host, s.config.HTTP.Port)
	log.Printf("   - gRPC: %s:%s", s.config.GRPC.Host, s.config.GRPC.Port)
	log.Printf("   - Storage: %s (%s)", s.config.Storage.Type, s.config.Storage.DataDir)
	log.Println("   - Webhook endpoints: /webhook/{tag}")

	// TODO: Implement actual HTTP server startup
	log.Printf("HTTP server would listen on %s:%s", s.config.HTTP.Host, s.config.HTTP.Port)

	return nil
}

// Stop stops the MCP server
func (s *MCPServer) Stop(ctx context.Context) error {
	log.Println("Stopping MCP server...")

	// TODO: Implement actual server shutdown logic including webhook system
	log.Println("MCP server stopped")

	return nil
}
