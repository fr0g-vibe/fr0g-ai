package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/config"
	"google.golang.org/grpc"
)

func main() {
	log.Println("🌉 Starting fr0g.ai Bridge Service")
	log.Println("==================================")
	
	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	// Validate configuration
	if validationErrors := cfg.Validate(); len(validationErrors) > 0 {
		log.Fatalf("Configuration validation failed: %v", validationErrors)
	}
	
	log.Println("✅ Configuration loaded and validated successfully")
	
	// Start gRPC server
	grpcServer := grpc.NewServer()
	// TODO: Register bridge services here
	
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port %d: %v", cfg.Server.GRPCPort, err)
	}
	
	go func() {
		log.Printf("✅ gRPC server starting on port %d", cfg.Server.GRPCPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()
	
	// Start HTTP server
	mux := http.NewServeMux()
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"fr0g-ai-bridge"}`))
	})
	
	// TODO: Add bridge-specific endpoints here
	
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.HTTPPort),
		Handler: mux,
	}
	
	go func() {
		log.Printf("✅ HTTP server starting on port %d", cfg.Server.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()
	
	log.Println("🎯 fr0g.ai Bridge is running...")
	
	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	
	log.Println("🛑 Shutting down servers...")
	
	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	
	// Shutdown HTTP server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
	
	// Shutdown gRPC server
	grpcServer.GracefulStop()
	
	log.Println("👋 fr0g.ai Bridge shutdown complete")
}
