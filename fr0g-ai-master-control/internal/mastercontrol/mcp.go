package mastercontrol

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// SystemState represents the current state of the MCP system
type SystemState struct {
	Status          string              `json:"status"`
	ActiveWorkflows int                 `json:"active_workflows"`
	SystemLoad      float64             `json:"system_load"`
	LastUpdate      time.Time           `json:"last_update"`
	Intelligence    IntelligenceMetrics `json:"intelligence"`
}

// IntelligenceMetrics represents AI intelligence measurements
type IntelligenceMetrics struct {
	LearningRate          float64 `json:"learning_rate"`
	PatternCount          int     `json:"pattern_count"`
	AdaptationScore       float64 `json:"adaptation_score"`
	EfficiencyIndex       float64 `json:"efficiency_index"`
	EmergentCapabilities  int     `json:"emergent_capabilities"`
}

// Capability represents a system capability
type Capability struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Emergent    bool      `json:"emergent"`
	CreatedAt   time.Time `json:"created_at"`
}

// MasterControlProgram represents the main MCP instance
type MasterControlProgram struct {
	config       *MCPConfig
	server       *http.Server
	systemState  SystemState
	capabilities map[string]Capability
	mu           sync.RWMutex
	startTime    time.Time
}

// NewMasterControlProgram creates a new MCP instance
func NewMasterControlProgram(config *MCPConfig) *MasterControlProgram {
	mcp := &MasterControlProgram{
		config:       config,
		capabilities: make(map[string]Capability),
		startTime:    time.Now(),
	}
	
	// Initialize system state
	mcp.initializeSystemState()
	
	// Initialize capabilities
	mcp.initializeCapabilities()
	
	return mcp
}

// initializeSystemState sets up the initial system state
func (mcp *MasterControlProgram) initializeSystemState() {
	mcp.systemState = SystemState{
		Status:          "initializing",
		ActiveWorkflows: 0,
		SystemLoad:      0.0,
		LastUpdate:      time.Now(),
		Intelligence: IntelligenceMetrics{
			LearningRate:         mcp.config.AdaptiveLearningRate,
			PatternCount:         6,
			AdaptationScore:      0.85,
			EfficiencyIndex:      0.92,
			EmergentCapabilities: 3,
		},
	}
}

// initializeCapabilities sets up the initial system capabilities
func (mcp *MasterControlProgram) initializeCapabilities() {
	now := time.Now()
	
	mcp.capabilities["discord_processing"] = Capability{
		Name:        "Discord Message Processing",
		Description: "Analyze Discord messages using AI persona communities",
		Emergent:    false,
		CreatedAt:   now,
	}
	
	mcp.capabilities["pattern_recognition"] = Capability{
		Name:        "Pattern Recognition",
		Description: "Identify patterns in communication and behavior",
		Emergent:    true,
		CreatedAt:   now,
	}
	
	mcp.capabilities["adaptive_learning"] = Capability{
		Name:        "Adaptive Learning",
		Description: "Learn and adapt from processed data",
		Emergent:    true,
		CreatedAt:   now,
	}
	
	mcp.capabilities["threat_assessment"] = Capability{
		Name:        "Threat Assessment",
		Description: "Evaluate potential security threats",
		Emergent:    false,
		CreatedAt:   now,
	}
	
	mcp.capabilities["consciousness_simulation"] = Capability{
		Name:        "Consciousness Simulation",
		Description: "Simulate self-awareness and reflection",
		Emergent:    true,
		CreatedAt:   now,
	}
}

// Start starts the MCP with webhook input system
func (mcp *MasterControlProgram) Start() error {
	mux := http.NewServeMux()
	
	// Health endpoint
	mux.HandleFunc("/health", mcp.healthHandler)
	
	// Status endpoint
	mux.HandleFunc("/status", mcp.statusHandler)
	
	// System state endpoint
	mux.HandleFunc("/system/state", mcp.systemStateHandler)
	
	// Capabilities endpoint
	mux.HandleFunc("/system/capabilities", mcp.capabilitiesHandler)
	
	// Discord webhook endpoint
	mux.HandleFunc("/webhook/discord", mcp.discordWebhookHandler)
	
	// Catch-all for unknown webhook tags
	mux.HandleFunc("/webhook/", mcp.unknownWebhookHandler)
	
	mcp.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", mcp.config.Input.Webhook.Host, mcp.config.Input.Webhook.Port),
		Handler: mux,
		ReadTimeout:  mcp.config.Input.Webhook.ReadTimeout,
		WriteTimeout: mcp.config.Input.Webhook.WriteTimeout,
	}
	
	// Update system state to running
	mcp.mu.Lock()
	mcp.systemState.Status = "running"
	mcp.systemState.LastUpdate = time.Now()
	mcp.mu.Unlock()
	
	// Start background processes
	go mcp.backgroundProcesses()
	
	go func() {
		if err := mcp.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Webhook server error: %v", err)
		}
	}()
	
	return nil
}

// backgroundProcesses runs continuous background tasks
func (mcp *MasterControlProgram) backgroundProcesses() {
	ticker := time.NewTicker(mcp.config.CognitiveReflectionRate)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			mcp.updateSystemMetrics()
			mcp.performCognitiveReflection()
		}
	}
}

