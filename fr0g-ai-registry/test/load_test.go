package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

// TestRegistryLoadUnderConcurrentServices tests registry performance with multiple services
func TestRegistryLoadUnderConcurrentServices(t *testing.T) {
	registryURL := "http://localhost:8500"
	client := &http.Client{Timeout: 5 * time.Second}

	t.Run("ConcurrentServiceRegistrations", func(t *testing.T) {
		numServices := 100
		var wg sync.WaitGroup
		var mu sync.Mutex
		var successCount, errorCount int

		start := time.Now()

		for i := 0; i < numServices; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				service := ServiceInfo{
					ID:      fmt.Sprintf("load-test-service-%d", id),
					Name:    fmt.Sprintf("load-test-service-%d", id),
					Address: "localhost",
					Port:    10000 + id,
					Tags:    []string{"load-test", "concurrent", fmt.Sprintf("batch-%d", id/10)},
					Meta: map[string]string{
						"test_id":     fmt.Sprintf("%d", id),
						"batch":       fmt.Sprintf("%d", id/10),
						"created_at":  time.Now().Format(time.RFC3339),
					},
				}

				payload, _ := json.Marshal(service)
				req, _ := http.NewRequest("PUT", registryURL+"/v1/agent/service/register", 
					bytes.NewBuffer(payload))
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				
				mu.Lock()
				if err != nil || resp.StatusCode != http.StatusOK {
					errorCount++
					if err != nil {
						t.Logf("Registration error for service %d: %v", id, err)
					} else {
						t.Logf("Registration failed for service %d: status %d", id, resp.StatusCode)
					}
				} else {
					successCount++
				}
				mu.Unlock()

				if resp != nil {
					resp.Body.Close()
				}
			}(i)
		}

		wg.Wait()
		duration := time.Since(start)

		t.Logf("Concurrent registration results:")
		t.Logf("  Total services: %d", numServices)
		t.Logf("  Successful: %d", successCount)
		t.Logf("  Failed: %d", errorCount)
		t.Logf("  Duration: %v", duration)
		t.Logf("  Rate: %.2f registrations/sec", float64(successCount)/duration.Seconds())

		if successCount < numServices*8/10 { // Allow 20% failure rate under load
			t.Errorf("Too many registration failures: %d/%d succeeded", successCount, numServices)
		}

		// Cleanup
		t.Logf("Cleaning up %d test services...", successCount)
		cleanupStart := time.Now()
		for i := 0; i < numServices; i++ {
			serviceID := fmt.Sprintf("load-test-service-%d", i)
			url := fmt.Sprintf("%s/v1/agent/service/deregister/%s", registryURL, serviceID)
			req, _ := http.NewRequest("PUT", url, nil)
			resp, _ := client.Do(req)
			if resp != nil {
				resp.Body.Close()
			}
		}
		t.Logf("Cleanup completed in %v", time.Since(cleanupStart))
	})

	t.Run("ConcurrentServiceDiscovery", func(t *testing.T) {
		// First register some services for discovery
		numTestServices := 20
		for i := 0; i < numTestServices; i++ {
			service := ServiceInfo{
				ID:      fmt.Sprintf("discovery-load-service-%d", i),
				Name:    fmt.Sprintf("discovery-load-service-%d", i),
				Address: "localhost",
				Port:    11000 + i,
				Tags:    []string{"discovery-load", "test"},
			}

			payload, _ := json.Marshal(service)
			req, _ := http.NewRequest("PUT", registryURL+"/v1/agent/service/register", 
				bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := client.Do(req)
			if resp != nil {
				resp.Body.Close()
			}
		}

		// Now test concurrent discovery
		numDiscoveryRequests := 200
		var wg sync.WaitGroup
		var mu sync.Mutex
		var successCount, errorCount int
		var totalLatency time.Duration

		start := time.Now()

		for i := 0; i < numDiscoveryRequests; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				requestStart := time.Now()
				resp, err := client.Get(registryURL + "/v1/catalog/services")
				requestDuration := time.Since(requestStart)

				mu.Lock()
				totalLatency += requestDuration
				if err != nil || resp.StatusCode != http.StatusOK {
					errorCount++
				} else {
					successCount++
				}
				mu.Unlock()

				if resp != nil {
					resp.Body.Close()
				}
			}(i)
		}

		wg.Wait()
		duration := time.Since(start)
		var avgLatency time.Duration
		if successCount > 0 {
			avgLatency = totalLatency / time.Duration(successCount)
		}

		t.Logf("Concurrent discovery results:")
		t.Logf("  Total requests: %d", numDiscoveryRequests)
		t.Logf("  Successful: %d", successCount)
		t.Logf("  Failed: %d", errorCount)
		t.Logf("  Duration: %v", duration)
		t.Logf("  Rate: %.2f requests/sec", float64(successCount)/duration.Seconds())
		t.Logf("  Average latency: %v", avgLatency)

		if avgLatency > 100*time.Millisecond {
			t.Logf("Warning: Average discovery latency is high: %v", avgLatency)
		}

		// Cleanup discovery test services
		for i := 0; i < numTestServices; i++ {
			serviceID := fmt.Sprintf("discovery-load-service-%d", i)
			url := fmt.Sprintf("%s/v1/agent/service/deregister/%s", registryURL, serviceID)
			req, _ := http.NewRequest("PUT", url, nil)
			resp, _ := client.Do(req)
			if resp != nil {
				resp.Body.Close()
			}
		}
	})

	t.Run("MixedOperationsLoad", func(t *testing.T) {
		// Test mixed operations: register, discover, health check, deregister
		numOperations := 300
		var wg sync.WaitGroup
		var mu sync.Mutex
		var opCounts = make(map[string]int)
		var opErrors = make(map[string]int)

		start := time.Now()

		for i := 0; i < numOperations; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				// Randomly choose operation type
				opType := []string{"register", "discover", "health", "deregister"}[id%4]
				
				mu.Lock()
				opCounts[opType]++
				mu.Unlock()

				var err error
				switch opType {
				case "register":
					err = performRegistration(client, registryURL, id)
				case "discover":
					err = performDiscovery(client, registryURL)
				case "health":
					err = performHealthCheck(client, registryURL)
				case "deregister":
					err = performDeregistration(client, registryURL, id)
				}

				if err != nil {
					mu.Lock()
					opErrors[opType]++
					mu.Unlock()
				}
			}(i)
		}

		wg.Wait()
		duration := time.Since(start)

		t.Logf("Mixed operations load test results:")
		t.Logf("  Total operations: %d", numOperations)
		t.Logf("  Duration: %v", duration)
		t.Logf("  Rate: %.2f operations/sec", float64(numOperations)/duration.Seconds())

		for opType, count := range opCounts {
			errors := opErrors[opType]
			successRate := float64(count-errors) / float64(count) * 100
			t.Logf("  %s: %d operations, %d errors (%.1f%% success)", 
				opType, count, errors, successRate)
		}
	})
}

