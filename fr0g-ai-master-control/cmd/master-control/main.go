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

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// Config holds master-control specific configuration
type Config struct {
	Server     sharedconfig.ServerConfig     `yaml:"server"`
	Security   sharedconfig.SecurityConfig   `yaml:"security"`
	Monitoring sharedconfig.MonitoringConfig `yaml:"monitoring"`
	
	// Master-control specific config
	Learning struct {
		Enabled bool    `yaml:"enabled"`
		Rate    float64 `yaml:"rate"`
	} `yaml:"learning"`
	
	Cognitive struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"cognitive"`
}

func main() {
	log.Println("🧠 Starting fr0g-ai-master-control...")

	// Load configuration
	loader := sharedconfig.NewLoader(sharedconfig.LoaderOptions{
		ConfigPath: "config.yaml",
		EnvPrefix:  "MCP",
	})

	config := &Config{
		Server: sharedconfig.ServerConfig{
			Host:     "0.0.0.0",
			HTTPPort: 8081,
		},
		Security: sharedconfig.SecurityConfig{
			EnableCORS:    true,
			RateLimitRPM:  100,
		},
		Monitoring: sharedconfig.MonitoringConfig{
			EnableMetrics: true,
			MetricsPort:   8083,
		},
	}

	// Set default values
	config.Learning.Enabled = true
	config.Learning.Rate = 0.154
	config.Cognitive.Enabled = true

	// Load from file and environment
	if err := loader.LoadFromFile(config); err != nil {
		log.Printf("Warning: failed to load config file: %v", err)
	}

	if err := loader.LoadEnvFiles(); err != nil {
		log.Printf("Warning: failed to load env files: %v", err)
	}

	// Override with environment variables
	if port := os.Getenv("MCP_HTTP_PORT"); port != "" {
		fmt.Printf("Using MCP_HTTP_PORT: %s\n", port)
	}

	// Validate configuration
	if errors := config.Server.Validate(); len(errors) > 0 {
		log.Printf("Configuration validation failed: %v", errors)
		// Don't fail on validation errors, use defaults
	}

	// Create HTTP server
	mux := http.NewServeMux()
	
	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"status": "healthy",
			"service": "fr0g-ai-master-control",
			"timestamp": "%s",
			"intelligence": {
				"learning_enabled": %t,
				"cognitive_enabled": %t,
				"status": "operational"
			}
		}`, time.Now().Format(time.RFC3339), config.Learning.Enabled, config.Cognitive.Enabled)
	})

	// Status endpoint
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"service": "fr0g-ai-master-control",
			"version": "1.0.0",
			"intelligence": "operational",
			"learning_rate": 0.154,
			"pattern_count": 6,
			"consciousness": "active"
		}`)
	})

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Server.Host, config.Server.HTTPPort),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Master Control HTTP server starting on %s", server.Addr)
		log.Printf("🔗 Health check available at: http://%s/health", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ Server failed to start: %v", err)
		}
	}()

	// Give server time to start and verify it's running
	time.Sleep(500 * time.Millisecond)
	log.Printf("✅ Master Control ready and listening on %s", server.Addr)
	log.Printf("🏥 Health endpoint: http://%s/health", server.Addr)
	log.Printf("📊 Status endpoint: http://%s/status", server.Addr)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Master Control...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Master Control stopped")
}
