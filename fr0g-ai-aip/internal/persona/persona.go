// Package persona provides AI persona management functionality.
//
// This package implements the core persona service that handles creation,
// retrieval, updating, and deletion of AI subject matter experts. Each
// persona consists of a name, topic, system prompt, and optional context
// and RAG (Retrieval-Augmented Generation) documents.
//
// The service supports both in-memory and file-based storage backends,
// with comprehensive validation and error handling. All operations are
// safe for concurrent use.
//
// Example usage:
//
//	storage := storage.NewMemoryStorage()
//	service := persona.NewService(storage)
//
//	p := &types.Persona{
//		Name:   "Go Expert",
//		Topic:  "Golang Programming",
//		Prompt: "You are an expert Go programmer with deep knowledge of best practices.",
//		Context: map[string]string{
//			"experience": "10 years",
//			"specialty":  "backend development",
//		},
//	}
//
//	err := service.CreatePersona(p)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// The package also provides identity management functionality, allowing
// creation of persona instances with rich demographic and behavioral
// attributes for community simulation and analysis.
package persona

import (
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/types"
)

// Persona is a type alias for backward compatibility.
// Use types.Persona directly in new code.
type Persona = types.Persona

// Global service instance for backward compatibility
var defaultService *Service

// SetDefaultService sets the default service instance
func SetDefaultService(service *Service) {
	defaultService = service
}

// Legacy functions for backward compatibility
func CreatePersona(p *types.Persona) error {
	if defaultService == nil {
		defaultService = NewService(storage.NewMemoryStorage())
	}
	return defaultService.CreatePersona(p)
}

func GetPersona(id string) (types.Persona, error) {
	if defaultService == nil {
		defaultService = NewService(storage.NewMemoryStorage())
	}
	return defaultService.GetPersona(id)
}

func ListPersonas() []types.Persona {
	if defaultService == nil {
		defaultService = NewService(storage.NewMemoryStorage())
	}
	personas, _ := defaultService.ListPersonas()
	return personas
}

func DeletePersona(id string) error {
	if defaultService == nil {
		defaultService = NewService(storage.NewMemoryStorage())
	}
	return defaultService.DeletePersona(id)
}

func UpdatePersona(id string, p types.Persona) error {
	if defaultService == nil {
		defaultService = NewService(storage.NewMemoryStorage())
	}
	return defaultService.UpdatePersona(id, p)
}
