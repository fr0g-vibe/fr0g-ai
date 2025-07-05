package outputs

import (
	"context"
	"fmt"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs/discord"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs/irc"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs/sms"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs/types"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs/voice"
)

// Manager manages all output processors
type Manager struct {
	config     *sharedconfig.Config
	processors map[string]types.OutputProcessor
}

// NewManager creates a new output manager
func NewManager(cfg *sharedconfig.Config) (*Manager, error) {
	mgr := &Manager{
		config:     cfg,
		processors: make(map[string]types.OutputProcessor),
	}

	// Initialize SMS processor
	smsProcessor := sms.NewProcessor(cfg)
	mgr.processors["sms"] = smsProcessor

	// Initialize IRC processor
	ircProcessor := irc.NewProcessor(cfg)
	mgr.processors["irc"] = ircProcessor

	// Initialize Discord processor
	discordProcessor := discord.NewProcessor(cfg)
	mgr.processors["discord"] = discordProcessor

	// Initialize Voice processor
	voiceProcessor := voice.NewProcessor(cfg)
	mgr.processors["voice"] = voiceProcessor

	// SMTP processor will be added when the package is created

	return mgr, nil
}

// Start starts the output manager
func (m *Manager) Start(ctx context.Context) error {
	return nil
}

// Stop stops the output manager
func (m *Manager) Stop() error {
	return nil
}

// ExecuteCommand executes an output command using the appropriate processor
func (m *Manager) ExecuteCommand(command *types.OutputCommand) (*types.OutputResult, error) {
	processor, exists := m.processors[command.Type]
	if !exists {
		return &types.OutputResult{
			CommandID:    command.ID,
			Success:      false,
			ErrorMessage: fmt.Sprintf("no processor found for type: %s", command.Type),
			Metadata:     map[string]string{"error": "unknown_processor_type"},
			CompletedAt:  time.Now(),
		}, nil
	}

	return processor.Process(command)
}

// RegisterProcessor registers a new output processor
func (m *Manager) RegisterProcessor(processor types.OutputProcessor) {
	m.processors[processor.GetType()] = processor
}

// GetProcessors returns all registered processors
func (m *Manager) GetProcessors() map[string]types.OutputProcessor {
	return m.processors
}

// GetStatus returns the status of all processors
func (m *Manager) GetStatus() map[string]interface{} {
	status := map[string]interface{}{
		"total_processors": len(m.processors),
		"processors":       make(map[string]interface{}),
	}

	for name, processor := range m.processors {
		status["processors"].(map[string]interface{})[name] = map[string]interface{}{
			"type":    processor.GetType(),
			"enabled": processor.IsEnabled(),
		}
	}

	return status
}
