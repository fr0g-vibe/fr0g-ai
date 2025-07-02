package input

import (
	"context"
	"log"
	"time"
)

// InputManager manages all input sources for the MCP
type InputManager struct {
	webhookManager *WebhookManager
	processors     map[string]WebhookProcessor
	config         *InputConfig
	ctx            context.Context
	cancel         context.CancelFunc
}

// InputConfig holds configuration for the input manager
type InputConfig struct {
	Webhook WebhookConfig `yaml:"webhook"`
	Discord DiscordProcessorConfig `yaml:"discord"`
}

// NewInputManager creates a new input manager
func NewInputManager(config *InputConfig) *InputManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &InputManager{
		webhookManager: NewWebhookManager(&config.Webhook),
		processors:     make(map[string]WebhookProcessor),
		config:         config,
		ctx:            ctx,
		cancel:         cancel,
	}
}

// Start begins input manager operation
func (im *InputManager) Start() error {
	log.Println("Input Manager: Starting input management...")
	
	// Start webhook manager
	if err := im.webhookManager.Start(); err != nil {
		return err
	}
	
	// Register default processors
	im.registerDefaultProcessors()
	
	log.Println("Input Manager: Input management started successfully")
	return nil
}

// Stop gracefully stops the input manager
func (im *InputManager) Stop() error {
	log.Println("Input Manager: Stopping input management...")
	
	im.cancel()
	
	if err := im.webhookManager.Stop(); err != nil {
		return err
	}
	
	log.Println("Input Manager: Input management stopped")
	return nil
}

// RegisterProcessor registers a webhook processor
func (im *InputManager) RegisterProcessor(processor WebhookProcessor) error {
	tag := processor.GetTag()
	im.processors[tag] = processor
	return im.webhookManager.RegisterProcessor(processor)
}

// GetWebhookManager returns the webhook manager
func (im *InputManager) GetWebhookManager() *WebhookManager {
	return im.webhookManager
}

// GetProcessors returns all registered processors
func (im *InputManager) GetProcessors() map[string]string {
	return im.webhookManager.GetRegisteredProcessors()
}

// registerDefaultProcessors registers the default webhook processors
func (im *InputManager) registerDefaultProcessors() {
	// Register Discord processor with mock client
	mockClient := NewMockAIPersonaCommunityClient()
	discordProcessor := NewDiscordWebhookProcessor(mockClient, &im.config.Discord)
	
	if err := im.RegisterProcessor(discordProcessor); err != nil {
		log.Printf("Input Manager: Failed to register Discord processor: %v", err)
	}
	
	// Additional processors can be registered here
}

// DefaultInputConfig returns a default input configuration
func DefaultInputConfig() *InputConfig {
	return &InputConfig{
		Webhook: WebhookConfig{
			Port:           8081,
			Host:           "0.0.0.0",
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxRequestSize: 1024 * 1024, // 1MB
			EnableLogging:  true,
			AllowedOrigins: []string{"*"},
		},
		Discord: DiscordProcessorConfig{
			CommunityTopic:    "general_discussion",
			PersonaCount:      3,
			ReviewTimeout:     60 * time.Second,
			RequiredConsensus: 0.7,
			EnableSentiment:   true,
			FilterKeywords:    []string{"spam", "inappropriate"},
		},
	}
}
