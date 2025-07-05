package lifehistory

import (
	"fmt"
	"strings"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/config"
	pb "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/grpc/pb"
)

// Processor handles life history attribute processing and validation
type Processor struct {
	config *config.ValidationConfig
}

// NewProcessor creates a new life history processor
func NewProcessor(cfg *config.ValidationConfig) *Processor {
	return &Processor{
		config: cfg,
	}
}

// ValidateLifeHistory validates life history data
func (p *Processor) ValidateLifeHistory(history *pb.LifeHistory) []config.ValidationError {
	var errors []config.ValidationError

	// Validate childhood traumas list
	if len(history.ChildhoodTraumas) > 15 {
		errors = append(errors, config.ValidationError{
			Field:   "childhood_traumas",
			Message: "too many childhood traumas specified (max 15)",
		})
	}

	// Validate adult traumas list
	if len(history.AdultTraumas) > 15 {
		errors = append(errors, config.ValidationError{
			Field:   "adult_traumas",
			Message: "too many adult traumas specified (max 15)",
		})
	}

	// Validate major events list
	if len(history.MajorEvents) > 25 {
		errors = append(errors, config.ValidationError{
			Field:   "major_events",
			Message: "too many major events specified (max 25)",
		})
	}

	// Validate education history list
	if len(history.EducationHistory) > 10 {
		errors = append(errors, config.ValidationError{
			Field:   "education_history",
			Message: "too many education entries specified (max 10)",
		})
	}

	// Validate career history list
	if len(history.CareerHistory) > 15 {
		errors = append(errors, config.ValidationError{
			Field:   "career_history",
			Message: "too many career entries specified (max 15)",
		})
	}

	// Validate individual major events
	for i, event := range history.MajorEvents {
		if eventErrors := p.validateLifeEvent(event, fmt.Sprintf("major_events[%d]", i)); len(eventErrors) > 0 {
			errors = append(errors, eventErrors...)
		}
	}

	// Validate individual education entries
	for i, education := range history.EducationHistory {
		if eduErrors := p.validateEducation(education, fmt.Sprintf("education_history[%d]", i)); len(eduErrors) > 0 {
			errors = append(errors, eduErrors...)
		}
	}

	// Validate individual career entries
	for i, career := range history.CareerHistory {
		if careerErrors := p.validateCareer(career, fmt.Sprintf("career_history[%d]", i)); len(careerErrors) > 0 {
			errors = append(errors, careerErrors...)
		}
	}

	return errors
}

// validateLifeEvent validates a single life event
func (p *Processor) validateLifeEvent(event *pb.LifeEvent, fieldPrefix string) []config.ValidationError {
	var errors []config.ValidationError

	if event == nil {
		errors = append(errors, config.ValidationError{
			Field:   fieldPrefix,
			Message: "life event cannot be nil",
		})
		return errors
	}

	// Validate event type
	if event.Type != "" {
		validTypes := []string{
			"birth", "death", "marriage", "divorce", "graduation", "job-change",
			"relocation", "illness", "accident", "achievement", "loss", "trauma",
			"milestone", "relationship", "family", "career", "education", "health",
		}
		if !p.isValidOption(event.Type, validTypes) {
			errors = append(errors, config.ValidationError{
				Field:   fieldPrefix + ".type",
				Message: "invalid event type",
			})
		}
	}

	// Validate impact level
	if event.Impact != "" {
		validImpacts := []string{"very-low", "low", "moderate", "high", "very-high", "life-changing"}
		if !p.isValidOption(event.Impact, validImpacts) {
			errors = append(errors, config.ValidationError{
				Field:   fieldPrefix + ".impact",
				Message: "invalid impact level",
			})
		}
	}

	// Validate description length
	if len(event.Description) > 500 {
		errors = append(errors, config.ValidationError{
			Field:   fieldPrefix + ".description",
			Message: "description too long (max 500 characters)",
		})
	}

	return errors
}

