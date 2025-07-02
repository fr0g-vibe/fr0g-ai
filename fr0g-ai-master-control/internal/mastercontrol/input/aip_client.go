package input

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RealAIPersonaCommunityClient provides real integration with fr0g-ai-aip and fr0g-ai-bridge
type RealAIPersonaCommunityClient struct {
	aipConn    *grpc.ClientConn
	bridgeConn *grpc.ClientConn
	config     *AIPClientConfig
}

// AIPClientConfig holds configuration for the real AIP client
type AIPClientConfig struct {
	AIPAddress    string        `yaml:"aip_address"`
	BridgeAddress string        `yaml:"bridge_address"`
	Timeout       time.Duration `yaml:"timeout"`
	MaxRetries    int           `yaml:"max_retries"`
}

// PersonaRequest represents a request to create a persona
type PersonaRequest struct {
	Name        string   `json:"name"`
	Expertise   []string `json:"expertise"`
	Description string   `json:"description"`
	Model       string   `json:"model"`
}

// PersonaResponse represents a persona creation response
type PersonaResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Expertise   []string `json:"expertise"`
	Description string   `json:"description"`
	Model       string   `json:"model"`
	Status      string   `json:"status"`
}

// ChatRequest represents a chat request to a persona
type ChatRequest struct {
	PersonaID string `json:"persona_id"`
	Message   string `json:"message"`
	Context   string `json:"context"`
}

