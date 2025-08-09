package model_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
	"github.com/GAKiknadze/payment_service/domain/organization/model"
)

// MockTariff для тестирования тарифов
type MockTariff struct {
	mock.Mock
}

func (m *MockTariff) CalculateCost(duration time.Duration) (valueobjects.Money, error) {
	args := m.Called(duration)
	return args.Get(0).(valueobjects.Money), args.Error(1)
}

// Генерация тестовых данных
func newTestOrganizationID() valueobjects.OrganizationID {
	return valueobjects.OrganizationID(uuid.New().String())
}

func newTestMoney(amount float64) valueobjects.Money {
	return valueobjects.NewMoneyFromFloat(amount)
}

func newTestDateTime(year int, month time.Month, day, hour, min, sec int) valueobjects.DateTime {
	return valueobjects.NewDateTime(time.Date(year, month, day, hour, min, sec, 0, time.UTC))
}

// Тесты пополнения баланса
func TestDeposit(t *testing.T) {
	orgID := newTestOrganizationID()
	clock := &MockClock{}
	now := newTestDateTime(2025, 1, 1, 10, 0, 0)
	clock.On("Now").Return(now)

	tariff := &MockTariff{}

	org, _ := model.NewOrganization(
		orgID,
		"Test Org",
		newTestMoney(1000),
		tariff,
		clock,
	)

	t.Run("valid deposit", func(t *testing.T) {
		// Сбрасываем события перед тестом
		org.PopEvents()

		err := org.Deposit(newTestMoney(500), clock)
		assert.NoError(t, err)
		assert.Equal(t, newTestMoney(1500), org.Balance())

		// Проверяем сгенерированное событие
		events := org.PopEvents()
		assert.Len(t, events, 1)

		balanceUpdated, ok := events[0].(model.BalanceUpdated)
		assert.True(t, ok)
		assert.Equal(t, orgID, balanceUpdated.OrganizationID)
		assert.Equal(t, newTestMoney(1000), balanceUpdated.PreviousBalance)
		assert.Equal(t, newTestMoney(1500), balanceUpdated.NewBalance)
		assert.Equal(t, newTestMoney(500), balanceUpdated.ChangeAmount)
		assert.Equal(t, "deposit", balanceUpdated.ChangeType)
		assert.Equal(t, now, balanceUpdated.Timestamp)
	})

	t.Run("negative amount", func(t *testing.T) {
		// Сбрасываем события
		org.PopEvents()

		err := org.Deposit(newTestMoney(-100), clock)
		assert.Error(t, err)
		assert.Equal(t, newTestMoney(1500), org.Balance())
		assert.Empty(t, org.PopEvents())
	})

	t.Run("zero amount", func(t *testing.T) {
		// Сбрасываем события
		org.PopEvents()

		err := org.Deposit(newTestMoney(0), clock)
		assert.Error(t, err)
		assert.Equal(t, newTestMoney(1500), org.Balance())
		assert.Empty(t, org.PopEvents())
	})

	t.Run("deposit to suspended org", func(t *testing.T) {
		// Симулируем приостановку
		org.suspendOrganization(clock)

		// Сбрасываем события
		org.PopEvents()

		err := org.Deposit(newTestMoney(100), clock)
		assert.ErrorIs(t, err, model.ErrOrganizationSuspended)
		assert.Equal(t, newTestMoney(0), org.Balance()) // Баланс не изменился
		assert.Empty(t, org.PopEvents())
	})

	t.Run("multiple deposits", func(t *testing.T) {
		// Сбрасываем события и баланс
		org, _ = model.NewOrganization(
			orgID,
			"Test Org",
			newTestMoney(0),
			tariff,
			clock,
		)

		depositAmount := newTestMoney(250)
		expectedBalance := newTestMoney(0)

		// Выполняем несколько пополнений
		for i := 0; i < 4; i++ {
			err := org.Deposit(depositAmount, clock)
			assert.NoError(t, err)
			expectedBalance = expectedBalance.Add(depositAmount)
			assert.Equal(t, expectedBalance, org.Balance())
		}

		assert.Equal(t, newTestMoney(1000), org.Balance())
	})
}

