package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

// BenchmarkServiceRegistration benchmarks service registration performance
func BenchmarkServiceRegistration(b *testing.B) {
	registryURL := "http://localhost:8500"
	client := &http.Client{Timeout: 5 * time.Second}

	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		service := ServiceInfo{
			ID:      fmt.Sprintf("bench-service-%d", i),
			Name:    fmt.Sprintf("bench-service-%d", i),
			Address: "localhost",
			Port:    8000 + (i % 1000), // Avoid port conflicts
			Tags:    []string{"benchmark", "test"},
		}

		payload, _ := json.Marshal(service)
		req, _ := http.NewRequest("PUT", registryURL+"/v1/agent/service/register", 
			bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("Registration failed: %v", err)
		}
		resp.Body.Close()
	}
}

// BenchmarkServiceDiscovery benchmarks service discovery performance
func BenchmarkServiceDiscovery(b *testing.B) {
	registryURL := "http://localhost:8500"
	client := &http.Client{Timeout: 5 * time.Second}

	// Register some services first
	for i := 0; i < 10; i++ {
		service := ServiceInfo{
			ID:      fmt.Sprintf("discovery-bench-service-%d", i),
			Name:    fmt.Sprintf("discovery-bench-service-%d", i),
			Address: "localhost",
			Port:    9000 + i,
			Tags:    []string{"benchmark", "discovery"},
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

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := client.Get(registryURL + "/v1/catalog/services")
		if err != nil {
			b.Fatalf("Discovery failed: %v", err)
		}
		resp.Body.Close()
	}
}

// BenchmarkServiceDeregistration benchmarks service deregistration performance
func BenchmarkServiceDeregistration(b *testing.B) {
	registryURL := "http://localhost:8500"
	client := &http.Client{Timeout: 5 * time.Second}

	// Pre-register services for deregistration
	serviceIDs := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		serviceID := fmt.Sprintf("dereg-bench-service-%d", i)
		serviceIDs[i] = serviceID
		
		service := ServiceInfo{
			ID:      serviceID,
			Name:    serviceID,
			Address: "localhost",
			Port:    7000 + (i % 1000),
			Tags:    []string{"benchmark", "deregistration"},
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

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		url := fmt.Sprintf("%s/v1/agent/service/deregister/%s", registryURL, serviceIDs[i])
		req, _ := http.NewRequest("PUT", url, nil)

		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("Deregistration failed: %v", err)
		}
		resp.Body.Close()
	}
}
