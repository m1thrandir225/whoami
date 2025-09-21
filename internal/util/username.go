package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

// GenerateUniqueUsername generates a unique username based on the provided name
// It sanitizes the name and adds a random suffix if needed to ensure uniqueness
func GenerateUniqueUsername(baseName string, isUsernameAvailable func(string) bool) string {
	if baseName == "" {
		baseName = "user"
	}

	// Sanitize the base name
	username := sanitizeUsername(baseName)

	// If the sanitized username is empty, use a default
	if username == "" {
		username = "user"
	}

	// Check if the base username is available
	if isUsernameAvailable(username) {
		return username
	}

	// Try with random suffixes
	for i := 0; i < 100; i++ { // Try up to 100 times
		suffix := generateRandomSuffix()
		candidate := fmt.Sprintf("%s_%s", username, suffix)

		if isUsernameAvailable(candidate) {
			return candidate
		}
	}

	// If all attempts fail, use a completely random username
	return generateRandomUsername()
}

// sanitizeUsername cleans and formats a username to be valid
func sanitizeUsername(name string) string {
	// Convert to lowercase
	username := strings.ToLower(name)

	// Remove special characters, keep only alphanumeric and underscores
	reg := regexp.MustCompile(`[^a-z0-9_]`)
	username = reg.ReplaceAllString(username, "")

	// Remove multiple consecutive underscores
	reg = regexp.MustCompile(`_+`)
	username = reg.ReplaceAllString(username, "_")

	// Remove leading/trailing underscores
	username = strings.Trim(username, "_")

	// Ensure minimum length
	if len(username) < 3 {
		username = "user"
	}

	// Ensure maximum length (database constraint is 100, but let's be safe)
	if len(username) > 50 {
		username = username[:50]
	}

	return username
}

// generateRandomSuffix creates a random suffix for usernames
func generateRandomSuffix() string {
	// Generate a random number between 1000 and 9999
	n, err := rand.Int(rand.Reader, big.NewInt(9000))
	if err != nil {
		// Fallback to a simple counter-based approach
		return "1234"
	}
	return fmt.Sprintf("%d", n.Int64()+1000)
}

// generateRandomUsername creates a completely random username
func generateRandomUsername() string {
	// Generate a random number
	n, err := rand.Int(rand.Reader, big.NewInt(999999999))
	if err != nil {
		return "user_12345"
	}
	return fmt.Sprintf("user_%d", n.Int64())
}
