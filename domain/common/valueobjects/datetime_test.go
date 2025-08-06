package valueobjects_test

import (
	"strings"
	"testing"
	"time"

	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
)

func TestNewDateTime(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "UTC time remains unchanged",
			input:    time.Date(2025, 5, 15, 10, 30, 0, 0, time.UTC),
			expected: "2025-05-15T10:30:00Z",
		},
		{
			name:     "Moscow time converted to UTC (+3)",
			input:    time.Date(2025, 5, 15, 13, 30, 0, 0, time.FixedZone("MSK", 3*3600)),
			expected: "2025-05-15T10:30:00Z",
		},
		{
			name:     "New York time converted to UTC (-4)",
			input:    time.Date(2025, 5, 15, 6, 30, 0, 0, time.FixedZone("EDT", -4*3600)),
			expected: "2025-05-15T10:30:00Z",
		},
		{
			name:     "Zero time handled correctly",
			input:    time.Time{},
			expected: "0001-01-01T00:00:00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt := valueobjects.NewDateTime(tt.input)
			if got := dt.String(); got != tt.expected {
				t.Errorf("NewDateTime().String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestDateTime_Format(t *testing.T) {
	tests := []struct {
		name        string
		layout      string
		timeFunc    func() time.Time
		expected    string
		contains    []string
		notContains []string
	}{
		{
			name:     "RFC3339 format",
			layout:   time.RFC3339,
			timeFunc: func() time.Time { return time.Date(2025, 5, 15, 10, 30, 45, 0, time.UTC) },
			expected: "2025-05-15T10:30:45Z",
		},
		{
			name:     "Custom format",
			layout:   "2006-01-02 15:04:05",
			timeFunc: func() time.Time { return time.Date(2025, 5, 15, 10, 30, 45, 0, time.UTC) },
			expected: "2025-05-15 10:30:45",
		},
		{
			name:     "Date only",
			layout:   "2006-01-02",
			timeFunc: func() time.Time { return time.Date(2025, 5, 15, 10, 30, 45, 0, time.UTC) },
			expected: "2025-05-15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt := valueobjects.NewDateTime(tt.timeFunc())
			got := dt.Format(tt.layout)

			if tt.expected != "" {
				if got != tt.expected {
					t.Errorf("Format(%q) = %q, want %q", tt.layout, got, tt.expected)
				}
			} else {
				for _, substr := range tt.contains {
					if !strings.Contains(got, substr) {
						t.Errorf("Format(\"\") = %q, should contain %q", got, substr)
					}
				}
				for _, substr := range tt.notContains {
					if strings.Contains(got, substr) {
						t.Errorf("Format(\"\") = %q, should NOT contain %q", got, substr)
					}
				}
			}
		})
	}
}

func TestDateTime_Equals(t *testing.T) {
	baseTime := time.Date(2025, 5, 15, 10, 30, 0, 123456789, time.UTC)
	dt1 := valueobjects.NewDateTime(baseTime)
	dt2 := valueobjects.NewDateTime(baseTime)
	dt3 := valueobjects.NewDateTime(baseTime.Add(time.Nanosecond))

	tests := []struct {
		name   string
		dt     valueobjects.DateTime
		other  valueobjects.DateTime
		expect bool
	}{
		{"Equal times", dt1, dt2, true},
		{"Different times (1ns)", dt1, dt3, false},
		{"Same object", dt1, dt1, true},
		{"Zero time comparison", valueobjects.NewDateTime(time.Time{}), valueobjects.NewDateTime(time.Time{}), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dt.Equals(tt.other); got != tt.expect {
				t.Errorf("Equals() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestDateTime_Before_After(t *testing.T) {
	base := valueobjects.NewDateTime(time.Date(2025, 5, 15, 10, 30, 0, 0, time.UTC))
	earlier := valueobjects.NewDateTime(base.Time().Add(-time.Hour))
	later := valueobjects.NewDateTime(base.Time().Add(time.Hour))

	tests := []struct {
		name   string
		dt     valueobjects.DateTime
		other  valueobjects.DateTime
		before bool
		after  bool
	}{
		{"Base before later", base, later, true, false},
		{"Later before base", later, base, false, true},
		{"Base after earlier", base, earlier, false, true},
		{"Earlier after base", earlier, base, true, false},
		{"Equal times", base, base, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dt.Before(tt.other); got != tt.before {
				t.Errorf("Before() = %v, want %v", got, tt.before)
			}
			if got := tt.dt.After(tt.other); got != tt.after {
				t.Errorf("After() = %v, want %v", got, tt.after)
			}
		})
	}
}

func TestDateTime_Add_Sub(t *testing.T) {
	base := valueobjects.NewDateTime(time.Date(2025, 5, 15, 10, 30, 0, 0, time.UTC))
	hour := time.Hour
	twoHours := 2 * time.Hour

	tests := []struct {
		name        string
		operation   func() valueobjects.DateTime
		expectedStr string
	}{
		{
			name: "Add one hour",
			operation: func() valueobjects.DateTime {
				return base.Add(hour)
			},
			expectedStr: "2025-05-15T11:30:00Z",
		},
		{
			name: "Add negative duration (subtract)",
			operation: func() valueobjects.DateTime {
				return base.Add(-hour)
			},
			expectedStr: "2025-05-15T09:30:00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.operation()
			if got := result.String(); got != tt.expectedStr {
				t.Errorf("%s = %q, want %q", tt.name, got, tt.expectedStr)
			}
		})
	}

	// Проверка Sub напрямую
	t.Run("Sub calculates correct duration", func(t *testing.T) {
		later := base.Add(twoHours)
		diff := later.Sub(base)
		if diff != twoHours {
			t.Errorf("Sub() duration = %v, want %v", diff, twoHours)
		}
	})
}

func TestDateTime_IsZero(t *testing.T) {
	tests := []struct {
		name   string
		dt     valueobjects.DateTime
		isZero bool
	}{
		{
			name:   "Zero time",
			dt:     valueobjects.NewDateTime(time.Time{}),
			isZero: true,
		},
		{
			name:   "Non-zero time",
			dt:     valueobjects.NewDateTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
			isZero: false,
		},
		{
			name:   "Now is not zero",
			dt:     valueobjects.Now(),
			isZero: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dt.IsZero(); got != tt.isZero {
				t.Errorf("IsZero() = %v, want %v", got, tt.isZero)
			}
		})
	}
}

func TestDateTime_Immutability(t *testing.T) {
	base := valueobjects.NewDateTime(time.Date(2025, 5, 15, 10, 30, 0, 0, time.UTC))
	originalString := base.String()

	// Выполняем операции, которые НЕ должны изменить исходный объект
	base.Add(time.Hour)
	base.Sub(base)
	base.Format(time.RFC3339)

	if current := base.String(); current != originalString {
		t.Errorf("DateTime should be immutable: changed from %q to %q", originalString, current)
	}
}
