package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
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
	
	// Review and validation fields
	ReviewStatus   string     `json:"review_status,omitempty"`
	ReviewedBy     string     `json:"reviewed_by,omitempty"`
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty"`
	ReviewComments string     `json:"review_comments,omitempty"`
	RequiresReview bool       `json:"requires_review"`
}

// ReviewableOutputCommand represents a command that needs review
type ReviewableOutputCommand struct {
	Command        *OutputCommand `json:"command"`
	RequiresReview bool           `json:"requires_review"`
	ReviewStatus   string         `json:"review_status"`
	QueuedAt       time.Time      `json:"queued_at"`
	ReviewDeadline *time.Time     `json:"review_deadline,omitempty"`
}

// ValidationIssue represents a validation problem
type ValidationIssue struct {
	Field      string `json:"field"`
	Issue      string `json:"issue"`
	Severity   string `json:"severity"` // "error", "warning", "info"
	Suggestion string `json:"suggestion,omitempty"`
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

// ProcessOutputCommand processes output commands from master-control with enhanced review capabilities
func (s *IOService) ProcessOutputCommand(ctx context.Context, command *OutputCommand) (*OutputResponse, error) {
	log.Printf("gRPC Service: Received output command %s of type %s", command.ID, command.Type)

	// Validate the command first
	validationIssues := s.validateOutputCommand(command)
	if hasErrors(validationIssues) {
		return &OutputResponse{
			CommandID: command.ID,
			Success:   false,
			Message:   "Command validation failed",
			Metadata: map[string]interface{}{
				"validation_issues": validationIssues,
				"rejected_at":       time.Now(),
			},
			Timestamp: time.Now(),
		}, nil
	}

	// Check if command requires review
	requiresReview := s.shouldRequireReview(command)
	if requiresReview {
		log.Printf("gRPC Service: Command %s requires review, queuing for review", command.ID)
		
		// Add to review queue instead of direct processing
		reviewCommand := &ReviewableOutputCommand{
			Command:        command,
			RequiresReview: true,
			ReviewStatus:   "pending",
			QueuedAt:       time.Now(),
		}
		
		if err := s.queueForReview(reviewCommand); err != nil {
			return &OutputResponse{
				CommandID: command.ID,
				Success:   false,
				Message:   fmt.Sprintf("Failed to queue command for review: %v", err),
				Timestamp: time.Now(),
			}, nil
		}

		return &OutputResponse{
			CommandID: command.ID,
			Success:   true,
			Message:   "Command queued for review",
			Metadata: map[string]interface{}{
				"requires_review": true,
				"review_status":   "pending",
				"queued_at":       time.Now(),
				"validation_issues": validationIssues, // Include warnings
			},
			Timestamp: time.Now(),
		}, nil
	}

	// Process immediately if no review required
	return s.processOutputCommandDirect(ctx, command, validationIssues)
}

// processOutputCommandDirect processes a command directly without review
func (s *IOService) processOutputCommandDirect(ctx context.Context, command *OutputCommand, validationIssues []ValidationIssue) (*OutputResponse, error) {
	// Convert to queue message and send to output queue
	message := &queue.Message{
		ID:          command.ID,
		Type:        command.Type,
		Source:      "master-control",
		Destination: command.Target,
		Content:     command.Content,
		Metadata:    command.Metadata,
		Timestamp:   time.Now(),
		Retries:     0,
		MaxRetries:  3,
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
	
	metadata := map[string]interface{}{
		"queued_at":         time.Now(),
		"auto_approved":     true,
		"validation_issues": validationIssues,
	}

	return &OutputResponse{
		CommandID: command.ID,
		Success:   true,
		Message:   "Command queued successfully",
		Metadata:  metadata,
		Timestamp: time.Now(),
	}, nil
}

// validateOutputCommand validates an output command
func (s *IOService) validateOutputCommand(command *OutputCommand) []ValidationIssue {
	var issues []ValidationIssue

	// Basic validation
	if command.ID == "" {
		issues = append(issues, ValidationIssue{
			Field:    "id",
			Issue:    "Command ID is required",
			Severity: "error",
		})
	}

	if command.Type == "" {
		issues = append(issues, ValidationIssue{
			Field:    "type",
			Issue:    "Command type is required",
			Severity: "error",
		})
	}

	if command.Target == "" {
		issues = append(issues, ValidationIssue{
			Field:    "target",
			Issue:    "Target is required",
			Severity: "error",
		})
	}

	if command.Content == "" {
		issues = append(issues, ValidationIssue{
			Field:      "content",
			Issue:      "Content is empty",
			Severity:   "warning",
			Suggestion: "Consider adding meaningful content",
		})
	}

	// Type-specific validation
	switch command.Type {
	case "sms":
		if len(command.Content) > 160 {
			issues = append(issues, ValidationIssue{
				Field:      "content",
				Issue:      "SMS content exceeds 160 characters",
				Severity:   "warning",
				Suggestion: "Consider splitting into multiple messages",
			})
		}
	case "email":
		// Basic email validation could be added here
		if command.Metadata["subject"] == "" {
			issues = append(issues, ValidationIssue{
				Field:      "metadata.subject",
				Issue:      "Email subject is missing",
				Severity:   "warning",
				Suggestion: "Add a subject line for better deliverability",
			})
		}
	}

	return issues
}

// shouldRequireReview determines if a command needs manual review
func (s *IOService) shouldRequireReview(command *OutputCommand) bool {
	// High priority commands might need review
	if command.Priority > 8 {
		return true
	}

	// Commands with sensitive content patterns
	sensitivePatterns := []string{"urgent", "emergency", "critical", "alert"}
	contentLower := strings.ToLower(command.Content)
	for _, pattern := range sensitivePatterns {
		if strings.Contains(contentLower, pattern) {
			return true
		}
	}

	// Commands to external targets might need review
	if command.Metadata["external"] == "true" {
		return true
	}

	// Large content might need review
	if len(command.Content) > 1000 {
		return true
	}

	return false
}

// queueForReview queues a command for manual review
func (s *IOService) queueForReview(reviewCommand *ReviewableOutputCommand) error {
	// This would integrate with a review system
	// For now, just log it
	log.Printf("gRPC Service: Command %s queued for review", reviewCommand.Command.ID)
	
	// In a real implementation, this would:
	// 1. Store the command in a review database
	// 2. Notify reviewers
	// 3. Set up review workflow
	
	return nil
}

// hasErrors checks if validation issues contain any errors
func hasErrors(issues []ValidationIssue) bool {
	for _, issue := range issues {
		if issue.Severity == "error" {
			return true
		}
	}
	return false
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
	content, ok := message.Content.(string)
	if !ok {
		return fmt.Errorf("message content is not a string")
	}

	event := &InputEvent{
		ID:        message.ID,
		Type:      message.Type,
		Source:    message.Source,
		Content:   content,
		Metadata:  message.Metadata,
		Timestamp: message.Timestamp,
		Priority:  0, // Default priority since queue.Message doesn't have Priority field
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
			ID:          outputCommand.ID,
			Type:        outputCommand.Type,
			Source:      "master-control-response",
			Destination: outputCommand.Target,
			Content:     outputCommand.Content,
			Metadata:    outputCommand.Metadata,
			Timestamp:   time.Now(),
			Retries:     0,
			MaxRetries:  3,
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
