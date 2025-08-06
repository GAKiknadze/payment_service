package valueobjects_test

import (
	"strings"
	"testing"

	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
	"github.com/stretchr/testify/assert"
)

func TestNewIdempotencyKey(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		// Позитивные тесты
		{
			name:        "Valid UUID in lowercase",
			input:       "123e4567-e89b-42d3-a456-426614174000",
			expected:    "123e4567-e89b-42d3-a456-426614174000",
			expectError: false,
		},
		{
			name:        "Valid UUID in uppercase",
			input:       "123E4567-E89B-42D3-A456-426614174000",
			expected:    "123e4567-e89b-42d3-a456-426614174000",
			expectError: false,
		},
		{
			name:        "Valid UUID with mixed case",
			input:       "123e4567-E89b-42d3-a456-426614174000",
			expected:    "123e4567-e89b-42d3-a456-426614174000",
			expectError: false,
		},
		{
			name:        "Valid UUID version 4",
			input:       "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			expected:    "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			expectError: false,
		},
		{
			name:        "Valid UUID with variant 8",
			input:       "123e4567-e89b-42d3-a456-826614174000",
			expected:    "123e4567-e89b-42d3-a456-826614174000",
			expectError: false,
		},
		{
			name:        "Valid UUID with variant 9",
			input:       "123e4567-e89b-42d3-a456-926614174000",
			expected:    "123e4567-e89b-42d3-a456-926614174000",
			expectError: false,
		},
		{
			name:        "Valid UUID with variant A",
			input:       "123e4567-e89b-42d3-a456-a26614174000",
			expected:    "123e4567-e89b-42d3-a456-a26614174000",
			expectError: false,
		},
		{
			name:        "Valid UUID with variant B",
			input:       "123e4567-e89b-42d3-a456-b26614174000",
			expected:    "123e4567-e89b-42d3-a456-b26614174000",
			expectError: false,
		},

		// Негативные тесты
		{
			name:        "Empty string",
			input:       "",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Too short",
			input:       "123e4567-e89b-42d3-a456-426614174",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Too long",
			input:       "123e4567-e89b-42d3-a456-4266141740000",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid characters",
			input:       "123e4567-e89b-42d3-a456-42661417400g",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid version (not 4)",
			input:       "123e4567-e89b-12d3-a456-426614174000",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid variant (not 8,9,A,B)",
			input:       "123e4567-e89b-42d3-7456-426614174000",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Missing hyphens",
			input:       "123e4567e89b42d3a456426614174000",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Extra hyphens",
			input:       "123e4567--e89b-42d3-a456-426614174000",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Hyphens in wrong positions",
			input:       "123-e4567-e89b-42d3-a456-426614174000",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Non-hex character",
			input:       "123e4567-e89b-42d3-a456-42661417400Z",
			expected:    "",
			expectError: true,
		},
		{
			name:        "UUID with spaces",
			input:       "123e4567-e89b-42d3-a456-426614174000 ",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := valueobjects.NewIdempotencyKey(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, valueobjects.ErrInvalidIdempotencyKey, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, key.String())
			}
		})
	}
}

func TestIdempotencyKey_Validate(t *testing.T) {
	tests := []struct {
		name        string
		key         valueobjects.IdempotencyKey
		expectError bool
	}{
		{
			name:        "Valid key in lowercase",
			key:         valueobjects.IdempotencyKey("123e4567-e89b-42d3-a456-426614174000"),
			expectError: false,
		},
		{
			name:        "Valid UUID version 4",
			key:         valueobjects.IdempotencyKey("f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			expectError: false,
		},
		{
			name:        "Valid UUID with variant 8",
			key:         valueobjects.IdempotencyKey("123e4567-e89b-42d3-a456-826614174000"),
			expectError: false,
		},
		{
			name:        "Invalid key - wrong version",
			key:         valueobjects.IdempotencyKey("123e4567-e89b-12d3-a456-426614174000"),
			expectError: true,
		},
		{
			name:        "Invalid key - wrong variant",
			key:         valueobjects.IdempotencyKey("123e4567-e89b-42d3-7456-426614174000"),
			expectError: true,
		},
		{
			name:        "Empty key",
			key:         valueobjects.IdempotencyKey(""),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.key.Validate()

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, valueobjects.ErrInvalidIdempotencyKey, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIdempotencyKey_Equals(t *testing.T) {
	validKey := valueobjects.IdempotencyKey("123e4567-e89b-42d3-a456-426614174000")
	equalKeyUpper := valueobjects.IdempotencyKey("123E4567-E89B-42D3-A456-426614174000")
	differentKey := valueobjects.IdempotencyKey("f47ac10b-58cc-4372-a567-0e02b2c3d479")

	tests := []struct {
		name     string
		key      valueobjects.IdempotencyKey
		other    valueobjects.IdempotencyKey
		expected bool
	}{
		{
			name:     "Equal keys (same case)",
			key:      validKey,
			other:    validKey,
			expected: true,
		},
		{
			name:     "Equal keys (different case)",
			key:      validKey,
			other:    equalKeyUpper,
			expected: true,
		},
		{
			name:     "Different keys",
			key:      validKey,
			other:    differentKey,
			expected: false,
		},
		{
			name:     "Compare with empty key",
			key:      validKey,
			other:    valueobjects.IdempotencyKey(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.key.Equals(tt.other)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIdempotencyKey_String(t *testing.T) {
	tests := []struct {
		name     string
		key      valueobjects.IdempotencyKey
		expected string
	}{
		{
			name:     "Valid key in lowercase",
			key:      valueobjects.IdempotencyKey("123e4567-e89b-42d3-a456-426614174000"),
			expected: "123e4567-e89b-42d3-a456-426614174000",
		},
		{
			name:     "Valid key in uppercase (stored as is)",
			key:      valueobjects.IdempotencyKey("123E4567-E89B-42D3-A456-426614174000"),
			expected: "123E4567-E89B-42D3-A456-426614174000",
		},
		{
			name:     "Empty key",
			key:      valueobjects.IdempotencyKey(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.key.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateIdempotencyKey(t *testing.T) {
	// Генерируем несколько ключей и проверяем их валидность
	for i := 0; i < 10; i++ {
		key := valueobjects.GenerateIdempotencyKey()

		// Проверяем, что ключ не пустой
		assert.NotEmpty(t, key.String())

		// Проверяем, что ключ проходит валидацию
		err := key.Validate()
		assert.NoError(t, err, "Сгенерированный ключ должен быть валидным")

		// Проверяем, что ключ имеет правильный формат
		keyStr := key.String()
		assert.True(t, strings.Contains(keyStr, "-"), "Ключ должен содержать дефисы")
		assert.Equal(t, 36, len(keyStr), "Длина UUID должна быть 36 символов")

		// Проверяем, что ключ в нижнем регистре (поскольку uuid.NewString() возвращает в нижнем регистре)
		assert.Equal(t, strings.ToLower(keyStr), keyStr, "Ключ должен быть в нижнем регистре")
	}
}

func TestGenerateIdempotencyKey_Uniqueness(t *testing.T) {
	// Генерируем много ключей и проверяем их уникальность
	count := 100
	keys := make(map[string]bool, count)

	for i := 0; i < count; i++ {
		key := valueobjects.GenerateIdempotencyKey()
		keys[key.String()] = true
	}

	// Проверяем, что все ключи уникальны
	assert.Equal(t, count, len(keys), "Все сгенерированные ключи должны быть уникальными")
}
