package slices

// Contains returns true if a string slice contains a given element, false otherwise
func Contains(slice []string, element string) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}
	return false
}
