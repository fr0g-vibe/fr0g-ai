# Shared Configuration Library

This package provides a standardized configuration system for all fr0g-ai subprojects.

## Features

- **Unified validation** - Consistent validation across all projects
- **Common configuration types** - Reusable config structures for server, security, storage, etc.
- **Flexible loading** - Support for YAML files, environment variables, and .env files
- **Proper error handling** - ValidationErrors type for collecting multiple validation errors
- **Easy integration** - Simple migration path from existing configurations

## Usage

### Basic Integration

```go
package config

import (
    sharedconfig "pkg/config"
)

// Your project-specific config
type Config struct {
    // Embed shared configs
    Server     sharedconfig.ServerConfig     `yaml:"server"`
    Security   sharedconfig.SecurityConfig   `yaml:"security"`
    Monitoring sharedconfig.MonitoringConfig `yaml:"monitoring"`
    
    // Project-specific fields
    MyCustomField string `yaml:"my_custom_field"`
}

// Implement validation
func (c *Config) Validate() sharedconfig.ValidationErrors {
    var errors []sharedconfig.ValidationError
    
    // Validate shared configs
    errors = append(errors, c.Server.Validate()...)
    errors = append(errors, c.Security.Validate()...)
    errors = append(errors, c.Monitoring.Validate()...)
    
    // Add custom validation
    if err := sharedconfig.ValidateRequired(c.MyCustomField, "my_custom_field"); err != nil {
        errors = append(errors, *err)
    }
    
    return sharedconfig.ValidationErrors(errors)
}

// Load configuration
func LoadConfig(configPath string) (*Config, error) {
    loader := sharedconfig.NewLoader(sharedconfig.LoaderOptions{
        ConfigPath: configPath,
        EnvPrefix:  "MYAPP",
    })
    
    // Load .env files
    if err := loader.LoadEnvFiles(); err != nil {
        return nil, err
    }
    
    // Create config with defaults
    config := &Config{
        Server: sharedconfig.ServerConfig{
            Host:     "0.0.0.0",
            HTTPPort: 8080,
            GRPCPort: 9090,
        },
        MyCustomField: "default_value",
    }
    
    // Load from file
    if err := loader.LoadFromFile(config); err != nil {
        return nil, err
    }
    
    // Override with environment variables
    config.MyCustomField = loader.GetEnvString("MY_CUSTOM_FIELD", config.MyCustomField)
    
    // Validate
    if errors := config.Validate(); errors.HasErrors() {
        return nil, errors
    }
    
    return config, nil
}
```

### Available Shared Types

#### ServerConfig
```yaml
server:
  host: "0.0.0.0"
  http_port: 8080
  grpc_port: 9090
  enable_tls: false
  cert_file: ""
  key_file: ""
  read_timeout: "30s"
  write_timeout: "30s"
  shutdown_timeout: "5s"
```

#### SecurityConfig
```yaml
security:
  enable_auth: false
  api_key: ""
  allowed_api_keys: []
  enable_cors: true
  allowed_origins: ["*"]
  rate_limit_rpm: 60
  require_api_key: false
```

#### MonitoringConfig
```yaml
monitoring:
  enable_metrics: true
  metrics_port: 8082
  health_check_interval: 30
  enable_tracing: false
```

#### StorageConfig
```yaml
storage:
  type: "memory"  # memory, file, redis, postgres
  data_dir: "./data"
```

### Validation Functions

The library provides many validation functions:

- `ValidatePort(port, fieldName)` - Validates port numbers
- `ValidateTimeout(duration, fieldName)` - Validates timeout durations
- `ValidateRequired(value, fieldName)` - Validates required fields
- `ValidateAPIKey(key, fieldName)` - Validates API key strength
- `ValidateModel(model)` - Validates model names
- `ValidateRole(role)` - Validates chat message roles
- `ValidateURL(url, fieldName)` - Validates URL format
- `ValidateRange(value, min, max, fieldName)` - Validates numeric ranges
- `ValidatePositive(value, fieldName)` - Validates positive numbers
- `ValidateNetworkAddress(address, fieldName)` - Validates network addresses

### Environment Variable Loading

The loader supports automatic environment variable loading with prefixes:

```go
loader := sharedconfig.NewLoader(sharedconfig.LoaderOptions{
    ConfigPath: "config.yaml",
    EnvPrefix:  "MYAPP",  // Will look for MYAPP_HTTP_PORT, etc.
})

// Get values with fallbacks
httpPort := loader.GetEnvInt("HTTP_PORT", 8080)
debug := loader.GetEnvBool("DEBUG", false)
timeout := loader.GetEnvDuration("TIMEOUT", 30*time.Second)
```

### Migration Guide

#### From fr0g-ai-aip
Replace custom validation with shared types:
```go
// Old
type ValidationError struct { ... }

// New
import sharedconfig "pkg/config"
type ValidationError = sharedconfig.ValidationError
```

#### From fr0g-ai-bridge
Use shared SecurityConfig:
```go
// Old
type SecurityConfig struct {
    EnableAuth bool `yaml:"enable_auth"`
    // ...
}

// New
type SecurityConfig struct {
    sharedconfig.SecurityConfig `yaml:",inline"`
    EnableReflection bool `yaml:"enable_reflection"`  // Project-specific
}
```

#### From fr0g-ai-master-control
Add validation methods:
```go
// Add to existing config
func (c *MCPConfig) Validate() sharedconfig.ValidationErrors {
    var errors []sharedconfig.ValidationError
    
    if err := sharedconfig.ValidatePositive(c.MaxConcurrentWorkflows, "max_concurrent_workflows"); err != nil {
        errors = append(errors, *err)
    }
    
    return sharedconfig.ValidationErrors(errors)
}
```

## Best Practices

1. **Always validate** - Call `Validate()` after loading configuration
2. **Use shared types** - Embed common configs instead of duplicating
3. **Provide defaults** - Set sensible defaults before loading from files
4. **Handle errors** - Check for validation errors and handle appropriately
5. **Document custom fields** - Add comments for project-specific configuration
6. **Use environment overrides** - Allow environment variables to override file settings

## Testing

```go
func TestConfigValidation(t *testing.T) {
    config := &Config{
        Server: sharedconfig.ServerConfig{
            HTTPPort: -1,  // Invalid
        },
    }
    
    errors := config.Validate()
    assert.True(t, errors.HasErrors())
    assert.Contains(t, errors.Error(), "port must be between 1 and 65535")
}
```
