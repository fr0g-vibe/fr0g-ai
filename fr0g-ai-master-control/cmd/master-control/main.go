package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-master-control/internal/mastercontrol"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/registry"
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

	// Initialize service registry client if enabled
	var registryClient *registry.RegistryClient
	if os.Getenv("SERVICE_REGISTRY_ENABLED") == "true" {
		registryURL := os.Getenv("SERVICE_REGISTRY_URL")
		if registryURL == "" {
			registryURL = "http://localhost:8500"
		}
		
		log.Printf("🔗 Initializing service registry client: %s", registryURL)
		registryClient = registry.NewRegistryClient(registryURL, nil)
		
		// Register service
		serviceName := os.Getenv("SERVICE_NAME")
		if serviceName == "" {
			serviceName = "fr0g-ai-master-control"
		}
		
		serviceID := os.Getenv("SERVICE_ID")
		if serviceID == "" {
			serviceID = serviceName + "-1"
		}
		
		httpPort := mcpConfig.Input.Webhook.Port
		grpcPort := 9093 // Default gRPC port for master control
		
		serviceInfo := &registry.ServiceInfo{
			ID:      serviceID,
			Name:    serviceName,
			Address: "localhost",
			Port:    httpPort,
			Tags:    []string{"ai", "master-control", "mcp"},
			Meta: map[string]string{
				"version":   "1.0.0",
				"http_port": strconv.Itoa(httpPort),
				"grpc_port": strconv.Itoa(grpcPort),
			},
			Check: &registry.HealthCheck{
				HTTP:     "http://localhost:" + strconv.Itoa(httpPort) + "/health",
				Interval: "30s",
				Timeout:  "10s",
			},
		}
		
		if err := registryClient.RegisterService(serviceInfo); err != nil {
			log.Printf("⚠️  Failed to register service: %v", err)
		} else {
			log.Printf("✅ Service registered successfully: %s", serviceID)
		}
	}

	log.Printf("COMPLETED Configuration loaded successfully")
	log.Printf("   - Learning Enabled: %v", mcpConfig.LearningEnabled)
	log.Printf("   - System Consciousness: %v", mcpConfig.SystemConsciousness)
	log.Printf("   - Emergent Capabilities: %v", mcpConfig.EmergentCapabilities)
	log.Printf("   - Max Concurrent Workflows: %d", mcpConfig.MaxConcurrentWorkflows)

	// Create Master Control Program
	log.Println("🧠 Initializing Master Control Program...")
	mcp := mastercontrol.NewMasterControlProgram(mcpConfig)

	// Start the MCP
	log.Println("STARTING Starting Master Control Program...")
	if err := mcp.Start(); err != nil {
		log.Fatalf("Failed to start MCP: %v", err)
	}

	log.Println("COMPLETED Master Control Program is now operational!")
	
	// Display system information
	systemState := mcp.GetSystemState()
	log.Printf("📊 System Status: %s", systemState.Status)
	log.Printf("📈 Active Workflows: %d", systemState.ActiveWorkflows)
	log.Printf("🧮 System Load: %.2f", systemState.SystemLoad)
	
	capabilities := mcp.GetCapabilities()
	log.Printf("TARGET System Capabilities: %d registered", len(capabilities))
	for id, cap := range capabilities {
		log.Printf("   - %s: %s (Emergent: %v)", id, cap.Name, cap.Emergent)
	}

	log.Printf("🧠 Intelligence Metrics:")
	log.Printf("   - Learning Rate: %.3f", systemState.Intelligence.LearningRate)
	log.Printf("   - Pattern Count: %d", systemState.Intelligence.PatternCount)
	log.Printf("   - Adaptation Score: %.3f", systemState.Intelligence.AdaptationScore)
	log.Printf("   - Efficiency Index: %.3f", systemState.Intelligence.EfficiencyIndex)
	log.Printf("   - Emergent Capabilities: %d", systemState.Intelligence.EmergentCapabilities)

	log.Printf("STARTING Master Control HTTP server ready on %s:%d", mcpConfig.Input.Webhook.Host, mcpConfig.Input.Webhook.Port)
	log.Printf("🔗 Health check: http://%s:%d/health", mcpConfig.Input.Webhook.Host, mcpConfig.Input.Webhook.Port)
	log.Printf("📊 Status endpoint: http://%s:%d/status", mcpConfig.Input.Webhook.Host, mcpConfig.Input.Webhook.Port)
	log.Printf("TARGET Discord webhook: http://%s:%d/webhook/discord", mcpConfig.Input.Webhook.Host, mcpConfig.Input.Webhook.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Master Control...")

	// Deregister from service registry
	if registryClient != nil {
		log.Println("🔗 Deregistering from service registry...")
		if err := registryClient.DeregisterService(); err != nil {
			log.Printf("⚠️  Failed to deregister service: %v", err)
		} else {
			log.Println("✅ Service deregistered successfully")
		}
	}

	// Graceful shutdown
	if err := mcp.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("COMPLETED Master Control stopped")
}
