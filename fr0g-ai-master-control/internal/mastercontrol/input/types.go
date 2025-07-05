package input

import (
	"context"
	"time"
)

// WebhookProcessor defines the interface for processing webhooks
type WebhookProcessor interface {
	ProcessWebhook(ctx context.Context, request *WebhookRequest) (*WebhookResponse, error)
	GetTag() string
	GetDescription() string
}

// WebhookRequest represents an incoming webhook request
type WebhookRequest struct {
	ID          string                 `json:"id"`
	Source      string                 `json:"source"`
	Tag         string                 `json:"tag"`
	Timestamp   time.Time              `json:"timestamp"`
	Headers     map[string]string      `json:"headers"`
	Body        interface{}            `json:"body"`
	Metadata    map[string]interface{} `json:"metadata"`
	ProcessedAt *time.Time             `json:"processed_at,omitempty"`
}

// WebhookResponse represents the response to a webhook
type WebhookResponse struct {
	Success   bool                   `json:"success"`
	Message   string                 `json:"message"`
	RequestID string                 `json:"request_id"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// WebhookConfig holds webhook manager configuration
type WebhookConfig struct {
	Port           int           `yaml:"port"`
	Host           string        `yaml:"host"`
	ReadTimeout    time.Duration `yaml:"read_timeout"`
	WriteTimeout   time.Duration `yaml:"write_timeout"`
	MaxRequestSize int64         `yaml:"max_request_size"`
	EnableLogging  bool          `yaml:"enable_logging"`
	AllowedOrigins []string      `yaml:"allowed_origins"`
}

// AIPersonaCommunityClient defines the interface for AI persona community interactions
type AIPersonaCommunityClient interface {
	CreateCommunity(ctx context.Context, topic string, personaCount int) (*Community, error)
	SubmitForReview(ctx context.Context, communityID string, content string) (*CommunityReview, error)
	GetReviewStatus(ctx context.Context, reviewID string) (*CommunityReview, error)
	GetCommunityMembers(ctx context.Context, communityID string) ([]PersonaInfo, error)
}

// Fr0gIOClient defines the interface for fr0g-ai-io service interactions
type Fr0gIOClient interface {
	SendOutputCommand(ctx context.Context, command *OutputCommand) (*OutputResponse, error)
	SendThreatAnalysisResult(ctx context.Context, result *ThreatAnalysisResult) error
	GetServiceStatus(ctx context.Context) (*ServiceStatus, error)
}

// Fr0gIOInputHandler defines the interface for handling input events from fr0g-ai-io
type Fr0gIOInputHandler interface {
	HandleInputEvent(ctx context.Context, event *InputEvent) (*InputEventResponse, error)
	HandleSMSMessage(ctx context.Context, message *SMSMessage) error
	HandleVoiceMessage(ctx context.Context, message *VoiceMessage) error
	HandleIRCMessage(ctx context.Context, message *IRCMessage) error
	HandleDiscordMessage(ctx context.Context, message *DiscordMessage) error
}

// Community represents an AI persona community
type Community struct {
	ID        string        `json:"id"`
	Topic     string        `json:"topic"`
	Members   []PersonaInfo `json:"members"`
	CreatedAt time.Time     `json:"created_at"`
	Status    string        `json:"status"`
}

// PersonaInfo represents information about an AI persona
type PersonaInfo struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Expertise   []string `json:"expertise"`
	Description string   `json:"description"`
	Model       string   `json:"model"`
}

// CommunityReview represents the AI community's review of content
type CommunityReview struct {
	ReviewID        string                 `json:"review_id"`
	Topic           string                 `json:"topic"`
	Content         string                 `json:"content"`
	PersonaReviews  []PersonaReview        `json:"persona_reviews"`
	Consensus       *Consensus             `json:"consensus"`
	Sentiment       *SentimentAnalysis     `json:"sentiment,omitempty"`
	Recommendations []string               `json:"recommendations"`
	Metadata        map[string]interface{} `json:"metadata"`
	CreatedAt       time.Time              `json:"created_at"`
	CompletedAt     *time.Time             `json:"completed_at,omitempty"`
}

// PersonaReview represents an individual AI persona's review
type PersonaReview struct {
	PersonaID   string                 `json:"persona_id"`
	PersonaName string                 `json:"persona_name"`
	Expertise   []string               `json:"expertise"`
	Review      string                 `json:"review"`
	Score       float64                `json:"score"`
	Confidence  float64                `json:"confidence"`
	Tags        []string               `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
	Timestamp   time.Time              `json:"timestamp"`
}

// Consensus represents the community consensus
type Consensus struct {
	OverallScore    float64  `json:"overall_score"`
	Agreement       float64  `json:"agreement"`
	Recommendation  string   `json:"recommendation"`
	KeyPoints       []string `json:"key_points"`
	Dissenting      []string `json:"dissenting,omitempty"`
	ConfidenceLevel float64  `json:"confidence_level"`
}

// SentimentAnalysis represents sentiment analysis results
type SentimentAnalysis struct {
	Overall      string             `json:"overall"`
	Score        float64            `json:"score"`
	Emotions     map[string]float64 `json:"emotions"`
	Toxicity     float64            `json:"toxicity"`
	Subjectivity float64            `json:"subjectivity"`
}