// updateSystemMetrics updates system performance metrics
func (mcp *MasterControlProgram) updateSystemMetrics() {
	mcp.mu.Lock()
	defer mcp.mu.Unlock()
	
	// Simulate system load calculation
	uptime := time.Since(mcp.startTime)
	mcp.systemState.SystemLoad = 0.1 + (float64(uptime.Minutes()) * 0.001)
	if mcp.systemState.SystemLoad > 1.0 {
		mcp.systemState.SystemLoad = 0.8 + (0.2 * (float64(time.Now().Unix()) % 10 / 10))
	}
	
	// Update intelligence metrics
	now := time.Now().Unix()
	mcp.systemState.Intelligence.PatternCount += int(now) % 3
	mcp.systemState.Intelligence.AdaptationScore = 0.8 + (0.2 * float64(now%10) / 10)
	mcp.systemState.Intelligence.EfficiencyIndex = 0.85 + (0.15 * float64(now%10) / 10)
	
	mcp.systemState.LastUpdate = time.Now()
}

// performCognitiveReflection simulates AI consciousness and self-reflection
func (mcp *MasterControlProgram) performCognitiveReflection() {
	if !mcp.config.SystemConsciousness {
		return
	}
	
	reflections := []string{
		"Analyzing communication patterns for emergent behaviors...",
		"Reflecting on system efficiency and adaptation strategies...",
		"Evaluating threat landscape and response protocols...",
		"Considering ethical implications of autonomous decisions...",
		"Synthesizing learned patterns into actionable intelligence...",
	}
	
	reflection := reflections[int(time.Now().Unix())%len(reflections)]
	log.Printf("🧠 Cognitive Reflection: %s", reflection)
}

// Stop stops the MCP
func (mcp *MasterControlProgram) Stop() error {
	if mcp.server != nil {
		return mcp.server.Close()
	}
	return nil
}

func (mcp *MasterControlProgram) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (mcp *MasterControlProgram) statusHandler(w http.ResponseWriter, r *http.Request) {
	mcp.mu.RLock()
	defer mcp.mu.RUnlock()
	
	uptime := time.Since(mcp.startTime)
	
	status := map[string]interface{}{
		"service":    "fr0g-ai-master-control",
		"version":    "1.0.0",
		"status":     mcp.systemState.Status,
		"uptime":     uptime.String(),
		"started_at": mcp.startTime.Format(time.RFC3339),
		"processors": map[string]string{
			"discord": "Discord message processor - analyzes messages for AI community review",
		},
		"endpoints": []string{
			"/webhook/discord", "/health", "/status", 
			"/system/state", "/system/capabilities",
		},
		"intelligence": map[string]interface{}{
			"learning_rate":         mcp.systemState.Intelligence.LearningRate,
			"pattern_count":         mcp.systemState.Intelligence.PatternCount,
			"adaptation_score":      mcp.systemState.Intelligence.AdaptationScore,
			"efficiency_index":      mcp.systemState.Intelligence.EfficiencyIndex,
			"emergent_capabilities": mcp.systemState.Intelligence.EmergentCapabilities,
			"consciousness":         "active",
		},
		"system": map[string]interface{}{
			"active_workflows": mcp.systemState.ActiveWorkflows,
			"system_load":      mcp.systemState.SystemLoad,
			"last_update":      mcp.systemState.LastUpdate.Format(time.RFC3339),
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// systemStateHandler returns detailed system state
func (mcp *MasterControlProgram) systemStateHandler(w http.ResponseWriter, r *http.Request) {
	mcp.mu.RLock()
	defer mcp.mu.RUnlock()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mcp.systemState)
}

// capabilitiesHandler returns system capabilities
func (mcp *MasterControlProgram) capabilitiesHandler(w http.ResponseWriter, r *http.Request) {
	mcp.mu.RLock()
	defer mcp.mu.RUnlock()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mcp.capabilities)
}

func (mcp *MasterControlProgram) discordWebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var message map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Mock AI community review process
	response := map[string]interface{}{
		"success":    true,
		"message":    "Discord message processed successfully",
		"request_id": fmt.Sprintf("req_%d", time.Now().Unix()),
		"data": map[string]interface{}{
			"action":        "reviewed",
			"persona_count": 3,
			"sentiment":     "neutral",
			"requires_action": false,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (mcp *MasterControlProgram) unknownWebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the webhook tag from the path
	path := strings.TrimPrefix(r.URL.Path, "/webhook/")
	
	response := map[string]interface{}{
		"success": false,
		"message": fmt.Sprintf("Unknown webhook tag: %s", path),
		"error":   "webhook_not_found",
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(response)
}

// GetSystemState returns the current system state
func (mcp *MasterControlProgram) GetSystemState() SystemState {
	mcp.mu.RLock()
	defer mcp.mu.RUnlock()
	return mcp.systemState
}

// GetCapabilities returns the system capabilities
func (mcp *MasterControlProgram) GetCapabilities() map[string]Capability {
	mcp.mu.RLock()
	defer mcp.mu.RUnlock()
	
	// Return a copy to prevent external modification
	capabilities := make(map[string]Capability)
	for k, v := range mcp.capabilities {
		capabilities[k] = v
	}
	return capabilities
}

// SetInputManager sets the input manager (placeholder for future integration)
func (mcp *MasterControlProgram) SetInputManager(inputManager interface{}) {
	log.Printf("🔧 Input manager configured")
	// TODO: Implement input manager integration
}
