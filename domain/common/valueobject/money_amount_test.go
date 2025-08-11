package valueobject_test

import (
	"testing"

	"github.com/GAKiknadze/payment_service/domain/common/valueobject"
	"github.com/shopspring/decimal"
)

func TestNewMoneyAmount_ValidAmountRUB(t *testing.T) {
	// Given - валюта RUB и валидная сумма
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	amount := decimal.NewFromFloat(100.50)

	// When - создаем MoneyAmount
	moneyAmount, err := valueobject.NewMoneyAmount(amount, currency)

	// Then - проверяем, что ошибка отсутствует и данные корректны
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !moneyAmount.Amount().Equal(amount) {
		t.Errorf("Expected amount %s, got %s", amount, moneyAmount.Amount())
	}

	if moneyAmount.Currency().Code() != "RUB" {
		t.Errorf("Expected currency RUB, got %s", moneyAmount.Currency().Code())
	}
}

func TestNewMoneyAmount_ValidAmountKZT(t *testing.T) {
	// Given - валюта KZT и валидная сумма
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyKZT)
	amount := decimal.NewFromFloat(1500.75)

	// When - создаем MoneyAmount
	moneyAmount, err := valueobject.NewMoneyAmount(amount, currency)

	// Then - проверяем, что ошибка отсутствует
	if err != nil {
		t.Fatalf("Expected no error for KZT, got: %v", err)
	}

	if moneyAmount.Currency().Code() != "KZT" {
		t.Errorf("Expected currency KZT, got %s", moneyAmount.Currency().Code())
	}
}

func TestNewMoneyAmount_InvalidDecimalPlacesRUB(t *testing.T) {
	// Given - валюта RUB и сумма с 3 знаками после запятой
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	amount, _ := decimal.NewFromString("100.555") // 3 знака после запятой для RUB

	// When - пытаемся создать MoneyAmount
	_, err := valueobject.NewMoneyAmount(amount, currency)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for invalid decimal places in RUB, got nil")
	}

	if err != valueobject.ErrInvalidDecimalPlaces {
		t.Errorf("Expected ErrInvalidDecimalPlaces, got %v", err)
	}
}

func TestNewMoneyAmount_InvalidDecimalPlacesKZT(t *testing.T) {
	// Given - валюта KZT и сумма с 3 знаками после запятой
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyKZT)
	amount, _ := decimal.NewFromString("1500.755") // 3 знака после запятой для KZT

	// When - пытаемся создать MoneyAmount
	_, err := valueobject.NewMoneyAmount(amount, currency)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for invalid decimal places in KZT, got nil")
	}

	if err != valueobject.ErrInvalidDecimalPlaces {
		t.Errorf("Expected ErrInvalidDecimalPlaces, got %v", err)
	}
}

func TestEquals_SameAmountSameCurrency(t *testing.T) {
	// Given - две одинаковые суммы одной валюты (RUB)
	currencyRUB, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	amount1, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(100.50), currencyRUB)
	amount2, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(100.50), currencyRUB)

	// When - проверяем равенство
	result := amount1.Equals(amount2)

	// Then - суммы должны быть равны
	if !result {
		t.Error("Expected amounts to be equal")
	}
}

func TestEquals_DifferentCurrencies(t *testing.T) {
	// Given - одинаковые суммы разных валют
	currencyRUB, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	currencyKZT, _ := valueobject.NewCurrency(valueobject.CurrencyKZT)
	amount1, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(100.50), currencyRUB)
	amount2, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(100.50), currencyKZT)

	// When - проверяем равенство
	result := amount1.Equals(amount2)

	// Then - суммы не должны быть равны из-за разной валюты
	if result {
		t.Error("Expected amounts to be different due to different currencies")
	}
}

func TestAdd_SameCurrencyRUB(t *testing.T) {
	// Given - две суммы в RUB
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	amount1, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(100.50), currency)
	amount2, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(50.25), currency)

	// When - складываем суммы
	result, err := amount1.Add(amount2)

	// Then - проверяем результат
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedAmount, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(150.75), currency)
	if !result.Equals(expectedAmount) {
		t.Errorf("Expected %s, got %s", expectedAmount.Format(), result.Format())
	}
}

func TestAdd_SameCurrencyKZT(t *testing.T) {
	// Given - две суммы в KZT
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyKZT)
	amount1, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(1500.50), currency)
	amount2, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(750.25), currency)

	// When - складываем суммы
	result, err := amount1.Add(amount2)

	// Then - проверяем результат
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedAmount, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(2250.75), currency)
	if !result.Equals(expectedAmount) {
		t.Errorf("Expected %s, got %s", expectedAmount.Format(), result.Format())
	}
}

func TestAdd_DifferentCurrencies(t *testing.T) {
	// Given - суммы разных валют
	currencyRUB, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	currencyKZT, _ := valueobject.NewCurrency(valueobject.CurrencyKZT)
	amountRUB, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(100.50), currencyRUB)
	amountKZT, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(1500.75), currencyKZT)

	// When - пытаемся сложить суммы
	_, err := amountRUB.Add(amountKZT)

	// Then - должна быть ошибка из-за разной валюты
	if err == nil {
		t.Fatal("Expected currency mismatch error, got nil")
	}
	if err != valueobject.ErrCurrencyMismatch {
		t.Errorf("Expected ErrCurrencyMismatch, got %v", err)
	}
}

func TestFormat_RUB(t *testing.T) {
	// Given - сумма в рублях
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	amount, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(100.50), currency)

	// When - форматируем сумму
	formatted := amount.Format()

	// Then - проверяем форматирование
	expected := "100.50 ₽"
	if formatted != expected {
		t.Errorf("Expected '%s', got '%s'", expected, formatted)
	}
}

func TestFormat_KZT(t *testing.T) {
	// Given - сумма в тенге
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyKZT)
	amount, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(1500.75), currency)

	// When - форматируем сумму
	formatted := amount.Format()

	// Then - проверяем форматирование
	expected := "1500.75 ₸"
	if formatted != expected {
		t.Errorf("Expected '%s', got '%s'", expected, formatted)
	}
}

func TestCanCover_EnoughFundsRUB(t *testing.T) {
	// Given - сумма, которая может покрыть другую сумму (RUB)
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyRUB)
	balance, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(100.50), currency)
	cost, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(50.25), currency)

	// When - проверяем, может ли баланс покрыть стоимость
	canCover, err := balance.CanCover(cost)

	// Then - проверяем результат
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !canCover {
		t.Error("Expected balance to cover the cost")
	}
}

func TestCanCover_EnoughFundsKZT(t *testing.T) {
	// Given - сумма, которая может покрыть другую сумму (KZT)
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyKZT)
	balance, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(1500.50), currency)
	cost, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(750.25), currency)

	// When - проверяем, может ли баланс покрыть стоимость
	canCover, err := balance.CanCover(cost)

	// Then - проверяем результат
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !canCover {
		t.Error("Expected balance to cover the cost")
	}
}

func TestSubtract_ValidKZT(t *testing.T) {
	// Given - две суммы в KZT, где первая больше второй
	currency, _ := valueobject.NewCurrency(valueobject.CurrencyKZT)
	amount1, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(1500.50), currency)
	amount2, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(750.25), currency)

	// When - вычитаем сумму
	result, err := amount1.Subtract(amount2)

	// Then - проверяем результат
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedAmount, _ := valueobject.NewMoneyAmount(decimal.NewFromFloat(750.25), currency)
	if !result.Equals(expectedAmount) {
		t.Errorf("Expected %s, got %s", expectedAmount.Format(), result.Format())
	}
}