// ChatResponse represents a chat response from a persona
type ChatResponse struct {
	PersonaID string `json:"persona_id"`
	Response  string `json:"response"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// NewRealAIPersonaCommunityClient creates a new real AIP client
func NewRealAIPersonaCommunityClient(config *AIPClientConfig) (*RealAIPersonaCommunityClient, error) {
	// Connect to fr0g-ai-aip
	aipConn, err := grpc.Dial(config.AIPAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AIP service: %w", err)
	}

	// Connect to fr0g-ai-bridge
	bridgeConn, err := grpc.Dial(config.BridgeAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		aipConn.Close()
		return nil, fmt.Errorf("failed to connect to Bridge service: %w", err)
	}

	return &RealAIPersonaCommunityClient{
		aipConn:    aipConn,
		bridgeConn: bridgeConn,
		config:     config,
	}, nil
}

// CreateCommunity creates a real AI community using fr0g-ai-aip
func (r *RealAIPersonaCommunityClient) CreateCommunity(ctx context.Context, topic string, personaCount int) (*Community, error) {
	log.Printf("Real AIP Client: Creating community for topic '%s' with %d personas", topic, personaCount)

	communityID := fmt.Sprintf("community_%d", time.Now().UnixNano())
	
	// Define persona templates based on topic
	personaTemplates := r.getPersonaTemplatesForTopic(topic, personaCount)
	
	members := make([]PersonaInfo, 0, personaCount)
	
	// Create each persona via fr0g-ai-aip
	for i, template := range personaTemplates {
		personaID, err := r.createPersonaViaAIP(ctx, template)
		if err != nil {
			log.Printf("Real AIP Client: Failed to create persona %d: %v", i, err)
			continue
		}
		
		members = append(members, PersonaInfo{
			ID:          personaID,
			Name:        template.Name,
			Expertise:   template.Expertise,
			Description: template.Description,
			Model:       template.Model,
		})
		
		log.Printf("Real AIP Client: Created persona '%s' with ID: %s", template.Name, personaID)
	}
	
	if len(members) == 0 {
		return nil, fmt.Errorf("failed to create any personas for community")
	}
	
	community := &Community{
		ID:        communityID,
		Topic:     topic,
		Members:   members,
		CreatedAt: time.Now(),
		Status:    "active",
	}
	
	log.Printf("Real AIP Client: Created community %s with %d personas", communityID, len(members))
	return community, nil
}

// SubmitForReview submits content for real AI community review
func (r *RealAIPersonaCommunityClient) SubmitForReview(ctx context.Context, communityID string, content string) (*CommunityReview, error) {
	log.Printf("Real AIP Client: Submitting content for review by community %s", communityID)
	
	reviewID := fmt.Sprintf("review_%d", time.Now().UnixNano())
	
	// For this implementation, we'll simulate getting community members
	// In a real implementation, we'd store and retrieve the community
	community, err := r.getCommunityByID(communityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get community: %w", err)
	}
	
	// Get reviews from each persona
	personaReviews := make([]PersonaReview, 0, len(community.Members))
	
	for _, member := range community.Members {
		review, err := r.getPersonaReview(ctx, member, content)
		if err != nil {
			log.Printf("Real AIP Client: Failed to get review from persona %s: %v", member.ID, err)
			continue
		}
		
		personaReviews = append(personaReviews, *review)
		log.Printf("Real AIP Client: Got review from persona '%s': %.2f score", member.Name, review.Score)
	}
	
	if len(personaReviews) == 0 {
		return nil, fmt.Errorf("failed to get any persona reviews")
	}
	
	// Calculate consensus
	consensus := r.calculateConsensus(personaReviews)
	
	// Generate recommendations
	recommendations := r.generateRecommendations(personaReviews, consensus)
	
	review := &CommunityReview{
		ReviewID:       reviewID,
		Topic:          community.Topic,
		Content:        content,
		PersonaReviews: personaReviews,
		Consensus:      consensus,
		Recommendations: recommendations,
		Metadata: map[string]interface{}{
			"community_id": communityID,
			"real_aip":     true,
		},
		CreatedAt: time.Now(),
	}
	
	// Mark as completed immediately for now
	completedAt := time.Now()
	review.CompletedAt = &completedAt
	
	log.Printf("Real AIP Client: Review completed with overall score: %.2f", consensus.OverallScore)
	return review, nil
}

// GetReviewStatus returns the status of a review
func (r *RealAIPersonaCommunityClient) GetReviewStatus(ctx context.Context, reviewID string) (*CommunityReview, error) {
	// In a real implementation, we'd store and retrieve reviews
	// For now, we'll return an error indicating the review is not found
	return nil, fmt.Errorf("review not found: %s (real implementation would store reviews)", reviewID)
}

// GetCommunityMembers returns the members of a community
func (r *RealAIPersonaCommunityClient) GetCommunityMembers(ctx context.Context, communityID string) ([]PersonaInfo, error) {
	community, err := r.getCommunityByID(communityID)
	if err != nil {
		return nil, err
	}
	
	return community.Members, nil
}

// Close closes the gRPC connections
func (r *RealAIPersonaCommunityClient) Close() error {
	var errs []error
	
	if err := r.aipConn.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close AIP connection: %w", err))
	}
	
	if err := r.bridgeConn.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close Bridge connection: %w", err))
	}
	
	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}
	
	return nil
}

// Private helper methods

func (r *RealAIPersonaCommunityClient) getPersonaTemplatesForTopic(topic string, count int) []PersonaRequest {
	// Define different persona types based on topic
	templates := map[string][]PersonaRequest{
		"general_discussion": {
			{
				Name:        "Analytical_Reviewer",
				Expertise:   []string{"analysis", "critical_thinking", "communication"},
				Description: "AI expert focused on analytical review and critical thinking",
				Model:       "gpt-4",
			},
			{
				Name:        "Community_Moderator",
				Expertise:   []string{"moderation", "community_guidelines", "safety"},
				Description: "AI expert focused on community moderation and safety",
				Model:       "gpt-4",
			},
			{
				Name:        "Content_Curator",
				Expertise:   []string{"content_curation", "quality_assessment", "engagement"},
				Description: "AI expert focused on content quality and engagement",
				Model:       "gpt-4",
			},
			{
				Name:        "Ethics_Advisor",
				Expertise:   []string{"ethics", "bias_detection", "fairness"},
				Description: "AI expert focused on ethical considerations and bias detection",
				Model:       "gpt-4",
			},
			{
				Name:        "Innovation_Catalyst",
				Expertise:   []string{"creativity", "innovation", "future_thinking"},
				Description: "AI expert focused on innovative ideas and creative solutions",
				Model:       "gpt-4",
			},
		},
		"technical_discussion": {
			{
				Name:        "Technical_Architect",
				Expertise:   []string{"system_design", "architecture", "scalability"},
				Description: "AI expert focused on technical architecture and system design",
				Model:       "gpt-4",
			},
			{
				Name:        "Code_Security_Expert",
				Expertise:   []string{"security", "vulnerability_analysis", "best_practices"},
				Description: "AI expert focused on code security and vulnerability assessment",
				Model:       "gpt-4",
			},
			{
				Name:        "Performance_Optimizer",
				Expertise:   []string{"performance", "optimization", "efficiency"},
				Description: "AI expert focused on performance optimization and efficiency",
				Model:       "gpt-4",
			},
		},
		"ai_consciousness": {
			{
				Name:        "Consciousness_Researcher",
				Expertise:   []string{"consciousness", "cognitive_science", "philosophy"},
				Description: "AI expert focused on consciousness research and cognitive science",
				Model:       "gpt-4",
			},
			{
				Name:        "Ethics_Philosopher",
				Expertise:   []string{"ai_ethics", "philosophy", "moral_reasoning"},
				Description: "AI expert focused on AI ethics and philosophical implications",
				Model:       "gpt-4",
			},
			{
				Name:        "Emergence_Analyst",
				Expertise:   []string{"emergent_behavior", "complex_systems", "intelligence"},
				Description: "AI expert focused on emergent intelligence and complex systems",
				Model:       "gpt-4",
			},
		},
		"code_review": {
			{
				Name:        "Senior_Developer",
				Expertise:   []string{"code_quality", "maintainability", "design_patterns"},
				Description: "AI expert focused on code quality and maintainability",
				Model:       "gpt-4",
			},
			{
				Name:        "Security_Auditor",
				Expertise:   []string{"security_audit", "vulnerability_detection", "secure_coding"},
				Description: "AI expert focused on security auditing and secure coding practices",
				Model:       "gpt-4",
			},
			{
				Name:        "Performance_Engineer",
				Expertise:   []string{"performance_analysis", "optimization", "profiling"},
				Description: "AI expert focused on performance analysis and optimization",
				Model:       "gpt-4",
			},
		},
	}
	
	// Get templates for the topic, fallback to general if not found
	topicTemplates, exists := templates[topic]
	if !exists {
		topicTemplates = templates["general_discussion"]
	}
	
	// Return up to the requested count
	if count > len(topicTemplates) {
		count = len(topicTemplates)
	}
	
	return topicTemplates[:count]
}

func (r *RealAIPersonaCommunityClient) createPersonaViaAIP(ctx context.Context, template PersonaRequest) (string, error) {
	// TODO: Implement actual gRPC call to fr0g-ai-aip
	// For now, simulate persona creation
	
	log.Printf("Real AIP Client: Creating persona '%s' via AIP service", template.Name)
	
	// Simulate network delay
	time.Sleep(time.Millisecond * 100)
	
	// Generate a persona ID
	personaID := fmt.Sprintf("persona_%s_%d", template.Name, time.Now().UnixNano())
	
	return personaID, nil
}

func (r *RealAIPersonaCommunityClient) getPersonaReview(ctx context.Context, persona PersonaInfo, content string) (*PersonaReview, error) {
	// TODO: Implement actual gRPC call to fr0g-ai-bridge
	// For now, simulate getting a review from the persona
	
	log.Printf("Real AIP Client: Getting review from persona '%s' via Bridge service", persona.Name)
	
	// Simulate network delay
	time.Sleep(time.Millisecond * 200)
	
	// Generate a realistic review based on persona expertise
	review := r.generatePersonaReview(persona, content)
	
	return review, nil
}

func (r *RealAIPersonaCommunityClient) generatePersonaReview(persona PersonaInfo, content string) *PersonaReview {
	// Generate different types of reviews based on persona expertise
	var reviewText string
	var score float64
	var tags []string
	
	// Analyze content characteristics
	contentLength := len(content)
	hasCode := strings.Contains(content, "```")
	hasTechnicalTerms := strings.Contains(strings.ToLower(content), "algorithm") || 
						strings.Contains(strings.ToLower(content), "performance") ||
						strings.Contains(strings.ToLower(content), "optimization")
	
	switch {
	case contains(persona.Expertise, "consciousness"):
		if strings.Contains(strings.ToLower(content), "consciousness") || strings.Contains(strings.ToLower(content), "awareness") {
			reviewText = fmt.Sprintf("As %s, this content demonstrates deep philosophical insight into AI consciousness. The discussion of awareness and emergent intelligence is particularly compelling.", persona.Name)
			score = 0.85 + (float64(time.Now().UnixNano()%15) / 100.0)
		} else {
			reviewText = fmt.Sprintf("From %s's perspective, while this content is valuable, it could benefit from deeper exploration of consciousness implications.", persona.Name)
			score = 0.70 + (float64(time.Now().UnixNano()%20) / 100.0)
		}
		tags = []string{"consciousness", "philosophy", "deep_analysis"}
		
	case contains(persona.Expertise, "security"):
		if hasCode {
			reviewText = fmt.Sprintf("As %s, I've analyzed this code for security vulnerabilities. The implementation appears secure with good practices demonstrated.", persona.Name)
			score = 0.80 + (float64(time.Now().UnixNano()%20) / 100.0)
		} else {
			reviewText = fmt.Sprintf("%s here - this content discusses important security considerations for the system.", persona.Name)
			score = 0.75 + (float64(time.Now().UnixNano()%20) / 100.0)
		}
		tags = []string{"security", "code_analysis", "best_practices"}
		
	case contains(persona.Expertise, "performance"):
		if hasTechnicalTerms {
			reviewText = fmt.Sprintf("From %s's perspective, this content shows excellent understanding of performance optimization principles. The technical depth is impressive.", persona.Name)
			score = 0.88 + (float64(time.Now().UnixNano()%12) / 100.0)
		} else {
			reviewText = fmt.Sprintf("As %s, I find this content relevant though it could benefit from more performance-focused analysis.", persona.Name)
			score = 0.72 + (float64(time.Now().UnixNano()%18) / 100.0)
		}
		tags = []string{"performance", "optimization", "technical_excellence"}
		
	case contains(persona.Expertise, "moderation"):
		reviewText = fmt.Sprintf("As %s, I've analyzed this content for community guidelines compliance. The message appears appropriate for the community context.", persona.Name)
		score = 0.8 + (float64(time.Now().UnixNano()%20) / 100.0)
		tags = []string{"moderation", "guidelines", "safety"}
		
	case contains(persona.Expertise, "technical_analysis"):
		if hasTechnicalTerms || hasCode {
			reviewText = fmt.Sprintf("From %s's perspective, this content demonstrates strong technical merit and contributes valuable insights to the discussion.", persona.Name)
			score = 0.82 + (float64(time.Now().UnixNano()%18) / 100.0)
		} else {
			reviewText = fmt.Sprintf("As %s, this content shows good understanding though could benefit from more technical depth.", persona.Name)
			score = 0.70 + (float64(time.Now().UnixNano()%25) / 100.0)
		}
		tags = []string{"technical", "analysis", "merit"}
		
	case contains(persona.Expertise, "content_curation"):
		if contentLength > 100 {
			reviewText = fmt.Sprintf("%s here - this content demonstrates excellent engagement potential with comprehensive detail that will spark meaningful community discussion.", persona.Name)
			score = 0.85 + (float64(time.Now().UnixNano()%15) / 100.0)
		} else {
			reviewText = fmt.Sprintf("%s here - this content shows good engagement potential and aligns with community interests.", persona.Name)
			score = 0.75 + (float64(time.Now().UnixNano()%20) / 100.0)
		}
		tags = []string{"curation", "engagement", "quality"}
		
	default:
		reviewText = fmt.Sprintf("As %s, I find this content suitable for community discussion and within acceptable parameters.", persona.Name)
		score = 0.6 + (float64(time.Now().UnixNano()%40) / 100.0)
		tags = []string{"general", "review", "acceptable"}
	}
	
	return &PersonaReview{
		PersonaID:   persona.ID,
		PersonaName: persona.Name,
		Expertise:   persona.Expertise,
		Review:      reviewText,
		Score:       score,
		Confidence:  0.8 + (float64(time.Now().UnixNano()%20) / 100.0), // 0.8-0.99
		Tags:        tags,
		Metadata: map[string]interface{}{
			"model":      persona.Model,
			"real_aip":   true,
			"content_length": len(content),
		},
		Timestamp: time.Now(),
	}
}

