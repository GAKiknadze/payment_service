package valueobject_test

import (
	"errors"
	"testing"
	"time"

	"github.com/GAKiknadze/payment_service/domain/common/valueobject"
	"github.com/shopspring/decimal"
)

func TestNewQuotaDefinition_ValidRecurringQuota(t *testing.T) {
	// Given - валидные параметры для периодической квоты
	resourceType := "tokens"
	limit := decimal.NewFromFloat(1000)
	unit := "count"
	isRecurring := true
	resetPeriod := 30 * 24 * time.Hour // 30 дней

	// When - создаем квоту
	quota, err := valueobject.NewQuotaDefinition(
		resourceType, limit, unit, isRecurring, resetPeriod,
	)

	// Then - проверяем, что ошибка отсутствует и данные корректны
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if quota.ResourceType() != resourceType {
		t.Errorf("Expected resource type %s, got %s", resourceType, quota.ResourceType())
	}

	if !quota.Limit().Equal(limit) {
		t.Errorf("Expected limit %s, got %s", limit, quota.Limit())
	}

	if quota.Unit() != unit {
		t.Errorf("Expected unit %s, got %s", unit, quota.Unit())
	}

	if quota.IsRecurring() != isRecurring {
		t.Errorf("Expected isRecurring %v, got %v", isRecurring, quota.IsRecurring())
	}

	if quota.ResetPeriod() != resetPeriod {
		t.Errorf("Expected resetPeriod %v, got %v", resetPeriod, quota.ResetPeriod())
	}
}

func TestNewQuotaDefinition_ValidNonRecurringQuota(t *testing.T) {
	// Given - валидные параметры для разовой квоты
	resourceType := "ssl_certificates"
	limit := decimal.NewFromFloat(1)
	unit := "count"
	isRecurring := false
	resetPeriod := time.Duration(0) // Для непериодических квот resetPeriod должен быть 0

	// When - создаем квоту
	quota, err := valueobject.NewQuotaDefinition(
		resourceType, limit, unit, isRecurring, resetPeriod,
	)

	// Then - проверяем, что ошибка отсутствует
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Проверяем, что квота действительно непериодическая
	if quota.IsRecurring() {
		t.Error("Expected non-recurring quota, but got recurring")
	}
}

func TestNewQuotaDefinition_EmptyResourceType(t *testing.T) {
	// Given - пустой тип ресурса
	resourceType := ""
	limit := decimal.NewFromFloat(1000)
	unit := "count"
	isRecurring := true
	resetPeriod := 30 * 24 * time.Hour

	// When - пытаемся создать квоту
	_, err := valueobject.NewQuotaDefinition(
		resourceType, limit, unit, isRecurring, resetPeriod,
	)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for empty resource type, got nil")
	}

	if err != valueobject.ErrInvalidResourceType {
		t.Errorf("Expected ErrInvalidResourceType, got %v", err)
	}
}

func TestNewQuotaDefinition_NonPositiveLimit(t *testing.T) {
	// Given - неположительный лимит
	resourceType := "tokens"
	limit := decimal.NewFromFloat(0)
	unit := "count"
	isRecurring := true
	resetPeriod := 30 * 24 * time.Hour

	// When - пытаемся создать квоту
	_, err := valueobject.NewQuotaDefinition(
		resourceType, limit, unit, isRecurring, resetPeriod,
	)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for non-positive limit, got nil")
	}

	if err != valueobject.ErrInvalidQuotaLimit {
		t.Errorf("Expected ErrInvalidQuotaLimit, got %v", err)
	}
}

func TestNewQuotaDefinition_NonRecurringWithResetPeriod(t *testing.T) {
	// Given - непериодическая квота с ненулевым resetPeriod
	resourceType := "ssl_certificates"
	limit := decimal.NewFromFloat(1)
	unit := "count"
	isRecurring := false
	resetPeriod := 30 * 24 * time.Hour

	// When - пытаемся создать квоту
	_, err := valueobject.NewQuotaDefinition(
		resourceType, limit, unit, isRecurring, resetPeriod,
	)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for non-recurring quota with reset period, got nil")
	}

	if err != valueobject.ErrNonRecurringWithReset {
		t.Errorf("Expected ErrNonRecurringWithReset, got %v", err)
	}
}

