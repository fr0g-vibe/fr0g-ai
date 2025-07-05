package mastercontrol

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// MasterControlProgram represents the main MCP instance
type MasterControlProgram struct {
	config *MCPConfig
	server *http.Server
}

// NewMasterControlProgram creates a new MCP instance
func NewMasterControlProgram(config *MCPConfig) *MasterControlProgram {
	return &MasterControlProgram{
		config: config,
	}
}

// Start starts the MCP with webhook input system
func (mcp *MasterControlProgram) Start() error {
	mux := http.NewServeMux()
	
	// Health endpoint
	mux.HandleFunc("/health", mcp.healthHandler)
	
	// Status endpoint
	mux.HandleFunc("/status", mcp.statusHandler)
	
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
	
	go func() {
		if err := mcp.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Webhook server error: %v", err)
		}
	}()
	
	return nil
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
	status := map[string]interface{}{
		"status": "running",
		"processors": map[string]string{
			"discord": "Discord message processor - analyzes messages for AI community review",
		},
		"uptime":    time.Now().Format(time.RFC3339),
		"endpoints": []string{"/webhook/discord", "/health", "/status"},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
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
