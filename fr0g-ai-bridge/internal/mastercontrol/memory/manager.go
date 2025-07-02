package memory

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// MemoryManager handles all memory operations for the MCP
type MemoryManager struct {
	// Memory stores
	shortTerm  map[string]*MemoryEntry
	longTerm   map[string]*MemoryEntry
	episodic   []*EpisodicMemory
	semantic   map[string]*SemanticMemory
	
	// Memory management
	mu         sync.RWMutex
	config     *MemoryConfig
	
	// Statistics
	stats      *MemoryStats
}

// MemoryEntry represents a single memory entry
type MemoryEntry struct {
	Key         string                 `json:"key"`
	Value       interface{}            `json:"value"`
	Type        string                 `json:"type"`
	Importance  float64               `json:"importance"`
	AccessCount int                   `json:"access_count"`
	CreatedAt   time.Time             `json:"created_at"`
	LastAccess  time.Time             `json:"last_access"`
	ExpiresAt   *time.Time            `json:"expires_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// EpisodicMemory represents memory of specific events/episodes
type EpisodicMemory struct {
	ID          string                 `json:"id"`
	Event       string                 `json:"event"`
	Context     map[string]interface{} `json:"context"`
	Participants []string              `json:"participants"`
	Outcome     string                 `json:"outcome"`
	Importance  float64               `json:"importance"`
	Timestamp   time.Time             `json:"timestamp"`
	Tags        []string              `json:"tags"`
}

// SemanticMemory represents conceptual knowledge
type SemanticMemory struct {
	Concept     string                 `json:"concept"`
	Definition  string                 `json:"definition"`
	Relations   map[string]float64     `json:"relations"`
	Examples    []string               `json:"examples"`
	Confidence  float64               `json:"confidence"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// MemoryStats tracks memory usage statistics
type MemoryStats struct {
	ShortTermCount    int       `json:"short_term_count"`
	LongTermCount     int       `json:"long_term_count"`
	EpisodicCount     int       `json:"episodic_count"`
	SemanticCount     int       `json:"semantic_count"`
	TotalMemoryUsage  int64     `json:"total_memory_usage"`
	LastCleanup       time.Time `json:"last_cleanup"`
	AccessesPerSecond float64   `json:"accesses_per_second"`
}

// MemoryConfig holds memory management configuration
type MemoryConfig struct {
	ShortTermTTL        time.Duration `yaml:"short_term_ttl"`
	LongTermTTL         time.Duration `yaml:"long_term_ttl"`
	MaxShortTermEntries int           `yaml:"max_short_term_entries"`
	MaxLongTermEntries  int           `yaml:"max_long_term_entries"`
	MaxEpisodicMemories int           `yaml:"max_episodic_memories"`
	MaxSemanticMemories int           `yaml:"max_semantic_memories"`
	CleanupInterval     time.Duration `yaml:"cleanup_interval"`
	ImportanceThreshold float64       `yaml:"importance_threshold"`
	CompressionEnabled  bool          `yaml:"compression_enabled"`
}

// NewMemoryManager creates a new memory manager
func NewMemoryManager(config *MemoryConfig) *MemoryManager {
	mm := &MemoryManager{
		shortTerm: make(map[string]*MemoryEntry),
		longTerm:  make(map[string]*MemoryEntry),
		episodic:  make([]*EpisodicMemory, 0),
		semantic:  make(map[string]*SemanticMemory),
		config:    config,
		stats: &MemoryStats{
			LastCleanup: time.Now(),
		},
	}
	
	return mm
}

// Start begins memory management operations
func (mm *MemoryManager) Start() error {
	log.Println("Memory Manager: Starting memory management processes...")
	
	// Start cleanup routine
	go mm.cleanupLoop()
	
	// Start statistics update routine
	go mm.statsLoop()
	
	log.Println("Memory Manager: Memory management processes started")
	return nil
}

// Stop gracefully stops the memory manager
func (mm *MemoryManager) Stop() error {
	log.Println("Memory Manager: Stopping memory management...")
	return nil
}

// Store stores a value in memory with the specified type
func (mm *MemoryManager) Store(key string, value interface{}) error {
	return mm.StoreWithType(key, value, "short_term")
}

// StoreWithType stores a value in memory with a specific memory type
func (mm *MemoryManager) StoreWithType(key string, value interface{}, memoryType string) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	entry := &MemoryEntry{
		Key:         key,
		Value:       value,
		Type:        memoryType,
		Importance:  mm.calculateImportance(value),
		AccessCount: 0,
		CreatedAt:   time.Now(),
		LastAccess:  time.Now(),
		Metadata:    make(map[string]interface{}),
	}
	
	// Set expiration based on type
	switch memoryType {
	case "short_term":
		expiresAt := time.Now().Add(mm.config.ShortTermTTL)
		entry.ExpiresAt = &expiresAt
		mm.shortTerm[key] = entry
	case "long_term":
		if mm.config.LongTermTTL > 0 {
			expiresAt := time.Now().Add(mm.config.LongTermTTL)
			entry.ExpiresAt = &expiresAt
		}
		mm.longTerm[key] = entry
	default:
		return fmt.Errorf("unknown memory type: %s", memoryType)
	}
	
	mm.updateStats()
	return nil
}