func TestNewQuotaDefinition_RecurringWithZeroResetPeriod(t *testing.T) {
	// Given - периодическая квота с нулевым resetPeriod
	resourceType := "tokens"
	limit := decimal.NewFromFloat(1000)
	unit := "count"
	isRecurring := true
	resetPeriod := time.Duration(0)

	// When - пытаемся создать квоту
	_, err := valueobject.NewQuotaDefinition(
		resourceType, limit, unit, isRecurring, resetPeriod,
	)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for recurring quota with zero reset period, got nil")
	}

	if err != valueobject.ErrInvalidResetPeriod {
		t.Errorf("Expected ErrInvalidResetPeriod, got %v", err)
	}
}

func TestIsWithinLimit_WithinLimit(t *testing.T) {
	// Given - квота и текущее использование
	quota, _ := valueobject.NewQuotaDefinition(
		"tokens", decimal.NewFromFloat(1000), "count", true, 30*24*time.Hour,
	)
	currentUsage := decimal.NewFromFloat(500)
	increment := decimal.NewFromFloat(400)

	// When - проверяем, умещается ли увеличение в лимит
	withinLimit := quota.IsWithinLimit(currentUsage, increment)

	// Then - увеличение должно умещаться в лимит
	if !withinLimit {
		t.Error("Expected increment to be within limit")
	}
}

func TestIsWithinLimit_ExceedsLimit(t *testing.T) {
	// Given - квота и текущее использование
	quota, _ := valueobject.NewQuotaDefinition(
		"tokens", decimal.NewFromFloat(1000), "count", true, 30*24*time.Hour,
	)
	currentUsage := decimal.NewFromFloat(900)
	increment := decimal.NewFromFloat(200)

	// When - проверяем, умещается ли увеличение в лимит
	withinLimit := quota.IsWithinLimit(currentUsage, increment)

	// Then - увеличение не должно умещаться в лимит
	if withinLimit {
		t.Error("Expected increment to exceed limit")
	}
}

func TestCalculateRemaining(t *testing.T) {
	// Given - квота и текущее использование
	quota, _ := valueobject.NewQuotaDefinition(
		"tokens", decimal.NewFromFloat(1000), "count", true, 30*24*time.Hour,
	)

	cases := []struct {
		name         string
		currentUsage decimal.Decimal
		expected     decimal.Decimal
	}{
		{
			"half used",
			decimal.NewFromFloat(500),
			decimal.NewFromFloat(500),
		},
		{
			"fully used",
			decimal.NewFromFloat(1000),
			decimal.Zero,
		},
		{
			"overused",
			decimal.NewFromFloat(1200),
			decimal.Zero,
		},
		{
			"not used",
			decimal.Zero,
			decimal.NewFromFloat(1000),
		},
	}

	// When & Then - проверяем оставшийся лимит для разных сценариев
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			remaining := quota.CalculateRemaining(tc.currentUsage)
			if !remaining.Equal(tc.expected) {
				t.Errorf("Expected remaining %s, got %s", tc.expected, remaining)
			}
		})
	}
}

func TestCanUse_RecurringQuota(t *testing.T) {
	// Given - периодическая квота
	quota, _ := valueobject.NewQuotaDefinition(
		"tokens", decimal.NewFromFloat(1000), "count", true, 30*24*time.Hour,
	)

	cases := []struct {
		name         string
		currentUsage decimal.Decimal
		amount       decimal.Decimal
		expected     bool
	}{
		{
			"within limit",
			decimal.NewFromFloat(500),
			decimal.NewFromFloat(400),
			true,
		},
		{
			"at limit",
			decimal.NewFromFloat(1000),
			decimal.NewFromFloat(0),
			false,
		},
		{
			"exceeds limit",
			decimal.NewFromFloat(900),
			decimal.NewFromFloat(200),
			false,
		},
		{
			"zero amount",
			decimal.NewFromFloat(500),
			decimal.Zero,
			false,
		},
	}

	// When & Then - проверяем использование для разных сценариев
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			canUse := quota.CanUse(tc.currentUsage, tc.amount)
			if canUse != tc.expected {
				t.Errorf("Expected canUse to be %v, got %v", tc.expected, canUse)
			}
		})
	}
}

