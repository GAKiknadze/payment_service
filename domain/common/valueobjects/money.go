package valueobjects

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidMoneyValue = errors.New("invalid money value: must be non-negative")
	ErrCurrencyMismatch  = errors.New("currency mismatch in money operation")
)

const scale = 10000

// Money представляет денежную сумму с точностью до копеек
// Value Object: неизменяемый, сравнивается по значению
type Money struct {
	amount   decimal.Decimal
	currency Currency
}

// NewMoney создает новый объект Money с валидацией
func NewMoney(amount decimal.Decimal, currency Currency) (Money, error) {
	if amount.LessThan(decimal.Zero) {
		return Money{}, ErrInvalidMoneyValue
	}

	// Округляем до 2 знаков после запятой (для большинства валют)
	rounded := amount.Round(2)
	return Money{
		amount:   rounded,
		currency: currency,
	}, nil
}

// MustMoney создает Money или паникует при ошибке (для констант)
func MustMoney(amount decimal.Decimal, currency Currency) Money {
	money, err := NewMoney(amount, currency)
	if err != nil {
		panic(err)
	}
	return money
}

// FromFloat создает Money из float64 с безопасным преобразованием
func FromFloat(amount float64, currency Currency) (Money, error) {
	// Избегаем проблем с плавающей точкой через строковое представление
	str := fmt.Sprintf("%.10f", amount)
	d, err := decimal.NewFromString(str)
	if err != nil {
		return Money{}, err
	}
	return NewMoney(d, currency)
}

// Amount возвращает внутреннее значение суммы
func (m Money) Amount() decimal.Decimal {
	return m.amount
}

// Currency возвращает валюту
func (m Money) Currency() Currency {
	return m.currency
}

// IsZero проверяет, является ли сумма нулевой
func (m Money) IsZero() bool {
	return m.amount.Equal(decimal.Zero)
}

// IsPositive проверяет, является ли сумма положительной
func (m Money) IsPositive() bool {
	return m.amount.GreaterThan(decimal.Zero)
}

// IsNegative проверяет, является ли сумма отрицательной
func (m Money) IsNegative() bool {
	return m.amount.LessThan(decimal.Zero)
}

// Equals сравнивает два Money по значению и валюте
func (m Money) Equals(other Money) bool {
	return m.currency == other.currency && m.amount.Equal(other.amount)
}

// IsApproximatelyEqual проверяет, равны ли два денежных значения с учетом допустимой погрешности
func (m Money) IsApproximatelyEqual(other Money, epsilon float64) bool {
	epsilonStr := strconv.FormatFloat(epsilon, 'f', -1, 64)
	epsilonDecimal, _ := decimal.NewFromString(epsilonStr)

	diff := m.amount.Sub(other.amount).Abs()
	return diff.Cmp(epsilonDecimal) <= 0
}

// LessThan
func (m Money) LessThan(other Money) bool {
	return m.amount.LessThanOrEqual(other.amount)
}

// Add складывает две суммы (только одинаковая валюта)
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	return NewMoney(m.amount.Add(other.amount), m.currency)
}

// Subtract вычитает сумму (только одинаковая валюта)
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	result := m.amount.Sub(other.amount)
	if result.LessThan(decimal.Zero) {
		return Money{}, errors.New("resulting money would be negative")
	}
	return NewMoney(result, m.currency)
}

// Multiply умножает сумму на коэффициент
func (m Money) Multiply(factor decimal.Decimal) (Money, error) {
	return NewMoney(m.amount.Mul(factor), m.currency)
}

// Divide делит сумму на коэффициент
func (m Money) Divide(divisor decimal.Decimal) (Money, error) {
	if divisor.Equal(decimal.Zero) {
		return Money{}, errors.New("division by zero")
	}
	return NewMoney(m.amount.Div(divisor), m.currency)
}

// String возвращает строковое представление
func (m Money) String() string {
	return fmt.Sprintf("%s %s", m.amount.String(), m.currency)
}

// Format возвращает форматированное строковое представление
func (m Money) Format() string {
	amountStr := m.amount.StringFixed(2)
	return fmt.Sprintf("%s %s", amountStr, m.currency)
}

// Negate возвращает отрицательное значение суммы
func (m Money) Negate() Money {
	return Money{
		amount:   m.amount.Neg(),
		currency: m.currency,
	}
}
