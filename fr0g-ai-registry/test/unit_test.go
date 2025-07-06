package test

import (
	"testing"
	"time"
)

// TestUnitServiceValidation tests service validation logic without network calls
func TestUnitServiceValidation(t *testing.T) {
	t.Run("ValidateServiceInfo", func(t *testing.T) {
		validService := ServiceInfo{
			ID:      "test-service-1",
			Name:    "test-service",
			Address: "localhost",
			Port:    8080,
			Tags:    []string{"api", "v1"},
		}

		// Test valid service
		if !isValidServiceInfo(validService) {
			t.Error("Expected valid service to pass validation")
		}

		// Test missing required fields
		invalidService := ServiceInfo{
			Name: "test-service",
			// Missing ID, Address, Port
		}

		if isValidServiceInfo(invalidService) {
			t.Error("Expected invalid service to fail validation")
		}

		// Test invalid port
		invalidPortService := ServiceInfo{
			ID:      "test-service-2",
			Name:    "test-service",
			Address: "localhost",
			Port:    -1, // Invalid port
		}

		if isValidServiceInfo(invalidPortService) {
			t.Error("Expected service with invalid port to fail validation")
		}
	})

	t.Run("GenerateServiceID", func(t *testing.T) {
		id1 := generateServiceID("test-service", "localhost", 8080)
		id2 := generateServiceID("test-service", "localhost", 8080)
		id3 := generateServiceID("different-service", "localhost", 8080)

		if id1 != id2 {
			t.Error("Expected same service info to generate same ID")
		}

		if id1 == id3 {
			t.Error("Expected different services to generate different IDs")
		}

		if id1 == "" {
			t.Error("Expected non-empty service ID")
		}
	})
}

// TestUnitRegistryLogic tests registry business logic without network dependencies
func TestUnitRegistryLogic(t *testing.T) {
	t.Run("ServiceRegistryOperations", func(t *testing.T) {
		registry := NewInMemoryRegistry()

		// Test service registration
		service := ServiceInfo{
			ID:      "test-service",
			Name:    "test-service",
			Address: "localhost",
			Port:    8080,
			Tags:    []string{"test"},
		}

		err := registry.RegisterService(service)
		if err != nil {
			t.Fatalf("Failed to register service: %v", err)
		}

		// Test service discovery
		services, err := registry.GetServices()
		if err != nil {
			t.Fatalf("Failed to get services: %v", err)
		}

		if len(services) != 1 {
			t.Errorf("Expected 1 service, got %d", len(services))
		}

		if services[0].ID != "test-service" {
			t.Errorf("Expected service ID 'test-service', got %s", services[0].ID)
		}

		// Test service deregistration
		err = registry.DeregisterService("test-service")
		if err != nil {
			t.Fatalf("Failed to deregister service: %v", err)
		}

		services, err = registry.GetServices()
		if err != nil {
			t.Fatalf("Failed to get services after deregistration: %v", err)
		}

		if len(services) != 0 {
			t.Errorf("Expected 0 services after deregistration, got %d", len(services))
		}
	})
}

// TestUnitPerformanceMetrics tests performance calculation without network calls
func TestUnitPerformanceMetrics(t *testing.T) {
	t.Run("CalculateMetrics", func(t *testing.T) {
		// Test with successful operations
		start := time.Now()
		time.Sleep(10 * time.Millisecond) // Simulate work
		duration := time.Since(start)

		successCount := 100
		errorCount := 5
		totalOps := successCount + errorCount

		rate := calculateOperationRate(successCount, duration)
		if rate <= 0 {
			t.Error("Expected positive operation rate")
		}

		successRate := calculateSuccessRate(successCount, totalOps)
		if successRate < 90 || successRate > 100 {
			t.Errorf("Expected success rate between 90-100%%, got %.2f%%", successRate)
		}

		// Test with zero operations (should not panic)
		zeroRate := calculateOperationRate(0, duration)
		if zeroRate != 0 {
			t.Errorf("Expected zero rate for zero operations, got %.2f", zeroRate)
		}

		// Test with zero duration (should not panic)
		zeroDurationRate := calculateOperationRate(100, 0)
		if zeroDurationRate != 0 {
			t.Errorf("Expected zero rate for zero duration, got %.2f", zeroDurationRate)
		}
	})
}

// Helper functions for unit tests (types are imported from types.go)

type InMemoryRegistry struct {
	services map[string]ServiceInfo
}

func NewInMemoryRegistry() *InMemoryRegistry {
	return &InMemoryRegistry{
		services: make(map[string]ServiceInfo),
	}
}

func (r *InMemoryRegistry) RegisterService(service ServiceInfo) error {
	r.services[service.ID] = service
	return nil
}

func (r *InMemoryRegistry) DeregisterService(serviceID string) error {
	delete(r.services, serviceID)
	return nil
}

func (r *InMemoryRegistry) GetServices() ([]ServiceInfo, error) {
	services := make([]ServiceInfo, 0, len(r.services))
	for _, service := range r.services {
		services = append(services, service)
	}
	return services, nil
}

func isValidServiceInfo(service ServiceInfo) bool {
	if service.ID == "" || service.Name == "" || service.Address == "" {
		return false
	}
	
	if service.Port <= 0 || service.Port > 65535 {
		return false
	}
	
	return true
}

func generateServiceID(name, address string, port int) string {
	return name + "-" + address + "-" + string(rune(port+'0'))
}

func calculateOperationRate(operations int, duration time.Duration) float64 {
	if duration.Seconds() <= 0 {
		return 0
	}
	return float64(operations) / duration.Seconds()
}

func calculateSuccessRate(successful, total int) float64 {
	if total <= 0 {
		return 0
	}
	return float64(successful) / float64(total) * 100
}
