package behavioral

import (
	"fmt"
	"strings"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/config"
	pb "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/grpc/pb"
)

// Processor handles behavioral tendencies attribute processing and validation
type Processor struct {
	config *config.ValidationConfig
}

// NewProcessor creates a new behavioral processor
func NewProcessor(cfg *config.ValidationConfig) *Processor {
	return &Processor{
		config: cfg,
	}
}

// ValidateBehavioralTendencies validates behavioral tendencies data
func (p *Processor) ValidateBehavioralTendencies(behavioral *pb.BehavioralTendencies) []config.ValidationError {
	var errors []config.ValidationError

	// Validate decision making style
	if behavioral.DecisionMaking != "" {
		validStyles := []string{
			"analytical", "intuitive", "collaborative", "decisive",
			"cautious", "impulsive", "data-driven", "gut-feeling",
			"consensus-seeking", "independent", "systematic", "flexible",
		}
		if !p.isValidOption(behavioral.DecisionMaking, validStyles) {
			errors = append(errors, config.ValidationError{
				Field:   "decision_making",
				Message: "invalid decision making style",
			})
		}
	}

	// Validate conflict resolution style
	if behavioral.ConflictResolution != "" {
		validStyles := []string{
			"collaborative", "competitive", "accommodating", "avoiding",
			"compromising", "assertive", "diplomatic", "direct",
			"mediating", "confrontational", "passive", "aggressive",
		}
		if !p.isValidOption(behavioral.ConflictResolution, validStyles) {
			errors = append(errors, config.ValidationError{
				Field:   "conflict_resolution",
				Message: "invalid conflict resolution style",
			})
		}
	}

	// Validate communication style
	if behavioral.CommunicationStyle != "" {
		validStyles := []string{
			"direct", "indirect", "formal", "informal", "assertive",
			"passive", "aggressive", "diplomatic", "expressive",
			"reserved", "verbose", "concise", "emotional", "logical",
		}
		if !p.isValidOption(behavioral.CommunicationStyle, validStyles) {
			errors = append(errors, config.ValidationError{
				Field:   "communication_style",
				Message: "invalid communication style",
			})
		}
	}

	// Validate leadership style
	if behavioral.LeadershipStyle != "" {
		validStyles := []string{
			"democratic", "autocratic", "laissez-faire", "transformational",
			"transactional", "servant", "authentic", "coaching",
			"visionary", "pacesetting", "commanding", "affiliative",
		}
		if !p.isValidOption(behavioral.LeadershipStyle, validStyles) {
			errors = append(errors, config.ValidationError{
				Field:   "leadership_style",
				Message: "invalid leadership style",
			})
		}
	}

	// Validate coping mechanisms list
	if len(behavioral.CopingMechanisms) > 20 {
		errors = append(errors, config.ValidationError{
			Field:   "coping_mechanisms",
			Message: "too many coping mechanisms specified (max 20)",
		})
	}

	// Validate stress response
	if behavioral.StressResponse != "" {
		validResponses := []string{
			"fight", "flight", "freeze", "fawn", "problem-solving",
			"emotion-focused", "avoidance", "seeking-support",
			"self-soothing", "distraction", "rumination", "action-oriented",
		}
		if !p.isValidOption(behavioral.StressResponse, validResponses) {
			errors = append(errors, config.ValidationError{
				Field:   "stress_response",
				Message: "invalid stress response",
			})
		}
	}

	return errors
}

