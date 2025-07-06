package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-master-control/internal/mastercontrol"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-master-control/internal/mastercontrol/memory"
)

func main() {
	fmt.Println("🧠 Memory Manager Demo")
	fmt.Println("=====================")
	
	// Create MCP config
	config := mastercontrol.DefaultMCPConfig()
	
	// Create memory manager
	memManager := mastercontrol.NewMemoryManager(config)
	
	// Start memory manager
	if err := memManager.Start(); err != nil {
		log.Fatalf("Failed to start memory manager: %v", err)
	}
	
	fmt.Println("COMPLETED Memory Manager started successfully")
	fmt.Println()
	
	// Demonstrate memory operations
	demonstrateMemoryOperations(memManager)
	
	// Stop memory manager
	if err := memManager.Stop(); err != nil {
		log.Printf("Error stopping memory manager: %v", err)
	}
	
	fmt.Println("👋 Memory Manager demo complete")
}

func demonstrateMemoryOperations(mm *memory.MemoryManager) {
	fmt.Println("CHECKING Demonstrating Memory Operations:")
	fmt.Println("-----------------------------------")
	
	// Store some short-term memories
	fmt.Println("📝 Storing short-term memories...")
	mm.Store("user_preference", map[string]interface{}{
		"theme": "dark",
		"language": "en",
		"notifications": true,
	})
	
	mm.Store("session_data", map[string]interface{}{
		"user_id": "user123",
		"login_time": time.Now(),
		"ip_address": "192.168.1.100",
	})
	
	mm.Store("temp_calculation", 42.5)
	
	// Store some long-term memories
	fmt.Println("📚 Storing long-term memories...")
	mm.StoreWithType("system_config", map[string]interface{}{
		"version": "1.0.0",
		"deployment": "production",
		"features": []string{"ai", "learning", "consciousness"},
	}, "long_term")
	
	// Store episodic memory
	fmt.Println("📖 Storing episodic memory...")
	episode := &memory.EpisodicMemory{
		ID: "episode_001",
		Event: "User asked about weather",
		Context: map[string]interface{}{
			"location": "San Francisco",
			"time": time.Now(),
			"user_mood": "curious",
		},
		Participants: []string{"user123", "ai_assistant"},
		Outcome: "Provided weather information successfully",
		Importance: 0.7,
		Timestamp: time.Now(),
		Tags: []string{"weather", "query", "successful"},
	}
	mm.StoreEpisode(episode)
	
	// Store semantic memory
	fmt.Println("🧩 Storing semantic memory...")
	concept := &memory.SemanticMemory{
		Concept: "weather",
		Definition: "Atmospheric conditions including temperature, humidity, precipitation, and wind",
		Relations: map[string]float64{
			"temperature": 0.9,
			"humidity": 0.8,
			"precipitation": 0.85,
			"forecast": 0.7,
		},
		Examples: []string{"sunny", "rainy", "cloudy", "stormy"},
		Confidence: 0.95,
	}
	mm.StoreConcept("weather", concept)
	
	fmt.Println()
	
	// Retrieve memories
	fmt.Println("CHECKING Retrieving memories...")
	
	if userPref, err := mm.Retrieve("user_preference"); err == nil {
		fmt.Printf("COMPLETED Retrieved user preferences: %+v\n", userPref)
	}
	
	if sessionData, err := mm.Retrieve("session_data"); err == nil {
		fmt.Printf("COMPLETED Retrieved session data: %+v\n", sessionData)
	}
	
	if tempCalc, err := mm.Retrieve("temp_calculation"); err == nil {
		fmt.Printf("COMPLETED Retrieved calculation: %v\n", tempCalc)
	}
	
	if sysConfig, err := mm.Retrieve("system_config"); err == nil {
		fmt.Printf("COMPLETED Retrieved system config: %+v\n", sysConfig)
	}
	
	fmt.Println()
	
	// Get episodic memories
	fmt.Println("📖 Retrieving episodic memories...")
	episodes := mm.GetEpisodes([]string{"weather"}, 5)
	for _, ep := range episodes {
		fmt.Printf("COMPLETED Episode: %s - %s (Importance: %.2f)\n", ep.ID, ep.Event, ep.Importance)
	}
	
	// Get semantic memory
	fmt.Println("🧩 Retrieving semantic memory...")
	if weatherConcept, err := mm.GetConcept("weather"); err == nil {
		fmt.Printf("COMPLETED Concept 'weather': %s (Confidence: %.2f)\n", weatherConcept.Definition, weatherConcept.Confidence)
		fmt.Printf("   Relations: %+v\n", weatherConcept.Relations)
	}
	
	fmt.Println()
	
	// Show memory statistics
	fmt.Println("📊 Memory Statistics:")
	stats := mm.GetStats()
	fmt.Printf("   - Short-term memories: %d\n", stats.ShortTermCount)
	fmt.Printf("   - Long-term memories: %d\n", stats.LongTermCount)
	fmt.Printf("   - Episodic memories: %d\n", stats.EpisodicCount)
	fmt.Printf("   - Semantic memories: %d\n", stats.SemanticCount)
	fmt.Printf("   - Total memory usage: %d units\n", stats.TotalMemoryUsage)
	fmt.Printf("   - Last cleanup: %s\n", stats.LastCleanup.Format("15:04:05"))
	
	fmt.Println()
	
	// Demonstrate memory promotion
	fmt.Println("⬆️  Promoting important memory to long-term...")
	if err := mm.PromoteToLongTerm("user_preference"); err == nil {
		fmt.Println("COMPLETED Successfully promoted user preferences to long-term memory")
	} else {
		fmt.Printf("FAILED Failed to promote memory: %v\n", err)
	}
	
	// Show updated stats
	fmt.Println("📊 Updated Memory Statistics:")
	stats = mm.GetStats()
	fmt.Printf("   - Short-term memories: %d\n", stats.ShortTermCount)
	fmt.Printf("   - Long-term memories: %d\n", stats.LongTermCount)
	
	fmt.Println()
}
package main

