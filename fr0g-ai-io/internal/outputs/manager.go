package outputs

import (
	"context"
	"fmt"
	"log"
	"sync"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs/discord"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs/irc"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs/sms"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs/voice"
)

// Manager manages all output senders
type Manager struct {
	config  *sharedconfig.Config
	senders map[string]Sender
	mu      sync.RWMutex
	isRunning bool
}

// Sender defines the interface for output senders
type Sender interface {
	Start(ctx context.Context) error
	Stop() error
	Send(message *OutputMessage) error
	GetStats() map[string]interface{}
	GetType() string
	IsEnabled() bool
}

// OutputMessage represents a message to be sent
type OutputMessage struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`        // "sms", "voice", "irc", "discord"
	Destination string                 `json:"destination"` // phone number, channel, user ID, etc.
	Content     string                 `json:"content"`
	Metadata    map[string]interface{} `json:"metadata"`
	Priority    int                    `json:"priority"`
	RetryCount  int                    `json:"retry_count"`
	MaxRetries  int                    `json:"max_retries"`
}

// NewManager creates a new output manager
func NewManager(cfg *sharedconfig.Config) (*Manager, error) {
	mgr := &Manager{
		config:  cfg,
		senders: make(map[string]Sender),
	}

	if err := mgr.initializeSenders(); err != nil {
		return nil, fmt.Errorf("failed to initialize senders: %w", err)
	}

	return mgr, nil
}

// initializeSenders creates and configures all senders
func (m *Manager) initializeSenders() error {
	// SMS Sender
	if m.config.SMS.Enabled && m.config.SMS.ResponseEnabled {
		smsSender := sms.NewSender(&m.config.SMS)
		m.senders["sms"] = smsSender
		log.Println("SMS sender initialized")
	}

	// Voice Sender
	if m.config.Voice.Enabled && m.config.Voice.ResponseEnabled {
		voiceSender := voice.NewSender(&m.config.Voice)
		m.senders["voice"] = voiceSender
		log.Println("Voice sender initialized")
	}

	// IRC Sender
	if m.config.IRC.Enabled && m.config.IRC.ResponseEnabled {
		ircSender := irc.NewSender(&m.config.IRC)
		m.senders["irc"] = ircSender
		log.Println("IRC sender initialized")
	}

	// Discord Sender
	if m.config.Discord.Enabled && m.config.Discord.ResponseEnabled {
		discordSender := discord.NewSender(&m.config.Discord)
		m.senders["discord"] = discordSender
		log.Println("Discord sender initialized")
	}

	return nil
}

// Start starts all enabled senders
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	if m.isRunning {
		m.mu.Unlock()
		return fmt.Errorf("output manager is already running")
	}
	m.isRunning = true
	m.mu.Unlock()

	log.Printf("Starting output manager with %d senders", len(m.senders))

	for name, sender := range m.senders {
		if sender.IsEnabled() {
			log.Printf("Starting sender: %s", name)
			if err := sender.Start(ctx); err != nil {
				log.Printf("Failed to start sender %s: %v", name, err)
			}
		}
	}

	log.Println("Output manager started successfully")
	return nil
}

// Stop stops all senders
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return fmt.Errorf("output manager is not running")
	}

	log.Println("Stopping output manager...")

	for name, sender := range m.senders {
		log.Printf("Stopping sender: %s", name)
		if err := sender.Stop(); err != nil {
			log.Printf("Error stopping sender %s: %v", name, err)
		}
	}

	m.isRunning = false
	log.Println("Output manager stopped")
	return nil
}

// Send sends a message through the appropriate sender
func (m *Manager) Send(message *OutputMessage) error {
	m.mu.RLock()
	sender, exists := m.senders[message.Type]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("sender type %s not found", message.Type)
	}

	return sender.Send(message)
}

// GetStatus returns status of all senders
func (m *Manager) GetStatus() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := make(map[string]interface{})
	status["is_running"] = m.isRunning
	status["sender_count"] = len(m.senders)

	senders := make(map[string]interface{})
	for name, sender := range m.senders {
		senders[name] = map[string]interface{}{
			"type":    sender.GetType(),
			"enabled": sender.IsEnabled(),
			"stats":   sender.GetStats(),
		}
	}
	status["senders"] = senders

	return status
}
