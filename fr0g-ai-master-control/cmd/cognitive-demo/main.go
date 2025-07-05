package main

import (
	"fmt"
	"log"
	"time"

	"fr0g-ai-master-control/internal/mastercontrol"
	"fr0g-ai-master-control/internal/mastercontrol/cognitive"
)

func main() {
	fmt.Println("🧠 Cognitive Engine Demo")
	fmt.Println("========================")

	// Create MCP config
	config := mastercontrol.DefaultMCPConfig()

	// Create dependencies
	memManager := mastercontrol.NewMemoryManager(config)
	learningEngine := mastercontrol.NewLearningEngine(config, memManager)
	cognitiveEngine := mastercontrol.NewCognitiveEngine(config, memManager, learningEngine)

	// Start components
	if err := memManager.Start(); err != nil {
		log.Fatalf("Failed to start memory manager: %v", err)
	}

	if err := learningEngine.Start(); err != nil {
		log.Fatalf("Failed to start learning engine: %v", err)
	}

	if err := cognitiveEngine.Start(); err != nil {
		log.Fatalf("Failed to start cognitive engine: %v", err)
	}

	fmt.Println("✅ Cognitive Engine started successfully")
	fmt.Println()

	// Demonstrate cognitive operations
	demonstrateCognitiveOperations(cognitiveEngine)

	// Stop components
	cognitiveEngine.Stop()
	learningEngine.Stop()
	memManager.Stop()

	fmt.Println("👋 Cognitive Engine demo complete")
}

func demonstrateCognitiveOperations(ce *cognitive.CognitiveEngine) {
	fmt.Println("🔍 Demonstrating Cognitive Operations:")
	fmt.Println("-------------------------------------")

	// Wait for cognitive processes to initialize
	time.Sleep(time.Second * 2)

	// Get system awareness
	fmt.Println("🌐 System Awareness:")
	awareness := ce.GetSystemAwareness()
	fmt.Printf("   - Awareness Level: %.3f\n", awareness.AwarenessLevel)
	fmt.Printf("   - Last Update: %s\n", awareness.LastUpdate.Format("15:04:05"))
	fmt.Printf("   - State History Length: %d\n", len(awareness.StateHistory))
	fmt.Printf("   - Component Map Size: %d\n", len(awareness.ComponentMap))
	fmt.Println()

	// Simulate system state for reflection
	fmt.Println("🤔 Triggering Self-Reflection...")
	systemState := map[string]interface{}{
		"status":        "running",
		"load":          0.65,
		"active_users":  42,
		"response_time": "150ms",
		"memory_usage":  "78%",
	}

	// Trigger reflection
	ce.Reflect(systemState)
	ce.Reflect(map[string]interface{}{
		"event":     "high_load_detected",
		"threshold": 0.8,
		"current":   0.85,
	})
	ce.Reflect(map[string]interface{}{
		"optimization": "cache_hit_rate_improved",
		"before":       0.65,
		"after":        0.82,
	})

	// Wait for reflections to be processed
	time.Sleep(time.Second * 1)

	// Get reflections
	fmt.Println("💭 Generated Reflections:")
	reflections := ce.GetReflections()
	for i, reflection := range reflections {
		fmt.Printf("   %d. [%s] %s\n", i+1, reflection.Type, reflection.Content)
		fmt.Printf("      Created: %s, Depth: %d\n", reflection.CreatedAt.Format("15:04:05"), reflection.Depth)
	}
	fmt.Println()

	// Get patterns
	fmt.Println("🔍 Recognized Patterns:")
	patterns := ce.GetPatternsMap()
	if len(patterns) == 0 {
		fmt.Println("   No patterns recognized yet (pattern recognition is ongoing)")
	} else {
		for id, pattern := range patterns {
			fmt.Printf("   - %s: %s (Confidence: %.2f)\n", id, pattern.Description, pattern.Confidence)
		}
	}
	fmt.Println()

	// Get insights
	fmt.Println("💡 Generated Insights:")
	insights := ce.GetInsights()
	if len(insights) == 0 {
		fmt.Println("   No insights generated yet (insight generation is ongoing)")
	} else {
		for i, insight := range insights {
			fmt.Printf("   %d. [%s] %s\n", i+1, insight.Type, insight.Content)
			fmt.Printf("      Confidence: %.2f, Impact: %s, Actionable: %v\n",
				insight.Confidence, insight.Impact, insight.Actionable)
		}
	}
	fmt.Println()

	fmt.Println("🧠 Cognitive processes are running in the background...")
	fmt.Println("   - Pattern recognition every 30 seconds")
	fmt.Println("   - Insight generation every 5 minutes")
	fmt.Println("   - Awareness updates every 30 seconds")
	fmt.Println()

	// Wait a bit to show ongoing processes
	fmt.Println("⏳ Waiting 5 seconds to observe cognitive processes...")
	time.Sleep(time.Second * 5)

	// Check updated awareness
	fmt.Println("🔄 Updated System Awareness:")
	awareness = ce.GetSystemAwareness()
	fmt.Printf("   - Awareness Level: %.3f\n", awareness.AwarenessLevel)
	fmt.Printf("   - Last Update: %s\n", awareness.LastUpdate.Format("15:04:05"))
	fmt.Printf("   - State History Length: %d\n", len(awareness.StateHistory))

	fmt.Println()
}
