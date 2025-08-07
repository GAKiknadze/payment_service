package model_test

import (
	"errors"
	"testing"
	"time"

	"github.com/GAKiknadze/payment_service/domain/common/interfaces"
	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
	"github.com/GAKiknadze/payment_service/domain/common/valueobjects/fixtures"
	"github.com/GAKiknadze/payment_service/domain/common/valueobjects/mocks"
	"github.com/GAKiknadze/payment_service/domain/tariff/event"
	"github.com/GAKiknadze/payment_service/domain/tariff/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тестовые константы
const (
	TestYear        = 2025
	TestMonth       = time.October
	TestDay         = 1
	TestHour        = 12
	TestMinute      = 0
	TestSecond      = 0
	DefaultTariffID = "default-tariff-id"
)

// Тестовые фикстуры
var (
	TestClock = mocks.NewFixedClock(time.Date(TestYear, TestMonth, TestDay, TestHour, TestMinute, TestSecond, 0, time.UTC))
)

// Вспомогательные функции
func newTestDateTime(year int, month time.Month, day, hour, min, sec int) valueobjects.DateTime {
	return valueobjects.NewDateTime(time.Date(year, month, day, hour, min, sec, 0, time.UTC))
}

// tariffFixture создает тестовый тариф с заданными параметрами
type tariffFixture struct {
	ID        valueobjects.TariffID
	Name      string
	Price     valueobjects.Money
	StartDate time.Time
	Duration  time.Duration
	Clock     interfaces.Clock
}

// Build создает тариф на основе фикстуры
func (f *tariffFixture) Build(t *testing.T) *model.Tariff {
	if f.ID == "" {
		f.ID = valueobjects.GenerateTariffID()
	}
	if f.Name == "" {
		f.Name = "Standard Plan"
	}
	if f.Price.IsZero() {
		f.Price = fixtures.NewTestMoneyRUB(1000)
	}
	if f.StartDate.IsZero() {
		f.StartDate = TestClock.Now()
	}
	if f.Duration == 0 {
		f.Duration = 730 * time.Hour // Месячный тариф по умолчанию
	}
	if f.Clock == nil {
		f.Clock = TestClock
	}

	endDate := f.StartDate.Add(f.Duration)
	period, err := valueobjects.NewTimeRange(f.StartDate, endDate)
	require.NoError(t, err, "Failed to create time range")

	tariff, err := model.NewTariff(f.ID, f.Name, f.Price, period, f.Clock)
	require.NoError(t, err, "Failed to create tariff")

	return tariff
}

// updatePrice обновляет цену тарифа и проверяет результат
func (f *tariffFixture) updatePrice(t *testing.T, tariff *model.Tariff, newPrice valueobjects.Money, effectiveDate valueobjects.DateTime) {
	err := tariff.UpdatePrice(newPrice, effectiveDate, f.Clock)
	require.NoError(t, err, "Price update failed")

	assert.Equal(t, newPrice, tariff.Price(), "Price was not updated correctly")
}

// assertCost проверяет расчет стоимости с учетом допустимой погрешности
func assertCost(t *testing.T, tariff *model.Tariff, duration time.Duration, expected valueobjects.Money) {
	cost, err := tariff.CalculateCost(duration)
	require.NoError(t, err, "Cost calculation failed")

	// Для финансовых вычислений используем погрешность 0.01 (1 копейка)
	assert.True(t, cost.IsApproximatelyEqual(expected, 0.01),
		"Expected cost %.2f, got %.2f (difference: %.2f)",
		expected.Amount().BigFloat(), cost.Amount().BigFloat(), cost.Amount().Sub(expected.Amount()).Abs().BigFloat())
}

// assertEventPublished проверяет, что определенное событие было опубликовано
func assertEventPublished[T any](t *testing.T, tariff *model.Tariff, eventType T, assertion func(event T)) {
	events := tariff.PopEvents()

	var found bool
	for _, e := range events {
		if typedEvent, ok := e.(T); ok {
			found = true
			assertion(typedEvent)
			break
		}
	}

	assert.True(t, found, "Expected event %T was not published", eventType)
}

// --- Тесты создания тарифа ---

func TestNewTariff(t *testing.T) {
	t.Run("valid creation", func(t *testing.T) {
		tariff := (&tariffFixture{
			ID:        valueobjects.GenerateTariffID(),
			Name:      "Standard Plan",
			Price:     fixtures.NewTestMoneyRUB(1000),
			StartDate: newTestDateTime(2025, 1, 1, 0, 0, 0).Time(),
			Duration:  365 * 24 * time.Hour,
			Clock:     TestClock,
		}).Build(t)

		assert.Equal(t, "Standard Plan", tariff.Name())
		assert.Equal(t, fixtures.NewTestMoneyRUB(1000), tariff.Price())
		assert.Empty(t, tariff.PopEvents())
	})

	// t.Run("empty name", func(t *testing.T) {
	// 	_, err := tariffFixture{
	// 		Name:  "",
	// 		Clock: TestClock,
	// 	}.Build(t)

	// 	assert.ErrorIs(t, err, model.ErrInvalidTariffName)
	// })

	// t.Run("zero duration period", func(t *testing.T) {
	// 	_, err := tariffFixture{
	// 		Duration: 0,
	// 		Clock:    TestClock,
	// 	}.Build(t)

	// 	assert.Error(t, err)
	// })

	// t.Run("expired period", func(t *testing.T) {
	// 	_, err := tariffFixture{
	// 		StartDate: newTestDateTime(2024, 1, 1, 0, 0, 0).Time(),
	// 		Duration:  365 * 24 * time.Hour,
	// 		Clock:     TestClock,
	// 	}.Build(t)

	// 	assert.ErrorIs(t, err, model.ErrPeriodAlreadyExpired)
	// })
}

// --- Тесты обновления цены ---

func TestUpdatePrice(t *testing.T) {
	t.Run("valid price update", func(t *testing.T) {
		tariff := (&tariffFixture{Clock: TestClock}).Build(t)
		effectiveDate := newTestDateTime(2025, time.October, 6, 0, 0, 0)

		tariff.UpdatePrice(fixtures.NewTestMoneyRUB(1200), effectiveDate, TestClock)

		assert.Equal(t, fixtures.NewTestMoneyRUB(1200), tariff.Price())
		assert.Equal(t, 1, tariff.Version())

		assertEventPublished(t, tariff, event.TariffPriceUpdated{}, func(e event.TariffPriceUpdated) {
			assert.Equal(t, tariff.ID(), e.TariffID)
			assert.Equal(t, fixtures.NewTestMoneyRUB(1000), e.OldPrice)
			assert.Equal(t, fixtures.NewTestMoneyRUB(1200), e.NewPrice)
			assert.Equal(t, effectiveDate, e.EffectiveDate)
		})
	})

	t.Run("past effective date", func(t *testing.T) {
		tariff := (&tariffFixture{Clock: TestClock}).Build(t)
		err := tariff.UpdatePrice(
			fixtures.NewTestMoneyRUB(1200),
			newTestDateTime(2025, 5, 1, 0, 0, 0),
			TestClock,
		)

		assert.Error(t, err)
		assert.Equal(t, fixtures.NewTestMoneyRUB(1000), tariff.Price())
		assert.Empty(t, tariff.PopEvents())
	})

	t.Run("multiple price updates", func(t *testing.T) {
		tariff := (&tariffFixture{Clock: TestClock}).Build(t)

		// Проверяем начальную цену
		assert.Equal(t, fixtures.NewTestMoneyRUB(1000), tariff.Price())

		// Первое обновление
		err := tariff.UpdatePrice(
			fixtures.NewTestMoneyRUB(1200),
			newTestDateTime(2025, time.October, 7, 13, 0, 0),
			TestClock,
		)
		assert.NoError(t, err)

		// Проверяем, что цена действительно обновилась
		assert.Equal(t, fixtures.NewTestMoneyRUB(1200), tariff.Price())
		// Очищаем события
		tariff.PopEvents()

		// Второе обновление
		err = tariff.UpdatePrice(
			fixtures.NewTestMoneyRUB(1500),
			newTestDateTime(2025, time.October, 7, 14, 0, 0),
			TestClock,
		)
		assert.NoError(t, err)

		// Проверяем финальную цену
		assert.Equal(t, fixtures.NewTestMoneyRUB(1500), tariff.Price())
		assert.Equal(t, 2, tariff.Version())

		// Проверяем события
		events := tariff.PopEvents()
		assert.Len(t, events, 1)

		priceUpdated, ok := events[0].(event.TariffPriceUpdated)
		assert.True(t, ok)
		assert.Equal(t, fixtures.NewTestMoneyRUB(1200), priceUpdated.OldPrice)
		assert.Equal(t, fixtures.NewTestMoneyRUB(1500), priceUpdated.NewPrice)
	})
}

// --- Тесты расчета стоимости ---

func TestCalculateCost(t *testing.T) {
	tests := []struct {
		name        string
		duration    time.Duration
		expected    valueobjects.Money
		tariffSetup func(*tariffFixture)
	}{
		{
			name:     "one hour",
			duration: time.Hour,
			expected: fixtures.NewTestMoneyRUB(1000.0 / 730.0),
		},
		{
			name:     "one day",
			duration: 24 * time.Hour,
			expected: fixtures.NewTestMoneyRUB((1000.0 / 730.0) * 24.0),
		},
		{
			name:     "one month",
			duration: 730 * time.Hour,
			expected: fixtures.NewTestMoneyRUB(1000),
		},
		{
			name:     "one year",
			duration: 730 * 12 * time.Hour,
			expected: fixtures.NewTestMoneyRUB(12000),
		},
		{
			name:     "one hour thirty minutes",
			duration: 1*time.Hour + 30*time.Minute,
			expected: fixtures.NewTestMoneyRUB(1.5 * 1000.0 / 730.0),
		},
		// {
		// 	name:     "negative duration",
		// 	duration: -time.Hour,
		// 	expected: valueobjects.Money{},
		// 	tariffSetup: func(f *tariffFixture) {
		// 		f.SetupErrorTest = true
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := tariffFixture{Clock: TestClock}
			if tt.tariffSetup != nil {
				tt.tariffSetup(&fixture)
			}

			tariff := fixture.Build(t)

			if tt.name == "negative duration" {
				cost, err := tariff.CalculateCost(tt.duration)
				assert.Error(t, err)
				assert.True(t, errors.Is(err, valueobjects.ErrInvalidBillingPeriod))
				assert.Equal(t, valueobjects.Money{}, cost)
				return
			}

			assertCost(t, tariff, tt.duration, tt.expected)
		})
	}
}

