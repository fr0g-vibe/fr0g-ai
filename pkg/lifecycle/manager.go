package lifecycle

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/registry"
)

// ServiceConfig represents the configuration for a service
type ServiceConfig struct {
	Name       string
	ID         string
	HTTPPort   int
	GRPCPort   int
	Tags       []string
	Meta       map[string]string
	HealthPath string
}

// LifecycleManager manages service registration, health monitoring, and graceful shutdown
type LifecycleManager struct {
	config         ServiceConfig
	registryClient *registry.RegistryClient
	shutdownChan   chan os.Signal
	ctx            context.Context
	cancel         context.CancelFunc
}

// NewLifecycleManager creates a new lifecycle manager
func NewLifecycleManager(config ServiceConfig) *LifecycleManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &LifecycleManager{
		config:       config,
		shutdownChan: make(chan os.Signal, 1),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start initializes service registration and starts lifecycle management
func (lm *LifecycleManager) Start() error {
	// Initialize service registry if enabled
	if os.Getenv("SERVICE_REGISTRY_ENABLED") == "true" {
		registryURL := os.Getenv("SERVICE_REGISTRY_URL")
		if registryURL == "" {
			registryURL = "http://localhost:8500"
		}
		
		log.Printf("🔗 Initializing service registry client: %s", registryURL)
		lm.registryClient = registry.NewRegistryClient(registryURL, nil)
		
		// Register service
		if err := lm.registerService(); err != nil {
			log.Printf("⚠️  Failed to register service: %v", err)
			return err
		}
		
		log.Printf("✅ Service registered successfully: %s", lm.config.ID)
	}
	
	// Setup graceful shutdown
	signal.Notify(lm.shutdownChan, syscall.SIGINT, syscall.SIGTERM)
	
	return nil
}

// registerService registers the service with the registry
func (lm *LifecycleManager) registerService() error {
	healthPath := lm.config.HealthPath
	if healthPath == "" {
		healthPath = "/health"
	}
	
	serviceInfo := &registry.ServiceInfo{
		ID:      lm.config.ID,
		Name:    lm.config.Name,
		Address: "localhost",
		Port:    lm.config.HTTPPort,
		Tags:    lm.config.Tags,
		Meta:    lm.config.Meta,
		Check: &registry.HealthCheck{
			HTTP:     fmt.Sprintf("http://localhost:%d%s", lm.config.HTTPPort, healthPath),
			Interval: "30s",
			Timeout:  "10s",
		},
	}
	
	// Ensure metadata includes port information
	if serviceInfo.Meta == nil {
		serviceInfo.Meta = make(map[string]string)
	}
	serviceInfo.Meta["http_port"] = strconv.Itoa(lm.config.HTTPPort)
	serviceInfo.Meta["grpc_port"] = strconv.Itoa(lm.config.GRPCPort)
	
	return lm.registryClient.RegisterService(serviceInfo)
}

// WaitForShutdown blocks until a shutdown signal is received
func (lm *LifecycleManager) WaitForShutdown() {
	<-lm.shutdownChan
	log.Println("🛑 Shutdown signal received")
}

// Shutdown performs graceful shutdown including service deregistration
func (lm *LifecycleManager) Shutdown() error {
	log.Println("🛑 Starting graceful shutdown...")
	
	// Deregister from service registry
	if lm.registryClient != nil {
		log.Println("🔗 Deregistering from service registry...")
		if err := lm.registryClient.DeregisterService(); err != nil {
			log.Printf("⚠️  Failed to deregister service: %v", err)
			return err
		}
		log.Println("✅ Service deregistered successfully")
	}
	
	// Cancel context to stop background operations
	lm.cancel()
	
	log.Println("✅ Graceful shutdown completed")
	return nil
}

// GetServiceConfig creates a service config from environment variables
func GetServiceConfig(defaultName string, defaultHTTPPort, defaultGRPCPort int) ServiceConfig {
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = defaultName
	}
	
	serviceID := os.Getenv("SERVICE_ID")
	if serviceID == "" {
		serviceID = serviceName + "-1"
	}
	
	httpPort := defaultHTTPPort
	if portStr := os.Getenv("HTTP_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			httpPort = port
		}
	}
	
	grpcPort := defaultGRPCPort
	if portStr := os.Getenv("GRPC_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			grpcPort = port
		}
	}
	
	return ServiceConfig{
		Name:       serviceName,
		ID:         serviceID,
		HTTPPort:   httpPort,
		GRPCPort:   grpcPort,
		Tags:       []string{"ai", "fr0g"},
		Meta:       map[string]string{"version": "1.0.0"},
		HealthPath: "/health",
	}
}
