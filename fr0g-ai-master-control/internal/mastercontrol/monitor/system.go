package monitor

import (
	"context"
	"log"
	"runtime"
	"sync"
	"time"
)

// SystemMonitor handles system monitoring for the MCP
type SystemMonitor struct {
	metrics *SystemMetrics
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
	config  *MonitorConfig
}

// SystemMetrics holds current system metrics
type SystemMetrics struct {
	CPUUsage        float64                    `json:"cpu_usage"`
	MemoryUsage     float64                    `json:"memory_usage"`
	GoroutineCount  int                        `json:"goroutine_count"`
	HeapSize        uint64                     `json:"heap_size"`
	SystemLoad      float64                    `json:"system_load"`
	LastUpdate      time.Time                  `json:"last_update"`
	ComponentHealth map[string]ComponentHealth `json:"component_health"`
}

// ComponentHealth represents health status of a component
type ComponentHealth struct {
	Status       string        `json:"status"` // "healthy", "warning", "critical"
	LastCheck    time.Time     `json:"last_check"`
	ResponseTime time.Duration `json:"response_time"`
	ErrorCount   int           `json:"error_count"`
	Uptime       time.Duration `json:"uptime"`
}

// MonitorConfig holds monitoring configuration
type MonitorConfig struct {
	UpdateInterval      time.Duration   `yaml:"update_interval"`
	HealthCheckInterval time.Duration   `yaml:"health_check_interval"`
	AlertThresholds     AlertThresholds `yaml:"alert_thresholds"`
}

// AlertThresholds defines when to alert on metrics
type AlertThresholds struct {
	CPUWarning     float64 `yaml:"cpu_warning"`
	CPUCritical    float64 `yaml:"cpu_critical"`
	MemoryWarning  float64 `yaml:"memory_warning"`
	MemoryCritical float64 `yaml:"memory_critical"`
}

// NewSystemMonitor creates a new system monitor
func NewSystemMonitor() *SystemMonitor {
	ctx, cancel := context.WithCancel(context.Background())

	config := &MonitorConfig{
		UpdateInterval:      time.Second * 5,
		HealthCheckInterval: time.Second * 30,
		AlertThresholds: AlertThresholds{
			CPUWarning:     70.0,
			CPUCritical:    90.0,
			MemoryWarning:  80.0,
			MemoryCritical: 95.0,
		},
	}

	return &SystemMonitor{
		metrics: &SystemMetrics{
			ComponentHealth: make(map[string]ComponentHealth),
			LastUpdate:      time.Now(),
		},
		ctx:    ctx,
		cancel: cancel,
		config: config,
	}
}

// Start begins system monitoring
func (sm *SystemMonitor) Start() error {
	log.Println("System Monitor: Starting real-time monitoring processes...")

	// Start metrics collection loop
	go sm.metricsLoop()

	// Start health check loop
	go sm.healthCheckLoop()

	log.Println("System Monitor: Real-time monitoring started successfully")
	return nil
}

// Stop gracefully stops the system monitor
func (sm *SystemMonitor) Stop() error {
	log.Println("System Monitor: Stopping monitoring processes...")
	sm.cancel()
	return nil
}

// GetSystemLoad returns current system load (real implementation)
func (sm *SystemMonitor) GetSystemLoad() float64 {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.metrics.SystemLoad
}

// GetMetrics returns current system metrics
func (sm *SystemMonitor) GetMetrics() *SystemMetrics {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Return a copy
	metrics := *sm.metrics
	return &metrics
}

// RegisterComponent registers a component for health monitoring
func (sm *SystemMonitor) RegisterComponent(name string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.metrics.ComponentHealth[name] = ComponentHealth{
		Status:    "healthy",
		LastCheck: time.Now(),
		Uptime:    0,
	}

	log.Printf("System Monitor: Registered component '%s' for health monitoring", name)
}

// UpdateComponentHealth updates health status for a component
func (sm *SystemMonitor) UpdateComponentHealth(name string, status string, responseTime time.Duration) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if health, exists := sm.metrics.ComponentHealth[name]; exists {
		health.Status = status
		health.LastCheck = time.Now()
		health.ResponseTime = responseTime
		if status != "healthy" {
			health.ErrorCount++
		}
		sm.metrics.ComponentHealth[name] = health
	}
}

