package uuids

import "github.com/google/uuid"

// NewString returns new UUID as a string.
func NewString() string {
	return uuid.NewString()
}

// IsUUID validates if value is a valid UUID string.
func IsUUID(value string) bool {
	if value == "" {
		return false
	}

	_, err := uuid.Parse(value)

	return err == nil
}
