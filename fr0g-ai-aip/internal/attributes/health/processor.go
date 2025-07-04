package health

import (
	"fmt"
	"strings"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/config"
	pb "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/grpc/pb"
)

// Processor handles health attribute processing and validation
type Processor struct {
	config *config.ValidationConfig
}

// NewProcessor creates a new health processor
func NewProcessor(cfg *config.ValidationConfig) *Processor {
	return &Processor{
		config: cfg,
	}
}

// ValidateHealth validates health data
func (p *Processor) ValidateHealth(health *pb.Health) []config.ValidationError {
	var errors []config.ValidationError

	// Validate physical health status
	if health.PhysicalHealth != "" {
		validStatuses := []string{
			"excellent", "very-good", "good", "fair", "poor",
			"disabled", "chronic-condition", "recovering",
		}
		if !p.isValidOption(health.PhysicalHealth, validStatuses) {
			errors = append(errors, config.ValidationError{
				Field:   "physical_health",
				Message: "invalid physical health status",
			})
		}
	}

	// Validate mental health status
	if health.MentalHealth != "" {
		validStatuses := []string{
			"excellent", "very-good", "good", "fair", "poor",
			"seeking-help", "in-therapy", "medicated", "stable",
		}
		if !p.isValidOption(health.MentalHealth, validStatuses) {
			errors = append(errors, config.ValidationError{
				Field:   "mental_health",
				Message: "invalid mental health status",
			})
		}
	}

	// Validate fitness level
	if health.FitnessLevel != "" {
		validLevels := []string{
			"sedentary", "lightly-active", "moderately-active", 
			"very-active", "extremely-active", "athlete",
		}
		if !p.isValidOption(health.FitnessLevel, validLevels) {
			errors = append(errors, config.ValidationError{
				Field:   "fitness_level",
				Message: "invalid fitness level",
			})
		}
	}

	// Validate exercise frequency
	if health.ExerciseFrequency != "" {
		validFrequencies := []string{
			"never", "rarely", "weekly", "few-times-week", 
			"daily", "multiple-daily", "professional",
		}
		if !p.isValidOption(health.ExerciseFrequency, validFrequencies) {
			errors = append(errors, config.ValidationError{
				Field:   "exercise_frequency",
				Message: "invalid exercise frequency",
			})
		}
	}

	// Validate diet type
	if health.DietType != "" {
		validDiets := []string{
			"omnivore", "vegetarian", "vegan", "pescatarian",
			"keto", "paleo", "mediterranean", "low-carb",
			"low-fat", "intermittent-fasting", "raw", "other",
		}
		if !p.isValidOption(health.DietType, validDiets) {
			errors = append(errors, config.ValidationError{
				Field:   "diet_type",
				Message: "invalid diet type",
			})
		}
	}

	// Validate sleep quality
	if health.SleepQuality != "" {
		validQualities := []string{
			"excellent", "good", "fair", "poor", "very-poor",
			"insomnia", "sleep-disorder", "irregular",
		}
		if !p.isValidOption(health.SleepQuality, validQualities) {
			errors = append(errors, config.ValidationError{
				Field:   "sleep_quality",
				Message: "invalid sleep quality",
			})
		}
	}

	// Validate stress level
	if health.StressLevel != "" {
		validLevels := []string{
			"very-low", "low", "moderate", "high", "very-high",
			"chronic", "managed", "overwhelming",
		}
		if !p.isValidOption(health.StressLevel, validLevels) {
			errors = append(errors, config.ValidationError{
				Field:   "stress_level",
				Message: "invalid stress level",
			})
		}
	}

	// Validate substance use
	if health.SubstanceUse != "" {
		validUse := []string{
			"none", "social-drinking", "regular-drinking", "heavy-drinking",
			"social-smoking", "regular-smoking", "heavy-smoking",
			"recreational-drugs", "prescription-drugs", "recovering",
		}
		if !p.isValidOption(health.SubstanceUse, validUse) {
			errors = append(errors, config.ValidationError{
				Field:   "substance_use",
				Message: "invalid substance use",
			})
		}
	}

	// Validate medical conditions list
	if len(health.MedicalConditions) > 20 {
		errors = append(errors, config.ValidationError{
			Field:   "medical_conditions",
			Message: "too many medical conditions specified (max 20)",
		})
	}

	// Validate medications list
	if len(health.Medications) > 25 {
		errors = append(errors, config.ValidationError{
			Field:   "medications",
			Message: "too many medications specified (max 25)",
		})
	}

	// Validate allergies list
	if len(health.Allergies) > 15 {
		errors = append(errors, config.ValidationError{
			Field:   "allergies",
			Message: "too many allergies specified (max 15)",
		})
	}

	return errors
}

