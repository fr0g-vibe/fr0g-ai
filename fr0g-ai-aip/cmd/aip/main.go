package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/config"
	"github.com/fr0g-vibe/fr0g-ai/pkg/lifecycle"
)

func main() {
	log.Println("🤖 Starting fr0g-ai-aip service...")

	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	log.Printf("✅ Configuration loaded successfully")
	log.Printf("   - HTTP Port: %s", cfg.HTTP.Port)
	log.Printf("   - GRPC Port: %s", cfg.GRPC.Port)
	log.Printf("   - Storage Type: %s", cfg.Storage.Type)
	log.Printf("   - Service Registry: %v", cfg.ServiceRegistry.Enabled)

	// Parse ports for lifecycle manager
	httpPort, _ := strconv.Atoi(cfg.HTTP.Port)
	grpcPort, _ := strconv.Atoi(cfg.GRPC.Port)

	// Initialize lifecycle manager
	serviceConfig := lifecycle.ServiceConfig{
		Name:       cfg.ServiceRegistry.ServiceName,
		ID:         cfg.ServiceRegistry.ServiceID,
		HTTPPort:   httpPort,
		GRPCPort:   grpcPort,
		Tags:       []string{"ai", "personas", "identities", "aip"},
		Meta:       map[string]string{"version": "1.0.0", "service": "aip"},
		HealthPath: "/health",
	}

	lifecycleManager := lifecycle.NewLifecycleManager(serviceConfig)

	// Start lifecycle management (service registration)
	if err := lifecycleManager.Start(); err != nil {
		log.Printf("⚠️  Failed to start lifecycle management: %v", err)
	}

	// Create HTTP server with health endpoint
	mux := http.NewServeMux()
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"fr0g-ai-aip","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
	})

	// Personas endpoint (placeholder)
	mux.HandleFunc("/personas", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"personas":[],"message":"AIP service is running"}`)
	})

	// Status endpoint
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"operational","uptime":"%s","config":{"storage":"%s","registry":"%v"}}`, 
			time.Since(time.Now()).String(), cfg.Storage.Type, cfg.ServiceRegistry.Enabled)
	})

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:      mux,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	// Start HTTP server in goroutine
	go func() {
		log.Printf("🚀 HTTP server starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	log.Printf("✅ fr0g-ai-aip service is now operational!")
	log.Printf("🔗 Health check: http://%s:%s/health", cfg.HTTP.Host, cfg.HTTP.Port)
	log.Printf("📊 Status endpoint: http://%s:%s/status", cfg.HTTP.Host, cfg.HTTP.Port)
	log.Printf("🤖 Personas endpoint: http://%s:%s/personas", cfg.HTTP.Host, cfg.HTTP.Port)

	// Wait for shutdown signal
	lifecycleManager.WaitForShutdown()

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

	log.Println("✅ fr0g-ai-aip service stopped gracefully")
}