// InputEvent represents an input event from fr0g-ai-io
type InputEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"` // "sms", "voice", "irc", "discord"
	Source    string                 `json:"source"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
	Priority  int                    `json:"priority"`
}

// InputEventResponse represents the response to an input event
type InputEventResponse struct {
	EventID     string                 `json:"event_id"`
	Processed   bool                   `json:"processed"`
	Actions     []OutputAction         `json:"actions"`
	Analysis    *ThreatAnalysisResult  `json:"analysis,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
	ProcessedAt time.Time              `json:"processed_at"`
}


// OutputCommand represents a command to send to fr0g-ai-io for output processing
type OutputCommand struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"` // "sms", "voice", "irc", "discord"
	Target   string                 `json:"target"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
	Priority int                    `json:"priority"`
	
	// Enhanced output command fields for review workflow
	ReviewStatus     string     `json:"review_status,omitempty"`     // "pending", "approved", "rejected", "auto_approved"
	ReviewedBy       string     `json:"reviewed_by,omitempty"`       // Reviewer identifier
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty"`       // Review timestamp
	ReviewComments   string     `json:"review_comments,omitempty"`   // Review feedback
	RequiresReview   bool       `json:"requires_review"`             // Whether manual review is needed
	ReviewDeadline   *time.Time `json:"review_deadline,omitempty"`   // When review must be completed
	AutoApprovalRule string     `json:"auto_approval_rule,omitempty"` // Rule that auto-approved this command
}

// OutputAction represents an action to be taken by fr0g-ai-io
type OutputAction struct {
	Type     string                 `json:"type"`
	Target   string                 `json:"target"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
}

// OutputResponse represents the response from fr0g-ai-io after processing an output command
type OutputResponse struct {
	CommandID string                 `json:"command_id"`
	Success   bool                   `json:"success"`
	Message   string                 `json:"message"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
	
	// Enhanced response fields for better tracking
	ProcessingTime   time.Duration          `json:"processing_time,omitempty"`
	ValidationIssues []ValidationIssue      `json:"validation_issues,omitempty"`
	DeliveryStatus   string                 `json:"delivery_status,omitempty"` // "queued", "sent", "delivered", "failed"
	RetryCount       int                    `json:"retry_count,omitempty"`
	NextRetryAt      *time.Time             `json:"next_retry_at,omitempty"`
}

// ValidationIssue represents a validation problem with an output command
type ValidationIssue struct {
	Field      string `json:"field"`
	Issue      string `json:"issue"`
	Severity   string `json:"severity"` // "error", "warning", "info"
	Suggestion string `json:"suggestion,omitempty"`
}

// OutputReviewRequest represents a request to review an output command
type OutputReviewRequest struct {
	CommandID      string    `json:"command_id"`
	ReviewerID     string    `json:"reviewer_id"`
	ReviewerType   string    `json:"reviewer_type"` // "human", "ai", "system"
	RequestedAt    time.Time `json:"requested_at"`
	Priority       int       `json:"priority"`
	ReviewDeadline time.Time `json:"review_deadline"`
	Context        string    `json:"context,omitempty"`
}

// OutputReviewResponse represents the response to an output review
type OutputReviewResponse struct {
	CommandID      string     `json:"command_id"`
	ReviewerID     string     `json:"reviewer_id"`
	Status         string     `json:"status"`         // "approved", "rejected", "needs_changes"
	Comments       string     `json:"comments"`
	Confidence     float64    `json:"confidence,omitempty"` // For AI reviews
	ReviewCriteria []string   `json:"review_criteria"`
	CompletedAt    time.Time  `json:"completed_at"`
	Modifications  []string   `json:"modifications,omitempty"` // Suggested changes
}

// ThreatAnalysisResult represents the result of threat analysis processing
type ThreatAnalysisResult struct {
	EventID       string                 `json:"event_id"`
	ThreatLevel   string                 `json:"threat_level"` // "low", "medium", "high", "critical"
	ThreatScore   float64                `json:"threat_score"` // 0.0-1.0
	ThreatTypes   []string               `json:"threat_types"`
	Indicators    []ThreatIndicator      `json:"indicators"`
	Mitigation    []string               `json:"mitigation"`
	Confidence    float64                `json:"confidence"`
	Analysis      string                 `json:"analysis"`
	Metadata      map[string]interface{} `json:"metadata"`
	AnalyzedAt    time.Time              `json:"analyzed_at"`
	RecommendedActions []OutputAction    `json:"recommended_actions"`
}

// ThreatIndicator represents a specific threat indicator
type ThreatIndicator struct {
	Type        string  `json:"type"`
	Value       string  `json:"value"`
	Confidence  float64 `json:"confidence"`
	Description string  `json:"description"`
}

// ServiceStatus represents the status of fr0g-ai-io service
type ServiceStatus struct {
	ServiceName   string    `json:"service_name"`
	Status        string    `json:"status"`
	Version       string    `json:"version"`
	Uptime        string    `json:"uptime"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
}
