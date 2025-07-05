package preferences

import (
	"fmt"
	"strings"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/config"
	pb "github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/grpc/pb"
)

// Processor handles preferences attribute processing and validation
type Processor struct {
	config *config.ValidationConfig
}

// NewProcessor creates a new preferences processor
func NewProcessor(cfg *config.ValidationConfig) *Processor {
	return &Processor{
		config: cfg,
	}
}

// ValidatePreferences validates preferences data
func (p *Processor) ValidatePreferences(prefs *pb.Preferences) []config.ValidationError {
	var errors []config.ValidationError

	// Validate hobbies list
	if len(prefs.Hobbies) > 30 {
		errors = append(errors, config.ValidationError{
			Field:   "hobbies",
			Message: "too many hobbies specified (max 30)",
		})
	}

	// Validate interests list
	if len(prefs.Interests) > 40 {
		errors = append(errors, config.ValidationError{
			Field:   "interests",
			Message: "too many interests specified (max 40)",
		})
	}

	// Validate favorite foods list
	if len(prefs.FavoriteFoods) > 25 {
		errors = append(errors, config.ValidationError{
			Field:   "favorite_foods",
			Message: "too many favorite foods specified (max 25)",
		})
	}

	// Validate favorite music list
	if len(prefs.FavoriteMusic) > 20 {
		errors = append(errors, config.ValidationError{
			Field:   "favorite_music",
			Message: "too many favorite music entries specified (max 20)",
		})
	}

	// Validate favorite movies list
	if len(prefs.FavoriteMovies) > 20 {
		errors = append(errors, config.ValidationError{
			Field:   "favorite_movies",
			Message: "too many favorite movies specified (max 20)",
		})
	}

	// Validate favorite books list
	if len(prefs.FavoriteBooks) > 20 {
		errors = append(errors, config.ValidationError{
			Field:   "favorite_books",
			Message: "too many favorite books specified (max 20)",
		})
	}

	// Validate technology use
	if prefs.TechnologyUse != "" {
		validLevels := []string{
			"minimal", "basic", "moderate", "advanced", "expert",
			"early-adopter", "tech-savvy", "digital-native", "luddite",
		}
		if !p.isValidOption(prefs.TechnologyUse, validLevels) {
			errors = append(errors, config.ValidationError{
				Field:   "technology_use",
				Message: "invalid technology use level",
			})
		}
	}

	// Validate travel style
	if prefs.TravelStyle != "" {
		validStyles := []string{
			"luxury", "budget", "backpacking", "business", "family",
			"adventure", "cultural", "relaxation", "eco-tourism",
			"solo", "group", "domestic", "international", "none",
		}
		if !p.isValidOption(prefs.TravelStyle, validStyles) {
			errors = append(errors, config.ValidationError{
				Field:   "travel_style",
				Message: "invalid travel style",
			})
		}
	}

	return errors
}

