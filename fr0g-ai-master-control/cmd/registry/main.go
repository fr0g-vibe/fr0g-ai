package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	Health   string            `json:"health"`
	LastSeen time.Time         `json:"last_seen"`
}

// Registry manages service registration
type Registry struct {
	services map[string]*ServiceInfo
	mu       sync.RWMutex
}

// NewRegistry creates a new service registry
func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]*ServiceInfo),
	}
}

// Register adds a service to the registry
func (r *Registry) Register(service *ServiceInfo) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	service.LastSeen = time.Now()
	service.Health = "passing"
	r.services[service.ID] = service
	
	log.Printf("Registered service: %s (%s)", service.Name, service.ID)
}

// Deregister removes a service from the registry
func (r *Registry) Deregister(serviceID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if service, exists := r.services[serviceID]; exists {
		delete(r.services, serviceID)
		log.Printf("Deregistered service: %s (%s)", service.Name, serviceID)
	}
}

// GetServices returns all registered services
func (r *Registry) GetServices() map[string]*ServiceInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	services := make(map[string]*ServiceInfo)
	for id, service := range r.services {
		services[id] = service
	}
	return services
}

// GetService returns a specific service by ID
func (r *Registry) GetService(serviceID string) (*ServiceInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	service, exists := r.services[serviceID]
	return service, exists
}

// HTTP Handlers
func (r *Registry) registerHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var service ServiceInfo
	if err := json.NewDecoder(req.Body).Decode(&service); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	r.Register(&service)
	w.WriteHeader(http.StatusOK)
}

func (r *Registry) deregisterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	vars := mux.Vars(req)
	serviceID := vars["serviceId"]
	
	r.Deregister(serviceID)
	w.WriteHeader(http.StatusOK)
}

func (r *Registry) servicesHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	services := r.GetServices()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (r *Registry) healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "fr0g-ai-registry",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func main() {
	registry := NewRegistry()
	
	// Get port from environment or use default
	port := os.Getenv("REGISTRY_PORT")
	if port == "" {
		port = "8500"
	}
	
	host := os.Getenv("REGISTRY_HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	
	// Setup routes
	router := mux.NewRouter()
	
	// Service registration endpoints
	router.HandleFunc("/v1/agent/service/register", registry.registerHandler).Methods("PUT")
	router.HandleFunc("/v1/agent/service/deregister/{serviceId}", registry.deregisterHandler).Methods("PUT")
	router.HandleFunc("/v1/catalog/services", registry.servicesHandler).Methods("GET")
	router.HandleFunc("/v1/health/service/{serviceId}", registry.servicesHandler).Methods("GET")
	
	// Health check endpoint
	router.HandleFunc("/health", registry.healthHandler).Methods("GET")
	
	// Start server
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("🔍 Starting fr0g.ai service registry on %s", addr)
	
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("Failed to start registry server:", err)
	}
}
