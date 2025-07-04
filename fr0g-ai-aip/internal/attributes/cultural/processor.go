package cultural

import (
	"fmt"
	"strings"

	"github.com/fr0g-ai/fr0g-ai-aip/internal/config"
	pb "github.com/fr0g-ai/fr0g-ai-aip/internal/grpc/pb"
)

// Processor handles cultural and religious attribute processing and validation
type Processor struct {
	config *config.ValidationConfig
}

// NewProcessor creates a new cultural processor
func NewProcessor(cfg *config.ValidationConfig) *Processor {
	return &Processor{
		config: cfg,
	}
}

// ValidateCulturalReligious validates cultural and religious data
func (p *Processor) ValidateCulturalReligious(cultural *pb.CulturalReligious) []config.ValidationError {
	var errors []config.ValidationError

	// Validate spirituality level
	if cultural.Spirituality != "" {
		validLevels := []string{
			"very-low", "low", "moderate", "high", "very-high",
			"none", "spiritual-not-religious", "religious", "deeply-religious",
		}
		if !p.isValidOption(cultural.Spirituality, validLevels) {
			errors = append(errors, config.ValidationError{
				Field:   "spirituality",
				Message: "invalid spirituality level",
			})
		}
	}

	// Validate traditions list
	if len(cultural.Traditions) > 25 {
		errors = append(errors, config.ValidationError{
			Field:   "traditions",
			Message: "too many traditions specified (max 25)",
		})
	}

	// Validate holidays list
	if len(cultural.Holidays) > 30 {
		errors = append(errors, config.ValidationError{
			Field:   "holidays",
			Message: "too many holidays specified (max 30)",
		})
	}

	// Validate dietary restrictions
	if len(cultural.DietaryRestrictions) > 15 {
		errors = append(errors, config.ValidationError{
			Field:   "dietary_restrictions",
			Message: "too many dietary restrictions specified (max 15)",
		})
	}

	// Validate individual dietary restrictions
	for _, restriction := range cultural.DietaryRestrictions {
		if !p.isValidDietaryRestriction(restriction) {
			errors = append(errors, config.ValidationError{
				Field:   "dietary_restrictions",
				Message: fmt.Sprintf("invalid dietary restriction: %s", restriction),
			})
		}
	}

	return errors
}

// ProcessCulturalReligious processes and enriches cultural/religious data
func (p *Processor) ProcessCulturalReligious(cultural *pb.CulturalReligious) (*pb.CulturalReligious, error) {
	if cultural == nil {
		return nil, fmt.Errorf("cultural religious data cannot be nil")
	}

	// Validate first
	if validationErrors := p.ValidateCulturalReligious(cultural); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Create processed copy
	processed := &pb.CulturalReligious{
		Religion:            p.normalizeString(cultural.Religion),
		Spirituality:        p.normalizeString(cultural.Spirituality),
		CulturalBackground:  p.normalizeString(cultural.CulturalBackground),
		Traditions:          p.normalizeStringSlice(cultural.Traditions),
		Holidays:            p.normalizeStringSlice(cultural.Holidays),
		DietaryRestrictions: p.normalizeDietaryRestrictions(cultural.DietaryRestrictions),
	}

	return processed, nil
}

// GetReligiousProfile returns religious/spiritual profile information
func (p *Processor) GetReligiousProfile(cultural *pb.CulturalReligious) map[string]interface{} {
	profile := make(map[string]interface{})

	if cultural.Religion != "" {
		profile["religion"] = cultural.Religion
		profile["religious_family"] = p.getReligiousFamily(cultural.Religion)
	}

	if cultural.Spirituality != "" {
		profile["spirituality_level"] = cultural.Spirituality
		profile["spiritual_practices"] = p.getSpiritualPractices(cultural.Spirituality)
	}

	if len(cultural.Traditions) > 0 {
		profile["tradition_categories"] = p.categorizeTraditions(cultural.Traditions)
	}

	if len(cultural.Holidays) > 0 {
		profile["holiday_types"] = p.categorizeHolidays(cultural.Holidays)
	}

	return profile
}

