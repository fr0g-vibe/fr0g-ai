package workflow

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// WorkflowEngine handles dynamic workflow generation and execution
type WorkflowEngine struct {
	config             *WorkflowConfig
	activeWorkflows    map[string]*SampleWorkflow
	completedWorkflows []string
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
func NewWorkflowEngine(config *WorkflowConfig) *WorkflowEngine {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkflowEngine{
		config:             config,
		activeWorkflows:    make(map[string]*SampleWorkflow),
		completedWorkflows: make([]string, 0),
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
