package outputs

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
	
	// Review and validation fields
	ReviewStatus   string     `json:"review_status,omitempty"`   // "pending", "approved", "rejected", "auto_approved"
	ReviewedBy     string     `json:"reviewed_by,omitempty"`     // User or system that reviewed
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty"`     // When review was completed
	ReviewComments string     `json:"review_comments,omitempty"` // Review feedback
	RequiresReview bool       `json:"requires_review"`           // Whether this command needs manual review
}

// OutputResult represents the result of executing an output command
type OutputResult struct {
	CommandID    string            `json:"command_id"`
	Success      bool              `json:"success"`
	ErrorMessage string            `json:"error_message,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	CompletedAt  time.Time         `json:"completed_at"`
	
	// Enhanced result tracking
	AttemptCount    int                    `json:"attempt_count"`
	ProcessingTime  time.Duration          `json:"processing_time"`
	ValidationIssues []ValidationIssue     `json:"validation_issues,omitempty"`
	DeliveryStatus  string                 `json:"delivery_status,omitempty"` // "sent", "delivered", "failed", "bounced"
	DeliveryDetails map[string]interface{} `json:"delivery_details,omitempty"`
}

// ValidationIssue represents a validation problem with an output command
type ValidationIssue struct {
	Field       string `json:"field"`
	Issue       string `json:"issue"`
	Severity    string `json:"severity"` // "error", "warning", "info"
	Suggestion  string `json:"suggestion,omitempty"`
}

// OutputReview represents a review of an output command
type OutputReview struct {
	ID             string    `json:"id"`
	CommandID      string    `json:"command_id"`
	ReviewerID     string    `json:"reviewer_id"`
	ReviewerType   string    `json:"reviewer_type"` // "human", "ai", "system"
	Status         string    `json:"status"`        // "approved", "rejected", "needs_changes"
	Comments       string    `json:"comments"`
	Confidence     float64   `json:"confidence,omitempty"` // For AI reviews
	ReviewCriteria []string  `json:"review_criteria"`
	CreatedAt      time.Time `json:"created_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

// OutputMetrics represents metrics for output processing
type OutputMetrics struct {
	TotalCommands     int64         `json:"total_commands"`
	SuccessfulOutputs int64         `json:"successful_outputs"`
	FailedOutputs     int64         `json:"failed_outputs"`
	PendingReviews    int64         `json:"pending_reviews"`
	AverageProcessTime time.Duration `json:"average_process_time"`
	ByType            map[string]int64 `json:"by_type"`
	ByStatus          map[string]int64 `json:"by_status"`
	LastUpdated       time.Time     `json:"last_updated"`
}

// OutputProcessor defines the interface for output processors
type OutputProcessor interface {
	Process(command *OutputCommand) (*OutputResult, error)
	GetType() string
	IsEnabled() bool
	
	// Enhanced processor capabilities
	ValidateCommand(command *OutputCommand) []ValidationIssue
	GetSupportedTargets() []string
	GetProcessingCapabilities() ProcessorCapabilities
}

// ProcessorCapabilities describes what an output processor can do
type ProcessorCapabilities struct {
	SupportsDeliveryTracking bool     `json:"supports_delivery_tracking"`
	SupportsRichContent      bool     `json:"supports_rich_content"`
	MaxContentLength         int      `json:"max_content_length"`
	SupportedContentTypes    []string `json:"supported_content_types"`
	RequiresAuthentication   bool     `json:"requires_authentication"`
	RateLimits              *RateLimit `json:"rate_limits,omitempty"`
}

// RateLimit defines rate limiting for a processor
type RateLimit struct {
	RequestsPerMinute int `json:"requests_per_minute"`
	RequestsPerHour   int `json:"requests_per_hour"`
	RequestsPerDay    int `json:"requests_per_day"`
}

// OutputReviewer defines the interface for output reviewers
type OutputReviewer interface {
	ReviewCommand(command *OutputCommand) (*OutputReview, error)
	GetReviewerType() string
	GetReviewCriteria() []string
	CanReviewType(outputType string) bool
}

// OutputValidator defines the interface for output validators
type OutputValidator interface {
	ValidateCommand(command *OutputCommand) []ValidationIssue
	GetValidationRules() []ValidationRule
	GetSupportedTypes() []string
}

// ValidationRule represents a validation rule
type ValidationRule struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Enabled     bool   `json:"enabled"`
}
