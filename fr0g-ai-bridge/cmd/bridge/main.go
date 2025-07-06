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

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/config"
	"github.com/fr0g-vibe/fr0g-ai/pkg/lifecycle"
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

	// Initialize lifecycle manager
	serviceConfig := lifecycle.ServiceConfig{
		Name:       cfg.ServiceRegistry.ServiceName,
		ID:         cfg.ServiceRegistry.ServiceID,
		HTTPPort:   cfg.Server.HTTPPort,
		GRPCPort:   cfg.Server.GRPCPort,
		Tags:       []string{"ai", "bridge", "openwebui", "chat"},
		Meta:       map[string]string{"version": "1.0.0", "service": "bridge"},
		HealthPath: "/health",
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

	// Shutdown lifecycle manager (deregister service)
	if err := lifecycleManager.Shutdown(); err != nil {
		log.Printf("Lifecycle manager shutdown error: %v", err)
	}

	log.Println("✅ fr0g-ai-bridge service stopped gracefully")
}
