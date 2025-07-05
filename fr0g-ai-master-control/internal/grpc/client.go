package grpc

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"fr0g-ai-master-control/internal/mastercontrol/input"
)

// Fr0gIOGRPCClient implements the Fr0gIOClient interface
type Fr0gIOGRPCClient struct {
	conn         *grpc.ClientConn
	config       *ClientConfig
	eventHandler input.Fr0gIOInputHandler
	mu           sync.RWMutex
	isConnected  bool
}

// ClientConfig holds gRPC client configuration
type ClientConfig struct {
	Host        string
	Port        int
	Timeout     time.Duration
	MaxRetries  int
	ServiceName string
}

// NewFr0gIOGRPCClient creates a new gRPC client for fr0g-ai-io
func NewFr0gIOGRPCClient(config *ClientConfig) (*Fr0gIOGRPCClient, error) {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to fr0g-ai-io service at %s: %w", address, err)
	}

	client := &Fr0gIOGRPCClient{
		conn:        conn,
		config:      config,
		isConnected: true,
	}

	log.Printf("gRPC Client: Connected to fr0g-ai-io service at %s", address)
	return client, nil
}

// SetInputHandler sets the input event handler for processing events from fr0g-ai-io
func (c *Fr0gIOGRPCClient) SetInputHandler(handler input.Fr0gIOInputHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.eventHandler = handler
	log.Printf("gRPC Client: Input handler registered for fr0g-ai-io events")
}

// StartInputEventListener starts listening for input events from fr0g-ai-io
func (c *Fr0gIOGRPCClient) StartInputEventListener(ctx context.Context) error {
	c.mu.RLock()
	handler := c.eventHandler
	c.mu.RUnlock()

	if handler == nil {
		return fmt.Errorf("no input handler registered")
	}

	log.Printf("gRPC Client: Starting input event listener for fr0g-ai-io")
	
	// This would typically be a streaming gRPC call to receive events
	// For now, we'll simulate with a ticker for demonstration
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Printf("gRPC Client: Input event listener stopped")
				return
			case <-ticker.C:
				// In a real implementation, this would receive actual events from fr0g-ai-io
				// For now, we'll just log that we're listening
				log.Printf("gRPC Client: Listening for input events from fr0g-ai-io...")
			}
		}
	}()

	return nil
}

// ProcessInputEvent processes an input event received from fr0g-ai-io
func (c *Fr0gIOGRPCClient) ProcessInputEvent(ctx context.Context, event *input.InputEvent) (*input.InputEventResponse, error) {
	c.mu.RLock()
	handler := c.eventHandler
	c.mu.RUnlock()

	if handler == nil {
		return nil, fmt.Errorf("no input handler registered")
	}

	log.Printf("gRPC Client: Processing input event %s of type %s from fr0g-ai-io", event.ID, event.Type)
	
	response, err := handler.HandleInputEvent(ctx, event)
	if err != nil {
		log.Printf("gRPC Client: Error processing input event %s: %v", event.ID, err)
		return nil, err
	}

	log.Printf("gRPC Client: Successfully processed input event %s", event.ID)
	return response, nil
}

// SendOutputCommand sends an output command to fr0g-ai-io
func (c *Fr0gIOGRPCClient) SendOutputCommand(ctx context.Context, command *input.OutputCommand) (*input.OutputResponse, error) {
	if c.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.config.Timeout)
		defer cancel()
	}

	// For now, simulate the gRPC call with logging
	log.Printf("gRPC Client: Sending output command %s of type %s to %s", command.ID, command.Type, command.Target)

	// Simulate successful response
	return &input.OutputResponse{
		CommandID: command.ID,
		Success:   true,
		Message:   "Command processed successfully (simulated)",
		Metadata:  map[string]interface{}{"grpc_response": true, "simulated": true},
		Timestamp: time.Now(),
	}, nil
}

// SendThreatAnalysisResult sends threat analysis results to fr0g-ai-io
func (c *Fr0gIOGRPCClient) SendThreatAnalysisResult(ctx context.Context, result *input.ThreatAnalysisResult) error {
	if c.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.config.Timeout)
		defer cancel()
	}

	// For now, simulate the gRPC call with logging
	log.Printf("gRPC Client: Sending threat analysis for event %s with threat level %s (simulated)", 
		result.EventID, result.ThreatLevel)

	return nil
}

// GetServiceStatus gets the status of fr0g-ai-io service
func (c *Fr0gIOGRPCClient) GetServiceStatus(ctx context.Context) (*input.ServiceStatus, error) {
	if c.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.config.Timeout)
		defer cancel()
	}

	// For now, simulate the gRPC call with logging
	log.Printf("gRPC Client: Getting service status (simulated)")

	return &input.ServiceStatus{
		ServiceName:   c.config.ServiceName,
		Status:        "healthy",
		Version:       "1.0.0",
		Uptime:        "24h30m",
		LastHeartbeat: time.Now(),
	}, nil
}

// IsConnected returns the connection status
func (c *Fr0gIOGRPCClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isConnected
}

// Close closes the gRPC connection
func (c *Fr0gIOGRPCClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.isConnected = false
	if c.conn != nil {
		log.Printf("gRPC Client: Closing connection to fr0g-ai-io service")
		return c.conn.Close()
	}
	return nil
}

