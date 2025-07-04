package political

import (
	"fmt"
	"strings"

	"github.com/fr0g-ai/fr0g-ai-aip/internal/config"
	pb "github.com/fr0g-ai/fr0g-ai-aip/internal/grpc/pb"
)

// Processor handles political and social attribute processing and validation
type Processor struct {
	config *config.ValidationConfig
}

// NewProcessor creates a new political processor
func NewProcessor(cfg *config.ValidationConfig) *Processor {
	return &Processor{
		config: cfg,
	}
}

// ValidatePoliticalSocial validates political and social data
func (p *Processor) ValidatePoliticalSocial(political *pb.PoliticalSocial) []config.ValidationError {
	var errors []config.ValidationError

	// Validate political leaning
	if political.PoliticalLeaning != "" {
		validLeanings := []string{
			"far-left", "left", "center-left", "center", "center-right", 
			"right", "far-right", "libertarian", "authoritarian", 
			"progressive", "conservative", "moderate", "independent", "apolitical",
		}
		if !p.isValidOption(political.PoliticalLeaning, validLeanings) {
			errors = append(errors, config.ValidationError{
				Field:   "political_leaning",
				Message: "invalid political leaning",
			})
		}
	}

	// Validate activism list
	if len(political.Activism) > 20 {
		errors = append(errors, config.ValidationError{
			Field:   "activism",
			Message: "too many activism entries specified (max 20)",
		})
	}

	// Validate social groups list
	if len(political.SocialGroups) > 25 {
		errors = append(errors, config.ValidationError{
			Field:   "social_groups",
			Message: "too many social groups specified (max 25)",
		})
	}

	// Validate causes list
	if len(political.Causes) > 30 {
		errors = append(errors, config.ValidationError{
			Field:   "causes",
			Message: "too many causes specified (max 30)",
		})
	}

	// Validate voting history
	if political.VotingHistory != "" {
		validHistory := []string{
			"never-voted", "rarely-votes", "sometimes-votes", 
			"usually-votes", "always-votes", "eligible-but-never-voted",
			"not-eligible", "first-time-voter",
		}
		if !p.isValidOption(political.VotingHistory, validHistory) {
			errors = append(errors, config.ValidationError{
				Field:   "voting_history",
				Message: "invalid voting history",
			})
		}
	}

	// Validate media consumption list
	if len(political.MediaConsumption) > 15 {
		errors = append(errors, config.ValidationError{
			Field:   "media_consumption",
			Message: "too many media consumption entries specified (max 15)",
		})
	}

	return errors
}

// ProcessPoliticalSocial processes and enriches political/social data
func (p *Processor) ProcessPoliticalSocial(political *pb.PoliticalSocial) (*pb.PoliticalSocial, error) {
	if political == nil {
		return nil, fmt.Errorf("political social data cannot be nil")
	}

	// Validate first
	if validationErrors := p.ValidatePoliticalSocial(political); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Create processed copy
	processed := &pb.PoliticalSocial{
		PoliticalLeaning: p.normalizeString(political.PoliticalLeaning),
		Activism:         p.normalizeStringSlice(political.Activism),
		SocialGroups:     p.normalizeStringSlice(political.SocialGroups),
		Causes:           p.normalizeStringSlice(political.Causes),
		VotingHistory:    p.normalizeString(political.VotingHistory),
		MediaConsumption: p.normalizeStringSlice(political.MediaConsumption),
	}

	return processed, nil
}

// GetPoliticalProfile returns political profile information
func (p *Processor) GetPoliticalProfile(political *pb.PoliticalSocial) map[string]interface{} {
	profile := make(map[string]interface{})

	if political.PoliticalLeaning != "" {
		profile["political_leaning"] = political.PoliticalLeaning
		profile["political_spectrum"] = p.getPoliticalSpectrum(political.PoliticalLeaning)
		profile["political_engagement"] = p.getPoliticalEngagement(political)
	}

	if political.VotingHistory != "" {
		profile["voting_history"] = political.VotingHistory
		profile["voting_likelihood"] = p.getVotingLikelihood(political.VotingHistory)
	}

	if len(political.Activism) > 0 {
		profile["activism_categories"] = p.categorizeActivism(political.Activism)
		profile["activism_level"] = p.getActivismLevel(political.Activism)
	}

	if len(political.Causes) > 0 {
		profile["cause_categories"] = p.categorizeCauses(political.Causes)
	}

	return profile
}

