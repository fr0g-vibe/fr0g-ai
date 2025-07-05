package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
	"github.com/fr0g-vibe/fr0g-ai/internal/server"
)

func main() {
	log.Println("🎛️  fr0g.ai Master Control Program")
	log.Println("==================================")
	
	// Load configuration using shared config system
	cfg := sharedconfig.LoadConfig()
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}
	
	log.Println("✅ Configuration loaded and validated successfully")
	
	// Create MCP server
	log.Println("🧠 Initializing Master Control Program...")
	mcpServer := NewMCPServer(cfg)
	
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Start the server
	log.Println("🚀 Starting Master Control Program...")
	if err := mcpServer.Start(ctx); err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
	}
	
	log.Println("✅ Master Control Program is now operational!")
	log.Printf("   - HTTP Server: http://localhost:%s", cfg.HTTP.Port)
	log.Printf("   - Health Check: http://localhost:%s/health", cfg.HTTP.Port)
	
	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	log.Println("🎯 Master Control Program is running...")
	log.Println("   Press Ctrl+C to shutdown gracefully")
	
	// Wait for shutdown signal
	<-sigChan
	
	log.Println("🛑 Shutdown signal received...")
	
	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	
	// Stop the server
	if err := mcpServer.Stop(shutdownCtx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
	
	log.Println("👋 Master Control Program shutdown complete")
}

// MCPServer represents the main MCP server
type MCPServer struct {
	config     *sharedconfig.Config
	httpServer *server.HTTPServer
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(cfg *sharedconfig.Config) *MCPServer {
	return &MCPServer{
		config:     cfg,
		httpServer: server.NewHTTPServer(cfg.HTTP.Port),
	}
}

// Start starts the MCP server
func (s *MCPServer) Start(ctx context.Context) error {
	log.Println("🔧 Starting MCP server components...")
	
	// Start HTTP server in a goroutine
	go func() {
		if err := s.httpServer.Start(); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()
	
	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)
	
	// Log enabled processors based on config
	log.Println("📡 Input processors status:")
	// Note: These config fields need to be added to shared config
	// For now, commenting out until config structure is updated
	/*
	if s.config.ESMTP.Enabled {
		log.Printf("📧 ESMTP processor enabled on port %s", s.config.ESMTP.Port)
	}
	if s.config.Discord.Enabled {
		log.Println("💬 Discord processor enabled")
	}
	if s.config.SMS.Enabled {
		log.Println("📱 SMS processor enabled")
	}
	if s.config.Voice.Enabled {
		log.Println("🎤 Voice processor enabled")
	}
	if s.config.IRC.Enabled {
		log.Println("💭 IRC processor enabled")
	}
	*/
	
	return nil
}

// Stop stops the MCP server
func (s *MCPServer) Stop(ctx context.Context) error {
	log.Println("🛑 Stopping MCP server...")
	
	if s.httpServer != nil {
		return s.httpServer.Stop()
	}
	
	return nil
}
