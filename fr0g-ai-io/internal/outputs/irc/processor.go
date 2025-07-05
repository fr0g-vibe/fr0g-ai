package irc

import (
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs/types"
	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// Processor handles IRC output operations
type Processor struct {
	config    *sharedconfig.Config
	isEnabled bool
}

// NewProcessor creates a new IRC output processor
func NewProcessor(cfg *sharedconfig.Config) *Processor {
	return &Processor{
		config:    cfg,
		isEnabled: false, // TODO: Enable when IRC config is available
	}
}

// Process executes an IRC output command
func (p *Processor) Process(command *types.OutputCommand) (*types.OutputResult, error) {
	if !p.isEnabled {
		return &types.OutputResult{
			CommandID:    command.ID,
			Success:      false,
			ErrorMessage: "IRC processor disabled - missing IRC configuration",
			Metadata:     map[string]string{"processor": "irc"},
			CompletedAt:  time.Now(),
		}, nil
	}

	// TODO: Implement actual IRC message sending
	return &types.OutputResult{
		CommandID:   command.ID,
		Success:     true,
		Metadata: map[string]string{
			"processor": "irc",
			"target":    command.Target,
		},
		CompletedAt: time.Now(),
	}, nil
}

// GetType returns the processor type
func (p *Processor) GetType() string {
	return "irc"
}

// IsEnabled returns whether the processor is enabled
func (p *Processor) IsEnabled() bool {
	return p.isEnabled
}