// GetSocialProfile returns social engagement information
func (p *Processor) GetSocialProfile(political *pb.PoliticalSocial) map[string]interface{} {
	profile := make(map[string]interface{})

	if len(political.SocialGroups) > 0 {
		profile["social_group_categories"] = p.categorizeSocialGroups(political.SocialGroups)
		profile["social_engagement"] = p.getSocialEngagement(political.SocialGroups)
	}

	if len(political.MediaConsumption) > 0 {
		profile["media_categories"] = p.categorizeMedia(political.MediaConsumption)
		profile["media_diversity"] = p.getMediaDiversity(political.MediaConsumption)
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

func (p *Processor) getPoliticalSpectrum(leaning string) string {
	spectrumMap := map[string]string{
		"far-left":      "left",
		"left":          "left",
		"center-left":   "center-left",
		"center":        "center",
		"center-right":  "center-right",
		"right":         "right",
		"far-right":     "right",
		"libertarian":   "libertarian",
		"authoritarian": "authoritarian",
		"progressive":   "left",
		"conservative":  "right",
		"moderate":      "center",
		"independent":   "independent",
		"apolitical":    "apolitical",
	}

	normalized := p.normalizeString(leaning)
	if spectrum, exists := spectrumMap[normalized]; exists {
		return spectrum
	}
	return "unknown"
}

func (p *Processor) getPoliticalEngagement(political *pb.PoliticalSocial) string {
	score := 0

	// Political leaning specified
	if political.PoliticalLeaning != "" && political.PoliticalLeaning != "apolitical" {
		score += 1
	}

	// Has activism
	if len(political.Activism) > 0 {
		score += 2
	}

	// Has causes
	if len(political.Causes) > 0 {
		score += 1
	}

	// Voting history
	switch p.normalizeString(political.VotingHistory) {
	case "always-votes":
		score += 3
	case "usually-votes":
		score += 2
	case "sometimes-votes":
		score += 1
	}

	switch {
	case score >= 6:
		return "very-high"
	case score >= 4:
		return "high"
	case score >= 2:
		return "moderate"
	case score >= 1:
		return "low"
	default:
		return "very-low"
	}
}

func (p *Processor) getVotingLikelihood(history string) string {
	likelihoodMap := map[string]string{
		"always-votes":           "very-high",
		"usually-votes":          "high",
		"sometimes-votes":        "moderate",
		"rarely-votes":           "low",
		"never-voted":            "very-low",
		"eligible-but-never-voted": "very-low",
		"not-eligible":           "not-applicable",
		"first-time-voter":       "unknown",
	}

	normalized := p.normalizeString(history)
	if likelihood, exists := likelihoodMap[normalized]; exists {
		return likelihood
	}
	return "unknown"
}

func (p *Processor) categorizeActivism(activism []string) map[string][]string {
	categories := make(map[string][]string)

	activismCategories := map[string]string{
		"voting":           "electoral",
		"campaigning":      "electoral",
		"canvassing":       "electoral",
		"phone-banking":    "electoral",
		"protesting":       "direct-action",
		"marching":         "direct-action",
		"civil-disobedience": "direct-action",
		"boycotting":       "economic",
		"divesting":        "economic",
		"fundraising":      "financial",
		"donating":         "financial",
		"volunteering":     "community",
		"organizing":       "community",
		"advocacy":         "advocacy",
		"lobbying":         "advocacy",
		"petitioning":      "advocacy",
		"social-media":     "digital",
		"blogging":         "digital",
		"online-organizing": "digital",
	}

	for _, activity := range activism {
		normalized := p.normalizeString(activity)
		if category, exists := activismCategories[normalized]; exists {
			categories[category] = append(categories[category], activity)
		} else {
			categories["other"] = append(categories["other"], activity)
		}
	}

	return categories
}

func (p *Processor) getActivismLevel(activism []string) string {
	count := len(activism)
	switch {
	case count == 0:
		return "none"
	case count <= 2:
		return "low"
	case count <= 5:
		return "moderate"
	case count <= 10:
		return "high"
	default:
		return "very-high"
	}
}

func (p *Processor) categorizeCauses(causes []string) map[string][]string {
	categories := make(map[string][]string)

	causeCategories := map[string]string{
		"environment":        "environmental",
		"climate-change":     "environmental",
		"conservation":       "environmental",
		"human-rights":       "social-justice",
		"civil-rights":       "social-justice",
		"equality":           "social-justice",
		"justice":            "social-justice",
		"education":          "social-services",
		"healthcare":         "social-services",
		"poverty":            "social-services",
		"homelessness":       "social-services",
		"immigration":        "policy",
		"gun-control":        "policy",
		"abortion":           "policy",
		"taxation":           "economic",
		"labor":              "economic",
		"minimum-wage":       "economic",
		"animal-rights":      "animal-welfare",
		"animal-welfare":     "animal-welfare",
		"veterans":           "community",
		"seniors":            "community",
		"children":           "community",
		"technology":         "technology",
		"privacy":            "technology",
		"internet-freedom":   "technology",
	}

	for _, cause := range causes {
		normalized := p.normalizeString(cause)
		if category, exists := causeCategories[normalized]; exists {
			categories[category] = append(categories[category], cause)
		} else {
			categories["other"] = append(categories["other"], cause)
		}
	}

	return categories
}

func (p *Processor) categorizeSocialGroups(groups []string) map[string][]string {
	categories := make(map[string][]string)

	groupCategories := map[string]string{
		"political-party":    "political",
		"campaign":           "political",
		"union":              "professional",
		"professional":       "professional",
		"trade-association":  "professional",
		"religious":          "religious",
		"church":             "religious",
		"mosque":             "religious",
		"synagogue":          "religious",
		"community":          "community",
		"neighborhood":       "community",
		"civic":              "civic",
		"volunteer":          "civic",
		"charity":            "charitable",
		"nonprofit":          "charitable",
		"advocacy":           "advocacy",
		"activist":           "advocacy",
		"hobby":              "recreational",
		"sports":             "recreational",
		"cultural":           "cultural",
		"ethnic":             "cultural",
		"online":             "digital",
		"social-media":       "digital",
	}

	for _, group := range groups {
		normalized := p.normalizeString(group)
		found := false
		for key, category := range groupCategories {
			if strings.Contains(normalized, key) {
				categories[category] = append(categories[category], group)
				found = true
				break
			}
		}
		if !found {
			categories["other"] = append(categories["other"], group)
		}
	}

	return categories
}

func (p *Processor) getSocialEngagement(groups []string) string {
	count := len(groups)
	switch {
	case count == 0:
		return "none"
	case count <= 2:
		return "low"
	case count <= 5:
		return "moderate"
	case count <= 10:
		return "high"
	default:
		return "very-high"
	}
}

func (p *Processor) categorizeMedia(media []string) map[string][]string {
	categories := make(map[string][]string)

	mediaCategories := map[string]string{
		"newspaper":      "traditional",
		"magazine":       "traditional",
		"television":     "traditional",
		"radio":          "traditional",
		"social-media":   "digital",
		"blog":           "digital",
		"podcast":        "digital",
		"youtube":        "digital",
		"twitter":        "social",
		"facebook":       "social",
		"instagram":      "social",
		"reddit":         "social",
		"news-website":   "online",
		"online":         "online",
		"streaming":      "streaming",
		"netflix":        "streaming",
	}

	for _, medium := range media {
		normalized := p.normalizeString(medium)
		found := false
		for key, category := range mediaCategories {
			if strings.Contains(normalized, key) {
				categories[category] = append(categories[category], medium)
				found = true
				break
			}
		}
		if !found {
			categories["other"] = append(categories["other"], medium)
		}
	}

	return categories
}

func (p *Processor) getMediaDiversity(media []string) string {
	categories := p.categorizeMedia(media)
	categoryCount := len(categories)

	switch {
	case categoryCount >= 5:
		return "very-high"
	case categoryCount >= 4:
		return "high"
	case categoryCount >= 2:
		return "moderate"
	case categoryCount >= 1:
		return "low"
	default:
		return "none"
	}
}
