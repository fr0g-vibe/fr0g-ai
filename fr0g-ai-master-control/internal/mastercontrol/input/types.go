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
