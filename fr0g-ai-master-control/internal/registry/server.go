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
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Address     string            `json:"address"`
	Port        int               `json:"port"`
	Tags        []string          `json:"tags"`
	Meta        map[string]string `json:"meta"`
	Health      HealthStatus      `json:"health"`
	RegisteredAt time.Time        `json:"registered_at"`
	LastSeen    time.Time         `json:"last_seen"`
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
	
	service.RegisteredAt = time.Now()
	service.LastSeen = time.Now()
	service.Health = HealthStatus{
		Status:      "unknown",
		LastChecked: time.Now(),
	}
	
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
		service.Health = HealthStatus{
			Status:      "critical",
			Output:      fmt.Sprintf("Health check failed: %v", err),
			LastChecked: time.Now(),
		}
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		service.Health = HealthStatus{
			Status:      "passing",
			Output:      "Health check passed",
			LastChecked: time.Now(),
		}
		service.LastSeen = time.Now()
	} else {
		service.Health = HealthStatus{
			Status:      "warning",
			Output:      fmt.Sprintf("Health check returned status %d", resp.StatusCode),
			LastChecked: time.Now(),
		}
	}
}

// Stop stops the registry
func (r *Registry) Stop() {
	close(r.stopChan)
}

// Server provides HTTP API for the service registry
type Server struct {
	registry *Registry
	router   *mux.Router
}

// NewServer creates a new registry server
func NewServer() *Server {
	registry := NewRegistry()
	server := &Server{
		registry: registry,
		router:   mux.NewRouter(),
	}
	
	server.setupRoutes()
	return server
}

// setupRoutes configures HTTP routes
func (s *Server) setupRoutes() {
	api := s.router.PathPrefix("/v1").Subrouter()
	
	// Service registration
	api.HandleFunc("/agent/service/register", s.registerServiceHandler).Methods("PUT")
	api.HandleFunc("/agent/service/deregister/{serviceId}", s.deregisterServiceHandler).Methods("PUT")
	
	// Service discovery
	api.HandleFunc("/catalog/services", s.listServicesHandler).Methods("GET")
	api.HandleFunc("/catalog/service/{serviceName}", s.getServiceHandler).Methods("GET")
	api.HandleFunc("/health/service/{serviceId}", s.getServiceHealthHandler).Methods("GET")
	
	// Health endpoint
	s.router.HandleFunc("/health", s.healthHandler).Methods("GET")
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
	vars := mux.Vars(r)
	serviceID := vars["serviceId"]
	
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
	vars := mux.Vars(r)
	serviceName := vars["serviceName"]
	
	services := s.registry.GetServicesByName(serviceName)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (s *Server) getServiceHealthHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID := vars["serviceId"]
	
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
		"status":           "healthy",
		"timestamp":        time.Now(),
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
		Handler: s.router,
	}
	
	log.Printf("Service registry starting on %s", addr)
	return server.ListenAndServe()
}
