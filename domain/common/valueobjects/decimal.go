package valueobjects

import (
	"github.com/shopspring/decimal"
)

// NewDecimal создает Decimal из float64 с безопасным преобразованием
func NewDecimal(value float64) decimal.Decimal {
	return decimal.NewFromFloatWithExponent(value, -2)
}

// NewDecimalFromInt создает Decimal из целого числа
func NewDecimalFromInt(value int64) decimal.Decimal {
	return decimal.New(value, 0)
}
