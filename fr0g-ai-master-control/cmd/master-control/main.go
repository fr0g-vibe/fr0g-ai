package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fr0g-ai-master-control/internal/mastercontrol"
)


func main() {
	log.Println("🧠 Starting fr0g-ai-master-control...")

	// Load MCP configuration
	mcpConfig, err := mastercontrol.LoadMCPConfig("")
	if err != nil {
		log.Fatalf("Failed to load MCP config: %v", err)
	}

	// Override port from environment if set
	if port := os.Getenv("MCP_HTTP_PORT"); port != "" {
		log.Printf("Using MCP_HTTP_PORT: %s", port)
		// Note: Would need to parse port string to int if we want to override
		// For now, just log it
	}

	log.Printf("✅ Configuration loaded successfully")
	log.Printf("   - Learning Enabled: %v", mcpConfig.LearningEnabled)
	log.Printf("   - System Consciousness: %v", mcpConfig.SystemConsciousness)
	log.Printf("   - Emergent Capabilities: %v", mcpConfig.EmergentCapabilities)
	log.Printf("   - Max Concurrent Workflows: %d", mcpConfig.MaxConcurrentWorkflows)

	// Create Master Control Program
	log.Println("🧠 Initializing Master Control Program...")
	mcp := mastercontrol.NewMasterControlProgram(mcpConfig)

	// Start the MCP
	log.Println("🚀 Starting Master Control Program...")
	if err := mcp.Start(); err != nil {
		log.Fatalf("Failed to start MCP: %v", err)
	}

	log.Println("✅ Master Control Program is now operational!")
	
	// Display system information
	systemState := mcp.GetSystemState()
	log.Printf("📊 System Status: %s", systemState.Status)
	log.Printf("📈 Active Workflows: %d", systemState.ActiveWorkflows)
	log.Printf("🧮 System Load: %.2f", systemState.SystemLoad)
	
	capabilities := mcp.GetCapabilities()
	log.Printf("🎯 System Capabilities: %d registered", len(capabilities))
	for id, cap := range capabilities {
		log.Printf("   - %s: %s (Emergent: %v)", id, cap.Name, cap.Emergent)
	}

	log.Printf("🧠 Intelligence Metrics:")
	log.Printf("   - Learning Rate: %.3f", systemState.Intelligence.LearningRate)
	log.Printf("   - Pattern Count: %d", systemState.Intelligence.PatternCount)
	log.Printf("   - Adaptation Score: %.3f", systemState.Intelligence.AdaptationScore)
	log.Printf("   - Efficiency Index: %.3f", systemState.Intelligence.EfficiencyIndex)
	log.Printf("   - Emergent Capabilities: %d", systemState.Intelligence.EmergentCapabilities)

	log.Printf("🚀 Master Control HTTP server ready on %s:%d", mcpConfig.Input.Webhook.Host, mcpConfig.Input.Webhook.Port)
	log.Printf("🔗 Health check: http://%s:%d/health", mcpConfig.Input.Webhook.Host, mcpConfig.Input.Webhook.Port)
	log.Printf("📊 Status endpoint: http://%s:%d/status", mcpConfig.Input.Webhook.Host, mcpConfig.Input.Webhook.Port)
	log.Printf("🎯 Discord webhook: http://%s:%d/webhook/discord", mcpConfig.Input.Webhook.Host, mcpConfig.Input.Webhook.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Master Control...")

	// Graceful shutdown
	if err := mcp.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("✅ Master Control stopped")
}
