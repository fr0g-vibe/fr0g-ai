package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/types"
	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents a collection of validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// Error implements the error interface
func (ve ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "validation failed"
	}
	
	if len(ve.Errors) == 1 {
		return fmt.Sprintf("validation failed: %s: %s", ve.Errors[0].Field, ve.Errors[0].Message)
	}
	
	return fmt.Sprintf("validation failed with %d errors", len(ve.Errors))
}

// ValidatePersona validates a persona struct using shared validation
func ValidatePersona(p *types.Persona) error {
	if p == nil {
		return sharedconfig.ValidationErrors{sharedconfig.ValidationError{
			Field:   "persona",
			Message: "persona cannot be nil",
		}}
	}

	var errors []sharedconfig.ValidationError

	// Validate required fields
	if err := sharedconfig.ValidateRequired(p.Name, "name"); err != nil {
		errors = append(errors, *err)
	}
	if err := sharedconfig.ValidateRequired(p.Topic, "topic"); err != nil {
		errors = append(errors, *err)
	}
	if err := sharedconfig.ValidateRequired(p.Prompt, "prompt"); err != nil {
		errors = append(errors, *err)
	}

	// Validate field lengths
	if err := sharedconfig.ValidateLength(p.Name, 1, 100, "name"); err != nil {
		errors = append(errors, *err)
	}
	if err := sharedconfig.ValidateLength(p.Topic, 1, 200, "topic"); err != nil {
		errors = append(errors, *err)
	}
	if err := sharedconfig.ValidateLength(p.Prompt, 1, 10000, "prompt"); err != nil {
		errors = append(errors, *err)
	}

	// Validate context keys and values
	for key, value := range p.Context {
		if err := sharedconfig.ValidateRequired(key, "context.key"); err != nil {
			errors = append(errors, *err)
		}
		if err := sharedconfig.ValidateLength(key, 1, 50, "context.key"); err != nil {
			errors = append(errors, *err)
		}
		if err := sharedconfig.ValidateLength(value, 0, 500, "context.value"); err != nil {
			errors = append(errors, *err)
		}
	}

	// Validate RAG field
	if p.Rag != nil {
		for i, rag := range p.Rag {
			if err := sharedconfig.ValidateRequired(rag, fmt.Sprintf("rag[%d]", i)); err != nil {
				errors = append(errors, *err)
			}
		}
	}

	if len(errors) > 0 {
		return sharedconfig.ValidationErrors(errors)
	}

	return nil
}

// ValidationMiddleware returns an HTTP middleware that validates request bodies
func ValidationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only validate POST and PUT requests with JSON content
		if (r.Method == http.MethodPost || r.Method == http.MethodPut) &&
			strings.Contains(r.Header.Get("Content-Type"), "application/json") {

			var p types.Persona
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				http.Error(w, "Invalid JSON format", http.StatusBadRequest)
				return
			}

			if err := ValidatePersona(&p); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				
				// Handle both ValidationErrors and regular errors
				if validationErrors, ok := err.(sharedconfig.ValidationErrors); ok {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error":   "Validation failed",
						"details": validationErrors,
					})
				} else {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"error":   "Validation failed",
						"details": err.Error(),
					})
				}
				return
			}

			// Store validated persona in request context for handlers to use
			// For now, we'll let the handler re-decode, but this could be optimized
		}

		next.ServeHTTP(w, r)
	}
}

// SanitizePersona sanitizes persona input by trimming whitespace
func SanitizePersona(p *types.Persona) {
	if p == nil {
		return
	}

	p.Name = strings.TrimSpace(p.Name)
	p.Topic = strings.TrimSpace(p.Topic)
	p.Prompt = strings.TrimSpace(p.Prompt)

	// Initialize context if nil
	if p.Context == nil {
		p.Context = make(map[string]string)
	}

	// Sanitize context
	for key, value := range p.Context {
		delete(p.Context, key)
		cleanKey := strings.TrimSpace(key)
		cleanValue := strings.TrimSpace(value)
		if cleanKey != "" {
			p.Context[cleanKey] = cleanValue
		}
	}

	// Sanitize RAG field
	if p.Rag != nil {
		var cleanRAG []string
		for _, rag := range p.Rag {
			if clean := strings.TrimSpace(rag); clean != "" {
				cleanRAG = append(cleanRAG, clean)
			}
		}
		p.Rag = cleanRAG
	}
}