// --- Тесты для TariffPeriod Value Object ---

func TestTariffPeriod(t *testing.T) {
	t.Run("valid period creation", func(t *testing.T) {
		start := newTestDateTime(2025, 1, 1, 0, 0, 0)
		end := newTestDateTime(2025, 12, 31, 23, 59, 59)
		period, err := valueobjects.NewTimeRange(start.Time(), end.Time())

		require.NoError(t, err)
		assert.True(t, period.Contains(start))
		assert.True(t, period.Contains(end))
		assert.True(t, period.Contains(newTestDateTime(2025, 6, 15, 12, 0, 0)))
		assert.Equal(t, 364*24*time.Hour+23*time.Hour+59*time.Minute+59*time.Second, period.Duration())
	})

	t.Run("invalid period (end before start)", func(t *testing.T) {
		start := newTestDateTime(2025, 12, 31, 23, 59, 59)
		end := newTestDateTime(2025, 1, 1, 0, 0, 0)
		period, err := valueobjects.NewTimeRange(start.Time(), end.Time())

		assert.Error(t, err)
		assert.Equal(t, valueobjects.TimeRange{}, period)
	})

	t.Run("edge cases", func(t *testing.T) {
		// Очень короткий период (1 секунда)
		start := newTestDateTime(2025, 1, 1, 0, 0, 0)
		end := newTestDateTime(2025, 1, 1, 0, 0, 1)
		period, err := valueobjects.NewTimeRange(start.Time(), end.Time())

		require.NoError(t, err)
		assert.Equal(t, time.Second, period.Duration())

		// Проверка активности на границах
		assert.True(t, period.Contains(start))
		assert.True(t, period.Contains(end))
		assert.False(t, period.Contains(newTestDateTime(2022, 12, 31, 23, 59, 59)))
		assert.False(t, period.Contains(newTestDateTime(2025, 1, 1, 0, 0, 2)))
	})
}

