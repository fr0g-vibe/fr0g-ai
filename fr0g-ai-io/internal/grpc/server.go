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
)

type Server struct {
	config     *sharedconfig.Config
	grpcServer *grpc.Server
	listener   net.Listener
}

func NewServer(cfg *sharedconfig.Config) *Server {
	return &Server{
		config: cfg,
	}
}

// HealthCheck implements a basic health check
func (s *Server) HealthCheck() map[string]string {
	return map[string]string{
		"status":  "healthy",
		"service": "fr0g-ai-io",
		"uptime":  time.Now().Format(time.RFC3339),
	}
}

func (s *Server) Start() error {
	if s.config == nil || s.config.GRPC.Port == "" {
		return fmt.Errorf("invalid server configuration")
	}

	addr := fmt.Sprintf(":%s", s.config.GRPC.Port)
	
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	s.listener = listener

	s.grpcServer = grpc.NewServer()
	
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
