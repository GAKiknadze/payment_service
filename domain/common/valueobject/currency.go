package valueobject

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

// CurrencyType - enum для кодов валют
type CurrencyType string

const (
	CurrencyRUB CurrencyType = "RUB"
	CurrencyKZT CurrencyType = "KZT"
)

var ErrUnsupportedCurrencyType = errors.New("unsupported currency type")

// Currency - Value Object для работы с валютой
type Currency struct {
	code          CurrencyType
	symbol        string
	decimalPlaces int32
	isSupported   bool
}

// NewCurrency - фабричный метод для создания валюты
func NewCurrency(currencyType CurrencyType) (Currency, error) {
	switch currencyType {
	case CurrencyRUB:
		return Currency{
			code:          CurrencyRUB,
			symbol:        "₽",
			decimalPlaces: 2,
			isSupported:   true,
		}, nil
	case CurrencyKZT:
		return Currency{
			code:          CurrencyKZT,
			symbol:        "₸",
			decimalPlaces: 2,
			isSupported:   false,
		}, nil
	default:
		return Currency{}, ErrUnsupportedCurrencyType
	}
}

func (c Currency) Code() string {
	return string(c.code)
}

func (c Currency) Symbol() string {
	return c.symbol
}

func (c Currency) DecimalPlaces() int32 {
	return c.decimalPlaces
}

func (c Currency) IsSupported() bool {
	return c.isSupported
}

func (c Currency) FormatAmount(amount decimal.Decimal) string {
	// Округляем до нужного количества десятичных знаков
	rounded := amount.Round(c.decimalPlaces)

	// Форматируем число с разделителями тысяч
	formatted := rounded.StringFixed(c.decimalPlaces)

	return fmt.Sprintf("%s %s", formatted, c.symbol)
}

// IsValidAmount - проверяет, является ли сумма валидной для этой валюты
func (c Currency) IsValidAmount(amount decimal.Decimal) bool {
	// Проверяем, что сумма не отрицательная
	if amount.Cmp(decimal.Zero) < 0 {
		return false
	}

	// Проверяем, что сумма точно кратна минимальной единице валюты
	minUnit := decimal.New(1, -c.decimalPlaces)
	remainder := amount.Mod(minUnit)

	return remainder.Equal(decimal.Zero)
}
