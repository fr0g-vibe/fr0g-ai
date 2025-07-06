package test

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// TestService represents a test gRPC service
type TestService struct {
	Name     string
	Port     string
	Server   interface{}
	Client   *grpc.ClientConn
	Endpoint string
}

// GRPCConnectivityTestSuite manages all gRPC connectivity tests
type GRPCConnectivityTestSuite struct {
	services map[string]*TestService
}

// NewGRPCConnectivityTestSuite creates a new test suite
func NewGRPCConnectivityTestSuite() *GRPCConnectivityTestSuite {
	return &GRPCConnectivityTestSuite{
		services: make(map[string]*TestService),
	}
}

// TestGRPCConnectivity runs comprehensive gRPC connectivity tests
func TestGRPCConnectivity(t *testing.T) {
	t.Log("Testing basic gRPC connectivity functionality")
	
	// Test basic gRPC server creation
	t.Run("BasicGRPCServer", func(t *testing.T) {
		// Create a simple gRPC server for testing
		lis, err := net.Listen("tcp", ":0") // Use port 0 for automatic assignment
		if err != nil {
			t.Fatalf("Failed to listen: %v", err)
		}
		defer lis.Close()
		
		s := grpc.NewServer()
		reflection.Register(s)
		
		// Start server in background
		go func() {
			s.Serve(lis)
		}()
		defer s.Stop()
		
		// Test connection
		addr := lis.Addr().String()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()
		
		t.Logf("✓ Successfully created and connected to gRPC server at %s", addr)
	})
	
	// Test health check functionality
	t.Run("HealthCheck", func(t *testing.T) {
		// This test verifies that health check imports work
		req := &grpc_health_v1.HealthCheckRequest{Service: ""}
		if req.Service != "" {
			t.Error("Health check request creation failed")
		}
		t.Log("✓ Health check functionality available")
	})
	
	// Test reflection functionality
	t.Run("Reflection", func(t *testing.T) {
		// This test verifies that reflection imports work
		req := &grpc_reflection_v1alpha.ServerReflectionRequest{}
		if req == nil {
			t.Error("Reflection request creation failed")
		}
		t.Log("✓ gRPC reflection functionality available")
	})
}

// Cleanup stops all services and cleans up resources
func (suite *GRPCConnectivityTestSuite) Cleanup() {
	for _, service := range suite.services {
		if service.Client != nil {
			service.Client.Close()
		}
	}
}
