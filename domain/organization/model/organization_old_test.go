// domain/organization/model/organization_test.go
package model_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClock для тестирования времени
type MockClock struct {
	mock.Mock
}

func (m *MockClock) Now() valueobjects.DateTime {
	args := m.Called()
	return args.Get(0).(valueobjects.DateTime)
}

// MockTariff для тестирования тарифов
type MockTariff struct {
	mock.Mock
}

func (m *MockTariff) CalculateCost(duration time.Duration) (valueobjects.Money, error) {
	args := m.Called(duration)
	return args.Get(0).(valueobjects.Money), args.Error(1)
}

func (m *MockTariff) ID() valueobjects.TariffID {
	args := m.Called()
	return args.Get(0).(valueobjects.TariffID)
}

// Генерация тестовых данных
func newTestOrganizationID() valueobjects.OrganizationID {
	return valueobjects.OrganizationID(uuid.New().String())
}

func newTestTariffID() valueobjects.TariffID {
	return valueobjects.TariffID(uuid.New().String())
}

func newTestMoney(amount float64) valueobjects.Money {
	return valueobjects.NewMoneyFromFloat(amount)
}

func newTestDateTime(year int, month time.Month, day, hour, min, sec int) valueobjects.DateTime {
	return valueobjects.NewDateTime(time.Date(year, month, day, hour, min, sec, 0, time.UTC))
}

// Тесты создания организации
func TestNewOrganization(t *testing.T) {
	t.Run("valid creation", func(t *testing.T) {
		orgID := newTestOrganizationID()
		tariffID := newTestTariffID()
		clock := &MockClock{}
		now := newTestDateTime(2025, 1, 1, 10, 0, 0)
		clock.On("Now").Return(now)

		tariff := &MockTariff{}
		tariff.On("ID").Return(tariffID)

		org, err := model.NewOrganization(
			orgID,
			"Test Org",
			newTestMoney(1000),
			tariff,
			clock,
		)

		assert.NoError(t, err)
		assert.Equal(t, orgID, org.ID())
		assert.Equal(t, "Test Org", org.Name())
		assert.Equal(t, newTestMoney(1000), org.Balance())
		assert.Equal(t, model.StatusActive, org.Status())
		assert.False(t, org.IsSuspended())
		assert.Equal(t, now, org.BillingInfo().LastBillingTime)
		assert.Equal(t, newTestDateTime(2025, 2, 1, 10, 0, 0), org.BillingInfo().NextBillingTime)

		// Проверяем, что событий еще нет (только после операций)
		assert.Empty(t, org.PopEvents())
	})

	t.Run("empty name", func(t *testing.T) {
		org, err := model.NewOrganization(
			newTestOrganizationID(),
			"",
			newTestMoney(1000),
			&MockTariff{},
			&MockClock{},
		)
		assert.ErrorIs(t, err, model.ErrInvalidOrganizationName)
		assert.Nil(t, org)
	})

	t.Run("negative balance", func(t *testing.T) {
		org, err := model.NewOrganization(
			newTestOrganizationID(),
			"Test Org",
			newTestMoney(-100),
			&MockTariff{},
			&MockClock{},
		)
		assert.Error(t, err)
		assert.Nil(t, org)
	})

	t.Run("nil tariff", func(t *testing.T) {
		org, err := model.NewOrganization(
			newTestOrganizationID(),
			"Test Org",
			newTestMoney(1000),
			nil,
			&MockClock{},
		)
		assert.Error(t, err)
		assert.Nil(t, org)
	})
}

