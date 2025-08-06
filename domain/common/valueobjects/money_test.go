package valueobjects_test

import (
	"testing"

	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Вспомогательные функции для безопасного создания Money в тестах
func mustMoney(amount decimal.Decimal, currency valueobjects.Currency) valueobjects.Money {
	money, err := valueobjects.NewMoney(amount, currency)
	if err != nil {
		panic("invalid test money: " + err.Error())
	}
	return money
}

func mustMoneyFromFloat(amount float64, currency valueobjects.Currency) valueobjects.Money {
	money, err := valueobjects.FromFloat(amount, currency)
	if err != nil {
		panic("invalid test money: " + err.Error())
	}
	return money
}

func TestNewMoney(t *testing.T) {
	rub := valueobjects.RUB

	tests := []struct {
		name        string
		amount      decimal.Decimal
		currency    valueobjects.Currency
		expectedStr string
		expectError bool
	}{
		{
			name:        "Valid positive amount",
			amount:      decimal.NewFromFloat(100.50),
			currency:    rub,
			expectedStr: "100.50 RUB",
			expectError: false,
		},
		{
			name:        "Valid zero amount",
			amount:      decimal.Zero,
			currency:    rub,
			expectedStr: "0.00 RUB",
			expectError: false,
		},
		{
			name:        "Valid amount with more than 2 decimal places",
			amount:      decimal.NewFromFloat(100.555),
			currency:    rub,
			expectedStr: "100.56 RUB",
			expectError: false,
		},
		{
			name:        "Negative amount",
			amount:      decimal.NewFromFloat(-50.25),
			currency:    rub,
			expectedStr: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := valueobjects.NewMoney(tt.amount, tt.currency)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, valueobjects.ErrInvalidMoneyValue, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStr, money.Format())
			}
		})
	}
}

func TestMustMoney(t *testing.T) {
	rub := valueobjects.RUB

	t.Run("Valid amount does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			valueobjects.MustMoney(decimal.NewFromFloat(100.50), rub)
		})
	})

	t.Run("Negative amount panics", func(t *testing.T) {
		assert.Panics(t, func() {
			valueobjects.MustMoney(decimal.NewFromFloat(-50.25), rub)
		})
	})
}

func TestFromFloat(t *testing.T) {
	rub := valueobjects.RUB

	tests := []struct {
		name        string
		amount      float64
		currency    valueobjects.Currency
		expectedStr string
		expectError bool
	}{
		{
			name:        "Valid float amount",
			amount:      100.50,
			currency:    rub,
			expectedStr: "100.50 RUB",
			expectError: false,
		},
		{
			name:        "Float with many decimals",
			amount:      100.555,
			currency:    rub,
			expectedStr: "100.56 RUB",
			expectError: false,
		},
		{
			name:        "Float representation error",
			amount:      0.1 + 0.2,
			currency:    rub,
			expectedStr: "0.30 RUB",
			expectError: false,
		},
		{
			name:        "Negative float",
			amount:      -50.25,
			currency:    rub,
			expectedStr: "",
			expectError: true,
		},
		{
			name:        "Zero float",
			amount:      0.0,
			currency:    rub,
			expectedStr: "0.00 RUB",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := valueobjects.FromFloat(tt.amount, tt.currency)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, valueobjects.ErrInvalidMoneyValue, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStr, money.Format())
			}
		})
	}
}

func TestMoney_Amount(t *testing.T) {
	rub := valueobjects.RUB
	money := mustMoney(decimal.NewFromFloat(100.50), rub)

	assert.True(t, money.Amount().Equal(decimal.NewFromFloat(100.50)))
}

func TestMoney_Currency(t *testing.T) {
	rub := valueobjects.RUB
	money := mustMoney(decimal.NewFromFloat(100.50), rub)

	assert.Equal(t, rub, money.Currency())
}

