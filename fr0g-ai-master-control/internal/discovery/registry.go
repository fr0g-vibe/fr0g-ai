package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ServiceInfo represents information about a registered service
type ServiceInfo struct {
	Name      string            `json:"name"`
	ID        string            `json:"id"`
	Address   string            `json:"address"`
	Port      int               `json:"port"`
	Health    string            `json:"health"`
	Metadata  map[string]string `json:"metadata"`
	LastSeen  time.Time         `json:"last_seen"`
	Tags      []string          `json:"tags"`
}

// ServiceRegistry manages service discovery
type ServiceRegistry struct {
	services map[string]*ServiceInfo
	mu       sync.RWMutex
	config   *RegistryConfig
}

// RegistryConfig holds registry configuration
type RegistryConfig struct {
	Port            int           `yaml:"port"`
	Host            string        `yaml:"host"`
	HealthInterval  time.Duration `yaml:"health_interval"`
	ServiceTTL      time.Duration `yaml:"service_ttl"`
	EnableHTTPAPI   bool          `yaml:"enable_http_api"`
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry(config *RegistryConfig) *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]*ServiceInfo),
		config:   config,
	}
}

// Start begins the service registry
func (sr *ServiceRegistry) Start(ctx context.Context) error {
	log.Printf("Service Registry: Starting on %s:%d", sr.config.Host, sr.config.Port)
	
	// Start HTTP API if enabled
	if sr.config.EnableHTTPAPI {
		go sr.startHTTPAPI(ctx)
	}
	
	// Start cleanup routine
	go sr.cleanupLoop(ctx)
	
	return nil
}

// Register registers a service
func (sr *ServiceRegistry) Register(service *ServiceInfo) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	
	service.LastSeen = time.Now()
	sr.services[service.ID] = service
	
	log.Printf("Service Registry: Registered service %s (%s) at %s:%d", 
		service.Name, service.ID, service.Address, service.Port)
	
	return nil
}

// Deregister removes a service
func (sr *ServiceRegistry) Deregister(serviceID string) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	
	if service, exists := sr.services[serviceID]; exists {
		delete(sr.services, serviceID)
		log.Printf("Service Registry: Deregistered service %s (%s)", service.Name, serviceID)
	}
	
	return nil
}

// Discover finds services by name
func (sr *ServiceRegistry) Discover(serviceName string) ([]*ServiceInfo, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	
	var services []*ServiceInfo
	for _, service := range sr.services {
		if service.Name == serviceName && service.Health == "healthy" {
			services = append(services, service)
		}
	}
	
	return services, nil
}

// GetService gets a specific service by ID
func (sr *ServiceRegistry) GetService(serviceID string) (*ServiceInfo, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	
	if service, exists := sr.services[serviceID]; exists {
		return service, nil
	}
	
	return nil, fmt.Errorf("service not found: %s", serviceID)
}

// ListServices returns all registered services
func (sr *ServiceRegistry) ListServices() []*ServiceInfo {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	
	services := make([]*ServiceInfo, 0, len(sr.services))
	for _, service := range sr.services {
		services = append(services, service)
	}
	
	return services
}

// UpdateHealth updates service health status
func (sr *ServiceRegistry) UpdateHealth(serviceID, health string) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	
	if service, exists := sr.services[serviceID]; exists {
		service.Health = health
		service.LastSeen = time.Now()
		return nil
	}
	
	return fmt.Errorf("service not found: %s", serviceID)
}

// startHTTPAPI starts the HTTP API for service discovery
func (sr *ServiceRegistry) startHTTPAPI(ctx context.Context) {
	mux := http.NewServeMux()
	
	// Register endpoints
	mux.HandleFunc("/services", sr.handleServices)
	mux.HandleFunc("/services/", sr.handleServiceByName)
	mux.HandleFunc("/health", sr.handleHealth)
	
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", sr.config.Host, sr.config.Port),
		Handler: mux,
	}
	
	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()
	
	log.Printf("Service Registry: HTTP API listening on %s:%d", sr.config.Host, sr.config.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("Service Registry: HTTP API error: %v", err)
	}
}

// cleanupLoop removes stale services
func (sr *ServiceRegistry) cleanupLoop(ctx context.Context) {
	ticker := time.NewTicker(sr.config.HealthInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sr.cleanup()
		}
	}
}

// cleanup removes stale services
func (sr *ServiceRegistry) cleanup() {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	
	now := time.Now()
	for id, service := range sr.services {
		if now.Sub(service.LastSeen) > sr.config.ServiceTTL {
			delete(sr.services, id)
			log.Printf("Service Registry: Removed stale service %s (%s)", service.Name, id)
		}
	}
}

// HTTP handlers
func (sr *ServiceRegistry) handleServices(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		services := sr.ListServices()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(services)
	case "POST":
		var service ServiceInfo
		if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := sr.Register(&service); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (sr *ServiceRegistry) handleServiceByName(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Path[len("/services/"):]
	services, err := sr.Discover(serviceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (sr *ServiceRegistry) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":        "healthy",
		"service_count": len(sr.services),
		"timestamp":     time.Now(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}