// Retrieve retrieves a value from memory
func (mm *MemoryManager) Retrieve(key string) (interface{}, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	// Check short-term memory first
	if entry, exists := mm.shortTerm[key]; exists {
		if !mm.isExpired(entry) {
			entry.AccessCount++
			entry.LastAccess = time.Now()
			return entry.Value, nil
		}
		// Remove expired entry
		delete(mm.shortTerm, key)
	}
	
	// Check long-term memory
	if entry, exists := mm.longTerm[key]; exists {
		if !mm.isExpired(entry) {
			entry.AccessCount++
			entry.LastAccess = time.Now()
			return entry.Value, nil
		}
		// Remove expired entry
		delete(mm.longTerm, key)
	}
	
	return nil, fmt.Errorf("key not found: %s", key)
}

// StoreEpisode stores an episodic memory
func (mm *MemoryManager) StoreEpisode(episode *EpisodicMemory) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	mm.episodic = append(mm.episodic, episode)
	
	// Limit episodic memories
	if len(mm.episodic) > mm.config.MaxEpisodicMemories {
		// Remove least important episodes
		mm.episodic = mm.removeOldestEpisodes(mm.episodic, mm.config.MaxEpisodicMemories)
	}
	
	mm.updateStats()
	return nil
}

// StoreConcept stores semantic memory
func (mm *MemoryManager) StoreConcept(concept string, memory *SemanticMemory) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	memory.UpdatedAt = time.Now()
	if memory.CreatedAt.IsZero() {
		memory.CreatedAt = time.Now()
	}
	
	mm.semantic[concept] = memory
	
	// Limit semantic memories
	if len(mm.semantic) > mm.config.MaxSemanticMemories {
		mm.cleanupSemanticMemory()
	}
	
	mm.updateStats()
	return nil
}

// GetEpisodes retrieves episodic memories matching criteria
func (mm *MemoryManager) GetEpisodes(tags []string, limit int) []*EpisodicMemory {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	var matching []*EpisodicMemory
	
	for _, episode := range mm.episodic {
		if mm.episodeMatchesTags(episode, tags) {
			matching = append(matching, episode)
		}
		
		if len(matching) >= limit {
			break
		}
	}
	
	return matching
}

// GetConcept retrieves semantic memory for a concept
func (mm *MemoryManager) GetConcept(concept string) (*SemanticMemory, error) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	if memory, exists := mm.semantic[concept]; exists {
		return memory, nil
	}
	
	return nil, fmt.Errorf("concept not found: %s", concept)
}

// GetPatterns returns patterns from memory (implements interface)
func (mm *MemoryManager) GetPatterns() []interface{} {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	var patterns []interface{}
	
	// Extract patterns from episodic memory
	for _, episode := range mm.episodic {
		if episode.Importance > mm.config.ImportanceThreshold {
			patterns = append(patterns, episode)
		}
	}
	
	return patterns
}

// GetStats returns memory usage statistics
func (mm *MemoryManager) GetStats() *MemoryStats {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	stats := *mm.stats
	return &stats
}

