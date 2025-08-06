package valueobjects_test

import (
	"testing"
	"time"

	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
	"github.com/GAKiknadze/payment_service/domain/common/valueobjects/mocks"
)

func TestDate_Equals(t *testing.T) {
	validDate, _ := valueobjects.NewDate(2025, time.May, 15)
	sameDate, _ := valueobjects.NewDate(2025, time.May, 15)
	differentYear, _ := valueobjects.NewDate(2024, time.May, 15)
	differentMonth, _ := valueobjects.NewDate(2025, time.June, 15)
	differentDay, _ := valueobjects.NewDate(2025, time.May, 16)

	tests := []struct {
		name   string
		date   valueobjects.Date
		other  valueobjects.Date
		expect bool
	}{
		{
			name:   "Equal dates",
			date:   validDate,
			other:  sameDate,
			expect: true,
		},
		{
			name:   "Different year",
			date:   validDate,
			other:  differentYear,
			expect: false,
		},
		{
			name:   "Different month",
			date:   validDate,
			other:  differentMonth,
			expect: false,
		},
		{
			name:   "Different day",
			date:   validDate,
			other:  differentDay,
			expect: false,
		},
		{
			name:   "Same date compared to itself",
			date:   validDate,
			other:  validDate,
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.date.Equals(tt.other); got != tt.expect {
				t.Errorf("Equals() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestToday(t *testing.T) {
	janExpected, _ := valueobjects.NewDate(2025, time.January, 1)
	decExpected, _ := valueobjects.NewDate(2025, time.December, 31)
	leapExpected, _ := valueobjects.NewDate(2020, time.February, 29)
	middleExpected, _ := valueobjects.NewDate(2025, time.May, 15)

	tests := []struct {
		name     string
		now      time.Time
		expected valueobjects.Date
	}{
		{
			name:     "January 1, 2025",
			now:      time.Date(2025, time.January, 1, 12, 30, 0, 0, time.UTC),
			expected: janExpected,
		},
		{
			name:     "December 31, 2025",
			now:      time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC),
			expected: decExpected,
		},
		{
			name:     "Leap day 2020",
			now:      time.Date(2020, time.February, 29, 0, 0, 0, 0, time.UTC),
			expected: leapExpected,
		},
		{
			name:     "Middle of month",
			now:      time.Date(2025, time.May, 15, 8, 45, 30, 0, time.UTC),
			expected: middleExpected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clock := mocks.NewFixedClock(tt.now)
			got := valueobjects.Today(clock)

			if !got.Equals(tt.expected) {
				t.Errorf("Today() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestToday_UsesClockImplementation(t *testing.T) {
	// Проверяем, что Today действительно использует переданный Clock
	// а не системное время
	now := time.Now().UTC()
	future := now.AddDate(0, 0, 10)
	clock := mocks.NewFixedClock(future)

	today := valueobjects.Today(clock)
	expected, _ := valueobjects.NewDate(future.Year(), future.Month(), future.Day())

	if !today.Equals(expected) {
		t.Errorf("Today() = %v, want %v", today, expected)
	}
}
