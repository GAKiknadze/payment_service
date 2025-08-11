package valueobject

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidQuotaLimit     = errors.New("quota limit must be greater than zero")
	ErrInvalidResourceType   = errors.New("resource type cannot be empty")
	ErrInvalidUnit           = errors.New("unit cannot be empty")
	ErrInvalidResetPeriod    = errors.New("reset period must be positive for recurring quotas")
	ErrNonRecurringWithReset = errors.New("non-recurring quota cannot have reset period")
)

// QuotaDefinition представляет определение квоты для тарифа
type QuotaDefinition struct {
	resourceType string
	limit        decimal.Decimal
	unit         string
	isRecurring  bool
	resetPeriod  time.Duration
}

// NewQuotaDefinition - фабричный метод для создания квоты с валидацией
func NewQuotaDefinition(
	resourceType string,
	limit decimal.Decimal,
	unit string,
	isRecurring bool,
	resetPeriod time.Duration,
) (QuotaDefinition, error) {
	// Валидация обязательных полей
	if resourceType == "" {
		return QuotaDefinition{}, ErrInvalidResourceType
	}

	if unit == "" {
		return QuotaDefinition{}, ErrInvalidUnit
	}

	// Проверка лимита
	if limit.Cmp(decimal.Zero) <= 0 {
		return QuotaDefinition{}, ErrInvalidQuotaLimit
	}

	// Проверка resetPeriod для периодических квот
	if isRecurring && resetPeriod <= 0 {
		return QuotaDefinition{}, ErrInvalidResetPeriod
	}

	// Проверка resetPeriod для непериодических квот
	if !isRecurring && resetPeriod > 0 {
		return QuotaDefinition{}, ErrNonRecurringWithReset
	}

	return QuotaDefinition{
		resourceType: resourceType,
		limit:        limit,
		unit:         unit,
		isRecurring:  isRecurring,
		resetPeriod:  resetPeriod,
	}, nil
}

// ResourceType возвращает тип ресурса
func (qd QuotaDefinition) ResourceType() string {
	return qd.resourceType
}

// Limit возвращает максимальный лимит квоты
func (qd QuotaDefinition) Limit() decimal.Decimal {
	return qd.limit
}

// Unit возвращает единицу измерения
func (qd QuotaDefinition) Unit() string {
	return qd.unit
}

// IsRecurring проверяет, является ли квота периодической
func (qd QuotaDefinition) IsRecurring() bool {
	return qd.isRecurring
}

// ResetPeriod возвращает период сброса для периодических квот
func (qd QuotaDefinition) ResetPeriod() time.Duration {
	return qd.resetPeriod
}

// IsWithinLimit проверяет, умещается ли запрашиваемое увеличение в лимит
func (qd QuotaDefinition) IsWithinLimit(currentUsage decimal.Decimal, increment decimal.Decimal) bool {
	if increment.Cmp(decimal.Zero) <= 0 {
		return false
	}

	newUsage := currentUsage.Add(increment)
	return newUsage.Cmp(qd.limit) <= 0
}

// CalculateRemaining возвращает оставшееся количество ресурса
func (qd QuotaDefinition) CalculateRemaining(currentUsage decimal.Decimal) decimal.Decimal {
	if currentUsage.Cmp(qd.limit) >= 0 {
		return decimal.Zero
	}

	return qd.limit.Sub(currentUsage)
}

// CanUse проверяет, можно ли использовать указанное количество ресурса
func (qd QuotaDefinition) CanUse(currentUsage decimal.Decimal, amount decimal.Decimal) bool {
	if amount.Cmp(decimal.Zero) <= 0 {
		return false
	}

	// Для разовых квот проверяем, что amount <= limit
	if !qd.isRecurring {
		return amount.Cmp(qd.limit) <= 0
	}

	// Для периодических квот проверяем текущее использование
	return qd.IsWithinLimit(currentUsage, amount)
}

// FormatLimit возвращает отформатированное представление лимита
func (qd QuotaDefinition) FormatLimit() string {
	return fmt.Sprintf("%s %s", qd.limit.String(), qd.unit)
}

// FormatRemaining возвращает отформатированное представление оставшегося лимита
func (qd QuotaDefinition) FormatRemaining(currentUsage decimal.Decimal) string {
	remaining := qd.CalculateRemaining(currentUsage)
	return fmt.Sprintf("%s %s", remaining.String(), qd.unit)
}

// FormatUsage возвращает строку использования квоты (например, "500/1000 tokens")
func (qd QuotaDefinition) FormatUsage(currentUsage decimal.Decimal) string {
	used := currentUsage
	if used.Cmp(qd.limit) > 0 {
		used = qd.limit
	}

	return fmt.Sprintf("%s/%s %s", used.String(), qd.limit.String(), qd.unit)
}

// Equals проверяет равенство двух квот
func (qd QuotaDefinition) Equals(other QuotaDefinition) bool {
	return qd.resourceType == other.resourceType &&
		qd.limit.Equal(other.limit) &&
		qd.unit == other.unit &&
		qd.isRecurring == other.isRecurring &&
		qd.resetPeriod == other.resetPeriod
}

// Validate проверяет корректность квоты
func (qd QuotaDefinition) Validate() error {
	if qd.resourceType == "" {
		return ErrInvalidResourceType
	}

	if qd.unit == "" {
		return ErrInvalidUnit
	}

	if qd.limit.Cmp(decimal.Zero) <= 0 {
		return ErrInvalidQuotaLimit
	}

	if qd.isRecurring && qd.resetPeriod <= 0 {
		return ErrInvalidResetPeriod
	}

	if !qd.isRecurring && qd.resetPeriod > 0 {
		return ErrNonRecurringWithReset
	}

	return nil
}

// IsExceeded проверяет, превышен ли лимит квоты
func (qd QuotaDefinition) IsExceeded(currentUsage decimal.Decimal) bool {
	return currentUsage.Cmp(qd.limit) >= 0
}

// UsagePercentage возвращает процент использования квоты
func (qd QuotaDefinition) UsagePercentage(currentUsage decimal.Decimal) float64 {
	if qd.limit.Cmp(decimal.Zero) == 0 {
		return 0
	}

	// Если текущее использование превышает лимит, возвращаем 100%
	if currentUsage.Cmp(qd.limit) >= 0 {
		return 100.0
	}

	percentage := currentUsage.Div(qd.limit).Mul(decimal.NewFromFloat(100))
	return percentage.InexactFloat64()
}

// NeedsReset проверяет, нужно ли сбросить квоту на основе времени
func (qd QuotaDefinition) NeedsReset(lastReset time.Time) bool {
	if !qd.isRecurring {
		return false
	}

	return time.Since(lastReset) >= qd.resetPeriod
}

// NextResetTime возвращает время следующего сброса квоты
func (qd QuotaDefinition) NextResetTime(lastReset time.Time) time.Time {
	if !qd.isRecurring {
		return time.Time{}
	}

	return lastReset.Add(qd.resetPeriod)
}

// NewQuotaDefinitionForTest создает квоту без валидации (только для тестов)
func NewQuotaDefinitionForTest(
	resourceType string,
	limit decimal.Decimal,
	unit string,
	isRecurring bool,
	resetPeriod time.Duration,
) QuotaDefinition {
	return QuotaDefinition{
		resourceType: resourceType,
		limit:        limit,
		unit:         unit,
		isRecurring:  isRecurring,
		resetPeriod:  resetPeriod,
	}
}
