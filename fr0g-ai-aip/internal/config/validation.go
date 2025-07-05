package config

import (
	"fmt"

	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
)

// ValidationError represents a validation error (alias to shared type)
type ValidationError = sharedconfig.ValidationError

// ValidationErrors represents multiple validation errors (alias to shared type)
type ValidationErrors = sharedconfig.ValidationErrors

// ValidatePersonaData validates persona-specific data using shared validation
func ValidatePersonaData(name, topic, prompt string) sharedconfig.ValidationErrors {
	var errors []sharedconfig.ValidationError
	
	// Validate persona name
	if err := sharedconfig.ValidateRequired(name, "persona.name"); err != nil {
		errors = append(errors, *err)
	}
	if err := sharedconfig.ValidateLength(name, 1, 100, "persona.name"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate persona topic
	if err := sharedconfig.ValidateRequired(topic, "persona.topic"); err != nil {
		errors = append(errors, *err)
	}
	if err := sharedconfig.ValidateLength(topic, 1, 200, "persona.topic"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate persona prompt
	if err := sharedconfig.ValidateRequired(prompt, "persona.prompt"); err != nil {
		errors = append(errors, *err)
	}
	if err := sharedconfig.ValidateLength(prompt, 1, 8000, "persona.prompt"); err != nil {
		errors = append(errors, *err)
	}
	
	return sharedconfig.ValidationErrors(errors)
}

// ValidateIdentityData validates identity-specific data using shared validation
func ValidateIdentityData(name, description string, attributes map[string]string) sharedconfig.ValidationErrors {
	var errors []sharedconfig.ValidationError
	
	// Validate identity name
	if err := sharedconfig.ValidateRequired(name, "identity.name"); err != nil {
		errors = append(errors, *err)
	}
	if err := sharedconfig.ValidateLength(name, 1, 100, "identity.name"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate identity description
	if err := sharedconfig.ValidateLength(description, 0, 1000, "identity.description"); err != nil {
		errors = append(errors, *err)
	}
	
	// Validate attributes map
	if len(attributes) > 50 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "identity.attributes",
			Message: "too many attributes (max 50)",
		})
	}
	
	for key, value := range attributes {
		if err := sharedconfig.ValidateLength(key, 1, 50, fmt.Sprintf("identity.attributes.%s.key", key)); err != nil {
			errors = append(errors, *err)
		}
		if err := sharedconfig.ValidateLength(value, 0, 500, fmt.Sprintf("identity.attributes.%s.value", key)); err != nil {
			errors = append(errors, *err)
		}
	}
	
	return sharedconfig.ValidationErrors(errors)
}
