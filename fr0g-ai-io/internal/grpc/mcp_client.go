package grpc

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InputEvent represents an input event to be processed by master-control
type InputEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
	Priority  int                    `json:"priority"`
}

// InputEventResponse represents the response from master-control
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

// MCPGRPCClient implements the MCPClient interface for communicating with master-control
type MCPGRPCClient struct {
	conn        *grpc.ClientConn
	config      *MCPClientConfig
	mu          sync.RWMutex
	isConnected bool
}

// MCPClientConfig holds master-control client configuration
type MCPClientConfig struct {
	Host        string        `yaml:"host"`
	Port        int           `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	MaxRetries  int           `yaml:"max_retries"`
	ServiceName string        `yaml:"service_name"`
}

// NewMCPGRPCClient creates a new master-control gRPC client
func NewMCPGRPCClient(config *MCPClientConfig) (*MCPGRPCClient, error) {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master-control at %s: %w", address, err)
	}

	client := &MCPGRPCClient{
		conn:        conn,
		config:      config,
		isConnected: true,
	}

	log.Printf("MCP Client: Connected to master-control at %s", address)
	return client, nil
}

// SendInputEvent sends an input event to master-control
func (c *MCPGRPCClient) SendInputEvent(ctx context.Context, event *InputEvent) (*InputEventResponse, error) {
	c.mu.RLock()
	connected := c.isConnected
	c.mu.RUnlock()

	if !connected {
		return nil, fmt.Errorf("not connected to master-control")
	}

	if c.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.config.Timeout)
		defer cancel()
	}

	// For now, simulate the gRPC call with logging
	// In a real implementation, this would use the actual protobuf service
	log.Printf("MCP Client: Sending input event %s of type %s to master-control", event.ID, event.Type)

	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	// Generate a simulated response based on event type
	response := &InputEventResponse{
		EventID:     event.ID,
		Processed:   true,
		ProcessedAt: time.Now(),
		Actions:     []OutputAction{},
		Metadata: map[string]interface{}{
			"processed_by": "master-control",
			"simulated":    true,
		},
	}

	// Generate appropriate response actions based on event type
	switch event.Type {
	case "sms":
		response.Actions = append(response.Actions, OutputAction{
			Type:    "sms",
			Target:  event.Source,
			Content: fmt.Sprintf("Auto-reply: Received your SMS message: %s", event.Content),
			Metadata: map[string]interface{}{
				"response_type": "auto_reply",
				"original_id":   event.ID,
			},
		})

	case "voice":
		response.Actions = append(response.Actions, OutputAction{
			Type:    "voice",
			Target:  event.Source,
			Content: fmt.Sprintf("Voice message processed: %s", event.Content),
			Metadata: map[string]interface{}{
				"response_type": "voice_response",
				"original_id":   event.ID,
			},
		})

	case "irc":
		channel := ""
		if event.Metadata != nil {
			if ch, ok := event.Metadata["channel"].(string); ok {
				channel = ch
			}
		}
		
		target := event.Source
		if channel != "" {
			target = channel
		}

		response.Actions = append(response.Actions, OutputAction{
			Type:    "irc",
			Target:  target,
			Content: fmt.Sprintf("IRC: %s", event.Content),
			Metadata: map[string]interface{}{
				"response_type": "irc_response",
				"original_id":   event.ID,
			},
		})

	case "discord":
		channelID := ""
		if event.Metadata != nil {
			if ch, ok := event.Metadata["channel_id"].(string); ok {
				channelID = ch
			}
		}

		response.Actions = append(response.Actions, OutputAction{
			Type:    "discord",
			Target:  channelID,
			Content: fmt.Sprintf("Discord: Processed message from %s", event.Source),
			Metadata: map[string]interface{}{
				"response_type": "discord_response",
				"original_id":   event.ID,
			},
		})
	}

	log.Printf("MCP Client: Successfully processed input event %s with %d actions", 
		event.ID, len(response.Actions))
	return response, nil
}

// IsConnected returns the connection status
func (c *MCPGRPCClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isConnected
}

// Close closes the connection to master-control
func (c *MCPGRPCClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.isConnected = false
	if c.conn != nil {
		log.Printf("MCP Client: Closing connection to master-control")
		return c.conn.Close()
	}
	return nil
}

// Reconnect attempts to reconnect to master-control
func (c *MCPGRPCClient) Reconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		c.conn.Close()
	}

	address := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to reconnect to master-control at %s: %w", address, err)
	}

	c.conn = conn
	c.isConnected = true

	log.Printf("MCP Client: Reconnected to master-control at %s", address)
	return nil
}

// GetStatus returns client status information
func (c *MCPGRPCClient) GetStatus() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]interface{}{
		"connected":    c.isConnected,
		"service_name": c.config.ServiceName,
		"address":      fmt.Sprintf("%s:%d", c.config.Host, c.config.Port),
		"timeout":      c.config.Timeout.String(),
	}
}
