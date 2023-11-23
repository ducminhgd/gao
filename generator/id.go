package generator

import (
	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// NewNanoID generates a new NanoID.
//
// It returns a string and an error.
func NewNanoID() (string, error) {
	return gonanoid.New()
}

// GenerateNanoID generates a nano ID using the given alphabet and size.
//
// Parameters:
// - alphabet: the characters to be used for generating the ID.
// - size: the length of the ID to be generated.
//
// Returns:
// - string: the generated nano ID.
// - error: an error if the generation process fails.
func GenerateNanoID(alphabet string, size int) (string, error) {
	return gonanoid.Generate(alphabet, size)
}

// NewUUID generates a new UUID.
//
// It does not take any parameters.
// It returns a string.
func NewUUID() string {
	return uuid.New().String()
}
