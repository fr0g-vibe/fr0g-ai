package outputs

import (
	"time"
)

// Re-export types for backward compatibility
type OutputCommand = types.OutputCommand
type OutputResult = types.OutputResult
type OutputProcessor = types.OutputProcessor

// Legacy types that may be used elsewhere
type OutputCommandLegacy struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Target    string            `json:"target"`
	Content   string            `json:"content"`
	Metadata  map[string]string `json:"metadata"`
	Priority  int               `json:"priority"`
	CreatedAt time.Time         `json:"created_at"`
}

type OutputResultLegacy struct {
	CommandID    string            `json:"command_id"`
	Success      bool              `json:"success"`
	ErrorMessage string            `json:"error_message,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	CompletedAt  time.Time         `json:"completed_at"`
}
