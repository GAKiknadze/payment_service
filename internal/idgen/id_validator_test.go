package idgen_test

import (
	"testing"

	"github.com/GAKiknadze/payment_service/internal/idgen"
)

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		// Позитивные тесты
		{
			name:  "Valid UUID in lowercase",
			input: "123e4567-e89b-42d3-a456-426614174000",
			valid: true,
		},
		{
			name:  "Valid UUID in uppercase",
			input: "123E4567-E89B-42D3-A456-426614174000",
			valid: true,
		},
		{
			name:  "Valid UUID with mixed case",
			input: "123e4567-E89b-42d3-a456-426614174000",
			valid: true,
		},
		{
			name:  "Valid UUID version 4",
			input: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			valid: true,
		},
		{
			name:  "Valid UUID with variant 8",
			input: "123e4567-e89b-42d3-a456-826614174000",
			valid: true,
		},
		{
			name:  "Valid UUID with variant 9",
			input: "123e4567-e89b-42d3-a456-926614174000",
			valid: true,
		},
		{
			name:  "Valid UUID with variant A",
			input: "123e4567-e89b-42d3-a456-a26614174000",
			valid: true,
		},
		{
			name:  "Valid UUID with variant B",
			input: "123e4567-e89b-42d3-a456-b26614174000",
			valid: true,
		},

		// Негативные тесты
		{
			name:  "Empty string",
			input: "",
			valid: false,
		},
		{
			name:  "Too short",
			input: "123e4567-e89b-42d3-a456-426614174",
			valid: false,
		},
		{
			name:  "Too long",
			input: "123e4567-e89b-42d3-a456-4266141740000",
			valid: false,
		},
		{
			name:  "Invalid characters",
			input: "123e4567-e89b-42d3-a456-42661417400g",
			valid: false,
		},
		{
			name:  "Invalid version (not 4) - third group starts with 1",
			input: "123e4567-e89b-12d3-a456-426614174000",
			valid: false,
		},
		{
			name:  "Invalid version (not 4) - third group starts with 5",
			input: "123e4567-e89b-52d3-a456-426614174000",
			valid: false,
		},
		{
			name:  "Invalid variant (not 8,9,A,B) - fourth group starts with 7",
			input: "123e4567-e89b-42d3-7456-426614174000",
			valid: false,
		},
		{
			name:  "Invalid variant (not 8,9,A,B) - fourth group starts with c",
			input: "123e4567-e89b-42d3-c456-426614174000",
			valid: false,
		},
		{
			name:  "Missing hyphens",
			input: "123e4567e89b42d3a456426614174000",
			valid: false,
		},
		{
			name:  "Extra hyphens",
			input: "123e4567--e89b-42d3-a456-426614174000",
			valid: false,
		},
		{
			name:  "Hyphens in wrong positions",
			input: "123-e4567-e89b-42d3-a456-426614174000",
			valid: false,
		},
		{
			name:  "Non-hex character",
			input: "123e4567-e89b-42d3-a456-42661417400Z",
			valid: false,
		},
		{
			name:  "UUID with spaces",
			input: "123e4567-e89b-42d3-a456-426614174000 ",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := idgen.ValidateUUID(tt.input); got != tt.valid {
				t.Errorf("ValidateUUID(%q) = %v, want %v", tt.input, got, tt.valid)
			}
		})
	}
}

