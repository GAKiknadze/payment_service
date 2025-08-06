package valueobjects_test

import (
	"testing"

	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
)

func TestCurrencyIsValid(t *testing.T) {
	tests := []struct {
		name     string
		currency valueobjects.Currency
		expected bool
	}{
		{"RUB uppercase", "RUB", true},
		{"RUB lowercase", "rub", true},
		{"RUB mixed case", "RuB", true},
		{"Invalid currency", "USD", false},
		{"Empty string", "", false},
		{"Whitespace", " RUB ", false},
		{"Partial match", "RUBB", false},
		{"Leading/trailing spaces", "  rub  ", false},
		{"Different currency", "EUR", false},
		{"Special characters", "!@#$%", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.currency.IsValid(); got != tt.expected {
				t.Errorf("IsValid() = %v, want %v for currency %q", got, tt.expected, tt.currency)
			}
		})
	}
}

func TestCurrencyString(t *testing.T) {
	tests := []struct {
		name     string
		currency valueobjects.Currency
		expected string
	}{
		{"RUB constant", valueobjects.RUB, "RUB"},
		{"Lowercase currency", "usd", "usd"},
		{"Empty string", "", ""},
		{"Whitespace", "  ", "  "},
		{"Mixed case", "Rub", "Rub"},
		{"Special characters", "!@#$%", "!@#$%"},
		{"Leading/trailing spaces", "  rub  ", "  rub  "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.currency.String(); got != tt.expected {
				t.Errorf("String() = %q, want %q for currency %q", got, tt.expected, tt.currency)
			}
		})
	}
}
