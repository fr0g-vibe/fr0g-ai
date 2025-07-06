package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-master-control/internal/mastercontrol/input"
	"gopkg.in/yaml.v2"
)

func main() {
	var configPath = flag.String("config", "configs/esmtp.yaml", "Path to configuration file")
	var webhookMode = flag.Bool("webhook", false, "Run as webhook processor instead of SMTP server")
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if *webhookMode {
		runWebhookMode(config)
	} else {
		runSMTPMode(config)
	}
}

func runWebhookMode(config *input.ESMTPConfig) {
	// Create AI community client
	aiClient := &input.MockAIPersonaCommunityClient{} // Use mock for now

	// Create ESMTP processor as webhook processor
	processor, err := input.NewESMTPProcessor(config, aiClient)
	if err != nil {
		log.Fatalf("Failed to create ESMTP processor: %v", err)
	}

	// Create webhook manager
	webhookConfig := &input.WebhookConfig{
		Port:           8090,
		Host:           "0.0.0.0",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxRequestSize: 10 * 1024 * 1024,
		EnableLogging:  true,
		AllowedOrigins: []string{"*"},
	}

	manager, err := input.NewWebhookManager(webhookConfig)
	if err != nil {
		log.Fatalf("Failed to create webhook manager: %v", err)
	}

	// Register ESMTP processor
	if err := manager.RegisterProcessor(processor); err != nil {
		log.Fatalf("Failed to register ESMTP processor: %v", err)
	}

	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start webhook manager
	if err := manager.Start(ctx); err != nil {
		log.Fatalf("Failed to start webhook manager: %v", err)
	}

	printBanner("Webhook Mode")

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("🛑 Shutdown signal received, stopping ESMTP Webhook Processor...")

	// Graceful shutdown
	cancel()
	if err := manager.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("COMPLETED fr0g.ai ESMTP Webhook Processor stopped")
}

func runSMTPMode(config *input.ESMTPConfig) {
	// Create AI community client
	aiClient := &input.MockAIPersonaCommunityClient{} // Use mock for now

	// Create ESMTP processor
	processor, err := input.NewESMTPProcessor(config, aiClient)
	if err != nil {
		log.Fatalf("Failed to create ESMTP processor: %v", err)
	}

	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the processor
	if err := processor.Start(ctx); err != nil {
		log.Fatalf("Failed to start ESMTP processor: %v", err)
	}

	printBanner("SMTP Server Mode")

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("🛑 Shutdown signal received, stopping ESMTP Threat Vector Interceptor...")

	// Graceful shutdown
	cancel()
	if err := processor.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("COMPLETED fr0g.ai ESMTP Interceptor stopped")
}

func loadConfig(path string) (*input.ESMTPConfig, error) {
	// Default configuration
	config := &input.ESMTPConfig{
		Host:              "0.0.0.0",
		Port:              2525,
		TLSPort:           2465,
		Hostname:          "fr0g-ai-interceptor.local",
		MaxMessageSize:    10 * 1024 * 1024, // 10MB
		Timeout:           5 * time.Minute,
		EnableTLS:         false,
		CertFile:          "",
		KeyFile:           "",
		CommunityTopic:    "email-threat-analysis",
		PersonaCount:      5,
		ReviewTimeout:     2 * time.Minute,
		RequiredConsensus: 0.7,
	}

	// Try to load from file
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

func printBanner(mode string) {
	banner := `
	╔═══════════════════════════════════════════════════════════════╗
	║                                                               ║
	║    fr0g.ai fr0g.ai - ESMTP Threat Vector Interceptor v1.0         ║
	║                                                               ║
	║    "Eliminating human-computer interaction vulnerabilities"   ║
	║                                                               ║
	║    📧 Email Intelligence Gathering: ACTIVE                   ║
	║    🛡️  Threat Analysis Engine: ONLINE                        ║
	║    🧠 AI Community Review: CONNECTED                         ║
	║    🔧 Mode: %-50s║
	║                                                               ║
	╚═══════════════════════════════════════════════════════════════╝
	`
	log.Printf(banner, mode)
}
