package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/queue"
	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// IOService implements the gRPC service for fr0g-ai-io
type IOService struct {
	inputQueue  queue.Queue
	outputQueue queue.Queue
	config      *sharedconfig.GRPCConfig
	server      *grpc.Server
	listener    net.Listener
	mu          sync.RWMutex
	isRunning   bool
	
	// Master Control client for sending events
	mcpClient MCPClient
}

// MCPClient defines interface for communicating with master-control
type MCPClient interface {
	SendInputEvent(ctx context.Context, event *InputEvent) (*InputEventResponse, error)
	IsConnected() bool
	Close() error
}

// InputEvent represents an input event to send to master-control
type InputEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"` // "sms", "voice", "irc", "discord"
	Source    string                 `json:"source"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
	Priority  int                    `json:"priority"`
}

// InputEventResponse represents response from master-control
type InputEventResponse struct {
	EventID     string                 `json:"event_id"`
	Processed   bool                   `json:"processed"`
	Actions     []OutputAction         `json:"actions"`
	Metadata    map[string]interface{} `json:"metadata"`
	ProcessedAt time.Time              `json:"processed_at"`
}

// OutputAction represents an action to be taken
type OutputAction struct {
	Type     string                 `json:"type"`
	Target   string                 `json:"target"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
}

// OutputCommand represents a command from master-control
type OutputCommand struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Target   string                 `json:"target"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
	Priority int                    `json:"priority"`
}

// OutputResponse represents response to output command
type OutputResponse struct {
	CommandID string                 `json:"command_id"`
	Success   bool                   `json:"success"`
	Message   string                 `json:"message"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
}

// ThreatAnalysisResult represents threat analysis from master-control
type ThreatAnalysisResult struct {
	EventID     string    `json:"event_id"`
	ThreatLevel string    `json:"threat_level"`
	ThreatScore float64   `json:"threat_score"`
	Analysis    string    `json:"analysis"`
	AnalyzedAt  time.Time `json:"analyzed_at"`
}

// NewIOService creates a new gRPC I/O service
func NewIOService(inputQueue, outputQueue queue.Queue, config *sharedconfig.GRPCConfig) *IOService {
	return &IOService{
		inputQueue:  inputQueue,
		outputQueue: outputQueue,
		config:      config,
	}
}

// SetMCPClient sets the master-control client
func (s *IOService) SetMCPClient(client MCPClient) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mcpClient = client
	log.Printf("gRPC Service: Master-control client configured")
}

// Start starts the gRPC service
func (s *IOService) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("gRPC service is already running")
	}

	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", address, err)
	}

	s.listener = listener
	s.server = grpc.NewServer()

	// Register the service (this would use actual protobuf definitions)
	// For now, we'll set up the server structure
	
	// Enable reflection for debugging
	reflection.Register(s.server)

	s.isRunning = true

	go func() {
		log.Printf("gRPC Service: Starting server on %s", address)
		if err := s.server.Serve(s.listener); err != nil {
			log.Printf("gRPC Service: Server error: %v", err)
		}
	}()

	log.Printf("gRPC Service: Started successfully on %s", address)
	return nil
}

// Stop stops the gRPC service
func (s *IOService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	log.Printf("gRPC Service: Stopping...")

	if s.server != nil {
		s.server.GracefulStop()
	}

	if s.listener != nil {
		s.listener.Close()
	}

	if s.mcpClient != nil {
		s.mcpClient.Close()
	}

	s.isRunning = false
	log.Printf("gRPC Service: Stopped")
	return nil
}

