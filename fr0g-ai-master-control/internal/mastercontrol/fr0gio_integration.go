package mastercontrol

import (
	"context"
	"fmt"
	"log"
	"time"

	"fr0g-ai-master-control/internal/grpc"
	"fr0g-ai-master-control/internal/mastercontrol/input"
	"fr0g-ai-master-control/internal/mastercontrol/workflow"
)

// Fr0gIOIntegration manages the integration with fr0g-ai-io service
type Fr0gIOIntegration struct {
	config       *MCPConfig
	grpcClient   input.Fr0gIOClient
	grpcServer   *grpc.Fr0gIOInputServer
	inputHandler input.Fr0gIOInputHandler
	workflowEngine *workflow.WorkflowEngine
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewFr0gIOIntegration creates a new fr0g-ai-io integration
func NewFr0gIOIntegration(config *MCPConfig) (*Fr0gIOIntegration, error) {
	ctx, cancel := context.WithCancel(context.Background())

	integration := &Fr0gIOIntegration{
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}

	// Initialize components if fr0g-ai-io is enabled
	if config.Fr0gIOService.Enabled {
		if err := integration.initializeComponents(); err != nil {
			cancel()
			return nil, fmt.Errorf("failed to initialize fr0g-ai-io integration: %w", err)
		}
	}

	return integration, nil
}

// initializeComponents initializes all integration components
func (f *Fr0gIOIntegration) initializeComponents() error {
	// Initialize gRPC client for sending commands to fr0g-ai-io
	clientConfig := &grpc.ClientConfig{
		Host:        f.config.Fr0gIOService.GRPCHost,
		Port:        f.config.Fr0gIOService.GRPCPort,
		Timeout:     30 * time.Second,
		MaxRetries:  3,
		ServiceName: f.config.Fr0gIOService.ServiceName,
	}

	var err error
	f.grpcClient, err = grpc.NewFr0gIOGRPCClient(clientConfig)
	if err != nil {
		return fmt.Errorf("failed to create gRPC client: %w", err)
	}

	// Initialize workflow engine with fr0g-ai-io client
	workflowConfig := &workflow.WorkflowConfig{
		MaxConcurrentWorkflows: f.config.MaxConcurrentWorkflows,
		WorkflowTimeout:        5 * time.Minute,
		AutoStartWorkflows:     true,
		WorkflowInterval:       2 * time.Minute,
	}

	// Create input handler
	f.inputHandler = input.NewFr0gIOInputHandler(nil) // Will be set after workflow engine creation

	// Create workflow engine
	f.workflowEngine = workflow.NewWorkflowEngine(workflowConfig, f.grpcClient, f.inputHandler)

	// Update input handler with workflow engine reference
	if handlerImpl, ok := f.inputHandler.(*input.Fr0gIOInputHandlerImpl); ok {
		handlerImpl.SetWorkflowEngine(f.workflowEngine)
	}

	// Initialize gRPC server for receiving input events from fr0g-ai-io
	serverConfig := &grpc.ServerConfig{
		Host: f.config.Input.GRPC.Host,
		Port: f.config.Input.GRPC.Port,
	}

	f.grpcServer = grpc.NewFr0gIOInputServer(serverConfig, f.inputHandler)

	log.Printf("Fr0g-AI-IO Integration: Initialized with client %s:%d and server %s:%d",
		clientConfig.Host, clientConfig.Port, serverConfig.Host, serverConfig.Port)

	return nil
}

// Start starts the fr0g-ai-io integration
func (f *Fr0gIOIntegration) Start() error {
	if !f.config.Fr0gIOService.Enabled {
		log.Println("Fr0g-AI-IO Integration: Service disabled, skipping startup")
		return nil
	}

	log.Println("Fr0g-AI-IO Integration: Starting integration services...")

	// Start workflow engine
	if f.workflowEngine != nil {
		if err := f.workflowEngine.Start(); err != nil {
			return fmt.Errorf("failed to start workflow engine: %w", err)
		}
	}

	// Start gRPC server for input events
	if f.grpcServer != nil {
		if err := f.grpcServer.Start(); err != nil {
			return fmt.Errorf("failed to start gRPC input server: %w", err)
		}
	}

	// Test connection to fr0g-ai-io service
	go f.testServiceConnection()

	// Start health monitoring
	go f.monitorServiceHealth()

	log.Println("Fr0g-AI-IO Integration: All services started successfully")
	return nil
}

// Stop stops the fr0g-ai-io integration
func (f *Fr0gIOIntegration) Stop() error {
	log.Println("Fr0g-AI-IO Integration: Stopping integration services...")

	f.cancel()

	// Stop workflow engine
	if f.workflowEngine != nil {
		if err := f.workflowEngine.Stop(); err != nil {
			log.Printf("Fr0g-AI-IO Integration: Error stopping workflow engine: %v", err)
		}
	}

	// Stop gRPC server
	if f.grpcServer != nil {
		f.grpcServer.Stop()
	}

	// Close gRPC client
	if f.grpcClient != nil {
		if closer, ok := f.grpcClient.(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				log.Printf("Fr0g-AI-IO Integration: Error closing gRPC client: %v", err)
			}
		}
	}

	log.Println("Fr0g-AI-IO Integration: All services stopped")
	return nil
}

// testServiceConnection tests the connection to fr0g-ai-io service
func (f *Fr0gIOIntegration) testServiceConnection() {
	if f.grpcClient == nil {
		return
	}

	ctx, cancel := context.WithTimeout(f.ctx, 10*time.Second)
	defer cancel()

	status, err := f.grpcClient.GetServiceStatus(ctx)
	if err != nil {
		log.Printf("Fr0g-AI-IO Integration: Failed to connect to service: %v", err)
		return
	}

	log.Printf("Fr0g-AI-IO Integration: Successfully connected to %s (status: %s, version: %s)",
		status.ServiceName, status.Status, status.Version)
}

// monitorServiceHealth monitors the health of fr0g-ai-io service
func (f *Fr0gIOIntegration) monitorServiceHealth() {
	if f.grpcClient == nil {
		return
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-f.ctx.Done():
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(f.ctx, 5*time.Second)
			
			status, err := f.grpcClient.GetServiceStatus(ctx)
			if err != nil {
				log.Printf("Fr0g-AI-IO Integration: Health check failed: %v", err)
			} else {
				log.Printf("Fr0g-AI-IO Integration: Health check passed - %s is %s", 
					status.ServiceName, status.Status)
			}
			
			cancel()
		}
	}
}

// GetWorkflowEngine returns the workflow engine for external access
func (f *Fr0gIOIntegration) GetWorkflowEngine() *workflow.WorkflowEngine {
	return f.workflowEngine
}

// GetGRPCClient returns the gRPC client for external access
func (f *Fr0gIOIntegration) GetGRPCClient() input.Fr0gIOClient {
	return f.grpcClient
}

// GetInputHandler returns the input handler for external access
func (f *Fr0gIOIntegration) GetInputHandler() input.Fr0gIOInputHandler {
	return f.inputHandler
}
