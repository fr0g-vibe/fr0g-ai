package workflow

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-master-control/internal/mastercontrol/input"
)

// WorkflowEngine handles dynamic workflow generation and execution
type WorkflowEngine struct {
	config             *WorkflowConfig
	activeWorkflows    map[string]*SampleWorkflow
	completedWorkflows []string
	fr0gIOClient       input.Fr0gIOClient
	inputHandler       input.Fr0gIOInputHandler
	mu                 sync.RWMutex
	ctx                context.Context
	cancel             context.CancelFunc
}

// WorkflowConfig holds workflow engine configuration
type WorkflowConfig struct {
	MaxConcurrentWorkflows int           `yaml:"max_concurrent_workflows"`
	WorkflowTimeout        time.Duration `yaml:"workflow_timeout"`
	AutoStartWorkflows     bool          `yaml:"auto_start_workflows"`
	WorkflowInterval       time.Duration `yaml:"workflow_interval"`
}

// SampleWorkflow represents a demonstration workflow
type SampleWorkflow struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Steps       []WorkflowStep         `json:"steps"`
	Status      string                 `json:"status"`
	Progress    float64                `json:"progress"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
}

// WorkflowStep represents a single step in a workflow
type WorkflowStep struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	Duration    time.Duration          `json:"duration"`
	Output      map[string]interface{} `json:"output"`
	Error       string                 `json:"error,omitempty"`
}

// NewWorkflowEngine creates a new workflow engine
func NewWorkflowEngine(config *WorkflowConfig, fr0gIOClient input.Fr0gIOClient, inputHandler input.Fr0gIOInputHandler) *WorkflowEngine {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkflowEngine{
		config:             config,
		activeWorkflows:    make(map[string]*SampleWorkflow),
		completedWorkflows: make([]string, 0),
		fr0gIOClient:       fr0gIOClient,
		inputHandler:       inputHandler,
		ctx:                ctx,
		cancel:             cancel,
	}
}

// Start begins workflow engine operation
func (we *WorkflowEngine) Start() error {
	log.Println("Workflow Engine: Starting workflow processes...")

	// Set default values if not configured
	if we.config.WorkflowInterval == 0 {
		we.config.WorkflowInterval = 2 * time.Minute
	}
	if we.config.MaxConcurrentWorkflows == 0 {
		we.config.MaxConcurrentWorkflows = 3
	}

	// Always start workflows for demo
	go we.workflowManagementLoop()

	return nil
}

// Stop gracefully stops the workflow engine
func (we *WorkflowEngine) Stop() error {
	log.Println("Workflow Engine: Stopping workflow processes...")
	we.cancel()
	return nil
}

// GetActiveWorkflowCount returns the number of active workflows
func (we *WorkflowEngine) GetActiveWorkflowCount() int {
	we.mu.RLock()
	defer we.mu.RUnlock()
	return len(we.activeWorkflows)
}

// GetCompletedWorkflowCount returns the number of completed workflows
func (we *WorkflowEngine) GetCompletedWorkflowCount() int {
	we.mu.RLock()
	defer we.mu.RUnlock()
	return len(we.completedWorkflows)
}

// CreateSampleWorkflows creates demonstration workflows to show system capabilities
func (we *WorkflowEngine) CreateSampleWorkflows() []interface{} {
	workflows := []*SampleWorkflow{
		{
			ID:          fmt.Sprintf("workflow_cognitive_analysis_%d", time.Now().UnixNano()),
			Name:        "Cognitive System Analysis",
			Description: "Analyze system cognitive capabilities and generate insights",
			Status:      "ready",
			Progress:    0.0,
			Metadata: map[string]interface{}{
				"type":      "cognitive",
				"priority":  "high",
				"automated": true,
			},
			Steps: []WorkflowStep{
				{
					ID:          "step_1",
					Name:        "Pattern Recognition",
					Description: "Identify behavioral patterns in system operation",
					Status:      "pending",
				},
				{
					ID:          "step_2",
					Name:        "Insight Generation",
					Description: "Generate actionable insights from recognized patterns",
					Status:      "pending",
				},
				{
					ID:          "step_3",
					Name:        "Self-Reflection",
					Description: "Perform meta-cognitive analysis of system state",
					Status:      "pending",
				},
				{
					ID:          "step_4",
					Name:        "Knowledge Integration",
					Description: "Integrate new knowledge into long-term memory",
					Status:      "pending",
				},
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          fmt.Sprintf("workflow_system_optimization_%d", time.Now().UnixNano()),
			Name:        "System Optimization",
			Description: "Optimize system performance based on learned patterns",
			Status:      "ready",
			Progress:    0.0,
			Metadata: map[string]interface{}{
				"type":      "optimization",
				"priority":  "medium",
				"automated": true,
			},
			Steps: []WorkflowStep{
				{
					ID:          "step_1",
					Name:        "Performance Analysis",
					Description: "Analyze current system performance metrics",
					Status:      "pending",
				},
				{
					ID:          "step_2",
					Name:        "Bottleneck Identification",
					Description: "Identify performance bottlenecks and inefficiencies",
					Status:      "pending",
				},
				{
					ID:          "step_3",
					Name:        "Optimization Strategy",
					Description: "Develop optimization strategy based on analysis",
					Status:      "pending",
				},
			},
			CreatedAt: time.Now(),
		},
	}

	// Convert to []interface{} for interface compatibility
	result := make([]interface{}, len(workflows))
	for i, workflow := range workflows {
		result[i] = workflow
	}

	return result
}

// ExecuteSampleWorkflow executes a sample workflow to demonstrate capabilities
func (we *WorkflowEngine) ExecuteSampleWorkflow(ctx context.Context, workflow *SampleWorkflow) error {
	log.Printf("Workflow Engine: Starting workflow '%s'", workflow.Name)

	workflow.Status = "running"
	startTime := time.Now()
	workflow.StartedAt = &startTime

	totalSteps := len(workflow.Steps)

	for i, step := range workflow.Steps {
		select {
		case <-ctx.Done():
			workflow.Status = "cancelled"
			return ctx.Err()
		default:
			// Execute step
			log.Printf("Workflow Engine: Executing step '%s'", step.Name)

			step.Status = "running"
			stepStart := time.Now()

			// Simulate step execution
			time.Sleep(time.Millisecond * 500) // Simulate work

			step.Duration = time.Since(stepStart)
			step.Status = "completed"
			step.Output = map[string]interface{}{
				"execution_time": step.Duration.String(),
				"success":        true,
				"step_id":        step.ID,
			}

			workflow.Steps[i] = step
			workflow.Progress = float64(i+1) / float64(totalSteps)

			log.Printf("Workflow Engine: Step '%s' completed in %v", step.Name, step.Duration)
		}
	}

	workflow.Status = "completed"
	completedTime := time.Now()
	workflow.CompletedAt = &completedTime
	workflow.Progress = 1.0

	log.Printf("Workflow Engine: Workflow '%s' completed successfully in %v",
		workflow.Name, completedTime.Sub(*workflow.StartedAt))

	return nil
}

// workflowManagementLoop manages automatic workflow execution
func (we *WorkflowEngine) workflowManagementLoop() {
	ticker := time.NewTicker(we.config.WorkflowInterval)
	defer ticker.Stop()

	// Start initial workflows
	we.startSampleWorkflows()

	for {
		select {
		case <-we.ctx.Done():
			return
		case <-ticker.C:
			we.manageContinuousWorkflows()
		}
	}
}

// startSampleWorkflows starts the initial set of sample workflows
func (we *WorkflowEngine) startSampleWorkflows() {
	workflows := we.CreateSampleWorkflows()

	for _, workflow := range workflows {
		if len(we.activeWorkflows) < we.config.MaxConcurrentWorkflows {
			if sampleWorkflow, ok := workflow.(*SampleWorkflow); ok {
				we.startWorkflow(sampleWorkflow)
			}
		}
	}
}

// startWorkflow starts a single workflow
func (we *WorkflowEngine) startWorkflow(workflow *SampleWorkflow) {
	we.mu.Lock()
	we.activeWorkflows[workflow.ID] = workflow
	we.mu.Unlock()

	go func() {
		defer func() {
			we.mu.Lock()
			delete(we.activeWorkflows, workflow.ID)
			we.completedWorkflows = append(we.completedWorkflows, workflow.ID)
			we.mu.Unlock()
		}()

		if err := we.ExecuteSampleWorkflow(we.ctx, workflow); err != nil {
			log.Printf("Workflow Engine: Workflow %s failed: %v", workflow.Name, err)
		}
	}()
}

// manageContinuousWorkflows manages ongoing workflow execution
func (we *WorkflowEngine) manageContinuousWorkflows() {
	we.mu.Lock()
	activeCount := len(we.activeWorkflows)
	we.mu.Unlock()

	// Start new workflows if we have capacity
	if activeCount < we.config.MaxConcurrentWorkflows {
		workflows := we.CreateSampleWorkflows()

		for _, workflow := range workflows {
			if len(we.activeWorkflows) < we.config.MaxConcurrentWorkflows {
				if sampleWorkflow, ok := workflow.(*SampleWorkflow); ok {
					we.startWorkflow(sampleWorkflow)
					break // Start one at a time
				}
			}
		}
	}
}

// ProcessInputEvent processes an input event from fr0g-ai-io and triggers appropriate workflows
func (we *WorkflowEngine) ProcessInputEvent(ctx context.Context, event *input.InputEvent) (*input.InputEventResponse, error) {
	log.Printf("Workflow Engine: Processing input event %s of type %s", event.ID, event.Type)

	// Create a specialized workflow for this input event
	workflow := we.createInputEventWorkflow(event)
	
	// Start the workflow
	we.startWorkflow(workflow)

	// Perform immediate threat analysis
	threatResult, err := we.performThreatAnalysis(ctx, event)
	if err != nil {
		log.Printf("Workflow Engine: Threat analysis failed for event %s: %v", event.ID, err)
	}

	// Generate response actions based on analysis
	actions := we.generateResponseActions(event, threatResult)

	// Send threat analysis result back to fr0g-ai-io if available
	if threatResult != nil && we.fr0gIOClient != nil {
		go func() {
			if err := we.fr0gIOClient.SendThreatAnalysisResult(context.Background(), threatResult); err != nil {
				log.Printf("Workflow Engine: Failed to send threat analysis result: %v", err)
			}
		}()
	}

	// Execute output actions through fr0g-ai-io
	for _, action := range actions {
		if we.fr0gIOClient != nil {
			command := &input.OutputCommand{
				ID:       fmt.Sprintf("cmd_%s_%d", event.ID, time.Now().UnixNano()),
				Type:     action.Type,
				Target:   action.Target,
				Content:  action.Content,
				Metadata: action.Metadata,
				Priority: event.Priority,
			}

			go func(cmd *input.OutputCommand) {
				if _, err := we.fr0gIOClient.SendOutputCommand(context.Background(), cmd); err != nil {
					log.Printf("Workflow Engine: Failed to send output command: %v", err)
				}
			}(command)
		}
	}

	return &input.InputEventResponse{
		EventID:     event.ID,
		Processed:   true,
		Actions:     actions,
		Analysis:    threatResult,
		Metadata:    map[string]interface{}{"workflow_id": workflow.ID},
		ProcessedAt: time.Now(),
	}, nil
}

// createInputEventWorkflow creates a specialized workflow for processing input events
func (we *WorkflowEngine) createInputEventWorkflow(event *input.InputEvent) *SampleWorkflow {
	return &SampleWorkflow{
		ID:          fmt.Sprintf("workflow_input_event_%s_%d", event.Type, time.Now().UnixNano()),
		Name:        fmt.Sprintf("Input Event Processing - %s", event.Type),
		Description: fmt.Sprintf("Process %s input event from %s", event.Type, event.Source),
		Status:      "ready",
		Progress:    0.0,
		Metadata: map[string]interface{}{
			"type":         "input_processing",
			"input_type":   event.Type,
			"input_source": event.Source,
			"event_id":     event.ID,
			"priority":     event.Priority,
		},
		Steps: []WorkflowStep{
			{
				ID:          "step_1",
				Name:        "Content Analysis",
				Description: "Analyze input content for patterns and context",
				Status:      "pending",
			},
			{
				ID:          "step_2",
				Name:        "Threat Assessment",
				Description: "Assess potential security threats in the input",
				Status:      "pending",
			},
			{
				ID:          "step_3",
				Name:        "Response Generation",
				Description: "Generate appropriate response actions",
				Status:      "pending",
			},
			{
				ID:          "step_4",
				Name:        "Learning Integration",
				Description: "Integrate insights into system knowledge base",
				Status:      "pending",
			},
		},
		CreatedAt: time.Now(),
	}
}

// performThreatAnalysis performs threat analysis on input events
func (we *WorkflowEngine) performThreatAnalysis(ctx context.Context, event *input.InputEvent) (*input.ThreatAnalysisResult, error) {
	// Simulate threat analysis processing
	threatScore := 0.0
	threatLevel := "low"
	threatTypes := []string{}
	indicators := []input.ThreatIndicator{}

	// Basic content analysis for demonstration
	content := event.Content
	if len(content) > 1000 {
		threatScore += 0.2
		threatTypes = append(threatTypes, "suspicious_length")
	}

	// Check for suspicious patterns (simplified)
	suspiciousPatterns := []string{"hack", "attack", "exploit", "malware", "phishing"}
	for _, pattern := range suspiciousPatterns {
		if contains(content, pattern) {
			threatScore += 0.3
			threatTypes = append(threatTypes, "suspicious_content")
			indicators = append(indicators, input.ThreatIndicator{
				Type:        "keyword",
				Value:       pattern,
				Confidence:  0.8,
				Description: fmt.Sprintf("Suspicious keyword detected: %s", pattern),
			})
		}
	}

	// Determine threat level based on score
	if threatScore >= 0.7 {
		threatLevel = "critical"
	} else if threatScore >= 0.5 {
		threatLevel = "high"
	} else if threatScore >= 0.3 {
		threatLevel = "medium"
	}

	// Generate mitigation recommendations
	mitigation := []string{}
	if threatScore > 0.5 {
		mitigation = append(mitigation, "Monitor source closely")
		mitigation = append(mitigation, "Apply content filtering")
	}
	if threatScore > 0.7 {
		mitigation = append(mitigation, "Block source temporarily")
		mitigation = append(mitigation, "Alert security team")
	}

	return &input.ThreatAnalysisResult{
		EventID:     event.ID,
		ThreatLevel: threatLevel,
		ThreatScore: threatScore,
		ThreatTypes: threatTypes,
		Indicators:  indicators,
		Mitigation:  mitigation,
		Confidence:  0.85,
		Analysis:    fmt.Sprintf("Automated threat analysis completed for %s event", event.Type),
		Metadata: map[string]interface{}{
			"analysis_version": "1.0",
			"processing_time":  time.Now().Format(time.RFC3339),
		},
		AnalyzedAt: time.Now(),
		RecommendedActions: we.generateThreatResponseActions(event, threatLevel, threatScore),
	}, nil
}

// generateResponseActions generates appropriate response actions for input events
func (we *WorkflowEngine) generateResponseActions(event *input.InputEvent, threatResult *input.ThreatAnalysisResult) []input.OutputAction {
	actions := []input.OutputAction{}

	// Generate acknowledgment response
	actions = append(actions, input.OutputAction{
		Type:    event.Type,
		Target:  event.Source,
		Content: "Message received and processed by Master Control Program",
		Metadata: map[string]interface{}{
			"response_type": "acknowledgment",
			"event_id":      event.ID,
		},
	})

	// Add threat-specific actions if needed
	if threatResult != nil && threatResult.ThreatScore > 0.5 {
		actions = append(actions, input.OutputAction{
			Type:    "alert",
			Target:  "security_team",
			Content: fmt.Sprintf("High threat detected in %s from %s: %s", event.Type, event.Source, threatResult.Analysis),
			Metadata: map[string]interface{}{
				"threat_level": threatResult.ThreatLevel,
				"threat_score": threatResult.ThreatScore,
				"event_id":     event.ID,
			},
		})
	}

	return actions
}

// generateThreatResponseActions generates response actions based on threat analysis
func (we *WorkflowEngine) generateThreatResponseActions(event *input.InputEvent, threatLevel string, threatScore float64) []input.OutputAction {
	actions := []input.OutputAction{}

	if threatScore > 0.7 {
		// Critical threat - immediate response
		actions = append(actions, input.OutputAction{
			Type:    "alert",
			Target:  "security_ops",
			Content: fmt.Sprintf("CRITICAL THREAT DETECTED: %s from %s requires immediate attention", event.Type, event.Source),
			Metadata: map[string]interface{}{
				"priority":     "critical",
				"threat_level": threatLevel,
				"event_id":     event.ID,
			},
		})
	} else if threatScore > 0.3 {
		// Medium threat - monitoring response
		actions = append(actions, input.OutputAction{
			Type:    "log",
			Target:  "security_log",
			Content: fmt.Sprintf("Threat detected in %s from %s - monitoring initiated", event.Type, event.Source),
			Metadata: map[string]interface{}{
				"threat_level": threatLevel,
				"event_id":     event.ID,
			},
		})
	}

	return actions
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if toLower(s[i+j]) != toLower(substr[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + ('a' - 'A')
	}
	return c
}
