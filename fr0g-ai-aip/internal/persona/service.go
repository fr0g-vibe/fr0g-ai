package persona

import (
	"fmt"
	"time"

	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/attributes/behavioral"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/attributes/cultural"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/attributes/demographics"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/attributes/health"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/attributes/lifehistory"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/attributes/political"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/attributes/preferences"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/attributes/psychographics"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/config"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/storage"
	"github.com/fr0g-vibe/fr0g-ai/fr0g-ai-aip/internal/types"
	sharedconfig "github.com/fr0g-vibe/fr0g-ai/pkg/config"
	"github.com/google/uuid"
)

// Service provides core business logic for persona and identity management
type Service struct {
	storage                 storage.Storage
	config                  *config.Config
	demographicsProcessor   *demographics.Processor
	psychographicsProcessor *psychographics.Processor
	lifeHistoryProcessor    *lifehistory.Processor
	preferencesProcessor    *preferences.Processor
	culturalProcessor       *cultural.Processor
	politicalProcessor      *political.Processor
	healthProcessor         *health.Processor
	behavioralProcessor     *behavioral.Processor
}

// NewService creates a new persona service with all attribute processors
func NewService(storage storage.Storage) *Service {
	cfg := config.Load()

	return &Service{
		storage:                 storage,
		config:                  cfg,
		demographicsProcessor:   demographics.NewProcessor(&cfg.Validation),
		psychographicsProcessor: psychographics.NewProcessor(&cfg.Validation),
		lifeHistoryProcessor:    lifehistory.NewProcessor(&cfg.Validation),
		preferencesProcessor:    preferences.NewProcessor(&cfg.Validation),
		culturalProcessor:       cultural.NewProcessor(&cfg.Validation),
		politicalProcessor:      political.NewProcessor(&cfg.Validation),
		healthProcessor:         health.NewProcessor(&cfg.Validation),
		behavioralProcessor:     behavioral.NewProcessor(&cfg.Validation),
	}
}

// GetStorage returns the underlying storage interface
func (s *Service) GetStorage() storage.Storage {
	return s.storage
}

// PERSONA OPERATIONS

