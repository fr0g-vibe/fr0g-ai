package monitor

import (
	"log"
	"math/rand"
)

// SystemMonitor handles system monitoring for the MCP
type SystemMonitor struct {
}

// NewSystemMonitor creates a new system monitor
func NewSystemMonitor() *SystemMonitor {
	return &SystemMonitor{}
}

// Start begins system monitoring
func (sm *SystemMonitor) Start() error {
	log.Println("System Monitor: Starting monitoring processes...")
	return nil
}

// Stop gracefully stops the system monitor
func (sm *SystemMonitor) Stop() error {
	log.Println("System Monitor: Stopping monitoring processes...")
	return nil
}

// GetSystemLoad returns current system load (mock implementation)
func (sm *SystemMonitor) GetSystemLoad() float64 {
	// Return a mock system load value
	return rand.Float64() * 0.1 // Low load for demo
}
