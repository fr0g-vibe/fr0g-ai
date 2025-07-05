package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-master-control/internal/discovery"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-master-control/internal/mastercontrol"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-master-control/internal/mastercontrol/input"
)

func main() {
	fmt.Println("🎛️  fr0g.ai Master Control Program Demo")
	fmt.Println("=====================================")

	// Initialize service discovery client
	discoveryClient := discovery.NewClient(&discovery.ClientConfig{
		RegistryURL:    "http://localhost:8500", // Default registry URL
		ServiceName:    "fr0g-ai-master-control",
		ServiceID:      "mcp-001",
		ServiceAddress: "localhost",
		ServicePort:    8081,
		HealthInterval: 30 * time.Second,
		Tags:           []string{"mcp", "master-control", "ai"},
		Metadata: map[string]string{
			"version": "1.0.0",
			"type":    "master-control",
		},
	})

	// Start service discovery
	if err := discoveryClient.Start(); err != nil {
		log.Printf("Warning: Failed to start service discovery: %v", err)
	}

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

	// Setup input processors (Discord and ESMTP)
	setupInputProcessors(mcp)

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

	// Stop service discovery
	if err := discoveryClient.Stop(); err != nil {
		log.Printf("Error stopping service discovery: %v", err)
	}

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

func setupInputProcessors(mcp *mastercontrol.MasterControlProgram) {
	fmt.Println("🔧 Setting up input processors...")

	// Create real AI client
	aiClientConfig := &input.AIClientConfig{
		AIPServiceURL:    "http://localhost:8080", // fr0g-ai-aip service
		BridgeServiceURL: "http://localhost:8082", // fr0g-ai-bridge service
		Timeout:          30 * time.Second,
	}
	aiClient := input.NewRealAIPersonaCommunityClient(aiClientConfig)

	// Discord processor configuration
	discordConfig := &input.DiscordProcessorConfig{
		CommunityTopic:    "general_discussion",
		PersonaCount:      3,
		ReviewTimeout:     time.Minute * 2,
		RequiredConsensus: 0.7,
		EnableSentiment:   true,
		FilterKeywords:    []string{},
	}

	// ESMTP processor configuration
	esmtpConfig := &input.ESMTPConfig{
		Server:            "0.0.0.0",
		Port:              2525,
		UseSSL:            false,
		Domain:            "fr0g-ai-interceptor.local",
		CommunityTopic:    "email-threat-analysis",
		PersonaCount:      5,
		ReviewTimeout:     2 * time.Minute,
		RequiredConsensus: 0.7,
	}

	// SMS processor configuration
	smsConfig := &input.SMSConfig{
		Provider:          "google_voice",
		AccountSID:        "demo-account-sid",
		AuthToken:         "demo-auth-token",
		WebhookURL:        "https://fr0g-ai.local/webhook/sms",
		PhoneNumber:       "+1-555-FR0G-AI",
		CommunityTopic:    "sms-threat-analysis",
		PersonaCount:      4,
		ReviewTimeout:     time.Minute * 2,
		RequiredConsensus: 0.75,
		TrustedNumbers:    []string{},
		BlockedNumbers:    []string{"spam", "telemarketer"},
	}

	// Voice processor configuration
	voiceConfig := &input.VoiceConfig{
		Provider:             "google_voice",
		APIKey:               "demo-api-key",
		APISecret:            "demo-api-secret",
		WebhookURL:           "https://fr0g-ai.local/webhook/voice",
		CommunityTopic:       "voice-threat-analysis",
		PersonaCount:         5,
		ReviewTimeout:        time.Minute * 3,
		RequiredConsensus:    0.8,
		TrustedNumbers:       []string{},
		BlockedNumbers:       []string{"robocaller", "spam"},
		MaxRecordingDuration: time.Minute * 10,
		SupportedFormats:     []string{"wav", "mp3"},
		TranscriptionEnabled: true,
		SentimentAnalysis:    true,
		VoiceprintAnalysis:   false,
		AudioStoragePath:     "/tmp/audio",
	}

	// IRC processor configuration
	ircConfig := &input.IRCConfig{
		Server:            "irc.libera.chat",
		Port:              6697,
		UseSSL:            true,
		Nickname:          "fr0g-ai-monitor",
		Username:          "fr0gai",
		RealName:          "fr0g.ai Security Monitor",
		Channels:          []string{"#security", "#ai", "#fr0g-ai"},
		CommunityTopic:    "irc-threat-analysis",
		PersonaCount:      4,
		ReviewTimeout:     time.Minute * 2,
		RequiredConsensus: 0.7,
		MonitorPrivateMsg: true,
		IgnoredNicks:      []string{"bot", "service"},
		TrustedNicks:      []string{"admin", "moderator"},
		AutoJoinChannels:  true,
	}

	// Webhook manager configuration
	webhookConfig := &input.WebhookConfig{
		Port:           8081,
		Host:           "0.0.0.0",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxRequestSize: 1024 * 1024, // 1MB
		EnableLogging:  true,
		AllowedOrigins: []string{"*"},
	}

	// Create input manager configuration
	inputManagerConfig := &input.InputManagerConfig{
		WebhookConfig: webhookConfig,
		Discord:       discordConfig,
		ESMTP:         esmtpConfig,
		SMS:           smsConfig,
		Voice:         voiceConfig,
		IRC:           ircConfig,
	}

	// Create and configure input manager
	inputManager, err := input.NewInputManager(inputManagerConfig, aiClient)
	if err != nil {
		log.Printf("Failed to create input manager: %v", err)
		return
	}

	// Set input manager in MCP
	mcp.SetInputManager(inputManager)

	fmt.Println("✅ Input processors configured (Discord + ESMTP + SMS + Voice + IRC)")
}
