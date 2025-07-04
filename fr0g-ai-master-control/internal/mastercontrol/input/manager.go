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
	SMS           *SMSConfig             `yaml:"sms"`
	Voice         *VoiceConfig           `yaml:"voice"`
	IRC           *IRCConfig             `yaml:"irc"`
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
	
	// Register SMS processor if enabled
	if im.config.SMS != nil {
		smsProcessor, err := NewSMSProcessor(im.config.SMS, im.aiClient)
		if err != nil {
			return fmt.Errorf("failed to create SMS processor: %w", err)
		}
		
		if err := im.webhookManager.RegisterProcessor(smsProcessor); err != nil {
			return fmt.Errorf("failed to register SMS processor: %w", err)
		}
		log.Printf("Input Manager: Registered SMS processor")
	}
	
	// Register Voice processor if enabled
	if im.config.Voice != nil {
		voiceProcessor, err := NewVoiceProcessor(im.config.Voice, im.aiClient)
		if err != nil {
			return fmt.Errorf("failed to create Voice processor: %w", err)
		}
		
		if err := im.webhookManager.RegisterProcessor(voiceProcessor); err != nil {
			return fmt.Errorf("failed to register Voice processor: %w", err)
		}
		log.Printf("Input Manager: Registered Voice processor")
	}
	
	// Register IRC processor if enabled
	if im.config.IRC != nil {
		ircProcessor, err := NewIRCProcessor(im.config.IRC, im.aiClient)
		if err != nil {
			return fmt.Errorf("failed to create IRC processor: %w", err)
		}
		
		if err := im.webhookManager.RegisterProcessor(ircProcessor); err != nil {
			return fmt.Errorf("failed to register IRC processor: %w", err)
		}
		log.Printf("Input Manager: Registered IRC processor")
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
			Server:            "0.0.0.0",
			Port:              2525,
			UseSSL:            false,
			Domain:            "fr0g-ai-interceptor.local",
			CommunityTopic:    "email-threat-analysis",
			PersonaCount:      5,
			ReviewTimeout:     2 * time.Minute,
			RequiredConsensus: 0.7,
		},
	}
}
