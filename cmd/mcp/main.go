package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/internal/config"
	"github.com/fr0g-vibe/fr0g-ai/internal/server"
)

func main() {
	fmt.Println("🎛️  fr0g.ai Master Control Program")
	fmt.Println("==================================")
	
	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	fmt.Printf("✅ Configuration loaded successfully\n")
	fmt.Printf("   - HTTP Port: %d\n", cfg.HTTPPort)
	fmt.Printf("   - Log Level: %s\n", cfg.LogLevel)
	fmt.Printf("   - Learning Enabled: %v\n", cfg.LearningEnabled)
	fmt.Printf("   - System Consciousness: %v\n", cfg.SystemConsciousness)
	fmt.Printf("   - Emergent Capabilities: %v\n", cfg.EmergentCapabilities)
	fmt.Println()
	
	// Create MCP server
	fmt.Println("🧠 Initializing Master Control Program...")
	mcpServer := NewMCPServer(cfg)
	
	// Start the server
	fmt.Println("🚀 Starting Master Control Program...")
	if err := mcpServer.Start(); err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
	}
	
	fmt.Println("✅ Master Control Program is now operational!")
	fmt.Printf("   - HTTP Server: http://localhost:%d\n", cfg.HTTPPort)
	fmt.Printf("   - Health Check: http://localhost:%d/health\n", cfg.HTTPPort)
	fmt.Println()
	
	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	fmt.Println("🎯 Master Control Program is running...")
	fmt.Println("   Press Ctrl+C to shutdown gracefully")
	fmt.Println()
	
	// Wait for shutdown signal
	<-c
	
	fmt.Println("\n🛑 Shutdown signal received...")
	
	// Stop the server
	if err := mcpServer.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
	
	fmt.Println("👋 Master Control Program shutdown complete")
}

// MCPServer represents the main MCP server
type MCPServer struct {
	config     *config.Config
	httpServer *server.HTTPServer
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(cfg *config.Config) *MCPServer {
	return &MCPServer{
		config:     cfg,
		httpServer: server.NewHTTPServer(cfg.HTTPPort),
	}
}

// Start starts the MCP server
func (s *MCPServer) Start() error {
	fmt.Println("🔧 Starting MCP server components...")
	
	// Start HTTP server in a goroutine
	go func() {
		if err := s.httpServer.Start(); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()
	
	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)
	
	// Initialize input processors
	if s.config.ESMTPEnabled {
		fmt.Printf("📧 ESMTP processor enabled on port %d\n", s.config.ESMTPPort)
	}
	if s.config.DiscordEnabled {
		fmt.Println("💬 Discord processor enabled")
	}
	if s.config.SMSEnabled {
		fmt.Println("📱 SMS processor enabled")
	}
	if s.config.VoiceEnabled {
		fmt.Println("🎤 Voice processor enabled")
	}
	if s.config.IRCEnabled {
		fmt.Println("💭 IRC processor enabled")
	}
	
	return nil
}

// Stop stops the MCP server
func (s *MCPServer) Stop() error {
	fmt.Println("🛑 Stopping MCP server...")
	
	if s.httpServer != nil {
		return s.httpServer.Stop()
	}
	
	return nil
}
