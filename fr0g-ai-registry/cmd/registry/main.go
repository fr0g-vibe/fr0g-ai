package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-registry/internal/cache"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-registry/internal/metrics"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-registry/internal/storage"
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

// Registry manages service registration with Redis persistence and caching
type Registry struct {
	services map[string]*ServiceInfo
	mu       sync.RWMutex
	storage  *storage.RedisStorage
	cache    *cache.LRUCache
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewRegistry creates a new service registry with Redis and caching
func NewRegistry() *Registry {
	ctx, cancel := context.WithCancel(context.Background())
	
	// Initialize Redis storage
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if db, err := strconv.Atoi(dbStr); err == nil {
			redisDB = db
		}
	}
	
	redisStorage := storage.NewRedisStorage(redisAddr, redisPassword, redisDB)
	
	// Initialize LRU cache (1000 entries, 30 second TTL)
	lruCache, err := cache.NewLRUCache(1000, 30*time.Second)
	if err != nil {
		log.Fatal("Failed to create LRU cache:", err)
	}
	
	registry := &Registry{
		services: make(map[string]*ServiceInfo),
		storage:  redisStorage,
		cache:    lruCache,
		ctx:      ctx,
		cancel:   cancel,
	}
	
	// Load services from Redis on startup
	if err := registry.loadFromRedis(); err != nil {
		log.Printf("Warning: Failed to load services from Redis: %v", err)
	}
	
	// Start health monitoring
	go registry.healthMonitor()
	
	return registry
}

// loadFromRedis loads all services from Redis into memory
func (r *Registry) loadFromRedis() error {
	services, err := r.storage.GetAllServices(r.ctx)
	if err != nil {
		return err
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	for id, service := range services {
		r.services[id] = service
		log.Printf("Loaded service from Redis: %s (%s)", service.Name, service.ID)
	}
	
	log.Printf("Loaded %d services from Redis", len(services))
	return nil
}

// Register adds a service to the registry with Redis persistence
func (r *Registry) Register(service *ServiceInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	service.LastSeen = time.Now()
	service.Health = "passing"
	
	// Save to Redis first
	if err := r.storage.SaveService(r.ctx, service); err != nil {
		metrics.RedisOperations.WithLabelValues("save", "error").Inc()
		return fmt.Errorf("failed to persist service: %w", err)
	}
	metrics.RedisOperations.WithLabelValues("save", "success").Inc()
	
	// Update in-memory registry
	r.services[service.ID] = service
	
	// Invalidate cache entries
	r.cache.Delete("all_services")
	r.cache.Delete("service:" + service.ID)
	
	// Update metrics
	metrics.ServiceRegistrations.WithLabelValues(service.Name).Inc()
	r.updateServiceMetrics()
	
	log.Printf("Registered service: %s (%s)", service.Name, service.ID)
	return nil
}

// Deregister removes a service from the registry with Redis cleanup
func (r *Registry) Deregister(serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	service, exists := r.services[serviceID]
	if !exists {
		return fmt.Errorf("service not found: %s", serviceID)
	}
	
	// Remove from Redis
	if err := r.storage.DeleteService(r.ctx, serviceID); err != nil {
		metrics.RedisOperations.WithLabelValues("delete", "error").Inc()
		return fmt.Errorf("failed to remove service from storage: %w", err)
	}
	metrics.RedisOperations.WithLabelValues("delete", "success").Inc()
	
	// Remove from in-memory registry
	delete(r.services, serviceID)
	
	// Invalidate cache entries
	r.cache.Delete("all_services")
	r.cache.Delete("service:" + serviceID)
	
	// Update metrics
	metrics.ServiceDeregistrations.WithLabelValues(service.Name).Inc()
	r.updateServiceMetrics()
	
	log.Printf("Deregistered service: %s (%s)", service.Name, serviceID)
	return nil
}

// GetServices returns all registered services with caching
func (r *Registry) GetServices() map[string]*ServiceInfo {
	// Try cache first
	if cached, found := r.cache.Get("all_services"); found {
		metrics.CacheHits.Inc()
		return cached.(map[string]*ServiceInfo)
	}
	metrics.CacheMisses.Inc()
	
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	services := make(map[string]*ServiceInfo)
	for id, service := range r.services {
		services[id] = service
	}
	
	// Cache the result
	r.cache.Set("all_services", services)
	
	return services
}

// GetService returns a specific service by ID with caching
func (r *Registry) GetService(serviceID string) (*ServiceInfo, bool) {
	cacheKey := "service:" + serviceID
	
	// Try cache first
	if cached, found := r.cache.Get(cacheKey); found {
		metrics.CacheHits.Inc()
		if cached == nil {
			return nil, false
		}
		return cached.(*ServiceInfo), true
	}
	metrics.CacheMisses.Inc()
	
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	service, exists := r.services[serviceID]
	
	// Cache the result (including nil for not found)
	if exists {
		r.cache.Set(cacheKey, service)
	} else {
		r.cache.Set(cacheKey, nil)
	}
	
	return service, exists
}

// updateServiceMetrics updates Prometheus metrics for service counts
func (r *Registry) updateServiceMetrics() {
	healthCounts := make(map[string]int)
	for _, service := range r.services {
		healthCounts[service.Health]++
	}
	
	for status, count := range healthCounts {
		metrics.ServicesTotal.WithLabelValues(status).Set(float64(count))
	}
}

// healthMonitor periodically checks service health
func (r *Registry) healthMonitor() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ticker.C:
			r.checkServiceHealth()
		}
	}
}

