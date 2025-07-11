package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/client"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/config"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-bridge/internal/models"
	"github.com/gorilla/mux"
)

// RESTServer handles REST API requests
type RESTServer struct {
	client      *client.OpenWebUIClient
	router      *mux.Router
	config      *config.Config
	rateLimiter *RateLimiter
}

// RateLimiter implements a simple rate limiting mechanism
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    requestsPerMinute,
		window:   time.Minute,
	}
}

// Allow checks if a request from the given IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	// Clean old requests
	if requests, exists := rl.requests[ip]; exists {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if now.Sub(reqTime) < rl.window {
				validRequests = append(validRequests, reqTime)
			}
		}
		rl.requests[ip] = validRequests
	}

	// Check if under limit
	if len(rl.requests[ip]) >= rl.limit {
		return false
	}

	// Add current request
	rl.requests[ip] = append(rl.requests[ip], now)
	return true
}

// NewRESTServer creates a new REST server
func NewRESTServer(openWebUIClient *client.OpenWebUIClient, cfg *config.Config) *RESTServer {
	server := &RESTServer{
		client:      openWebUIClient,
		router:      mux.NewRouter(),
		config:      cfg,
		rateLimiter: NewRateLimiter(cfg.Security.RateLimitRPM),
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures the REST API routes
func (s *RESTServer) setupRoutes() {
	// Health check endpoint
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")

	// Chat completion endpoint (OpenAI compatible)
	s.router.HandleFunc("/v1/chat/completions", s.handleChatCompletion).Methods("POST")
	s.router.HandleFunc("/api/chat/completions", s.handleChatCompletion).Methods("POST") // Legacy support

	// Legacy simple chat endpoint
	s.router.HandleFunc("/api/v1/chat", s.handleSimpleChat).Methods("POST")

	// Models endpoint
	s.router.HandleFunc("/api/v1/models", s.handleModels).Methods("GET")

	// Add middleware
	s.router.Use(s.securityHeadersMiddleware)
	s.router.Use(s.requestSizeLimitMiddleware)
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.rateLimitMiddleware)
	if s.config.Security.RequireAPIKey {
		s.router.Use(s.apiKeyMiddleware)
	}
	if s.config.Security.EnableCORS {
		s.router.Use(s.corsMiddleware)
	}
}

// GetRouter returns the configured router
func (s *RESTServer) GetRouter() *mux.Router {
	return s.router
}

// handleHealth handles health check requests with comprehensive validation
func (s *RESTServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Initialize response with required fields
	response := models.HealthResponse{
		Time:    time.Now().UTC(), // Use UTC for consistency
		Version: "1.0.0",
	}

	// Check OpenWebUI health
	err := s.client.HealthCheck(ctx)

	if err != nil {
		response.Status = "unhealthy"
		response.Error = s.sanitizeErrorMessage(err.Error())
		response.Details = s.buildHealthDetails(false, err)
		log.Printf("Health check failed: %v", err)
	} else {
		response.Status = "healthy"
		response.Details = s.buildHealthDetails(true, nil)
	}

	// Validate response format and get appropriate status code
	statusCode, validationErr := response.ValidateForStatusCode()
	if validationErr != nil {
		log.Printf("Health response validation failed: %v", validationErr)
		s.writeError(w, http.StatusInternalServerError, "Health check validation failed", validationErr)
		return
	}

	// Set headers and status code
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("X-Health-Check-Time", response.Time.Format(time.RFC3339))
	w.WriteHeader(statusCode)

	// Encode and send response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode health response: %v", err)
		return
	}

	// Log health check result
	log.Printf("Health check completed: status=%s, code=%d", response.Status, statusCode)
}

// buildHealthDetails creates the details map for health responses
func (s *RESTServer) buildHealthDetails(healthy bool, err error) map[string]interface{} {
	details := map[string]interface{}{
		"service":         "fr0g-ai-bridge",
		"openwebui_url":   s.config.OpenWebUI.BaseURL,
		"has_api_key":     s.config.OpenWebUI.APIKey != "",
		"timeout_seconds": s.config.OpenWebUI.Timeout,
		"timestamp":       time.Now().UTC().Format(time.RFC3339),
	}

	if healthy {
		details["authenticated"] = true
		details["connection"] = "ok"
	} else {
		details["connection"] = "failed"
		if err != nil {
			details["error_type"] = s.categorizeError(err)
		}
	}

	return details
}

// sanitizeErrorMessage removes sensitive information from error messages
func (s *RESTServer) sanitizeErrorMessage(errMsg string) string {
	// Remove potential sensitive information
	sanitized := strings.ReplaceAll(errMsg, s.config.OpenWebUI.APIKey, "[REDACTED]")

	// Truncate if too long
	if len(sanitized) > 500 {
		sanitized = sanitized[:497] + "..."
	}

	return sanitized
}

// categorizeError categorizes errors for better debugging
func (s *RESTServer) categorizeError(err error) string {
	errStr := strings.ToLower(err.Error())

	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline") {
		return "timeout"
	}
	if strings.Contains(errStr, "connection") || strings.Contains(errStr, "dial") {
		return "connection"
	}
	if strings.Contains(errStr, "unauthorized") || strings.Contains(errStr, "401") {
		return "authentication"
	}
	if strings.Contains(errStr, "not found") || strings.Contains(errStr, "404") {
		return "not_found"
	}

	return "unknown"
}

