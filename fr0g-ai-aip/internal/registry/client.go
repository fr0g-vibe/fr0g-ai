package registry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ServiceInfo represents service registration information
type ServiceInfo struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Address string            `json:"address"`
	Port    int               `json:"port"`
	Tags    []string          `json:"tags,omitempty"`
	Meta    map[string]string `json:"meta,omitempty"`
	Check   *HealthCheck      `json:"check,omitempty"`
}

// HealthCheck represents a health check configuration
type HealthCheck struct {
	HTTP     string `json:"http,omitempty"`
	Interval string `json:"interval,omitempty"`
	Timeout  string `json:"timeout,omitempty"`
}

// RegistryClient handles service registry operations
type RegistryClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *logrus.Logger
	
	// Service registration info
	serviceInfo *ServiceInfo
	registered  bool
	mu          sync.RWMutex
	
	// Background health updates
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewRegistryClient creates a new registry client
func NewRegistryClient(registryURL string, logger *logrus.Logger) *RegistryClient {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &RegistryClient{
		baseURL: registryURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}
}

// RegisterService registers this service with the registry
func (rc *RegistryClient) RegisterService(serviceInfo *ServiceInfo) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	
	// Store service info
	rc.serviceInfo = serviceInfo
	
	// Register with registry
	if err := rc.registerWithRegistry(serviceInfo); err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}
	
	rc.registered = true
	
	// Start background health updates
	rc.wg.Add(1)
	go rc.healthUpdateLoop()
	
	rc.logger.WithFields(logrus.Fields{
		"service_id":   serviceInfo.ID,
		"service_name": serviceInfo.Name,
		"address":      serviceInfo.Address,
		"port":         serviceInfo.Port,
	}).Info("Service registered with registry")
	
	return nil
}

// DeregisterService removes this service from the registry
func (rc *RegistryClient) DeregisterService() error {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	
	if !rc.registered || rc.serviceInfo == nil {
		return nil
	}
	
	// Stop background updates
	rc.cancel()
	rc.wg.Wait()
	
	// Deregister from registry
	url := fmt.Sprintf("%s/v1/agent/service/deregister/%s", rc.baseURL, rc.serviceInfo.ID)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create deregister request: %w", err)
	}
	
	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("deregister request failed with status: %d", resp.StatusCode)
	}
	
	rc.registered = false
	
	rc.logger.WithField("service_id", rc.serviceInfo.ID).Info("Service deregistered from registry")
	
	return nil
}

// GetServiceEndpoint discovers an endpoint for a service
func (rc *RegistryClient) GetServiceEndpoint(serviceName string) (string, error) {
	url := fmt.Sprintf("%s/v1/catalog/service/%s", rc.baseURL, serviceName)
	
	resp, err := rc.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to query service: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("service query failed with status: %d", resp.StatusCode)
	}
	
	var services []ServiceInfo
	if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
		return "", fmt.Errorf("failed to decode service response: %w", err)
	}
	
	if len(services) == 0 {
		return "", fmt.Errorf("no healthy instances found for service: %s", serviceName)
	}
	
	// Return the first healthy service
	service := services[0]
	endpoint := fmt.Sprintf("http://%s:%d", service.Address, service.Port)
	
	return endpoint, nil
}

// GetHealthyServices returns all healthy instances of a service
func (rc *RegistryClient) GetHealthyServices(serviceName string) ([]*ServiceInfo, error) {
	url := fmt.Sprintf("%s/v1/health/service/%s?passing=true", rc.baseURL, serviceName)
	
	resp, err := rc.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to query healthy services: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("healthy services query failed with status: %d", resp.StatusCode)
	}
	
	var healthResponse []struct {
		Service *ServiceInfo `json:"Service"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&healthResponse); err != nil {
		return nil, fmt.Errorf("failed to decode healthy services response: %w", err)
	}
	
	services := make([]*ServiceInfo, len(healthResponse))
	for i, item := range healthResponse {
		services[i] = item.Service
	}
	
	return services, nil
}

// registerWithRegistry performs the actual registration
func (rc *RegistryClient) registerWithRegistry(serviceInfo *ServiceInfo) error {
	url := fmt.Sprintf("%s/v1/agent/service/register", rc.baseURL)
	
	payload, err := json.Marshal(serviceInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal service info: %w", err)
	}
	
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create register request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("register request failed with status: %d", resp.StatusCode)
	}
	
	return nil
}

// healthUpdateLoop periodically updates service health
func (rc *RegistryClient) healthUpdateLoop() {
	defer rc.wg.Done()
	
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-rc.ctx.Done():
			return
		case <-ticker.C:
			rc.updateHealth()
		}
	}
}

// updateHealth updates the service health status
func (rc *RegistryClient) updateHealth() {
	rc.mu.RLock()
	serviceInfo := rc.serviceInfo
	rc.mu.RUnlock()
	
	if serviceInfo == nil {
		return
	}
	
	// Re-register to update health status
	if err := rc.registerWithRegistry(serviceInfo); err != nil {
		rc.logger.WithError(err).Warn("Failed to update service health")
	}
}

// Shutdown gracefully shuts down the registry client
func (rc *RegistryClient) Shutdown() error {
	return rc.DeregisterService()
}

// NewRegistryClient creates a new registry client with URL
func NewRegistryClient(registryURL string) (*RegistryClient, error) {
	return NewRegistryClient(registryURL, nil), nil
}

// RegisterService registers a service with the registry
func (rc *RegistryClient) RegisterService(serviceName, httpPort, grpcPort string) error {
	serviceInfo := &ServiceInfo{
		ID:      serviceName + "-" + httpPort,
		Name:    serviceName,
		Address: "localhost",
		Port:    8080, // Default HTTP port
		Tags:    []string{"http", "grpc"},
		Meta: map[string]string{
			"http_port": httpPort,
			"grpc_port": grpcPort,
		},
	}
	
	return rc.RegisterService(serviceInfo)
}
