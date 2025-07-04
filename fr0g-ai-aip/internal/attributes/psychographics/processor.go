package psychographics

import (
	"fmt"
	"strings"

	"github.com/fr0g-ai/fr0g-ai-aip/internal/config"
	pb "github.com/fr0g-ai/fr0g-ai-aip/internal/grpc/pb"
)

// Processor handles psychographics attribute processing and validation
type Processor struct {
	config *config.ValidationConfig
}

// NewProcessor creates a new psychographics processor
func NewProcessor(cfg *config.ValidationConfig) *Processor {
	return &Processor{
		config: cfg,
	}
}

// ValidatePsychographics validates psychographics data
func (p *Processor) ValidatePsychographics(psycho *pb.Psychographics) []config.ValidationError {
	var errors []config.ValidationError

	// Validate personality if present
	if psycho.Personality != nil {
		personalityErrors := p.validatePersonality(psycho.Personality)
		errors = append(errors, personalityErrors...)
	}

	// Validate cognitive style
	if psycho.CognitiveStyle != "" {
		validStyles := []string{
			"analytical", "intuitive", "systematic", "creative", 
			"logical", "holistic", "detail-oriented", "big-picture",
		}
		if !p.isValidOption(psycho.CognitiveStyle, validStyles) {
			errors = append(errors, config.ValidationError{
				Field:   "cognitive_style",
				Message: "invalid cognitive style",
			})
		}
	}

	// Validate learning style
	if psycho.LearningStyle != "" {
		validStyles := []string{
			"visual", "auditory", "kinesthetic", "reading-writing",
			"multimodal", "experiential", "theoretical", "practical",
		}
		if !p.isValidOption(psycho.LearningStyle, validStyles) {
			errors = append(errors, config.ValidationError{
				Field:   "learning_style",
				Message: "invalid learning style",
			})
		}
	}

	// Validate risk tolerance
	if psycho.RiskTolerance != "" {
		validLevels := []string{"very-low", "low", "moderate", "high", "very-high"}
		if !p.isValidOption(psycho.RiskTolerance, validLevels) {
			errors = append(errors, config.ValidationError{
				Field:   "risk_tolerance",
				Message: "invalid risk tolerance level",
			})
		}
	}

	// Validate openness to change (0.0-1.0)
	if psycho.OpennessToChange < 0.0 || psycho.OpennessToChange > 1.0 {
		errors = append(errors, config.ValidationError{
			Field:   "openness_to_change",
			Message: "openness to change must be between 0.0 and 1.0",
		})
	}

	// Validate values and beliefs
	if len(psycho.Values) > 20 {
		errors = append(errors, config.ValidationError{
			Field:   "values",
			Message: "too many values specified (max 20)",
		})
	}

	if len(psycho.CoreBeliefs) > 15 {
		errors = append(errors, config.ValidationError{
			Field:   "core_beliefs",
			Message: "too many core beliefs specified (max 15)",
		})
	}

	return errors
}

// validatePersonality validates Big Five personality traits
func (p *Processor) validatePersonality(personality *pb.Personality) []config.ValidationError {
	var errors []config.ValidationError

	traits := map[string]float64{
		"openness":          personality.Openness,
		"conscientiousness": personality.Conscientiousness,
		"extraversion":      personality.Extraversion,
		"agreeableness":     personality.Agreeableness,
		"neuroticism":       personality.Neuroticism,
	}

	for trait, value := range traits {
		if value < 0.0 || value > 1.0 {
			errors = append(errors, config.ValidationError{
				Field:   trait,
				Message: fmt.Sprintf("%s must be between 0.0 and 1.0", trait),
			})
		}
	}

	return errors
}

// ProcessPsychographics processes and enriches psychographics data
func (p *Processor) ProcessPsychographics(psycho *pb.Psychographics) (*pb.Psychographics, error) {
	if psycho == nil {
		return nil, fmt.Errorf("psychographics cannot be nil")
	}

	// Validate first
	if validationErrors := p.ValidatePsychographics(psycho); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Create processed copy
	processed := &pb.Psychographics{
		Personality:      psycho.Personality,
		Values:           p.normalizeStringSlice(psycho.Values),
		CoreBeliefs:      p.normalizeStringSlice(psycho.CoreBeliefs),
		CognitiveStyle:   p.normalizeString(psycho.CognitiveStyle),
		LearningStyle:    p.normalizeString(psycho.LearningStyle),
		RiskTolerance:    p.normalizeString(psycho.RiskTolerance),
		OpennessToChange: psycho.OpennessToChange,
	}

	return processed, nil
}

