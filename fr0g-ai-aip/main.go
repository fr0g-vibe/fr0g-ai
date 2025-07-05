package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/api"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/config"
	grpcserver "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/grpc"
	pb "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Start gRPC server
	grpcServer := grpc.NewServer()
	grpcPersonaServer := grpcserver.NewServer()
	pb.RegisterPersonaServiceServer(grpcServer, grpcPersonaServer)

	grpcPort := cfg.GetString("GRPC_PORT", "50051")
	grpcListener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port %s: %v", grpcPort, err)
	}

	go func() {
		log.Printf("gRPC server starting on port %s", grpcPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Start REST API server
	restServer := api.NewServer(cfg, nil)
	restPort := cfg.GetString("REST_PORT", "8080")
	
	httpServer := &http.Server{
		Addr:    ":" + restPort,
		Handler: restServer,
	}

	go func() {
		log.Printf("REST API server starting on port %s", restPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("REST server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down servers...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown REST server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("REST server shutdown error: %v", err)
	}

	// Shutdown gRPC server
	grpcServer.GracefulStop()

	log.Println("Servers stopped")
}
