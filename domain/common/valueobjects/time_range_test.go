package valueobjects_test

import (
	"testing"
	"time"

	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
	"github.com/stretchr/testify/assert"
)

func TestNewTimeRange(t *testing.T) {
	now := time.Now()
	future := now.Add(1 * time.Hour)
	past := now.Add(-1 * time.Hour)
	sameTime := now

	tests := []struct {
		name        string
		start       time.Time
		end         time.Time
		expectError bool
	}{
		{
			name:        "Valid time range",
			start:       past,
			end:         future,
			expectError: false,
		},
		{
			name:        "Start equals end",
			start:       now,
			end:         sameTime,
			expectError: true,
		},
		{
			name:        "Start after end",
			start:       future,
			end:         past,
			expectError: true,
		},
		{
			name:        "Very short interval",
			start:       now,
			end:         now.Add(1 * time.Nanosecond),
			expectError: false,
		},
		{
			name:        "Long interval",
			start:       now.Add(-365 * 24 * time.Hour),
			end:         now.Add(365 * 24 * time.Hour),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := valueobjects.NewTimeRange(tt.start, tt.end)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, valueobjects.ErrInvalidTimeRange, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTimeRange_Start_End(t *testing.T) {
	now := time.Now()
	future := now.Add(1 * time.Hour)

	timeRange, err := valueobjects.NewTimeRange(now, future)
	assert.NoError(t, err)

	assert.Equal(t, valueobjects.NewDateTime(now).String(), timeRange.Start().String())
	assert.Equal(t, valueobjects.NewDateTime(future).String(), timeRange.End().String())
}

func TestTimeRange_Duration(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		start    time.Time
		end      time.Time
		expected time.Duration
	}{
		{
			name:     "One hour interval",
			start:    now,
			end:      now.Add(1 * time.Hour),
			expected: 1 * time.Hour,
		},
		{
			name:     "One minute interval",
			start:    now,
			end:      now.Add(1 * time.Minute),
			expected: 1 * time.Minute,
		},
		{
			name:     "One day interval",
			start:    now,
			end:      now.Add(24 * time.Hour),
			expected: 24 * time.Hour,
		},
		{
			name:     "Very short interval",
			start:    now,
			end:      now.Add(500 * time.Millisecond),
			expected: 500 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeRange, err := valueobjects.NewTimeRange(tt.start, tt.end)
			assert.NoError(t, err)

			// Допускаем небольшую погрешность из-за преобразования в UTC
			duration := timeRange.Duration()
			assert.InDelta(t, float64(tt.expected), float64(duration), float64(1*time.Millisecond))
		})
	}
}

func TestTimeRange_Contains(t *testing.T) {
	now := time.Now()
	start := now.Add(-1 * time.Hour)
	end := now.Add(1 * time.Hour)

	timeRange, err := valueobjects.NewTimeRange(start, end)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		moment   time.Time
		contains bool
	}{
		{
			name:     "Moment in the middle",
			moment:   now,
			contains: true,
		},
		{
			name:     "Moment at start",
			moment:   start,
			contains: true,
		},
		{
			name:     "Moment at end",
			moment:   end,
			contains: true,
		},
		{
			name:     "Moment before start",
			moment:   start.Add(-1 * time.Minute),
			contains: false,
		},
		{
			name:     "Moment after end",
			moment:   end.Add(1 * time.Minute),
			contains: false,
		},
		{
			name:     "Moment very close to start",
			moment:   start.Add(1 * time.Nanosecond),
			contains: true,
		},
		{
			name:     "Moment very close to end",
			moment:   end.Add(-1 * time.Nanosecond),
			contains: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contains := timeRange.Contains(valueobjects.NewDateTime(tt.moment))
			assert.Equal(t, tt.contains, contains)
		})
	}
}

func TestTimeRange_Equals(t *testing.T) {
	now := time.Now()
	future := now.Add(1 * time.Hour)
	past := now.Add(-1 * time.Hour)

	timeRange1, err := valueobjects.NewTimeRange(past, now)
	assert.NoError(t, err)

	timeRange2, err := valueobjects.NewTimeRange(past, now)
	assert.NoError(t, err)

	timeRange3, err := valueobjects.NewTimeRange(past, future)
	assert.NoError(t, err)

	timeRange4, err := valueobjects.NewTimeRange(now, future)
	assert.NoError(t, err)

	tests := []struct {
		name  string
		tr1   valueobjects.TimeRange
		tr2   valueobjects.TimeRange
		equal bool
	}{
		{
			name:  "Equal time ranges",
			tr1:   timeRange1,
			tr2:   timeRange2,
			equal: true,
		},
		{
			name:  "Different start times",
			tr1:   timeRange1,
			tr2:   timeRange4,
			equal: false,
		},
		{
			name:  "Different end times",
			tr1:   timeRange1,
			tr2:   timeRange3,
			equal: false,
		},
		{
			name:  "Same object",
			tr1:   timeRange1,
			tr2:   timeRange1,
			equal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.equal, tt.tr1.Equals(tt.tr2))
		})
	}
}

