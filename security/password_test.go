package security

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/pbkdf2"
)

func TestParseIteration(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "valid iteration",
			input:    10000,
			expected: 10000,
		},
		{
			name:     "iteration below min",
			input:    999,
			expected: DefaultIteration,
		},
		{
			name:     "iteration above max",
			input:    1000001,
			expected: DefaultIteration,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ParseIteration(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestParseKeyLength(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "valid key length",
			input:    32,
			expected: 32,
		},
		{
			name:     "key length below min",
			input:    7,
			expected: DefaultKeyLength,
		},
		{
			name:     "key length above max",
			input:    1025,
			expected: DefaultKeyLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ParseKeyLength(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestGetHashFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hashFunc func() hash.Hash
	}{
		{
			name:     "md5",
			input:    "md5",
			expected: "md5",
			hashFunc: md5.New,
		},
		{
			name:     "sha1",
			input:    "sha1",
			expected: "sha1",
			hashFunc: sha1.New,
		},
		{
			name:     "sha256",
			input:    "sha256",
			expected: "sha256",
			hashFunc: sha256.New,
		},
		{
			name:     "sha512",
			input:    "sha512",
			expected: "sha512",
			hashFunc: sha512.New,
		},
		{
			name:     "unknown",
			input:    "unknown",
			expected: "sha512",
			hashFunc: sha512.New,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualAlgo, actualHashFunc := GetHashFunc(tt.input)
			assert.Equal(t, tt.expected, actualAlgo)
			assert.NotNil(t, actualHashFunc)
			assert.Implements(t, (*hash.Hash)(nil), actualHashFunc())
		})
	}
}
func TestGeneratePBKDF2(t *testing.T) {
	type Given struct {
		password  string
		salt      string
		iteration int
		keyLength int
		hashAlgo  string
	}
	type Want struct {
		iteration int
		keyLength int
		hashAlgo  string
	}
	tests := []struct {
		name     string
		input    Given
		expected Want
	}{
		{
			name:     "valid parameters",
			input:    Given{password: "password", salt: "salt", iteration: 10_000, keyLength: 32, hashAlgo: "sha512"},
			expected: Want{iteration: 10_000, keyLength: 32, hashAlgo: "sha512"},
		},
		{
			name:     "invalid iteration and keylength",
			input:    Given{password: "password", salt: "salt", iteration: 0, keyLength: 0, hashAlgo: "sha512"},
			expected: Want{iteration: DefaultIteration, keyLength: DefaultKeyLength, hashAlgo: "sha512"},
		},
		{
			name:     "invalid hash algorithm",
			input:    Given{password: "password", salt: "salt", iteration: 10_000, keyLength: 32, hashAlgo: "unknown"},
			expected: Want{iteration: 10_000, keyLength: 32, hashAlgo: "sha512"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashAlgo, hashFunc := GetHashFunc(tt.input.hashAlgo)
			expectedHash := pbkdf2.Key([]byte(tt.input.password), []byte(tt.input.salt), ParseIteration(tt.input.iteration), ParseKeyLength(tt.input.keyLength), hashFunc)
			expectedHashStr := fmt.Sprintf("pbkdf2_%s$%s$%d$%d$%x", hashAlgo, tt.input.salt, ParseIteration(tt.input.iteration), ParseKeyLength(tt.input.keyLength), expectedHash)

			actual := GeneratePBKDF2(tt.input.password, tt.input.salt, tt.input.iteration, tt.input.keyLength, tt.input.hashAlgo)
			assert.Equal(t, expectedHashStr, actual)
		})
	}
}

func TestValidatePBKDF2(t *testing.T) {
	type Given struct {
		password string
		hash     string
	}
	tests := []struct {
		name     string
		input    Given
		expected bool
	}{
		{
			name:     "valid password and hash",
			input:    Given{password: "password", hash: "pbkdf2_sha512$salt$10000$32$72629a41b076e588fba8c71ca37fadc9acdc8e7321b9cb4ea55fd0bf9fe8ed72"},
			expected: true,
		},
		{
			name:     "empty password",
			input:    Given{password: "", hash: "pbkdf2_sha512$salt$10000$32$926889d936bb7c21a92599b03d648a5c05c82e202aa5ba8d421196df05509d47"},
			expected: true,
		},
		{
			name:     "invalid password",
			input:    Given{password: "wrong password", hash: "pbkdf2_sha512$salt$10000$32$9a8c1f3f1cf8f1f7f1f7f1f7f1f7f1f1f"},
			expected: false,
		},
		{
			name:     "invalid format - no prefix",
			input:    Given{password: "password", hash: "sha512$salt$10000$32$72629a41b076e588fba8c71ca37fadc9acdc8e7321b9cb4ea55fd0bf9fe8ed72"},
			expected: false,
		},
		{
			name:     "invalid format - wrong prefix",
			input:    Given{password: "password", hash: "pbkdf3_sha512$salt$10000$32$72629a41b076e588fba8c71ca37fadc9acdc8e7321b9cb4ea55fd0bf9fe8ed72"},
			expected: false,
		},
		{
			name:     "invalid format - wrong algorithm",
			input:    Given{password: "password", hash: "pbkdf2_unknown$salt$10000$32$72629a41b076e588fba8c71ca37fadc9acdc8e7321b9cb4ea55fd0bf9fe8ed72"},
			expected: false,
		},
		{
			name:     "invalid format - empty salt",
			input:    Given{password: "password", hash: "pbkdf2_sha512$$10000$32$72629a41b076e588fba8c71ca37fadc9acdc8e7321b9cb4ea55fd0bf9fe8ed72"},
			expected: false,
		},
		{
			name:     "invalid format - empty hash",
			input:    Given{password: "password", hash: "pbkdf2_sha512$salt$10000$32$"},
			expected: false,
		},
		{
			name:     "invalid format - wrong key length",
			input:    Given{password: "password", hash: "pbkdf2_sha512$salt$10000$0$72629a41b076e588fba8c71ca37fadc9acdc8e7321b9cb4ea55fd0bf9fe8ed72"},
			expected: false,
		},
		{
			name:     "invalid format - wrong key iteration",
			input:    Given{password: "password", hash: "pbkdf2_sha512$salt$10001$32$72629a41b076e588fba8c71ca37fadc9acdc8e7321b9cb4ea55fd0bf9fe8ed72"},
			expected: false,
		},
		{
			name:     "invalid format - empty string",
			input:    Given{password: "password", hash: ""},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ValidatePBKDF2(tt.input.password, tt.input.hash)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
