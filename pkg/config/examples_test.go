package config_test

import (
	"testing"
	"time"

	sharedconfig "pkg/config"
)

// Example of how to integrate shared config in a subproject
type ExampleConfig struct {
	// Embed shared configurations
	Server     sharedconfig.ServerConfig     `yaml:"server"`
	Security   sharedconfig.SecurityConfig   `yaml:"security"`
	Monitoring sharedconfig.MonitoringConfig `yaml:"monitoring"`
	Storage    sharedconfig.StorageConfig    `yaml:"storage"`
	
	// Project-specific fields
	ProjectName    string        `yaml:"project_name"`
	MaxWorkers     int           `yaml:"max_workers"`
	ProcessTimeout time.Duration `yaml:"process_timeout"`
	Features       []string      `yaml:"features"`
}

// Validate implements validation for the example config
func (c *ExampleConfig) Validate() sharedconfig.ValidationErrors {
	var errors []sharedconfig.ValidationError
	
	// Validate shared configurations
	errors = append(errors, c.Server.Validate()...)
	errors = append(errors, c.Security.Validate()...)
	errors = append(errors, c.Monitoring.Validate()...)
	errors = append(errors, c.Storage.Validate()...)
	
	// Validate project-specific fields
	if err := sharedconfig.ValidateRequired(c.ProjectName, "project_name"); err != nil {
		errors = append(errors, *err)
	}
	
	if err := sharedconfig.ValidatePositive(c.MaxWorkers, "max_workers"); err != nil {
		errors = append(errors, *err)
	}
	
	if c.ProcessTimeout > 0 {
		if err := sharedconfig.ValidateTimeout(c.ProcessTimeout, "process_timeout"); err != nil {
			errors = append(errors, *err)
		}
	}
	
	if err := sharedconfig.ValidateStringSliceNotEmpty(c.Features, "features"); err != nil {
		errors = append(errors, *err)
	}
	
	return sharedconfig.ValidationErrors(errors)
}

// LoadExampleConfig demonstrates how to load configuration
func LoadExampleConfig(configPath string) (*ExampleConfig, error) {
	// Create loader
	loader := sharedconfig.NewLoader(sharedconfig.LoaderOptions{
		ConfigPath: configPath,
		EnvPrefix:  "EXAMPLE",
		EnvFilePaths: []string{
			".env",
			"../example.env",
		},
	})
	
	// Load environment files
	if err := loader.LoadEnvFiles(); err != nil {
		return nil, err
	}
	
	// Create config with defaults
	config := &ExampleConfig{
		Server: sharedconfig.ServerConfig{
			Host:            "0.0.0.0",
			HTTPPort:        8080,
			GRPCPort:        9090,
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			ShutdownTimeout: 5 * time.Second,
		},
		Security: sharedconfig.SecurityConfig{
			EnableCORS:     true,
			AllowedOrigins: []string{"*"},
			RateLimitRPM:   60,
		},
		Monitoring: sharedconfig.MonitoringConfig{
			EnableMetrics:       true,
			MetricsPort:         8082,
			HealthCheckInterval: 30,
		},
		Storage: sharedconfig.StorageConfig{
			Type:    "memory",
			DataDir: "./data",
		},
		ProjectName:    "example-service",
		MaxWorkers:     10,
		ProcessTimeout: 5 * time.Minute,
		Features:       []string{"feature1", "feature2"},
	}
	
	// Load from file
	if err := loader.LoadFromFile(config); err != nil {
		return nil, err
	}
	
	// Override with environment variables
	config.ProjectName = loader.GetEnvString("PROJECT_NAME", config.ProjectName)
	config.MaxWorkers = loader.GetEnvInt("MAX_WORKERS", config.MaxWorkers)
	config.ProcessTimeout = loader.GetEnvDuration("PROCESS_TIMEOUT", config.ProcessTimeout)
	
	// Validate final configuration
	if errors := config.Validate(); errors.HasErrors() {
		return nil, errors
	}
	
	return config, nil
}

func TestExampleConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  ExampleConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: ExampleConfig{
				Server: sharedconfig.ServerConfig{
					HTTPPort: 8080,
					GRPCPort: 9090,
				},
				Storage: sharedconfig.StorageConfig{
					Type: "memory",
				},
				ProjectName: "test-project",
				MaxWorkers:  5,
				Features:    []string{"feature1"},
			},
			wantErr: false,
		},
		{
			name: "invalid port",
			config: ExampleConfig{
				Server: sharedconfig.ServerConfig{
					HTTPPort: -1,
					GRPCPort: 9090,
				},
				Storage: sharedconfig.StorageConfig{
					Type: "memory",
				},
				ProjectName: "test-project",
				MaxWorkers:  5,
				Features:    []string{"feature1"},
			},
			wantErr: true,
			errMsg:  "port must be between 1 and 65535",
		},
		{
			name: "missing project name",
			config: ExampleConfig{
				Server: sharedconfig.ServerConfig{
					HTTPPort: 8080,
					GRPCPort: 9090,
				},
				Storage: sharedconfig.StorageConfig{
					Type: "memory",
				},
				MaxWorkers: 5,
				Features:   []string{"feature1"},
			},
			wantErr: true,
			errMsg:  "project_name: field is required",
		},
		{
			name: "invalid max workers",
			config: ExampleConfig{
				Server: sharedconfig.ServerConfig{
					HTTPPort: 8080,
					GRPCPort: 9090,
				},
				Storage: sharedconfig.StorageConfig{
					Type: "memory",
				},
				ProjectName: "test-project",
				MaxWorkers:  -1,
				Features:    []string{"feature1"},
			},
			wantErr: true,
			errMsg:  "max_workers: value must be positive",
		},
		{
			name: "empty features",
			config: ExampleConfig{
				Server: sharedconfig.ServerConfig{
					HTTPPort: 8080,
					GRPCPort: 9090,
				},
				Storage: sharedconfig.StorageConfig{
					Type: "memory",
				},
				ProjectName: "test-project",
				MaxWorkers:  5,
				Features:    []string{},
			},
			wantErr: true,
			errMsg:  "features: at least one item is required",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := tt.config.Validate()
			
			if tt.wantErr {
				if !errors.HasErrors() {
					t.Errorf("expected validation errors, got none")
					return
				}
				
				if tt.errMsg != "" && !contains(errors.Error(), tt.errMsg) {
					t.Errorf("expected error message to contain %q, got %q", tt.errMsg, errors.Error())
				}
			} else {
				if errors.HasErrors() {
					t.Errorf("unexpected validation errors: %v", errors)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())))
}