// handleChatCompletion handles OpenAI-compatible chat completion requests
func (s *RESTServer) handleChatCompletion(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req models.ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := ValidateChatCompletionRequest(&ChatCompletionRequest{
		Model:       req.Model,
		Messages:    convertToValidationMessages(req.Messages),
		Temperature: req.Temperature,
		MaxTokens:   convertMaxTokens(req.MaxTokens),
		Stream:      req.Stream != nil && *req.Stream,
	}); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	// Validate persona prompt if provided
	if err := ValidatePersonaPrompt(&req.PersonaPrompt); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid persona prompt", err)
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

// SimpleChatRequest represents a simple chat request
type SimpleChatRequest struct {
	Message string `json:"message"`
	Model   string `json:"model,omitempty"`
}

// SimpleChatResponse represents a simple chat response
type SimpleChatResponse struct {
	Response  string    `json:"response"`
	Model     string    `json:"model"`
	Timestamp time.Time `json:"timestamp"`
}

// handleSimpleChat handles simple chat requests (legacy endpoint)
func (s *RESTServer) handleSimpleChat(w http.ResponseWriter, r *http.Request) {
	var req SimpleChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request format", err)
		return
	}

	if req.Message == "" {
		s.writeError(w, http.StatusBadRequest, "Message is required", nil)
		return
	}

	// Validate message length and content
	if len(req.Message) > 10000 {
		s.writeError(w, http.StatusBadRequest, "Message too long (max 10000 characters)", nil)
		return
	}

	// Sanitize input
	req.Message = strings.TrimSpace(req.Message)
	if req.Model != "" {
		req.Model = strings.TrimSpace(req.Model)
	}

	// Bridge request to OpenWebUI using legacy method
	response, err := s.client.SendMessage(req.Message, req.Model)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to process request", err)
		return
	}

	chatResp := SimpleChatResponse{
		Response:  response,
		Model:     req.Model,
		Timestamp: time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chatResp)
}

// handleModels handles requests for available AI models
func (s *RESTServer) handleModels(w http.ResponseWriter, r *http.Request) {
	models, err := s.client.GetModels()
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to fetch models", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"models":    models,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// convertToValidationMessages converts models.ChatMessage to validation.Message
func convertToValidationMessages(messages []models.ChatMessage) []Message {
	result := make([]Message, len(messages))
	for i, msg := range messages {
		result[i] = Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return result
}

// convertMaxTokens converts *int to *int32 for validation
func convertMaxTokens(maxTokens *int) *int32 {
	if maxTokens == nil {
		return nil
	}
	val := int32(*maxTokens)
	return &val
}

// writeError writes an error response
func (s *RESTServer) writeError(w http.ResponseWriter, statusCode int, message string, err error) {
	log.Printf("API Error: %s - %v", message, err)

	errorResp := models.ErrorResponse{
		Error: message,
		Code:  statusCode,
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
		origin := "*"
		if len(s.config.Security.AllowedOrigins) > 0 && s.config.Security.AllowedOrigins[0] != "*" {
			origin = s.config.Security.AllowedOrigins[0] // Simplified for now
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// rateLimitMiddleware implements rate limiting
func (s *RESTServer) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract client IP
		clientIP := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			clientIP = forwarded
		}

		if !s.rateLimiter.Allow(clientIP) {
			s.writeError(w, http.StatusTooManyRequests, "Rate limit exceeded", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs all requests
func (s *RESTServer) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		log.Printf("[%s] %s %s - %d - %v",
			r.Method, r.RequestURI, r.RemoteAddr, wrapped.statusCode, duration)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// securityHeadersMiddleware adds security headers to all responses
func (s *RESTServer) securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")
		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Enable XSS protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		// Enforce HTTPS (when behind a proxy)
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		// Prevent referrer leakage
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		// Content Security Policy
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'none'; object-src 'none'")

		next.ServeHTTP(w, r)
	})
}

// requestSizeLimitMiddleware limits the size of request bodies
func (s *RESTServer) requestSizeLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Limit request body size to 1MB
		r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
		next.ServeHTTP(w, r)
	})
}

// apiKeyMiddleware validates API keys
func (s *RESTServer) apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip API key check for health endpoint
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		// Extract API key from Authorization header or X-API-Key header
		var apiKey string
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			apiKey = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			apiKey = r.Header.Get("X-API-Key")
		}

		if apiKey == "" {
			s.writeError(w, http.StatusUnauthorized, "API key required", nil)
			return
		}

		// Validate API key
		if !s.isValidAPIKey(apiKey) {
			s.writeError(w, http.StatusUnauthorized, "Invalid API key", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// isValidAPIKey checks if the provided API key is valid
func (s *RESTServer) isValidAPIKey(apiKey string) bool {
	if len(s.config.Security.AllowedAPIKeys) == 0 {
		return true // No API keys configured, allow all
	}

	for _, validKey := range s.config.Security.AllowedAPIKeys {
		if apiKey == validKey {
			return true
		}
	}
	return false
}

// isValidRole checks if a role is valid
func isValidRole(role string) bool {
	validRoles := []string{"user", "assistant", "system", "function"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}
