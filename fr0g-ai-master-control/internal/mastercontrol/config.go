package mastercontrol

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
	sharedconfig "pkg/config"
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

// Validate validates the MCP configuration
func (c *MCPConfig) Validate() sharedconfig.ValidationErrors {
	var errors []sharedconfig.ValidationError
	
	if err := sharedconfig.ValidatePositive(c.MaxConcurrentWorkflows, "max_concurrent_workflows"); err != nil {
		errors = append(errors, *err)
	}
	
	if err := sharedconfig.ValidateRange(c.AdaptationThreshold, 0, 1, "adaptation_threshold"); err != nil {
		errors = append(errors, *err)
	}
	
	if err := sharedconfig.ValidateTimeout(c.MemoryRetention, "memory_retention"); err != nil {
		errors = append(errors, *err)
	}
	
	if err := sharedconfig.ValidateTimeout(c.HealthCheckInterval, "health_check_interval"); err != nil {
		errors = append(errors, *err)
	}
	
	if err := sharedconfig.ValidateTimeout(c.MetricsInterval, "metrics_interval"); err != nil {
		errors = append(errors, *err)
	}
	
	return sharedconfig.ValidationErrors(errors)
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
	
	// Validate the loaded configuration
	if errors := config.Validate(); errors.HasErrors() {
		return nil, errors
	}
	
	return config, nil
}
