package idgen

import (
	"regexp"
	"strings"
)

var (
	uuidRegex    = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	shortIDRegex = regexp.MustCompile(`^[A-Za-z0-9]{1,}$`)
)

// ValidateUUID проверяет, является ли строка валидным UUID
func ValidateUUID(id string) bool {
	return uuidRegex.MatchString(strings.ToLower(id))
}

// ValidateShortID проверяет, является ли строка валидным коротким ID
func ValidateShortID(id string, length int) bool {
	return shortIDRegex.MatchString(id) && len(id) == length
}

// ValidatePrefixedID проверяет формат префиксного ID
func ValidatePrefixedID(id, prefix string, length int) bool {
	parts := strings.Split(id, "-")
	if len(parts) != 2 {
		return false
	}

	return strings.EqualFold(parts[0], prefix) &&
		ValidateShortID(parts[1], length)
}