import (
	"fmt"
	"log"
	"time"

	"fr0g-ai-master-control/internal/mastercontrol"
	"fr0g-ai-master-control/internal/mastercontrol/memory"
)

func main() {
	fmt.Println("🧠 Memory Manager Demo")
	fmt.Println("=====================")
	
	// Create MCP config
	config := mastercontrol.DefaultMCPConfig()
	
	// Create memory manager
	memManager := mastercontrol.NewMemoryManager(config)
	
	// Start memory manager
	if err := memManager.Start(); err != nil {
		log.Fatalf("Failed to start memory manager: %v", err)
	}
	
	fmt.Println("COMPLETED Memory Manager started successfully")
	fmt.Println()
	
	// Demonstrate memory operations
	demonstrateMemoryOperations(memManager)
	
	// Stop memory manager
	if err := memManager.Stop(); err != nil {
		log.Printf("Error stopping memory manager: %v", err)
	}
	
	fmt.Println("👋 Memory Manager demo complete")
}

func demonstrateMemoryOperations(mm *memory.MemoryManager) {
	fmt.Println("CHECKING Demonstrating Memory Operations:")
	fmt.Println("-----------------------------------")
	
	// Store some short-term memories
	fmt.Println("📝 Storing short-term memories...")
	mm.Store("user_preference", map[string]interface{}{
		"theme": "dark",
		"language": "en",
		"notifications": true,
	})
	
	mm.Store("session_data", map[string]interface{}{
		"user_id": "user123",
		"login_time": time.Now(),
		"ip_address": "192.168.1.100",
	})
	
	mm.Store("temp_calculation", 42.5)
	
	// Store some long-term memories
	fmt.Println("📚 Storing long-term memories...")
	mm.StoreWithType("system_config", map[string]interface{}{
		"version": "1.0.0",
		"deployment": "production",
		"features": []string{"ai", "learning", "consciousness"},
	}, "long_term")
	
	// Store episodic memory
	fmt.Println("📖 Storing episodic memory...")
	episode := &memory.EpisodicMemory{
		ID: "episode_001",
		Event: "User asked about weather",
		Context: map[string]interface{}{
			"location": "San Francisco",
			"time": time.Now(),
			"user_mood": "curious",
		},
		Participants: []string{"user123", "ai_assistant"},
		Outcome: "Provided weather information successfully",
		Importance: 0.7,
		Timestamp: time.Now(),
		Tags: []string{"weather", "query", "successful"},
	}
	mm.StoreEpisode(episode)
	
	// Store semantic memory
	fmt.Println("🧩 Storing semantic memory...")
	concept := &memory.SemanticMemory{
		Concept: "weather",
		Definition: "Atmospheric conditions including temperature, humidity, precipitation, and wind",
		Relations: map[string]float64{
			"temperature": 0.9,
			"humidity": 0.8,
			"precipitation": 0.85,
			"forecast": 0.7,
		},
		Examples: []string{"sunny", "rainy", "cloudy", "stormy"},
		Confidence: 0.95,
	}
	mm.StoreConcept("weather", concept)
	
	fmt.Println()
	
	// Retrieve memories
	fmt.Println("CHECKING Retrieving memories...")
	
	if userPref, err := mm.Retrieve("user_preference"); err == nil {
		fmt.Printf("COMPLETED Retrieved user preferences: %+v\n", userPref)
	}
	
	if sessionData, err := mm.Retrieve("session_data"); err == nil {
		fmt.Printf("COMPLETED Retrieved session data: %+v\n", sessionData)
	}
	
	if tempCalc, err := mm.Retrieve("temp_calculation"); err == nil {
		fmt.Printf("COMPLETED Retrieved calculation: %v\n", tempCalc)
	}
	
	if sysConfig, err := mm.Retrieve("system_config"); err == nil {
		fmt.Printf("COMPLETED Retrieved system config: %+v\n", sysConfig)
	}
	
	fmt.Println()
	
	// Get episodic memories
	fmt.Println("📖 Retrieving episodic memories...")
	episodes := mm.GetEpisodes([]string{"weather"}, 5)
	for _, ep := range episodes {
		fmt.Printf("COMPLETED Episode: %s - %s (Importance: %.2f)\n", ep.ID, ep.Event, ep.Importance)
	}
	
	// Get semantic memory
	fmt.Println("🧩 Retrieving semantic memory...")
	if weatherConcept, err := mm.GetConcept("weather"); err == nil {
		fmt.Printf("COMPLETED Concept 'weather': %s (Confidence: %.2f)\n", weatherConcept.Definition, weatherConcept.Confidence)
		fmt.Printf("   Relations: %+v\n", weatherConcept.Relations)
	}
	
	fmt.Println()
	
	// Show memory statistics
	fmt.Println("📊 Memory Statistics:")
	stats := mm.GetStats()
	fmt.Printf("   - Short-term memories: %d\n", stats.ShortTermCount)
	fmt.Printf("   - Long-term memories: %d\n", stats.LongTermCount)
	fmt.Printf("   - Episodic memories: %d\n", stats.EpisodicCount)
	fmt.Printf("   - Semantic memories: %d\n", stats.SemanticCount)
	fmt.Printf("   - Total memory usage: %d units\n", stats.TotalMemoryUsage)
	fmt.Printf("   - Last cleanup: %s\n", stats.LastCleanup.Format("15:04:05"))
	
	fmt.Println()
	
	// Demonstrate memory promotion
	fmt.Println("⬆️  Promoting important memory to long-term...")
	if err := mm.PromoteToLongTerm("user_preference"); err == nil {
		fmt.Println("COMPLETED Successfully promoted user preferences to long-term memory")
	} else {
		fmt.Printf("FAILED Failed to promote memory: %v\n", err)
	}
	
	// Show updated stats
	fmt.Println("📊 Updated Memory Statistics:")
	stats = mm.GetStats()
	fmt.Printf("   - Short-term memories: %d\n", stats.ShortTermCount)
	fmt.Printf("   - Long-term memories: %d\n", stats.LongTermCount)
	
	fmt.Println()
}
