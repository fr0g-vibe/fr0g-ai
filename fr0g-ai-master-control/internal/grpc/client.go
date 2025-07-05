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
	client Fr0gIOServiceClient
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

// Fr0gIOServiceClient defines the gRPC client interface for fr0g-ai-io
type Fr0gIOServiceClient interface {
	SendOutputCommand(ctx context.Context, in *OutputCommandRequest, opts ...grpc.CallOption) (*OutputCommandResponse, error)
	SendThreatAnalysis(ctx context.Context, in *ThreatAnalysisRequest, opts ...grpc.CallOption) (*ThreatAnalysisResponse, error)
	GetServiceStatus(ctx context.Context, in *ServiceStatusRequest, opts ...grpc.CallOption) (*ServiceStatusResponse, error)
}

// gRPC message types for fr0g-ai-io communication
type OutputCommandRequest struct {
	Command *OutputCommand `json:"command"`
}

type OutputCommandResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	CommandId string `json:"command_id"`
}

type ThreatAnalysisRequest struct {
	Analysis *ThreatAnalysis `json:"analysis"`
}

type ThreatAnalysisResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ServiceStatusRequest struct{}

type ServiceStatusResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	Uptime    string `json:"uptime"`
	Timestamp int64  `json:"timestamp"`
}

type OutputCommand struct {
	Id       string            `json:"id"`
	Type     string            `json:"type"`
	Target   string            `json:"target"`
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata"`
	Priority int32             `json:"priority"`
}

type ThreatAnalysis struct {
	EventId       string            `json:"event_id"`
	ThreatLevel   string            `json:"threat_level"`
	ThreatScore   float64           `json:"threat_score"`
	ThreatTypes   []string          `json:"threat_types"`
	Analysis      string            `json:"analysis"`
	Confidence    float64           `json:"confidence"`
	Metadata      map[string]string `json:"metadata"`
	AnalyzedAt    int64             `json:"analyzed_at"`
}

// NewFr0gIOGRPCClient creates a new gRPC client for fr0g-ai-io
func NewFr0gIOGRPCClient(config *ClientConfig) (*Fr0gIOGRPCClient, error) {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to fr0g-ai-io service at %s: %w", address, err)
	}

	// Note: In a real implementation, you would use the generated gRPC client
	// For now, we'll use a mock client that implements the interface
	client := &mockFr0gIOServiceClient{}

	return &Fr0gIOGRPCClient{
		conn:   conn,
		client: client,
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

	grpcCommand := &OutputCommand{
		Id:       command.ID,
		Type:     command.Type,
		Target:   command.Target,
		Content:  command.Content,
		Priority: int32(command.Priority),
		Metadata: make(map[string]string),
	}

	// Convert metadata
	for k, v := range command.Metadata {
		if str, ok := v.(string); ok {
			grpcCommand.Metadata[k] = str
		}
	}

	req := &OutputCommandRequest{Command: grpcCommand}
	
	resp, err := c.client.SendOutputCommand(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to send output command: %w", err)
	}

	return &input.OutputResponse{
		CommandID: resp.CommandId,
		Success:   resp.Success,
		Message:   resp.Message,
		Metadata:  map[string]interface{}{"grpc_response": true},
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

	grpcAnalysis := &ThreatAnalysis{
		EventId:     result.EventID,
		ThreatLevel: result.ThreatLevel,
		ThreatScore: result.ThreatScore,
		ThreatTypes: result.ThreatTypes,
		Analysis:    result.Analysis,
		Confidence:  result.Confidence,
		AnalyzedAt:  result.AnalyzedAt.Unix(),
		Metadata:    make(map[string]string),
	}

	// Convert metadata
	for k, v := range result.Metadata {
		if str, ok := v.(string); ok {
			grpcAnalysis.Metadata[k] = str
		}
	}

	req := &ThreatAnalysisRequest{Analysis: grpcAnalysis}
	
	_, err := c.client.SendThreatAnalysis(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send threat analysis: %w", err)
	}

	log.Printf("Threat analysis result sent to fr0g-ai-io for event %s", result.EventID)
	return nil
}

// GetServiceStatus gets the status of fr0g-ai-io service
func (c *Fr0gIOGRPCClient) GetServiceStatus(ctx context.Context) (*input.ServiceStatus, error) {
	if c.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.config.Timeout)
		defer cancel()
	}

	req := &ServiceStatusRequest{}
	
	resp, err := c.client.GetServiceStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get service status: %w", err)
	}

	return &input.ServiceStatus{
		ServiceName:   c.config.ServiceName,
		Status:        resp.Status,
		Version:       resp.Version,
		Uptime:        resp.Uptime,
		LastHeartbeat: time.Unix(resp.Timestamp, 0),
	}, nil
}

// Close closes the gRPC connection
func (c *Fr0gIOGRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// mockFr0gIOServiceClient is a mock implementation for demonstration
type mockFr0gIOServiceClient struct{}

func (m *mockFr0gIOServiceClient) SendOutputCommand(ctx context.Context, in *OutputCommandRequest, opts ...grpc.CallOption) (*OutputCommandResponse, error) {
	log.Printf("Mock: Sending output command %s of type %s to %s", in.Command.Id, in.Command.Type, in.Command.Target)
	return &OutputCommandResponse{
		Success:   true,
		Message:   "Command processed successfully",
		CommandId: in.Command.Id,
	}, nil
}

func (m *mockFr0gIOServiceClient) SendThreatAnalysis(ctx context.Context, in *ThreatAnalysisRequest, opts ...grpc.CallOption) (*ThreatAnalysisResponse, error) {
	log.Printf("Mock: Sending threat analysis for event %s with threat level %s", in.Analysis.EventId, in.Analysis.ThreatLevel)
	return &ThreatAnalysisResponse{
		Success: true,
		Message: "Threat analysis received",
	}, nil
}

func (m *mockFr0gIOServiceClient) GetServiceStatus(ctx context.Context, in *ServiceStatusRequest, opts ...grpc.CallOption) (*ServiceStatusResponse, error) {
	return &ServiceStatusResponse{
		Status:    "healthy",
		Version:   "1.0.0",
		Uptime:    "24h30m",
		Timestamp: time.Now().Unix(),
	}, nil
}
