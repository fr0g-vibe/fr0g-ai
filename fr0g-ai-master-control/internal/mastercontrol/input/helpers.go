package input

import "strconv"

// Helper functions for parsing webhook data

// getStringFromMap safely extracts a string value from a map
func getStringFromMap(m map[string]interface{}, key string) string {
	if value, ok := m[key]; ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// getBoolFromMap safely extracts a boolean value from a map
func getBoolFromMap(m map[string]interface{}, key string) bool {
	if value, ok := m[key]; ok {
		if b, ok := value.(bool); ok {
			return b
		}
		// Try to parse string as bool
		if str, ok := value.(string); ok {
			if b, err := strconv.ParseBool(str); err == nil {
				return b
			}
		}
	}
	return false
}

// getIntFromMap safely extracts an integer value from a map
func getIntFromMap(m map[string]interface{}, key string) int {
	if value, ok := m[key]; ok {
		if i, ok := value.(int); ok {
			return i
		}
		if f, ok := value.(float64); ok {
			return int(f)
		}
		// Try to parse string as int
		if str, ok := value.(string); ok {
			if i, err := strconv.Atoi(str); err == nil {
				return i
			}
		}
	}
	return 0
}

// getFloat64FromMap safely extracts a float64 value from a map
func getFloat64FromMap(m map[string]interface{}, key string) float64 {
	if value, ok := m[key]; ok {
		if f, ok := value.(float64); ok {
			return f
		}
		if i, ok := value.(int); ok {
			return float64(i)
		}
		// Try to parse string as float
		if str, ok := value.(string); ok {
			if f, err := strconv.ParseFloat(str, 64); err == nil {
				return f
			}
		}
	}
	return 0.0
}