func TestCanUse_NonRecurringQuota(t *testing.T) {
	// Given - разовая квота
	quota, _ := valueobject.NewQuotaDefinition(
		"ssl_certificates", decimal.NewFromFloat(1), "count", false, 0,
	)

	cases := []struct {
		name     string
		amount   decimal.Decimal
		expected bool
	}{
		{
			"valid amount",
			decimal.NewFromFloat(1),
			true,
		},
		{
			"less than limit",
			decimal.NewFromFloat(0.5),
			true,
		},
		{
			"exceeds limit",
			decimal.NewFromFloat(2),
			false,
		},
		{
			"zero amount",
			decimal.Zero,
			false,
		},
	}

	// When & Then - проверяем использование для разовой квоты
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Для разовых квот текущее использование не учитывается
			canUse := quota.CanUse(decimal.Zero, tc.amount)
			if canUse != tc.expected {
				t.Errorf("Expected canUse to be %v, got %v", tc.expected, canUse)
			}
		})
	}
}

func TestFormatUsage(t *testing.T) {
	// Given - квота
	quota, _ := valueobject.NewQuotaDefinition(
		"tokens", decimal.NewFromFloat(1000), "count", true, 30*24*time.Hour,
	)

	cases := []struct {
		name         string
		currentUsage decimal.Decimal
		expected     string
	}{
		{
			"half used",
			decimal.NewFromFloat(500),
			"500/1000 count",
		},
		{
			"fully used",
			decimal.NewFromFloat(1000),
			"1000/1000 count",
		},
		{
			"overused",
			decimal.NewFromFloat(1200),
			"1000/1000 count",
		},
		{
			"not used",
			decimal.Zero,
			"0/1000 count",
		},
	}

	// When & Then - проверяем форматирование использования
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			formatted := quota.FormatUsage(tc.currentUsage)
			if formatted != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, formatted)
			}
		})
	}
}

func TestIsExceeded(t *testing.T) {
	// Given - квота
	quota, _ := valueobject.NewQuotaDefinition(
		"tokens", decimal.NewFromFloat(1000), "count", true, 30*24*time.Hour,
	)

	cases := []struct {
		name         string
		currentUsage decimal.Decimal
		expected     bool
	}{
		{
			"below limit",
			decimal.NewFromFloat(500),
			false,
		},
		{
			"at limit",
			decimal.NewFromFloat(1000),
			true,
		},
		{
			"over limit",
			decimal.NewFromFloat(1200),
			true,
		},
		{
			"not used",
			decimal.Zero,
			false,
		},
	}

	// When & Then - проверяем превышение лимита
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			exceeded := quota.IsExceeded(tc.currentUsage)
			if exceeded != tc.expected {
				t.Errorf("Expected exceeded to be %v, got %v", tc.expected, exceeded)
			}
		})
	}
}

func TestUsagePercentage(t *testing.T) {
	// Given - квота
	quota, _ := valueobject.NewQuotaDefinition(
		"tokens", decimal.NewFromFloat(1000), "count", true, 30*24*time.Hour,
	)

	cases := []struct {
		name         string
		currentUsage decimal.Decimal
		expected     float64
	}{
		{
			"half used",
			decimal.NewFromFloat(500),
			50.0,
		},
		{
			"fully used",
			decimal.NewFromFloat(1000),
			100.0,
		},
		{
			"overused",
			decimal.NewFromFloat(1200),
			100.0,
		},
		{
			"not used",
			decimal.Zero,
			0.0,
		},
	}

	// When & Then - проверяем процент использования
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			percentage := quota.UsagePercentage(tc.currentUsage)
			if percentage != tc.expected {
				t.Errorf("Expected percentage %f, got %f", tc.expected, percentage)
			}
		})
	}
}

