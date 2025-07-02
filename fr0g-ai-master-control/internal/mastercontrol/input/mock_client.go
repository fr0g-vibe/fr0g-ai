package input

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// MockAIPersonaCommunityClient provides a mock implementation for testing
type MockAIPersonaCommunityClient struct {
	communities map[string]*Community
	reviews     map[string]*CommunityReview
}

// NewMockAIPersonaCommunityClient creates a new mock client
func NewMockAIPersonaCommunityClient() *MockAIPersonaCommunityClient {
	return &MockAIPersonaCommunityClient{
		communities: make(map[string]*Community),
		reviews:     make(map[string]*CommunityReview),
	}
}

// CreateCommunity creates a mock AI community
func (m *MockAIPersonaCommunityClient) CreateCommunity(ctx context.Context, topic string, personaCount int) (*Community, error) {
	communityID := fmt.Sprintf("community_%d", time.Now().UnixNano())
	
	// Generate mock personas
	members := make([]PersonaInfo, personaCount)
	for i := 0; i < personaCount; i++ {
		members[i] = PersonaInfo{
			ID:          fmt.Sprintf("persona_%d_%d", i, time.Now().UnixNano()),
			Name:        fmt.Sprintf("Expert_%d", i+1),
			Expertise:   []string{topic, "analysis", "communication"},
			Description: fmt.Sprintf("AI expert specializing in %s", topic),
			Model:       "gpt-4",
		}
	}
	
	community := &Community{
		ID:        communityID,
		Topic:     topic,
		Members:   members,
		CreatedAt: time.Now(),
		Status:    "active",
	}
	
	m.communities[communityID] = community
	return community, nil
}

// SubmitForReview submits content for mock AI community review
func (m *MockAIPersonaCommunityClient) SubmitForReview(ctx context.Context, communityID string, content string) (*CommunityReview, error) {
	community, exists := m.communities[communityID]
	if !exists {
		return nil, fmt.Errorf("community not found: %s", communityID)
	}
	
	reviewID := fmt.Sprintf("review_%d", time.Now().UnixNano())
	
	// Create mock persona reviews
	personaReviews := make([]PersonaReview, len(community.Members))
	for i, member := range community.Members {
		personaReviews[i] = PersonaReview{
			PersonaID:   member.ID,
			PersonaName: member.Name,
			Expertise:   member.Expertise,
			Review:      m.generateMockReview(content, member.Name),
			Score:       0.5 + rand.Float64()*0.5, // Random score between 0.5-1.0
			Confidence:  0.7 + rand.Float64()*0.3, // Random confidence 0.7-1.0
			Tags:        []string{"analyzed", "reviewed"},
			Metadata:    map[string]interface{}{"model": member.Model},
			Timestamp:   time.Now(),
		}
	}
	
	// Calculate mock consensus
	totalScore := 0.0
	for _, review := range personaReviews {
		totalScore += review.Score
	}
	avgScore := totalScore / float64(len(personaReviews))
	
	consensus := &Consensus{
		OverallScore:    avgScore,
		Agreement:       0.8, // Mock agreement level
		Recommendation:  m.generateRecommendation(avgScore),
		KeyPoints:       []string{"Content analyzed", "Multiple perspectives considered"},
		ConfidenceLevel: 0.85,
	}
	
	review := &CommunityReview{
		ReviewID:       reviewID,
		Topic:          community.Topic,
		Content:        content,
		PersonaReviews: personaReviews,
		Consensus:      consensus,
		Recommendations: []string{"Consider context", "Monitor for updates"},
		Metadata: map[string]interface{}{
			"community_id": communityID,
			"mock":         true,
		},
		CreatedAt: time.Now(),
	}
	
	m.reviews[reviewID] = review
	
	// Simulate processing time
	go func() {
		time.Sleep(time.Second * 2) // Simulate 2-second processing
		completedAt := time.Now()
		review.CompletedAt = &completedAt
	}()
	
	return review, nil
}

// GetReviewStatus returns the status of a mock review
func (m *MockAIPersonaCommunityClient) GetReviewStatus(ctx context.Context, reviewID string) (*CommunityReview, error) {
	review, exists := m.reviews[reviewID]
	if !exists {
		return nil, fmt.Errorf("review not found: %s", reviewID)
	}
	
	return review, nil
}

// GetCommunityMembers returns the members of a mock community
func (m *MockAIPersonaCommunityClient) GetCommunityMembers(ctx context.Context, communityID string) ([]PersonaInfo, error) {
	community, exists := m.communities[communityID]
	if !exists {
		return nil, fmt.Errorf("community not found: %s", communityID)
	}
	
	return community.Members, nil
}

// generateMockReview generates a mock review based on content and persona
func (m *MockAIPersonaCommunityClient) generateMockReview(content, personaName string) string {
	reviews := []string{
		fmt.Sprintf("As %s, I find this content interesting and worth further discussion.", personaName),
		fmt.Sprintf("%s here - this content shows good potential but needs more context.", personaName),
		fmt.Sprintf("From %s's perspective, this content aligns well with current trends.", personaName),
		fmt.Sprintf("%s analysis: The content demonstrates clear thinking and good structure.", personaName),
		fmt.Sprintf("As %s, I recommend careful consideration of the implications discussed.", personaName),
	}
	
	return reviews[rand.Intn(len(reviews))]
}

// generateRecommendation generates a recommendation based on score
func (m *MockAIPersonaCommunityClient) generateRecommendation(score float64) string {
	switch {
	case score >= 0.8:
		return "Highly recommended for approval"
	case score >= 0.6:
		return "Recommended with minor considerations"
	case score >= 0.4:
		return "Requires additional review"
	default:
		return "Not recommended in current form"
	}
}
