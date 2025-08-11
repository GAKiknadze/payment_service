package valueobject

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidAmount        = errors.New("amount must be non-negative")
	ErrCurrencyMismatch     = errors.New("currency mismatch")
	ErrInvalidDecimalPlaces = errors.New("amount has invalid decimal places for currency")
)

type MoneyAmount struct {
	amount   decimal.Decimal
	currency Currency
}

// NewMoneyAmount создает новый объект MoneyAmount с проверкой валидности
func NewMoneyAmount(amount decimal.Decimal, currency Currency) (MoneyAmount, error) {
	// Проверяем, что сумма не отрицательная
	if amount.Cmp(decimal.Zero) < 0 {
		return MoneyAmount{}, ErrInvalidAmount
	}

	// Проверяем соответствие минимальной единице валюты
	if !currency.IsValidAmount(amount) {
		return MoneyAmount{}, ErrInvalidDecimalPlaces
	}

	return MoneyAmount{
		amount:   amount,
		currency: currency,
	}, nil
}

// Amount возвращает сумму как decimal.Decimal
func (ma MoneyAmount) Amount() decimal.Decimal {
	return ma.amount
}

// Currency возвращает валюту суммы
func (ma MoneyAmount) Currency() Currency {
	return ma.currency
}

// Format возвращает отформатированное строковое представление суммы
func (ma MoneyAmount) Format() string {
	return ma.currency.FormatAmount(ma.amount)
}

// IsValid проверяет, является ли сумма валидной
func (m MoneyAmount) IsValid() bool {
	return m.amount.Cmp(decimal.Zero) >= 0 && m.currency.IsValidAmount(m.amount)
}

// Equals проверяет равенство двух MoneyAmount
func (m MoneyAmount) Equals(other MoneyAmount) bool {
	return m.currency.Code() == other.currency.Code() &&
		m.amount.Equal(other.amount)
}

// GreaterThan проверяет, больше ли текущая сумма, чем другая
func (m MoneyAmount) GreaterThan(other MoneyAmount) (bool, error) {
	if m.currency.Code() != other.currency.Code() {
		return false, ErrCurrencyMismatch
	}
	return m.amount.Cmp(other.amount) > 0, nil
}

// GreaterThanOrEqual проверяет, больше или равна ли текущая сумма другой
func (m MoneyAmount) GreaterThanOrEqual(other MoneyAmount) (bool, error) {
	if m.currency.Code() != other.currency.Code() {
		return false, ErrCurrencyMismatch
	}
	return m.amount.Cmp(other.amount) >= 0, nil
}

// LessThan проверяет, меньше ли текущая сумма, чем другая
func (m MoneyAmount) LessThan(other MoneyAmount) (bool, error) {
	if m.currency.Code() != other.currency.Code() {
		return false, ErrCurrencyMismatch
	}
	return m.amount.Cmp(other.amount) < 0, nil
}

// Add складывает две суммы одной валюты
func (m MoneyAmount) Add(other MoneyAmount) (MoneyAmount, error) {
	if m.currency.Code() != other.currency.Code() {
		return MoneyAmount{}, ErrCurrencyMismatch
	}

	newAmount := m.amount.Add(other.amount)
	return NewMoneyAmount(newAmount, m.currency)
}

// Subtract вычитает другую сумму из текущей
func (m MoneyAmount) Subtract(other MoneyAmount) (MoneyAmount, error) {
	if m.currency.Code() != other.currency.Code() {
		return MoneyAmount{}, ErrCurrencyMismatch
	}

	newAmount := m.amount.Sub(other.amount)
	if newAmount.Cmp(decimal.Zero) < 0 {
		return MoneyAmount{}, ErrInvalidAmount
	}

	return NewMoneyAmount(newAmount, m.currency)
}

// CanCover проверяет, может ли текущая сумма покрыть другую сумму
func (m MoneyAmount) CanCover(other MoneyAmount) (bool, error) {
	if m.currency.Code() != other.currency.Code() {
		return false, ErrCurrencyMismatch
	}

	return m.amount.Cmp(other.amount) >= 0, nil
}

// Для тестов в обход проверок
func NewMoneyAmountForTest(amount decimal.Decimal, currency Currency) MoneyAmount {
	return MoneyAmount{
		amount:   amount,
		currency: currency,
	}
}
