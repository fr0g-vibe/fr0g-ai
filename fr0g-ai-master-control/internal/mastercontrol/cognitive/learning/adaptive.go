package learning

import (
	"math"
	"sync"
	"time"
)

// AdaptiveLearning implements adaptive learning algorithms
type AdaptiveLearning struct {
	learningRate     float64
	momentum         float64
	adaptationFactor float64
	experiences      []Experience
	patterns         map[string]*LearnedPattern
	mu               sync.RWMutex
}

// Experience represents a learning experience
type Experience struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Input       map[string]interface{} `json:"input"`
	Output      map[string]interface{} `json:"output"`
	Feedback    float64                `json:"feedback"` // -1.0 to 1.0
	Timestamp   time.Time              `json:"timestamp"`
	Context     map[string]interface{} `json:"context"`
	Importance  float64                `json:"importance"`
}

// LearnedPattern represents a pattern learned from experiences
type LearnedPattern struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"`
	Triggers     []string               `json:"triggers"`
	Responses    []string               `json:"responses"`
	Confidence   float64                `json:"confidence"`
	UsageCount   int                    `json:"usage_count"`
	SuccessRate  float64                `json:"success_rate"`
	LastUsed     time.Time              `json:"last_used"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// NewAdaptiveLearning creates a new adaptive learning system
func NewAdaptiveLearning() *AdaptiveLearning {
	return &AdaptiveLearning{
		learningRate:     0.1,
		momentum:         0.9,
		adaptationFactor: 0.05,
		experiences:      make([]Experience, 0),
		patterns:         make(map[string]*LearnedPattern),
	}
}

// AddExperience adds a new learning experience
func (al *AdaptiveLearning) AddExperience(exp Experience) {
	al.mu.Lock()
	defer al.mu.Unlock()
	
	exp.Timestamp = time.Now()
	al.experiences = append(al.experiences, exp)
	
	// Trigger learning from this experience
	al.learnFromExperience(exp)
	
	// Maintain experience buffer size
	if len(al.experiences) > 1000 {
		al.experiences = al.experiences[100:] // Keep most recent 900
	}
}

// GetLearningRate returns current learning rate
func (al *AdaptiveLearning) GetLearningRate() float64 {
	al.mu.RLock()
	defer al.mu.RUnlock()
	return al.learningRate
}

// GetAdaptationScore calculates adaptation score based on recent performance
func (al *AdaptiveLearning) GetAdaptationScore() float64 {
	al.mu.RLock()
	defer al.mu.RUnlock()
	
	if len(al.experiences) == 0 {
		return 0.0
	}
	
	// Calculate adaptation based on recent feedback trends
	recentExperiences := al.getRecentExperiences(50)
	if len(recentExperiences) == 0 {
		return 0.0
	}
	
	totalFeedback := 0.0
	improvementTrend := 0.0
	
	for i, exp := range recentExperiences {
		totalFeedback += exp.Feedback
		
		// Calculate improvement trend
		if i > 0 {
			if exp.Feedback > recentExperiences[i-1].Feedback {
				improvementTrend += 0.1
			} else if exp.Feedback < recentExperiences[i-1].Feedback {
				improvementTrend -= 0.05
			}
		}
	}
	
	avgFeedback := totalFeedback / float64(len(recentExperiences))
	
	// Normalize to 0-1 range
	adaptationScore := (avgFeedback + 1.0) / 2.0 // Convert from [-1,1] to [0,1]
	adaptationScore += improvementTrend * 0.1    // Add trend bonus
	
	if adaptationScore > 1.0 {
		adaptationScore = 1.0
	}
	if adaptationScore < 0.0 {
		adaptationScore = 0.0
	}
	
	return adaptationScore
}

// learnFromExperience processes an experience and updates learning
func (al *AdaptiveLearning) learnFromExperience(exp Experience) {
	// Update learning rate based on feedback
	al.updateLearningRate(exp.Feedback)
}

// updateLearningRate adjusts learning rate based on feedback
func (al *AdaptiveLearning) updateLearningRate(feedback float64) {
	// Positive feedback increases learning rate, negative decreases it
	adjustment := feedback * al.adaptationFactor
	
	newRate := al.learningRate + adjustment
	
	// Keep learning rate in reasonable bounds
	if newRate > 0.5 {
		newRate = 0.5
	}
	if newRate < 0.01 {
		newRate = 0.01
	}
	
	al.learningRate = newRate
}

// getRecentExperiences returns the most recent N experiences
func (al *AdaptiveLearning) getRecentExperiences(n int) []Experience {
	if len(al.experiences) == 0 {
		return []Experience{}
	}
	
	start := len(al.experiences) - n
	if start < 0 {
		start = 0
	}
	
	return al.experiences[start:]
}
