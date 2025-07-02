package memory

import (
	"context"
	"log"
	"sync"
)

// MemoryManager handles different types of memory storage
type MemoryManager struct {
	config       *MemoryConfig
	shortTerm    *MemoryStore
	longTerm     *MemoryStore
	episodic     *MemoryStore
	semantic     *MemoryStore
	mu           sync.RWMutex
}

// MemoryConfig holds memory configuration
type MemoryConfig struct {
	ShortTermCapacity int
	LongTermCapacity  int
	EpisodicCapacity  int
	SemanticCapacity  int
}

// MemoryStore represents a memory storage unit
type MemoryStore struct {
	capacity int
	items    []interface{}
	mu       sync.RWMutex
}

// NewMemoryManager creates a new memory manager
func NewMemoryManager(config *MCPConfig) *MemoryManager {
	return &MemoryManager{
		shortTerm: NewMemoryStore(1000),
		longTerm:  NewMemoryStore(10000),
		episodic:  NewMemoryStore(5000),
		semantic:  NewMemoryStore(15000),
	}
}

// NewMemoryStore creates a new memory store
func NewMemoryStore(capacity int) *MemoryStore {
	return &MemoryStore{
		capacity: capacity,
		items:    make([]interface{}, 0, capacity),
	}
}

// Start begins memory management
func (mm *MemoryManager) Start() error {
	log.Println("Memory Manager: Starting memory management processes...")
	log.Println("Memory Manager: Memory management processes started")
	return nil
}

// Stop gracefully stops memory management
func (mm *MemoryManager) Stop() error {
	return nil
}

// Store stores an item in the specified memory type
func (mm *MemoryManager) Store(memoryType string, item interface{}) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	switch memoryType {
	case "short_term":
		mm.shortTerm.Store(item)
	case "long_term":
		mm.longTerm.Store(item)
	case "episodic":
		mm.episodic.Store(item)
	case "semantic":
		mm.semantic.Store(item)
	}
}

// Store stores an item in the memory store
func (ms *MemoryStore) Store(item interface{}) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	
	if len(ms.items) >= ms.capacity {
		// Remove oldest item
		ms.items = ms.items[1:]
	}
	
	ms.items = append(ms.items, item)
}
