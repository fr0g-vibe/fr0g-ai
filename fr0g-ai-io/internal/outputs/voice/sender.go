package voice

import (
	"context"
	"fmt"
	"log"
	"sync"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs/types"
)

// Sender handles Voice message sending
type Sender struct {
	config    *sharedconfig.VoiceConfig
	mu        sync.RWMutex
	isRunning bool
	sentCount int64
	failCount int64
}

// NewSender creates a new Voice sender
func NewSender(cfg *sharedconfig.VoiceConfig) *Sender {
	return &Sender{
		config: cfg,
	}
}

// GetType returns the sender type
func (s *Sender) GetType() string {
	return "voice"
}

// IsEnabled returns whether the sender is enabled
func (s *Sender) IsEnabled() bool {
	return s.config.Enabled && s.config.ResponseEnabled
}

// Start starts the Voice sender
func (s *Sender) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return fmt.Errorf("Voice sender is already running")
	}
	s.isRunning = true
	s.mu.Unlock()

	log.Println("Voice sender started")
	return nil
}

// Stop stops the Voice sender
func (s *Sender) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return fmt.Errorf("Voice sender is not running")
	}

	s.isRunning = false
	log.Println("Voice sender stopped")
	return nil
}

// Send sends a Voice message
func (s *Sender) Send(message *types.OutputCommand) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return fmt.Errorf("Voice sender is not running")
	}

	// TODO: Implement actual voice call/message sending
	// This would integrate with voice APIs like Twilio Voice API
	log.Printf("Sending Voice message to %s: %s", message.Target, message.Content)

	s.sentCount++
	return nil
}

// GetStats returns sender statistics
func (s *Sender) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"sent_count": s.sentCount,
		"fail_count": s.failCount,
		"is_running": s.isRunning,
	}
}