// --- Тесты для edge cases ---

func TestTariffEdgeCases(t *testing.T) {
	t.Run("very long period", func(t *testing.T) {
		// 10 лет без високосных дней
		tenYears := 10 * 365 * 24 * time.Hour
		oneYear := 365 * 24 * time.Hour

		tariff := (&tariffFixture{
			Price:    fixtures.NewTestMoneyRUB(10000),
			Duration: tenYears,
			Clock:    TestClock,
		}).Build(t)

		assertCost(t, tariff, oneYear, fixtures.NewTestMoneyRUB(1000))
	})

	t.Run("very short period", func(t *testing.T) {
		tariff := (&tariffFixture{
			Price:    fixtures.NewTestMoneyRUB(1),
			Duration: time.Minute,
			Clock:    TestClock,
		}).Build(t)

		assertCost(t, tariff, 30*time.Second, fixtures.NewTestMoneyRUB(0.5))
	})

	t.Run("leap years calculation", func(t *testing.T) {
		leapClock := mocks.NewFixedClock(time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC))

		tariff := (&tariffFixture{
			Price:    fixtures.NewTestMoneyRUB(366),
			Duration: 366 * 24 * time.Hour, // Високосный год
			Clock:    leapClock,
		}).Build(t)

		assertCost(t, tariff, 365*24*time.Hour, fixtures.NewTestMoneyRUB(365))
	})

	t.Run("high precision money", func(t *testing.T) {
		yearDuration := 365 * 24 * time.Hour
		expectedHourlyCost := 1234.5678 / yearDuration.Hours()

		tariff := (&tariffFixture{
			Price:    fixtures.NewTestMoneyRUB(1234.5678),
			Duration: yearDuration,
			Clock:    TestClock,
		}).Build(t)

		assertCost(t, tariff, time.Hour, fixtures.NewTestMoneyRUB(expectedHourlyCost))
	})

	t.Run("multiple price updates", func(t *testing.T) {
		tariff := (&tariffFixture{Clock: TestClock}).Build(t)
		initialVersion := tariff.Version()

		// Выполняем несколько обновлений
		err := tariff.UpdatePrice(
			fixtures.NewTestMoneyRUB(1200),
			newTestDateTime(2025, time.October, 6, 0, 0, 0),
			TestClock,
		)
		assert.NoError(t, err)

		err = tariff.UpdatePrice(
			fixtures.NewTestMoneyRUB(1500),
			newTestDateTime(2025, time.October, 7, 0, 0, 0),
			TestClock,
		)
		assert.NoError(t, err)

		err = tariff.UpdatePrice(
			fixtures.NewTestMoneyRUB(1800),
			newTestDateTime(2025, time.October, 8, 0, 0, 0),
			TestClock,
		)
		assert.NoError(t, err)

		assert.Equal(t, initialVersion+3, tariff.Version())
		assert.Equal(t, fixtures.NewTestMoneyRUB(1800), tariff.Price())

		events := tariff.PopEvents()
		assert.Len(t, events, 3)

		// Проверяем последнее событие
		lastEvent, ok := events[2].(event.TariffPriceUpdated)
		assert.True(t, ok)
		assert.True(t, lastEvent.OldPrice.IsApproximatelyEqual(fixtures.NewTestMoneyRUB(1500), 0.01))
		assert.True(t, lastEvent.NewPrice.IsApproximatelyEqual(fixtures.NewTestMoneyRUB(1800), 0.01))
	})
}
