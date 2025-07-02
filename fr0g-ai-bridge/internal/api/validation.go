package api

// isValidRole checks if the role is one of the allowed values
func isValidRole(role string) bool {
	validRoles := []string{"system", "user", "assistant", "function"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}