func TestNeedsReset(t *testing.T) {
	// Given - периодическая квота
	quota, _ := valueobject.NewQuotaDefinition(
		"tokens", decimal.NewFromFloat(1000), "count", true, 30*24*time.Hour,
	)

	// Создаем время сброса
	lastReset := time.Now().Add(-25 * 24 * time.Hour) // 25 дней назад

	// When - проверяем, нужно ли сбросить квоту
	needsReset := quota.NeedsReset(lastReset)

	// Then - сброс еще не нужен (менее 30 дней)
	if needsReset {
		t.Error("Expected no reset needed")
	}

	// Given - время сброса 35 дней назад
	lastReset = time.Now().Add(-35 * 24 * time.Hour)

	// When - проверяем, нужно ли сбросить квоту
	needsReset = quota.NeedsReset(lastReset)

	// Then - сброс нужен (более 30 дней)
	if !needsReset {
		t.Error("Expected reset needed")
	}
}

func TestNeedsReset_NonRecurringQuota(t *testing.T) {
	// Given - разовая квота
	quota, _ := valueobject.NewQuotaDefinition(
		"ssl_certificates", decimal.NewFromFloat(1), "count", false, 0,
	)

	// When - проверяем, нужно ли сбросить квоту
	needsReset := quota.NeedsReset(time.Now().Add(-100 * 24 * time.Hour))

	// Then - для разовой квоты сброс никогда не нужен
	if needsReset {
		t.Error("Expected no reset needed for non-recurring quota")
	}
}

func TestNextResetTime(t *testing.T) {
	// Given - периодическая квота
	quota, _ := valueobject.NewQuotaDefinition(
		"tokens", decimal.NewFromFloat(1000), "count", true, 30*24*time.Hour,
	)

	// Создаем время последнего сброса
	lastReset := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	// When - получаем время следующего сброса
	nextReset := quota.NextResetTime(lastReset)

	// Then - проверяем, что время следующего сброса верное
	expected := time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)
	if !nextReset.Equal(expected) {
		t.Errorf("Expected next reset time %v, got %v", expected, nextReset)
	}
}

func TestNextResetTime_NonRecurringQuota(t *testing.T) {
	// Given - разовая квота
	quota, _ := valueobject.NewQuotaDefinition(
		"ssl_certificates", decimal.NewFromFloat(1), "count", false, 0,
	)

	// When - получаем время следующего сброса
	nextReset := quota.NextResetTime(time.Now())

	// Then - для разовой квоты время сброса должно быть нулевым
	if !nextReset.IsZero() {
		t.Error("Expected zero time for non-recurring quota")
	}
}

func TestEquals(t *testing.T) {
	// Given - две одинаковые квоты
	quota1, _ := valueobject.NewQuotaDefinition(
		"tokens", decimal.NewFromFloat(1000), "count", true, 30*24*time.Hour,
	)
	quota2, _ := valueobject.NewQuotaDefinition(
		"tokens", decimal.NewFromFloat(1000), "count", true, 30*24*time.Hour,
	)

	// When - проверяем равенство
	equal := quota1.Equals(quota2)

	// Then - квоты должны быть равны
	if !equal {
		t.Error("Expected quotas to be equal")
	}

	// Given - две разные квоты
	quota3, _ := valueobject.NewQuotaDefinition(
		"tokens", decimal.NewFromFloat(500), "count", true, 30*24*time.Hour,
	)

	// When - проверяем равенство
	equal = quota1.Equals(quota3)

	// Then - квоты не должны быть равны
	if equal {
		t.Error("Expected quotas to be different")
	}
}

func TestValidate(t *testing.T) {
	// Given - валидная квота
	quota, _ := valueobject.NewQuotaDefinition(
		"tokens", decimal.NewFromFloat(1000), "count", true, 30*24*time.Hour,
	)

	// When - проверяем валидность
	err := quota.Validate()

	// Then - ошибки быть не должно
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Given - невалидная квота (пустой тип ресурса)
	invalidQuota := valueobject.NewQuotaDefinitionForTest(
		"", decimal.NewFromFloat(1000), "count", true, 30*24*time.Hour,
	)

	// When - проверяем валидность
	err = invalidQuota.Validate()

	// Then - должна быть ошибка
	if err == nil {
		t.Fatal("Expected error for invalid quota, got nil")
	}

	// Проверяем, что это именно ошибка пустого типа ресурса
	if !errors.Is(err, valueobject.ErrInvalidResourceType) {
		t.Errorf("Expected ErrInvalidResourceType, got %v", err)
	}
}
