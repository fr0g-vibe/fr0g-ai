package mastercontrol

import "time"

// MCPConfig represents the Master Control Program configuration
type MCPConfig struct {
	Input InputConfig `yaml:"input"`
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
