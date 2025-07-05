package io

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"fr0g-ai-master-control/internal/grpc"
	"fr0g-ai-master-control/internal/mastercontrol/input"
)

// Manager handles bidirectional I/O communication with fr0g-ai-io service
type Manager struct {
	client          *grpc.Fr0gIOGRPCClient
	eventChan       chan *input.InputEvent
	outputChan      chan *input.OutputCommand
	threatChan      chan *input.ThreatAnalysisResult
	config          *Config
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	mu              sync.RWMutex
	isRunning       bool
	eventProcessors map[string]EventProcessor
}

// Config holds I/O manager configuration
type Config struct {
	InputEventBufferSize  int
	OutputCommandTimeout  time.Duration
	ThreatAnalysisEnabled bool
	MaxConcurrentEvents   int
}

// EventProcessor defines the interface for processing different types of input events
type EventProcessor interface {
	ProcessEvent(ctx context.Context, event *input.InputEvent) (*input.InputEventResponse, error)
	GetEventType() string
}

// NewManager creates a new I/O manager
func NewManager(client *grpc.Fr0gIOGRPCClient, config *Config) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		client:          client,
		eventChan:       make(chan *input.InputEvent, config.InputEventBufferSize),
		outputChan:      make(chan *input.OutputCommand, config.InputEventBufferSize),
		threatChan:      make(chan *input.ThreatAnalysisResult, config.InputEventBufferSize),
		config:          config,
		ctx:             ctx,
		cancel:          cancel,
		eventProcessors: make(map[string]EventProcessor),
	}
}

// Start begins I/O manager operation
func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isRunning {
		return fmt.Errorf("I/O manager is already running")
	}

	log.Printf("I/O Manager: Starting bidirectional communication with fr0g-ai-io")

	// Set this manager as the input handler for the gRPC client
	m.client.SetInputHandler(m)

	// Start input event listener
	if err := m.client.StartInputEventListener(m.ctx); err != nil {
		return fmt.Errorf("failed to start input event listener: %w", err)
	}

	// Start worker goroutines
	m.wg.Add(3)
	go m.processInputEvents()
	go m.processOutputCommands()
	go m.processThreatAnalysis()

	m.isRunning = true
	log.Printf("I/O Manager: Started successfully")
	return nil
}

// Stop gracefully stops the I/O manager
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isRunning {
		return nil
	}

	log.Printf("I/O Manager: Stopping...")

	m.cancel()
	m.wg.Wait()

	close(m.eventChan)
	close(m.outputChan)
	close(m.threatChan)

	m.isRunning = false
	log.Printf("I/O Manager: Stopped")
	return nil
}

// RegisterEventProcessor registers a processor for a specific event type
func (m *Manager) RegisterEventProcessor(processor EventProcessor) {
	m.mu.Lock()
	defer m.mu.Unlock()

	eventType := processor.GetEventType()
	m.eventProcessors[eventType] = processor
	log.Printf("I/O Manager: Registered event processor for type '%s'", eventType)
}

