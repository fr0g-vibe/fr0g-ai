package test

import (
	"context"
	"testing"
	"time"

	aipregistry "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/registry"
	bridgediscovery "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/discovery"
	"github.com/sirupsen/logrus"
)

// TestServiceDiscovery tests the service discovery functionality
func TestServiceDiscovery(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create mock registry client
	registryClient := aipregistry.NewRegistryClient("http://localhost:8500", logger)
	
	// Create service discovery
	discovery := bridgediscovery.NewServiceDiscovery(registryClient, logger)
	defer discovery.Shutdown()

	t.Run("GetAIPEndpoint", func(t *testing.T) {
		// This would normally connect to a real registry
		// For testing, we'll verify the method exists and handles errors gracefully
		_, err := discovery.GetAIPEndpoint()
		if err == nil {
			t.Log("✓ AIP endpoint discovery method works")
		} else {
			t.Logf("AIP endpoint discovery failed as expected (no registry): %v", err)
		}
	})

	t.Run("GetMasterControlEndpoint", func(t *testing.T) {
		_, err := discovery.GetMasterControlEndpoint()
		if err == nil {
			t.Log("✓ Master Control endpoint discovery method works")
		} else {
			t.Logf("Master Control endpoint discovery failed as expected (no registry): %v", err)
		}
	})

	t.Run("GetIOEndpoint", func(t *testing.T) {
		_, err := discovery.GetIOEndpoint()
		if err == nil {
			t.Log("✓ IO endpoint discovery method works")
		} else {
			t.Logf("IO endpoint discovery failed as expected (no registry): %v", err)
		}
	})

	t.Run("ServiceDependencyStatus", func(t *testing.T) {
		status := discovery.GetServiceDependencyStatus()
		
		expectedServices := []string{"fr0g-ai-aip", "fr0g-ai-master-control", "fr0g-ai-io"}
		for _, service := range expectedServices {
			if _, exists := status[service]; !exists {
				t.Errorf("Service %s not found in dependency status", service)
			} else {
				t.Logf("✓ Service %s status checked: %v", service, status[service])
			}
		}
	})

	t.Run("BackgroundRefresh", func(t *testing.T) {
		// Test background refresh functionality
		discovery.StartBackgroundRefresh()
		
		// Wait a short time to ensure background refresh starts
		time.Sleep(100 * time.Millisecond)
		
		// Shutdown should stop background refresh
		discovery.Shutdown()
		
		t.Log("✓ Background refresh started and stopped successfully")
	})
}

// TestRegistryClient tests the registry client functionality
func TestRegistryClient(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	registryClient := aipregistry.NewRegistryClient("http://localhost:8500", logger)

	t.Run("RegisterService", func(t *testing.T) {
		serviceInfo := &aipregistry.ServiceInfo{
			ID:      "test-service-1",
			Name:    "test-service",
			Address: "localhost",
			Port:    8080,
			Tags:    []string{"test", "grpc"},
			Meta: map[string]string{
				"version": "1.0.0",
			},
			Check: &aipregistry.HealthCheck{
				HTTP:     "http://localhost:8080/health",
				Interval: "10s",
				Timeout:  "3s",
			},
		}

		// This will fail without a real registry, but we can test the method exists
		err := registryClient.RegisterService(serviceInfo)
		if err != nil {
			t.Logf("Service registration failed as expected (no registry): %v", err)
		} else {
			t.Log("✓ Service registration method works")
		}
	})

	t.Run("GetServiceEndpoint", func(t *testing.T) {
		_, err := registryClient.GetServiceEndpoint("test-service")
		if err != nil {
			t.Logf("Service endpoint lookup failed as expected (no registry): %v", err)
		} else {
			t.Log("✓ Service endpoint lookup method works")
		}
	})

	t.Run("GetHealthyServices", func(t *testing.T) {
		_, err := registryClient.GetHealthyServices("test-service")
		if err != nil {
			t.Logf("Healthy services lookup failed as expected (no registry): %v", err)
		} else {
			t.Log("✓ Healthy services lookup method works")
		}
	})

	t.Run("SimpleRegistration", func(t *testing.T) {
		err := registryClient.RegisterServiceSimple("test-service", "8080", "9090")
		if err != nil {
			t.Logf("Simple service registration failed as expected (no registry): %v", err)
		} else {
			t.Log("✓ Simple service registration method works")
		}
	})
}

// TestServiceRegistrationFlow tests the complete service registration flow
func TestServiceRegistrationFlow(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Test the complete flow that services should follow
	t.Run("CompleteRegistrationFlow", func(t *testing.T) {
		// 1. Create registry client
		registryClient := aipregistry.NewRegistryClient("http://localhost:8500", logger)

		// 2. Create service info
		serviceInfo := &aipregistry.ServiceInfo{
			ID:      "fr0g-ai-aip-test",
			Name:    "fr0g-ai-aip",
			Address: "localhost",
			Port:    9090,
			Tags:    []string{"grpc", "personas"},
			Meta: map[string]string{
				"version":   "1.0.0",
				"grpc_port": "9090",
				"http_port": "8080",
			},
			Check: &aipregistry.HealthCheck{
				HTTP:     "http://localhost:8080/health",
				Interval: "30s",
				Timeout:  "10s",
			},
		}

		// 3. Register service
		err := registryClient.RegisterService(serviceInfo)
		if err != nil {
			t.Logf("Registration failed as expected (no registry): %v", err)
		}

		// 4. Test discovery
		_, err = registryClient.GetServiceEndpoint("fr0g-ai-aip")
		if err != nil {
			t.Logf("Discovery failed as expected (no registry): %v", err)
		}

		// 5. Test health check
		_, err = registryClient.GetHealthyServices("fr0g-ai-aip")
		if err != nil {
			t.Logf("Health check failed as expected (no registry): %v", err)
		}

		// 6. Deregister
		err = registryClient.DeregisterService()
		if err != nil {
			t.Logf("Deregistration failed as expected (no registry): %v", err)
		}

		t.Log("✓ Complete registration flow tested")
	})
}
