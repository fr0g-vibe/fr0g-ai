package input

import (
	"context"
	"fmt"
	"log"
	"time"
)

// InputManager manages all input processors (Discord, ESMTP, etc.)
type InputManager struct {
	webhookManager *WebhookManager
	aiClient       AIPersonaCommunityClient
	config         *InputManagerConfig
}

// InputManagerConfig holds configuration for the input manager
type InputManagerConfig struct {
	WebhookConfig *WebhookConfig         `yaml:"webhook"`
	Discord       *DiscordProcessorConfig `yaml:"discord"`
	ESMTP         *ESMTPConfig           `yaml:"esmtp"`
}

// NewInputManager creates a new input manager
func NewInputManager(config *InputManagerConfig, aiClient AIPersonaCommunityClient) (*InputManager, error) {
	// Create webhook manager
	webhookManager := NewWebhookManager(config.WebhookConfig)
	
	return &InputManager{
		webhookManager: webhookManager,
		aiClient:       aiClient,
		config:         config,
	}, nil
}

// Start initializes and starts all input processors
func (im *InputManager) Start(ctx context.Context) error {
	log.Println("Input Manager: Starting input management...")
	
	// Start webhook manager
	if err := im.webhookManager.Start(); err != nil {
		return fmt.Errorf("failed to start webhook manager: %w", err)
	}
	
	// Register Discord processor if enabled
	if im.config.Discord != nil {
		discordProcessor := NewDiscordWebhookProcessor(im.aiClient, im.config.Discord)
		if err := im.webhookManager.RegisterProcessor(discordProcessor); err != nil {
			return fmt.Errorf("failed to register Discord processor: %w", err)
		}
		log.Printf("Input Manager: Registered Discord processor")
	}
	
	// Register ESMTP processor if enabled
	if im.config.ESMTP != nil {
		esmtpProcessor, err := NewESMTPProcessor(im.config.ESMTP, im.aiClient)
		if err != nil {
			return fmt.Errorf("failed to create ESMTP processor: %w", err)
		}
		
		if err := im.webhookManager.RegisterProcessor(esmtpProcessor); err != nil {
			return fmt.Errorf("failed to register ESMTP processor: %w", err)
		}
		log.Printf("Input Manager: Registered ESMTP processor")
	}
	
	log.Println("Input Manager: Input management started successfully")
	return nil
}

// Stop gracefully stops all input processors
func (im *InputManager) Stop() error {
	log.Println("Input Manager: Stopping input management...")
	
	if err := im.webhookManager.Stop(); err != nil {
		return fmt.Errorf("failed to stop webhook manager: %w", err)
	}
	
	log.Println("Input Manager: Input management stopped")
	return nil
}

// GetWebhookManager returns the webhook manager for external access
func (im *InputManager) GetWebhookManager() *WebhookManager {
	return im.webhookManager
}

// DefaultInputManagerConfig returns a default input manager configuration
func DefaultInputManagerConfig() *InputManagerConfig {
	return &InputManagerConfig{
		WebhookConfig: &WebhookConfig{
			Port:           8081,
			Host:           "0.0.0.0",
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxRequestSize: 1024 * 1024, // 1MB
			EnableLogging:  true,
			AllowedOrigins: []string{"*"},
		},
		Discord: &DiscordProcessorConfig{
			CommunityTopic:    "general_discussion",
			PersonaCount:      3,
			ReviewTimeout:     60 * time.Second,
			RequiredConsensus: 0.7,
			EnableSentiment:   true,
			FilterKeywords:    []string{"spam", "inappropriate"},
		},
		ESMTP: &ESMTPConfig{
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
		},
	}
}