// metricsLoop continuously updates system metrics
func (sm *SystemMonitor) metricsLoop() {
	ticker := time.NewTicker(sm.config.UpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-sm.ctx.Done():
			return
		case <-ticker.C:
			sm.updateMetrics()
		}
	}
}

// healthCheckLoop performs periodic health checks
func (sm *SystemMonitor) healthCheckLoop() {
	ticker := time.NewTicker(sm.config.HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-sm.ctx.Done():
			return
		case <-ticker.C:
			sm.performHealthChecks()
		}
	}
}

// updateMetrics collects current system metrics
func (sm *SystemMonitor) updateMetrics() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Calculate CPU usage (simplified)
	sm.metrics.CPUUsage = sm.calculateCPUUsage()

	// Memory usage
	sm.metrics.MemoryUsage = float64(memStats.Alloc) / float64(memStats.Sys) * 100

	// Goroutine count
	sm.metrics.GoroutineCount = runtime.NumGoroutine()

	// Heap size
	sm.metrics.HeapSize = memStats.HeapAlloc

	// System load (composite metric)
	sm.metrics.SystemLoad = sm.calculateSystemLoad()

	sm.metrics.LastUpdate = time.Now()

	// Check for alerts
	sm.checkAlerts()
}

// calculateCPUUsage calculates CPU usage (simplified implementation)
func (sm *SystemMonitor) calculateCPUUsage() float64 {
	// This is a simplified CPU calculation
	// In production, you'd want to use proper CPU monitoring
	goroutines := float64(runtime.NumGoroutine())

	// Estimate CPU usage based on goroutine activity
	cpuUsage := (goroutines / 100.0) * 10.0
	if cpuUsage > 100.0 {
		cpuUsage = 100.0
	}

	return cpuUsage
}

// calculateSystemLoad calculates overall system load
func (sm *SystemMonitor) calculateSystemLoad() float64 {
	// Composite load based on CPU, memory, and goroutines
	cpuWeight := 0.4
	memoryWeight := 0.4
	goroutineWeight := 0.2

	normalizedGoroutines := float64(sm.metrics.GoroutineCount) / 1000.0
	if normalizedGoroutines > 1.0 {
		normalizedGoroutines = 1.0
	}

	load := (sm.metrics.CPUUsage/100.0)*cpuWeight +
		(sm.metrics.MemoryUsage/100.0)*memoryWeight +
		normalizedGoroutines*goroutineWeight

	return load
}

// performHealthChecks checks health of all registered components
func (sm *SystemMonitor) performHealthChecks() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for name, health := range sm.metrics.ComponentHealth {
		// Update uptime
		health.Uptime = time.Since(health.LastCheck)

		// Simple health check based on last update time
		timeSinceLastCheck := time.Since(health.LastCheck)
		if timeSinceLastCheck > time.Minute*5 {
			health.Status = "warning"
		}
		if timeSinceLastCheck > time.Minute*10 {
			health.Status = "critical"
		}

		sm.metrics.ComponentHealth[name] = health
	}
}

// checkAlerts checks if any metrics exceed alert thresholds
func (sm *SystemMonitor) checkAlerts() {
	if sm.metrics.CPUUsage > sm.config.AlertThresholds.CPUCritical {
		log.Printf("System Monitor: CRITICAL - CPU usage at %.2f%%", sm.metrics.CPUUsage)
	} else if sm.metrics.CPUUsage > sm.config.AlertThresholds.CPUWarning {
		log.Printf("System Monitor: WARNING - CPU usage at %.2f%%", sm.metrics.CPUUsage)
	}

	if sm.metrics.MemoryUsage > sm.config.AlertThresholds.MemoryCritical {
		log.Printf("System Monitor: CRITICAL - Memory usage at %.2f%%", sm.metrics.MemoryUsage)
	} else if sm.metrics.MemoryUsage > sm.config.AlertThresholds.MemoryWarning {
		log.Printf("System Monitor: WARNING - Memory usage at %.2f%%", sm.metrics.MemoryUsage)
	}
}
