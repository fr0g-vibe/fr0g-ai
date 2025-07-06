package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-master-control/internal/mastercontrol"
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

	// Create MCP service
	log.Println("Initializing Master Control Program...")
	
	// Convert shared config to MCP config
	mcpConfig := &mastercontrol.MCPConfig{
		Input: mastercontrol.InputConfig{
			Webhook: mastercontrol.WebhookConfig{
				Host:         cfg.HTTP.Host,
				Port:         8082, // Use port 8082 for MCP
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,
			},
		},
		AdaptiveLearningRate:     0.154,
		CognitiveReflectionRate:  30 * time.Second,
		SystemConsciousness:      true,
	}
	
	mcp := mastercontrol.NewMasterControlProgram(mcpConfig)

	// Start the MCP service
	log.Println("Starting Master Control Program...")
	if err := mcp.Start(); err != nil {
		log.Fatalf("Failed to start MCP: %v", err)
	}

	log.Println("Master Control Program is now operational!")
	log.Printf("   - HTTP Server: http://%s:%d", mcpConfig.Input.Webhook.Host, mcpConfig.Input.Webhook.Port)
	log.Printf("   - Health Check: http://%s:%d/health", mcpConfig.Input.Webhook.Host, mcpConfig.Input.Webhook.Port)

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Master Control Program is running...")
	log.Println("   Press Ctrl+C to shutdown gracefully")

	// Wait for shutdown signal
	<-sigChan

	log.Println("Shutdown signal received...")

	// Stop the MCP service
	if err := mcp.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("Master Control Program shutdown complete")
}