func TestMoney_IsZero_IsPositive_IsNegative(t *testing.T) {
	rub := valueobjects.RUB
	zeroMoney := mustMoney(decimal.Zero, rub)
	positiveMoney := mustMoney(decimal.NewFromFloat(100.50), rub)

	// IsZero
	assert.True(t, zeroMoney.IsZero())
	assert.False(t, positiveMoney.IsZero())

	// IsPositive
	assert.False(t, zeroMoney.IsPositive())
	assert.True(t, positiveMoney.IsPositive())

	// IsNegative
	assert.False(t, zeroMoney.IsNegative())
	assert.False(t, positiveMoney.IsNegative())
}

func TestMoney_Equals(t *testing.T) {
	rub := valueobjects.RUB
	money1 := mustMoney(decimal.NewFromFloat(100.50), rub)
	money2 := mustMoney(decimal.NewFromFloat(100.50), rub)
	money3 := mustMoney(decimal.NewFromFloat(200.00), rub)

	assert.True(t, money1.Equals(money2))
	assert.False(t, money1.Equals(money3))
	assert.True(t, money1.Equals(money1)) // Сравнение с самим собой
}

func TestMoney_LessThan(t *testing.T) {
	rub := valueobjects.RUB
	money1 := mustMoney(decimal.NewFromFloat(100.50), rub)
	money2 := mustMoney(decimal.NewFromFloat(200.00), rub)
	money3 := mustMoney(decimal.NewFromFloat(100.50), rub)

	assert.True(t, money1.LessThan(money2))
	assert.False(t, money2.LessThan(money1))
	assert.True(t, money1.LessThan(money1)) // LessThan включает равенство
	assert.True(t, money1.LessThan(money3))
}

func TestMoney_Add(t *testing.T) {
	rub := valueobjects.RUB
	money1 := mustMoney(decimal.NewFromFloat(100.50), rub)
	money2 := mustMoney(decimal.NewFromFloat(50.25), rub)
	expected := mustMoney(decimal.NewFromFloat(150.75), rub)

	result, err := money1.Add(money2)
	assert.NoError(t, err)
	assert.Equal(t, expected.Format(), result.Format())
}

func TestMoney_Subtract(t *testing.T) {
	rub := valueobjects.RUB
	money1 := mustMoney(decimal.NewFromFloat(100.50), rub)
	money2 := mustMoney(decimal.NewFromFloat(50.25), rub)
	money3 := mustMoney(decimal.NewFromFloat(150.00), rub)

	t.Run("Subtract resulting in positive amount", func(t *testing.T) {
		result, err := money1.Subtract(money2)
		assert.NoError(t, err)
		assert.Equal(t, "50.25 RUB", result.Format())
	})

	t.Run("Subtract resulting in negative amount", func(t *testing.T) {
		_, err := money2.Subtract(money3)
		assert.Error(t, err)
		assert.Equal(t, "resulting money would be negative", err.Error())
	})
}

func TestMoney_Multiply(t *testing.T) {
	rub := valueobjects.RUB
	money := mustMoney(decimal.NewFromFloat(100.50), rub)

	t.Run("Multiply by factor", func(t *testing.T) {
		result, err := money.Multiply(decimal.NewFromFloat(2.5))
		assert.NoError(t, err)
		assert.Equal(t, "251.25 RUB", result.Format())
	})

	t.Run("Multiply with rounding", func(t *testing.T) {
		money2 := mustMoney(decimal.NewFromFloat(10.00), rub)
		result2, err := money2.Multiply(decimal.NewFromFloat(0.3333))
		assert.NoError(t, err)
		assert.Equal(t, "3.33 RUB", result2.Format())
	})
}

func TestMoney_Divide(t *testing.T) {
	rub := valueobjects.RUB
	money := mustMoney(decimal.NewFromFloat(100.50), rub)

	t.Run("Valid division", func(t *testing.T) {
		result, err := money.Divide(decimal.NewFromFloat(2.0))
		assert.NoError(t, err)
		assert.Equal(t, "50.25 RUB", result.Format())
	})

	t.Run("Division by zero", func(t *testing.T) {
		_, err := money.Divide(decimal.Zero)
		assert.Error(t, err)
		assert.Equal(t, "division by zero", err.Error())
	})
}

func TestMoney_String(t *testing.T) {
	rub := valueobjects.RUB
	money := mustMoney(decimal.NewFromFloat(100.50), rub)

	assert.Equal(t, "100.5 RUB", money.String())
}