// Тесты проверки баланса
func TestCheckBalance(t *testing.T) {
	orgID := newTestOrganizationID()
	clock := &MockClock{}
	now := newTestDateTime(2025, 1, 1, 10, 0, 0)
	clock.On("Now").Return(now)

	tariff := &MockTariff{}

	org, _ := model.NewOrganization(
		orgID,
		"Test Org",
		newTestMoney(1000),
		tariff,
		clock,
	)

	t.Run("sufficient balance", func(t *testing.T) {
		assert.True(t, org.CheckBalance(newTestMoney(500)))
		assert.True(t, org.CheckBalance(newTestMoney(1000)))
	})

	t.Run("insufficient balance", func(t *testing.T) {
		assert.False(t, org.CheckBalance(newTestMoney(1500)))
		assert.False(t, org.CheckBalance(newTestMoney(1000.01)))
	})

	t.Run("exact balance", func(t *testing.T) {
		assert.True(t, org.CheckBalance(newTestMoney(1000)))
	})

	t.Run("zero amount check", func(t *testing.T) {
		assert.True(t, org.CheckBalance(newTestMoney(0)))
	})

	t.Run("negative amount check", func(t *testing.T) {
		assert.Panics(t, func() {
			org.CheckBalance(newTestMoney(-100))
		})
	})

	t.Run("suspended organization", func(t *testing.T) {
		// Приостанавливаем организацию
		org.suspendOrganization(clock)

		assert.False(t, org.CheckBalance(newTestMoney(500)))
		assert.False(t, org.CheckBalance(newTestMoney(0)))
	})

	t.Run("active after resume", func(t *testing.T) {
		// Приостанавливаем организацию
		org.suspendOrganization(clock)

		// Пополняем и возобновляем
		org.Deposit(newTestMoney(500), clock)
		org.Resume(clock)

		assert.True(t, org.CheckBalance(newTestMoney(250)))
	})
}

// Тесты для edge cases с балансом
func TestBalanceEdgeCases(t *testing.T) {
	orgID := newTestOrganizationID()
	clock := &MockClock{}
	now := newTestDateTime(2025, 1, 1, 10, 0, 0)
	clock.On("Now").Return(now)

	tariff := &MockTariff{}
	tariff.On("CalculateCost", mock.Anything).Return(newTestMoney(100), nil)

	t.Run("balance precision", func(t *testing.T) {
		// Создаем организацию с балансом, требующим высокой точности
		org, _ := model.NewOrganization(
			orgID,
			"Precision Test",
			newTestMoney(1000.01),
			tariff,
			clock,
		)

		// Проверяем, что точность сохраняется
		assert.Equal(t, newTestMoney(1000.01), org.Balance())

		// Выполняем операции с дробными суммами
		org.Deposit(newTestMoney(0.01), clock)
		assert.Equal(t, newTestMoney(1000.02), org.Balance())

		billingPeriod, _ := tariff_model.NewTariffPeriod(
			newTestDateTime(2025, 1, 1, 0, 0, 0).Time(),
			newTestDateTime(2025, 1, 1, 10, 0, 0).Time(),
		)
		org.ProcessBilling(billingPeriod, clock)

		// Проверяем, что точность сохраняется после списания
		assert.Equal(t, newTestMoney(900.02), org.Balance())
	})

	t.Run("large amounts", func(t *testing.T) {
		// Создаем организацию с большим балансом
		org, _ := model.NewOrganization(
			orgID,
			"Large Amount Test",
			newTestMoney(1000000000.0),
			tariff,
			clock,
		)

		// Проверяем операции с большими суммами
		org.Deposit(newTestMoney(500000000.0), clock)
		assert.Equal(t, newTestMoney(1500000000.0), org.Balance())

		billingPeriod, _ := tariff_model.NewTariffPeriod(
			newTestDateTime(2025, 1, 1, 0, 0, 0).Time(),
			newTestDateTime(2025, 1, 1, 10, 0, 0).Time(),
		)
		org.ProcessBilling(billingPeriod, clock)

		assert.Equal(t, newTestMoney(1499999900.0), org.Balance())
	})

	t.Run("zero balance transitions", func(t *testing.T) {
		org, _ := model.NewOrganization(
			orgID,
			"Zero Balance Test",
			newTestMoney(100),
			tariff,
			clock,
		)

		billingPeriod, _ := tariff_model.NewTariffPeriod(
			newTestDateTime(2025, 1, 1, 0, 0, 0).Time(),
			newTestDateTime(2025, 1, 1, 10, 0, 0).Time(),
		)

		// Списываем до нуля
		org.ProcessBilling(billingPeriod, clock)
		assert.Equal(t, newTestMoney(0), org.Balance())
		assert.True(t, org.IsSuspended())

		// Пополняем до 1
		org.Deposit(newTestMoney(1), clock)
		assert.Equal(t, newTestMoney(1), org.Balance())
		assert.True(t, org.IsSuspended()) // Все еще приостановлена

		// Возобновляем
		org.Resume(clock)
		assert.False(t, org.IsSuspended())

		// Снова списываем до нуля
		org.ProcessBilling(billingPeriod, clock)
		assert.Equal(t, newTestMoney(0), org.Balance())
		assert.True(t, org.IsSuspended())
	})

	t.Run("version increments on deposit", func(t *testing.T) {
		org, _ := model.NewOrganization(
			orgID,
			"Version Test",
			newTestMoney(1000),
			tariff,
			clock,
		)

		initialVersion := org.Version()

		// Пополняем баланс
		org.Deposit(newTestMoney(500), clock)
		assert.Equal(t, initialVersion+1, org.Version())

		// Еще одно пополнение
		org.Deposit(newTestMoney(250), clock)
		assert.Equal(t, initialVersion+2, org.Version())
	})
}
