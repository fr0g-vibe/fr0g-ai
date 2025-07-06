package cache

import (
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

// CacheEntry represents a cached service discovery result
type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

// LRUCache provides thread-safe LRU caching with TTL
type LRUCache struct {
	cache *lru.Cache[string, *CacheEntry]
	mu    sync.RWMutex
	ttl   time.Duration
}

// NewLRUCache creates a new LRU cache with specified size and TTL
func NewLRUCache(size int, ttl time.Duration) (*LRUCache, error) {
	cache, err := lru.New[string, *CacheEntry](size)
	if err != nil {
		return nil, err
	}

	lruCache := &LRUCache{
		cache: cache,
		ttl:   ttl,
	}

	// Start cleanup goroutine
	go lruCache.cleanup()

	return lruCache, nil
}

// Get retrieves a value from cache if not expired
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache.Get(key)
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		c.cache.Remove(key)
		return nil, false
	}

	return entry.Data, true
}

// Set stores a value in cache with TTL
func (c *LRUCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := &CacheEntry{
		Data:      value,
		ExpiresAt: time.Now().Add(c.ttl),
	}

	c.cache.Add(key, entry)
}

// Delete removes a key from cache
func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Remove(key)
}

// Clear removes all entries from cache
func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Purge()
}

// cleanup removes expired entries periodically
func (c *LRUCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		keys := c.cache.Keys()
		for _, key := range keys {
			if entry, exists := c.cache.Peek(key); exists {
				if time.Now().After(entry.ExpiresAt) {
					c.cache.Remove(key)
				}
			}
		}
		c.mu.Unlock()
	}
}
