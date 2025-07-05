package processors

import (
	"context"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/processors/esmtp"
)

// InputEvent represents an incoming event from external sources
type InputEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
	Priority  int                    `json:"priority"`
}

// InputEventResponse represents the response after processing an input event
type InputEventResponse struct {
	EventID     string                 `json:"event_id"`
	Processed   bool                   `json:"processed"`
	Actions     []OutputAction         `json:"actions"`
	Metadata    map[string]interface{} `json:"metadata"`
	ProcessedAt time.Time              `json:"processed_at"`
}

// OutputAction represents an action to be taken as a result of processing
type OutputAction struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Target    string                 `json:"target"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata"`
	Priority  int                    `json:"priority"`
	CreatedAt time.Time              `json:"created_at"`
}

// Manager manages all input processors
type Manager struct {
	config     *sharedconfig.Config
	processors map[string]InputProcessor
}

// InputProcessor defines the interface for input processors
type InputProcessor interface {
	Process(event *InputEvent) (*InputEventResponse, error)
	GetType() string
	IsEnabled() bool
}

// NewManager creates a new processor manager
func NewManager(cfg *sharedconfig.Config) (*Manager, error) {
	mgr := &Manager{
		config:     cfg,
		processors: make(map[string]InputProcessor),
	}

	// Initialize ESMTP processor if configured
	if cfg.ESMTP.Enabled {
		esmtpProcessor := esmtp.NewProcessor(&cfg.ESMTP)
		mgr.processors["esmtp"] = esmtpProcessor
	}

	// Initialize other processors here as they're implemented
	
	return mgr, nil
}

// ProcessEvent processes an input event using the appropriate processor
func (m *Manager) ProcessEvent(event *InputEvent) (*InputEventResponse, error) {
	processor, exists := m.processors[event.Type]
	if !exists {
		// Return a basic response for unknown types
		return &InputEventResponse{
			EventID:     event.ID,
			Processed:   false,
			Actions:     []OutputAction{},
			Metadata:    map[string]interface{}{"error": "no processor found for type: " + event.Type},
			ProcessedAt: time.Now(),
		}, nil
	}

	return processor.Process(event)
}

// RegisterProcessor registers a new processor
func (m *Manager) RegisterProcessor(processor InputProcessor) {
	m.processors[processor.GetType()] = processor
}

// GetProcessors returns all registered processors
func (m *Manager) GetProcessors() map[string]InputProcessor {
	return m.processors
}

// Start starts the processor manager
func (m *Manager) Start(ctx context.Context) error {
	return nil
}

// Stop stops the processor manager
func (m *Manager) Stop() error {
	return nil
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
