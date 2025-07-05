package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// ServiceInfo represents a registered service
type ServiceInfo struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Tags     []string          `json:"tags,omitempty"`
	Meta     map[string]string `json:"meta,omitempty"`
	Health   string            `json:"health"` // "passing", "warning", "critical"
	LastSeen time.Time         `json:"last_seen"`
}

// RegistryConfig holds configuration for the service registry
type RegistryConfig struct {
	Port           int
	Host           string
	HealthInterval time.Duration
	ServiceTTL     time.Duration
	EnableHTTPAPI  bool
}

// ServiceRegistry manages service registration and discovery
type ServiceRegistry struct {
	config   *RegistryConfig
	services map[string]*ServiceInfo
	mu       sync.RWMutex
	server   *http.Server
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry(config *RegistryConfig) *ServiceRegistry {
	return &ServiceRegistry{
		config:   config,
		services: make(map[string]*ServiceInfo),
	}
}

// Start starts the service registry HTTP server
func (sr *ServiceRegistry) Start(ctx context.Context) error {
	if !sr.config.EnableHTTPAPI {
		return nil
	}

	router := mux.NewRouter()
	
	// Health check endpoint
	router.HandleFunc("/health", sr.healthHandler).Methods("GET")
	
	// Service registration endpoints
	router.HandleFunc("/v1/agent/service/register", sr.registerHandler).Methods("PUT")
	router.HandleFunc("/v1/agent/service/deregister/{serviceId}", sr.deregisterHandler).Methods("PUT")
	
	// Service discovery endpoints
	router.HandleFunc("/v1/health/service/{service}", sr.serviceHealthHandler).Methods("GET")
	router.HandleFunc("/v1/catalog/services", sr.servicesHandler).Methods("GET")
	router.HandleFunc("/v1/catalog/service/{service}", sr.serviceHandler).Methods("GET")

	sr.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", sr.config.Host, sr.config.Port),
		Handler: router,
	}

	// Start cleanup routine
	go sr.cleanupRoutine(ctx)

	go func() {
		if err := sr.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Registry server error: %v", err)
		}
	}()

	return nil
}

// Stop stops the service registry
func (sr *ServiceRegistry) Stop(ctx context.Context) error {
	if sr.server != nil {
		return sr.server.Shutdown(ctx)
	}
	return nil
}

// Health check handler
func (sr *ServiceRegistry) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Service registration handler
func (sr *ServiceRegistry) registerHandler(w http.ResponseWriter, r *http.Request) {
	var service ServiceInfo
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if service.ID == "" || service.Name == "" {
		http.Error(w, "Service ID and Name are required", http.StatusBadRequest)
		return
	}

	service.Health = "passing"
	service.LastSeen = time.Now()

	sr.mu.Lock()
	sr.services[service.ID] = &service
	sr.mu.Unlock()

	log.Printf("Registered service: %s (%s)", service.Name, service.ID)
	w.WriteHeader(http.StatusOK)
}

// Service deregistration handler
func (sr *ServiceRegistry) deregisterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := vars["serviceId"]

	sr.mu.Lock()
	delete(sr.services, serviceID)
	sr.mu.Unlock()

	log.Printf("Deregistered service: %s", serviceID)
	w.WriteHeader(http.StatusOK)
}

// Service health handler
func (sr *ServiceRegistry) serviceHealthHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["service"]

	sr.mu.RLock()
	var healthyServices []*ServiceInfo
	for _, service := range sr.services {
		if service.Name == serviceName && service.Health == "passing" {
			healthyServices = append(healthyServices, service)
		}
	}
	sr.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(healthyServices)
}

// Services list handler
func (sr *ServiceRegistry) servicesHandler(w http.ResponseWriter, r *http.Request) {
	sr.mu.RLock()
	services := make(map[string][]string)
	for _, service := range sr.services {
		services[service.Name] = service.Tags
	}
	sr.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

// Service instances handler
func (sr *ServiceRegistry) serviceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["service"]

	sr.mu.RLock()
	var serviceInstances []*ServiceInfo
	for _, service := range sr.services {
		if service.Name == serviceName {
			serviceInstances = append(serviceInstances, service)
		}
	}
	sr.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serviceInstances)
}

// Cleanup routine to remove stale services
func (sr *ServiceRegistry) cleanupRoutine(ctx context.Context) {
	ticker := time.NewTicker(sr.config.HealthInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sr.cleanupStaleServices()
		}
	}
}

// Remove services that haven't been seen recently
func (sr *ServiceRegistry) cleanupStaleServices() {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	now := time.Now()
	for id, service := range sr.services {
		if now.Sub(service.LastSeen) > sr.config.ServiceTTL {
			log.Printf("Removing stale service: %s (%s)", service.Name, id)
			delete(sr.services, id)
		}
	}
}
