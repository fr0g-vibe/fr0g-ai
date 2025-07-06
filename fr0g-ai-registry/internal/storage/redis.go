package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStorage provides Redis-backed persistent storage
type RedisStorage struct {
	client *redis.Client
	prefix string
}

// ServiceInfo represents a registered service (duplicate from main for storage)
type ServiceInfo struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Tags     []string          `json:"tags,omitempty"`
	Meta     map[string]string `json:"meta,omitempty"`
	Health   string            `json:"health"`
	LastSeen time.Time         `json:"last_seen"`
}

// NewRedisStorage creates a new Redis storage backend
func NewRedisStorage(addr, password string, db int) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     20,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	return &RedisStorage{
		client: client,
		prefix: "fr0g:registry:service:",
	}
}

// Ping checks Redis connectivity
func (r *RedisStorage) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// SaveService persists a service to Redis
func (r *RedisStorage) SaveService(ctx context.Context, service *ServiceInfo) error {
	data, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service: %w", err)
	}

	key := r.prefix + service.ID
	return r.client.Set(ctx, key, data, 0).Err()
}

// GetService retrieves a service from Redis
func (r *RedisStorage) GetService(ctx context.Context, serviceID string) (*ServiceInfo, error) {
	key := r.prefix + serviceID
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Service not found
		}
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	var service ServiceInfo
	if err := json.Unmarshal([]byte(data), &service); err != nil {
		return nil, fmt.Errorf("failed to unmarshal service: %w", err)
	}

	return &service, nil
}

// GetAllServices retrieves all services from Redis
func (r *RedisStorage) GetAllServices(ctx context.Context) (map[string]*ServiceInfo, error) {
	pattern := r.prefix + "*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get service keys: %w", err)
	}

	services := make(map[string]*ServiceInfo)
	if len(keys) == 0 {
		return services, nil
	}

	// Use pipeline for better performance
	pipe := r.client.Pipeline()
	cmds := make(map[string]*redis.StringCmd)
	
	for _, key := range keys {
		cmds[key] = pipe.Get(ctx, key)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute pipeline: %w", err)
	}

	for _, cmd := range cmds {
		data, err := cmd.Result()
		if err != nil {
			continue // Skip failed entries
		}

		var service ServiceInfo
		if err := json.Unmarshal([]byte(data), &service); err != nil {
			continue // Skip invalid entries
		}

		services[service.ID] = &service
	}

	return services, nil
}

// DeleteService removes a service from Redis
func (r *RedisStorage) DeleteService(ctx context.Context, serviceID string) error {
	key := r.prefix + serviceID
	return r.client.Del(ctx, key).Err()
}

// Close closes the Redis connection
func (r *RedisStorage) Close() error {
	return r.client.Close()
}