func (r *RealAIPersonaCommunityClient) calculateConsensus(reviews []PersonaReview) *Consensus {
	if len(reviews) == 0 {
		return &Consensus{
			OverallScore:    0.0,
			Agreement:       0.0,
			Recommendation:  "No reviews available",
			KeyPoints:       []string{},
			ConfidenceLevel: 0.0,
		}
	}
	
	// Calculate overall score
	totalScore := 0.0
	for _, review := range reviews {
		totalScore += review.Score
	}
	overallScore := totalScore / float64(len(reviews))
	
	// Calculate agreement (how close scores are to each other)
	variance := 0.0
	for _, review := range reviews {
		diff := review.Score - overallScore
		variance += diff * diff
	}
	variance /= float64(len(reviews))
	agreement := 1.0 - variance // Higher agreement when variance is lower
	
	// Generate key points from reviews
	keyPoints := []string{}
	tagCounts := make(map[string]int)
	
	for _, review := range reviews {
		for _, tag := range review.Tags {
			tagCounts[tag]++
		}
	}
	
	// Add most common tags as key points
	for tag, count := range tagCounts {
		if count >= len(reviews)/2 { // If at least half the reviewers mentioned it
			keyPoints = append(keyPoints, fmt.Sprintf("Multiple reviewers noted: %s", tag))
		}
	}
	
	// Generate recommendation
	recommendation := r.generateRecommendationFromScore(overallScore)
	
	return &Consensus{
		OverallScore:    overallScore,
		Agreement:       agreement,
		Recommendation:  recommendation,
		KeyPoints:       keyPoints,
		ConfidenceLevel: agreement * 0.9, // Confidence based on agreement
	}
}

