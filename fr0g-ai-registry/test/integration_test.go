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


// TestRegistryIntegration tests the complete service registration and discovery workflow
func TestRegistryIntegration(t *testing.T) {
	registryURL := "http://localhost:8500"
	
	// Test services that should integrate with the registry
	testServices := []ServiceInfo{
		{
			ID:      "fr0g-ai-aip-1",
			Name:    "fr0g-ai-aip",
			Address: "localhost",
			Port:    8080,
			Tags:    []string{"ai", "persona", "identity"},
			Meta: map[string]string{
				"version": "1.0.0",
				"env":     "test",
			},
			Check: &HealthCheck{
				HTTP:     "http://localhost:8080/health",
				Interval: "10s",
				Timeout:  "3s",
			},
		},
		{
			ID:      "fr0g-ai-bridge-1",
			Name:    "fr0g-ai-bridge",
			Address: "localhost",
			Port:    8081,
			Tags:    []string{"ai", "bridge", "llm"},
			Meta: map[string]string{
				"version": "1.0.0",
				"env":     "test",
			},
			Check: &HealthCheck{
				HTTP:     "http://localhost:8081/health",
				Interval: "10s",
				Timeout:  "3s",
			},
		},
		{
			ID:      "fr0g-ai-io-1",
			Name:    "fr0g-ai-io",
			Address: "localhost",
			Port:    8082,
			Tags:    []string{"io", "input", "output"},
			Meta: map[string]string{
				"version": "1.0.0",
				"env":     "test",
			},
			Check: &HealthCheck{
				HTTP:     "http://localhost:8082/health",
				Interval: "10s",
				Timeout:  "3s",
			},
		},
		{
			ID:      "fr0g-ai-master-control-1",
			Name:    "fr0g-ai-master-control",
			Address: "localhost",
			Port:    8083,
			Tags:    []string{"control", "orchestration", "workflow"},
			Meta: map[string]string{
				"version": "1.0.0",
				"env":     "test",
			},
			Check: &HealthCheck{
				HTTP:     "http://localhost:8083/health",
				Interval: "10s",
				Timeout:  "3s",
			},
		},
	}

	t.Run("RegistryHealthCheck", func(t *testing.T) {
		testRegistryHealth(t, registryURL)
	})

	t.Run("ServiceRegistration", func(t *testing.T) {
		for _, service := range testServices {
			t.Run(service.Name, func(t *testing.T) {
				testServiceRegistration(t, registryURL, service)
			})
		}
	})

	t.Run("ServiceDiscovery", func(t *testing.T) {
		testServiceDiscovery(t, registryURL, testServices)
	})

	t.Run("ServiceHealth", func(t *testing.T) {
		for _, service := range testServices {
			t.Run(service.Name, func(t *testing.T) {
				testServiceHealth(t, registryURL, service)
			})
		}
	})

	t.Run("ServiceDeregistration", func(t *testing.T) {
		for _, service := range testServices {
			t.Run(service.Name, func(t *testing.T) {
				testServiceDeregistration(t, registryURL, service)
			})
		}
	})
}

// testRegistryHealth verifies the registry service is running and healthy
func testRegistryHealth(t *testing.T, registryURL string) {
	resp, err := http.Get(registryURL + "/health")
	if err != nil {
		t.Fatalf("Failed to connect to registry: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Registry health check failed: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read health response: %v", err)
	}

	t.Logf("Registry health check passed: %s", string(body))
}

