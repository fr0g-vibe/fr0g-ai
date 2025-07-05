package queue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// Manager manages message queues for I/O operations
type Manager struct {
	config    *sharedconfig.QueueConfig
	inputQ    Queue
	outputQ   Queue
	mu        sync.RWMutex
	isRunning bool
}

// Queue defines the interface for message queues
type Queue interface {
	Enqueue(message *Message) error
	Dequeue() (*Message, error)
	Size() int
	IsEmpty() bool
	Clear() error
}

// Message represents a queued message
type Message struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`        // "input" or "output"
	Source      string                 `json:"source"`      // "sms", "voice", "irc", etc.
	Destination string                 `json:"destination"` // for output messages
	Content     interface{}            `json:"content"`
	Metadata    map[string]interface{} `json:"metadata"`
	Timestamp   time.Time              `json:"timestamp"`
	Retries     int                    `json:"retries"`
	MaxRetries  int                    `json:"max_retries"`
	Priority    int                    `json:"priority"` // Higher number = higher priority
}

// NewManager creates a new queue manager
func NewManager(cfg *sharedconfig.QueueConfig) (*Manager, error) {
	mgr := &Manager{
		config: cfg,
	}

	// Initialize queues based on configuration
	if err := mgr.initializeQueues(); err != nil {
		return nil, fmt.Errorf("failed to initialize queues: %w", err)
	}

	return mgr, nil
}

// initializeQueues creates input and output queues
func (m *Manager) initializeQueues() error {
	switch m.config.Type {
	case "memory":
		m.inputQ = NewMemoryQueue(m.config.MaxSize)
		m.outputQ = NewMemoryQueue(m.config.MaxSize)
	case "redis":
		// TODO: Implement Redis queue
		return fmt.Errorf("redis queue not implemented yet")
	case "rabbitmq":
		// TODO: Implement RabbitMQ queue
		return fmt.Errorf("rabbitmq queue not implemented yet")
	default:
		return fmt.Errorf("unsupported queue type: %s", m.config.Type)
	}

	return nil
}

// Start starts the queue manager
func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	if m.isRunning {
		m.mu.Unlock()
		return fmt.Errorf("queue manager is already running")
	}
	m.isRunning = true
	m.mu.Unlock()

	log.Printf("Starting queue manager with type: %s", m.config.Type)

	// Start queue processing goroutines
	go m.processInputQueue(ctx)
	go m.processOutputQueue(ctx)

	log.Println("Queue manager started successfully")
	return nil
}

// Stop stops the queue manager
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return fmt.Errorf("queue manager is not running")
	}

	log.Println("Stopping queue manager...")

	m.isRunning = false
	log.Println("Queue manager stopped")
	return nil
}

// EnqueueInput adds a message to the input queue
func (m *Manager) EnqueueInput(message *Message) error {
	message.Type = "input"
	message.Timestamp = time.Now()
	return m.inputQ.Enqueue(message)
}

// EnqueueOutput adds a message to the output queue
func (m *Manager) EnqueueOutput(message *Message) error {
	message.Type = "output"
	message.Timestamp = time.Now()
	return m.outputQ.Enqueue(message)
}

// DequeueInput removes a message from the input queue
func (m *Manager) DequeueInput() (*Message, error) {
	return m.inputQ.Dequeue()
}

// DequeueOutput removes a message from the output queue
func (m *Manager) DequeueOutput() (*Message, error) {
	return m.outputQ.Dequeue()
}

// GetStatus returns queue status
func (m *Manager) GetStatus() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"is_running":    m.isRunning,
		"type":          m.config.Type,
		"input_size":    m.inputQ.Size(),
		"output_size":   m.outputQ.Size(),
		"input_empty":   m.inputQ.IsEmpty(),
		"output_empty":  m.outputQ.IsEmpty(),
	}
}

// GetStats returns detailed queue statistics
func (m *Manager) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"input_queue": map[string]interface{}{
			"size":     m.inputQ.Size(),
			"is_empty": m.inputQ.IsEmpty(),
		},
		"output_queue": map[string]interface{}{
			"size":     m.outputQ.Size(),
			"is_empty": m.outputQ.IsEmpty(),
		},
		"config": map[string]interface{}{
			"type":           m.config.Type,
			"max_size":       m.config.MaxSize,
			"retry_attempts": m.config.RetryAttempts,
			"retry_delay":    m.config.RetryDelay.String(),
		},
	}
}

// processInputQueue processes messages from the input queue
func (m *Manager) processInputQueue(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !m.isRunning {
				return
			}

			// Process input messages
			for !m.inputQ.IsEmpty() {
				message, err := m.inputQ.Dequeue()
				if err != nil {
					log.Printf("Error dequeuing input message: %v", err)
					continue
				}

				// TODO: Send message to master-control for analysis
				log.Printf("Processing input message: %s from %s", message.ID, message.Source)
			}
		}
	}
}

// processOutputQueue processes messages from the output queue
func (m *Manager) processOutputQueue(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !m.isRunning {
				return
			}

			// Process output messages
			for !m.outputQ.IsEmpty() {
				message, err := m.outputQ.Dequeue()
				if err != nil {
					log.Printf("Error dequeuing output message: %v", err)
					continue
				}

				// TODO: Send message to appropriate output processor
				log.Printf("Processing output message: %s to %s", message.ID, message.Destination)
			}
		}
	}
}
