package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/types/known/timestamppb"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
	pb "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/pb"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/processors"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-io/internal/outputs"
)

type Server struct {
	pb.UnimplementedIOServiceServer
	config         *sharedconfig.Config
	grpcServer     *grpc.Server
	listener       net.Listener
	processorMgr   *processors.Manager
	outputMgr      *outputs.Manager
}

func NewServer(cfg *sharedconfig.Config, processorMgr *processors.Manager, outputMgr *outputs.Manager) *Server {
	return &Server{
		config:       cfg,
		processorMgr: processorMgr,
		outputMgr:    outputMgr,
	}
}

// HealthCheck implements the gRPC health check
func (s *Server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status:  "healthy",
		Version: "1.0.0",
		Details: map[string]string{
			"service": "fr0g-ai-io",
			"uptime":  time.Now().Format(time.RFC3339),
		},
	}, nil
}

// ProcessInputEvent processes an incoming input event
func (s *Server) ProcessInputEvent(ctx context.Context, req *pb.InputEvent) (*pb.InputEventResponse, error) {
	log.Printf("Processing input event: %s from %s", req.Type, req.Source)
	
	// Convert protobuf to internal event format
	event := &processors.InputEvent{
		ID:        req.Id,
		Type:      req.Type,
		Source:    req.Source,
		Content:   req.Content,
		Metadata:  req.Metadata,
		Timestamp: req.Timestamp.AsTime(),
		Priority:  int(req.Priority),
	}
	
	// Process the event
	result, err := s.processorMgr.ProcessEvent(event)
	if err != nil {
		return nil, fmt.Errorf("failed to process event: %w", err)
	}
	
	// Convert result to protobuf format
	var actions []*pb.OutputCommand
	for _, action := range result.Actions {
		actions = append(actions, &pb.OutputCommand{
			Id:        action.ID,
			Type:      action.Type,
			Target:    action.Target,
			Content:   action.Content,
			Metadata:  action.Metadata,
			Priority:  int32(action.Priority),
			CreatedAt: timestamppb.New(action.CreatedAt),
		})
	}
	
	return &pb.InputEventResponse{
		EventId:     result.EventID,
		Processed:   result.Processed,
		Actions:     actions,
		Metadata:    result.Metadata,
		ProcessedAt: timestamppb.New(result.ProcessedAt),
	}, nil
}

// ExecuteOutputCommand executes an output command
func (s *Server) ExecuteOutputCommand(ctx context.Context, req *pb.OutputCommand) (*pb.OutputResult, error) {
	log.Printf("Executing output command: %s to %s", req.Type, req.Target)
	
	// Convert protobuf to internal command format
	command := &outputs.OutputCommand{
		ID:        req.Id,
		Type:      req.Type,
		Target:    req.Target,
		Content:   req.Content,
		Metadata:  req.Metadata,
		Priority:  int(req.Priority),
		CreatedAt: req.CreatedAt.AsTime(),
	}
	
	// Execute the command
	result, err := s.outputMgr.ExecuteCommand(command)
	if err != nil {
		return &pb.OutputResult{
			CommandId:     req.Id,
			Success:       false,
			ErrorMessage:  err.Error(),
			CompletedAt:   timestamppb.New(time.Now()),
		}, nil
	}
	
	return &pb.OutputResult{
		CommandId:   result.CommandID,
		Success:     result.Success,
		ErrorMessage: result.ErrorMessage,
		Metadata:    result.Metadata,
		CompletedAt: timestamppb.New(result.CompletedAt),
	}, nil
}

// StreamInputEvents handles streaming input events to master-control
func (s *Server) StreamInputEvents(stream pb.IOService_StreamInputEventsServer) error {
	log.Println("Starting input event stream")
	
	for {
		event, err := stream.Recv()
		if err != nil {
			log.Printf("Input stream error: %v", err)
			return err
		}
		
		// Process event and send analysis result
		// This would typically forward to master-control for analysis
		analysisResult := &pb.AnalysisResult{
			EventId:      event.Id,
			AnalysisType: "basic",
			Results: map[string]string{
				"processed": "true",
				"timestamp": time.Now().Format(time.RFC3339),
			},
			AnalyzedAt: timestamppb.New(time.Now()),
		}
		
		if err := stream.Send(analysisResult); err != nil {
			log.Printf("Failed to send analysis result: %v", err)
			return err
		}
	}
}

// StreamOutputCommands handles streaming output commands from master-control
func (s *Server) StreamOutputCommands(stream pb.IOService_StreamOutputCommandsServer) error {
	log.Println("Starting output command stream")
	
	for {
		command, err := stream.Recv()
		if err != nil {
			log.Printf("Output stream error: %v", err)
			return err
		}
		
		// Execute command and send result
		result, err := s.ExecuteOutputCommand(stream.Context(), command)
		if err != nil {
			log.Printf("Failed to execute streamed command: %v", err)
			result = &pb.OutputResult{
				CommandId:     command.Id,
				Success:       false,
				ErrorMessage:  err.Error(),
				CompletedAt:   timestamppb.New(time.Now()),
			}
		}
		
		if err := stream.Send(result); err != nil {
			log.Printf("Failed to send command result: %v", err)
			return err
		}
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.GRPC.Port)
	
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	s.listener = listener

	s.grpcServer = grpc.NewServer()
	
	// Register I/O service
	pb.RegisterIOServiceServer(s.grpcServer, s)
	
	// Register health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s.grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	log.Printf("gRPC server starting on %s", addr)
	
	go func() {
		if err := s.grpcServer.Serve(listener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
	if s.listener != nil {
		s.listener.Close()
	}
	return nil
}
