package types

import (
	"time"
)

// OutputCommand represents a command to send output via I/O channels
type OutputCommand struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"` // "sms", "voice", "irc", "discord", "email"
	Target    string            `json:"target"`
	Content   string            `json:"content"`
	Metadata  map[string]string `json:"metadata"`
	Priority  int               `json:"priority"`
	CreatedAt time.Time         `json:"created_at"`
}

// OutputResult represents the result of executing an output command
type OutputResult struct {
	CommandID    string            `json:"command_id"`
	Success      bool              `json:"success"`
	ErrorMessage string            `json:"error_message,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	CompletedAt  time.Time         `json:"completed_at"`
}

// OutputProcessor defines the interface for output processors
type OutputProcessor interface {
	Process(command *OutputCommand) (*OutputResult, error)
	GetType() string
	IsEnabled() bool
}
