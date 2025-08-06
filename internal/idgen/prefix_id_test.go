package idgen_test

import (
	"strings"
	"testing"

	"github.com/GAKiknadze/payment_service/internal/idgen"
	"github.com/stretchr/testify/assert"
)

// Вспомогательная функция для проверки корректности короткого ID
func isValidShortID(id string) bool {
	if id == "" {
		return false
	}

	for _, r := range id {
		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}

func TestGeneratePrefixedID(t *testing.T) {
	tests := []struct {
		name        string
		prefix      string
		length      int
		expectedLen int
	}{
		{
			name:        "Standard prefix with uppercase",
			prefix:      "TAR",
			length:      8,
			expectedLen: 12, // 3 (TAR) + 1 (-) + 8
		},
		{
			name:        "Prefix with lowercase letters",
			prefix:      "tar",
			length:      6,
			expectedLen: 10, // 3 (TAR) + 1 (-) + 6
		},
		{
			name:        "Prefix with numbers",
			prefix:      "TAR123",
			length:      10,
			expectedLen: 17, // 6 (TAR123) + 1 (-) + 10
		},
		{
			name:        "Prefix with special characters",
			prefix:      "T@R!",
			length:      5,
			expectedLen: 8, // 2 (TR) + 1 (-) + 5
		},
		{
			name:        "Prefix with mixed case and numbers",
			prefix:      "Us3r",
			length:      7,
			expectedLen: 12, // 4 (US3R) + 1 (-) + 7
		},
		{
			name:        "Empty prefix",
			prefix:      "",
			length:      4,
			expectedLen: 7, // 2 (ID) + 1 (-) + 4
		},
		{
			name:        "Prefix with only special characters",
			prefix:      "!@#$",
			length:      12,
			expectedLen: 15, // 2 (ID) + 1 (-) + 12
		},
		{
			name:        "Prefix with Unicode characters",
			prefix:      "Тест",
			length:      3,
			expectedLen: 6, // 2 (ID) + 1 (-) + 3
		},
		{
			name:        "Long prefix",
			prefix:      "PAYMENTSERVICE",
			length:      15,
			expectedLen: 30, // 14 (PAYMENTSERVICE) + 1 (-) + 15
		},
		{
			name:        "Minimum length",
			prefix:      "ID",
			length:      1,
			expectedLen: 4, // 2 (ID) + 1 (-) + 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := idgen.GeneratePrefixedID(tt.prefix, tt.length)

			// Проверяем длину результата
			assert.Equal(t, tt.expectedLen, len(result),
				"Длина должна быть %d, получено %d", tt.expectedLen, len(result))

			// Проверяем структуру: префикс + '-' + shortID
			parts := strings.Split(result, "-")
			assert.Equal(t, 2, len(parts),
				"Должны быть ровно 2 части, разделенные дефисом")

			// Проверяем префикс
			expectedPrefix := tt.prefix
			if expectedPrefix == "" || !containsValidChars(expectedPrefix) {
				expectedPrefix = "ID"
			} else {
				expectedPrefix = cleanPrefix(expectedPrefix)
			}

			assert.Equal(t, strings.ToUpper(expectedPrefix), parts[0],
				"Префикс должен быть %q, получен %q", strings.ToUpper(expectedPrefix), parts[0])

			// Проверяем длину короткой части
			assert.Equal(t, tt.length, len(parts[1]),
				"Длина короткой части должна быть %d, получено %d", tt.length, len(parts[1]))

			// Проверяем, что короткая часть содержит только допустимые символы
			assert.True(t, isValidShortID(parts[1]),
				"Короткая часть содержит недопустимые символы: %q", parts[1])
		})
	}
}

// Вспомогательная функция для очистки префикса (как в оригинальном коде)
func cleanPrefix(prefix string) string {
	clean := strings.Map(func(r rune) rune {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, strings.ToUpper(prefix))

	if clean == "" {
		return "ID"
	}
	return clean
}

// Вспомогательная функция для проверки наличия допустимых символов в префиксе
func containsValidChars(prefix string) bool {
	for _, r := range prefix {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return true
		}
	}
	return false
}

func TestGeneratePrefixedID_EdgeCases(t *testing.T) {
	// Проверка с нулевой длиной
	t.Run("Zero length", func(t *testing.T) {
		result := idgen.GeneratePrefixedID("TAR", 0)
		assert.Equal(t, "TAR-", result)
		assert.Equal(t, 4, len(result))
	})

	// Проверка с отрицательной длиной
	t.Run("Negative length", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative length")
			} else if msg, ok := r.(string); !ok || !strings.Contains(msg, "negative length") {
				t.Errorf("unexpected panic message: %v", r)
			}
		}()
		idgen.GeneratePrefixedID("TAR", -5)
	})

	// Проверка с очень длинным префиксом
	t.Run("Very long prefix", func(t *testing.T) {
		longPrefix := strings.Repeat("A", 100)
		result := idgen.GeneratePrefixedID(longPrefix, 8)

		parts := strings.Split(result, "-")
		assert.Equal(t, 100, len(parts[0]))
		assert.Equal(t, 8, len(parts[1]))
	})

	// Проверка с максимальной длиной
	t.Run("Maximum length", func(t *testing.T) {
		result := idgen.GeneratePrefixedID("TAR", 1000)
		parts := strings.Split(result, "-")
		assert.Equal(t, 1000, len(parts[1]))
	})
}
