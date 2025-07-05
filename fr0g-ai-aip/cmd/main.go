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
	grpcserver "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/grpc"
	pb "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/grpc/pb"
	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration using shared config system
	loader := sharedconfig.NewLoader(sharedconfig.LoaderOptions{
		ConfigPath: "",
		EnvPrefix:  "FR0G_AIP",
	})
	
	// Load environment files
	if err := loader.LoadEnvFiles(); err != nil {
		log.Printf("Warning: failed to load env files: %v", err)
	}
	
	// Create default config
	cfg := &sharedconfig.Config{
		HTTP: sharedconfig.HTTPConfig{
			Port: "8080",
		},
		GRPC: sharedconfig.GRPCConfig{
			Port: "9090",
		},
	}
	
	// Load from file
	if err := loader.LoadFromFile(cfg); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	log.Println("🚀 Starting fr0g.ai AIP servers...")
	
	// Start gRPC server
	grpcServer := grpc.NewServer()
	grpcPersonaServer := grpcserver.NewServer()
	pb.RegisterPersonaServiceServer(grpcServer, grpcPersonaServer)

	grpcPort := cfg.GRPC.Port
	grpcListener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port %s: %v", grpcPort, err)
	}

	go func() {
		log.Printf("✅ gRPC server starting on port %s", grpcPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Start REST API server
	restServer := api.NewServer(cfg, nil)
	restPort := cfg.HTTP.Port
	
	httpServer := &http.Server{
		Addr:    ":" + restPort,
		Handler: restServer,
	}

	go func() {
		log.Printf("✅ REST API server starting on port %s", restPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("REST server error: %v", err)
		}
	}()
	
	log.Println("🎯 fr0g.ai AIP is running...")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("🛑 Shutting down servers...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown REST server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("REST server shutdown error: %v", err)
	}

	// Shutdown gRPC server
	grpcServer.GracefulStop()

	log.Println("👋 fr0g.ai AIP shutdown complete")
}
