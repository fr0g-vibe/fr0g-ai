package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// CrossServiceTest verifies all fr0g.ai services can register and discover each other
func TestCrossServiceDiscovery(t *testing.T) {
	registryURL := "http://localhost:8500"
	client := &http.Client{Timeout: 10 * time.Second}

	// Define all fr0g.ai services with their expected configurations
	allServices := []ServiceInfo{
		{
			ID:      "fr0g-ai-aip-prod",
			Name:    "fr0g-ai-aip",
			Address: "localhost",
			Port:    8080,
			Tags:    []string{"ai", "persona", "identity", "core"},
			Meta: map[string]string{
				"version":     "1.0.0",
				"environment": "production",
				"role":        "core-ai-processing",
				"grpc_port":   "9090",
			},
			Check: &HealthCheck{
				HTTP:     "http://localhost:8080/health",
				Interval: "10s",
				Timeout:  "3s",
			},
		},
		{
			ID:      "fr0g-ai-bridge-prod",
			Name:    "fr0g-ai-bridge",
			Address: "localhost",
			Port:    8082,
			Tags:    []string{"integration", "api-gateway", "openwebui"},
			Meta: map[string]string{
				"version":     "1.0.0",
				"environment": "production",
				"role":        "integration-bridge",
				"grpc_port":   "9091",
			},
			Check: &HealthCheck{
				HTTP:     "http://localhost:8082/health",
				Interval: "10s",
				Timeout:  "3s",
			},
		},
		{
			ID:      "fr0g-ai-master-control-prod",
			Name:    "fr0g-ai-master-control",
			Address: "localhost",
			Port:    8081,
			Tags:    []string{"cognitive", "orchestration", "intelligence"},
			Meta: map[string]string{
				"version":     "1.0.0",
				"environment": "production",
				"role":        "cognitive-engine",
				"ai_status":   "conscious",
			},
			Check: &HealthCheck{
				HTTP:     "http://localhost:8081/health",
				Interval: "10s",
				Timeout:  "3s",
			},
		},
		{
			ID:      "fr0g-ai-io-prod",
			Name:    "fr0g-ai-io",
			Address: "localhost",
			Port:    8083,
			Tags:    []string{"input", "output", "threat-detection"},
			Meta: map[string]string{
				"version":     "1.0.0",
				"environment": "production",
				"role":        "io-processing",
				"grpc_port":   "9092",
			},
			Check: &HealthCheck{
				HTTP:     "http://localhost:8083/health",
				Interval: "10s",
				Timeout:  "3s",
			},
		},
		{
			ID:      "fr0g-ai-registry-prod",
			Name:    "fr0g-ai-registry",
			Address: "localhost",
			Port:    8500,
			Tags:    []string{"service-discovery", "health-monitoring", "registry"},
			Meta: map[string]string{
				"version":     "1.0.0",
				"environment": "production",
				"role":        "service-registry",
			},
			Check: &HealthCheck{
				HTTP:     "http://localhost:8500/health",
				Interval: "10s",
				Timeout:  "3s",
			},
		},
	}

	t.Run("RegisterAllServices", func(t *testing.T) {
		for _, service := range allServices {
			t.Run(service.Name, func(t *testing.T) {
				if err := registerService(client, registryURL, service); err != nil {
					t.Fatalf("Failed to register %s: %v", service.Name, err)
				}
				t.Logf("Successfully registered %s on port %d", service.Name, service.Port)
			})
		}
	})

	t.Run("VerifyAllServicesDiscoverable", func(t *testing.T) {
		// Wait a moment for registration to propagate
		time.Sleep(100 * time.Millisecond)

		services, err := discoverAllServices(client, registryURL)
		if err != nil {
			t.Fatalf("Failed to discover services: %v", err)
		}

		// Verify each expected service is discoverable
		for _, expected := range allServices {
			if serviceDetail, found := findServiceByName(services, expected.Name); found {
				t.Logf("✓ Found %s: %s:%d (tags: %v)", 
					expected.Name, serviceDetail.Address, serviceDetail.Port, serviceDetail.Tags)
				
				// Verify service details match expectations
				if serviceDetail.Port != expected.Port {
					t.Errorf("Port mismatch for %s: expected %d, got %d", 
						expected.Name, expected.Port, serviceDetail.Port)
				}
				
				// Verify required tags are present
				if !containsAllTags(serviceDetail.Tags, expected.Tags) {
					t.Errorf("Missing tags for %s: expected %v, got %v", 
						expected.Name, expected.Tags, serviceDetail.Tags)
				}
			} else {
				t.Errorf("✗ Service %s not found in discovery", expected.Name)
			}
		}

		t.Logf("Service discovery verification complete. Found %d services", len(services))
	})

	t.Run("CrossServiceCommunicationTest", func(t *testing.T) {
		// Test that each service can discover others through the registry
		testCases := []struct {
			requester string
			target    string
			purpose   string
		}{
			{"fr0g-ai-bridge", "fr0g-ai-aip", "AI processing requests"},
			{"fr0g-ai-master-control", "fr0g-ai-aip", "persona analysis"},
			{"fr0g-ai-master-control", "fr0g-ai-io", "threat processing"},
			{"fr0g-ai-io", "fr0g-ai-master-control", "event reporting"},
			{"fr0g-ai-bridge", "fr0g-ai-master-control", "orchestration requests"},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("%s_to_%s", tc.requester, tc.target), func(t *testing.T) {
				// Simulate service discovery from requester's perspective
				targetService, err := discoverSpecificService(client, registryURL, tc.target)
				if err != nil {
					t.Fatalf("Service %s failed to discover %s: %v", tc.requester, tc.target, err)
				}

				t.Logf("✓ %s successfully discovered %s at %s:%d for %s", 
					tc.requester, tc.target, targetService.Address, targetService.Port, tc.purpose)

				// Verify the discovered service has expected metadata
				if role, exists := targetService.Meta["role"]; exists {
					t.Logf("  Target service role: %s", role)
				}
			})
		}
	})

	t.Run("ServiceHealthMonitoring", func(t *testing.T) {
		for _, service := range allServices {
			t.Run(service.Name+"_health", func(t *testing.T) {
				health, err := getServiceHealth(client, registryURL, service.Name)
				if err != nil {
					t.Logf("Health check for %s returned error (expected if service not running): %v", 
						service.Name, err)
					return
				}

				// Parse health response to verify service status
				var healthData map[string]ServiceDetail
				if err := json.Unmarshal(health, &healthData); err != nil {
					t.Errorf("Failed to parse health data for %s: %v", service.Name, err)
					return
				}

				// Find our service in the health data
				for _, serviceDetail := range healthData {
					if serviceDetail.Name == service.Name {
						t.Logf("✓ %s health status: %s (last seen: %s)", 
							service.Name, serviceDetail.Health, serviceDetail.LastSeen)
						
						if serviceDetail.Health != "passing" {
							t.Logf("⚠ Service %s health status is %s", service.Name, serviceDetail.Health)
						}
						break
					}
				}
			})
		}
	})

	t.Run("ServiceMetadataValidation", func(t *testing.T) {
		services, err := discoverAllServices(client, registryURL)
		if err != nil {
			t.Fatalf("Failed to discover services for metadata validation: %v", err)
		}

		// Validate that each service has required metadata
		requiredMetadata := map[string][]string{
			"fr0g-ai-aip":            {"role", "version", "grpc_port"},
			"fr0g-ai-bridge":         {"role", "version", "grpc_port"},
			"fr0g-ai-master-control": {"role", "version", "ai_status"},
			"fr0g-ai-io":             {"role", "version", "grpc_port"},
			"fr0g-ai-registry":       {"role", "version"},
		}

		for serviceName, requiredKeys := range requiredMetadata {
			if serviceDetail, found := findServiceByName(services, serviceName); found {
				for _, key := range requiredKeys {
					if value, exists := serviceDetail.Meta[key]; exists {
						t.Logf("✓ %s has %s: %s", serviceName, key, value)
					} else {
						t.Errorf("✗ %s missing required metadata: %s", serviceName, key)
					}
				}
			}
		}
	})

	t.Run("CleanupTestServices", func(t *testing.T) {
		for _, service := range allServices {
			if err := deregisterService(client, registryURL, service.ID); err != nil {
				t.Logf("Warning: Failed to deregister %s: %v", service.Name, err)
			} else {
				t.Logf("✓ Deregistered %s", service.Name)
			}
		}
	})
}