// validateEducation validates a single education entry
func (p *Processor) validateEducation(education *pb.Education, fieldPrefix string) []config.ValidationError {
	var errors []config.ValidationError

	if education == nil {
		errors = append(errors, config.ValidationError{
			Field:   fieldPrefix,
			Message: "education entry cannot be nil",
		})
		return errors
	}

	// Validate education level
	if education.Level != "" {
		validLevels := []string{
			"elementary", "middle-school", "high-school", "ged", "some-college",
			"associate", "bachelor", "master", "doctorate", "professional",
			"certificate", "diploma", "trade-school", "bootcamp", "online-course",
		}
		if !p.isValidOption(education.Level, validLevels) {
			errors = append(errors, config.ValidationError{
				Field:   fieldPrefix + ".level",
				Message: "invalid education level",
			})
		}
	}

	// Validate field of study length
	if len(education.Field) > 100 {
		errors = append(errors, config.ValidationError{
			Field:   fieldPrefix + ".field",
			Message: "field too long (max 100 characters)",
		})
	}

	// Validate institution length
	if len(education.Institution) > 150 {
		errors = append(errors, config.ValidationError{
			Field:   fieldPrefix + ".institution",
			Message: "institution name too long (max 150 characters)",
		})
	}

	return errors
}

// validateCareer validates a single career entry
func (p *Processor) validateCareer(career *pb.Career, fieldPrefix string) []config.ValidationError {
	var errors []config.ValidationError

	if career == nil {
		errors = append(errors, config.ValidationError{
			Field:   fieldPrefix,
			Message: "career entry cannot be nil",
		})
		return errors
	}

	// Validate job title length
	if len(career.Title) > 100 {
		errors = append(errors, config.ValidationError{
			Field:   fieldPrefix + ".title",
			Message: "title too long (max 100 characters)",
		})
	}

	// Validate company length
	if len(career.Company) > 150 {
		errors = append(errors, config.ValidationError{
			Field:   fieldPrefix + ".company",
			Message: "company name too long (max 150 characters)",
		})
	}

	// Validate industry
	if career.Industry != "" {
		validIndustries := []string{
			"technology", "healthcare", "finance", "education", "retail", "manufacturing",
			"construction", "transportation", "hospitality", "entertainment", "media",
			"government", "nonprofit", "consulting", "legal", "real-estate", "agriculture",
			"energy", "telecommunications", "automotive", "aerospace", "pharmaceutical",
			"biotechnology", "food-service", "logistics", "insurance", "banking", "other",
		}
		if !p.isValidOption(career.Industry, validIndustries) {
			errors = append(errors, config.ValidationError{
				Field:   fieldPrefix + ".industry",
				Message: "invalid industry",
			})
		}
	}

	// Note: Career protobuf doesn't have description field
	// Only validate existing fields: title, industry, company, dates, salary

	return errors
}

// ProcessLifeHistory processes and enriches life history data
func (p *Processor) ProcessLifeHistory(history *pb.LifeHistory) (*pb.LifeHistory, error) {
	if history == nil {
		return nil, fmt.Errorf("life history cannot be nil")
	}

	// Validate first
	if validationErrors := p.ValidateLifeHistory(history); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Create processed copy
	processed := &pb.LifeHistory{
		ChildhoodTraumas:  p.normalizeStringSlice(history.ChildhoodTraumas),
		AdultTraumas:      p.normalizeStringSlice(history.AdultTraumas),
		MajorEvents:       p.processMajorEvents(history.MajorEvents),
		EducationHistory:  p.processEducationHistory(history.EducationHistory),
		CareerHistory:     p.processCareerHistory(history.CareerHistory),
	}

	return processed, nil
}

