package mastercontrol

import (
	"fr0g-ai-bridge/internal/mastercontrol/cognitive"
	"fr0g-ai-bridge/internal/mastercontrol/learning"
	"fr0g-ai-bridge/internal/mastercontrol/memory"
	"fr0g-ai-bridge/internal/mastercontrol/monitor"
	"fr0g-ai-bridge/internal/mastercontrol/orchestrator"
	"fr0g-ai-bridge/internal/mastercontrol/workflow"
	"time"
)

// NewMemoryManager creates a new memory manager with default config
func NewMemoryManager(config *MCPConfig) *memory.MemoryManager {
	memConfig := &memory.MemoryConfig{
		ShortTermTTL:        time.Hour,
		LongTermTTL:         time.Hour * 24 * 7, // 1 week
		MaxShortTermEntries: 1000,
		MaxLongTermEntries:  10000,
		MaxEpisodicMemories: 1000,
		MaxSemanticMemories: 5000,
		CleanupInterval:     time.Minute * 15,
		ImportanceThreshold: 0.5,
		CompressionEnabled:  true,
	}
	
	return memory.NewMemoryManager(memConfig)
}

// NewLearningEngine creates a new learning engine with default config
func NewLearningEngine(config *MCPConfig, mem *memory.MemoryManager) *learning.LearningEngine {
	learningConfig := &learning.LearningConfig{
		LearningRate:    0.01,
		AdaptationSpeed: 0.1,
		UpdateInterval:  time.Minute * 5,
	}
	
	return learning.NewLearningEngine(learningConfig, mem)
}

// NewCognitiveEngine creates a new cognitive engine with default config
func NewCognitiveEngine(config *MCPConfig, mem *memory.MemoryManager, learn *learning.LearningEngine) *cognitive.CognitiveEngine {
	cogConfig := &cognitive.CognitiveConfig{
		PatternRecognitionEnabled:  true,
		InsightGenerationEnabled:   true,
		ReflectionEnabled:          config.SystemConsciousness,
		AwarenessUpdateInterval:    time.Second * 30,
		PatternConfidenceThreshold: 0.7,
		MaxPatterns:               1000,
		MaxInsights:               500,
		MaxReflections:            100,
	}
	
	return cognitive.NewCognitiveEngine(cogConfig, mem, learn)
}

// NewSystemMonitor creates a new system monitor with default config
func NewSystemMonitor(config *MCPConfig) *monitor.SystemMonitor {
	monConfig := &monitor.MonitorConfig{
		HealthCheckInterval: config.HealthCheckInterval,
		MetricsInterval:     config.MetricsInterval,
	}
	
	return monitor.NewSystemMonitor(monConfig)
}

// NewWorkflowEngine creates a new workflow engine with default config
func NewWorkflowEngine(config *MCPConfig) *workflow.WorkflowEngine {
	workflowConfig := &workflow.WorkflowConfig{
		MaxConcurrentWorkflows: config.MaxConcurrentWorkflows,
		WorkflowTimeout:        time.Minute * 30,
	}
	
	return workflow.NewWorkflowEngine(workflowConfig)
}

// NewStrategyOrchestrator creates a new strategy orchestrator with default config
func NewStrategyOrchestrator(config *MCPConfig, cog *cognitive.CognitiveEngine, wf *workflow.WorkflowEngine) *orchestrator.StrategyOrchestrator {
	orchConfig := &orchestrator.OrchestratorConfig{
		ResourceOptimization: config.ResourceOptimization,
		PredictiveManagement: config.PredictiveManagement,
	}
	
	return orchestrator.NewStrategyOrchestrator(orchConfig, cog, wf)
}
