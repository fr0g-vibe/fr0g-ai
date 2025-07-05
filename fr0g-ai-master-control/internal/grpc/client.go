package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"fr0g-ai-master-control/internal/mastercontrol/input"
)

// Fr0gIOGRPCClient implements the Fr0gIOClient interface
type Fr0gIOGRPCClient struct {
	conn   *grpc.ClientConn
	config *ClientConfig
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

	return &Fr0gIOGRPCClient{
		conn:   conn,
		config: config,
	}, nil
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

// Close closes the gRPC connection
func (c *Fr0gIOGRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

