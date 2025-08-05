package idgen

import "github.com/google/uuid"

// GenerateUUID генерирует UUID v4
func GenerateUUID() string {
	return uuid.NewString()
}
