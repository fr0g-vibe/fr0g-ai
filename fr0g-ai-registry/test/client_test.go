package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// RegistryClient provides a simple client for testing registry interactions
type RegistryClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewRegistryClient creates a new registry client for testing
func NewRegistryClient(baseURL string) *RegistryClient {
	return &RegistryClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// RegisterService registers a service with the registry
func (rc *RegistryClient) RegisterService(service ServiceInfo) error {
	payload, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service: %w", err)
	}

	req, err := http.NewRequest("PUT", rc.baseURL+"/v1/agent/service/register", 
		bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("registration failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeregisterService deregisters a service from the registry
func (rc *RegistryClient) DeregisterService(serviceID string) error {
	url := fmt.Sprintf("%s/v1/agent/service/deregister/%s", rc.baseURL, serviceID)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("deregistration failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// DiscoverServices discovers all registered services
func (rc *RegistryClient) DiscoverServices() (map[string][]string, error) {
	resp, err := rc.httpClient.Get(rc.baseURL + "/v1/catalog/services")
	if err != nil {
		return nil, fmt.Errorf("failed to discover services: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("discovery failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var services map[string][]string
	if err := json.Unmarshal(body, &services); err != nil {
		return nil, fmt.Errorf("failed to unmarshal services: %w", err)
	}

	return services, nil
}

// GetServiceHealth gets the health status of a service
func (rc *RegistryClient) GetServiceHealth(serviceName string) ([]byte, error) {
	url := fmt.Sprintf("%s/v1/health/service/%s", rc.baseURL, serviceName)
	resp, err := rc.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get service health: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}

// TestRegistryClient tests the registry client functionality
func TestRegistryClient(t *testing.T) {
	client := NewRegistryClient("http://localhost:8500")

	testService := ServiceInfo{
		ID:      "test-client-service",
		Name:    "test-client-service",
		Address: "localhost",
		Port:    9999,
		Tags:    []string{"test", "client"},
		Meta: map[string]string{
			"test": "true",
		},
	}

	t.Run("RegisterService", func(t *testing.T) {
		err := client.RegisterService(testService)
		if err != nil {
			t.Fatalf("Failed to register service: %v", err)
		}
		t.Log("Service registered successfully")
	})

	t.Run("DiscoverServices", func(t *testing.T) {
		services, err := client.DiscoverServices()
		if err != nil {
			t.Fatalf("Failed to discover services: %v", err)
		}

		if _, exists := services[testService.Name]; !exists {
			t.Errorf("Registered service not found in discovery")
		}

		t.Logf("Discovered %d services", len(services))
	})

	t.Run("GetServiceHealth", func(t *testing.T) {
		health, err := client.GetServiceHealth(testService.Name)
		if err != nil {
			t.Logf("Health check failed (expected if service not running): %v", err)
		} else {
			t.Logf("Service health: %s", string(health))
		}
	})

	t.Run("DeregisterService", func(t *testing.T) {
		err := client.DeregisterService(testService.ID)
		if err != nil {
			t.Fatalf("Failed to deregister service: %v", err)
		}
		t.Log("Service deregistered successfully")
	})
}