func TestMoney_Format(t *testing.T) {
	rub := valueobjects.RUB
	money := mustMoney(decimal.NewFromFloat(100.50), rub)

	assert.Equal(t, "100.50 RUB", money.Format())

	// Проверка форматирования с разным количеством знаков
	money2 := mustMoney(decimal.NewFromFloat(100.555), rub)
	assert.Equal(t, "100.56 RUB", money2.Format())
}

func TestMoney_Negate(t *testing.T) {
	rub := valueobjects.RUB
	money := mustMoney(decimal.NewFromFloat(100.50), rub)
	negated := money.Negate()

	assert.True(t, negated.Amount().Equal(decimal.NewFromFloat(-100.50)))
	assert.Equal(t, rub, negated.Currency())

	// Проверка, что отрицание нуля остается нулем
	zeroMoney := mustMoney(decimal.Zero, rub)
	negatedZero := zeroMoney.Negate()
	assert.True(t, negatedZero.Amount().Equal(decimal.Zero))
}

func TestMoney_Immutability(t *testing.T) {
	rub := valueobjects.RUB
	money := mustMoney(decimal.NewFromFloat(100.50), rub)
	originalString := money.Format()

	// Проверяем, что операции не изменяют исходный объект
	money.Add(mustMoney(decimal.NewFromFloat(50), rub))
	money.Subtract(mustMoney(decimal.NewFromFloat(50), rub))
	money.Multiply(decimal.NewFromFloat(2))
	money.Divide(decimal.NewFromFloat(2))
	money.Negate()

	assert.Equal(t, originalString, money.Format(), "Money должен быть неизменяемым")
}

func TestMoney_Rounding(t *testing.T) {
	rub := valueobjects.RUB

	t.Run("Rounding during creation", func(t *testing.T) {
		money, err := valueobjects.NewMoney(decimal.NewFromFloat(100.555), rub)
		require.NoError(t, err)
		assert.Equal(t, "100.56 RUB", money.Format())
	})

	t.Run("Rounding in operations", func(t *testing.T) {
		money2 := mustMoney(decimal.NewFromFloat(1.03), rub)
		money3 := mustMoney(decimal.NewFromFloat(0.01), rub)

		result, err := money2.Add(money3)
		require.NoError(t, err)
		assert.Equal(t, "1.04 RUB", result.Format())

		// Проверяем округление при умножении
		factor := decimal.NewFromFloat(1.333)
		result2, err := money2.Multiply(factor)
		require.NoError(t, err)
		assert.Equal(t, "1.37 RUB", result2.Format())
	})
}

func TestMoney_EdgeCases(t *testing.T) {
	rub := valueobjects.RUB

	t.Run("Very large amount", func(t *testing.T) {
		largeAmount := decimal.New(1000000000, 0) // 1e9
		largeMoney := mustMoney(largeAmount, rub)

		result, err := largeMoney.Add(largeMoney)
		require.NoError(t, err)
		assert.Equal(t, "2000000000.00 RUB", result.Format())
	})

	t.Run("Very small amount", func(t *testing.T) {
		smallAmount := decimal.New(1, -10) // 0.0000000001
		smallMoney, err := valueobjects.NewMoney(smallAmount, rub)
		require.NoError(t, err)

		// Должно округлиться до 0.00
		assert.Equal(t, "0.00 RUB", smallMoney.Format())

		// Проверка сравнения с нулевой суммой
		zeroMoney := mustMoney(decimal.Zero, rub)
		assert.True(t, smallMoney.Equals(zeroMoney))
	})

	t.Run("Precision with many operations", func(t *testing.T) {
		// Проверяем сохранение точности после цепочки операций
		base := mustMoney(decimal.NewFromFloat(10.00), rub)

		// 10.00 * 0.3333 = 3.33
		result, err := base.Multiply(decimal.NewFromFloat(0.3333))
		require.NoError(t, err)

		// 3.33 * 3 = 9.99
		result, err = result.Multiply(decimal.NewFromFloat(3))
		require.NoError(t, err)

		assert.Equal(t, "9.99 RUB", result.Format())
	})
}