// Тесты списания
func TestProcessBilling(t *testing.T) {
	orgID := newTestOrganizationID()
	clock := &MockClock{}
	now := newTestDateTime(2025, 1, 1, 10, 0, 0)
	clock.On("Now").Return(now)

	tariff := &MockTariff{}
	tariff.On("CalculateCost", mock.Anything).Return(newTestMoney(100), nil)

	org, _ := model.NewOrganization(
		orgID,
		"Test Org",
		newTestMoney(1000),
		tariff,
		clock,
	)

	billingPeriod, _ := tariff_model.NewTariffPeriod(
		newTestDateTime(2025, 1, 1, 0, 0, 0).Time(),
		newTestDateTime(2025, 1, 1, 10, 0, 0).Time(),
	)

	t.Run("valid billing", func(t *testing.T) {
		// Сбрасываем события перед тестом
		org.PopEvents()

		err := org.ProcessBilling(billingPeriod, clock)
		assert.NoError(t, err)
		assert.Equal(t, newTestMoney(900), org.Balance())

		// Проверяем события
		events := org.PopEvents()
		assert.Len(t, events, 2)

		// Событие обновления баланса
		balanceUpdated, ok := events[0].(model.BalanceUpdated)
		assert.True(t, ok)
		assert.Equal(t, newTestMoney(1000), balanceUpdated.PreviousBalance)
		assert.Equal(t, newTestMoney(900), balanceUpdated.NewBalance)
		assert.Equal(t, newTestMoney(-100), balanceUpdated.ChangeAmount)

		// Событие списания
		billingProcessed, ok := events[1].(model.BillingProcessed)
		assert.True(t, ok)
		assert.Equal(t, newTestMoney(100), billingProcessed.Amount)
		assert.Equal(t, billingPeriod.Start(), billingProcessed.PeriodStart)
		assert.Equal(t, billingPeriod.End(), billingProcessed.PeriodEnd)
	})

	t.Run("insufficient balance", func(t *testing.T) {
		// Создаем организацию с небольшим балансом
		org, _ := model.NewOrganization(
			orgID,
			"Test Org",
			newTestMoney(50),
			tariff,
			clock,
		)

		err := org.ProcessBilling(billingPeriod, clock)
		assert.ErrorIs(t, err, model.ErrInsufficientBalance)
		assert.Equal(t, newTestMoney(50), org.Balance())
		assert.Empty(t, org.PopEvents())
	})

	t.Run("zero duration period", func(t *testing.T) {
		zeroPeriod, _ := tariff_model.NewTariffPeriod(
			now.Time(),
			now.Time(),
		)

		err := org.ProcessBilling(zeroPeriod, clock)
		assert.ErrorIs(t, err, model.ErrInvalidBillingPeriod)
	})

	t.Run("suspended organization", func(t *testing.T) {
		// Симулируем приостановку
		org.suspendOrganization(clock)

		err := org.ProcessBilling(billingPeriod, clock)
		assert.ErrorIs(t, err, model.ErrOrganizationSuspended)
	})
}

// Тесты приостановки и возобновления
func TestSuspendAndResume(t *testing.T) {
	orgID := newTestOrganizationID()
	clock := &MockClock{}
	now := newTestDateTime(2025, 1, 1, 10, 0, 0)
	clock.On("Now").Return(now)

	tariff := &MockTariff{}
	tariff.On("CalculateCost", mock.Anything).Return(newTestMoney(1000), nil)

	org, _ := model.NewOrganization(
		orgID,
		"Test Org",
		newTestMoney(1000),
		tariff,
		clock,
	)

	billingPeriod, _ := tariff_model.NewTariffPeriod(
		newTestDateTime(2025, 1, 1, 0, 0, 0).Time(),
		newTestDateTime(2025, 1, 1, 10, 0, 0).Time(),
	)

	t.Run("auto-suspend on zero balance", func(t *testing.T) {
		// Сбрасываем события
		org.PopEvents()

		// Выполняем списание, которое приведет к нулевому балансу
		err := org.ProcessBilling(billingPeriod, clock)
		assert.NoError(t, err)
		assert.Equal(t, newTestMoney(0), org.Balance())
		assert.True(t, org.IsSuspended())
		assert.Equal(t, model.StatusSuspended, org.Status())

		// Проверяем событие приостановки
		events := org.PopEvents()
		assert.Len(t, events, 3) // BalanceUpdated, BillingProcessed, OrganizationSuspended

		suspendedEvent, ok := events[2].(model.OrganizationSuspended)
		assert.True(t, ok)
		assert.Equal(t, orgID, suspendedEvent.OrganizationID)
		assert.Equal(t, now, suspendedEvent.Timestamp)
	})

	t.Run("resume with positive balance", func(t *testing.T) {
		// Пополняем баланс
		org.Deposit(newTestMoney(500), clock)

		// Сбрасываем события
		org.PopEvents()

		err := org.Resume(clock)
		assert.NoError(t, err)
		assert.False(t, org.IsSuspended())
		assert.Equal(t, model.StatusActive, org.Status())

		// Проверяем событие возобновления
		events := org.PopEvents()
		assert.Len(t, events, 1)

		resumedEvent, ok := events[0].(model.OrganizationResumed)
		assert.True(t, ok)
		assert.Equal(t, orgID, resumedEvent.OrganizationID)
	})

	t.Run("resume with zero balance", func(t *testing.T) {
		// Создаем приостановленную организацию с нулевым балансом
		org, _ := model.NewOrganization(
			orgID,
			"Test Org",
			newTestMoney(0),
			tariff,
			clock,
		)
		org.suspendOrganization(clock)

		err := org.Resume(clock)
		assert.Error(t, err)
		assert.True(t, org.IsSuspended())
	})
}