// checkServiceHealth updates service health status
func (r *Registry) checkServiceHealth() {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	now := time.Now()
	for _, service := range r.services {
		// Mark services as unhealthy if not seen for 2 minutes
		if now.Sub(service.LastSeen) > 2*time.Minute {
			if service.Health != "critical" {
				service.Health = "critical"
				// Persist health change to Redis
				r.storage.SaveService(r.ctx, service)
				// Invalidate cache
				r.cache.Delete("all_services")
				r.cache.Delete("service:" + service.ID)
			}
		}
	}
	
	r.updateServiceMetrics()
}

// HTTP Handlers with performance optimizations
func (r *Registry) registerHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	defer func() {
		metrics.DiscoveryLatency.WithLabelValues("register").Observe(time.Since(start).Seconds())
		metrics.DiscoveryRequests.WithLabelValues("register").Inc()
	}()
	
	if req.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var service ServiceInfo
	if err := json.NewDecoder(req.Body).Decode(&service); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	if err := r.Register(&service); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

func (r *Registry) deregisterHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	defer func() {
		metrics.DiscoveryLatency.WithLabelValues("deregister").Observe(time.Since(start).Seconds())
		metrics.DiscoveryRequests.WithLabelValues("deregister").Inc()
	}()
	
	if req.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	vars := mux.Vars(req)
	serviceID := vars["serviceId"]
	
	if err := r.Deregister(serviceID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

func (r *Registry) servicesHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	defer func() {
		metrics.DiscoveryLatency.WithLabelValues("services").Observe(time.Since(start).Seconds())
		metrics.DiscoveryRequests.WithLabelValues("services").Inc()
	}()
	
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	services := r.GetServices()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func (r *Registry) healthServiceHandler(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	defer func() {
		metrics.DiscoveryLatency.WithLabelValues("health").Observe(time.Since(start).Seconds())
		metrics.DiscoveryRequests.WithLabelValues("health").Inc()
	}()
	
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	vars := mux.Vars(req)
	serviceID := vars["serviceId"]
	
	service, exists := r.GetService(serviceID)
	if !exists {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"service": service,
		"checks": []map[string]string{
			{
				"status": service.Health,
				"output": fmt.Sprintf("Service %s is %s", service.Name, service.Health),
			},
		},
	})
}

func (r *Registry) healthHandler(w http.ResponseWriter, req *http.Request) {
	// Check Redis connectivity
	redisHealthy := true
	if err := r.storage.Ping(r.ctx); err != nil {
		redisHealthy = false
	}
	
	status := "healthy"
	if !redisHealthy {
		status = "degraded"
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      status,
		"service":     "fr0g-ai-registry",
		"timestamp":   time.Now().Format(time.RFC3339),
		"redis":       redisHealthy,
		"services":    len(r.services),
	})
}

// Shutdown gracefully shuts down the registry
func (r *Registry) Shutdown() {
	log.Println("Shutting down registry...")
	r.cancel()
	if r.storage != nil {
		r.storage.Close()
	}
}

func main() {
	registry := NewRegistry()
	
	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		registry.Shutdown()
		os.Exit(0)
	}()
	
	// Get port from environment or use default
	port := os.Getenv("REGISTRY_PORT")
	if port == "" {
		port = "8500"
	}
	
	host := os.Getenv("REGISTRY_HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	
	// Setup routes with optimized handlers
	router := mux.NewRouter()
	
	// Service registration endpoints
	router.HandleFunc("/v1/agent/service/register", registry.registerHandler).Methods("PUT")
	router.HandleFunc("/v1/agent/service/deregister/{serviceId}", registry.deregisterHandler).Methods("PUT")
	router.HandleFunc("/v1/catalog/services", registry.servicesHandler).Methods("GET")
	router.HandleFunc("/v1/health/service/{serviceId}", registry.healthServiceHandler).Methods("GET")
	
	// Health check endpoint
	router.HandleFunc("/health", registry.healthHandler).Methods("GET")
	
	// Prometheus metrics endpoint
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")
	
	// Start server with optimized settings
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	log.Printf("Starting optimized fr0g.ai service registry on %s", server.Addr)
	log.Printf("Redis: %s, Cache: LRU(1000, 30s), Metrics: /metrics", os.Getenv("REDIS_ADDR"))
	
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start registry server:", err)
	}
}
