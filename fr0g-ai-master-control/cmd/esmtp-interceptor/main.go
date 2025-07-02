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
	flag.Parse()

	// Load configuration
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create ESMTP processor
	processor, err := input.NewESMTPProcessor(config)
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

	// ASCII art banner
	printBanner()

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

	log.Println("✅ fr0g.ai ESMTP Interceptor stopped")
}

func loadConfig(path string) (*input.ESMTPConfig, error) {
	// Default configuration
	config := &input.ESMTPConfig{
		Host:           "0.0.0.0",
		Port:           2525,
		TLSPort:        2465,
		Hostname:       "fr0g-ai-interceptor.local",
		MaxMessageSize: 10 * 1024 * 1024, // 10MB
		Timeout:        5 * time.Minute,
		MCPAddress:     "localhost:9092", // Master Control Protocol address
		EnableTLS:      false,
		CertFile:       "",
		KeyFile:        "",
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

func printBanner() {
	banner := `
	╔═══════════════════════════════════════════════════════════════╗
	║                                                               ║
	║    🐸 fr0g.ai - ESMTP Threat Vector Interceptor v1.0         ║
	║                                                               ║
	║    "Eliminating human-computer interaction vulnerabilities"   ║
	║                                                               ║
	║    📧 Email Intelligence Gathering: ACTIVE                   ║
	║    🛡️  Threat Analysis Engine: ONLINE                        ║
	║    🧠 Master Control Protocol: CONNECTED                     ║
	║                                                               ║
	╚═══════════════════════════════════════════════════════════════╝
	`
	log.Println(banner)
}