// CreatePersona creates a new persona with validation
func (s *Service) CreatePersona(persona *types.Persona) error {
	if persona == nil {
		return fmt.Errorf("persona cannot be nil")
	}

	// Validate persona
	if err := s.validatePersona(persona); err != nil {
		return err
	}

	// Generate ID if not provided
	if persona.Id == "" {
		persona.Id = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	persona.CreatedAt = now
	persona.UpdatedAt = now

	// Store persona
	return s.storage.Create(persona)
}

// GetPersona retrieves a persona by ID
func (s *Service) GetPersona(id string) (types.Persona, error) {
	if id == "" {
		return types.Persona{}, fmt.Errorf("persona ID is required")
	}

	return s.storage.Get(id)
}

// ListPersonas returns all personas
func (s *Service) ListPersonas() ([]types.Persona, error) {
	return s.storage.List()
}

// UpdatePersona updates an existing persona
func (s *Service) UpdatePersona(id string, persona types.Persona) error {
	if id == "" {
		return fmt.Errorf("persona ID is required")
	}

	// Validate persona
	if err := s.validatePersona(&persona); err != nil {
		return err
	}

	// Ensure ID matches
	persona.Id = id
	persona.UpdatedAt = time.Now()

	return s.storage.Update(id, persona)
}

// DeletePersona removes a persona by ID
func (s *Service) DeletePersona(id string) error {
	if id == "" {
		return fmt.Errorf("persona ID is required")
	}

	return s.storage.Delete(id)
}

// IDENTITY OPERATIONS

// CreateIdentity creates a new identity with rich attribute processing
func (s *Service) CreateIdentity(identity *types.Identity) error {
	if identity == nil {
		return fmt.Errorf("identity cannot be nil")
	}

	// Validate identity
	if err := s.validateIdentity(identity); err != nil {
		return err
	}

	// Generate ID if not provided
	if identity.Id == "" {
		identity.Id = uuid.New().String()
	}

	// Process rich attributes if provided
	if identity.RichAttributes != nil {
		if err := s.processRichAttributes(identity.RichAttributes); err != nil {
			return fmt.Errorf("rich attributes validation failed: %w", err)
		}
	}

	// Set timestamps
	now := time.Now()
	identity.CreatedAt = now
	identity.UpdatedAt = now

	// Store identity
	return s.storage.CreateIdentity(identity)
}

// GetIdentity retrieves an identity by ID
func (s *Service) GetIdentity(id string) (types.Identity, error) {
	if id == "" {
		return types.Identity{}, fmt.Errorf("identity ID is required")
	}

	return s.storage.GetIdentity(id)
}

// ListIdentities returns identities with optional filtering
func (s *Service) ListIdentities(filter *types.IdentityFilter) ([]types.Identity, error) {
	return s.storage.ListIdentities(filter)
}

// UpdateIdentity updates an existing identity
func (s *Service) UpdateIdentity(id string, identity types.Identity) error {
	if id == "" {
		return fmt.Errorf("identity ID is required")
	}

	// Validate identity
	if err := s.validateIdentity(&identity); err != nil {
		return err
	}

	// Process rich attributes if provided
	if identity.RichAttributes != nil {
		if err := s.processRichAttributes(identity.RichAttributes); err != nil {
			return fmt.Errorf("rich attributes validation failed: %w", err)
		}
	}

	// Ensure ID matches
	identity.Id = id
	identity.UpdatedAt = time.Now()

	return s.storage.UpdateIdentity(id, identity)
}

// DeleteIdentity removes an identity by ID
func (s *Service) DeleteIdentity(id string) error {
	if id == "" {
		return fmt.Errorf("identity ID is required")
	}

	return s.storage.DeleteIdentity(id)
}

// GetIdentityWithPersona retrieves an identity with its associated persona
func (s *Service) GetIdentityWithPersona(id string) (types.IdentityWithPersona, error) {
	if id == "" {
		return types.IdentityWithPersona{}, fmt.Errorf("identity ID is required")
	}

	identity, err := s.storage.GetIdentity(id)
	if err != nil {
		return types.IdentityWithPersona{}, err
	}

	persona, err := s.storage.Get(identity.PersonaId)
	if err != nil {
		return types.IdentityWithPersona{}, fmt.Errorf("failed to get associated persona: %w", err)
	}

	return types.IdentityWithPersona{
		Identity: identity,
		Persona:  persona,
	}, nil
}

// VALIDATION METHODS

// validatePersona validates a persona using business rules
func (s *Service) validatePersona(persona *types.Persona) error {
	var errors []sharedconfig.ValidationError

	// Required fields
	if persona.Name == "" {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "name",
			Message: "persona name is required",
		})
	}

	if persona.Topic == "" {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "topic",
			Message: "persona topic is required",
		})
	}

	if persona.Prompt == "" {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "prompt",
			Message: "persona prompt is required",
		})
	}

	// Length validations
	if len(persona.Name) > 100 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "name",
			Message: "persona name must be 100 characters or less",
		})
	}

	if len(persona.Topic) > 200 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "topic",
			Message: "persona topic must be 200 characters or less",
		})
	}

	if len(persona.Prompt) > 10000 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "prompt",
			Message: "persona prompt must be 10000 characters or less",
		})
	}

	// Context validation
	if len(persona.Context) > 50 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "context",
			Message: "persona context cannot have more than 50 entries",
		})
	}

	for key, value := range persona.Context {
		if len(key) > 100 {
			errors = append(errors, sharedconfig.ValidationError{
				Field:   "context." + key,
				Message: "context key must be 100 characters or less",
			})
		}
		if len(value) > 1000 {
			errors = append(errors, sharedconfig.ValidationError{
				Field:   "context." + key,
				Message: "context value must be 1000 characters or less",
			})
		}
	}

	// RAG validation
	if len(persona.Rag) > 100 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "rag",
			Message: "persona RAG cannot have more than 100 entries",
		})
	}

	for i, rag := range persona.Rag {
		if len(rag) > 500 {
			errors = append(errors, sharedconfig.ValidationError{
				Field:   fmt.Sprintf("rag[%d]", i),
				Message: "RAG entry must be 500 characters or less",
			})
		}
	}

	if len(errors) > 0 {
		return sharedconfig.ValidationErrors(errors)
	}

	return nil
}

// validateIdentity validates an identity using business rules
func (s *Service) validateIdentity(identity *types.Identity) error {
	var errors []sharedconfig.ValidationError

	// Required fields
	if identity.PersonaId == "" {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "persona_id",
			Message: "persona ID is required",
		})
	}

	if identity.Name == "" {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "name",
			Message: "identity name is required",
		})
	}

	// Length validations
	if len(identity.Name) > 100 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "name",
			Message: "identity name must be 100 characters or less",
		})
	}

	if len(identity.Description) > 1000 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "description",
			Message: "identity description must be 1000 characters or less",
		})
	}

	if len(identity.Background) > 5000 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "background",
			Message: "identity background must be 5000 characters or less",
		})
	}

	// Tags validation
	if len(identity.Tags) > 20 {
		errors = append(errors, sharedconfig.ValidationError{
			Field:   "tags",
			Message: "identity cannot have more than 20 tags",
		})
	}

	for i, tag := range identity.Tags {
		if len(tag) > 50 {
			errors = append(errors, sharedconfig.ValidationError{
				Field:   fmt.Sprintf("tags[%d]", i),
				Message: "tag must be 50 characters or less",
			})
		}
	}

	// Verify persona exists
	if identity.PersonaId != "" {
		if _, err := s.storage.Get(identity.PersonaId); err != nil {
			errors = append(errors, sharedconfig.ValidationError{
				Field:   "persona_id",
				Message: "referenced persona does not exist",
			})
		}
	}

	if len(errors) > 0 {
		return sharedconfig.ValidationErrors(errors)
	}

	return nil
}

