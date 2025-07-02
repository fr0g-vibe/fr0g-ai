package mastercontrol

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// MCPConfig holds the Master Control Program configuration
type MCPConfig struct {
	// Core settings
	LearningEnabled         bool `yaml:"learning_enabled"`
	SystemConsciousness     bool `yaml:"system_consciousness"`
	EmergentCapabilities    bool `yaml:"emergent_capabilities"`
	MaxConcurrentWorkflows  int  `yaml:"max_concurrent_workflows"`
	
	// Legacy settings for compatibility
	AdaptationThreshold float64       `yaml:"adaptation_threshold"`
	MemoryRetention     time.Duration `yaml:"memory_retention"`
	HealthCheckInterval time.Duration `yaml:"health_check_interval"`
	MetricsInterval     time.Duration `yaml:"metrics_interval"`
	ResourceOptimization bool         `yaml:"resource_optimization"`
	PredictiveManagement bool         `yaml:"predictive_management"`
}

// DefaultMCPConfig returns a default configuration for the Master Control Program
func DefaultMCPConfig() *MCPConfig {
	return &MCPConfig{
		// Intelligence settings
		LearningEnabled:     true,
		AdaptationThreshold: 0.7,
		MemoryRetention:     time.Hour * 24 * 30, // 30 days
		
		// Monitoring settings
		HealthCheckInterval: time.Second * 30,
		MetricsInterval:     time.Minute * 5,
		
		// Orchestration settings
		MaxConcurrentWorkflows: 10,
		ResourceOptimization:   true,
		PredictiveManagement:   true,
		
		// System settings
		SystemConsciousness:  true,
		EmergentCapabilities: true,
	}
}

// LoadMCPConfig loads MCP configuration from various sources
func LoadMCPConfig(configPath string) (*MCPConfig, error) {
	config := DefaultMCPConfig()
	
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			data, err := os.ReadFile(configPath)
			if err != nil {
				return nil, err
			}
			
			if err := yaml.Unmarshal(data, config); err != nil {
				return nil, err
			}
		}
	}
	
	return config, nil
}
