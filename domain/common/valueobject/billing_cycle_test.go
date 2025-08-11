package valueobject_test

import (
	"testing"
	"time"

	"github.com/GAKiknadze/payment_service/domain/common/valueobject"
)

func TestCalculateNextBillingDate_Hourly(t *testing.T) {
	cycle, _ := valueobject.NewBillingCycle(valueobject.BillingCycleHourly)
	now := time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC)

	next := mustCalculateNextDate(cycle, now)

	expected := now.Add(time.Hour)
	if !next.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, next)
	}
}

func TestCalculateNextBillingDate_Monthly(t *testing.T) {
	cases := []struct {
		name        string
		currentDate time.Time
		expected    time.Time
	}{
		{
			"middle of month",
			time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC),
			time.Date(2023, 2, 15, 10, 30, 0, 0, time.UTC),
		},
		{
			"end of January",
			time.Date(2023, 1, 31, 10, 30, 0, 0, time.UTC),
			time.Date(2023, 2, 28, 10, 30, 0, 0, time.UTC),
		},
		{
			"leap year February 29",
			time.Date(2020, 2, 29, 10, 30, 0, 0, time.UTC),
			time.Date(2020, 3, 29, 10, 30, 0, 0, time.UTC),
		},
		{
			"non-leap year after February 29",
			time.Date(2021, 2, 28, 10, 30, 0, 0, time.UTC),
			time.Date(2021, 3, 28, 10, 30, 0, 0, time.UTC),
		},
	}

	cycle, _ := valueobject.NewBillingCycle(valueobject.BillingCycleMonthly)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			next := mustCalculateNextDate(cycle, tc.currentDate)
			if !next.Equal(tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, next)
			}
		})
	}
}

func TestCalculateNextBillingDate_OneTime(t *testing.T) {
	cycle, _ := valueobject.NewBillingCycle(valueobject.BillingCycleOneTime)
	next, err := cycle.CalculateNextBillingDate(time.Now())

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !next.IsZero() {
		t.Errorf("Expected zero date for OneTime cycle")
	}
}

// Вспомогательная функция для упрощения тестов
func mustCalculateNextDate(cycle valueobject.BillingCycle, date time.Time) time.Time {
	next, err := cycle.CalculateNextBillingDate(date)
	if err != nil {
		panic("Failed to calculate next billing date: " + err.Error())
	}
	return next
}