// Тесты окончательного завершения
func TestTerminate(t *testing.T) {
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

	t.Run("terminate active organization", func(t *testing.T) {
		// Сбрасываем события
		org.PopEvents()

		err := org.Terminate(clock)
		assert.NoError(t, err)
		assert.True(t, org.IsSuspended())
		assert.Equal(t, model.StatusTerminated, org.Status())

		// Проверяем событие
		events := org.PopEvents()
		assert.Len(t, events, 1)

		terminatedEvent, ok := events[0].(model.OrganizationTerminated)
		assert.True(t, ok)
		assert.Equal(t, orgID, terminatedEvent.OrganizationID)
		assert.Equal(t, newTestMoney(1000), terminatedEvent.FinalBalance)
	})

	t.Run("terminate already terminated", func(t *testing.T) {
		// Сначала завершаем организацию
		org.Terminate(clock)

		err := org.Terminate(clock)
		assert.Error(t, err)
	})
}

// Тесты расчета следующего времени списания
func TestCalculateNextBillingTime(t *testing.T) {
	clock := &MockClock{}
	now := newTestDateTime(2025, 1, 1, 10, 0, 0)
	clock.On("Now").Return(now)

	tariff := &MockTariff{}

	org, _ := model.NewOrganization(
		newTestOrganizationID(),
		"Test Org",
		newTestMoney(1000),
		tariff,
		clock,
	)

	t.Run("last day of month", func(t *testing.T) {
		// 31 января
		lastBilling := newTestDateTime(2025, 1, 31, 10, 0, 0)
		org.SetBillingInfo(model.BillingInfo{
			LastBillingTime: lastBilling,
			NextBillingTime: lastBilling,
		})

		// Должно стать 28 февраля (так как 31 февраля не существует)
		expected := newTestDateTime(2025, 2, 28, 10, 0, 0)
		nextBilling := model.CalculateNextBillingTime(lastBilling, tariff)
		assert.Equal(t, expected, nextBilling)
	})

	t.Run("30th day of month", func(t *testing.T) {
		// 30 апреля
		lastBilling := newTestDateTime(2025, 4, 30, 10, 0, 0)
		org.SetBillingInfo(model.BillingInfo{
			LastBillingTime: lastBilling,
			NextBillingTime: lastBilling,
		})

		// Должно стать 31 мая (так как 30 мая существует)
		expected := newTestDateTime(2025, 5, 31, 10, 0, 0)
		nextBilling := model.CalculateNextBillingTime(lastBilling, tariff)
		assert.Equal(t, expected, nextBilling)
	})

	t.Run("leap year", func(t *testing.T) {
		// 29 февраля високосного года
		lastBilling := newTestDateTime(2024, 2, 29, 10, 0, 0)
		org.SetBillingInfo(model.BillingInfo{
			LastBillingTime: lastBilling,
			NextBillingTime: lastBilling,
		})

		// Должно стать 28 марта
		expected := newTestDateTime(2024, 3, 28, 10, 0, 0)
		nextBilling := model.CalculateNextBillingTime(lastBilling, tariff)
		assert.Equal(t, expected, nextBilling)
	})
}

