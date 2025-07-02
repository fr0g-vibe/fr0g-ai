package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/fr0g-vibe/fr0g-ai-bridge/internal/client"
	"github.com/fr0g-vibe/fr0g-ai-bridge/internal/models"
)

// RESTServer handles REST API requests
type RESTServer struct {
	client *client.OpenWebUIClient
	router *mux.Router
}

// NewRESTServer creates a new REST server
func NewRESTServer(openWebUIClient *client.OpenWebUIClient) *RESTServer {
	server := &RESTServer{
		client: openWebUIClient,
		router: mux.NewRouter(),
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures the REST API routes
func (s *RESTServer) setupRoutes() {
	// Health check endpoint
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")

	// Chat completion endpoint
	s.router.HandleFunc("/api/chat/completions", s.handleChatCompletion).Methods("POST")

	// Add CORS middleware
	s.router.Use(s.corsMiddleware)
}

// GetRouter returns the configured router
func (s *RESTServer) GetRouter() *mux.Router {
	return s.router
}

// handleHealth handles health check requests
func (s *RESTServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Check OpenWebUI health
	err := s.client.HealthCheck(ctx)
	
	response := models.HealthResponse{
		Time:    time.Now(),
		Version: "1.0.0",
	}

	if err != nil {
		response.Status = "unhealthy"
		w.WriteHeader(http.StatusServiceUnavailable)
		log.Printf("Health check failed: %v", err)
	} else {
		response.Status = "healthy"
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleChatCompletion handles chat completion requests
func (s *RESTServer) handleChatCompletion(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req models.ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := s.validateChatCompletionRequest(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	// Forward to OpenWebUI
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	resp, err := s.client.ChatCompletion(ctx, &req)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to process chat completion", err)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// validateChatCompletionRequest validates the chat completion request
func (s *RESTServer) validateChatCompletionRequest(req *models.ChatCompletionRequest) error {
	if req.Model == "" {
		return fmt.Errorf("model is required")
	}
	if len(req.Messages) == 0 {
		return fmt.Errorf("messages are required")
	}
	for i, msg := range req.Messages {
		if msg.Role == "" {
			return fmt.Errorf("message %d: role is required", i)
		}
		if msg.Content == "" {
			return fmt.Errorf("message %d: content is required", i)
		}
	}
	return nil
}

// writeError writes an error response
func (s *RESTServer) writeError(w http.ResponseWriter, statusCode int, message string, err error) {
	log.Printf("API Error: %s - %v", message, err)
	
	errorResp := models.ErrorResponse{
		Error:   message,
		Code:    statusCode,
	}
	
	if err != nil {
		errorResp.Message = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResp)
}

// corsMiddleware adds CORS headers
func (s *RESTServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/fr0g-vibe/fr0g-ai-bridge/internal/client"
)

// RESTServer handles HTTP REST API requests
type RESTServer struct {
	openWebUIClient *client.OpenWebUIClient
	router          *mux.Router
}

// NewRESTServer creates a new REST server instance
func NewRESTServer(openWebUIClient *client.OpenWebUIClient) *RESTServer {
	server := &RESTServer{
		openWebUIClient: openWebUIClient,
		router:          mux.NewRouter(),
	}
	
	server.setupRoutes()
	return server
}

// GetRouter returns the configured router
func (s *RESTServer) GetRouter() *mux.Router {
	return s.router
}

// setupRoutes configures all REST API routes
func (s *RESTServer) setupRoutes() {
	// Health check endpoint
	s.router.HandleFunc("/health", s.healthHandler).Methods("GET")
	
	// AI chat endpoint - bridge to external AI systems
	s.router.HandleFunc("/api/v1/chat", s.chatHandler).Methods("POST")
	
	// AI models endpoint - list available models
	s.router.HandleFunc("/api/v1/models", s.modelsHandler).Methods("GET")
}

// ChatRequest represents an incoming chat request
type ChatRequest struct {
	Message string `json:"message"`
	Model   string `json:"model,omitempty"`
}

// ChatResponse represents a chat response
type ChatResponse struct {
	Response  string    `json:"response"`
	Model     string    `json:"model"`
	Timestamp time.Time `json:"timestamp"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// healthHandler handles health check requests
func (s *RESTServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "fr0g-ai-bridge",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// chatHandler handles chat requests and bridges to external AI systems
func (s *RESTServer) chatHandler(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, http.StatusBadRequest, "Invalid request format", err.Error())
		return
	}

	if req.Message == "" {
		s.sendError(w, http.StatusBadRequest, "Message is required", "")
		return
	}

	// Bridge request to OpenWebUI or other external AI system
	response, err := s.openWebUIClient.SendMessage(req.Message, req.Model)
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Failed to process request", err.Error())
		return
	}

	chatResp := ChatResponse{
		Response:  response,
		Model:     req.Model,
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chatResp)
}

// modelsHandler handles requests for available AI models
func (s *RESTServer) modelsHandler(w http.ResponseWriter, r *http.Request) {
	models, err := s.openWebUIClient.GetModels()
	if err != nil {
		s.sendError(w, http.StatusInternalServerError, "Failed to fetch models", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"models": models,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// sendError sends a standardized error response
func (s *RESTServer) sendError(w http.ResponseWriter, statusCode int, message, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	errorResp := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Code:    statusCode,
		Message: message,
	}
	
	if details != "" {
		errorResp.Message += ": " + details
	}
	
	json.NewEncoder(w).Encode(errorResp)
}