// processMajorEvents processes major life events
func (p *Processor) processMajorEvents(events []*pb.LifeEvent) []*pb.LifeEvent {
	var processed []*pb.LifeEvent
	for _, event := range events {
		if event != nil {
			processedEvent := &pb.LifeEvent{
				Type:        p.normalizeString(event.Type),
				Description: strings.TrimSpace(event.Description),
				Age:         event.Age,
				Date:        event.Date,
				Impact:      p.normalizeString(event.Impact),
			}
			processed = append(processed, processedEvent)
		}
	}
	return processed
}

// processEducationHistory processes education history
func (p *Processor) processEducationHistory(education []*pb.Education) []*pb.Education {
	var processed []*pb.Education
	for _, edu := range education {
		if edu != nil {
			processedEdu := &pb.Education{
				Level:       p.normalizeString(edu.Level),
				Field:       strings.TrimSpace(edu.Field),
				Institution: strings.TrimSpace(edu.Institution),
				Graduation:  edu.Graduation,
				Performance: p.normalizeString(edu.Performance),
			}
			processed = append(processed, processedEdu)
		}
	}
	return processed
}

// processCareerHistory processes career history
func (p *Processor) processCareerHistory(careers []*pb.Career) []*pb.Career {
	var processed []*pb.Career
	for _, career := range careers {
		if career != nil {
			processedCareer := &pb.Career{
				Title:     strings.TrimSpace(career.Title),
				Industry:  strings.TrimSpace(career.Industry),
				Company:   strings.TrimSpace(career.Company),
				StartDate: career.StartDate,
				EndDate:   career.EndDate,
				IsCurrent: career.IsCurrent,
				Salary:    strings.TrimSpace(career.Salary),
			}
			processed = append(processed, processedCareer)
		}
	}
	return processed
}

