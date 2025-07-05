package discord

import (
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs/types"
	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// Processor handles Discord output operations
type Processor struct {
	config    *sharedconfig.Config
	isEnabled bool
}

// NewProcessor creates a new Discord output processor
func NewProcessor(cfg *sharedconfig.Config) *Processor {
	return &Processor{
		config:    cfg,
		isEnabled: false, // TODO: Enable when Discord config is available
	}
}

// Process executes a Discord output command
func (p *Processor) Process(command *types.OutputCommand) (*types.OutputResult, error) {
	if !p.isEnabled {
		return &types.OutputResult{
			CommandID:    command.ID,
			Success:      false,
			ErrorMessage: "Discord processor disabled - missing Discord configuration",
			Metadata:     map[string]string{"processor": "discord"},
			CompletedAt:  time.Now(),
		}, nil
	}

	// TODO: Implement actual Discord message sending
	return &types.OutputResult{
		CommandID:   command.ID,
		Success:     true,
		Metadata: map[string]string{
			"processor": "discord",
			"target":    command.Target,
		},
		CompletedAt: time.Now(),
	}, nil
}

// GetType returns the processor type
func (p *Processor) GetType() string {
	return "discord"
}

// IsEnabled returns whether the processor is enabled
func (p *Processor) IsEnabled() bool {
	return p.isEnabled
}
