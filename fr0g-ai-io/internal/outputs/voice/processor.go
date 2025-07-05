package voice

import (
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs/types"
	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// Processor handles Voice output operations
type Processor struct {
	config    *sharedconfig.Config
	isEnabled bool
}

// NewProcessor creates a new Voice output processor
func NewProcessor(cfg *sharedconfig.Config) *Processor {
	return &Processor{
		config:    cfg,
		isEnabled: false, // TODO: Enable when Voice config is available
	}
}

// Process executes a Voice output command
func (p *Processor) Process(command *types.OutputCommand) (*types.OutputResult, error) {
	if !p.isEnabled {
		return &types.OutputResult{
			CommandID:    command.ID,
			Success:      false,
			ErrorMessage: "Voice processor disabled - missing Voice configuration",
			Metadata:     map[string]string{"processor": "voice"},
			CompletedAt:  time.Now(),
		}, nil
	}

	// TODO: Implement actual Voice call/message sending
	return &types.OutputResult{
		CommandID:   command.ID,
		Success:     true,
		Metadata: map[string]string{
			"processor": "voice",
			"target":    command.Target,
		},
		CompletedAt: time.Now(),
	}, nil
}

// GetType returns the processor type
func (p *Processor) GetType() string {
	return "voice"
}

// IsEnabled returns whether the processor is enabled
func (p *Processor) IsEnabled() bool {
	return p.isEnabled
}
