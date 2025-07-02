package input

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// DiscordWebhookProcessor handles Discord webhook messages for AI community review
type DiscordWebhookProcessor struct {
	aiCommunityClient AIPersonaCommunityClient
	config            *DiscordProcessorConfig
}

// DiscordProcessorConfig holds configuration for Discord webhook processing
type DiscordProcessorConfig struct {
	CommunityTopic     string   `yaml:"community_topic"`
	PersonaCount       int      `yaml:"persona_count"`
	ReviewTimeout      time.Duration `yaml:"review_timeout"`
	RequiredConsensus  float64  `yaml:"required_consensus"`
	EnableSentiment    bool     `yaml:"enable_sentiment"`
	FilterKeywords     []string `yaml:"filter_keywords"`
}

// DiscordMessage represents a Discord message from webhook
type DiscordMessage struct {
	ID          string    `json:"id"`
	Content     string    `json:"content"`
	Author      Author    `json:"author"`
	ChannelID   string    `json:"channel_id"`
	GuildID     string    `json:"guild_id"`
	Timestamp   time.Time `json:"timestamp"`
	Attachments []Attachment `json:"attachments"`
	Embeds      []Embed   `json:"embeds"`
}

// Author represents a Discord message author
type Author struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Bot      bool   `json:"bot"`
}

// Attachment represents a Discord message attachment
type Attachment struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Size     int    `json:"size"`
}

// Embed represents a Discord message embed
type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Color       int    `json:"color"`
}


// NewDiscordWebhookProcessor creates a new Discord webhook processor
func NewDiscordWebhookProcessor(client AIPersonaCommunityClient, config *DiscordProcessorConfig) *DiscordWebhookProcessor {
	return &DiscordWebhookProcessor{
		aiCommunityClient: client,
		config:            config,
	}
}

// ProcessWebhook processes Discord webhook requests
func (dp *DiscordWebhookProcessor) ProcessWebhook(ctx context.Context, request *WebhookRequest) (*WebhookResponse, error) {
	log.Printf("Discord Processor: Processing webhook for Discord message, request ID: %s", request.ID)
	
	// Parse Discord message from webhook body
	discordMsg, err := dp.parseDiscordMessage(request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Discord message: %w", err)
	}
	
	// Filter content if needed
	if dp.shouldFilterContent(discordMsg.Content) {
		log.Printf("Discord Processor: Content filtered for request %s", request.ID)
		return &WebhookResponse{
			Success:   true,
			Message:   "Content filtered - no review needed",
			RequestID: request.ID,
			Data: map[string]interface{}{
				"action": "filtered",
				"reason": "content_filter",
			},
			Timestamp: time.Now(),
		}, nil
	}
	
	// Determine appropriate topic based on content
	topic := dp.determineTopicFromContent(discordMsg.Content)
	
	// Create AI community for review
	community, err := dp.aiCommunityClient.CreateCommunity(ctx, topic, dp.config.PersonaCount)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI community: %w", err)
	}
	
	log.Printf("Discord Processor: Created AI community %s for topic '%s' (detected from content)", community.ID, topic)
	
	// Submit content for review
	review, err := dp.aiCommunityClient.SubmitForReview(ctx, community.ID, discordMsg.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to submit content for review: %w", err)
	}
	
	log.Printf("Discord Processor: Submitted content for review, review ID: %s", review.ReviewID)
	
	// For the real AIP client, the review is completed immediately
	// In a production system, this would be asynchronous
	finalReview := review
	
	// Process review results
	action := dp.determineAction(finalReview)
	
	log.Printf("Discord Processor: Review completed for request %s, action: %s", request.ID, action)
	log.Printf("Discord Processor: AI Community Review Summary:")
	log.Printf("  - Overall Score: %.2f", finalReview.Consensus.OverallScore)
	log.Printf("  - Recommendation: %s", finalReview.Consensus.Recommendation)
	log.Printf("  - Persona Reviews: %d", len(finalReview.PersonaReviews))
	
	for i, personaReview := range finalReview.PersonaReviews {
		log.Printf("  - Persona %d (%s): %.2f - %s", i+1, personaReview.PersonaName, personaReview.Score, personaReview.Review)
	}
	
	return &WebhookResponse{
		Success:   true,
		Message:   "Discord message reviewed by AI community",
		RequestID: request.ID,
		Data: map[string]interface{}{
			"discord_message": discordMsg,
			"community_id":    community.ID,
			"review":          finalReview,
			"action":          action,
			"persona_count":   len(finalReview.PersonaReviews),
		},
		Timestamp: time.Now(),
	}, nil
}

// GetTag returns the processor tag
func (dp *DiscordWebhookProcessor) GetTag() string {
	return "discord"
}

// GetDescription returns the processor description
func (dp *DiscordWebhookProcessor) GetDescription() string {
	clientType := "mock"
	if _, ok := dp.aiCommunityClient.(*RealAIPersonaCommunityClient); ok {
		clientType = "real"
	}
	return fmt.Sprintf("Discord message processor (%s AIP client) for AI community review on topic: %s", clientType, dp.config.CommunityTopic)
}