// GetPersonalityProfile returns a textual description of personality
func (p *Processor) GetPersonalityProfile(personality *pb.Personality) string {
	if personality == nil {
		return "No personality data available"
	}

	var profile []string

	// Openness
	switch {
	case personality.Openness >= 0.7:
		profile = append(profile, "highly creative and open to new experiences")
	case personality.Openness >= 0.3:
		profile = append(profile, "moderately open to new ideas")
	default:
		profile = append(profile, "prefers familiar and conventional approaches")
	}

	// Conscientiousness
	switch {
	case personality.Conscientiousness >= 0.7:
		profile = append(profile, "highly organized and disciplined")
	case personality.Conscientiousness >= 0.3:
		profile = append(profile, "moderately organized")
	default:
		profile = append(profile, "flexible and spontaneous")
	}

	// Extraversion
	switch {
	case personality.Extraversion >= 0.7:
		profile = append(profile, "highly social and outgoing")
	case personality.Extraversion >= 0.3:
		profile = append(profile, "moderately social")
	default:
		profile = append(profile, "introverted and reserved")
	}

	// Agreeableness
	switch {
	case personality.Agreeableness >= 0.7:
		profile = append(profile, "highly cooperative and trusting")
	case personality.Agreeableness >= 0.3:
		profile = append(profile, "moderately agreeable")
	default:
		profile = append(profile, "competitive and skeptical")
	}

	// Neuroticism
	switch {
	case personality.Neuroticism >= 0.7:
		profile = append(profile, "emotionally sensitive")
	case personality.Neuroticism >= 0.3:
		profile = append(profile, "moderately emotionally stable")
	default:
		profile = append(profile, "emotionally resilient")
	}

	return strings.Join(profile, ", ")
}

