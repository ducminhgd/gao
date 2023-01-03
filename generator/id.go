// This package is for ID Generator
// Default NanoID size is 21 and runes are "_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
package generator

import (
	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// NewNanoID creates a new Nano ID string with default size and alphabet runes
func NewNanoID() (string, error) {
	return gonanoid.New()
}

// GenerateNanoID generates a new Nano ID string with custom alphabet characters and custom size
func GenerateNanoID(alphabet string, size int) (string, error) {
	return gonanoid.Generate(alphabet, size)
}

// NewUUID generates a new UUID v4 string
func NewUUID() string {
	return uuid.New().String()
}
