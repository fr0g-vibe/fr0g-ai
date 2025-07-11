package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/community"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/config"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/middleware"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/persona"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/registry"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/types"
	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// Server holds the HTTP server configuration and dependencies
type Server struct {
	config           *config.Config
	service          *persona.Service
	communityService *community.Service
	registryClient   *registry.RegistryClient
	server           *http.Server
}

// NewServer creates a new HTTP server instance
func NewServer(cfg *config.Config, service *persona.Service, registryClient *registry.RegistryClient) *Server {
	s := &Server{
		config:         cfg,
		service:        service,
		registryClient: registryClient,
	}
	if service != nil {
		s.communityService = community.NewService(service.GetStorage())
	}
	return s
}

// ServeHTTP implements http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", s.healthHandler)

	// Persona endpoints
	mux.HandleFunc("/personas", s.personasHandler)
	mux.HandleFunc("/personas/", s.personaHandler)

	// Identity endpoints
	mux.HandleFunc("/identities", s.identitiesHandler)
	mux.HandleFunc("/identities/", s.identityHandler)

	// Community endpoints
	mux.HandleFunc("/communities", s.communitiesHandler)
	mux.HandleFunc("/communities/", s.communityHandler)
	mux.HandleFunc("/communities/generate", s.generateCommunityHandler)

	mux.ServeHTTP(w, r)
}

// Start starts the HTTP server with graceful shutdown support
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", s.healthHandler)

	// Persona endpoints
	mux.HandleFunc("/personas", s.personasHandler)
	mux.HandleFunc("/personas/", s.personaHandler)

	// Identity endpoints
	mux.HandleFunc("/identities", s.identitiesHandler)
	mux.HandleFunc("/identities/", s.identityHandler)

	// Community endpoints
	mux.HandleFunc("/communities", s.communitiesHandler)
	mux.HandleFunc("/communities/", s.communityHandler)
	mux.HandleFunc("/communities/generate", s.generateCommunityHandler)

	// Apply middleware
	var handler http.Handler = mux

	// Add CORS middleware
	handler = middleware.CORSMiddleware(handler)

	// Add authentication middleware if enabled
	if s.config.Security.EnableAuth {
		handler = middleware.AuthMiddleware(s.config.Security.APIKey)(handler)
	}

	s.server = &http.Server{
		Addr:         ":" + s.config.HTTP.Port,
		Handler:      handler,
		ReadTimeout:  s.config.HTTP.ReadTimeout,
		WriteTimeout: s.config.HTTP.WriteTimeout,
	}

	if s.config.HTTP.EnableTLS {
		return s.server.ListenAndServeTLS(s.config.HTTP.CertFile, s.config.HTTP.KeyFile)
	}

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	// Deregister from service registry
	if s.registryClient != nil {
		if err := s.registryClient.Shutdown(); err != nil {
			fmt.Printf("Warning: Failed to deregister from service registry: %v\n", err)
		}
	}
	
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}

// handleError provides consistent error response handling
func (s *Server) handleError(w http.ResponseWriter, err error, defaultStatus int) {
	w.Header().Set("Content-Type", "application/json")
	
	if validationErr, ok := err.(sharedconfig.ValidationErrors); ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "validation_failed",
			"message": "Input validation failed",
			"details": validationErr,
		})
		return
	}

	// Handle validation error strings (reject whitespace-only inputs)
	if err != nil && (err.Error() == "name is required" || err.Error() == "persona_id is required") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "validation_failed",
			"message": err.Error(),
		})
		return
	}

	// Handle other specific error types here
	w.WriteHeader(defaultStatus)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   "request_failed",
		"message": err.Error(),
	})
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Enhanced health check
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0.0", // This could be injected at build time
		"storage":   s.config.Storage.Type,
	}

	// Check storage health safely
	var personaCount int
	var isHealthy = true
	
	if s.service != nil {
		if personas, err := s.service.ListPersonas(); err != nil {
			health["status"] = "unhealthy"
			health["storage_error"] = err.Error()
			isHealthy = false
			personaCount = 0
		} else {
			personaCount = len(personas)
		}
	} else {
		health["status"] = "unhealthy"
		health["storage_error"] = "service not initialized"
		isHealthy = false
		personaCount = 0
	}

	health["persona_count"] = personaCount
	
	// Add gRPC reflection status for debugging - check environment variable at runtime
	enableReflection := os.Getenv("GRPC_ENABLE_REFLECTION") == "true"
	
	// Also check config if environment variable is not set
	if !enableReflection && s.config != nil && s.config.GRPC.EnableReflection {
		enableReflection = true
	}
	
	if enableReflection {
		health["grpc_reflection"] = "enabled"
		health["grpc_reflection_warning"] = "Reflection should be disabled in production"
	} else {
		health["grpc_reflection"] = "disabled"
	}
	
	// Set appropriate HTTP status code
	if !isHealthy {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	
	json.NewEncoder(w).Encode(health)
}