// parseDiscordMessage parses the Discord message from webhook body
func (dp *DiscordWebhookProcessor) parseDiscordMessage(body interface{}) (*DiscordMessage, error) {
	// This is a placeholder implementation
	// In a real implementation, you would parse the actual Discord webhook format
	
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid webhook body format")
	}
	
	content, _ := bodyMap["content"].(string)
	if content == "" {
		return nil, fmt.Errorf("no content found in Discord message")
	}
	
	// Create a simplified Discord message
	return &DiscordMessage{
		ID:        fmt.Sprintf("msg_%d", time.Now().UnixNano()),
		Content:   content,
		Author:    Author{Username: "unknown"},
		Timestamp: time.Now(),
	}, nil
}

// shouldFilterContent determines if content should be filtered
func (dp *DiscordWebhookProcessor) shouldFilterContent(content string) bool {
	// Only filter truly harmful content
	harmfulKeywords := []string{"spam", "abuse", "harmful", "toxic", "hate"}
	
	for _, keyword := range harmfulKeywords {
		if containsString(content, keyword) {
			return true
		}
	}
	
	// Allow most content through for AI community review
	// The AI personas will provide the real content analysis
	return false
}

// waitForReviewCompletion waits for the AI community review to complete
func (dp *DiscordWebhookProcessor) waitForReviewCompletion(ctx context.Context, reviewID string) (*CommunityReview, error) {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("review timeout exceeded")
		case <-ticker.C:
			review, err := dp.aiCommunityClient.GetReviewStatus(ctx, reviewID)
			if err != nil {
				return nil, err
			}
			
			if review.CompletedAt != nil {
				return review, nil
			}
			
			log.Printf("Discord Processor: Review %s still in progress...", reviewID)
		}
	}
}

// determineAction determines what action to take based on review results
func (dp *DiscordWebhookProcessor) determineAction(review *CommunityReview) string {
	if review.Consensus == nil {
		return "no_consensus"
	}
	
	// Simple action determination based on consensus score
	switch {
	case review.Consensus.OverallScore >= 0.8:
		return "approve"
	case review.Consensus.OverallScore >= 0.6:
		return "review_required"
	case review.Consensus.OverallScore >= 0.4:
		return "flag_for_attention"
	default:
		return "reject"
	}
}

// determineTopicFromContent analyzes content to determine the best AI community topic
func (dp *DiscordWebhookProcessor) determineTopicFromContent(content string) string {
	lowerContent := strings.ToLower(content)
	
	log.Printf("Discord Processor: Analyzing content for topic detection: '%s'", content[:min(100, len(content))])
	log.Printf("Discord Processor: Content analysis - hasCode: %v, hasAlgorithm: %v, hasConsciousness: %v", 
		strings.Contains(lowerContent, "```"),
		strings.Contains(lowerContent, "algorithm"),
		strings.Contains(lowerContent, "consciousness"))
	
	// Check for code review content first (most specific)
	if strings.Contains(lowerContent, "```") && strings.Contains(lowerContent, "review") {
		log.Printf("Discord Processor: Detected code_review topic (contains code blocks + review)")
		return "code_review"
	}
	
	// Check for security-related code content
	if strings.Contains(lowerContent, "```") && 
	   (strings.Contains(lowerContent, "security") || 
	    strings.Contains(lowerContent, "vulnerability") ||
	    strings.Contains(lowerContent, "eval(")) {
		log.Printf("Discord Processor: Detected code_review topic (security-related code)")
		return "code_review"
	}
	
	// Check for technical/algorithm content
	if strings.Contains(lowerContent, "algorithm") ||
	   strings.Contains(lowerContent, "performance") ||
	   strings.Contains(lowerContent, "optimization") ||
	   strings.Contains(lowerContent, "complexity") ||
	   strings.Contains(lowerContent, "o(n") ||
	   strings.Contains(lowerContent, "scalability") {
		log.Printf("Discord Processor: Detected technical_discussion topic")
		return "technical_discussion"
	}
	
	// Check for AI consciousness content
	if strings.Contains(lowerContent, "consciousness") ||
	   strings.Contains(lowerContent, "awareness") ||
	   strings.Contains(lowerContent, "emergent") ||
	   strings.Contains(lowerContent, "cognitive architecture") ||
	   strings.Contains(lowerContent, "self-awareness") ||
	   strings.Contains(lowerContent, "subjective experience") ||
	   strings.Contains(lowerContent, "ai personas") {
		log.Printf("Discord Processor: Detected ai_consciousness topic")
		return "ai_consciousness"
	}
	
	// Check for general code content
	if strings.Contains(lowerContent, "```") ||
	   strings.Contains(lowerContent, "function") ||
	   strings.Contains(lowerContent, "code") {
		log.Printf("Discord Processor: Detected technical_discussion topic (general code)")
		return "technical_discussion"
	}
	
	// Default to general discussion
	log.Printf("Discord Processor: Using default general_discussion topic")
	return "general_discussion"
}

// Utility function for string contains check
func containsString(text, substring string) bool {
	// Proper case-insensitive contains check
	if text == "" || substring == "" {
		return false
	}
	
	// Convert both to lowercase for case-insensitive comparison
	lowerText := strings.ToLower(text)
	lowerSubstring := strings.ToLower(substring)
	
	return strings.Contains(lowerText, lowerSubstring)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