// processRichAttributes validates and processes all rich attributes
func (s *Service) processRichAttributes(attrs *types.RichAttributes) error {
	var allErrors []sharedconfig.ValidationError

	// Process Demographics
	if attrs.Demographics != nil {
		if errors := s.demographicsProcessor.ValidateDemographics(attrs.Demographics); len(errors) > 0 {
			allErrors = append(allErrors, errors...)
		}
	}

	// Process Psychographics
	if attrs.Psychographics != nil {
		if errors := s.psychographicsProcessor.ValidatePsychographics(attrs.Psychographics); len(errors) > 0 {
			allErrors = append(allErrors, errors...)
		}
	}

	// Process LifeHistory
	if attrs.LifeHistory != nil {
		if errors := s.lifeHistoryProcessor.ValidateLifeHistory(attrs.LifeHistory); len(errors) > 0 {
			allErrors = append(allErrors, errors...)
		}
	}

	// Process Preferences
	if attrs.Preferences != nil {
		if errors := s.preferencesProcessor.ValidatePreferences(attrs.Preferences); len(errors) > 0 {
			allErrors = append(allErrors, errors...)
		}
	}

	// Process CulturalReligious
	if attrs.CulturalReligious != nil {
		if errors := s.culturalProcessor.ValidateCulturalReligious(attrs.CulturalReligious); len(errors) > 0 {
			allErrors = append(allErrors, errors...)
		}
	}

	// Process PoliticalSocial
	if attrs.PoliticalSocial != nil {
		if errors := s.politicalProcessor.ValidatePoliticalSocial(attrs.PoliticalSocial); len(errors) > 0 {
			allErrors = append(allErrors, errors...)
		}
	}

	// Process Health
	if attrs.Health != nil {
		if errors := s.healthProcessor.ValidateHealth(attrs.Health); len(errors) > 0 {
			allErrors = append(allErrors, errors...)
		}
	}

	// Process BehavioralTendencies
	if attrs.BehavioralTendencies != nil {
		if errors := s.behavioralProcessor.ValidateBehavioralTendencies(attrs.BehavioralTendencies); len(errors) > 0 {
			allErrors = append(allErrors, errors...)
		}
	}

	if len(allErrors) > 0 {
		return sharedconfig.ValidationErrors(allErrors)
	}

	return nil
}

// ANALYSIS AND INSIGHTS METHODS

// GetPersonaInsights generates insights about a persona's identities
func (s *Service) GetPersonaInsights(personaID string) (*types.PersonaInsights, error) {
	if personaID == "" {
		return nil, fmt.Errorf("persona ID is required")
	}

	// Get persona
	persona, err := s.storage.Get(personaID)
	if err != nil {
		return nil, fmt.Errorf("persona not found: %w", err)
	}

	// Get all identities for this persona
	filter := &types.IdentityFilter{PersonaID: personaID}
	identities, err := s.storage.ListIdentities(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get identities: %w", err)
	}

	insights := &types.PersonaInsights{
		PersonaID:     personaID,
		PersonaName:   persona.Name,
		IdentityCount: len(identities),
		GeneratedAt:   time.Now(),
	}

	// Analyze demographics distribution
	insights.DemographicsDistribution = s.analyzeDemographicsDistribution(identities)

	// Analyze personality patterns
	insights.PersonalityPatterns = s.analyzePersonalityPatterns(identities)

	// Analyze cultural diversity
	insights.CulturalDiversity = s.analyzeCulturalDiversity(identities)

	return insights, nil
}

