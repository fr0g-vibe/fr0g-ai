package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
		// Convert storage.ServiceInfo to main.ServiceInfo
		mainService := &ServiceInfo{
			ID:       service.ID,
			Name:     service.Name,
			Address:  service.Address,
			Port:     service.Port,
			Tags:     service.Tags,
			Meta:     service.Meta,
			Health:   service.Health,
			LastSeen: service.LastSeen,
		}
		r.services[id] = mainService
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
	
	// Convert to storage.ServiceInfo and save to Redis first
	storageService := &storage.ServiceInfo{
		ID:       service.ID,
		Name:     service.Name,
		Address:  service.Address,
		Port:     service.Port,
		Tags:     service.Tags,
		Meta:     service.Meta,
		Health:   service.Health,
		LastSeen: service.LastSeen,
	}
	if err := r.storage.SaveService(r.ctx, storageService); err != nil {
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
				// Convert to storage.ServiceInfo and persist health change to Redis
				storageService := &storage.ServiceInfo{
					ID:       service.ID,
					Name:     service.Name,
					Address:  service.Address,
					Port:     service.Port,
					Tags:     service.Tags,
					Meta:     service.Meta,
					Health:   service.Health,
					LastSeen: service.LastSeen,
				}
				r.storage.SaveService(r.ctx, storageService)
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
	
	// CRITICAL: Add debug logging to verify handler is called
	log.Printf("REGISTER HANDLER CALLED: Method=%s, URL=%s", req.Method, req.URL.Path)
	
	if req.Method != http.MethodPut && req.Method != http.MethodPost {
		log.Printf("REGISTER HANDLER: Method not allowed: %s", req.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}
	
	log.Printf("REGISTER HANDLER: Processing %s request", req.Method)
	
	var service ServiceInfo
	log.Printf("REGISTER HANDLER: Reading request body...")
	
	// Read the entire body first to debug
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("REGISTER HANDLER: Failed to read body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to read request body: " + err.Error()})
		return
	}
	
	log.Printf("REGISTER HANDLER: Request body: %s", string(body))
	
	if len(body) == 0 {
		log.Printf("REGISTER HANDLER: Empty request body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Request body is empty"})
		return
	}
	
	// Parse JSON from the body
	if err := json.Unmarshal(body, &service); err != nil {
		log.Printf("REGISTER HANDLER: JSON decode error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON: " + err.Error()})
		return
	}
	
	log.Printf("REGISTER HANDLER: Decoded service: ID=%s, Name=%s, Address=%s, Port=%d", 
		service.ID, service.Name, service.Address, service.Port)
	
	// Validate required fields
	if service.ID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Service ID is required"})
		return
	}
	
	if service.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Service name is required"})
		return
	}
	
	if service.Address == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Service address is required"})
		return
	}
	
	if service.Port == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Service port is required"})
		return
	}
	
	log.Printf("REGISTER HANDLER: Calling Register method...")
	if err := r.Register(&service); err != nil {
		log.Printf("REGISTER HANDLER: Register failed: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	log.Printf("REGISTER HANDLER: Service registered successfully: %s", service.ID)
	
	// Return success response with service details
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":  "success",
		"message": "Service registered successfully",
		"service": service,
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("REGISTER HANDLER: Failed to encode response: %v", err)
	} else {
		log.Printf("REGISTER HANDLER: Response sent successfully")
	}
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}
	
	services := r.GetServices()
	
	// Ensure we always return a valid JSON object, even if empty
	if services == nil {
		services = make(map[string]*ServiceInfo)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
	
	// Service registration endpoints - CRITICAL: Remove method restrictions to ensure handlers are called
	router.HandleFunc("/v1/agent/service/register", registry.registerHandler)
	router.HandleFunc("/v1/agent/service/deregister/{serviceId}", registry.deregisterHandler)
	router.HandleFunc("/v1/catalog/services", registry.servicesHandler)
	router.HandleFunc("/v1/health/service/{serviceId}", registry.healthServiceHandler)
	
	// Health check endpoint
	router.HandleFunc("/health", registry.healthHandler)
	
	// Prometheus metrics endpoint
	router.Handle("/metrics", promhttp.Handler())
	
	// Start server with optimized settings
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	log.Printf("🚀 Starting optimized fr0g.ai service registry on %s", server.Addr)
	log.Printf("🔧 Redis: %s, Cache: LRU(1000, 30s), Metrics: /metrics", os.Getenv("REDIS_ADDR"))
	log.Printf("📋 Endpoints:")
	log.Printf("   - Service Registration: PUT %s/v1/agent/service/register", server.Addr)
	log.Printf("   - Service Deregistration: PUT %s/v1/agent/service/deregister/{serviceId}", server.Addr)
	log.Printf("   - Service Discovery: GET %s/v1/catalog/services", server.Addr)
	log.Printf("   - Health Check: GET %s/health", server.Addr)
	log.Printf("   - Metrics: GET %s/metrics", server.Addr)
	log.Printf("✅ Service registry ready for automatic lifecycle management")
	
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start registry server:", err)
	}
}
