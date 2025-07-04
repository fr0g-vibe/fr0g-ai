package demographics

import (
	"fmt"
	"strings"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/config"
	pb "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/grpc/pb"
)

// Processor handles demographics attribute processing and validation
type Processor struct {
	config *config.ValidationConfig
}

// NewProcessor creates a new demographics processor
func NewProcessor(cfg *config.ValidationConfig) *Processor {
	return &Processor{
		config: cfg,
	}
}

// ValidateDemographics validates demographics data
func (p *Processor) ValidateDemographics(demo *pb.Demographics) []config.ValidationError {
	var errors []config.ValidationError

	// Validate age
	if demo.Age < 0 || demo.Age > 150 {
		errors = append(errors, config.ValidationError{
			Field:   "age",
			Message: "age must be between 0 and 150",
		})
	}

	// Validate gender (allow common values or custom)
	if demo.Gender != "" {
		validGenders := []string{"male", "female", "non-binary", "other", "prefer-not-to-say"}
		if !p.isValidOption(demo.Gender, validGenders) && !p.isCustomValue(demo.Gender) {
			errors = append(errors, config.ValidationError{
				Field:   "gender",
				Message: "invalid gender value",
			})
		}
	}

	// Validate education level
	if demo.Education != "" {
		validEducation := []string{
			"elementary", "high-school", "some-college", "bachelor", 
			"master", "doctorate", "professional", "other",
		}
		if !p.isValidOption(demo.Education, validEducation) {
			errors = append(errors, config.ValidationError{
				Field:   "education",
				Message: "invalid education level",
			})
		}
	}

	// Validate socioeconomic status
	if demo.SocioeconomicStatus != "" {
		validStatus := []string{"low", "lower-middle", "middle", "upper-middle", "high"}
		if !p.isValidOption(demo.SocioeconomicStatus, validStatus) {
			errors = append(errors, config.ValidationError{
				Field:   "socioeconomic_status",
				Message: "invalid socioeconomic status",
			})
		}
	}

	// Validate location if present
	if demo.Location != nil {
		locationErrors := p.validateLocation(demo.Location)
		errors = append(errors, locationErrors...)
	}

	return errors
}

// ProcessDemographics processes and enriches demographics data
func (p *Processor) ProcessDemographics(demo *pb.Demographics) (*pb.Demographics, error) {
	if demo == nil {
		return nil, fmt.Errorf("demographics cannot be nil")
	}

	// Validate first
	if validationErrors := p.ValidateDemographics(demo); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Create processed copy
	processed := &pb.Demographics{
		Age:                 demo.Age,
		Gender:              p.normalizeString(demo.Gender),
		Ethnicity:           p.normalizeString(demo.Ethnicity),
		Nationality:         p.normalizeString(demo.Nationality),
		Education:           p.normalizeString(demo.Education),
		Occupation:          p.normalizeString(demo.Occupation),
		SocioeconomicStatus: p.normalizeString(demo.SocioeconomicStatus),
		Location:            demo.Location,
	}

	// Process location if present
	if demo.Location != nil {
		processedLocation, err := p.processLocation(demo.Location)
		if err != nil {
			return nil, fmt.Errorf("failed to process location: %w", err)
		}
		processed.Location = processedLocation
	}

	return processed, nil
}

// validateLocation validates location data
func (p *Processor) validateLocation(loc *pb.Location) []config.ValidationError {
	var errors []config.ValidationError

	// Validate urban/rural classification
	if loc.UrbanRural != "" {
		validTypes := []string{"urban", "suburban", "rural"}
		if !p.isValidOption(loc.UrbanRural, validTypes) {
			errors = append(errors, config.ValidationError{
				Field:   "urban_rural",
				Message: "must be 'urban', 'suburban', or 'rural'",
			})
		}
	}

	// Basic timezone validation (simplified)
	if loc.Timezone != "" && !p.isValidTimezone(loc.Timezone) {
		errors = append(errors, config.ValidationError{
			Field:   "timezone",
			Message: "invalid timezone format",
		})
	}

	return errors
}

// processLocation processes location data
func (p *Processor) processLocation(loc *pb.Location) (*pb.Location, error) {
	return &pb.Location{
		Country:    p.normalizeString(loc.Country),
		Region:     p.normalizeString(loc.Region),
		City:       p.normalizeString(loc.City),
		UrbanRural: p.normalizeString(loc.UrbanRural),
		Timezone:   strings.TrimSpace(loc.Timezone),
	}, nil
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

func (p *Processor) isCustomValue(value string) bool {
	// Allow custom values that are reasonable length and don't contain special chars
	normalized := strings.TrimSpace(strings.ToLower(value))
	return len(normalized) > 0 && len(normalized) <= 50 && !strings.ContainsAny(normalized, "<>{}[]")
}

func (p *Processor) normalizeString(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func (p *Processor) isValidTimezone(tz string) bool {
	// Simplified timezone validation - in production, use proper timezone library
	return strings.Contains(tz, "/") || strings.HasPrefix(tz, "UTC") || strings.HasPrefix(tz, "GMT")
}

// GetAgeGroup returns age group classification
func (p *Processor) GetAgeGroup(age int32) string {
	switch {
	case age < 13:
		return "child"
	case age < 20:
		return "teenager"
	case age < 30:
		return "young-adult"
	case age < 50:
		return "adult"
	case age < 65:
		return "middle-aged"
	default:
		return "senior"
	}
}

// GetGenerationCohort returns generation classification
func (p *Processor) GetGenerationCohort(age int32) string {
	// Approximate birth year based on current age (2024)
	birthYear := 2024 - int(age)
	
	switch {
	case birthYear >= 2010:
		return "gen-alpha"
	case birthYear >= 1997:
		return "gen-z"
	case birthYear >= 1981:
		return "millennial"
	case birthYear >= 1965:
		return "gen-x"
	case birthYear >= 1946:
		return "boomer"
	default:
		return "silent-generation"
	}
}