// GetIdentityProfile generates a comprehensive profile for an identity
func (s *Service) GetIdentityProfile(identityID string) (*types.IdentityProfile, error) {
	if identityID == "" {
		return nil, fmt.Errorf("identity ID is required")
	}

	identity, err := s.storage.GetIdentity(identityID)
	if err != nil {
		return nil, fmt.Errorf("identity not found: %w", err)
	}

	profile := &types.IdentityProfile{
		IdentityID:  identityID,
		Name:        identity.Name,
		GeneratedAt: time.Now(),
	}

	// Generate profiles from each attribute processor
	if identity.RichAttributes != nil {
		if identity.RichAttributes.Demographics != nil {
			profile.DemographicProfile = s.demographicsProcessor.GetDemographicProfile(identity.RichAttributes.Demographics)
		}

		if identity.RichAttributes.Psychographics != nil {
			cognitiveProfile := s.psychographicsProcessor.GetCognitiveProfile(identity.RichAttributes.Psychographics)
			profile.PsychographicProfile = make(map[string]interface{})
			for k, v := range cognitiveProfile {
				profile.PsychographicProfile[k] = v
			}
		}

		if identity.RichAttributes.Preferences != nil {
			profile.PreferenceProfile = s.preferencesProcessor.GetPreferenceProfile(identity.RichAttributes.Preferences)
		}

		if identity.RichAttributes.CulturalReligious != nil {
			profile.CulturalProfile = s.culturalProcessor.GetCulturalProfile(identity.RichAttributes.CulturalReligious)
		}

		if identity.RichAttributes.PoliticalSocial != nil {
			profile.PoliticalProfile = s.politicalProcessor.GetPoliticalProfile(identity.RichAttributes.PoliticalSocial)
		}

		if identity.RichAttributes.Health != nil {
			profile.HealthProfile = s.healthProcessor.GetHealthProfile(identity.RichAttributes.Health)
		}

		if identity.RichAttributes.BehavioralTendencies != nil {
			profile.BehavioralProfile = s.behavioralProcessor.GetBehavioralProfile(identity.RichAttributes.BehavioralTendencies)
		}
	}

	return profile, nil
}

// HELPER METHODS FOR ANALYSIS

func (s *Service) analyzeDemographicsDistribution(identities []types.Identity) map[string]interface{} {
	distribution := make(map[string]interface{})

	ageGroups := make(map[string]int)
	genders := make(map[string]int)
	educationLevels := make(map[string]int)

	for _, identity := range identities {
		if identity.RichAttributes != nil && identity.RichAttributes.Demographics != nil {
			demo := identity.RichAttributes.Demographics

			// Age groups
			if demo.Age > 0 {
				ageGroup := s.getAgeGroup(demo.Age)
				ageGroups[ageGroup]++
			}

			// Gender distribution
			if demo.Gender != "" {
				genders[demo.Gender]++
			}

			// Education levels
			if demo.Education != "" {
				educationLevels[demo.Education]++
			}
		}
	}

	distribution["age_groups"] = ageGroups
	distribution["genders"] = genders
	distribution["education_levels"] = educationLevels

	return distribution
}

func (s *Service) analyzePersonalityPatterns(identities []types.Identity) map[string]interface{} {
	patterns := make(map[string]interface{})

	var opennessSum, conscientiousnessSum, extraversionSum, agreeablenessSum, neuroticismSum float64
	var count int

	for _, identity := range identities {
		if identity.RichAttributes != nil &&
			identity.RichAttributes.Psychographics != nil &&
			identity.RichAttributes.Psychographics.Personality != nil {

			p := identity.RichAttributes.Psychographics.Personality
			opennessSum += p.Openness
			conscientiousnessSum += p.Conscientiousness
			extraversionSum += p.Extraversion
			agreeablenessSum += p.Agreeableness
			neuroticismSum += p.Neuroticism
			count++
		}
	}

	if count > 0 {
		patterns["average_openness"] = opennessSum / float64(count)
		patterns["average_conscientiousness"] = conscientiousnessSum / float64(count)
		patterns["average_extraversion"] = extraversionSum / float64(count)
		patterns["average_agreeableness"] = agreeablenessSum / float64(count)
		patterns["average_neuroticism"] = neuroticismSum / float64(count)
		patterns["sample_size"] = count
	}

	return patterns
}

func (s *Service) analyzeCulturalDiversity(identities []types.Identity) map[string]interface{} {
	diversity := make(map[string]interface{})

	religions := make(map[string]int)
	cultures := make(map[string]int)

	for _, identity := range identities {
		if identity.RichAttributes != nil && identity.RichAttributes.CulturalReligious != nil {
			cultural := identity.RichAttributes.CulturalReligious

			if cultural.Religion != "" {
				religions[cultural.Religion]++
			}

			if cultural.CulturalBackground != "" {
				cultures[cultural.CulturalBackground]++
			}
		}
	}

	diversity["religions"] = religions
	diversity["cultural_backgrounds"] = cultures
	diversity["religious_diversity"] = len(religions)
	diversity["cultural_diversity"] = len(cultures)

	return diversity
}

func (s *Service) getAgeGroup(age int32) string {
	switch {
	case age < 18:
		return "minor"
	case age < 25:
		return "young-adult"
	case age < 35:
		return "adult"
	case age < 50:
		return "middle-aged"
	case age < 65:
		return "mature"
	default:
		return "senior"
	}
}
