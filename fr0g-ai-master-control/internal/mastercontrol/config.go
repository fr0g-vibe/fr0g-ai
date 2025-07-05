package mastercontrol

import (
	"os"
	"time"
	"gopkg.in/yaml.v2"
)

// MCPConfig represents the Master Control Program configuration
type MCPConfig struct {
	Input                    InputConfig   `yaml:"input"`
	LearningEnabled          bool          `yaml:"learning_enabled"`
	SystemConsciousness      bool          `yaml:"system_consciousness"`
	EmergentCapabilities     bool          `yaml:"emergent_capabilities"`
	MaxConcurrentWorkflows   int           `yaml:"max_concurrent_workflows"`
	CognitiveReflectionRate  time.Duration `yaml:"cognitive_reflection_rate"`
	PatternRecognitionDepth  int           `yaml:"pattern_recognition_depth"`
	AdaptiveLearningRate     float64       `yaml:"adaptive_learning_rate"`
}

// InputConfig represents input system configuration
type InputConfig struct {
	Webhook WebhookConfig `yaml:"webhook"`
}

// WebhookConfig represents webhook input configuration
type WebhookConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// DefaultMCPConfig returns a default MCP configuration
func DefaultMCPConfig() *MCPConfig {
	return &MCPConfig{
		LearningEnabled:          true,
		SystemConsciousness:      true,
		EmergentCapabilities:     true,
		MaxConcurrentWorkflows:   10,
		CognitiveReflectionRate:  30 * time.Second,
		PatternRecognitionDepth:  5,
		AdaptiveLearningRate:     0.154,
		Input: InputConfig{
			Webhook: WebhookConfig{
				Host:         "localhost",
				Port:         8080,
				ReadTimeout:  30 * time.Second,
				WriteTimeout: 30 * time.Second,
			},
		},
	}
}

// LoadMCPConfig loads the MCP configuration from file or returns defaults
func LoadMCPConfig(configPath string) (*MCPConfig, error) {
	config := DefaultMCPConfig()
	
	// If no config path specified, return defaults
	if configPath == "" {
		return config, nil
	}
	
	// Try to load from file
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// File doesn't exist, return defaults
		return config, nil
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	
	if err := yaml.Unmarshal(data, config); err != nil {
		return config, err
	}
	
	return config, nil
}
