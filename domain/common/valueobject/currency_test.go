package valueobject_test

import (
	"testing"

	"github.com/GAKiknadze/payment_service/domain/common/valueobject"
	"github.com/shopspring/decimal"
)

func TestNewCurrency_SupportedCurrency(t *testing.T) {
	// Given - поддерживаемый тип валюты
	currencyType := valueobject.CurrencyRUB

	// When - создаем валюту
	currency, err := valueobject.NewCurrency(currencyType)

	// Then - проверяем, что ошибка отсутствует и данные корректны
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if currency.Code() != "RUB" {
		t.Errorf("Expected code RUB, got %s", currency.Code())
	}

	if currency.Symbol() != "₽" {
		t.Errorf("Expected symbol ₽, got %s", currency.Symbol())
	}

	if currency.DecimalPlaces() != 2 {
		t.Errorf("Expected 2 decimal places, got %d", currency.DecimalPlaces())
	}

	if !currency.IsSupported() {
		t.Error("Expected currency to be supported")
	}
}

func TestNewCurrency_UnsupportedCurrency(t *testing.T) {
	// Given - неподдерживаемый тип валюты
	currencyType := valueobject.CurrencyType("XYZ")

	// When - пытаемся создать валюту
	_, err := valueobject.NewCurrency(currencyType)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for unsupported currency type, got nil")
	}

	if err != valueobject.ErrUnsupportedCurrencyType {
		t.Errorf("Expected ErrUnsupportedCurrencyType, got %v", err)
	}
}

func TestFormatAmount_WholeNumber(t *testing.T) {
	// Given - валюта RUB и целое число
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	amount := decimal.NewFromFloat(100)

	// When - форматируем сумму
	formatted := currency.FormatAmount(amount)

	// Then - проверяем форматирование
	expected := "100.00 ₽"
	if formatted != expected {
		t.Errorf("Expected '%s', got '%s'", expected, formatted)
	}
}

func TestFormatAmount_DecimalNumber(t *testing.T) {
	// Given - валюта RUB и число с дробной частью
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	amount := decimal.NewFromFloat(100.50)

	// When - форматируем сумму
	formatted := currency.FormatAmount(amount)

	// Then - проверяем форматирование
	expected := "100.50 ₽"
	if formatted != expected {
		t.Errorf("Expected '%s', got '%s'", expected, formatted)
	}
}

func TestFormatAmount_RoundingRequired(t *testing.T) {
	// Given - валюта RUB и число, требующее округления
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	amount := decimal.NewFromFloat(100.555)

	// When - форматируем сумму
	formatted := currency.FormatAmount(amount)

	// Then - проверяем, что число округлено до 2 знаков
	expected := "100.56 ₽"
	if formatted != expected {
		t.Errorf("Expected '%s', got '%s'", expected, formatted)
	}
}

func TestFormatAmount_ZeroAmount(t *testing.T) {
	// Given - валюта RUB и нулевое значение
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	amount := decimal.Zero

	// When - форматируем сумму
	formatted := currency.FormatAmount(amount)

	// Then - проверяем форматирование
	expected := "0.00 ₽"
	if formatted != expected {
		t.Errorf("Expected '%s', got '%s'", expected, formatted)
	}
}

func TestIsValidAmount_ValidForRUB(t *testing.T) {
	// Given - валюта RUB и валидная сумма
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	validAmount := decimal.NewFromFloat(100.50) // 2 знака после запятой

	// When - проверяем валидность
	valid := currency.IsValidAmount(validAmount)

	// Then - сумма должна быть валидной
	if !valid {
		t.Error("Expected amount to be valid")
	}
}

func TestIsValidAmount_InvalidForRUB(t *testing.T) {
	// Given - валюта RUB и невалидная сумма
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	invalidAmount := decimal.NewFromFloat(100.555) // 3 знака после запятой

	// When - проверяем валидность
	valid := currency.IsValidAmount(invalidAmount)

	// Then - сумма должна быть невалидной
	if valid {
		t.Error("Expected amount to be invalid")
	}
}

func TestIsValidAmount_NegativeAmount(t *testing.T) {
	// Given - валюта RUB и отрицательная сумма
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	negativeAmount := decimal.NewFromFloat(-100.50)

	// When - проверяем валидность
	valid := currency.IsValidAmount(negativeAmount)

	// Then - отрицательная сумма должна быть невалидной
	if valid {
		t.Error("Expected negative amount to be invalid")
	}
}

func TestIsValidAmount_ZeroAmount(t *testing.T) {
	// Given - валюта RUB и нулевая сумма
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	zeroAmount := decimal.Zero

	// When - проверяем валидность
	valid := currency.IsValidAmount(zeroAmount)

	// Then - нулевая сумма должна быть валидной
	if !valid {
		t.Error("Expected zero amount to be valid")
	}
}

func TestIsValidAmount_BoundaryValidAmount(t *testing.T) {
	// Given - валюта RUB и сумма, равная минимальной единице
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	minUnit := decimal.New(1, -2) // 0.01

	// When - проверяем валидность
	valid := currency.IsValidAmount(minUnit)

	// Then - минимальная единица должна быть валидной
	if !valid {
		t.Error("Expected minimal unit amount to be valid")
	}
}