// ProcessHealth processes and enriches health data
func (p *Processor) ProcessHealth(health *pb.Health) (*pb.Health, error) {
	if health == nil {
		return nil, fmt.Errorf("health data cannot be nil")
	}

	// Validate first
	if validationErrors := p.ValidateHealth(health); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Create processed copy
	processed := &pb.Health{
		PhysicalHealth:      p.normalizeString(health.PhysicalHealth),
		MentalHealth:        p.normalizeString(health.MentalHealth),
		FitnessLevel:        p.normalizeString(health.FitnessLevel),
		ExerciseFrequency:   p.normalizeString(health.ExerciseFrequency),
		DietType:            p.normalizeString(health.DietType),
		SleepQuality:        p.normalizeString(health.SleepQuality),
		StressLevel:         p.normalizeString(health.StressLevel),
		SubstanceUse:        p.normalizeString(health.SubstanceUse),
		MedicalConditions:   p.normalizeStringSlice(health.MedicalConditions),
		Medications:         p.normalizeStringSlice(health.Medications),
		Allergies:           p.normalizeStringSlice(health.Allergies),
	}

	return processed, nil
}

// GetHealthProfile returns comprehensive health profile information
func (p *Processor) GetHealthProfile(health *pb.Health) map[string]interface{} {
	profile := make(map[string]interface{})

	// Overall health assessment
	profile["overall_health"] = p.getOverallHealth(health)
	profile["health_risk_factors"] = p.getHealthRiskFactors(health)
	profile["wellness_score"] = p.getWellnessScore(health)

	// Physical health
	if health.PhysicalHealth != "" {
		profile["physical_health"] = health.PhysicalHealth
	}

	// Mental health
	if health.MentalHealth != "" {
		profile["mental_health"] = health.MentalHealth
		profile["mental_health_support"] = p.getMentalHealthSupport(health.MentalHealth)
	}

	// Fitness and activity
	if health.FitnessLevel != "" || health.ExerciseFrequency != "" {
		profile["fitness_profile"] = p.getFitnessProfile(health)
	}

	// Lifestyle factors
	profile["lifestyle_factors"] = p.getLifestyleFactors(health)

	// Medical information
	if len(health.MedicalConditions) > 0 {
		profile["condition_categories"] = p.categorizeConditions(health.MedicalConditions)
		profile["chronic_conditions"] = p.hasChronicConditions(health.MedicalConditions)
	}

	if len(health.Allergies) > 0 {
		profile["allergy_categories"] = p.categorizeAllergies(health.Allergies)
		profile["allergy_severity"] = p.getAllergySeverity(health.Allergies)
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

func (p *Processor) getOverallHealth(health *pb.Health) string {
	score := 0
	factors := 0

	// Physical health
	if health.PhysicalHealth != "" {
		factors++
		switch p.normalizeString(health.PhysicalHealth) {
		case "excellent":
			score += 5
		case "very-good":
			score += 4
		case "good":
			score += 3
		case "fair":
			score += 2
		case "poor":
			score += 1
		}
	}

	// Mental health
	if health.MentalHealth != "" {
		factors++
		switch p.normalizeString(health.MentalHealth) {
		case "excellent":
			score += 5
		case "very-good", "stable":
			score += 4
		case "good":
			score += 3
		case "fair", "seeking-help", "in-therapy":
			score += 2
		case "poor":
			score += 1
		}
	}

	// Fitness level
	if health.FitnessLevel != "" {
		factors++
		switch p.normalizeString(health.FitnessLevel) {
		case "athlete", "extremely-active":
			score += 5
		case "very-active":
			score += 4
		case "moderately-active":
			score += 3
		case "lightly-active":
			score += 2
		case "sedentary":
			score += 1
		}
	}

	if factors == 0 {
		return "unknown"
	}

	average := float64(score) / float64(factors)
	switch {
	case average >= 4.5:
		return "excellent"
	case average >= 3.5:
		return "very-good"
	case average >= 2.5:
		return "good"
	case average >= 1.5:
		return "fair"
	default:
		return "poor"
	}
}

func (p *Processor) getHealthRiskFactors(health *pb.Health) []string {
	var risks []string

	// Substance use risks
	substanceUse := p.normalizeString(health.SubstanceUse)
	if strings.Contains(substanceUse, "heavy") {
		risks = append(risks, "heavy-substance-use")
	} else if strings.Contains(substanceUse, "regular") && !strings.Contains(substanceUse, "prescription") {
		risks = append(risks, "regular-substance-use")
	}

	// Sedentary lifestyle
	if p.normalizeString(health.FitnessLevel) == "sedentary" {
		risks = append(risks, "sedentary-lifestyle")
	}

	// Poor sleep
	sleepQuality := p.normalizeString(health.SleepQuality)
	if sleepQuality == "poor" || sleepQuality == "very-poor" || strings.Contains(sleepQuality, "insomnia") {
		risks = append(risks, "poor-sleep")
	}

	// High stress
	stressLevel := p.normalizeString(health.StressLevel)
	if stressLevel == "high" || stressLevel == "very-high" || stressLevel == "chronic" {
		risks = append(risks, "high-stress")
	}

	// Chronic conditions
	if p.hasChronicConditions(health.MedicalConditions) {
		risks = append(risks, "chronic-conditions")
	}

	return risks
}

func (p *Processor) getWellnessScore(health *pb.Health) int {
	score := 50 // Base score

	// Physical health
	switch p.normalizeString(health.PhysicalHealth) {
	case "excellent":
		score += 20
	case "very-good":
		score += 15
	case "good":
		score += 10
	case "fair":
		score += 5
	case "poor":
		score -= 10
	}

	// Mental health
	switch p.normalizeString(health.MentalHealth) {
	case "excellent":
		score += 20
	case "very-good", "stable":
		score += 15
	case "good":
		score += 10
	case "fair":
		score += 5
	case "poor":
		score -= 10
	}

	// Fitness level
	switch p.normalizeString(health.FitnessLevel) {
	case "athlete", "extremely-active":
		score += 15
	case "very-active":
		score += 10
	case "moderately-active":
		score += 5
	case "sedentary":
		score -= 10
	}

	// Apply risk factor penalties
	risks := p.getHealthRiskFactors(health)
	score -= len(risks) * 5

	// Ensure score is within bounds
	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
	}

	return score
}

func (p *Processor) getMentalHealthSupport(mentalHealth string) string {
	normalized := p.normalizeString(mentalHealth)
	switch {
	case strings.Contains(normalized, "therapy"):
		return "professional-therapy"
	case strings.Contains(normalized, "medicated"):
		return "medication"
	case strings.Contains(normalized, "seeking"):
		return "seeking-help"
	case normalized == "stable":
		return "stable-managed"
	default:
		return "unknown"
	}
}

func (p *Processor) getFitnessProfile(health *pb.Health) map[string]string {
	profile := make(map[string]string)

	if health.FitnessLevel != "" {
		profile["fitness_level"] = health.FitnessLevel
	}

	if health.ExerciseFrequency != "" {
		profile["exercise_frequency"] = health.ExerciseFrequency
		profile["activity_consistency"] = p.getActivityConsistency(health.ExerciseFrequency)
	}

	return profile
}

func (p *Processor) getActivityConsistency(frequency string) string {
	normalized := p.normalizeString(frequency)
	switch normalized {
	case "daily", "multiple-daily", "professional":
		return "very-consistent"
	case "few-times-week":
		return "consistent"
	case "weekly":
		return "somewhat-consistent"
	case "rarely":
		return "inconsistent"
	case "never":
		return "none"
	default:
		return "unknown"
	}
}

func (p *Processor) getLifestyleFactors(health *pb.Health) map[string]string {
	factors := make(map[string]string)

	if health.DietType != "" {
		factors["diet_type"] = health.DietType
		factors["diet_category"] = p.getDietCategory(health.DietType)
	}

	if health.SleepQuality != "" {
		factors["sleep_quality"] = health.SleepQuality
	}

	if health.StressLevel != "" {
		factors["stress_level"] = health.StressLevel
		factors["stress_management"] = p.getStressManagement(health.StressLevel)
	}

	if health.SubstanceUse != "" {
		factors["substance_use"] = health.SubstanceUse
		factors["substance_risk"] = p.getSubstanceRisk(health.SubstanceUse)
	}

	return factors
}

func (p *Processor) getDietCategory(dietType string) string {
	normalized := p.normalizeString(dietType)
	switch {
	case normalized == "vegan" || normalized == "vegetarian" || normalized == "pescatarian":
		return "plant-based"
	case normalized == "keto" || normalized == "paleo" || normalized == "low-carb":
		return "low-carb"
	case normalized == "mediterranean":
		return "balanced"
	case normalized == "raw":
		return "specialized"
	default:
		return "standard"
	}
}

func (p *Processor) getStressManagement(stressLevel string) string {
	normalized := p.normalizeString(stressLevel)
	switch normalized {
	case "managed":
		return "well-managed"
	case "very-low", "low":
		return "low-stress"
	case "chronic", "overwhelming":
		return "needs-intervention"
	default:
		return "moderate-management"
	}
}

func (p *Processor) getSubstanceRisk(substanceUse string) string {
	normalized := p.normalizeString(substanceUse)
	switch {
	case normalized == "none":
		return "no-risk"
	case strings.Contains(normalized, "social"):
		return "low-risk"
	case strings.Contains(normalized, "regular"):
		return "moderate-risk"
	case strings.Contains(normalized, "heavy"):
		return "high-risk"
	case normalized == "recovering":
		return "recovery"
	default:
		return "unknown"
	}
}

func (p *Processor) categorizeConditions(conditions []string) map[string][]string {
	categories := make(map[string][]string)

	conditionCategories := map[string]string{
		"diabetes":     "metabolic",
		"hypertension": "cardiovascular",
		"heart":        "cardiovascular",
		"asthma":       "respiratory",
		"copd":         "respiratory",
		"arthritis":    "musculoskeletal",
		"back":         "musculoskeletal",
		"depression":   "mental-health",
		"anxiety":      "mental-health",
		"bipolar":      "mental-health",
		"cancer":       "oncological",
		"tumor":        "oncological",
		"kidney":       "renal",
		"liver":        "hepatic",
		"thyroid":      "endocrine",
		"migraine":     "neurological",
		"epilepsy":     "neurological",
	}

	for _, condition := range conditions {
		normalized := p.normalizeString(condition)
		found := false
		for key, category := range conditionCategories {
			if strings.Contains(normalized, key) {
				categories[category] = append(categories[category], condition)
				found = true
				break
			}
		}
		if !found {
			categories["other"] = append(categories["other"], condition)
		}
	}

	return categories
}

func (p *Processor) hasChronicConditions(conditions []string) bool {
	chronicConditions := []string{
		"diabetes", "hypertension", "heart", "asthma", "copd", 
		"arthritis", "depression", "anxiety", "bipolar", "cancer",
		"kidney", "liver", "thyroid", "migraine", "epilepsy",
	}

	for _, condition := range conditions {
		normalized := p.normalizeString(condition)
		for _, chronic := range chronicConditions {
			if strings.Contains(normalized, chronic) {
				return true
			}
		}
	}
	return false
}

func (p *Processor) categorizeAllergies(allergies []string) map[string][]string {
	categories := make(map[string][]string)

	allergyCategories := map[string]string{
		"peanut":    "food",
		"nut":       "food",
		"dairy":     "food",
		"egg":       "food",
		"shellfish": "food",
		"soy":       "food",
		"wheat":     "food",
		"gluten":    "food",
		"pollen":    "environmental",
		"dust":      "environmental",
		"mold":      "environmental",
		"pet":       "environmental",
		"cat":       "environmental",
		"dog":       "environmental",
		"penicillin": "medication",
		"aspirin":   "medication",
		"latex":     "contact",
		"nickel":    "contact",
	}

	for _, allergy := range allergies {
		normalized := p.normalizeString(allergy)
		found := false
		for key, category := range allergyCategories {
			if strings.Contains(normalized, key) {
				categories[category] = append(categories[category], allergy)
				found = true
				break
			}
		}
		if !found {
			categories["other"] = append(categories["other"], allergy)
		}
	}

	return categories
}

func (p *Processor) getAllergySeverity(allergies []string) string {
	count := len(allergies)
	
	// Check for severe allergies
	severeAllergies := []string{"peanut", "shellfish", "penicillin"}
	for _, allergy := range allergies {
		normalized := p.normalizeString(allergy)
		for _, severe := range severeAllergies {
			if strings.Contains(normalized, severe) {
				return "severe"
			}
		}
	}

	switch {
	case count >= 5:
		return "multiple"
	case count >= 3:
		return "moderate"
	case count >= 1:
		return "mild"
	default:
		return "none"
	}
}
