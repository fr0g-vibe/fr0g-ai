package discovery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/registry"
	"github.com/sirupsen/logrus"
)

// ServiceDiscovery manages service discovery and endpoint resolution
type ServiceDiscovery struct {
	registryClient *registry.RegistryClient
	logger         *logrus.Logger
	
	// Endpoint cache
	endpointCache map[string]string
	cacheMutex    sync.RWMutex
	cacheExpiry   time.Duration
	lastUpdate    map[string]time.Time
	
	// Background refresh
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewServiceDiscovery creates a new service discovery manager
func NewServiceDiscovery(registryClient *registry.RegistryClient, logger *logrus.Logger) *ServiceDiscovery {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &ServiceDiscovery{
		registryClient: registryClient,
		logger:         logger,
		endpointCache:  make(map[string]string),
		lastUpdate:     make(map[string]time.Time),
		cacheExpiry:    60 * time.Second,
		ctx:            ctx,
		cancel:         cancel,
	}
}

// GetAIPEndpoint returns the endpoint for the AIP service
func (sd *ServiceDiscovery) GetAIPEndpoint() (string, error) {
	return sd.getServiceEndpoint("fr0g-ai-aip")
}

// GetMasterControlEndpoint returns the endpoint for the Master Control service
func (sd *ServiceDiscovery) GetMasterControlEndpoint() (string, error) {
	return sd.getServiceEndpoint("fr0g-ai-master-control")
}

// GetIOEndpoint returns the endpoint for the I/O service
func (sd *ServiceDiscovery) GetIOEndpoint() (string, error) {
	return sd.getServiceEndpoint("fr0g-ai-io")
}

// getServiceEndpoint gets an endpoint for any service with caching
func (sd *ServiceDiscovery) getServiceEndpoint(serviceName string) (string, error) {
	// Check cache first
	sd.cacheMutex.RLock()
	if endpoint, exists := sd.endpointCache[serviceName]; exists {
		if lastUpdate, hasUpdate := sd.lastUpdate[serviceName]; hasUpdate {
			if time.Since(lastUpdate) < sd.cacheExpiry {
				sd.cacheMutex.RUnlock()
				return endpoint, nil
			}
		}
	}
	sd.cacheMutex.RUnlock()
	
	// Fetch from registry
	endpoint, err := sd.registryClient.GetServiceEndpoint(serviceName)
	if err != nil {
		return "", fmt.Errorf("failed to discover %s service: %w", serviceName, err)
	}
	
	// Update cache
	sd.cacheMutex.Lock()
	sd.endpointCache[serviceName] = endpoint
	sd.lastUpdate[serviceName] = time.Now()
	sd.cacheMutex.Unlock()
	
	sd.logger.WithFields(logrus.Fields{
		"service":  serviceName,
		"endpoint": endpoint,
	}).Debug("Service endpoint discovered")
	
	return endpoint, nil
}

// GetHealthyServices returns all healthy services for a given name
func (sd *ServiceDiscovery) GetHealthyServices(serviceName string) ([]*registry.ServiceInfo, error) {
	return sd.registryClient.GetHealthyServices(serviceName)
}

// StartBackgroundRefresh starts background endpoint refresh
func (sd *ServiceDiscovery) StartBackgroundRefresh() {
	sd.wg.Add(1)
	go sd.refreshLoop()
}

// Shutdown stops the service discovery manager
func (sd *ServiceDiscovery) Shutdown() {
	sd.cancel()
	sd.wg.Wait()
}

// refreshLoop periodically refreshes cached endpoints
func (sd *ServiceDiscovery) refreshLoop() {
	defer sd.wg.Done()
	
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-sd.ctx.Done():
			return
		case <-ticker.C:
			sd.refreshEndpoints()
		}
	}
}

// refreshEndpoints refreshes all cached endpoints
func (sd *ServiceDiscovery) refreshEndpoints() {
	sd.cacheMutex.RLock()
	services := make([]string, 0, len(sd.endpointCache))
	for serviceName := range sd.endpointCache {
		services = append(services, serviceName)
	}
	sd.cacheMutex.RUnlock()
	
	for _, serviceName := range services {
		if endpoint, err := sd.registryClient.GetServiceEndpoint(serviceName); err == nil {
			sd.cacheMutex.Lock()
			sd.endpointCache[serviceName] = endpoint
			sd.lastUpdate[serviceName] = time.Now()
			sd.cacheMutex.Unlock()
		} else {
			sd.logger.WithFields(logrus.Fields{
				"service": serviceName,
				"error":   err,
			}).Warn("Failed to refresh service endpoint")
		}
	}
}

// CheckServiceHealth checks if a service is healthy
func (sd *ServiceDiscovery) CheckServiceHealth(serviceName string) (bool, error) {
	services, err := sd.registryClient.GetHealthyServices(serviceName)
	if err != nil {
		return false, err
	}
	
	return len(services) > 0, nil
}

// GetServiceDependencyStatus returns the health status of all service dependencies
func (sd *ServiceDiscovery) GetServiceDependencyStatus() map[string]bool {
	dependencies := []string{"fr0g-ai-aip", "fr0g-ai-master-control", "fr0g-ai-io"}
	status := make(map[string]bool)
	
	for _, dep := range dependencies {
		healthy, err := sd.CheckServiceHealth(dep)
		if err != nil {
			sd.logger.WithFields(logrus.Fields{
				"service": dep,
				"error":   err,
			}).Debug("Failed to check service health")
			status[dep] = false
		} else {
			status[dep] = healthy
		}
	}
	
	return status
}
