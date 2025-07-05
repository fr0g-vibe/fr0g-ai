package processors

import (
	"context"
	"fmt"
	"log"
	"sync"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/processors/sms"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/processors/voice"
)

// Manager manages all I/O processors
type Manager struct {
	config     *sharedconfig.Config
	processors map[string]Processor
	mu         sync.RWMutex
	isRunning  bool
}

// Processor defines the interface for I/O processors
type Processor interface {
	Start(ctx context.Context) error
	Stop() error
	GetStats() map[string]interface{}
	GetType() string
	IsEnabled() bool
}

// NewManager creates a new processor manager
func NewManager(cfg *sharedconfig.Config) (*Manager, error) {
	mgr := &Manager{
		config:     cfg,
		processors: make(map[string]Processor),
	}

	// Initialize processors based on configuration
	if err := mgr.initializeProcessors(); err != nil {
		return nil, fmt.Errorf("failed to initialize processors: %w", err)
	}

	return mgr, nil
}

// initializeProcessors creates and configures all processors
func (m *Manager) initializeProcessors() error {
	// SMS Processor
	if m.config.SMS.Enabled {
		smsProcessor := sms.NewProcessor(&m.config.SMS)
		m.processors["sms"] = smsProcessor
		log.Println("SMS processor initialized")
	}

	// Voice Processor
	if m.config.Voice.Enabled {
		voiceProcessor := voice.NewProcessor(&m.config.Voice)
		m.processors["voice"] = voiceProcessor
		log.Println("Voice processor initialized")
	}

	// IRC Processor
	if m.config.IRC.Enabled {
		// TODO: Create IRC processor
		log.Println("IRC processor would be initialized here")
	}

	// ESMTP Processor
	if m.config.ESMTP.Enabled {
		// TODO: Create ESMTP processor
		log.Println("ESMTP processor would be initialized here")
	}

	// Discord Processor
	if m.config.Discord.Enabled {
		// TODO: Create Discord processor
		log.Println("Discord processor would be initialized here")
	}

	return nil
}

// Start starts all enabled processors
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	if m.isRunning {
		m.mu.Unlock()
		return fmt.Errorf("processor manager is already running")
	}
	m.isRunning = true
	m.mu.Unlock()

	log.Printf("Starting processor manager with %d processors", len(m.processors))

	// Start all processors
	for name, processor := range m.processors {
		if processor.IsEnabled() {
			log.Printf("Starting processor: %s", name)
			if err := processor.Start(ctx); err != nil {
				log.Printf("Failed to start processor %s: %v", name, err)
				// Continue starting other processors
			}
		}
	}

	log.Println("Processor manager started successfully")
	return nil
}

// Stop stops all processors
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return fmt.Errorf("processor manager is not running")
	}

	log.Println("Stopping processor manager...")

	// Stop all processors
	for name, processor := range m.processors {
		log.Printf("Stopping processor: %s", name)
		if err := processor.Stop(); err != nil {
			log.Printf("Error stopping processor %s: %v", name, err)
		}
	}

	m.isRunning = false
	log.Println("Processor manager stopped")
	return nil
}

// GetStatus returns status of all processors
func (m *Manager) GetStatus() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := make(map[string]interface{})
	status["is_running"] = m.isRunning
	status["processor_count"] = len(m.processors)

	processors := make(map[string]interface{})
	for name, processor := range m.processors {
		processors[name] = map[string]interface{}{
			"type":    processor.GetType(),
			"enabled": processor.IsEnabled(),
			"stats":   processor.GetStats(),
		}
	}
	status["processors"] = processors

	return status
}

// GetProcessor returns a specific processor by name
func (m *Manager) GetProcessor(name string) (Processor, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	processor, exists := m.processors[name]
	return processor, exists
}

// AddProcessor adds a processor to the manager
func (m *Manager) AddProcessor(name string, processor Processor) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.processors[name]; exists {
		return fmt.Errorf("processor %s already exists", name)
	}

	m.processors[name] = processor
	log.Printf("Added processor: %s (type: %s)", name, processor.GetType())
	return nil
}

// RemoveProcessor removes a processor from the manager
func (m *Manager) RemoveProcessor(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	processor, exists := m.processors[name]
	if !exists {
		return fmt.Errorf("processor %s not found", name)
	}

	// Stop processor if running
	if err := processor.Stop(); err != nil {
		log.Printf("Error stopping processor %s during removal: %v", name, err)
	}

	delete(m.processors, name)
	log.Printf("Removed processor: %s", name)
	return nil
}
