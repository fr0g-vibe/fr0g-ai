package patterns

import (
	"fmt"
	"sync"
	"time"
)

// PatternRecognition implements pattern recognition algorithms
type PatternRecognition struct {
	patterns    map[string]*RecognizedPattern
	dataStreams map[string]*DataStream
	mu          sync.RWMutex
}

// RecognizedPattern represents a recognized pattern
type RecognizedPattern struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Confidence  float64                `json:"confidence"`
	Frequency   int                    `json:"frequency"`
	LastSeen    time.Time              `json:"last_seen"`
	FirstSeen   time.Time              `json:"first_seen"`
	Triggers    []string               `json:"triggers"`
	Indicators  map[string]interface{} `json:"indicators"`
	Metadata    map[string]interface{} `json:"metadata"`
	Strength    float64                `json:"strength"`
}

// DataStream represents a stream of data for pattern analysis
type DataStream struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Data       []DataPoint `json:"data"`
	LastUpdate time.Time   `json:"last_update"`
}

// DataPoint represents a single data point in a stream
type DataPoint struct {
	Timestamp time.Time              `json:"timestamp"`
	Value     interface{}            `json:"value"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// NewPatternRecognition creates a new pattern recognition system
func NewPatternRecognition() *PatternRecognition {
	return &PatternRecognition{
		patterns:    make(map[string]*RecognizedPattern),
		dataStreams: make(map[string]*DataStream),
	}
}

// AddDataPoint adds a data point to a stream for pattern analysis
func (pr *PatternRecognition) AddDataPoint(streamID, streamType string, value interface{}, metadata map[string]interface{}) {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	
	stream, exists := pr.dataStreams[streamID]
	if !exists {
		stream = &DataStream{
			ID:   streamID,
			Type: streamType,
			Data: make([]DataPoint, 0),
		}
		pr.dataStreams[streamID] = stream
	}
	
	dataPoint := DataPoint{
		Timestamp: time.Now(),
		Value:     value,
		Metadata:  metadata,
	}
	
	stream.Data = append(stream.Data, dataPoint)
	stream.LastUpdate = time.Now()
	
	// Simple pattern recognition - detect frequent values
	pr.analyzeStream(stream)
}

// GetPatternCount returns the number of recognized patterns
func (pr *PatternRecognition) GetPatternCount() int {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	return len(pr.patterns)
}

// GetTopPatterns returns the top N patterns by confidence
func (pr *PatternRecognition) GetTopPatterns(n int) []*RecognizedPattern {
	pr.mu.RLock()
	defer pr.mu.RUnlock()
	
	patterns := make([]*RecognizedPattern, 0, len(pr.patterns))
	for _, pattern := range pr.patterns {
		patterns = append(patterns, pattern)
	}
	
	// Simple sorting by confidence
	if n > len(patterns) {
		n = len(patterns)
	}
	
	return patterns[:n]
}

// analyzeStream performs simple pattern recognition on a data stream
func (pr *PatternRecognition) analyzeStream(stream *DataStream) {
	if len(stream.Data) < 3 {
		return
	}
	
	// Simple frequency-based pattern detection
	valueFreq := make(map[string]int)
	for _, point := range stream.Data {
		key := fmt.Sprintf("%v", point.Value)
		valueFreq[key]++
	}
	
	// Create patterns for frequent values
	for value, freq := range valueFreq {
		if freq >= 2 { // Minimum frequency threshold
			patternID := fmt.Sprintf("freq_%s_%s", stream.ID, value)
			
			if _, exists := pr.patterns[patternID]; !exists {
				pattern := &RecognizedPattern{
					ID:          patternID,
					Type:        "frequency",
					Name:        fmt.Sprintf("Frequent Value: %s", value),
					Description: fmt.Sprintf("Value '%s' appears frequently in stream %s", value, stream.ID),
					Confidence:  float64(freq) / float64(len(stream.Data)),
					Frequency:   freq,
					FirstSeen:   time.Now(),
					LastSeen:    time.Now(),
					Triggers:    []string{value},
					Indicators: map[string]interface{}{
						"frequency":    freq,
						"total_points": len(stream.Data),
						"value":        value,
					},
					Metadata: map[string]interface{}{
						"recognizer": "frequency",
						"stream_id":  stream.ID,
					},
					Strength: float64(freq) / float64(len(stream.Data)),
				}
				
				pr.patterns[patternID] = pattern
			}
		}
	}
}
