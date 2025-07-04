package api

import "fmt"

// isValidRole checks if the role is one of the allowed values
func isValidRole(role string) bool {
	validRoles := []string{"system", "user", "assistant", "function"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// ValidateChatCompletionRequest validates a chat completion request
func ValidateChatCompletionRequest(req interface{}) error {
	// This function can be used by both REST and gRPC handlers
	// Implementation depends on the request type
	return nil
}

// ValidateMessage validates a single chat message
func ValidateMessage(role, content string) error {
	if role == "" {
		return fmt.Errorf("role is required")
	}
	if content == "" {
		return fmt.Errorf("content is required")
	}
	if len(content) > 10000 {
		return fmt.Errorf("content too long (max 10000 characters)")
	}
	if !isValidRole(role) {
		return fmt.Errorf("invalid role: %s", role)
	}
	return nil
}
