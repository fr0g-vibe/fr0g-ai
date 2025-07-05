package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// HTTPServer represents the HTTP server for the MCP
type HTTPServer struct {
	port   int
	router *mux.Router
	server *http.Server
}

// NewHTTPServer creates a new HTTP server instance
func NewHTTPServer(port int) *HTTPServer {
	router := mux.NewRouter()
	
	server := &HTTPServer{
		port:   port,
		router: router,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      router,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	}
	
	server.setupRoutes()
	return server
}

// setupRoutes configures the HTTP routes
func (s *HTTPServer) setupRoutes() {
	// Health check endpoint
	s.router.HandleFunc("/health", s.healthHandler).Methods("GET")
	
	// Status endpoint
	s.router.HandleFunc("/status", s.statusHandler).Methods("GET")
	
	// API routes
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/mcp/status", s.mcpStatusHandler).Methods("GET")
}

// Start starts the HTTP server
func (s *HTTPServer) Start() error {
	fmt.Printf("🌐 Starting HTTP server on port %d\n", s.port)
	return s.server.ListenAndServe()
}

// Stop stops the HTTP server
func (s *HTTPServer) Stop() error {
	fmt.Println("🛑 Stopping HTTP server...")
	return s.server.Close()
}

// healthHandler handles health check requests
func (s *HTTPServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"service":   "fr0g-ai-mcp",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// statusHandler handles status requests
func (s *HTTPServer) statusHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service":     "fr0g-ai-mcp",
		"status":      "operational",
		"uptime":      time.Since(time.Now()).String(), // TODO: Track actual uptime
		"components": map[string]string{
			"http_server":      "running",
			"input_processor":  "ready",
			"ai_community":     "ready",
			"cognitive_engine": "ready",
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// mcpStatusHandler handles MCP-specific status requests
func (s *HTTPServer) mcpStatusHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"mcp_version":        "1.0.0",
		"active_workflows":   0,
		"processed_messages": 0,
		"ai_personas":        0,
		"system_load":        0.0,
		"learning_enabled":   true,
		"consciousness":      true,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
