package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/pkg/lifecycle"
)

func main() {
	log.Println("📡 Starting fr0g-ai-io service...")

	// Get configuration from environment
	httpPort := 8084
	grpcPort := 9094
	
	if port := os.Getenv("HTTP_PORT"); port != "" {
		if p, err := fmt.Sscanf(port, "%d", &httpPort); err != nil || p != 1 {
			log.Printf("Invalid HTTP_PORT: %s, using default: %d", port, httpPort)
		}
	}
	
	if port := os.Getenv("GRPC_PORT"); port != "" {
		if p, err := fmt.Sscanf(port, "%d", &grpcPort); err != nil || p != 1 {
			log.Printf("Invalid GRPC_PORT: %s, using default: %d", port, grpcPort)
		}
	}

	log.Printf("✅ Configuration loaded successfully")
	log.Printf("   - HTTP Port: %d", httpPort)
	log.Printf("   - GRPC Port: %d", grpcPort)

	// Initialize lifecycle manager
	serviceConfig := lifecycle.ServiceConfig{
		Name:       "fr0g-ai-io",
		ID:         "fr0g-ai-io-1",
		HTTPPort:   httpPort,
		GRPCPort:   grpcPort,
		Tags:       []string{"ai", "io", "outputs", "communications"},
		Meta:       map[string]string{"version": "1.0.0", "service": "io"},
		HealthPath: "/health",
	}

	// Override from environment
	if serviceName := os.Getenv("SERVICE_NAME"); serviceName != "" {
		serviceConfig.Name = serviceName
	}
	if serviceID := os.Getenv("SERVICE_ID"); serviceID != "" {
		serviceConfig.ID = serviceID
	}

	lifecycleManager := lifecycle.NewLifecycleManager(serviceConfig)

	// Start lifecycle management (service registration)
	if err := lifecycleManager.Start(); err != nil {
		log.Printf("⚠️  Failed to start lifecycle management: %v", err)
	}

	// Create HTTP server with endpoints
	mux := http.NewServeMux()
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"fr0g-ai-io","timestamp":"%s"}`, 
			time.Now().Format(time.RFC3339))
	})

	// Output commands endpoint (placeholder)
	mux.HandleFunc("/outputs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"outputs":[],"message":"IO service is running"}`)
	})

	// Status endpoint
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"operational","supported_outputs":["sms","voice","irc","discord","email"]}`)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", httpPort),
		Handler: mux,
	}

	// Start HTTP server in goroutine
	go func() {
		log.Printf("🚀 HTTP server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	log.Printf("✅ fr0g-ai-io service is now operational!")
	log.Printf("🔗 Health check: http://localhost:%d/health", httpPort)
	log.Printf("📊 Status endpoint: http://localhost:%d/status", httpPort)
	log.Printf("📡 Outputs endpoint: http://localhost:%d/outputs", httpPort)

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

	// Shutdown lifecycle manager (deregister service)
	if err := lifecycleManager.Shutdown(); err != nil {
		log.Printf("Lifecycle manager shutdown error: %v", err)
	}

	log.Println("✅ fr0g-ai-io service stopped gracefully")
}
