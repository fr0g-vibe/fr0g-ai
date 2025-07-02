package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fr0g-ai-bridge/internal/mastercontrol"
)

func main() {
	fmt.Println("🎛️  fr0g.ai Master Control Program Demo")
	fmt.Println("=====================================")
	
	// Load MCP configuration
	config, err := mastercontrol.LoadMCPConfig("")
	if err != nil {
		log.Fatalf("Failed to load MCP config: %v", err)
	}
	
	fmt.Printf("✅ Configuration loaded successfully\n")
	fmt.Printf("   - Learning Enabled: %v\n", config.LearningEnabled)
	fmt.Printf("   - System Consciousness: %v\n", config.SystemConsciousness)
	fmt.Printf("   - Emergent Capabilities: %v\n", config.EmergentCapabilities)
	fmt.Printf("   - Max Concurrent Workflows: %d\n", config.MaxConcurrentWorkflows)
	fmt.Println()
	
	// Create Master Control Program
	fmt.Println("🧠 Initializing Master Control Program...")
	mcp := mastercontrol.NewMasterControlProgram(config)
	
	// Start the MCP
	fmt.Println("🚀 Starting Master Control Program...")
	if err := mcp.Start(); err != nil {
		log.Fatalf("Failed to start MCP: %v", err)
	}
	
	fmt.Println("✅ Master Control Program is now operational!")
	fmt.Println()
	
	// Demonstrate MCP functionality
	demonstrateMCPFunctionality(mcp)
	
	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	fmt.Println("🎯 Master Control Program is running...")
	fmt.Println("   Press Ctrl+C to shutdown gracefully")
	fmt.Println()
	
	// Wait for shutdown signal
	<-c
	
	fmt.Println("\n🛑 Shutdown signal received...")
	
	// Stop the MCP
	if err := mcp.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
	
	fmt.Println("👋 Master Control Program shutdown complete")
}

func demonstrateMCPFunctionality(mcp *mastercontrol.MasterControlProgram) {
	fmt.Println("🔍 Demonstrating MCP Functionality:")
	fmt.Println("-----------------------------------")
	
	// Wait a moment for components to initialize
	time.Sleep(time.Second * 2)
	
	// Get system state
	systemState := mcp.GetSystemState()
	fmt.Printf("📊 System Status: %s\n", systemState.Status)
	fmt.Printf("📈 Active Workflows: %d\n", systemState.ActiveWorkflows)
	fmt.Printf("🧮 System Load: %.2f\n", systemState.SystemLoad)
	fmt.Printf("🕒 Last Update: %s\n", systemState.LastUpdate.Format("15:04:05"))
	fmt.Println()
	
	// Get capabilities
	capabilities := mcp.GetCapabilities()
	fmt.Printf("🎯 System Capabilities: %d registered\n", len(capabilities))
	for id, cap := range capabilities {
		fmt.Printf("   - %s: %s (Emergent: %v)\n", id, cap.Name, cap.Emergent)
	}
	fmt.Println()
	
	// Intelligence metrics
	fmt.Printf("🧠 Intelligence Metrics:\n")
	fmt.Printf("   - Learning Rate: %.3f\n", systemState.Intelligence.LearningRate)
	fmt.Printf("   - Pattern Count: %d\n", systemState.Intelligence.PatternCount)
	fmt.Printf("   - Adaptation Score: %.3f\n", systemState.Intelligence.AdaptationScore)
	fmt.Printf("   - Efficiency Index: %.3f\n", systemState.Intelligence.EfficiencyIndex)
	fmt.Printf("   - Emergent Capabilities: %d\n", systemState.Intelligence.EmergentCapabilities)
	fmt.Println()
	
	fmt.Println("🎭 The MCP will now demonstrate its consciousness...")
	fmt.Println("   Watch the logs for cognitive reflections and system awareness updates")
	fmt.Println("   The system will continuously learn and adapt while running")
	fmt.Println()
}
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fr0g-ai-bridge/internal/mastercontrol"
)

