package test

// ServiceInfo represents a service registration
type ServiceInfo struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Address string            `json:"address"`
	Port    int               `json:"port"`
	Tags    []string          `json:"tags,omitempty"`
	Meta    map[string]string `json:"meta,omitempty"`
	Check   *HealthCheck      `json:"check,omitempty"`
}

// HealthCheck represents a health check configuration
type HealthCheck struct {
	HTTP     string `json:"http,omitempty"`
	Interval string `json:"interval,omitempty"`
	Timeout  string `json:"timeout,omitempty"`
}

// ServiceDetail represents the detailed service information returned by the registry
type ServiceDetail struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Tags     []string          `json:"tags"`
	Meta     map[string]string `json:"meta"`
	Health   string            `json:"health"`
	LastSeen string            `json:"last_seen"`
}