// ProcessOutputCommand processes output commands from master-control
func (s *IOService) ProcessOutputCommand(ctx context.Context, command *OutputCommand) (*OutputResponse, error) {
	log.Printf("gRPC Service: Received output command %s of type %s", command.ID, command.Type)

	// Convert to queue message and send to output queue
	message := &queue.Message{
		ID:        command.ID,
		Type:      command.Type,
		Source:    "master-control",
		Target:    command.Target,
		Content:   command.Content,
		Metadata:  command.Metadata,
		Priority:  command.Priority,
		Timestamp: time.Now(),
	}

	if err := s.outputQueue.Enqueue(message); err != nil {
		return &OutputResponse{
			CommandID: command.ID,
			Success:   false,
			Message:   fmt.Sprintf("Failed to queue command: %v", err),
			Timestamp: time.Now(),
		}, nil
	}

	log.Printf("gRPC Service: Successfully queued output command %s", command.ID)
	return &OutputResponse{
		CommandID: command.ID,
		Success:   true,
		Message:   "Command queued successfully",
		Metadata: map[string]interface{}{
			"queued_at": time.Now(),
		},
		Timestamp: time.Now(),
	}, nil
}

// ProcessThreatAnalysis processes threat analysis results from master-control
func (s *IOService) ProcessThreatAnalysis(ctx context.Context, result *ThreatAnalysisResult) error {
	log.Printf("gRPC Service: Received threat analysis for event %s with level %s", 
		result.EventID, result.ThreatLevel)

	// Handle threat analysis (could trigger alerts, logging, etc.)
	// For now, just log it
	log.Printf("gRPC Service: Threat Analysis - Event: %s, Level: %s, Score: %.2f", 
		result.EventID, result.ThreatLevel, result.ThreatScore)

	return nil
}

// SendInputEvent sends an input event to master-control
func (s *IOService) SendInputEvent(ctx context.Context, event *InputEvent) (*InputEventResponse, error) {
	s.mu.RLock()
	client := s.mcpClient
	s.mu.RUnlock()

	if client == nil {
		return nil, fmt.Errorf("master-control client not configured")
	}

	if !client.IsConnected() {
		return nil, fmt.Errorf("master-control client not connected")
	}

	log.Printf("gRPC Service: Sending input event %s to master-control", event.ID)

	response, err := client.SendInputEvent(ctx, event)
	if err != nil {
		log.Printf("gRPC Service: Error sending input event %s: %v", event.ID, err)
		return nil, err
	}

	log.Printf("gRPC Service: Successfully sent input event %s to master-control", event.ID)
	return response, nil
}

// ProcessInputMessage processes an input message and sends it to master-control
func (s *IOService) ProcessInputMessage(message *queue.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Convert queue message to input event
	event := &InputEvent{
		ID:        message.ID,
		Type:      message.Type,
		Source:    message.Source,
		Content:   message.Content,
		Metadata:  message.Metadata,
		Timestamp: message.Timestamp,
		Priority:  message.Priority,
	}

	response, err := s.SendInputEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to send input event to master-control: %w", err)
	}

	// Process any actions returned from master-control
	for _, action := range response.Actions {
		outputCommand := &OutputCommand{
			ID:       fmt.Sprintf("action_%s_%d", event.ID, time.Now().UnixNano()),
			Type:     action.Type,
			Target:   action.Target,
			Content:  action.Content,
			Metadata: action.Metadata,
			Priority: event.Priority,
		}

		// Send to output queue for processing
		outputMessage := &queue.Message{
			ID:        outputCommand.ID,
			Type:      outputCommand.Type,
			Source:    "master-control-response",
			Target:    outputCommand.Target,
			Content:   outputCommand.Content,
			Metadata:  outputCommand.Metadata,
			Priority:  outputCommand.Priority,
			Timestamp: time.Now(),
		}

		if err := s.outputQueue.Enqueue(outputMessage); err != nil {
			log.Printf("gRPC Service: Error queuing output action: %v", err)
		} else {
			log.Printf("gRPC Service: Queued output action %s", outputCommand.ID)
		}
	}

	return nil
}

// GetStatus returns the service status
func (s *IOService) GetStatus() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status := map[string]interface{}{
		"running":     s.isRunning,
		"address":     fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		"mcp_client":  s.mcpClient != nil,
	}

	if s.mcpClient != nil {
		status["mcp_connected"] = s.mcpClient.IsConnected()
	}

	return status
}