// ProcessBehavioralTendencies processes and enriches behavioral data
func (p *Processor) ProcessBehavioralTendencies(behavioral *pb.BehavioralTendencies) (*pb.BehavioralTendencies, error) {
	if behavioral == nil {
		return nil, fmt.Errorf("behavioral tendencies cannot be nil")
	}

	// Validate first
	if validationErrors := p.ValidateBehavioralTendencies(behavioral); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Create processed copy
	processed := &pb.BehavioralTendencies{
		DecisionMaking:     p.normalizeString(behavioral.DecisionMaking),
		ConflictResolution: p.normalizeString(behavioral.ConflictResolution),
		CommunicationStyle: p.normalizeString(behavioral.CommunicationStyle),
		LeadershipStyle:    p.normalizeString(behavioral.LeadershipStyle),
		CopingMechanisms:   p.normalizeStringSlice(behavioral.CopingMechanisms),
		StressResponse:     p.normalizeString(behavioral.StressResponse),
	}

	return processed, nil
}

// GetBehavioralProfile returns comprehensive behavioral profile
func (p *Processor) GetBehavioralProfile(behavioral *pb.BehavioralTendencies) map[string]interface{} {
	profile := make(map[string]interface{})

	// Decision making analysis
	if behavioral.DecisionMaking != "" {
		profile["decision_making"] = behavioral.DecisionMaking
		profile["decision_style_category"] = p.getDecisionStyleCategory(behavioral.DecisionMaking)
	}

	// Conflict resolution analysis
	if behavioral.ConflictResolution != "" {
		profile["conflict_resolution"] = behavioral.ConflictResolution
		profile["conflict_approach"] = p.getConflictApproach(behavioral.ConflictResolution)
	}

	// Communication analysis
	if behavioral.CommunicationStyle != "" {
		profile["communication_style"] = behavioral.CommunicationStyle
		profile["communication_effectiveness"] = p.getCommunicationEffectiveness(behavioral.CommunicationStyle)
	}

	// Leadership analysis
	if behavioral.LeadershipStyle != "" {
		profile["leadership_style"] = behavioral.LeadershipStyle
		profile["leadership_approach"] = p.getLeadershipApproach(behavioral.LeadershipStyle)
	}

	// Coping and stress analysis
	if len(behavioral.CopingMechanisms) > 0 {
		profile["coping_categories"] = p.categorizeCopingMechanisms(behavioral.CopingMechanisms)
		profile["coping_effectiveness"] = p.getCopingEffectiveness(behavioral.CopingMechanisms)
	}

	if behavioral.StressResponse != "" {
		profile["stress_response"] = behavioral.StressResponse
		profile["stress_management_style"] = p.getStressManagementStyle(behavioral.StressResponse)
	}

	// Overall behavioral assessment
	profile["behavioral_adaptability"] = p.getBehavioralAdaptability(behavioral)
	profile["interpersonal_style"] = p.getInterpersonalStyle(behavioral)

	return profile
}

// Helper methods

func (p *Processor) isValidOption(value string, validOptions []string) bool {
	normalized := p.normalizeString(value)
	for _, option := range validOptions {
		if normalized == option {
			return true
		}
	}
	return false
}