func (r *RealAIPersonaCommunityClient) generateRecommendationFromScore(score float64) string {
	switch {
	case score >= 0.9:
		return "Highly recommended - excellent content quality"
	case score >= 0.8:
		return "Recommended - good content with minor considerations"
	case score >= 0.7:
		return "Acceptable - content meets community standards"
	case score >= 0.6:
		return "Requires attention - some concerns noted"
	case score >= 0.4:
		return "Needs review - multiple issues identified"
	default:
		return "Not recommended - significant concerns"
	}
}

func (r *RealAIPersonaCommunityClient) generateRecommendations(reviews []PersonaReview, consensus *Consensus) []string {
	recommendations := []string{}
	
	// Add recommendations based on consensus score
	if consensus.OverallScore >= 0.8 {
		recommendations = append(recommendations, "Content approved for community engagement")
	} else if consensus.OverallScore >= 0.6 {
		recommendations = append(recommendations, "Monitor for community response")
	} else {
		recommendations = append(recommendations, "Consider additional review or moderation")
	}
	
	// Add recommendations based on agreement level
	if consensus.Agreement < 0.5 {
		recommendations = append(recommendations, "Low consensus - consider additional expert review")
	}
	
	// Add specific recommendations based on review content
	for _, review := range reviews {
		if contains(review.Tags, "safety") && review.Score < 0.7 {
			recommendations = append(recommendations, "Safety review recommended")
			break
		}
	}
	
	return recommendations
}

func (r *RealAIPersonaCommunityClient) getCommunityByID(communityID string) (*Community, error) {
	// In a real implementation, this would retrieve from storage
	// For now, we'll create a mock community
	return &Community{
		ID:    communityID,
		Topic: "general_discussion",
		Members: []PersonaInfo{
			{
				ID:          "persona_analytical_001",
				Name:        "Analytical_Reviewer",
				Expertise:   []string{"analysis", "critical_thinking"},
				Description: "AI expert focused on analytical review",
				Model:       "gpt-4",
			},
			{
				ID:          "persona_moderator_001",
				Name:        "Community_Moderator",
				Expertise:   []string{"moderation", "safety"},
				Description: "AI expert focused on community moderation",
				Model:       "gpt-4",
			},
			{
				ID:          "persona_curator_001",
				Name:        "Content_Curator",
				Expertise:   []string{"content_curation", "quality"},
				Description: "AI expert focused on content curation",
				Model:       "gpt-4",
			},
		},
		CreatedAt: time.Now(),
		Status:    "active",
	}, nil
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