// testServiceRegistration tests registering a service with the registry
func testServiceRegistration(t *testing.T, registryURL string, service ServiceInfo) {
	payload, err := json.Marshal(service)
	if err != nil {
		t.Fatalf("Failed to marshal service info: %v", err)
	}

	resp, err := http.NewRequest("PUT", registryURL+"/v1/agent/service/register", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create registration request: %v", err)
	}
	resp.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(resp)
	if err != nil {
		t.Fatalf("Failed to register service: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		t.Fatalf("Service registration failed: status %d, body: %s", response.StatusCode, string(body))
	}

	t.Logf("Successfully registered service: %s", service.Name)
}

// testServiceDiscovery tests discovering registered services
func testServiceDiscovery(t *testing.T, registryURL string, expectedServices []ServiceInfo) {
	resp, err := http.Get(registryURL + "/v1/catalog/services")
	if err != nil {
		t.Fatalf("Failed to discover services: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Service discovery failed: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read discovery response: %v", err)
	}

	// Try to unmarshal as detailed service information first
	var detailedServices map[string]ServiceDetail
	if err := json.Unmarshal(body, &detailedServices); err == nil {
		// Handle detailed service response format
		foundServices := make(map[string]bool)
		for _, serviceDetail := range detailedServices {
			foundServices[serviceDetail.Name] = true
		}

		// Verify expected services are discoverable
		for _, expected := range expectedServices {
			if foundServices[expected.Name] {
				t.Logf("Found service %s in detailed discovery", expected.Name)
			} else {
				t.Errorf("Expected service %s not found in discovery", expected.Name)
			}
		}

		t.Logf("Service discovery completed. Found %d services", len(detailedServices))
		return
	}

	// Fallback to simple service list format
	var services map[string][]string
	if err := json.Unmarshal(body, &services); err != nil {
		t.Fatalf("Failed to unmarshal services in any expected format: %v", err)
	}

	// Verify expected services are discoverable
	for _, expected := range expectedServices {
		if tags, exists := services[expected.Name]; exists {
			t.Logf("Found service %s with tags: %v", expected.Name, tags)
		} else {
			t.Errorf("Expected service %s not found in discovery", expected.Name)
		}
	}

	t.Logf("Service discovery completed. Found %d services", len(services))
}

// testServiceHealth tests health checking for registered services
func testServiceHealth(t *testing.T, registryURL string, service ServiceInfo) {
	url := fmt.Sprintf("%s/v1/health/service/%s", registryURL, service.Name)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to check service health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Service health check returned status %d for %s: %s", resp.StatusCode, service.Name, string(body))
		return // Health check may fail if service isn't actually running
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read health response: %v", err)
	}

	t.Logf("Health check for %s: %s", service.Name, string(body))
}

// testServiceDeregistration tests deregistering services
func testServiceDeregistration(t *testing.T, registryURL string, service ServiceInfo) {
	url := fmt.Sprintf("%s/v1/agent/service/deregister/%s", registryURL, service.ID)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		t.Fatalf("Failed to create deregistration request: %v", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to deregister service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Service deregistration failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	t.Logf("Successfully deregistered service: %s", service.Name)
}

// TestRegistryPerformance tests registry performance under load
func TestRegistryPerformance(t *testing.T) {
	registryURL := "http://localhost:8500"
	
	// Test concurrent service registrations
	t.Run("ConcurrentRegistrations", func(t *testing.T) {
		numServices := 50
		done := make(chan bool, numServices)
		
		start := time.Now()
		
		for i := 0; i < numServices; i++ {
			go func(id int) {
				service := ServiceInfo{
					ID:      fmt.Sprintf("test-service-%d", id),
					Name:    fmt.Sprintf("test-service-%d", id),
					Address: "localhost",
					Port:    8000 + id,
					Tags:    []string{"test", "performance"},
				}
				
				payload, _ := json.Marshal(service)
				req, _ := http.NewRequest("PUT", registryURL+"/v1/agent/service/register", bytes.NewBuffer(payload))
				req.Header.Set("Content-Type", "application/json")
				
				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					t.Errorf("Failed to register service %d: %v", id, err)
				} else {
					resp.Body.Close()
				}
				
				done <- true
			}(i)
		}
		
		// Wait for all registrations to complete
		for i := 0; i < numServices; i++ {
			<-done
		}
		
		duration := time.Since(start)
		t.Logf("Registered %d services in %v (%.2f services/sec)", 
			numServices, duration, float64(numServices)/duration.Seconds())
		
		// Clean up test services
		for i := 0; i < numServices; i++ {
			serviceID := fmt.Sprintf("test-service-%d", i)
			url := fmt.Sprintf("%s/v1/agent/service/deregister/%s", registryURL, serviceID)
			req, _ := http.NewRequest("PUT", url, nil)
			client := &http.Client{Timeout: 5 * time.Second}
			resp, _ := client.Do(req)
			if resp != nil {
				resp.Body.Close()
			}
		}
	})
	
	// Test service discovery performance
	t.Run("DiscoveryPerformance", func(t *testing.T) {
		numRequests := 100
		start := time.Now()
		
		for i := 0; i < numRequests; i++ {
			resp, err := http.Get(registryURL + "/v1/catalog/services")
			if err != nil {
				t.Errorf("Discovery request %d failed: %v", i, err)
				continue
			}
			resp.Body.Close()
		}
		
		duration := time.Since(start)
		avgLatency := duration / time.Duration(numRequests)
		t.Logf("Completed %d discovery requests in %v (avg latency: %v)", 
			numRequests, duration, avgLatency)
	})
}
