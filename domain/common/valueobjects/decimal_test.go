package valueobjects_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
	"github.com/shopspring/decimal"
)

func TestNewDecimal(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{
			name:     "Whole number",
			input:    100.0,
			expected: "100.00",
		},
		{
			name:     "Two decimal places",
			input:    123.45,
			expected: "123.45",
		},
		{
			name:     "Three decimal places - rounding down",
			input:    123.454,
			expected: "123.45",
		},
		{
			name:     "Three decimal places - rounding up",
			input:    123.455,
			expected: "123.46",
		},
		{
			name:     "Negative number",
			input:    -50.75,
			expected: "-50.75",
		},
		{
			name:     "Zero",
			input:    0.0,
			expected: "0.00",
		},
		{
			name:     "Small fraction",
			input:    0.001,
			expected: "0.00",
		},
		{
			name:     "Small fraction rounding up",
			input:    0.006,
			expected: "0.01",
		},
		{
			name:     "Large number",
			input:    123456789.12,
			expected: "123456789.12",
		},
		{
			name:     "Very small number",
			input:    0.0001,
			expected: "0.00",
		},
		{
			name:     "Precision test with 0.1",
			input:    0.1,
			expected: "0.10",
		},
		{
			name:     "Precision test with 0.2",
			input:    0.2,
			expected: "0.20",
		},
		{
			name:     "Precision test with 0.3",
			input:    0.3,
			expected: "0.30",
		},
		{
			name:     "Floating point representation error",
			input:    0.1 + 0.2,
			expected: "0.30",
		},
		{
			name:     "Max float64 value",
			input:    math.MaxFloat64,
			expected: "1.7976931348623157e+308",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := valueobjects.NewDecimal(tt.input)
			expected := decimal.RequireFromString(tt.expected)

			if !result.Equal(expected) {
				t.Errorf("NewDecimal(%v) = %s, want %s",
					tt.input,
					result.String(),
					expected.String())
			}
		})
	}
}

func TestNewDecimal_RoundingRules(t *testing.T) {
	// Проверяем правила округления при разных значениях третьего знака
	tests := []struct {
		input    float64
		expected string
	}{
		{1.234, "1.23"},
		{1.235, "1.24"},
		{1.236, "1.24"},
		{2.784, "2.78"},
		{2.785, "2.78"},
		{2.786, "2.79"},
		{0.004, "0.00"},
		{0.005, "0.00"},
		{0.006, "0.01"},
		{99.994, "99.99"},
		{99.995, "100.00"},
		{99.996, "100.00"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Rounding case #%d", i+1), func(t *testing.T) {
			result := valueobjects.NewDecimal(tt.input)

			expected := decimal.RequireFromString(tt.expected)

			if !result.Equal(expected) {
				t.Errorf("NewDecimal(%v) = %s, want %s", tt.input, result.String(), expected.String())
			}
		})
	}
}

func TestNewDecimalFromInt(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{
			name:     "Positive integer",
			input:    100,
			expected: "100",
		},
		{
			name:     "Negative integer",
			input:    -50,
			expected: "-50",
		},
		{
			name:     "Zero",
			input:    0,
			expected: "0",
		},
		{
			name:     "Large positive integer",
			input:    123456789,
			expected: "123456789",
		},
		{
			name:     "Large negative integer",
			input:    -987654321,
			expected: "-987654321",
		},
		{
			name:     "Max int64",
			input:    math.MaxInt64,
			expected: "9223372036854775807",
		},
		{
			name:     "Min int64",
			input:    math.MinInt64,
			expected: "-9223372036854775808",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := valueobjects.NewDecimalFromInt(tt.input)
			if result.String() != tt.expected {
				t.Errorf("NewDecimalFromInt(%d) = %s, want %s", tt.input, result.String(), tt.expected)
			}
		})
	}
}
