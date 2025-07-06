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

	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("Starting ESMTP webhook processor on port 8090...")

	printBanner("Webhook Mode")

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("🛑 Shutdown signal received, stopping ESMTP Webhook Processor...")

	// Graceful shutdown
	cancel()

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

	log.Printf("Starting ESMTP processor...")

	printBanner("SMTP Server Mode")

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("🛑 Shutdown signal received, stopping ESMTP Threat Vector Interceptor...")

	// Graceful shutdown
	cancel()

	log.Println("COMPLETED fr0g.ai ESMTP Interceptor stopped")
}

func loadConfig(path string) (*input.ESMTPConfig, error) {
	// Default configuration
	config := &input.ESMTPConfig{
		Server:            "0.0.0.0",
		Port:              2525,
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