// PromoteToLongTerm promotes a short-term memory to long-term
func (mm *MemoryManager) PromoteToLongTerm(key string) error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	if entry, exists := mm.shortTerm[key]; exists {
		// Remove from short-term
		delete(mm.shortTerm, key)
		
		// Add to long-term
		entry.Type = "long_term"
		if mm.config.LongTermTTL > 0 {
			expiresAt := time.Now().Add(mm.config.LongTermTTL)
			entry.ExpiresAt = &expiresAt
		} else {
			entry.ExpiresAt = nil
		}
		
		mm.longTerm[key] = entry
		mm.updateStats()
		
		return nil
	}
	
	return fmt.Errorf("key not found in short-term memory: %s", key)
}

// Private methods

func (mm *MemoryManager) cleanupLoop() {
	ticker := time.NewTicker(mm.config.CleanupInterval)
	defer ticker.Stop()
	
	for range ticker.C {
		mm.cleanup()
	}
}

func (mm *MemoryManager) statsLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		mm.updateStats()
	}
}

func (mm *MemoryManager) cleanup() {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	now := time.Now()
	
	// Cleanup short-term memory
	for key, entry := range mm.shortTerm {
		if mm.isExpired(entry) {
			delete(mm.shortTerm, key)
		}
	}
	
	// Cleanup long-term memory
	for key, entry := range mm.longTerm {
		if mm.isExpired(entry) {
			delete(mm.longTerm, key)
		}
	}
	
	// Cleanup episodic memory based on importance
	if len(mm.episodic) > mm.config.MaxEpisodicMemories {
		mm.episodic = mm.removeOldestEpisodes(mm.episodic, mm.config.MaxEpisodicMemories)
	}
	
	mm.stats.LastCleanup = now
	mm.updateStats()
	
	log.Printf("Memory Manager: Cleanup completed - ST:%d LT:%d EP:%d SM:%d", 
		len(mm.shortTerm), len(mm.longTerm), len(mm.episodic), len(mm.semantic))
}

func (mm *MemoryManager) isExpired(entry *MemoryEntry) bool {
	if entry.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*entry.ExpiresAt)
}

func (mm *MemoryManager) calculateImportance(value interface{}) float64 {
	// Simple importance calculation based on data size and type
	// In a real implementation, this would be more sophisticated
	
	data, _ := json.Marshal(value)
	size := len(data)
	
	// Base importance on size (larger data might be more important)
	importance := float64(size) / 1000.0
	
	// Cap at 1.0
	if importance > 1.0 {
		importance = 1.0
	}
	
	return importance
}

func (mm *MemoryManager) updateStats() {
	mm.stats.ShortTermCount = len(mm.shortTerm)
	mm.stats.LongTermCount = len(mm.longTerm)
	mm.stats.EpisodicCount = len(mm.episodic)
	mm.stats.SemanticCount = len(mm.semantic)
	
	// Calculate total memory usage (simplified)
	mm.stats.TotalMemoryUsage = int64(mm.stats.ShortTermCount + mm.stats.LongTermCount + 
		mm.stats.EpisodicCount + mm.stats.SemanticCount)
}

func (mm *MemoryManager) removeOldestEpisodes(episodes []*EpisodicMemory, maxCount int) []*EpisodicMemory {
	if len(episodes) <= maxCount {
		return episodes
	}
	
	// Sort by importance (descending) and timestamp (descending)
	// Keep the most important and recent episodes
	// This is a simplified implementation
	
	return episodes[len(episodes)-maxCount:]
}

func (mm *MemoryManager) cleanupSemanticMemory() {
	// Remove least confident concepts
	// This is a simplified implementation
	
	if len(mm.semantic) <= mm.config.MaxSemanticMemories {
		return
	}
	
	// In a real implementation, we would sort by confidence and remove the least confident
	// For now, just remove one random concept
	for concept := range mm.semantic {
		delete(mm.semantic, concept)
		break
	}
}

func (mm *MemoryManager) episodeMatchesTags(episode *EpisodicMemory, tags []string) bool {
	if len(tags) == 0 {
		return true
	}
	
	for _, tag := range tags {
		for _, episodeTag := range episode.Tags {
			if tag == episodeTag {
				return true
			}
		}
	}
	
	return false
}