// Тесты версионирования
func TestVersioning(t *testing.T) {
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

	initialVersion := org.Version()

	t.Run("version increments on billing", func(t *testing.T) {
		billingPeriod, _ := tariff_model.NewTariffPeriod(
			newTestDateTime(2025, 1, 1, 0, 0, 0).Time(),
			newTestDateTime(2025, 1, 1, 10, 0, 0).Time(),
		)

		currentVersion := org.Version()
		org.ProcessBilling(billingPeriod, clock)
		assert.Equal(t, currentVersion+1, org.Version())
	})

	t.Run("version increments on suspend/resume", func(t *testing.T) {
		currentVersion := org.Version()

		// Приостанавливаем
		org.suspendOrganization(clock)
		assert.Equal(t, currentVersion+1, org.Version())

		// Возобновляем
		org.Resume(clock)
		assert.Equal(t, currentVersion+2, org.Version())
	})
}

// Тесты для edge cases
func TestEdgeCases(t *testing.T) {
	orgID := newTestOrganizationID()
	clock := &MockClock{}
	now := newTestDateTime(2025, 1, 1, 10, 0, 0)
	clock.On("Now").Return(now)

	tariff := &MockTariff{}

	t.Run("multiple billing operations", func(t *testing.T) {
		org, _ := model.NewOrganization(
			orgID,
			"Test Org",
			newTestMoney(1000),
			tariff,
			clock,
		)

		billingPeriod, _ := tariff_model.NewTariffPeriod(
			newTestDateTime(2025, 1, 1, 0, 0, 0).Time(),
			newTestDateTime(2025, 1, 1, 10, 0, 0).Time(),
		)

		// Выполняем несколько списаний
		for i := 0; i < 5; i++ {
			prevBalance := org.Balance()
			err := org.ProcessBilling(billingPeriod, clock)
			assert.NoError(t, err)
			assert.Equal(t, prevBalance.Subtract(newTestMoney(100)), org.Balance())
		}

		assert.Equal(t, newTestMoney(500), org.Balance())
	})

	t.Run("deposit after suspension", func(t *testing.T) {
		org, _ := model.NewOrganization(
			orgID,
			"Test Org",
			newTestMoney(100),
			tariff,
			clock,
		)

		billingPeriod, _ := tariff_model.NewTariffPeriod(
			newTestDateTime(2025, 1, 1, 0, 0, 0).Time(),
			newTestDateTime(2025, 1, 1, 10, 0, 0).Time(),
		)

		// Списываем до нуля (приостановка)
		org.ProcessBilling(billingPeriod, clock)
		assert.True(t, org.IsSuspended())

		// Пополняем
		org.Deposit(newTestMoney(500), clock)
		assert.True(t, org.IsSuspended()) // Все еще приостановлена

		// Возобновляем
		org.Resume(clock)
		assert.False(t, org.IsSuspended())
		assert.Equal(t, model.StatusActive, org.Status())
	})

	t.Run("concurrent access simulation", func(t *testing.T) {
		org, _ := model.NewOrganization(
			orgID,
			"Test Org",
			newTestMoney(1000),
			tariff,
			clock,
		)

		// Имитируем конкурентный доступ
		versionBefore := org.Version()

		// Копия агрегата для "другого" пользователя
		orgCopy := *org

		// Первый пользователь вносит изменения
		billingPeriod, _ := tariff_model.NewTariffPeriod(
			newTestDateTime(2025, 1, 1, 0, 0, 0).Time(),
			newTestDateTime(2025, 1, 1, 10, 0, 0).Time(),
		)
		org.ProcessBilling(billingPeriod, clock)
		assert.Equal(t, versionBefore+1, org.Version())

		// Второй пользователь пытается внести изменения с устаревшей версией
		assert.Equal(t, versionBefore, orgCopy.Version())
		assert.Equal(t, versionBefore+1, org.Version())
	})
}
