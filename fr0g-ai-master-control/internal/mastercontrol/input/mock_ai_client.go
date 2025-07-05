package input

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// MockAIPersonaCommunityClient provides a mock implementation for testing
type MockAIPersonaCommunityClient struct{}

// CreateCommunity creates a mock AI community
func (m *MockAIPersonaCommunityClient) CreateCommunity(ctx context.Context, topic string, personaCount int) (*Community, error) {
	// Create realistic personas based on topic
	personas := m.generatePersonas(topic, personaCount)

	community := &Community{
		ID:        fmt.Sprintf("community_%d", time.Now().UnixNano()),
		Topic:     topic,
		Members:   personas,
		CreatedAt: time.Now(),
		Status:    "active",
	}

	return community, nil
}

// SubmitForReview submits content for mock AI community review
func (m *MockAIPersonaCommunityClient) SubmitForReview(ctx context.Context, communityID string, content string) (*CommunityReview, error) {
	// Simulate processing time
	time.Sleep(time.Millisecond * 100)

	// Generate realistic persona reviews
	personaReviews := m.generatePersonaReviews(content)

	// Calculate consensus
	consensus := m.calculateConsensus(personaReviews)

	now := time.Now()
	review := &CommunityReview{
		ReviewID:        fmt.Sprintf("review_%d", time.Now().UnixNano()),
		Topic:           "threat_analysis",
		Content:         content,
		PersonaReviews:  personaReviews,
		Consensus:       consensus,
		Recommendations: m.generateRecommendations(consensus.OverallScore),
		Metadata:        map[string]interface{}{"mock": true},
		CreatedAt:       now,
		CompletedAt:     &now,
	}

	return review, nil
}

// GetReviewStatus gets the status of a review
func (m *MockAIPersonaCommunityClient) GetReviewStatus(ctx context.Context, reviewID string) (*CommunityReview, error) {
	return nil, fmt.Errorf("mock implementation - review status not available")
}

// GetCommunityMembers gets community members
func (m *MockAIPersonaCommunityClient) GetCommunityMembers(ctx context.Context, communityID string) ([]PersonaInfo, error) {
	return m.generatePersonas("general", 3), nil
}

// generatePersonas creates realistic AI personas based on topic
func (m *MockAIPersonaCommunityClient) generatePersonas(topic string, count int) []PersonaInfo {
	var personas []PersonaInfo

	// Define persona templates based on topic
	var templates []PersonaInfo

	switch {
	case strings.Contains(topic, "threat") || strings.Contains(topic, "security"):
		templates = []PersonaInfo{
			{ID: "security_analyst_001", Name: "Security Analyst", Expertise: []string{"threat_analysis", "cybersecurity"}, Description: "Expert in cybersecurity threat analysis"},
			{ID: "forensics_expert_001", Name: "Digital Forensics Expert", Expertise: []string{"digital_forensics", "malware_analysis"}, Description: "Specialist in digital forensics and malware analysis"},
			{ID: "network_security_001", Name: "Network Security Specialist", Expertise: []string{"network_security", "intrusion_detection"}, Description: "Expert in network security and intrusion detection"},
			{ID: "incident_responder_001", Name: "Incident Response Specialist", Expertise: []string{"incident_response", "threat_hunting"}, Description: "Specialist in incident response and threat hunting"},
			{ID: "risk_assessor_001", Name: "Risk Assessment Expert", Expertise: []string{"risk_assessment", "vulnerability_analysis"}, Description: "Expert in risk assessment and vulnerability analysis"},
		}
	case strings.Contains(topic, "email"):
		templates = []PersonaInfo{
			{ID: "email_security_001", Name: "Email Security Expert", Expertise: []string{"email_security", "phishing_detection"}, Description: "Expert in email security and phishing detection"},
			{ID: "spam_analyst_001", Name: "Spam Analysis Specialist", Expertise: []string{"spam_detection", "content_analysis"}, Description: "Specialist in spam detection and email content analysis"},
			{ID: "social_engineer_001", Name: "Social Engineering Expert", Expertise: []string{"social_engineering", "human_psychology"}, Description: "Expert in social engineering tactics and human psychology"},
		}
	default:
		templates = []PersonaInfo{
			{ID: "analytical_reviewer_001", Name: "Analytical Reviewer", Expertise: []string{"analysis", "critical_thinking"}, Description: "Expert in analytical thinking and content review"},
			{ID: "content_moderator_001", Name: "Content Moderator", Expertise: []string{"content_moderation", "policy_enforcement"}, Description: "Specialist in content moderation and policy enforcement"},
			{ID: "quality_assessor_001", Name: "Quality Assessor", Expertise: []string{"quality_assessment", "evaluation"}, Description: "Expert in quality assessment and evaluation"},
		}
	}

	// Select personas up to the requested count
	for i := 0; i < count && i < len(templates); i++ {
		persona := templates[i]
		persona.Model = "gpt-4"
		personas = append(personas, persona)
	}

	// If we need more personas than templates, create variations
	for len(personas) < count {
		basePersona := templates[len(personas)%len(templates)]
		variation := basePersona
		variation.ID = fmt.Sprintf("%s_%d", basePersona.ID, len(personas))
		variation.Name = fmt.Sprintf("%s %d", basePersona.Name, len(personas)+1)
		personas = append(personas, variation)
	}

	return personas
}

