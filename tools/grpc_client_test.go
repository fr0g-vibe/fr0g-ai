package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func main() {
	var (
		address = flag.String("addr", "localhost:9090", "gRPC server address")
		timeout = flag.Duration("timeout", 5*time.Second, "connection timeout")
		test    = flag.String("test", "health", "test type: health, reflection, or connectivity")
	)
	flag.Parse()

	fmt.Printf("Testing gRPC server at %s\n", *address)
	fmt.Printf("Test type: %s\n", *test)
	fmt.Printf("Timeout: %v\n", *timeout)
	fmt.Println("=" + fmt.Sprintf("%*s", 50, "="))

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// Connect to the gRPC server
	conn, err := grpc.DialContext(ctx, *address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	fmt.Printf("✓ Connected to %s\n", *address)
	fmt.Printf("Connection state: %s\n", conn.GetState())

	switch *test {
	case "health":
		testHealthCheck(ctx, conn)
	case "reflection":
		testReflection(ctx, conn)
	case "connectivity":
		testConnectivity(ctx, conn)
	case "all":
		testHealthCheck(ctx, conn)
		fmt.Println()
		testReflection(ctx, conn)
		fmt.Println()
		testConnectivity(ctx, conn)
	default:
		fmt.Printf("Unknown test type: %s\n", *test)
		fmt.Println("Available tests: health, reflection, connectivity, all")
	}
}

func testHealthCheck(ctx context.Context, conn *grpc.ClientConn) {
	fmt.Println("Testing Health Check...")
	fmt.Println("-" + fmt.Sprintf("%*s", 20, "-"))

	client := grpc_health_v1.NewHealthClient(conn)

	req := &grpc_health_v1.HealthCheckRequest{
		Service: "", // Empty string checks overall server health
	}

	resp, err := client.Check(ctx, req)
	if err != nil {
		fmt.Printf("✗ Health check failed: %v\n", err)
		return
	}

	fmt.Printf("✓ Health check passed\n")
	fmt.Printf("Status: %s\n", resp.Status)

	// Test watch functionality
	fmt.Println("Testing health watch...")
	watchCtx, watchCancel := context.WithTimeout(ctx, 2*time.Second)
	defer watchCancel()

	stream, err := client.Watch(watchCtx, req)
	if err != nil {
		fmt.Printf("✗ Health watch failed: %v\n", err)
		return
	}

	// Try to receive one message
	watchResp, err := stream.Recv()
	if err != nil {
		fmt.Printf("✗ Health watch receive failed: %v\n", err)
		return
	}

	fmt.Printf("✓ Health watch working\n")
	fmt.Printf("Watch status: %s\n", watchResp.Status)
}

func testReflection(ctx context.Context, conn *grpc.ClientConn) {
	fmt.Println("Testing gRPC Reflection...")
	fmt.Println("-" + fmt.Sprintf("%*s", 25, "-"))

	client := grpc_reflection_v1alpha.NewServerReflectionClient(conn)

	stream, err := client.ServerReflectionInfo(ctx)
	if err != nil {
		fmt.Printf("✗ Reflection stream failed: %v\n", err)
		return
	}

	// Request list of services
	req := &grpc_reflection_v1alpha.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_ListServices{
			ListServices: "",
		},
	}

	err = stream.Send(req)
	if err != nil {
		fmt.Printf("✗ Failed to send reflection request: %v\n", err)
		return
	}

	resp, err := stream.Recv()
	if err != nil {
		fmt.Printf("✗ Failed to receive reflection response: %v\n", err)
		return
	}

	if listResp := resp.GetListServicesResponse(); listResp != nil {
		fmt.Printf("✓ Reflection working\n")
		fmt.Printf("Available services:\n")
		for _, service := range listResp.Service {
			fmt.Printf("  - %s\n", service.Name)
		}
	} else {
		fmt.Printf("✗ Unexpected reflection response type\n")
	}
}

func testConnectivity(ctx context.Context, conn *grpc.ClientConn) {
	fmt.Println("Testing Basic Connectivity...")
	fmt.Println("-" + fmt.Sprintf("%*s", 28, "-"))

	// Test connection state
	state := conn.GetState()
	fmt.Printf("Connection state: %s\n", state)

	if state == grpc.Ready {
		fmt.Printf("✓ Connection is ready\n")
	} else if state == grpc.Connecting {
		fmt.Printf("⚠ Connection is still connecting\n")
		
		// Wait for connection to be ready
		if conn.WaitForStateChange(ctx, state) {
			newState := conn.GetState()
			fmt.Printf("New connection state: %s\n", newState)
			if newState == grpc.Ready {
				fmt.Printf("✓ Connection became ready\n")
			}
		}
	} else {
		fmt.Printf("✗ Connection is not ready: %s\n", state)
	}

	// Test basic RPC call (using health check as a simple test)
	fmt.Println("Testing basic RPC call...")
	client := grpc_health_v1.NewHealthClient(conn)
	
	_, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		fmt.Printf("✗ Basic RPC call failed: %v\n", err)
	} else {
		fmt.Printf("✓ Basic RPC call succeeded\n")
	}
}