// GetCulturalProfile returns cultural background information
func (p *Processor) GetCulturalProfile(cultural *pb.CulturalReligious) map[string]interface{} {
	profile := make(map[string]interface{})

	if cultural.CulturalBackground != "" {
		profile["cultural_background"] = cultural.CulturalBackground
		profile["cultural_region"] = p.getCulturalRegion(cultural.CulturalBackground)
	}

	if len(cultural.DietaryRestrictions) > 0 {
		profile["dietary_categories"] = p.categorizeDietaryRestrictions(cultural.DietaryRestrictions)
		profile["dietary_complexity"] = p.getDietaryComplexity(cultural.DietaryRestrictions)
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

func (p *Processor) isValidDietaryRestriction(restriction string) bool {
	validRestrictions := []string{
		"vegetarian", "vegan", "pescatarian", "halal", "kosher",
		"gluten-free", "dairy-free", "nut-free", "shellfish-free",
		"low-sodium", "low-sugar", "keto", "paleo", "raw",
		"organic-only", "no-pork", "no-beef", "no-alcohol",
	}
	return p.isValidOption(restriction, validRestrictions)
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

func (p *Processor) normalizeDietaryRestrictions(restrictions []string) []string {
	var normalized []string
	seen := make(map[string]bool)
	
	for _, restriction := range restrictions {
		norm := p.normalizeString(restriction)
		if norm != "" && !seen[norm] && p.isValidDietaryRestriction(norm) {
			normalized = append(normalized, norm)
			seen[norm] = true
		}
	}
	return normalized
}

func (p *Processor) getReligiousFamily(religion string) string {
	religionFamilies := map[string]string{
		"christianity": "abrahamic",
		"catholic":     "abrahamic",
		"protestant":   "abrahamic",
		"orthodox":     "abrahamic",
		"islam":        "abrahamic",
		"judaism":      "abrahamic",
		"hinduism":     "dharmic",
		"buddhism":     "dharmic",
		"sikhism":      "dharmic",
		"jainism":      "dharmic",
		"taoism":       "east-asian",
		"confucianism": "east-asian",
		"shinto":       "east-asian",
		"atheism":      "non-religious",
		"agnosticism":  "non-religious",
	}

	normalized := p.normalizeString(religion)
	if family, exists := religionFamilies[normalized]; exists {
		return family
	}
	return "other"
}

func (p *Processor) getSpiritualPractices(spirituality string) []string {
	practiceMap := map[string][]string{
		"very-high": {"meditation", "prayer", "ritual", "pilgrimage", "fasting"},
		"high":      {"meditation", "prayer", "ritual"},
		"moderate":  {"prayer", "reflection"},
		"low":       {"occasional-prayer"},
		"spiritual-not-religious": {"meditation", "mindfulness", "nature-connection"},
		"religious": {"prayer", "worship", "scripture-study"},
		"deeply-religious": {"prayer", "worship", "scripture-study", "ritual", "service"},
	}

	normalized := p.normalizeString(spirituality)
	if practices, exists := practiceMap[normalized]; exists {
		return practices
	}
	return []string{}
}

func (p *Processor) categorizeTraditions(traditions []string) map[string][]string {
	categories := make(map[string][]string)
	
	traditionCategories := map[string]string{
		"christmas":     "religious-holiday",
		"easter":        "religious-holiday",
		"ramadan":       "religious-holiday",
		"diwali":        "religious-holiday",
		"thanksgiving":  "cultural-holiday",
		"new-year":      "cultural-holiday",
		"wedding":       "life-event",
		"funeral":       "life-event",
		"birthday":      "life-event",
		"graduation":    "life-event",
		"family-dinner": "family",
		"storytelling":  "family",
		"cooking":       "cultural-practice",
		"music":         "cultural-practice",
		"dance":         "cultural-practice",
	}

	for _, tradition := range traditions {
		normalized := p.normalizeString(tradition)
		if category, exists := traditionCategories[normalized]; exists {
			categories[category] = append(categories[category], tradition)
		} else {
			categories["other"] = append(categories["other"], tradition)
		}
	}

	return categories
}

func (p *Processor) categorizeHolidays(holidays []string) map[string][]string {
	categories := make(map[string][]string)
	
	holidayCategories := map[string]string{
		"christmas":      "religious",
		"easter":         "religious",
		"eid":            "religious",
		"diwali":         "religious",
		"hanukkah":       "religious",
		"new-year":       "secular",
		"independence":   "national",
		"thanksgiving":   "national",
		"memorial-day":   "national",
		"labor-day":      "national",
		"valentine":      "cultural",
		"halloween":      "cultural",
		"mother-day":     "family",
		"father-day":     "family",
	}

	for _, holiday := range holidays {
		normalized := p.normalizeString(holiday)
		found := false
		for key, category := range holidayCategories {
			if strings.Contains(normalized, key) {
				categories[category] = append(categories[category], holiday)
				found = true
				break
			}
		}
		if !found {
			categories["other"] = append(categories["other"], holiday)
		}
	}

	return categories
}

func (p *Processor) getCulturalRegion(background string) string {
	regionMap := map[string]string{
		"american":     "north-america",
		"canadian":     "north-america",
		"mexican":      "north-america",
		"british":      "europe",
		"french":       "europe",
		"german":       "europe",
		"italian":      "europe",
		"spanish":      "europe",
		"chinese":      "east-asia",
		"japanese":     "east-asia",
		"korean":       "east-asia",
		"indian":       "south-asia",
		"pakistani":    "south-asia",
		"bangladeshi":  "south-asia",
		"arab":         "middle-east",
		"persian":      "middle-east",
		"turkish":      "middle-east",
		"african":      "africa",
		"nigerian":     "africa",
		"south-african": "africa",
		"brazilian":    "south-america",
		"argentinian":  "south-america",
		"australian":   "oceania",
	}

	normalized := p.normalizeString(background)
	for key, region := range regionMap {
		if strings.Contains(normalized, key) {
			return region
		}
	}
	return "other"
}

func (p *Processor) categorizeDietaryRestrictions(restrictions []string) map[string][]string {
	categories := make(map[string][]string)
	
	restrictionCategories := map[string]string{
		"vegetarian":   "ethical",
		"vegan":        "ethical",
		"halal":        "religious",
		"kosher":       "religious",
		"gluten-free":  "medical",
		"dairy-free":   "medical",
		"nut-free":     "medical",
		"shellfish-free": "medical",
		"low-sodium":   "health",
		"low-sugar":    "health",
		"keto":         "lifestyle",
		"paleo":        "lifestyle",
		"raw":          "lifestyle",
		"organic-only": "lifestyle",
	}

	for _, restriction := range restrictions {
		normalized := p.normalizeString(restriction)
		if category, exists := restrictionCategories[normalized]; exists {
			categories[category] = append(categories[category], restriction)
		} else {
			categories["other"] = append(categories["other"], restriction)
		}
	}

	return categories
}

func (p *Processor) getDietaryComplexity(restrictions []string) string {
	count := len(restrictions)
	switch {
	case count == 0:
		return "none"
	case count <= 2:
		return "simple"
	case count <= 5:
		return "moderate"
	default:
		return "complex"
	}
}