// GetCognitiveProfile returns cognitive processing preferences
func (p *Processor) GetCognitiveProfile(psycho *pb.Psychographics) map[string]string {
	profile := make(map[string]string)

	if psycho.CognitiveStyle != "" {
		profile["cognitive_style"] = psycho.CognitiveStyle
	}

	if psycho.LearningStyle != "" {
		profile["learning_style"] = psycho.LearningStyle
	}

	if psycho.RiskTolerance != "" {
		profile["risk_tolerance"] = psycho.RiskTolerance
	}

	// Derive decision-making style from personality
	if psycho.Personality != nil {
		if psycho.Personality.Conscientiousness > 0.6 && psycho.Personality.Openness < 0.4 {
			profile["decision_style"] = "systematic"
		} else if psycho.Personality.Openness > 0.6 && psycho.Personality.Conscientiousness < 0.4 {
			profile["decision_style"] = "intuitive"
		} else {
			profile["decision_style"] = "balanced"
		}
	}

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

// GetValueCategories categorizes values into broader themes
func (p *Processor) GetValueCategories(values []string) map[string][]string {
	categories := make(map[string][]string)
	
	valueCategories := map[string]string{
		"family":        "relationships",
		"friendship":    "relationships",
		"love":          "relationships",
		"community":     "relationships",
		"success":       "achievement",
		"achievement":   "achievement",
		"ambition":      "achievement",
		"excellence":    "achievement",
		"freedom":       "autonomy",
		"independence":  "autonomy",
		"self-reliance": "autonomy",
		"creativity":    "self-expression",
		"authenticity":  "self-expression",
		"individuality": "self-expression",
		"justice":       "moral",
		"fairness":      "moral",
		"honesty":       "moral",
		"integrity":     "moral",
		"security":      "stability",
		"stability":     "stability",
		"tradition":     "stability",
		"adventure":     "stimulation",
		"excitement":    "stimulation",
		"novelty":       "stimulation",
	}

	for _, value := range values {
		normalized := p.normalizeString(value)
		if category, exists := valueCategories[normalized]; exists {
			categories[category] = append(categories[category], value)
		} else {
			categories["other"] = append(categories["other"], value)
		}
	}

	return categories
}
package psychographics

import (
	"fmt"
	"strings"

	"github.com/fr0g-ai/fr0g-ai-aip/internal/config"
	pb "github.com/fr0g-ai/fr0g-ai-aip/internal/grpc/pb"
)

// Processor handles psychographics attribute processing and validation
type Processor struct {
	config *config.ValidationConfig
}

// NewProcessor creates a new psychographics processor
func NewProcessor(cfg *config.ValidationConfig) *Processor {
	return &Processor{
		config: cfg,
	}
}

// ValidatePsychographics validates psychographics data
func (p *Processor) ValidatePsychographics(psycho *pb.Psychographics) []config.ValidationError {
	var errors []config.ValidationError

	// Validate personality if present
	if psycho.Personality != nil {
		personalityErrors := p.validatePersonality(psycho.Personality)
		errors = append(errors, personalityErrors...)
	}

	// Validate cognitive style
	if psycho.CognitiveStyle != "" {
		validStyles := []string{
			"analytical", "intuitive", "systematic", "creative", 
			"logical", "holistic", "detail-oriented", "big-picture",
		}
		if !p.isValidOption(psycho.CognitiveStyle, validStyles) {
			errors = append(errors, config.ValidationError{
				Field:   "cognitive_style",
				Message: "invalid cognitive style",
			})
		}
	}

	// Validate learning style
	if psycho.LearningStyle != "" {
		validStyles := []string{
			"visual", "auditory", "kinesthetic", "reading-writing",
			"multimodal", "experiential", "theoretical", "practical",
		}
		if !p.isValidOption(psycho.LearningStyle, validStyles) {
			errors = append(errors, config.ValidationError{
				Field:   "learning_style",
				Message: "invalid learning style",
			})
		}
	}

	// Validate risk tolerance
	if psycho.RiskTolerance != "" {
		validLevels := []string{"very-low", "low", "moderate", "high", "very-high"}
		if !p.isValidOption(psycho.RiskTolerance, validLevels) {
			errors = append(errors, config.ValidationError{
				Field:   "risk_tolerance",
				Message: "invalid risk tolerance level",
			})
		}
	}

	// Validate openness to change (0.0-1.0)
	if psycho.OpennessToChange < 0.0 || psycho.OpennessToChange > 1.0 {
		errors = append(errors, config.ValidationError{
			Field:   "openness_to_change",
			Message: "openness to change must be between 0.0 and 1.0",
		})
	}

	// Validate values and beliefs
	if len(psycho.Values) > 20 {
		errors = append(errors, config.ValidationError{
			Field:   "values",
			Message: "too many values specified (max 20)",
		})
	}

	if len(psycho.CoreBeliefs) > 15 {
		errors = append(errors, config.ValidationError{
			Field:   "core_beliefs",
			Message: "too many core beliefs specified (max 15)",
		})
	}

	return errors
}

// validatePersonality validates Big Five personality traits
func (p *Processor) validatePersonality(personality *pb.Personality) []config.ValidationError {
	var errors []config.ValidationError

	traits := map[string]float64{
		"openness":          personality.Openness,
		"conscientiousness": personality.Conscientiousness,
		"extraversion":      personality.Extraversion,
		"agreeableness":     personality.Agreeableness,
		"neuroticism":       personality.Neuroticism,
	}

	for trait, value := range traits {
		if value < 0.0 || value > 1.0 {
			errors = append(errors, config.ValidationError{
				Field:   trait,
				Message: fmt.Sprintf("%s must be between 0.0 and 1.0", trait),
			})
		}
	}

	return errors
}

// ProcessPsychographics processes and enriches psychographics data
func (p *Processor) ProcessPsychographics(psycho *pb.Psychographics) (*pb.Psychographics, error) {
	if psycho == nil {
		return nil, fmt.Errorf("psychographics cannot be nil")
	}

	// Validate first
	if validationErrors := p.ValidatePsychographics(psycho); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Create processed copy
	processed := &pb.Psychographics{
		Personality:      psycho.Personality,
		Values:           p.normalizeStringSlice(psycho.Values),
		CoreBeliefs:      p.normalizeStringSlice(psycho.CoreBeliefs),
		CognitiveStyle:   p.normalizeString(psycho.CognitiveStyle),
		LearningStyle:    p.normalizeString(psycho.LearningStyle),
		RiskTolerance:    p.normalizeString(psycho.RiskTolerance),
		OpennessToChange: psycho.OpennessToChange,
	}

	return processed, nil
}

// GetPersonalityProfile returns a textual description of personality
func (p *Processor) GetPersonalityProfile(personality *pb.Personality) string {
	if personality == nil {
		return "No personality data available"
	}

	var profile []string

	// Openness
	switch {
	case personality.Openness >= 0.7:
		profile = append(profile, "highly creative and open to new experiences")
	case personality.Openness >= 0.3:
		profile = append(profile, "moderately open to new ideas")
	default:
		profile = append(profile, "prefers familiar and conventional approaches")
	}

	// Conscientiousness
	switch {
	case personality.Conscientiousness >= 0.7:
		profile = append(profile, "highly organized and disciplined")
	case personality.Conscientiousness >= 0.3:
		profile = append(profile, "moderately organized")
	default:
		profile = append(profile, "flexible and spontaneous")
	}

	// Extraversion
	switch {
	case personality.Extraversion >= 0.7:
		profile = append(profile, "highly social and outgoing")
	case personality.Extraversion >= 0.3:
		profile = append(profile, "moderately social")
	default:
		profile = append(profile, "introverted and reserved")
	}

	// Agreeableness
	switch {
	case personality.Agreeableness >= 0.7:
		profile = append(profile, "highly cooperative and trusting")
	case personality.Agreeableness >= 0.3:
		profile = append(profile, "moderately agreeable")
	default:
		profile = append(profile, "competitive and skeptical")
	}

	// Neuroticism
	switch {
	case personality.Neuroticism >= 0.7:
		profile = append(profile, "emotionally sensitive")
	case personality.Neuroticism >= 0.3:
		profile = append(profile, "moderately emotionally stable")
	default:
		profile = append(profile, "emotionally resilient")
	}

	return strings.Join(profile, ", ")
}

// GetCognitiveProfile returns cognitive processing preferences
func (p *Processor) GetCognitiveProfile(psycho *pb.Psychographics) map[string]string {
	profile := make(map[string]string)

	if psycho.CognitiveStyle != "" {
		profile["cognitive_style"] = psycho.CognitiveStyle
	}

	if psycho.LearningStyle != "" {
		profile["learning_style"] = psycho.LearningStyle
	}

	if psycho.RiskTolerance != "" {
		profile["risk_tolerance"] = psycho.RiskTolerance
	}

	// Derive decision-making style from personality
	if psycho.Personality != nil {
		if psycho.Personality.Conscientiousness > 0.6 && psycho.Personality.Openness < 0.4 {
			profile["decision_style"] = "systematic"
		} else if psycho.Personality.Openness > 0.6 && psycho.Personality.Conscientiousness < 0.4 {
			profile["decision_style"] = "intuitive"
		} else {
			profile["decision_style"] = "balanced"
		}
	}

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

// GetValueCategories categorizes values into broader themes
func (p *Processor) GetValueCategories(values []string) map[string][]string {
	categories := make(map[string][]string)
	
	valueCategories := map[string]string{
		"family":        "relationships",
		"friendship":    "relationships",
		"love":          "relationships",
		"community":     "relationships",
		"success":       "achievement",
		"achievement":   "achievement",
		"ambition":      "achievement",
		"excellence":    "achievement",
		"freedom":       "autonomy",
		"independence":  "autonomy",
		"self-reliance": "autonomy",
		"creativity":    "self-expression",
		"authenticity":  "self-expression",
		"individuality": "self-expression",
		"justice":       "moral",
		"fairness":      "moral",
		"honesty":       "moral",
		"integrity":     "moral",
		"security":      "stability",
		"stability":     "stability",
		"tradition":     "stability",
		"adventure":     "stimulation",
		"excitement":    "stimulation",
		"novelty":       "stimulation",
	}

	for _, value := range values {
		normalized := p.normalizeString(value)
		if category, exists := valueCategories[normalized]; exists {
			categories[category] = append(categories[category], value)
		} else {
			categories["other"] = append(categories["other"], value)
		}
	}

	return categories
}