// HandleInputEvent implements Fr0gIOInputHandler interface
func (m *Manager) HandleInputEvent(ctx context.Context, event *input.InputEvent) (*input.InputEventResponse, error) {
	log.Printf("I/O Manager: Received input event %s of type %s", event.ID, event.Type)

	// Send event to processing channel
	select {
	case m.eventChan <- event:
		// Event queued successfully
		return &input.InputEventResponse{
			EventID:     event.ID,
			Processed:   true,
			ProcessedAt: time.Now(),
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return nil, fmt.Errorf("event queue is full")
	}
}

// HandleSMSMessage implements Fr0gIOInputHandler interface
func (m *Manager) HandleSMSMessage(ctx context.Context, message *input.SMSMessage) error {
	event := &input.InputEvent{
		ID:        message.ID,
		Type:      "sms",
		Source:    message.From,
		Content:   message.Content,
		Timestamp: message.Timestamp,
		Metadata: map[string]interface{}{
			"to":           message.To,
			"provider":     message.Provider,
			"message_type": message.MessageType,
		},
	}

	_, err := m.HandleInputEvent(ctx, event)
	return err
}

// HandleVoiceMessage implements Fr0gIOInputHandler interface
func (m *Manager) HandleVoiceMessage(ctx context.Context, message *input.VoiceMessage) error {
	event := &input.InputEvent{
		ID:        message.ID,
		Type:      "voice",
		Source:    message.From,
		Content:   message.Transcription,
		Timestamp: message.Timestamp,
		Metadata: map[string]interface{}{
			"to":       message.To,
			"duration": message.Duration,
			"format":   message.Format,
		},
	}

	_, err := m.HandleInputEvent(ctx, event)
	return err
}

// HandleIRCMessage implements Fr0gIOInputHandler interface
func (m *Manager) HandleIRCMessage(ctx context.Context, message *input.IRCMessage) error {
	event := &input.InputEvent{
		ID:        message.ID,
		Type:      "irc",
		Source:    message.Nick,
		Content:   message.Content,
		Timestamp: message.Timestamp,
		Metadata: map[string]interface{}{
			"server":     message.Server,
			"channel":    message.Channel,
			"is_private": message.IsPrivate,
		},
	}

	_, err := m.HandleInputEvent(ctx, event)
	return err
}

// HandleDiscordMessage implements Fr0gIOInputHandler interface
func (m *Manager) HandleDiscordMessage(ctx context.Context, message *input.DiscordMessage) error {
	event := &input.InputEvent{
		ID:        message.ID,
		Type:      "discord",
		Source:    message.Username,
		Content:   message.Content,
		Timestamp: message.Timestamp,
		Metadata: map[string]interface{}{
			"guild_id":     message.GuildID,
			"channel_id":   message.ChannelID,
			"user_id":      message.UserID,
			"message_type": message.MessageType,
		},
	}

	_, err := m.HandleInputEvent(ctx, event)
	return err
}

// SendOutputCommand queues an output command to be sent to fr0g-ai-io
func (m *Manager) SendOutputCommand(command *input.OutputCommand) error {
	select {
	case m.outputChan <- command:
		log.Printf("I/O Manager: Queued output command %s", command.ID)
		return nil
	default:
		return fmt.Errorf("output command queue is full")
	}
}

// SendThreatAnalysis queues threat analysis results to be sent to fr0g-ai-io
func (m *Manager) SendThreatAnalysis(result *input.ThreatAnalysisResult) error {
	if !m.config.ThreatAnalysisEnabled {
		return nil
	}

	select {
	case m.threatChan <- result:
		log.Printf("I/O Manager: Queued threat analysis for event %s", result.EventID)
		return nil
	default:
		return fmt.Errorf("threat analysis queue is full")
	}
}

// Worker goroutines

func (m *Manager) processInputEvents() {
	defer m.wg.Done()
	log.Printf("I/O Manager: Started input event processor")

	for {
		select {
		case <-m.ctx.Done():
			return
		case event := <-m.eventChan:
			if event == nil {
				return
			}

			m.processEvent(event)
		}
	}
}

func (m *Manager) processOutputCommands() {
	defer m.wg.Done()
	log.Printf("I/O Manager: Started output command processor")

	for {
		select {
		case <-m.ctx.Done():
			return
		case command := <-m.outputChan:
			if command == nil {
				return
			}

			m.sendOutputCommand(command)
		}
	}
}

func (m *Manager) processThreatAnalysis() {
	defer m.wg.Done()
	log.Printf("I/O Manager: Started threat analysis processor")

	for {
		select {
		case <-m.ctx.Done():
			return
		case result := <-m.threatChan:
			if result == nil {
				return
			}

			m.sendThreatAnalysis(result)
		}
	}
}

// Helper methods

func (m *Manager) processEvent(event *input.InputEvent) {
	ctx, cancel := context.WithTimeout(m.ctx, 30*time.Second)
	defer cancel()

	m.mu.RLock()
	processor, exists := m.eventProcessors[event.Type]
	m.mu.RUnlock()

	if !exists {
		log.Printf("I/O Manager: No processor found for event type '%s', using default processing", event.Type)
		m.defaultEventProcessing(ctx, event)
		return
	}

	response, err := processor.ProcessEvent(ctx, event)
	if err != nil {
		log.Printf("I/O Manager: Error processing event %s: %v", event.ID, err)
		return
	}

	// Send any output actions from the response
	for _, action := range response.Actions {
		command := &input.OutputCommand{
			ID:       fmt.Sprintf("cmd_%s_%d", event.ID, time.Now().UnixNano()),
			Type:     action.Type,
			Target:   action.Target,
			Content:  action.Content,
			Metadata: action.Metadata,
			Priority: event.Priority,
		}

		if err := m.SendOutputCommand(command); err != nil {
			log.Printf("I/O Manager: Error queuing output command: %v", err)
		}
	}

	// Send threat analysis if available
	if response.Analysis != nil {
		if err := m.SendThreatAnalysis(response.Analysis); err != nil {
			log.Printf("I/O Manager: Error queuing threat analysis: %v", err)
		}
	}
}

func (m *Manager) defaultEventProcessing(ctx context.Context, event *input.InputEvent) {
	log.Printf("I/O Manager: Default processing for event %s: %s", event.ID, event.Content)
	
	// Simple echo response for demonstration
	command := &input.OutputCommand{
		ID:      fmt.Sprintf("echo_%s", event.ID),
		Type:    event.Type,
		Target:  event.Source,
		Content: fmt.Sprintf("Processed: %s", event.Content),
		Metadata: map[string]interface{}{
			"original_event_id": event.ID,
			"processing_type":   "default",
		},
		Priority: event.Priority,
	}

	if err := m.SendOutputCommand(command); err != nil {
		log.Printf("I/O Manager: Error queuing default response: %v", err)
	}
}

func (m *Manager) sendOutputCommand(command *input.OutputCommand) {
	ctx, cancel := context.WithTimeout(m.ctx, m.config.OutputCommandTimeout)
	defer cancel()

	response, err := m.client.SendOutputCommand(ctx, command)
	if err != nil {
		log.Printf("I/O Manager: Error sending output command %s: %v", command.ID, err)
		return
	}

	log.Printf("I/O Manager: Successfully sent output command %s: %s", command.ID, response.Message)
}

func (m *Manager) sendThreatAnalysis(result *input.ThreatAnalysisResult) {
	ctx, cancel := context.WithTimeout(m.ctx, m.config.OutputCommandTimeout)
	defer cancel()

	err := m.client.SendThreatAnalysisResult(ctx, result)
	if err != nil {
		log.Printf("I/O Manager: Error sending threat analysis for event %s: %v", result.EventID, err)
		return
	}

	log.Printf("I/O Manager: Successfully sent threat analysis for event %s (threat level: %s)", 
		result.EventID, result.ThreatLevel)
}

// GetStatus returns the current status of the I/O manager
func (m *Manager) GetStatus() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return map[string]interface{}{
		"running":           m.isRunning,
		"client_connected":  m.client.IsConnected(),
		"event_queue_size":  len(m.eventChan),
		"output_queue_size": len(m.outputChan),
		"threat_queue_size": len(m.threatChan),
		"processors":        len(m.eventProcessors),
	}
}
