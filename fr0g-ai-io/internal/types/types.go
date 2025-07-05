package types

import "time"

// InputEvent represents an incoming event from external sources
type InputEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
	Priority  int                    `json:"priority"`
}

// InputEventResponse represents the response after processing an input event
type InputEventResponse struct {
	EventID     string                 `json:"event_id"`
	Processed   bool                   `json:"processed"`
	Actions     []OutputAction         `json:"actions"`
	Metadata    map[string]interface{} `json:"metadata"`
	ProcessedAt time.Time              `json:"processed_at"`
}

// OutputAction represents an action to be taken as a result of processing
type OutputAction struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Target    string                 `json:"target"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata"`
	Priority  int                    `json:"priority"`
	CreatedAt time.Time              `json:"created_at"`
}