// generatePersonaReviews creates realistic persona reviews
func (m *MockAIPersonaCommunityClient) generatePersonaReviews(content string) []PersonaReview {
	contentLower := strings.ToLower(content)

	// Analyze content for threat indicators
	threatScore := m.calculateThreatScore(contentLower)

	reviews := []PersonaReview{
		{
			PersonaID:   "security_analyst_001",
			PersonaName: "Security Analyst",
			Expertise:   []string{"threat_analysis", "cybersecurity"},
			Review:      m.generateSecurityAnalystReview(content, threatScore),
			Score:       threatScore + (rand.Float64()-0.5)*0.1, // Add some variation
			Confidence:  0.85 + rand.Float64()*0.1,
			Tags:        []string{"security", "analysis"},
			Metadata:    map[string]interface{}{"model": "gpt-4"},
			Timestamp:   time.Now(),
		},
		{
			PersonaID:   "forensics_expert_001",
			PersonaName: "Digital Forensics Expert",
			Expertise:   []string{"digital_forensics", "malware_analysis"},
			Review:      m.generateForensicsReview(content, threatScore),
			Score:       threatScore + (rand.Float64()-0.5)*0.15,
			Confidence:  0.80 + rand.Float64()*0.15,
			Tags:        []string{"forensics", "malware"},
			Metadata:    map[string]interface{}{"model": "gpt-4"},
			Timestamp:   time.Now(),
		},
		{
			PersonaID:   "risk_assessor_001",
			PersonaName: "Risk Assessment Expert",
			Expertise:   []string{"risk_assessment", "vulnerability_analysis"},
			Review:      m.generateRiskAssessmentReview(content, threatScore),
			Score:       threatScore + (rand.Float64()-0.5)*0.12,
			Confidence:  0.88 + rand.Float64()*0.08,
			Tags:        []string{"risk", "assessment"},
			Metadata:    map[string]interface{}{"model": "gpt-4"},
			Timestamp:   time.Now(),
		},
	}

	return reviews
}

// calculateThreatScore analyzes content and calculates a threat score
func (m *MockAIPersonaCommunityClient) calculateThreatScore(content string) float64 {
	score := 0.0

	// High threat indicators
	highThreatKeywords := []string{
		"urgent", "click here", "verify", "suspended", "bitcoin", "cryptocurrency",
		"wire transfer", "bank account", "social security", "credit card",
		"malware", "virus", "trojan", "ransomware", "phishing", "scam",
		"download", "install", "execute", "run", "admin", "administrator",
		"password", "login", "credentials", "account", "security alert",
	}

	// Medium threat indicators
	mediumThreatKeywords := []string{
		"free", "winner", "congratulations", "prize", "lottery", "inheritance",
		"investment", "opportunity", "limited time", "act now", "call now",
		"microsoft", "apple", "google", "amazon", "paypal", "ebay",
		"irs", "fbi", "police", "government", "tax", "refund",
	}

	// Low threat indicators
	lowThreatKeywords := []string{
		"unsubscribe", "marketing", "newsletter", "promotion", "sale",
		"discount", "offer", "deal", "shopping", "product",
	}

	// Count keyword occurrences
	for _, keyword := range highThreatKeywords {
		if strings.Contains(content, keyword) {
			score += 0.15
		}
	}

	for _, keyword := range mediumThreatKeywords {
		if strings.Contains(content, keyword) {
			score += 0.08
		}
	}

	for _, keyword := range lowThreatKeywords {
		if strings.Contains(content, keyword) {
			score += 0.03
		}
	}

	// Check for suspicious patterns
	if strings.Contains(content, "http://") && !strings.Contains(content, "https://") {
		score += 0.1 // Insecure links
	}

	if strings.Count(content, "!") > 3 {
		score += 0.05 // Excessive exclamation marks
	}

	if strings.Contains(content, "URGENT") || strings.Contains(content, "IMMEDIATE") {
		score += 0.1 // All caps urgency
	}

	// Normalize score to 0-1 range
	if score > 1.0 {
		score = 0.95 // Cap at 95% to leave room for uncertainty
	}

	return score
}