func (p *Processor) normalizeString(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func (p *Processor) normalizeStringSlice(slice []string) []string {
	var normalized []string
	for _, s := range slice {
		if trimmed := strings.TrimSpace(s); trimmed != "" {
			normalized = append(normalized, strings.ToLower(trimmed))
		}
	}
	return normalized
}

func (p *Processor) getDecisionStyleCategory(decisionMaking string) string {
	normalized := p.normalizeString(decisionMaking)

	categoryMap := map[string]string{
		"analytical":        "rational",
		"data-driven":       "rational",
		"systematic":        "rational",
		"cautious":          "rational",
		"intuitive":         "intuitive",
		"gut-feeling":       "intuitive",
		"impulsive":         "intuitive",
		"flexible":          "adaptive",
		"collaborative":     "social",
		"consensus-seeking": "social",
		"decisive":          "directive",
		"independent":       "directive",
	}

	if category, exists := categoryMap[normalized]; exists {
		return category
	}
	return "unknown"
}

func (p *Processor) getConflictApproach(conflictResolution string) string {
	normalized := p.normalizeString(conflictResolution)

	approachMap := map[string]string{
		"collaborative":   "win-win",
		"compromising":    "win-win",
		"mediating":       "win-win",
		"diplomatic":      "win-win",
		"competitive":     "win-lose",
		"confrontational": "win-lose",
		"aggressive":      "win-lose",
		"accommodating":   "lose-win",
		"passive":         "lose-win",
		"avoiding":        "lose-lose",
		"assertive":       "balanced",
		"direct":          "balanced",
	}

	if approach, exists := approachMap[normalized]; exists {
		return approach
	}
	return "unknown"
}

func (p *Processor) getCommunicationEffectiveness(communicationStyle string) string {
	normalized := p.normalizeString(communicationStyle)

	// Rate effectiveness based on style characteristics
	highEffective := []string{"assertive", "direct", "diplomatic", "concise", "logical"}
	moderateEffective := []string{"formal", "informal", "expressive"}
	lowEffective := []string{"passive", "aggressive", "verbose", "indirect"}

	for _, style := range highEffective {
		if normalized == style {
			return "high"
		}
	}

	for _, style := range moderateEffective {
		if normalized == style {
			return "moderate"
		}
	}

	for _, style := range lowEffective {
		if normalized == style {
			return "needs-improvement"
		}
	}

	return "unknown"
}

func (p *Processor) getLeadershipApproach(leadershipStyle string) string {
	normalized := p.normalizeString(leadershipStyle)

	approachMap := map[string]string{
		"democratic":       "participative",
		"collaborative":    "participative",
		"servant":          "participative",
		"coaching":         "participative",
		"autocratic":       "directive",
		"commanding":       "directive",
		"pacesetting":      "directive",
		"transformational": "inspirational",
		"visionary":        "inspirational",
		"authentic":        "inspirational",
		"transactional":    "transactional",
		"laissez-faire":    "hands-off",
		"affiliative":      "relationship-focused",
	}

	if approach, exists := approachMap[normalized]; exists {
		return approach
	}
	return "unknown"
}

func (p *Processor) categorizeCopingMechanisms(mechanisms []string) map[string][]string {
	categories := make(map[string][]string)

	mechanismCategories := map[string]string{
		"exercise":        "physical",
		"sports":          "physical",
		"running":         "physical",
		"yoga":            "physical",
		"meditation":      "mindfulness",
		"mindfulness":     "mindfulness",
		"breathing":       "mindfulness",
		"relaxation":      "mindfulness",
		"talking":         "social",
		"friends":         "social",
		"family":          "social",
		"support-group":   "social",
		"therapy":         "professional",
		"counseling":      "professional",
		"medication":      "professional",
		"reading":         "cognitive",
		"journaling":      "cognitive",
		"problem-solving": "cognitive",
		"planning":        "cognitive",
		"music":           "creative",
		"art":             "creative",
		"writing":         "creative",
		"cooking":         "creative",
		"alcohol":         "substance",
		"smoking":         "substance",
		"drugs":           "substance",
		"shopping":        "behavioral",
		"eating":          "behavioral",
		"sleeping":        "behavioral",
		"avoidance":       "avoidance",
		"denial":          "avoidance",
		"isolation":       "avoidance",
	}

	for _, mechanism := range mechanisms {
		normalized := p.normalizeString(mechanism)
		found := false
		for key, category := range mechanismCategories {
			if strings.Contains(normalized, key) {
				categories[category] = append(categories[category], mechanism)
				found = true
				break
			}
		}
		if !found {
			categories["other"] = append(categories["other"], mechanism)
		}
	}

	return categories
}

func (p *Processor) getCopingEffectiveness(mechanisms []string) string {
	categories := p.categorizeCopingMechanisms(mechanisms)

	// Count healthy vs unhealthy coping mechanisms
	healthyCount := 0
	unhealthyCount := 0

	healthyCategories := []string{"physical", "mindfulness", "social", "professional", "cognitive", "creative"}
	unhealthyCategories := []string{"substance", "avoidance"}

	for _, category := range healthyCategories {
		if items, exists := categories[category]; exists {
			healthyCount += len(items)
		}
	}

	for _, category := range unhealthyCategories {
		if items, exists := categories[category]; exists {
			unhealthyCount += len(items)
		}
	}

	if healthyCount > unhealthyCount*2 {
		return "very-effective"
	} else if healthyCount > unhealthyCount {
		return "effective"
	} else if healthyCount == unhealthyCount {
		return "mixed"
	} else if unhealthyCount > healthyCount {
		return "concerning"
	}

	return "unknown"
}

func (p *Processor) getStressManagementStyle(stressResponse string) string {
	normalized := p.normalizeString(stressResponse)

	styleMap := map[string]string{
		"problem-solving": "proactive",
		"action-oriented": "proactive",
		"seeking-support": "social",
		"emotion-focused": "emotional",
		"self-soothing":   "emotional",
		"fight":           "confrontational",
		"flight":          "avoidant",
		"freeze":          "avoidant",
		"fawn":            "accommodating",
		"avoidance":       "avoidant",
		"distraction":     "avoidant",
		"rumination":      "maladaptive",
	}

	if style, exists := styleMap[normalized]; exists {
		return style
	}
	return "unknown"
}

func (p *Processor) getBehavioralAdaptability(behavioral *pb.BehavioralTendencies) string {
	adaptabilityScore := 0

	// Check for flexible decision making
	if strings.Contains(p.normalizeString(behavioral.DecisionMaking), "flexible") {
		adaptabilityScore += 2
	}

	// Check for collaborative conflict resolution
	collaborative := []string{"collaborative", "diplomatic", "compromising"}
	for _, style := range collaborative {
		if p.normalizeString(behavioral.ConflictResolution) == style {
			adaptabilityScore += 2
			break
		}
	}

	// Check for balanced communication
	balanced := []string{"assertive", "diplomatic"}
	for _, style := range balanced {
		if p.normalizeString(behavioral.CommunicationStyle) == style {
			adaptabilityScore += 1
			break
		}
	}

	// Check for adaptive leadership
	adaptive := []string{"democratic", "coaching", "transformational"}
	for _, style := range adaptive {
		if p.normalizeString(behavioral.LeadershipStyle) == style {
			adaptabilityScore += 1
			break
		}
	}

	switch {
	case adaptabilityScore >= 5:
		return "very-high"
	case adaptabilityScore >= 3:
		return "high"
	case adaptabilityScore >= 2:
		return "moderate"
	case adaptabilityScore >= 1:
		return "low"
	default:
		return "very-low"
	}
}

func (p *Processor) getInterpersonalStyle(behavioral *pb.BehavioralTendencies) string {
	// Analyze interpersonal tendencies across different behavioral aspects
	socialScore := 0

	// Collaborative decision making
	if strings.Contains(p.normalizeString(behavioral.DecisionMaking), "collaborative") {
		socialScore += 2
	}

	// Collaborative conflict resolution
	collaborative := []string{"collaborative", "diplomatic", "accommodating"}
	for _, style := range collaborative {
		if p.normalizeString(behavioral.ConflictResolution) == style {
			socialScore += 2
			break
		}
	}

	// Social communication styles
	social := []string{"diplomatic", "expressive", "informal"}
	for _, style := range social {
		if p.normalizeString(behavioral.CommunicationStyle) == style {
			socialScore += 1
			break
		}
	}

	// People-oriented leadership
	peopleOriented := []string{"democratic", "servant", "coaching", "affiliative"}
	for _, style := range peopleOriented {
		if p.normalizeString(behavioral.LeadershipStyle) == style {
			socialScore += 1
			break
		}
	}

	switch {
	case socialScore >= 5:
		return "highly-social"
	case socialScore >= 3:
		return "social"
	case socialScore >= 2:
		return "moderately-social"
	case socialScore >= 1:
		return "somewhat-social"
	default:
		return "independent"
	}
}
