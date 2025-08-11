package valueobject_test

import (
	"testing"

	"github.com/GAKiknadze/payment_service/domain/common/valueobject"
	"github.com/shopspring/decimal"
)

func TestNewPrice_ValidPrice(t *testing.T) {
	// Given - валидная сумма и валюта RUB
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	amount, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(100.50), currency)

	// When - создаем Price
	price, err := valueobject.NewPrice("price_123", amount, true)

	// Then - проверяем, что ошибка отсутствует и данные корректны
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if price.ID() != "price_123" {
		t.Errorf("Expected ID price_123, got %s", price.ID())
	}

	if !price.Amount().Equals(amount) {
		t.Error("Expected amounts to be equal")
	}

	if price.IsDefault() != true {
		t.Error("Expected price to be default")
	}
}

func TestNewPrice_InvalidMoneyAmount(t *testing.T) {
	// Given - поддерживаемая валюта
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)

	// Создаем MoneyAmount напрямую для тестирования
	invalidAmount := valueobject.NewMoneyAmountForTest(decimal.NewFromFloat(-100.50), currency)

	// When - пытаемся создать Price
	_, err := valueobject.NewPrice("price_123", invalidAmount, true)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for invalid money amount, got nil")
	}

	if err != valueobject.ErrInvalidPrice {
		t.Errorf("Expected ErrInvalidPrice, got %v", err)
	}
}

func TestNewPrice_InvalidDecimalPlaces(t *testing.T) {
	// Given - сумма с неверным количеством знаков для RUB
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	amount := valueobject.NewMoneyAmountForTest(decimal.NewFromFloat(100.555), currency) // 3 знака после запятой

	// When - пытаемся создать Price
	_, err := valueobject.NewPrice("price_123", amount, true)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for invalid decimal places, got nil")
	}

	if err != valueobject.ErrInvalidPrice {
		t.Errorf("Expected ErrInvalidPrice, got %v", err)
	}
}

func TestAmount_ReturnsCorrectMoneyAmount(t *testing.T) {
	// Given - валидная цена
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	originalAmount, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(100.50), currency)
	price, _ := valueobject.NewPrice("price_123", originalAmount, true)

	// When - получаем сумму через метод Amount()
	retrievedAmount := price.Amount()

	// Then - проверяем, что возвращается корректная MoneyAmount
	if !retrievedAmount.Equals(originalAmount) {
		t.Error("Expected retrieved amount to equal original amount")
	}
}

func TestCurrency_ReturnsCorrectCurrency(t *testing.T) {
	// Given - цена в RUB
	currencyRUB, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	amountRUB, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(100.50), currencyRUB)
	priceRUB, _ := valueobject.NewPrice("price_rub", amountRUB, true)

	// Given - цена в KZT
	currencyKZT, _ := valueobject.NewCurrency(valueobject.CurrencyKZT)
	amountKZT, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(1500.75), currencyKZT)
	priceKZT, _ := valueobject.NewPrice("price_kzt", amountKZT, true)

	// When & Then - проверяем валюту для RUB
	if priceRUB.Currency().Code() != "RUB" {
		t.Errorf("Expected RUB currency, got %s", priceRUB.Currency().Code())
	}

	// When & Then - проверяем валюту для KZT
	if priceKZT.Currency().Code() != "KZT" {
		t.Errorf("Expected KZT currency, got %s", priceKZT.Currency().Code())
	}
}

func TestFormat_ReturnsCorrectString(t *testing.T) {
	// Given - цена в RUB
	currencyRUB, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	amountRUB, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(100.50), currencyRUB)
	priceRUB, _ := valueobject.NewPrice("price_rub", amountRUB, true)

	// Given - цена в KZT
	currencyKZT, _ := valueobject.NewCurrency(valueobject.CurrencyKZT)
	amountKZT, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(1500.75), currencyKZT)
	priceKZT, _ := valueobject.NewPrice("price_kzt", amountKZT, true)

	// When & Then - проверяем форматирование для RUB
	expectedRUB := "100.50 ₽"
	if priceRUB.Format() != expectedRUB {
		t.Errorf("Expected '%s', got '%s'", expectedRUB, priceRUB.Format())
	}

	// When & Then - проверяем форматирование для KZT
	expectedKZT := "1500.75 ₸"
	if priceKZT.Format() != expectedKZT {
		t.Errorf("Expected '%s', got '%s'", expectedKZT, priceKZT.Format())
	}
}

func TestIsCompatibleWith_SameCurrency(t *testing.T) {
	// Given - цена в RUB и организация с валютой RUB
	currencyRUB, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	amountRUB, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(100.50), currencyRUB)
	price, _ := valueobject.NewPrice("price_rub", amountRUB, true)

	// When - проверяем совместимость
	compatible := price.IsCompatibleWith(currencyRUB)

	// Then - должна быть совместимой
	if !compatible {
		t.Error("Expected price to be compatible with same currency")
	}
}

func TestIsCompatibleWith_DifferentCurrency(t *testing.T) {
	// Given - цена в RUB и организация с валютой KZT
	currencyRUB, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	amountRUB, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(100.50), currencyRUB)
	price, _ := valueobject.NewPrice("price_rub", amountRUB, true)

	currencyKZT, _ := valueobject.NewCurrency(valueobject.CurrencyKZT)

	// When - проверяем совместимость
	compatible := price.IsCompatibleWith(currencyKZT)

	// Then - не должна быть совместимой
	if compatible {
		t.Error("Expected price not to be compatible with different currency")
	}
}

func TestIsDefault_Flag(t *testing.T) {
	// Given - цена с isDefault = true
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	amount, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(100.50), currency)
	defaultPrice, _ := valueobject.NewPrice("price_default", amount, true)

	// Given - цена с isDefault = false
	nonDefaultPrice, _ := valueobject.NewPrice("price_non_default", amount, false)

	// When & Then - проверяем флаг isDefault для defaultPrice
	if !defaultPrice.IsDefault() {
		t.Error("Expected default price to have isDefault = true")
	}

	// When & Then - проверяем флаг isDefault для nonDefaultPrice
	if nonDefaultPrice.IsDefault() {
		t.Error("Expected non-default price to have isDefault = false")
	}
}
