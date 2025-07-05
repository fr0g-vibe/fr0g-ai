package sms

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs"
)

// Sender handles SMS message sending
type Sender struct {
	config    *sharedconfig.SMSConfig
	mu        sync.RWMutex
	isRunning bool
	sentCount int64
	failCount int64
}

// NewSender creates a new SMS sender
func NewSender(cfg *sharedconfig.SMSConfig) *Sender {
	return &Sender{
		config: cfg,
	}
}

// GetType returns the sender type
func (s *Sender) GetType() string {
	return "sms"
}

// IsEnabled returns whether the sender is enabled
func (s *Sender) IsEnabled() bool {
	return s.config.Enabled && s.config.ResponseEnabled
}

// Start starts the SMS sender
func (s *Sender) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return fmt.Errorf("SMS sender is already running")
	}
	s.isRunning = true
	s.mu.Unlock()

	log.Println("SMS sender started")
	return nil
}

// Stop stops the SMS sender
func (s *Sender) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return fmt.Errorf("SMS sender is not running")
	}

	s.isRunning = false
	log.Println("SMS sender stopped")
	return nil
}

// Send sends an SMS message
func (s *Sender) Send(message *outputs.OutputMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return fmt.Errorf("SMS sender is not running")
	}

	// TODO: Implement actual SMS sending logic
	// This would integrate with SMS providers like Twilio, AWS SNS, etc.
	log.Printf("Sending SMS to %s: %s", message.Destination, message.Content)

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
