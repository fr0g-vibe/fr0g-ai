package irc

import (
	"context"
	"fmt"
	"log"
	"sync"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs"
)

// Sender handles IRC message sending
type Sender struct {
	config    *sharedconfig.IRCConfig
	mu        sync.RWMutex
	isRunning bool
	sentCount int64
	failCount int64
}

// NewSender creates a new IRC sender
func NewSender(cfg *sharedconfig.IRCConfig) *Sender {
	return &Sender{
		config: cfg,
	}
}

// GetType returns the sender type
func (s *Sender) GetType() string {
	return "irc"
}

// IsEnabled returns whether the sender is enabled
func (s *Sender) IsEnabled() bool {
	return s.config.Enabled && s.config.ResponseEnabled
}

// Start starts the IRC sender
func (s *Sender) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return fmt.Errorf("IRC sender is already running")
	}
	s.isRunning = true
	s.mu.Unlock()

	log.Println("IRC sender started")
	return nil
}

// Stop stops the IRC sender
func (s *Sender) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return fmt.Errorf("IRC sender is not running")
	}

	s.isRunning = false
	log.Println("IRC sender stopped")
	return nil
}

// Send sends an IRC message
func (s *Sender) Send(message *outputs.OutputMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return fmt.Errorf("IRC sender is not running")
	}

	// TODO: Implement actual IRC message sending
	// This would integrate with IRC libraries like go-ircevent
	log.Printf("Sending IRC message to %s: %s", message.Destination, message.Content)

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