// GetLifeHistoryProfile returns comprehensive life history profile
func (p *Processor) GetLifeHistoryProfile(history *pb.LifeHistory) map[string]interface{} {
	profile := make(map[string]interface{})

	// Trauma analysis
	if len(history.ChildhoodTraumas) > 0 || len(history.AdultTraumas) > 0 {
		profile["trauma_profile"] = p.getTraumaProfile(history)
	}

	// Major events analysis
	if len(history.MajorEvents) > 0 {
		profile["life_events"] = p.getLifeEventsProfile(history.MajorEvents)
	}

	// Education analysis
	if len(history.EducationHistory) > 0 {
		profile["education_profile"] = p.getEducationProfile(history.EducationHistory)
	}

	// Career analysis
	if len(history.CareerHistory) > 0 {
		profile["career_profile"] = p.getCareerProfile(history.CareerHistory)
	}

	// Overall life patterns
	profile["life_patterns"] = p.getLifePatterns(history)
	profile["resilience_indicators"] = p.getResilienceIndicators(history)

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

func (p *Processor) getTraumaProfile(history *pb.LifeHistory) map[string]interface{} {
	profile := make(map[string]interface{})

	// Categorize traumas
	if len(history.ChildhoodTraumas) > 0 {
		profile["childhood_trauma_categories"] = p.categorizeTraumas(history.ChildhoodTraumas)
		profile["childhood_trauma_severity"] = p.getTraumaSeverity(history.ChildhoodTraumas)
	}

	if len(history.AdultTraumas) > 0 {
		profile["adult_trauma_categories"] = p.categorizeTraumas(history.AdultTraumas)
		profile["adult_trauma_severity"] = p.getTraumaSeverity(history.AdultTraumas)
	}

	// Overall trauma assessment
	totalTraumas := len(history.ChildhoodTraumas) + len(history.AdultTraumas)
	profile["total_trauma_count"] = totalTraumas
	profile["trauma_burden"] = p.getTraumaBurden(totalTraumas)

	return profile
}

func (p *Processor) categorizeTraumas(traumas []string) map[string][]string {
	categories := make(map[string][]string)

	traumaCategories := map[string]string{
		"abuse":        "interpersonal",
		"neglect":      "interpersonal",
		"violence":     "interpersonal",
		"bullying":     "interpersonal",
		"death":        "loss",
		"loss":         "loss",
		"separation":   "loss",
		"divorce":      "loss",
		"accident":     "physical",
		"injury":       "physical",
		"illness":      "physical",
		"surgery":      "physical",
		"poverty":      "environmental",
		"homelessness": "environmental",
		"war":          "environmental",
		"disaster":     "environmental",
		"addiction":    "substance",
		"alcoholism":   "substance",
	}

	for _, trauma := range traumas {
		normalized := p.normalizeString(trauma)
		found := false
		for key, category := range traumaCategories {
			if strings.Contains(normalized, key) {
				categories[category] = append(categories[category], trauma)
				found = true
				break
			}
		}
		if !found {
			categories["other"] = append(categories["other"], trauma)
		}
	}

	return categories
}

func (p *Processor) getTraumaSeverity(traumas []string) string {
	count := len(traumas)
	switch {
	case count == 0:
		return "none"
	case count <= 2:
		return "mild"
	case count <= 4:
		return "moderate"
	case count <= 7:
		return "severe"
	default:
		return "extreme"
	}
}

func (p *Processor) getTraumaBurden(totalCount int) string {
	switch {
	case totalCount == 0:
		return "none"
	case totalCount <= 2:
		return "low"
	case totalCount <= 5:
		return "moderate"
	case totalCount <= 8:
		return "high"
	default:
		return "very-high"
	}
}

func (p *Processor) getLifeEventsProfile(events []*pb.LifeEvent) map[string]interface{} {
	profile := make(map[string]interface{})

	// Categorize events by type
	eventCategories := p.categorizeLifeEvents(events)
	profile["event_categories"] = eventCategories

	// Analyze impact distribution
	profile["impact_distribution"] = p.getImpactDistribution(events)

	// Timeline analysis
	profile["event_timeline"] = p.getEventTimeline(events)

	// Life phases analysis
	profile["life_phases"] = p.getLifePhases(events)

	return profile
}

func (p *Processor) categorizeLifeEvents(events []*pb.LifeEvent) map[string][]string {
	categories := make(map[string][]string)

	for _, event := range events {
		if event != nil && event.Type != "" {
			eventType := p.normalizeString(event.Type)
			categories[eventType] = append(categories[eventType], event.Description)
		}
	}

	return categories
}

func (p *Processor) getImpactDistribution(events []*pb.LifeEvent) map[string]int {
	distribution := make(map[string]int)

	for _, event := range events {
		if event != nil && event.Impact != "" {
			impact := p.normalizeString(event.Impact)
			distribution[impact]++
		}
	}

	return distribution
}

func (p *Processor) getEventTimeline(events []*pb.LifeEvent) map[string]interface{} {
	timeline := make(map[string]interface{})

	// Group events by decade
	decades := make(map[string][]string)
	for _, event := range events {
		if event != nil && event.Age > 0 {
			decade := fmt.Sprintf("%ds", (event.Age/10)*10)
			decades[decade] = append(decades[decade], event.Type)
		}
	}

	timeline["by_decade"] = decades
	timeline["total_events"] = len(events)

	return timeline
}

func (p *Processor) getLifePhases(events []*pb.LifeEvent) map[string][]string {
	phases := make(map[string][]string)

	for _, event := range events {
		if event != nil && event.Age > 0 {
			var phase string
			switch {
			case event.Age <= 12:
				phase = "childhood"
			case event.Age <= 17:
				phase = "adolescence"
			case event.Age <= 25:
				phase = "young-adult"
			case event.Age <= 40:
				phase = "adult"
			case event.Age <= 65:
				phase = "middle-age"
			default:
				phase = "senior"
			}
			phases[phase] = append(phases[phase], event.Type)
		}
	}

	return phases
}

func (p *Processor) getEducationProfile(education []*pb.Education) map[string]interface{} {
	profile := make(map[string]interface{})

	// Education levels achieved
	levels := make(map[string]int)
	fields := make(map[string]int)
	var totalYears int32
	var completedCount int

	for _, edu := range education {
		if edu != nil {
			if edu.Level != "" {
				level := p.normalizeString(edu.Level)
				levels[level]++
			}
			if edu.Field != "" {
				field := p.normalizeString(edu.Field)
				fields[field]++
			}
			// Note: StartYear/EndYear/Completed are not in protobuf Education message
			// Using available fields only
			completedCount++ // Assume all listed education is completed
		}
	}

	profile["education_levels"] = levels
	profile["fields_of_study"] = fields
	profile["total_education_years"] = totalYears
	profile["completion_rate"] = float64(completedCount) / float64(len(education))
	profile["highest_level"] = p.getHighestEducationLevel(education)
	profile["education_diversity"] = len(fields)

	return profile
}

func (p *Processor) getHighestEducationLevel(education []*pb.Education) string {
	levelHierarchy := map[string]int{
		"elementary":    1,
		"middle-school": 2,
		"high-school":   3,
		"ged":           3,
		"certificate":   4,
		"diploma":       4,
		"trade-school":  4,
		"some-college":  5,
		"associate":     6,
		"bachelor":      7,
		"master":        8,
		"professional":  9,
		"doctorate":     10,
	}

	highestLevel := ""
	highestRank := 0

	for _, edu := range education {
		if edu != nil && edu.Level != "" {
			level := p.normalizeString(edu.Level)
			if rank, exists := levelHierarchy[level]; exists && rank > highestRank {
				highestRank = rank
				highestLevel = level
			}
		}
	}

	return highestLevel
}

func (p *Processor) getCareerProfile(careers []*pb.Career) map[string]interface{} {
	profile := make(map[string]interface{})

	// Industry analysis
	industries := make(map[string]int)
	var totalYears int32
	var jobChanges int

	for _, career := range careers {
		if career != nil {
			if career.Industry != "" {
				industry := p.normalizeString(career.Industry)
				industries[industry]++
			}
			// Note: StartYear/EndYear are not in protobuf Career message
			// Using available fields only
			if !career.IsCurrent {
				jobChanges++
			}
		}
	}

	profile["industries"] = industries
	profile["total_career_years"] = totalYears
	profile["job_changes"] = jobChanges
	profile["career_stability"] = p.getCareerStability(jobChanges, len(careers))
	profile["industry_diversity"] = len(industries)
	profile["current_position"] = p.getCurrentPosition(careers)

	return profile
}

func (p *Processor) getCareerStability(jobChanges, totalJobs int) string {
	if totalJobs == 0 {
		return "unknown"
	}

	changeRate := float64(jobChanges) / float64(totalJobs)
	switch {
	case changeRate <= 0.2:
		return "very-stable"
	case changeRate <= 0.4:
		return "stable"
	case changeRate <= 0.6:
		return "moderate"
	case changeRate <= 0.8:
		return "unstable"
	default:
		return "very-unstable"
	}
}

func (p *Processor) getCurrentPosition(careers []*pb.Career) map[string]string {
	for _, career := range careers {
		if career != nil && career.IsCurrent {
			return map[string]string{
				"title":    career.Title,
				"company":  career.Company,
				"industry": career.Industry,
			}
		}
	}
	return nil
}

func (p *Processor) getLifePatterns(history *pb.LifeHistory) map[string]interface{} {
	patterns := make(map[string]interface{})

	// Resilience patterns
	patterns["adaptability"] = p.getAdaptabilityScore(history)
	patterns["growth_orientation"] = p.getGrowthOrientation(history)
	patterns["stability_seeking"] = p.getStabilityPattern(history)

	// Life complexity
	totalEvents := len(history.MajorEvents) + len(history.ChildhoodTraumas) + len(history.AdultTraumas)
	patterns["life_complexity"] = p.getLifeComplexity(totalEvents)

	return patterns
}

func (p *Processor) getAdaptabilityScore(history *pb.LifeHistory) string {
	score := 0

	// Education diversity indicates adaptability
	if len(history.EducationHistory) > 2 {
		score += 2
	}

	// Career changes can indicate adaptability
	if len(history.CareerHistory) > 1 {
		score += 1
	}

	// Overcoming traumas indicates resilience
	traumaCount := len(history.ChildhoodTraumas) + len(history.AdultTraumas)
	if traumaCount > 0 && len(history.MajorEvents) > traumaCount {
		score += 2
	}

	switch {
	case score >= 4:
		return "very-high"
	case score >= 3:
		return "high"
	case score >= 2:
		return "moderate"
	case score >= 1:
		return "low"
	default:
		return "very-low"
	}
}

func (p *Processor) getGrowthOrientation(history *pb.LifeHistory) string {
	// Count positive/growth events vs negative events
	positiveEvents := 0
	for _, event := range history.MajorEvents {
		if event != nil {
			eventType := p.normalizeString(event.EventType)
			if eventType == "achievement" || eventType == "graduation" || eventType == "milestone" {
				positiveEvents++
			}
		}
	}

	totalEvents := len(history.MajorEvents)
	if totalEvents == 0 {
		return "unknown"
	}

	ratio := float64(positiveEvents) / float64(totalEvents)
	switch {
	case ratio >= 0.6:
		return "high"
	case ratio >= 0.4:
		return "moderate"
	default:
		return "low"
	}
}

func (p *Processor) getStabilityPattern(history *pb.LifeHistory) string {
	// Analyze job stability and education completion
	stabilityScore := 0

	// Education completion rate
	if len(history.EducationHistory) > 0 {
		completed := 0
		for _, edu := range history.EducationHistory {
			if edu != nil && edu.Completed {
				completed++
			}
		}
		if float64(completed)/float64(len(history.EducationHistory)) >= 0.8 {
			stabilityScore += 2
		}
	}

	// Career stability
	if len(history.CareerHistory) > 0 {
		longTermJobs := 0
		for _, career := range history.CareerHistory {
			if career != nil && career.EndYear-career.StartYear >= 3 {
				longTermJobs++
			}
		}
		if float64(longTermJobs)/float64(len(history.CareerHistory)) >= 0.5 {
			stabilityScore += 2
		}
	}

	switch {
	case stabilityScore >= 3:
		return "high"
	case stabilityScore >= 2:
		return "moderate"
	default:
		return "low"
	}
}

func (p *Processor) getLifeComplexity(totalEvents int) string {
	switch {
	case totalEvents >= 20:
		return "very-complex"
	case totalEvents >= 15:
		return "complex"
	case totalEvents >= 10:
		return "moderate"
	case totalEvents >= 5:
		return "simple"
	default:
		return "very-simple"
	}
}

func (p *Processor) getResilienceIndicators(history *pb.LifeHistory) []string {
	var indicators []string

	// Education after trauma
	traumaCount := len(history.ChildhoodTraumas) + len(history.AdultTraumas)
	if traumaCount > 0 && len(history.EducationHistory) > 0 {
		indicators = append(indicators, "continued-education-despite-trauma")
	}

	// Career progression
	if len(history.CareerHistory) > 1 {
		indicators = append(indicators, "career-progression")
	}

	// Positive major events after trauma
	if traumaCount > 0 {
		for _, event := range history.MajorEvents {
			if event != nil {
				eventType := p.normalizeString(event.EventType)
				if eventType == "achievement" || eventType == "graduation" {
					indicators = append(indicators, "achievements-after-trauma")
					break
				}
			}
		}
	}

	// Long-term relationships/stability
	for _, event := range history.MajorEvents {
		if event != nil && p.normalizeString(event.EventType) == "marriage" {
			indicators = append(indicators, "long-term-relationships")
			break
		}
	}

	return indicators
}