// Helper functions for load testing

func performRegistration(client *http.Client, registryURL string, id int) error {
	service := ServiceInfo{
		ID:      fmt.Sprintf("mixed-test-service-%d", id),
		Name:    fmt.Sprintf("mixed-test-service-%d", id),
		Address: "localhost",
		Port:    12000 + id,
		Tags:    []string{"mixed-test"},
	}

	payload, _ := json.Marshal(service)
	req, _ := http.NewRequest("PUT", registryURL+"/v1/agent/service/register", 
		bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if resp != nil {
		resp.Body.Close()
	}
	return err
}

func performDiscovery(client *http.Client, registryURL string) error {
	resp, err := client.Get(registryURL + "/v1/catalog/services")
	if resp != nil {
		resp.Body.Close()
	}
	return err
}

func performHealthCheck(client *http.Client, registryURL string) error {
	resp, err := client.Get(registryURL + "/health")
	if resp != nil {
		resp.Body.Close()
	}
	return err
}

func performDeregistration(client *http.Client, registryURL string, id int) error {
	serviceID := fmt.Sprintf("mixed-test-service-%d", id)
	url := fmt.Sprintf("%s/v1/agent/service/deregister/%s", registryURL, serviceID)
	req, _ := http.NewRequest("PUT", url, nil)

	resp, err := client.Do(req)
	if resp != nil {
		resp.Body.Close()
	}
	return err
}
