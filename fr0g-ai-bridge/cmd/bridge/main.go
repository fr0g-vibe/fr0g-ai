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

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/config"
)

func main() {
	log.Println("🌉 Starting fr0g-ai-bridge service...")

	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate configuration
	if validationErrors := cfg.Validate(); len(validationErrors) > 0 {
		log.Fatalf("Configuration validation failed: %v", validationErrors)
	}

	log.Printf("✅ Configuration loaded successfully")
	log.Printf("   - HTTP Port: %d", cfg.Server.HTTPPort)
	log.Printf("   - GRPC Port: %d", cfg.Server.GRPCPort)
	log.Printf("   - OpenWebUI: %s", cfg.OpenWebUI.BaseURL)
	log.Printf("   - Service Registry: %v", cfg.ServiceRegistry.Enabled)

	// Simple service registration (removed lifecycle manager dependency)
	log.Printf("✅ Service configuration loaded")
	log.Printf("   - Service Name: %s", cfg.ServiceRegistry.ServiceName)
	log.Printf("   - Service ID: %s", cfg.ServiceRegistry.ServiceID)

	// Create HTTP server with endpoints
	mux := http.NewServeMux()
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"fr0g-ai-bridge","timestamp":"%s","openwebui":"%s"}`, 
			time.Now().Format(time.RFC3339), cfg.OpenWebUI.BaseURL)
	})

	// Chat completions endpoint (placeholder)
	mux.HandleFunc("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"Bridge service is running","openwebui_url":"%s"}`, cfg.OpenWebUI.BaseURL)
	})

	// Status endpoint
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"operational","openwebui":"%s","registry":"%v","cors":"%v"}`, 
			cfg.OpenWebUI.BaseURL, cfg.ServiceRegistry.Enabled, cfg.Security.EnableCORS)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.HTTPPort),
		Handler: mux,
	}

	// Start HTTP server in goroutine
	go func() {
		log.Printf("🚀 HTTP server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Start gRPC server for bridge
	go func() {
		grpcAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.GRPCPort)
		log.Printf("🚀 gRPC server starting on %s", grpcAddr)
		if err := startBridgeGRPCServer(grpcAddr); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Give servers time to start
	time.Sleep(2 * time.Second)

	log.Printf("✅ fr0g-ai-bridge service is now operational!")
	log.Printf("🔗 Health check: http://%s:%d/health", cfg.Server.Host, cfg.Server.HTTPPort)
	log.Printf("📊 Status endpoint: http://%s:%d/status", cfg.Server.Host, cfg.Server.HTTPPort)
	log.Printf("💬 Chat endpoint: http://%s:%d/v1/chat/completions", cfg.Server.Host, cfg.Server.HTTPPort)

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutdown signal received")

	// Graceful shutdown
	log.Println("🛑 Starting graceful shutdown...")

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Service cleanup completed
	log.Println("🛑 Service cleanup completed")

	log.Println("✅ fr0g-ai-bridge service stopped gracefully")
}

// startBridgeGRPCServer starts a gRPC server for the bridge service
func startBridgeGRPCServer(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	s := grpc.NewServer()

	// Register health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	// Enable reflection for testing
	reflection.Register(s)

	log.Printf("Bridge gRPC server listening on %s", addr)
	return s.Serve(lis)
}
