package input

import (
	"context"
	"fmt"
	"time"
)

// MockAIPersonaCommunityClient provides a mock implementation for testing
type MockAIPersonaCommunityClient struct{}

// NewMockAIPersonaCommunityClient creates a new mock AI client
func NewMockAIPersonaCommunityClient() *MockAIPersonaCommunityClient {
	return &MockAIPersonaCommunityClient{}
}

// CreateCommunity creates a mock AI community
func (m *MockAIPersonaCommunityClient) CreateCommunity(ctx context.Context, topic string, personaCount int) (*Community, error) {
	community := &Community{
		ID:        fmt.Sprintf("community_%d", time.Now().UnixNano()),
		Topic:     topic,
		Members:   make([]PersonaInfo, personaCount),
		CreatedAt: time.Now(),
		Status:    "active",
	}
	
	// Create mock personas
	for i := 0; i < personaCount; i++ {
		community.Members[i] = PersonaInfo{
			ID:          fmt.Sprintf("persona_%d_%d", i+1, time.Now().UnixNano()),
			Name:        fmt.Sprintf("AI Persona %d", i+1),
			Expertise:   []string{"threat-analysis", "security", "ai"},
			Description: fmt.Sprintf("Mock AI persona for %s analysis", topic),
			Model:       "mock-gpt-4",
		}
	}
	
	return community, nil
}

// SubmitForReview submits content for mock AI community review
func (m *MockAIPersonaCommunityClient) SubmitForReview(ctx context.Context, communityID string, content string) (*CommunityReview, error) {
	review := &CommunityReview{
		ReviewID:        fmt.Sprintf("review_%d", time.Now().UnixNano()),
		Topic:           "mock-analysis",
		Content:         content,
		PersonaReviews:  []PersonaReview{},
		Recommendations: []string{"Mock recommendation 1", "Mock recommendation 2"},
		Metadata:        make(map[string]interface{}),
		CreatedAt:       time.Now(),
	}
	
	// Create mock persona reviews
	personaCount := 3 // Default for mock
	totalScore := 0.0
	
	for i := 0; i < personaCount; i++ {
		score := 0.3 + (float64(i) * 0.2) // Scores: 0.3, 0.5, 0.7
		totalScore += score
		
		personaReview := PersonaReview{
			PersonaID:   fmt.Sprintf("persona_%d", i+1),
			PersonaName: fmt.Sprintf("Mock Persona %d", i+1),
			Expertise:   []string{"security", "analysis"},
			Review:      fmt.Sprintf("Mock review %d: This content appears to have moderate threat indicators.", i+1),
			Score:       score,
			Confidence:  0.8,
			Tags:        []string{"mock", "analysis"},
			Metadata:    make(map[string]interface{}),
			Timestamp:   time.Now(),
		}
		
		review.PersonaReviews = append(review.PersonaReviews, personaReview)
	}
	
	// Calculate consensus
	avgScore := totalScore / float64(personaCount)
	review.Consensus = &Consensus{
		OverallScore:    avgScore,
		Agreement:       0.75,
		Recommendation:  "moderate_threat",
		KeyPoints:       []string{"Mock key point 1", "Mock key point 2"},
		ConfidenceLevel: 0.8,
	}
	
	// Add sentiment analysis
	review.Sentiment = &SentimentAnalysis{
		Overall:      "neutral",
		Score:        0.0,
		Emotions:     map[string]float64{"neutral": 0.7, "concern": 0.3},
		Toxicity:     0.1,
		Subjectivity: 0.4,
	}
	
	completedAt := time.Now()
	review.CompletedAt = &completedAt
	
	return review, nil
}

// GetReviewStatus gets the status of a mock review
func (m *MockAIPersonaCommunityClient) GetReviewStatus(ctx context.Context, reviewID string) (*CommunityReview, error) {
	// Return a completed mock review
	return m.SubmitForReview(ctx, "mock-community", "mock content for status check")
}

// GetCommunityMembers gets mock community members
func (m *MockAIPersonaCommunityClient) GetCommunityMembers(ctx context.Context, communityID string) ([]PersonaInfo, error) {
	members := []PersonaInfo{
		{
			ID:          "persona_1",
			Name:        "Security Analyst AI",
			Expertise:   []string{"cybersecurity", "threat-analysis", "malware"},
			Description: "Specialized in cybersecurity threat analysis",
			Model:       "mock-gpt-4",
		},
		{
			ID:          "persona_2",
			Name:        "Social Engineering Expert AI",
			Expertise:   []string{"social-engineering", "phishing", "psychology"},
			Description: "Expert in social engineering and human psychology",
			Model:       "mock-gpt-4",
		},
		{
			ID:          "persona_3",
			Name:        "Network Security AI",
			Expertise:   []string{"network-security", "intrusion-detection", "forensics"},
			Description: "Specialized in network security and digital forensics",
			Model:       "mock-gpt-4",
		},
	}
	
	return members, nil
}
