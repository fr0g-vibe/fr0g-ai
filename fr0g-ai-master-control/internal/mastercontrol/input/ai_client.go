package input

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// RealAIPersonaCommunityClient provides a real implementation that connects to AI services
type RealAIPersonaCommunityClient struct {
	aipServiceURL    string
	bridgeServiceURL string
	httpClient       *http.Client
}

// AIClientConfig holds configuration for the AI client
type AIClientConfig struct {
	AIPServiceURL    string `yaml:"aip_service_url"`
	BridgeServiceURL string `yaml:"bridge_service_url"`
	Timeout          time.Duration `yaml:"timeout"`
}

// NewRealAIPersonaCommunityClient creates a new real AI client
func NewRealAIPersonaCommunityClient(config *AIClientConfig) *RealAIPersonaCommunityClient {
	return &RealAIPersonaCommunityClient{
		aipServiceURL:    config.AIPServiceURL,
		bridgeServiceURL: config.BridgeServiceURL,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// CreateCommunity creates an AI community by selecting personas from AIP service
func (r *RealAIPersonaCommunityClient) CreateCommunity(ctx context.Context, topic string, personaCount int) (*Community, error) {
	// Call AIP service to get available personas
	personas, err := r.getPersonasFromAIP(ctx, topic, personaCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get personas from AIP: %w", err)
	}
	
	community := &Community{
		ID:        fmt.Sprintf("community_%d", time.Now().UnixNano()),
		Topic:     topic,
		Members:   personas,
		CreatedAt: time.Now(),
		Status:    "active",
	}
	
	return community, nil
}

// SubmitForReview submits content for AI community review using the bridge service
func (r *RealAIPersonaCommunityClient) SubmitForReview(ctx context.Context, communityID string, content string) (*CommunityReview, error) {
	// Get community members first
	community, err := r.getCommunityByID(communityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get community: %w", err)
	}
	
	review := &CommunityReview{
		ReviewID:        fmt.Sprintf("review_%d", time.Now().UnixNano()),
		Topic:           community.Topic,
		Content:         content,
		PersonaReviews:  []PersonaReview{},
		Recommendations: []string{},
		Metadata:        make(map[string]interface{}),
		CreatedAt:       time.Now(),
	}
	
	// Submit to each persona via bridge service
	totalScore := 0.0
	for _, persona := range community.Members {
		personaReview, err := r.getPersonaReview(ctx, persona, content)
		if err != nil {
			// Log error but continue with other personas
			fmt.Printf("Warning: Failed to get review from persona %s: %v\n", persona.Name, err)
			continue
		}
		
		review.PersonaReviews = append(review.PersonaReviews, *personaReview)
		totalScore += personaReview.Score
	}
	
	// Calculate consensus
	if len(review.PersonaReviews) > 0 {
		avgScore := totalScore / float64(len(review.PersonaReviews))
		review.Consensus = &Consensus{
			OverallScore:    avgScore,
			Agreement:       r.calculateAgreement(review.PersonaReviews),
			Recommendation:  r.generateRecommendation(avgScore),
			KeyPoints:       r.extractKeyPoints(review.PersonaReviews),
			ConfidenceLevel: r.calculateConfidence(review.PersonaReviews),
		}
	}
	
	completedAt := time.Now()
	review.CompletedAt = &completedAt
	
	return review, nil
}

// GetReviewStatus gets the status of a review (placeholder for now)
func (r *RealAIPersonaCommunityClient) GetReviewStatus(ctx context.Context, reviewID string) (*CommunityReview, error) {
	// In a real implementation, this would query a database or cache
	return nil, fmt.Errorf("review status lookup not implemented yet")
}

// GetCommunityMembers gets community members
func (r *RealAIPersonaCommunityClient) GetCommunityMembers(ctx context.Context, communityID string) ([]PersonaInfo, error) {
	community, err := r.getCommunityByID(communityID)
	if err != nil {
		return nil, err
	}
	
	return community.Members, nil
}

// getPersonasFromAIP calls the AIP service to get personas
func (r *RealAIPersonaCommunityClient) getPersonasFromAIP(ctx context.Context, topic string, count int) ([]PersonaInfo, error) {
	// Create request to AIP service
	reqBody := map[string]interface{}{
		"topic": topic,
		"count": count,
		"filter": map[string]interface{}{
			"is_active": true,
		},
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", r.aipServiceURL+"/api/v1/personas/search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call AIP service: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AIP service returned status %d", resp.StatusCode)
	}
	
	var response struct {
		Personas []PersonaInfo `json:"personas"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return response.Personas, nil
}

// getPersonaReview gets a review from a specific persona via bridge service
func (r *RealAIPersonaCommunityClient) getPersonaReview(ctx context.Context, persona PersonaInfo, content string) (*PersonaReview, error) {
	// Create chat completion request
	reqBody := map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": fmt.Sprintf("You are %s, an AI expert in %s. Analyze the following content for threats and provide a detailed review.", persona.Name, persona.Expertise),
			},
			{
				"role":    "user",
				"content": content,
			},
		},
		"persona_prompt": fmt.Sprintf("Acting as %s with expertise in %s", persona.Name, persona.Expertise),
		"temperature":    0.7,
		"max_tokens":     1000,
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", r.bridgeServiceURL+"/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call bridge service: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bridge service returned status %d", resp.StatusCode)
	}
	
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI model")
	}
	
	reviewContent := response.Choices[0].Message.Content
	
	// Parse the review content to extract score and other metrics
	// This is a simplified implementation - in reality you'd want more sophisticated parsing
	score := r.extractScoreFromReview(reviewContent)
	confidence := r.extractConfidenceFromReview(reviewContent)
	
	return &PersonaReview{
		PersonaID:   persona.ID,
		PersonaName: persona.Name,
		Expertise:   persona.Expertise,
		Review:      reviewContent,
		Score:       score,
		Confidence:  confidence,
		Tags:        []string{"ai-generated", "threat-analysis"},
		Metadata:    map[string]interface{}{"model": "gpt-4"},
		Timestamp:   time.Now(),
	}, nil
}

// Helper methods for processing reviews
func (r *RealAIPersonaCommunityClient) extractScoreFromReview(review string) float64 {
	// Simplified scoring based on keywords - in reality this would be more sophisticated
	// You could use sentiment analysis, keyword matching, or train a specific model
	
	// Count threat indicators
	threatKeywords := []string{"threat", "malicious", "suspicious", "dangerous", "risk", "attack", "phishing", "malware"}
	safeKeywords := []string{"safe", "legitimate", "normal", "benign", "trusted", "secure"}
	
	threatCount := 0
	safeCount := 0
	
	for _, keyword := range threatKeywords {
		if contains(review, keyword) {
			threatCount++
		}
	}
	
	for _, keyword := range safeKeywords {
		if contains(review, keyword) {
			safeCount++
		}
	}
	
	// Calculate score (0.0 = safe, 1.0 = high threat)
	if threatCount == 0 && safeCount == 0 {
		return 0.5 // neutral
	}
	
	total := threatCount + safeCount
	return float64(threatCount) / float64(total)
}

func (r *RealAIPersonaCommunityClient) extractConfidenceFromReview(review string) float64 {
	// Simplified confidence calculation
	// In reality, this could be based on model confidence scores or other metrics
	return 0.8 // Default confidence
}

func (r *RealAIPersonaCommunityClient) calculateAgreement(reviews []PersonaReview) float64 {
	if len(reviews) < 2 {
		return 1.0
	}
	
	// Calculate variance in scores
	var sum, sumSquares float64
	for _, review := range reviews {
		sum += review.Score
		sumSquares += review.Score * review.Score
	}
	
	mean := sum / float64(len(reviews))
	variance := (sumSquares / float64(len(reviews))) - (mean * mean)
	
	// Convert variance to agreement (lower variance = higher agreement)
	agreement := 1.0 - variance
	if agreement < 0 {
		agreement = 0
	}
	
	return agreement
}

func (r *RealAIPersonaCommunityClient) generateRecommendation(score float64) string {
	switch {
	case score >= 0.8:
		return "High threat detected - immediate action recommended"
	case score >= 0.6:
		return "Moderate threat - further investigation recommended"
	case score >= 0.4:
		return "Low threat - monitoring recommended"
	default:
		return "Minimal threat - content appears safe"
	}
}

func (r *RealAIPersonaCommunityClient) extractKeyPoints(reviews []PersonaReview) []string {
	// Simplified key point extraction
	// In reality, this could use NLP techniques to extract common themes
	keyPoints := []string{}
	
	for _, review := range reviews {
		if len(review.Review) > 100 {
			// Extract first sentence as a key point
			sentences := splitSentences(review.Review)
			if len(sentences) > 0 {
				keyPoints = append(keyPoints, sentences[0])
			}
		}
	}
	
	return keyPoints
}

func (r *RealAIPersonaCommunityClient) calculateConfidence(reviews []PersonaReview) float64 {
	if len(reviews) == 0 {
		return 0.0
	}
	
	var totalConfidence float64
	for _, review := range reviews {
		totalConfidence += review.Confidence
	}
	
	return totalConfidence / float64(len(reviews))
}

func (r *RealAIPersonaCommunityClient) getCommunityByID(communityID string) (*Community, error) {
	// In a real implementation, this would query a database or cache
	// For now, return an error since we don't have persistence yet
	return nil, fmt.Errorf("community lookup not implemented yet - need to add persistence")
}

// Utility functions
func contains(text, substring string) bool {
	return len(text) >= len(substring) && 
		   (text == substring || 
		    (len(text) > len(substring) && 
		     (text[:len(substring)] == substring || 
		      text[len(text)-len(substring):] == substring ||
		      containsSubstring(text, substring))))
}

func containsSubstring(text, substring string) bool {
	for i := 0; i <= len(text)-len(substring); i++ {
		if text[i:i+len(substring)] == substring {
			return true
		}
	}
	return false
}

func splitSentences(text string) []string {
	// Simplified sentence splitting
	sentences := []string{}
	current := ""
	
	for _, char := range text {
		current += string(char)
		if char == '.' || char == '!' || char == '?' {
			sentences = append(sentences, current)
			current = ""
		}
	}
	
	if current != "" {
		sentences = append(sentences, current)
	}
	
	return sentences
}
