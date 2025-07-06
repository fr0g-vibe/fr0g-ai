package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/registry"
)

// Example client demonstrating service discovery usage
func main() {
	log.Println("🔍 Service Discovery Client Example")
	
	// Create registry client
	registryClient := registry.NewRegistryClient("http://localhost:8500", nil)
	
	// Discover AIP service
	fmt.Println("\n📋 Discovering AIP service...")
	aipEndpoint, err := registryClient.GetServiceEndpoint("fr0g-ai-aip")
	if err != nil {
		log.Printf("Failed to discover AIP service: %v", err)
	} else {
		fmt.Printf("✅ AIP service found at: %s\n", aipEndpoint)
		
		// Test connection to AIP service
		resp, err := http.Get(aipEndpoint + "/health")
		if err != nil {
			log.Printf("Failed to connect to AIP service: %v", err)
		} else {
			defer resp.Body.Close()
			fmt.Printf("✅ AIP service health check: %s\n", resp.Status)
		}
	}
	
	// Discover Bridge service
	fmt.Println("\n🌉 Discovering Bridge service...")
	bridgeEndpoint, err := registryClient.GetServiceEndpoint("fr0g-ai-bridge")
	if err != nil {
		log.Printf("Failed to discover Bridge service: %v", err)
	} else {
		fmt.Printf("✅ Bridge service found at: %s\n", bridgeEndpoint)
		
		// Test connection to Bridge service
		resp, err := http.Get(bridgeEndpoint + "/health")
		if err != nil {
			log.Printf("Failed to connect to Bridge service: %v", err)
		} else {
			defer resp.Body.Close()
			fmt.Printf("✅ Bridge service health check: %s\n", resp.Status)
		}
	}
	
	// Get all healthy services
	fmt.Println("\n🔍 Getting all healthy services...")
	services := []string{"fr0g-ai-aip", "fr0g-ai-bridge", "fr0g-ai-io", "fr0g-ai-master-control"}
	
	for _, serviceName := range services {
		healthyServices, err := registryClient.GetHealthyServices(serviceName)
		if err != nil {
			log.Printf("Failed to get healthy instances of %s: %v", serviceName, err)
			continue
		}
		
		fmt.Printf("✅ %s: %d healthy instances\n", serviceName, len(healthyServices))
		for _, service := range healthyServices {
			fmt.Printf("   - %s:%d (ID: %s)\n", service.Address, service.Port, service.ID)
		}
	}
	
	// Demonstrate service registry monitoring
	fmt.Println("\n📊 Monitoring service registry...")
	for i := 0; i < 3; i++ {
		resp, err := http.Get("http://localhost:8500/v1/catalog/services")
		if err != nil {
			log.Printf("Failed to query service registry: %v", err)
			continue
		}
		defer resp.Body.Close()
		
		var services map[string][]string
		if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
			log.Printf("Failed to decode services: %v", err)
			continue
		}
		
		fmt.Printf("📋 Iteration %d - Registered services: %d\n", i+1, len(services))
		for name, tags := range services {
			fmt.Printf("   - %s: %v\n", name, tags)
		}
		
		time.Sleep(5 * time.Second)
	}
	
	fmt.Println("\n✅ Service discovery client example completed!")
}
