package input

import (
	"context"
	"fmt"
	"log"
	"time"
)

// DiscordWebhookProcessor processes Discord webhook messages
type DiscordWebhookProcessor struct {
	aiClient AIPersonaCommunityClient
	config   *DiscordProcessorConfig
}

// DiscordProcessorConfig holds Discord processor configuration
type DiscordProcessorConfig struct {
	CommunityTopic    string        `yaml:"community_topic"`
	PersonaCount      int           `yaml:"persona_count"`
	ReviewTimeout     time.Duration `yaml:"review_timeout"`
	RequiredConsensus float64       `yaml:"required_consensus"`
	EnableSentiment   bool          `yaml:"enable_sentiment"`
	FilterKeywords    []string      `yaml:"filter_keywords"`
}

// DiscordMessage represents a Discord message
type DiscordMessage struct {
	ID          string                 `json:"id"`
	Content     string                 `json:"content"`
	Author      DiscordUser            `json:"author"`
	ChannelID   string                 `json:"channel_id"`
	GuildID     string                 `json:"guild_id"`
	Timestamp   time.Time              `json:"timestamp"`
	Attachments []DiscordAttachment    `json:"attachments"`
	Embeds      []DiscordEmbed         `json:"embeds"`
}

// DiscordUser represents a Discord user
type DiscordUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Bot      bool   `json:"bot"`
}

// DiscordAttachment represents a Discord attachment
type DiscordAttachment struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Size     int    `json:"size"`
}

// DiscordEmbed represents a Discord embed
type DiscordEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Color       int    `json:"color"`
}

// NewDiscordWebhookProcessor creates a new Discord webhook processor
func NewDiscordWebhookProcessor(aiClient AIPersonaCommunityClient, config *DiscordProcessorConfig) *DiscordWebhookProcessor {
	return &DiscordWebhookProcessor{
		aiClient: aiClient,
		config:   config,
	}
}

// GetTag returns the processor tag
func (d *DiscordWebhookProcessor) GetTag() string {
	return "discord"
}

// GetDescription returns the processor description
func (d *DiscordWebhookProcessor) GetDescription() string {
	return fmt.Sprintf("Discord message processor (mock AIP client) for AI community review on topic: %s", d.config.CommunityTopic)
}

// ProcessWebhook processes a Discord webhook
func (d *DiscordWebhookProcessor) ProcessWebhook(ctx context.Context, request *WebhookRequest) (*WebhookResponse, error) {
	log.Printf("Discord Processor: Processing Discord message webhook")
	
	// Parse Discord message from request body
	message, err := d.parseDiscordMessage(request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Discord message: %w", err)
	}
	
	// Create AI community for review
	community, err := d.aiClient.CreateCommunity(ctx, d.config.CommunityTopic, d.config.PersonaCount)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI community: %w", err)
	}
	
	// Submit message for AI community review
	review, err := d.aiClient.SubmitForReview(ctx, community.ID, message.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to submit for review: %w", err)
	}
	
	// Determine action based on consensus
	action := "review_required"
	if review.Consensus != nil && review.Consensus.OverallScore >= d.config.RequiredConsensus {
		action = "approve"
	}
	
	// Create response
	response := &WebhookResponse{
		Success:   true,
		Message:   "Discord message reviewed by AI community",
		RequestID: request.ID,
		Data: map[string]interface{}{
			"action":        action,
			"community_id":  community.ID,
			"discord_message": map[string]interface{}{
				"id":        fmt.Sprintf("msg_%d", time.Now().UnixNano()),
				"content":   message.Content,
				"author":    message.Author,
				"channel_id": message.ChannelID,
				"guild_id":  message.GuildID,
				"timestamp": time.Now(),
				"attachments": message.Attachments,
				"embeds":    message.Embeds,
			},
			"persona_count": d.config.PersonaCount,
			"review":       review,
		},
		Timestamp: time.Now(),
	}
	
	log.Printf("Discord Processor: Message processed successfully - Action: %s, Consensus: %.2f", 
		action, review.Consensus.OverallScore)
	
	return response, nil
}

// parseDiscordMessage parses a Discord message from the request body
func (d *DiscordWebhookProcessor) parseDiscordMessage(body interface{}) (*DiscordMessage, error) {
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid body format")
	}
	
	message := &DiscordMessage{
		Content:     getStringFromMap(bodyMap, "content"),
		ChannelID:   getStringFromMap(bodyMap, "channel_id"),
		GuildID:     getStringFromMap(bodyMap, "guild_id"),
		Timestamp:   time.Now(),
		Attachments: []DiscordAttachment{},
		Embeds:      []DiscordEmbed{},
	}
	
	// Parse author
	if authorData, ok := bodyMap["author"].(map[string]interface{}); ok {
		message.Author = DiscordUser{
			ID:       getStringFromMap(authorData, "id"),
			Username: getStringFromMap(authorData, "username"),
			Avatar:   getStringFromMap(authorData, "avatar"),
			Bot:      getBoolFromMap(authorData, "bot"),
		}
	}
	
	return message, nil
}