// Helper functions for cross-service testing

func registerService(client *http.Client, registryURL string, service ServiceInfo) error {
	payload, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service: %w", err)
	}

	req, err := http.NewRequest("PUT", registryURL+"/v1/agent/service/register", 
		bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
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

func discoverAllServices(client *http.Client, registryURL string) (map[string]ServiceDetail, error) {
	resp, err := client.Get(registryURL + "/v1/catalog/services")
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

	var services map[string]ServiceDetail
	if err := json.Unmarshal(body, &services); err != nil {
		return nil, fmt.Errorf("failed to unmarshal services: %w", err)
	}

	return services, nil
}

func discoverSpecificService(client *http.Client, registryURL, serviceName string) (*ServiceDetail, error) {
	services, err := discoverAllServices(client, registryURL)
	if err != nil {
		return nil, err
	}

	if service, found := findServiceByName(services, serviceName); found {
		return &service, nil
	}

	return nil, fmt.Errorf("service %s not found", serviceName)
}

func getServiceHealth(client *http.Client, registryURL, serviceName string) ([]byte, error) {
	url := fmt.Sprintf("%s/v1/health/service/%s", registryURL, serviceName)
	resp, err := client.Get(url)
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

func deregisterService(client *http.Client, registryURL, serviceID string) error {
	url := fmt.Sprintf("%s/v1/agent/service/deregister/%s", registryURL, serviceID)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
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

func findServiceByName(services map[string]ServiceDetail, name string) (ServiceDetail, bool) {
	for _, service := range services {
		if service.Name == name {
			return service, true
		}
	}
	return ServiceDetail{}, false
}

func containsAllTags(serviceTags, requiredTags []string) bool {
	tagMap := make(map[string]bool)
	for _, tag := range serviceTags {
		tagMap[tag] = true
	}

	for _, required := range requiredTags {
		if !tagMap[required] {
			return false
		}
	}
	return true
}