func TestTimeRange_UTCNormalization(t *testing.T) {
	// Создаем время в разных временных зонах
	moscow, _ := time.LoadLocation("Europe/Moscow")
	ny, _ := time.LoadLocation("America/New_York")

	moscowTime := time.Date(2025, 5, 15, 10, 30, 0, 0, moscow)
	nyTime := time.Date(2025, 5, 15, 6, 30, 0, 0, ny) // То же время в UTC

	// Создаем TimeRange с временем в разных временных зонах
	timeRange, err := valueobjects.NewTimeRange(moscowTime, nyTime.Add(1*time.Hour))
	assert.NoError(t, err)

	// Проверяем, что время нормализовано к UTC
	startUTC := moscowTime.UTC()
	endUTC := nyTime.Add(1 * time.Hour).UTC()

	assert.Equal(t, startUTC.Format(time.RFC3339), timeRange.Start().String())
	assert.Equal(t, endUTC.Format(time.RFC3339), timeRange.End().String())
}

func TestTimeRange_BorderCases(t *testing.T) {
	// Тест с нулевым временем
	zeroTime := time.Time{}
	future := time.Now().Add(1 * time.Hour)

	t.Run("Start is zero time", func(t *testing.T) {
		_, err := valueobjects.NewTimeRange(zeroTime, future)
		assert.NoError(t, err)
	})

	t.Run("End is zero time", func(t *testing.T) {
		_, err := valueobjects.NewTimeRange(future, zeroTime)
		assert.Error(t, err)
		assert.Equal(t, valueobjects.ErrInvalidTimeRange, err)
	})

	// Тест с максимальным временем
	maxTime := time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)

	t.Run("Very far future", func(t *testing.T) {
		_, err := valueobjects.NewTimeRange(time.Now(), maxTime)
		assert.NoError(t, err)
	})

	t.Run("Max time range", func(t *testing.T) {
		_, err := valueobjects.NewTimeRange(zeroTime, maxTime)
		assert.NoError(t, err)
	})
}

func TestTimeRange_Immutability(t *testing.T) {
	now := time.Now()
	future := now.Add(1 * time.Hour)

	timeRange, err := valueobjects.NewTimeRange(now, future)
	assert.NoError(t, err)

	originalStart := timeRange.Start().String()
	originalEnd := timeRange.End().String()

	// Проверяем, что операции с DateTime не изменяют исходный TimeRange
	timeRange.Start().Add(1 * time.Hour)
	timeRange.End().Sub(timeRange.Start())

	// Проверяем, что исходные значения не изменились
	assert.Equal(t, originalStart, timeRange.Start().String())
	assert.Equal(t, originalEnd, timeRange.End().String())
}

func TestTimeRange_ContainsEdgeCases(t *testing.T) {
	now := time.Now()
	start := now.Add(-1 * time.Hour)
	end := now.Add(1 * time.Hour)

	timeRange, err := valueobjects.NewTimeRange(start, end)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		moment   time.Time
		contains bool
	}{
		{
			name:     "Exactly at start",
			moment:   start,
			contains: true,
		},
		{
			name:     "Exactly at end",
			moment:   end,
			contains: true,
		},
		{
			name:     "One nanosecond after start",
			moment:   start.Add(1 * time.Nanosecond),
			contains: true,
		},
		{
			name:     "One nanosecond before end",
			moment:   end.Add(-1 * time.Nanosecond),
			contains: true,
		},
		{
			name:     "Very close to start",
			moment:   start.Add(1 * time.Microsecond),
			contains: true,
		},
		{
			name:     "Very close to end",
			moment:   end.Add(-1 * time.Microsecond),
			contains: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contains := timeRange.Contains(valueobjects.NewDateTime(tt.moment))
			assert.Equal(t, tt.contains, contains)
		})
	}
}

func TestTimeRange_DurationEdgeCases(t *testing.T) {
	now := time.Now()

	t.Run("Minimal duration", func(t *testing.T) {
		start := now
		end := now.Add(1 * time.Nanosecond)

		timeRange, err := valueobjects.NewTimeRange(start, end)
		assert.NoError(t, err)

		duration := timeRange.Duration()
		assert.True(t, duration > 0)
		assert.InDelta(t, float64(1*time.Nanosecond), float64(duration), float64(1*time.Nanosecond))
	})

	t.Run("Zero duration", func(t *testing.T) {
		// Пытаемся создать интервал с нулевой продолжительностью
		_, err := valueobjects.NewTimeRange(now, now)
		assert.Error(t, err)
		assert.Equal(t, valueobjects.ErrInvalidTimeRange, err)
	})

	t.Run("Long duration", func(t *testing.T) {
		start := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2038, 1, 19, 3, 14, 7, 0, time.UTC) // Максимальное время для 32-битных систем

		timeRange, err := valueobjects.NewTimeRange(start, end)
		assert.NoError(t, err)

		expectedDuration := end.Sub(start)
		actualDuration := timeRange.Duration()

		// Допускаем небольшую погрешность из-за преобразования в UTC
		assert.InDelta(t, float64(expectedDuration), float64(actualDuration), float64(1*time.Millisecond))
	})
}
