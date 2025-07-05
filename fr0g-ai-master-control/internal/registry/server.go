package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RegistryConfig holds configuration for the service registry
type RegistryConfig struct {
	Port           int
	Host           string
	HealthInterval time.Duration
	ServiceTTL     time.Duration
	EnableHTTPAPI  bool
}

// ServiceRegistry wraps the Server to match the expected API
type ServiceRegistry struct {
	server *Server
	config *RegistryConfig
}

// NewServiceRegistry creates a new service registry with the expected API
func NewServiceRegistry(config *RegistryConfig) *ServiceRegistry {
	return &ServiceRegistry{
		server: NewServer(),
		config: config,
	}
}

// Start starts the service registry server
func (sr *ServiceRegistry) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", sr.config.Host, sr.config.Port)
	return sr.server.Start(ctx, addr)
}

// ServiceInfo represents a registered service
type ServiceInfo struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Tags     []string          `json:"tags"`
	Meta     map[string]string `json:"meta"`
	Health   string            `json:"health"` // "passing", "warning", "critical"
	LastSeen time.Time         `json:"last_seen"`
}

// HealthStatus represents service health
type HealthStatus struct {
	Status      string    `json:"status"` // "passing", "warning", "critical"
	Output      string    `json:"output"`
	LastChecked time.Time `json:"last_checked"`
}

// Registry manages service registration and discovery
type Registry struct {
	services map[string]*ServiceInfo
	mu       sync.RWMutex

	// Health checking
	healthInterval time.Duration
	healthTimeout  time.Duration
	stopChan       chan struct{}
}

// NewRegistry creates a new service registry
func NewRegistry() *Registry {
	return &Registry{
		services:       make(map[string]*ServiceInfo),
		healthInterval: 30 * time.Second,
		healthTimeout:  10 * time.Second,
		stopChan:       make(chan struct{}),
	}
}

// RegisterService registers a new service
func (r *Registry) RegisterService(service *ServiceInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	service.LastSeen = time.Now()
	service.Health = "unknown"

	r.services[service.ID] = service
	log.Printf("Service registered: %s (%s)", service.Name, service.ID)
	return nil
}

// DeregisterService removes a service from registry
func (r *Registry) DeregisterService(serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if service, exists := r.services[serviceID]; exists {
		delete(r.services, serviceID)
		log.Printf("Service deregistered: %s (%s)", service.Name, serviceID)
		return nil
	}

	return fmt.Errorf("service not found: %s", serviceID)
}

// GetService retrieves a service by ID
func (r *Registry) GetService(serviceID string) (*ServiceInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if service, exists := r.services[serviceID]; exists {
		return service, nil
	}

	return nil, fmt.Errorf("service not found: %s", serviceID)
}

// GetServicesByName retrieves all services with a given name
func (r *Registry) GetServicesByName(serviceName string) []*ServiceInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var services []*ServiceInfo
	for _, service := range r.services {
		if service.Name == serviceName {
			services = append(services, service)
		}
	}

	return services
}

// GetAllServices returns all registered services
func (r *Registry) GetAllServices() []*ServiceInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var services []*ServiceInfo
	for _, service := range r.services {
		services = append(services, service)
	}

	return services
}

// StartHealthChecking begins health checking for all services
func (r *Registry) StartHealthChecking(ctx context.Context) {
	ticker := time.NewTicker(r.healthInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-r.stopChan:
			return
		case <-ticker.C:
			r.performHealthChecks()
		}
	}
}

// performHealthChecks checks health of all registered services
func (r *Registry) performHealthChecks() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, service := range r.services {
		go r.checkServiceHealth(service)
	}
}

// checkServiceHealth performs health check for a single service
func (r *Registry) checkServiceHealth(service *ServiceInfo) {
	healthURL := fmt.Sprintf("http://%s:%d/health", service.Address, service.Port)

	client := &http.Client{Timeout: r.healthTimeout}
	resp, err := client.Get(healthURL)

	r.mu.Lock()
	defer r.mu.Unlock()

	if err != nil {
		service.Health = "critical"
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		service.Health = "passing"
		service.LastSeen = time.Now()
	} else {
		service.Health = "warning"
	}
}

// Stop stops the registry
func (r *Registry) Stop() {
	close(r.stopChan)
}

// Server provides HTTP API for the service registry
type Server struct {
	registry *Registry
	mux      *http.ServeMux
}

// NewServer creates a new registry server
func NewServer() *Server {
	registry := NewRegistry()
	server := &Server{
		registry: registry,
		mux:      http.NewServeMux(),
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures HTTP routes
func (s *Server) setupRoutes() {
	// Service registration
	s.mux.HandleFunc("/v1/agent/service/register", s.registerServiceHandler)
	s.mux.HandleFunc("/v1/agent/service/deregister/", s.deregisterServiceHandler)

	// Service discovery
	s.mux.HandleFunc("/v1/catalog/services", s.listServicesHandler)
	s.mux.HandleFunc("/v1/catalog/service/", s.getServiceHandler)
	s.mux.HandleFunc("/v1/health/service/", s.getServiceHealthHandler)

	// Health endpoint
	s.mux.HandleFunc("/health", s.healthHandler)
}

// HTTP Handlers
func (s *Server) registerServiceHandler(w http.ResponseWriter, r *http.Request) {
	var service ServiceInfo
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := s.registry.RegisterService(&service); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) deregisterServiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract service ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/v1/agent/service/deregister/")
	serviceID := strings.TrimSuffix(path, "/")

	if serviceID == "" {
		http.Error(w, "Service ID required", http.StatusBadRequest)
		return
	}

	if err := s.registry.DeregisterService(serviceID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) listServicesHandler(w http.ResponseWriter, r *http.Request) {
	services := s.registry.GetAllServices()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (s *Server) getServiceHandler(w http.ResponseWriter, r *http.Request) {
	// Extract service name from URL path
	path := strings.TrimPrefix(r.URL.Path, "/v1/catalog/service/")
	serviceName := strings.TrimSuffix(path, "/")

	if serviceName == "" {
		http.Error(w, "Service name required", http.StatusBadRequest)
		return
	}

	services := s.registry.GetServicesByName(serviceName)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (s *Server) getServiceHealthHandler(w http.ResponseWriter, r *http.Request) {
	// Extract service ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/v1/health/service/")
	serviceID := strings.TrimSuffix(path, "/")

	if serviceID == "" {
		http.Error(w, "Service ID required", http.StatusBadRequest)
		return
	}

	service, err := s.registry.GetService(serviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(service.Health)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":              "healthy",
		"timestamp":           time.Now(),
		"registered_services": len(s.registry.services),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// Start starts the registry server
func (s *Server) Start(ctx context.Context, addr string) error {
	// Start health checking
	go s.registry.StartHealthChecking(ctx)

	server := &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}

	log.Printf("Service registry starting on %s", addr)
	return server.ListenAndServe()
}