// ProcessPreferences processes and enriches preferences data
func (p *Processor) ProcessPreferences(prefs *pb.Preferences) (*pb.Preferences, error) {
	if prefs == nil {
		return nil, fmt.Errorf("preferences cannot be nil")
	}

	// Validate first
	if validationErrors := p.ValidatePreferences(prefs); len(validationErrors) > 0 {
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	// Create processed copy
	processed := &pb.Preferences{
		Hobbies:        p.normalizeStringSlice(prefs.Hobbies),
		Interests:      p.normalizeStringSlice(prefs.Interests),
		FavoriteFoods:  p.normalizeStringSlice(prefs.FavoriteFoods),
		FavoriteMusic:  p.normalizeStringSlice(prefs.FavoriteMusic),
		FavoriteMovies: p.normalizeStringSlice(prefs.FavoriteMovies),
		FavoriteBooks:  p.normalizeStringSlice(prefs.FavoriteBooks),
		TechnologyUse:  p.normalizeString(prefs.TechnologyUse),
		TravelStyle:    p.normalizeString(prefs.TravelStyle),
	}

	return processed, nil
}

// GetPreferencesProfile returns comprehensive preferences profile
func (p *Processor) GetPreferencesProfile(prefs *pb.Preferences) map[string]interface{} {
	profile := make(map[string]interface{})

	// Categorize hobbies
	if len(prefs.Hobbies) > 0 {
		profile["hobby_categories"] = p.categorizeHobbies(prefs.Hobbies)
		profile["hobby_diversity"] = p.getHobbyDiversity(prefs.Hobbies)
	}

	// Categorize interests
	if len(prefs.Interests) > 0 {
		profile["interest_categories"] = p.categorizeInterests(prefs.Interests)
		profile["interest_breadth"] = p.getInterestBreadth(prefs.Interests)
	}

	// Food preferences
	if len(prefs.FavoriteFoods) > 0 {
		profile["food_categories"] = p.categorizeFoods(prefs.FavoriteFoods)
		profile["culinary_adventurousness"] = p.getCulinaryAdventurousness(prefs.FavoriteFoods)
	}

	// Entertainment preferences
	profile["entertainment_profile"] = p.getEntertainmentProfile(prefs)

	// Technology and travel
	if prefs.TechnologyUse != "" {
		profile["technology_use"] = prefs.TechnologyUse
		profile["digital_engagement"] = p.getDigitalEngagement(prefs.TechnologyUse)
	}

	if prefs.TravelStyle != "" {
		profile["travel_style"] = prefs.TravelStyle
		profile["travel_personality"] = p.getTravelPersonality(prefs.TravelStyle)
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

func (p *Processor) categorizeHobbies(hobbies []string) map[string][]string {
	categories := make(map[string][]string)

	hobbyCategories := map[string]string{
		"reading":      "intellectual",
		"writing":      "intellectual",
		"chess":        "intellectual",
		"puzzles":      "intellectual",
		"running":      "physical",
		"cycling":      "physical",
		"swimming":     "physical",
		"hiking":       "physical",
		"yoga":         "physical",
		"gym":          "physical",
		"painting":     "creative",
		"drawing":      "creative",
		"photography":  "creative",
		"music":        "creative",
		"singing":      "creative",
		"dancing":      "creative",
		"cooking":      "creative",
		"gardening":    "creative",
		"gaming":       "digital",
		"programming":  "digital",
		"blogging":     "digital",
		"social-media": "digital",
		"traveling":    "social",
		"volunteering": "social",
		"networking":   "social",
		"collecting":   "collecting",
		"crafting":     "hands-on",
		"woodworking":  "hands-on",
		"mechanics":    "hands-on",
	}

	for _, hobby := range hobbies {
		normalized := p.normalizeString(hobby)
		found := false
		for key, category := range hobbyCategories {
			if strings.Contains(normalized, key) {
				categories[category] = append(categories[category], hobby)
				found = true
				break
			}
		}
		if !found {
			categories["other"] = append(categories["other"], hobby)
		}
	}

	return categories
}

func (p *Processor) getHobbyDiversity(hobbies []string) string {
	categories := p.categorizeHobbies(hobbies)
	categoryCount := len(categories)

	switch {
	case categoryCount >= 6:
		return "very-diverse"
	case categoryCount >= 4:
		return "diverse"
	case categoryCount >= 2:
		return "moderate"
	case categoryCount >= 1:
		return "focused"
	default:
		return "none"
	}
}

func (p *Processor) categorizeInterests(interests []string) map[string][]string {
	categories := make(map[string][]string)

	interestCategories := map[string]string{
		"science":      "academic",
		"history":      "academic",
		"philosophy":   "academic",
		"psychology":   "academic",
		"literature":   "academic",
		"mathematics":  "academic",
		"technology":   "technology",
		"ai":           "technology",
		"robotics":     "technology",
		"space":        "technology",
		"environment":  "social",
		"politics":     "social",
		"economics":    "social",
		"culture":      "social",
		"art":          "creative",
		"design":       "creative",
		"fashion":      "creative",
		"architecture": "creative",
		"sports":       "physical",
		"fitness":      "physical",
		"health":       "physical",
		"nutrition":    "physical",
		"travel":       "lifestyle",
		"food":         "lifestyle",
		"wine":         "lifestyle",
		"nature":       "lifestyle",
		"animals":      "lifestyle",
		"business":     "professional",
		"finance":      "professional",
		"marketing":    "professional",
		"leadership":   "professional",
	}

	for _, interest := range interests {
		normalized := p.normalizeString(interest)
		found := false
		for key, category := range interestCategories {
			if strings.Contains(normalized, key) {
				categories[category] = append(categories[category], interest)
				found = true
				break
			}
		}
		if !found {
			categories["other"] = append(categories["other"], interest)
		}
	}

	return categories
}

func (p *Processor) getInterestBreadth(interests []string) string {
	categories := p.categorizeInterests(interests)
	categoryCount := len(categories)

	switch {
	case categoryCount >= 7:
		return "very-broad"
	case categoryCount >= 5:
		return "broad"
	case categoryCount >= 3:
		return "moderate"
	case categoryCount >= 1:
		return "narrow"
	default:
		return "none"
	}
}

func (p *Processor) categorizeFoods(foods []string) map[string][]string {
	categories := make(map[string][]string)

	foodCategories := map[string]string{
		"italian":    "cuisine",
		"chinese":    "cuisine",
		"japanese":   "cuisine",
		"mexican":    "cuisine",
		"indian":     "cuisine",
		"thai":       "cuisine",
		"french":     "cuisine",
		"greek":      "cuisine",
		"pizza":      "comfort",
		"burger":     "comfort",
		"pasta":      "comfort",
		"ice-cream":  "comfort",
		"chocolate":  "comfort",
		"sushi":      "sophisticated",
		"wine":       "sophisticated",
		"cheese":     "sophisticated",
		"seafood":    "sophisticated",
		"steak":      "sophisticated",
		"salad":      "healthy",
		"fruit":      "healthy",
		"vegetables": "healthy",
		"smoothie":   "healthy",
		"quinoa":     "healthy",
		"spicy":      "flavor-profile",
		"sweet":      "flavor-profile",
		"savory":     "flavor-profile",
		"bitter":     "flavor-profile",
	}

	for _, food := range foods {
		normalized := p.normalizeString(food)
		found := false
		for key, category := range foodCategories {
			if strings.Contains(normalized, key) {
				categories[category] = append(categories[category], food)
				found = true
				break
			}
		}
		if !found {
			categories["other"] = append(categories["other"], food)
		}
	}

	return categories
}

func (p *Processor) getCulinaryAdventurousness(foods []string) string {
	categories := p.categorizeFoods(foods)

	// Count exotic cuisines and sophisticated foods
	adventurousCount := 0
	if cuisines, exists := categories["cuisine"]; exists {
		adventurousCount += len(cuisines)
	}
	if sophisticated, exists := categories["sophisticated"]; exists {
		adventurousCount += len(sophisticated)
	}

	switch {
	case adventurousCount >= 8:
		return "very-adventurous"
	case adventurousCount >= 5:
		return "adventurous"
	case adventurousCount >= 2:
		return "moderate"
	case adventurousCount >= 1:
		return "somewhat-adventurous"
	default:
		return "traditional"
	}
}

func (p *Processor) getEntertainmentProfile(prefs *pb.Preferences) map[string]interface{} {
	profile := make(map[string]interface{})

	// Music preferences
	if len(prefs.FavoriteMusic) > 0 {
		profile["music_genres"] = p.categorizeMusicGenres(prefs.FavoriteMusic)
		profile["music_diversity"] = p.getMusicDiversity(prefs.FavoriteMusic)
	}

	// Movie preferences
	if len(prefs.FavoriteMovies) > 0 {
		profile["movie_genres"] = p.categorizeMovieGenres(prefs.FavoriteMovies)
		profile["movie_preferences"] = p.getMoviePreferences(prefs.FavoriteMovies)
	}

	// Book preferences
	if len(prefs.FavoriteBooks) > 0 {
		profile["book_genres"] = p.categorizeBookGenres(prefs.FavoriteBooks)
		profile["reading_preferences"] = p.getReadingPreferences(prefs.FavoriteBooks)
	}

	return profile
}

func (p *Processor) categorizeMusicGenres(music []string) map[string][]string {
	categories := make(map[string][]string)

	genreMap := map[string]string{
		"rock":        "rock",
		"pop":         "pop",
		"jazz":        "jazz",
		"classical":   "classical",
		"hip-hop":     "hip-hop",
		"rap":         "hip-hop",
		"country":     "country",
		"blues":       "blues",
		"electronic":  "electronic",
		"folk":        "folk",
		"reggae":      "reggae",
		"metal":       "metal",
		"punk":        "punk",
		"indie":       "indie",
		"alternative": "alternative",
	}

	for _, item := range music {
		normalized := p.normalizeString(item)
		found := false
		for key, genre := range genreMap {
			if strings.Contains(normalized, key) {
				categories[genre] = append(categories[genre], item)
				found = true
				break
			}
		}
		if !found {
			categories["other"] = append(categories["other"], item)
		}
	}

	return categories
}

func (p *Processor) getMusicDiversity(music []string) string {
	genres := p.categorizeMusicGenres(music)
	genreCount := len(genres)

	switch {
	case genreCount >= 6:
		return "very-diverse"
	case genreCount >= 4:
		return "diverse"
	case genreCount >= 2:
		return "moderate"
	case genreCount >= 1:
		return "focused"
	default:
		return "none"
	}
}

func (p *Processor) categorizeMovieGenres(movies []string) map[string][]string {
	categories := make(map[string][]string)

	genreMap := map[string]string{
		"action":      "action",
		"comedy":      "comedy",
		"drama":       "drama",
		"horror":      "horror",
		"thriller":    "thriller",
		"romance":     "romance",
		"sci-fi":      "sci-fi",
		"fantasy":     "fantasy",
		"documentary": "documentary",
		"animation":   "animation",
		"mystery":     "mystery",
		"crime":       "crime",
		"war":         "war",
		"western":     "western",
	}

	for _, movie := range movies {
		normalized := p.normalizeString(movie)
		found := false
		for key, genre := range genreMap {
			if strings.Contains(normalized, key) {
				categories[genre] = append(categories[genre], movie)
				found = true
				break
			}
		}
		if !found {
			categories["other"] = append(categories["other"], movie)
		}
	}

	return categories
}

func (p *Processor) getMoviePreferences(movies []string) string {
	genres := p.categorizeMovieGenres(movies)

	// Analyze preference patterns
	if len(genres["documentary"]) > 0 || len(genres["drama"]) > 0 {
		return "serious"
	} else if len(genres["comedy"]) > 0 || len(genres["animation"]) > 0 {
		return "light-hearted"
	} else if len(genres["action"]) > 0 || len(genres["thriller"]) > 0 {
		return "excitement-seeking"
	} else if len(genres["sci-fi"]) > 0 || len(genres["fantasy"]) > 0 {
		return "imaginative"
	}

	return "varied"
}

func (p *Processor) categorizeBookGenres(books []string) map[string][]string {
	categories := make(map[string][]string)

	genreMap := map[string]string{
		"fiction":     "fiction",
		"non-fiction": "non-fiction",
		"biography":   "biography",
		"history":     "history",
		"science":     "science",
		"philosophy":  "philosophy",
		"psychology":  "psychology",
		"business":    "business",
		"self-help":   "self-help",
		"romance":     "romance",
		"mystery":     "mystery",
		"thriller":    "thriller",
		"fantasy":     "fantasy",
		"sci-fi":      "sci-fi",
		"poetry":      "poetry",
	}

	for _, book := range books {
		normalized := p.normalizeString(book)
		found := false
		for key, genre := range genreMap {
			if strings.Contains(normalized, key) {
				categories[genre] = append(categories[genre], book)
				found = true
				break
			}
		}
		if !found {
			categories["other"] = append(categories["other"], book)
		}
	}

	return categories
}

func (p *Processor) getReadingPreferences(books []string) string {
	genres := p.categorizeBookGenres(books)

	// Analyze reading patterns
	if len(genres["non-fiction"]) > len(genres["fiction"]) {
		return "educational"
	} else if len(genres["fiction"]) > len(genres["non-fiction"]) {
		return "entertainment"
	} else if len(genres["self-help"]) > 0 || len(genres["business"]) > 0 {
		return "self-improvement"
	}

	return "balanced"
}

func (p *Processor) getDigitalEngagement(techUse string) string {
	normalized := p.normalizeString(techUse)

	engagementMap := map[string]string{
		"expert":         "very-high",
		"early-adopter":  "very-high",
		"tech-savvy":     "high",
		"digital-native": "high",
		"advanced":       "high",
		"moderate":       "moderate",
		"basic":          "low",
		"minimal":        "very-low",
		"luddite":        "very-low",
	}

	if engagement, exists := engagementMap[normalized]; exists {
		return engagement
	}
	return "unknown"
}

func (p *Processor) getTravelPersonality(travelStyle string) string {
	normalized := p.normalizeString(travelStyle)

	personalityMap := map[string]string{
		"luxury":        "comfort-seeker",
		"budget":        "value-conscious",
		"backpacking":   "adventurous",
		"business":      "efficient",
		"family":        "family-oriented",
		"adventure":     "thrill-seeker",
		"cultural":      "culture-enthusiast",
		"relaxation":    "relaxation-focused",
		"eco-tourism":   "environmentally-conscious",
		"solo":          "independent",
		"group":         "social",
		"domestic":      "local-explorer",
		"international": "global-explorer",
		"none":          "homebody",
	}

	if personality, exists := personalityMap[normalized]; exists {
		return personality
	}
	return "unknown"
}

// GetPreferenceProfile creates a comprehensive preference profile
func (p *Processor) GetPreferenceProfile(prefs *pb.Preferences) map[string]interface{} {
	return p.GetPreferencesProfile(prefs)
}

func (p *Processor) calculateInterestDiversity(interests []string) float64 {
	categories := p.categorizeInterests(interests)

	if len(interests) == 0 {
		return 0.0
	}

	// Calculate diversity as number of categories / total possible categories
	maxCategories := 10.0 // Number of predefined categories
	actualCategories := float64(len(categories))

	return actualCategories / maxCategories
}
