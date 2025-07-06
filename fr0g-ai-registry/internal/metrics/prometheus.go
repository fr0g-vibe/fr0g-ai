package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Service registration metrics
	ServiceRegistrations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fr0g_registry_service_registrations_total",
			Help: "Total number of service registrations",
		},
		[]string{"service_name"},
	)

	ServiceDeregistrations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fr0g_registry_service_deregistrations_total",
			Help: "Total number of service deregistrations",
		},
		[]string{"service_name"},
	)

	// Discovery metrics
	DiscoveryRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fr0g_registry_discovery_requests_total",
			Help: "Total number of service discovery requests",
		},
		[]string{"endpoint"},
	)

	DiscoveryLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "fr0g_registry_discovery_duration_seconds",
			Help:    "Service discovery request duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)

	// Cache metrics
	CacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "fr0g_registry_cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	CacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "fr0g_registry_cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	// Service health metrics
	ServicesTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fr0g_registry_services_total",
			Help: "Total number of registered services",
		},
		[]string{"health_status"},
	)

	// Redis metrics
	RedisOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fr0g_registry_redis_operations_total",
			Help: "Total number of Redis operations",
		},
		[]string{"operation", "status"},
	)
)