func (s *Server) personasHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		personas, err := s.service.ListPersonas()
		if err != nil {
			http.Error(w, "Failed to list personas", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(personas)

	case http.MethodPost:
		var p types.Persona
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := s.service.CreatePersona(&p); err != nil {
			// Check if it's a validation error
			if validationErr, ok := err.(sharedconfig.ValidationErrors); ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error":   "Validation failed",
					"details": validationErr,
				})
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) personaHandler(w http.ResponseWriter, r *http.Request) {
	// Extract persona ID from URL path
	id := r.URL.Path[len("/personas/"):]
	if id == "" {
		http.Error(w, "Persona ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		p, err := s.service.GetPersona(id)
		if err != nil {
			http.Error(w, "Persona not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)

	case http.MethodPut:
		var p types.Persona
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := s.service.UpdatePersona(id, p); err != nil {
			// Check if it's a validation error
			if validationErr, ok := err.(sharedconfig.ValidationErrors); ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error":   "Validation failed",
					"details": validationErr,
				})
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)

	case http.MethodDelete:
		if err := s.service.DeletePersona(id); err != nil {
			http.Error(w, "Persona not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) identitiesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Parse query parameters for filtering
		filter := &types.IdentityFilter{}
		if personaID := r.URL.Query().Get("persona_id"); personaID != "" {
			filter.PersonaID = personaID
		}
		if search := r.URL.Query().Get("search"); search != "" {
			filter.Search = search
		}
		if isActiveStr := r.URL.Query().Get("is_active"); isActiveStr != "" {
			if isActive := isActiveStr == "true"; isActiveStr == "true" || isActiveStr == "false" {
				filter.IsActive = &isActive
			}
		}

		// Get identities from storage
		identities, err := s.service.GetStorage().ListIdentities(filter)
		if err != nil {
			http.Error(w, "Failed to list identities", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(identities)

	case http.MethodPost:
		var identity types.Identity
		if err := json.NewDecoder(r.Body).Decode(&identity); err != nil {
			s.handleError(w, err, http.StatusBadRequest)
			return
		}

		// Validate required fields (reject empty and whitespace-only strings)
		if len(strings.TrimSpace(identity.Name)) == 0 {
			s.handleError(w, fmt.Errorf("name is required"), http.StatusBadRequest)
			return
		}
		if len(strings.TrimSpace(identity.PersonaId)) == 0 {
			s.handleError(w, fmt.Errorf("persona_id is required"), http.StatusBadRequest)
			return
		}

		// Validate that the persona exists
		if _, err := s.service.GetPersona(identity.PersonaId); err != nil {
			s.handleError(w, fmt.Errorf("persona_id '%s' does not exist", identity.PersonaId), http.StatusBadRequest)
			return
		}

		if err := s.service.CreateIdentity(&identity); err != nil {
			s.handleError(w, err, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(identity)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) identityHandler(w http.ResponseWriter, r *http.Request) {
	// Extract identity ID from URL path
	path := r.URL.Path[len("/identities/"):]
	if path == "" {
		http.Error(w, "Identity ID required", http.StatusBadRequest)
		return
	}

	// Handle special endpoints
	if path == "with-persona" {
		// Get all identities with personas
		identities, err := s.service.GetStorage().ListIdentities(nil)
		if err != nil {
			http.Error(w, "Failed to list identities", http.StatusInternalServerError)
			return
		}

		var identitiesWithPersonas []types.IdentityWithPersona
		for _, identity := range identities {
			persona, err := s.service.GetPersona(identity.PersonaId)
			if err != nil {
				continue // Skip if persona not found
			}
			identitiesWithPersonas = append(identitiesWithPersonas, types.IdentityWithPersona{
				Identity: identity,
				Persona:  persona,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(identitiesWithPersonas)
		return
	}

	switch r.Method {
	case http.MethodGet:
		identity, err := s.service.GetIdentity(path)
		if err != nil {
			http.Error(w, "Identity not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(identity)

	case http.MethodPut:
		var identity types.Identity
		if err := json.NewDecoder(r.Body).Decode(&identity); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := s.service.UpdateIdentity(path, identity); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(identity)

	case http.MethodDelete:
		if err := s.service.DeleteIdentity(path); err != nil {
			http.Error(w, "Identity not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) communitiesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Parse query parameters for filtering
		filter := &types.CommunityFilter{}
		if communityType := r.URL.Query().Get("type"); communityType != "" {
			filter.Type = communityType
		}
		if search := r.URL.Query().Get("search"); search != "" {
			filter.Search = search
		}

		// Get communities from storage
		communities, err := s.service.GetStorage().ListCommunities(filter)
		if err != nil {
			http.Error(w, "Failed to list communities", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(communities)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) communityHandler(w http.ResponseWriter, r *http.Request) {
	// Extract community ID from URL path
	path := r.URL.Path[len("/communities/"):]
	if path == "" {
		http.Error(w, "Community ID required", http.StatusBadRequest)
		return
	}

	// Handle special endpoints
	if path == "stats" {
		http.Error(w, "Community stats endpoint should be accessed via /communities/{id}/stats", http.StatusBadRequest)
		return
	}

	// Handle stats endpoint with proper path parsing
	if len(path) > 6 && path[len(path)-6:] == "/stats" {
		communityId := path[:len(path)-6]

		// Create community service instance
		communityService := s.getCommunityService()
		stats, err := communityService.GetCommunityStats(communityId)
		if err != nil {
			http.Error(w, "Failed to get community stats", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
		return
	}

	switch r.Method {
	case http.MethodGet:
		community, err := s.service.GetStorage().GetCommunity(path)
		if err != nil {
			http.Error(w, "Community not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(community)

	case http.MethodPut:
		var community types.Community
		if err := json.NewDecoder(r.Body).Decode(&community); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := s.service.GetStorage().UpdateCommunity(path, community); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(community)

	case http.MethodDelete:
		if err := s.service.GetStorage().DeleteCommunity(path); err != nil {
			http.Error(w, "Community not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) generateCommunityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name             string                          `json:"name"`
		Description      string                          `json:"description"`
		Type             string                          `json:"type"`
		TargetSize       int                             `json:"target_size"`
		GenerationConfig types.CommunityGenerationConfig `json:"generation_config"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create community service instance
	communityService := s.getCommunityService()

	// Generate community
	community, err := communityService.GenerateCommunity(
		req.GenerationConfig,
		req.Name,
		req.Description,
		req.Type,
		req.TargetSize,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate community: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(community)
}

// StartServer starts the HTTP API server (legacy function for backward compatibility)
func StartServer(port string) error {
	cfg := &config.Config{
		HTTP: sharedconfig.HTTPConfig{
			Port:         port,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Security: sharedconfig.SecurityConfig{
			EnableAuth: false,
		},
		Storage: sharedconfig.StorageConfig{
			Type: "memory",
		},
	}

	// Create a default service for legacy usage
	store := storage.NewMemoryStorage()
	service := persona.NewService(store)
	server := NewServer(cfg, service, nil) // No registry client for legacy usage
	return server.Start()
}

// getCommunityService returns the community service instance
func (s *Server) getCommunityService() *community.Service {
	return s.communityService
}

// StartServerWithConfig starts the HTTP server with full configuration
func StartServerWithConfig(cfg *config.Config, service *persona.Service, registryClient *registry.RegistryClient) error {
	if service == nil {
		return fmt.Errorf("persona service is required")
	}
	server := NewServer(cfg, service, registryClient)
	return server.Start()
}