func TestValidateShortID(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		length int
		valid  bool
	}{
		// Позитивные тесты
		{
			name:   "Valid short ID of length 1",
			input:  "a",
			length: 1,
			valid:  true,
		},
		{
			name:   "Valid short ID of length 5",
			input:  "abc12",
			length: 5,
			valid:  true,
		},
		{
			name:   "Valid short ID with uppercase",
			input:  "ABC12",
			length: 5,
			valid:  true,
		},
		{
			name:   "Valid short ID with max length",
			input:  "abcdefghijklmnopqrstuvwxyz1234567890",
			length: 36,
			valid:  true,
		},
		{
			name:   "Valid short ID with mixed case",
			input:  "AbC12",
			length: 5,
			valid:  true,
		},
		{
			name:   "Only numbers",
			input:  "12345",
			length: 5,
			valid:  true,
		},
		{
			name:   "Only letters",
			input:  "abcde",
			length: 5,
			valid:  true,
		},

		// Негативные тесты
		{
			name:   "Empty string",
			input:  "",
			length: 1,
			valid:  false,
		},
		{
			name:   "Too short",
			input:  "ab",
			length: 3,
			valid:  false,
		},
		{
			name:   "Too long",
			input:  "abc",
			length: 2,
			valid:  false,
		},
		{
			name:   "Contains space",
			input:  "ab c",
			length: 4,
			valid:  false,
		},
		{
			name:   "Contains special character",
			input:  "ab#c",
			length: 4,
			valid:  false,
		},
		{
			name:   "Contains underscore",
			input:  "ab_c",
			length: 4,
			valid:  false,
		},
		{
			name:   "Contains hyphen",
			input:  "ab-c",
			length: 4,
			valid:  false,
		},
		{
			name:   "Contains dot",
			input:  "ab.c",
			length: 4,
			valid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := idgen.ValidateShortID(tt.input, tt.length); got != tt.valid {
				t.Errorf("ValidateShortID(%q, %d) = %v, want %v",
					tt.input, tt.length, got, tt.valid)
			}
		})
	}
}

func TestValidatePrefixedID(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		prefix string
		length int
		valid  bool
	}{
		// Позитивные тесты
		{
			name:   "Valid prefixed ID with lowercase prefix",
			id:     "user-abc123",
			prefix: "user",
			length: 6,
			valid:  true,
		},
		{
			name:   "Valid prefixed ID with uppercase prefix",
			id:     "USER-abc123",
			prefix: "user",
			length: 6,
			valid:  true,
		},
		{
			name:   "Valid prefixed ID with mixed case prefix",
			id:     "UsEr-abc123",
			prefix: "user",
			length: 6,
			valid:  true,
		},
		{
			name:   "Valid prefixed ID with different prefix",
			id:     "order-123abc",
			prefix: "order",
			length: 6,
			valid:  true,
		},
		{
			name:   "Valid prefixed ID with prefix of different case",
			id:     "USER-abc123",
			prefix: "USER",
			length: 6,
			valid:  true,
		},
		{
			name:   "Valid prefixed ID with minimal length",
			id:     "user-a",
			prefix: "user",
			length: 1,
			valid:  true,
		},

		// Негативные тесты
		{
			name:   "Empty string",
			id:     "",
			prefix: "user",
			length: 6,
			valid:  false,
		},
		{
			name:   "No hyphen",
			id:     "userabc123",
			prefix: "user",
			length: 6,
			valid:  false,
		},
		{
			name:   "Multiple hyphens",
			id:     "user-abc-123",
			prefix: "user",
			length: 6,
			valid:  false,
		},
		{
			name:   "Wrong prefix",
			id:     "order-abc123",
			prefix: "user",
			length: 6,
			valid:  false,
		},
		{
			name:   "Wrong prefix case-sensitive check (should pass as EqualFold is used)",
			id:     "USER-abc123",
			prefix: "user",
			length: 6,
			valid:  true,
		},
		{
			name:   "Invalid short ID part",
			id:     "user-ab#c",
			prefix: "user",
			length: 4,
			valid:  false,
		},
		{
			name:   "Short ID part too short",
			id:     "user-abc",
			prefix: "user",
			length: 4,
			valid:  false,
		},
		{
			name:   "Short ID part too long",
			id:     "user-abc123",
			prefix: "user",
			length: 5,
			valid:  false,
		},
		{
			name:   "Hyphen at start",
			id:     "-user-abc123",
			prefix: "user",
			length: 6,
			valid:  false,
		},
		{
			name:   "Hyphen at end",
			id:     "user-",
			prefix: "user",
			length: 0,
			valid:  false,
		},
		{
			name:   "Only hyphen",
			id:     "-",
			prefix: "",
			length: 0,
			valid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := idgen.ValidatePrefixedID(tt.id, tt.prefix, tt.length); got != tt.valid {
				t.Errorf("ValidatePrefixedID(%q, %q, %d) = %v, want %v",
					tt.id, tt.prefix, tt.length, got, tt.valid)
			}
		})
	}
}
