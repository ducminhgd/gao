package security

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	DefaultIteration = 10_000
	MinIteration     = 1_000
	MaxIteration     = 1_000_000

	DefaultKeyLength = 32
	MinKeyLength     = 8
	MaxKeyLength     = 1_024
)

// ParseIteration returns the given iteration if it is within the valid range.
// Otherwise, it returns the default iteration.
func ParseIteration(i int) int {
	if i < MinIteration || i > MaxIteration {
		return DefaultIteration
	}
	return i
}

// ParseKeyLength returns the given key length if it is within the valid range.
// Otherwise, it returns the default key length.
func ParseKeyLength(k int) int {
	if k < MinKeyLength || k > MaxKeyLength {
		return DefaultKeyLength
	}
	return k
}

// GetHashFunc returns the given hash algorithm and a function to create a new instance of the hash algorithm.
//
// The given hash algorithm is case-insensitive. If the given hash algorithm is not supported, it defaults to "sha512".
//
// Supported hash algorithms are:
// - md5
// - sha1
// - sha256
// - sha512
func GetHashFunc(hashAlgo string) (string, func() hash.Hash) {
	switch strings.ToLower(hashAlgo) {
	case "md5":
		return "md5", md5.New
	case "sha1":
		return "sha1", sha1.New
	case "sha256":
		return "sha256", sha256.New
	case "sha512":
		return "sha512", sha512.New
	default:
		return "sha512", sha512.New
	}
}

// GeneratePBKDF2 generates a PBKDF2 key based on the given password, salt, iteration and key length.
//
// The iteration and key length are parsed and validated before being used.
// The hash algorithm is case-insensitive and defaults to "sha512" if not supported.
// Supported hash algorithms are "md5", "sha1", "sha256", and "sha512".
//
// The output string is in the format of "hashAlgo$salt$iteration$keyLength$hashedPassword".
func GeneratePBKDF2(password, salt string, iteration, keyLength int, hashAlgo string) string {
	algo, hashFunc := GetHashFunc(hashAlgo)
	i := ParseIteration(iteration)
	k := ParseKeyLength(keyLength)
	hashedStr := pbkdf2.Key([]byte(password), []byte(salt), i, k, hashFunc)

	return fmt.Sprintf("pbkdf2_%s$%s$%d$%d$%x", algo, salt, i, k, hashedStr)
}

// ValidatePBKDF2 validates the given hash string against the given password.
//
// The given hash string is expected to be in the format of "pbkdf2_<hashAlgo>$<salt>$<iteration>$<keyLength>$<hashedPassword>".
// The given password is hashed using the given hash string's algorithm, salt, iteration, and key length.
// If the generated hash string matches the given hash string, the function returns true. Otherwise, it returns false.
func ValidatePBKDF2(password, hash string) bool {
	if hash == "" {
		return false
	}
	parts := strings.Split(hash, "$")

	// validate format
	if len(parts) != 5 {
		return false
	}
	if !strings.HasPrefix(parts[0], "pbkdf2_") {
		return false
	}

	// validate algo
	algo := strings.TrimPrefix(parts[0], "pbkdf2_")
	if algo != "md5" && algo != "sha1" && algo != "sha256" && algo != "sha512" {
		return false
	}

	// validate iteration
	i, err := strconv.Atoi(parts[2])
	if err != nil {
		return false
	}
	if i < MinIteration || i > MaxIteration {
		return false
	}

	// validate key length
	k, err := strconv.Atoi(parts[3])
	if err != nil {
		return false
	}
	if k < MinKeyLength || k > MaxKeyLength {
		return false
	}

	// validate salt and hashed string
	if parts[1] == "" || parts[4] == "" {
		return false
	}

	pw := GeneratePBKDF2(password, parts[1], i, k, algo)
	return pw == hash
}