func main() {
	fmt.Println("🎛️  fr0g.ai Master Control Program Demo")
	fmt.Println("=====================================")
	
	// Load MCP configuration
	config, err := mastercontrol.LoadMCPConfig("")
	if err != nil {
		log.Fatalf("Failed to load MCP config: %v", err)
	}
	
	fmt.Printf("✅ Configuration loaded successfully\n")
	fmt.Printf("   - Learning Enabled: %v\n", config.LearningEnabled)
	fmt.Printf("   - System Consciousness: %v\n", config.SystemConsciousness)
	fmt.Printf("   - Emergent Capabilities: %v\n", config.EmergentCapabilities)
	fmt.Printf("   - Max Concurrent Workflows: %d\n", config.MaxConcurrentWorkflows)
	fmt.Println()
	
	// Create Master Control Program
	fmt.Println("🧠 Initializing Master Control Program...")
	mcp := mastercontrol.NewMasterControlProgram(config)
	
	// Start the MCP
	fmt.Println("🚀 Starting Master Control Program...")
	if err := mcp.Start(); err != nil {
		log.Fatalf("Failed to start MCP: %v", err)
	}
	
	fmt.Println("✅ Master Control Program is now operational!")
	fmt.Println()
	
	// Demonstrate MCP functionality
	demonstrateMCPFunctionality(mcp)
	
	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	fmt.Println("🎯 Master Control Program is running...")
	fmt.Println("   Press Ctrl+C to shutdown gracefully")
	fmt.Println()
	
	// Wait for shutdown signal
	<-c
	
	fmt.Println("\n🛑 Shutdown signal received...")
	
	// Stop the MCP
	if err := mcp.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
	
	fmt.Println("👋 Master Control Program shutdown complete")
}

func demonstrateMCPFunctionality(mcp *mastercontrol.MasterControlProgram) {
	fmt.Println("🔍 Demonstrating MCP Functionality:")
	fmt.Println("-----------------------------------")
	
	// Wait a moment for components to initialize
	time.Sleep(time.Second * 2)
	
	// Get system state
	systemState := mcp.GetSystemState()
	fmt.Printf("📊 System Status: %s\n", systemState.Status)
	fmt.Printf("📈 Active Workflows: %d\n", systemState.ActiveWorkflows)
	fmt.Printf("🧮 System Load: %.2f\n", systemState.SystemLoad)
	fmt.Printf("🕒 Last Update: %s\n", systemState.LastUpdate.Format("15:04:05"))
	fmt.Println()
	
	// Get capabilities
	capabilities := mcp.GetCapabilities()
	fmt.Printf("🎯 System Capabilities: %d registered\n", len(capabilities))
	for id, cap := range capabilities {
		fmt.Printf("   - %s: %s (Emergent: %v)\n", id, cap.Name, cap.Emergent)
	}
	fmt.Println()
	
	// Intelligence metrics
	fmt.Printf("🧠 Intelligence Metrics:\n")
	fmt.Printf("   - Learning Rate: %.3f\n", systemState.Intelligence.LearningRate)
	fmt.Printf("   - Pattern Count: %d\n", systemState.Intelligence.PatternCount)
	fmt.Printf("   - Adaptation Score: %.3f\n", systemState.Intelligence.AdaptationScore)
	fmt.Printf("   - Efficiency Index: %.3f\n", systemState.Intelligence.EfficiencyIndex)
	fmt.Printf("   - Emergent Capabilities: %d\n", systemState.Intelligence.EmergentCapabilities)
	fmt.Println()
	
	fmt.Println("🎭 The MCP will now demonstrate its consciousness...")
	fmt.Println("   Watch the logs for cognitive reflections and system awareness updates")
	fmt.Println("   The system will continuously learn and adapt while running")
	fmt.Println()
}