// generateSecurityAnalystReview creates a security analyst review
func (m *MockAIPersonaCommunityClient) generateSecurityAnalystReview(content string, threatScore float64) string {
	if threatScore > 0.8 {
		return fmt.Sprintf("HIGH THREAT DETECTED: This content exhibits multiple indicators of malicious intent including suspicious keywords and social engineering tactics. Threat score: %.2f. Recommend immediate blocking and further investigation.", threatScore)
	} else if threatScore > 0.6 {
		return fmt.Sprintf("MODERATE THREAT: Content shows concerning patterns that warrant caution. Multiple threat indicators present. Threat score: %.2f. Recommend enhanced monitoring and user education.", threatScore)
	} else if threatScore > 0.3 {
		return fmt.Sprintf("LOW THREAT: Some suspicious elements detected but overall risk appears manageable. Threat score: %.2f. Standard security protocols should be sufficient.", threatScore)
	} else {
		return fmt.Sprintf("MINIMAL THREAT: Content appears legitimate with no significant security concerns. Threat score: %.2f. No special action required.", threatScore)
	}
}

// generateForensicsReview creates a digital forensics review
func (m *MockAIPersonaCommunityClient) generateForensicsReview(content string, threatScore float64) string {
	if threatScore > 0.7 {
		return "Forensic analysis reveals patterns consistent with known attack vectors. Recommend preservation of evidence and detailed investigation of source attribution."
	} else if threatScore > 0.4 {
		return "Forensic indicators suggest potential malicious activity. Recommend monitoring and collection of additional evidence for pattern analysis."
	} else {
		return "Forensic analysis shows no significant indicators of malicious activity. Content appears to follow normal communication patterns."
	}
}

// generateRiskAssessmentReview creates a risk assessment review
func (m *MockAIPersonaCommunityClient) generateRiskAssessmentReview(content string, threatScore float64) string {
	if threatScore > 0.8 {
		return "CRITICAL RISK: High probability of successful attack if user interaction occurs. Recommend immediate containment and user notification."
	} else if threatScore > 0.5 {
		return "ELEVATED RISK: Moderate probability of user compromise. Recommend enhanced security measures and user awareness training."
	} else {
		return "ACCEPTABLE RISK: Low probability of security impact. Standard security controls are adequate for this content type."
	}
}

// calculateConsensus calculates consensus from persona reviews
func (m *MockAIPersonaCommunityClient) calculateConsensus(reviews []PersonaReview) *Consensus {
	if len(reviews) == 0 {
		return &Consensus{
			OverallScore:    0.0,
			Agreement:       0.0,
			Recommendation:  "No reviews available",
			KeyPoints:       []string{},
			ConfidenceLevel: 0.0,
		}
	}

	// Calculate average score
	totalScore := 0.0
	totalConfidence := 0.0
	for _, review := range reviews {
		totalScore += review.Score
		totalConfidence += review.Confidence
	}

	avgScore := totalScore / float64(len(reviews))
	avgConfidence := totalConfidence / float64(len(reviews))

	// Calculate agreement (inverse of variance)
	variance := 0.0
	for _, review := range reviews {
		diff := review.Score - avgScore
		variance += diff * diff
	}
	variance /= float64(len(reviews))
	agreement := 1.0 - variance // Higher variance = lower agreement

	// Generate recommendation
	var recommendation string
	if avgScore >= 0.8 {
		recommendation = "CRITICAL THREAT - Immediate action required"
	} else if avgScore >= 0.6 {
		recommendation = "HIGH THREAT - Enhanced security measures recommended"
	} else if avgScore >= 0.4 {
		recommendation = "MODERATE THREAT - Standard security protocols sufficient"
	} else if avgScore >= 0.2 {
		recommendation = "LOW THREAT - Minimal security concern"
	} else {
		recommendation = "MINIMAL THREAT - Content appears safe"
	}

	// Extract key points
	keyPoints := []string{
		fmt.Sprintf("Multiple reviewers noted: %s", m.getCommonTheme(avgScore)),
		fmt.Sprintf("Consensus confidence: %.1f%%", avgConfidence*100),
		fmt.Sprintf("Reviewer agreement: %.1f%%", agreement*100),
	}

	return &Consensus{
		OverallScore:    avgScore,
		Agreement:       agreement,
		Recommendation:  recommendation,
		KeyPoints:       keyPoints,
		ConfidenceLevel: avgConfidence,
	}
}

// getCommonTheme returns a common theme based on score
func (m *MockAIPersonaCommunityClient) getCommonTheme(score float64) string {
	if score > 0.7 {
		return "high threat indicators and malicious patterns"
	} else if score > 0.4 {
		return "suspicious elements requiring caution"
	} else {
		return "legitimate content with minimal security concerns"
	}
}

// generateRecommendations generates recommendations based on threat score
func (m *MockAIPersonaCommunityClient) generateRecommendations(score float64) []string {
	if score >= 0.8 {
		return []string{
			"Block content immediately",
			"Notify security team",
			"Investigate source attribution",
			"Update threat intelligence",
			"Review similar patterns",
		}
	} else if score >= 0.6 {
		return []string{
			"Flag for enhanced monitoring",
			"Implement additional security controls",
			"Educate users about threats",
			"Monitor for similar patterns",
		}
	} else if score >= 0.4 {
		return []string{
			"Apply standard security protocols",
			"Monitor user interactions",
			"Log for pattern analysis",
		}
	} else {
		return []string{
			"No special action required",
			"Continue normal monitoring",
		}
	}
}
