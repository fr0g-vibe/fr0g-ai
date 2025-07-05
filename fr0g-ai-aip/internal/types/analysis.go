package types

import "time"

// PersonaInsights provides analytical insights about a persona's identities
type PersonaInsights struct {
	PersonaID                string                 `json:"persona_id"`
	PersonaName              string                 `json:"persona_name"`
	IdentityCount            int                    `json:"identity_count"`
	DemographicsDistribution map[string]interface{} `json:"demographics_distribution"`
	PersonalityPatterns      map[string]interface{} `json:"personality_patterns"`
	CulturalDiversity        map[string]interface{} `json:"cultural_diversity"`
	GeneratedAt              time.Time              `json:"generated_at"`
}

// IdentityProfile provides a comprehensive profile for an identity
type IdentityProfile struct {
	IdentityID           string                 `json:"identity_id"`
	Name                 string                 `json:"name"`
	DemographicProfile   map[string]interface{} `json:"demographic_profile,omitempty"`
	PsychographicProfile map[string]interface{} `json:"psychographic_profile,omitempty"`
	PreferenceProfile    map[string]interface{} `json:"preference_profile,omitempty"`
	CulturalProfile      map[string]interface{} `json:"cultural_profile,omitempty"`
	PoliticalProfile     map[string]interface{} `json:"political_profile,omitempty"`
	HealthProfile        map[string]interface{} `json:"health_profile,omitempty"`
	BehavioralProfile    map[string]interface{} `json:"behavioral_profile,omitempty"`
	GeneratedAt          time.Time              `json:"generated_at"`
}
