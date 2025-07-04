package discovery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Client provides service discovery client functionality
type Client struct {
	registryURL string
	serviceInfo *ServiceInfo
	httpClient  *http.Client
	ctx         context.Context
	cancel      context.CancelFunc
}

// ClientConfig holds client configuration
type ClientConfig struct {
	RegistryURL     string        `yaml:"registry_url"`
	ServiceName     string        `yaml:"service_name"`
	ServiceID       string        `yaml:"service_id"`
	ServiceAddress  string        `yaml:"service_address"`
	ServicePort     int           `yaml:"service_port"`
	HealthInterval  time.Duration `yaml:"health_interval"`
	Tags            []string      `yaml:"tags"`
	Metadata        map[string]string `yaml:"metadata"`
}

// NewClient creates a new service discovery client
func NewClient(config *ClientConfig) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	
	serviceInfo := &ServiceInfo{
		Name:     config.ServiceName,
		ID:       config.ServiceID,
		Address:  config.ServiceAddress,
		Port:     config.ServicePort,
		Health:   "healthy",
		Tags:     config.Tags,
		Metadata: config.Metadata,
		LastSeen: time.Now(),
	}
	
	if serviceInfo.Metadata == nil {
		serviceInfo.Metadata = make(map[string]string)
	}
	
	return &Client{
		registryURL: config.RegistryURL,
		serviceInfo: serviceInfo,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start begins service discovery client operations
func (c *Client) Start() error {
	log.Printf("Service Discovery Client: Starting for service %s (%s)", 
		c.serviceInfo.Name, c.serviceInfo.ID)
	
	// Register service
	if err := c.register(); err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}
	
	// Start health reporting
	go c.healthReportLoop()
	
	return nil
}

// Stop gracefully stops the client
func (c *Client) Stop() error {
	log.Printf("Service Discovery Client: Stopping service %s (%s)", 
		c.serviceInfo.Name, c.serviceInfo.ID)
	
	c.cancel()
	
	// Deregister service
	return c.deregister()
}

// Discover finds services by name
func (c *Client) Discover(serviceName string) ([]*ServiceInfo, error) {
	url := fmt.Sprintf("%s/services/%s", c.registryURL, serviceName)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to discover services: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("discovery request failed with status: %d", resp.StatusCode)
	}
	
	var services []*ServiceInfo
	if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
		return nil, fmt.Errorf("failed to decode services: %w", err)
	}
	
	return services, nil
}

// GetServiceURL returns the URL for a service (load balanced)
func (c *Client) GetServiceURL(serviceName string) (string, error) {
	services, err := c.Discover(serviceName)
	if err != nil {
		return "", err
	}
	
	if len(services) == 0 {
		return "", fmt.Errorf("no healthy instances of service %s found", serviceName)
	}
	
	// Simple round-robin (use first available for now)
	service := services[0]
	return fmt.Sprintf("http://%s:%d", service.Address, service.Port), nil
}

// register registers the service with the registry
func (c *Client) register() error {
	url := fmt.Sprintf("%s/services", c.registryURL)
	
	data, err := json.Marshal(c.serviceInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal service info: %w", err)
	}
	
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("registration failed with status: %d", resp.StatusCode)
	}
	
	log.Printf("Service Discovery Client: Successfully registered service %s", c.serviceInfo.Name)
	return nil
}

// deregister removes the service from the registry
func (c *Client) deregister() error {
	url := fmt.Sprintf("%s/services/%s", c.registryURL, c.serviceInfo.ID)
	
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create deregister request: %w", err)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}
	defer resp.Body.Close()
	
	log.Printf("Service Discovery Client: Deregistered service %s", c.serviceInfo.Name)
	return nil
}

// healthReportLoop periodically reports health status
func (c *Client) healthReportLoop() {
	ticker := time.NewTicker(30 * time.Second) // Report health every 30 seconds
	defer ticker.Stop()
	
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if err := c.reportHealth(); err != nil {
				log.Printf("Service Discovery Client: Failed to report health: %v", err)
			}
		}
	}
}

// reportHealth reports current health status
func (c *Client) reportHealth() error {
	// Re-register to update LastSeen timestamp
	return c.register()
}

// UpdateMetadata updates service metadata
func (c *Client) UpdateMetadata(metadata map[string]string) error {
	for k, v := range metadata {
		c.serviceInfo.Metadata[k] = v
	}
	
	return c.register() // Re-register with updated metadata
}
