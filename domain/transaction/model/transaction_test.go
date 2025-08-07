// domain/transaction/model/transaction_test.go
package model_test

import (
	"strings"
	"testing"
	"time"

	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
	"github.com/GAKiknadze/payment_service/domain/common/valueobjects/fixtures"
	mocks "github.com/GAKiknadze/payment_service/domain/common/valueobjects/mocks"
	"github.com/GAKiknadze/payment_service/domain/transaction/event"
	"github.com/GAKiknadze/payment_service/domain/transaction/model"
	vo "github.com/GAKiknadze/payment_service/domain/transaction/valueobjects"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
)

// TestTransaction_NewTransaction проверяет создание транзакции
func TestTransaction_NewTransaction(t *testing.T) {
	// Настройка тестовых данных
	testTime := time.Date(2025, 10, 1, 12, 0, 0, 0, time.UTC)
	testOrgID := valueobjects.GenerateOrganizationID()
	validAmount := fixtures.NewTestMoneyRUB(100.00) // 100.00

	// Настройка мока времени
	mockClock := mocks.NewFixedClock(testTime)

	t.Run("valid transaction creation", func(t *testing.T) {
		idempotencyKey, _ := valueobjects.NewIdempotencyKey("valid-key-123")
		txID := valueobjects.GenerateTransactionID()

		tx, err := model.NewTransaction(
			txID,
			testOrgID,
			validAmount,
			vo.TransactionTypeDebit,
			idempotencyKey,
			mockClock,
		)

		assert.NoError(t, err)
		assert.NotNil(t, tx)
		assert.Equal(t, txID, tx.ID())
		assert.Equal(t, testOrgID, tx.OrganizationID())
		assert.Equal(t, validAmount, tx.Amount())
		assert.Equal(t, vo.TransactionTypeDebit, tx.TransactionType())
		assert.Equal(t, vo.StatusPending, tx.Status().Value())
		assert.Equal(t, valueobjects.NewDateTime(testTime), tx.CreatedAt())
		assert.Nil(t, tx.CompletedAt())
		assert.Equal(t, 1, tx.Version())

		// Проверка сгенерированного события
		events := tx.PopEvents()
		assert.Len(t, events, 1)
		createdEvent, ok := events[0].(event.TransactionCreated)
		assert.True(t, ok)
		assert.Equal(t, txID, createdEvent.TransactionID)
		assert.Equal(t, testOrgID, createdEvent.OrganizationID)
		assert.Equal(t, validAmount, createdEvent.Amount)
		assert.Equal(t, vo.TransactionTypeDebit, createdEvent.TransactionType)
		assert.Equal(t, valueobjects.NewDateTime(testTime), createdEvent.Timestamp)
	})

	t.Run("negative amount returns error", func(t *testing.T) {
		idempotencyKey, _ := valueobjects.NewIdempotencyKey("key-123")
		money, _ := valueobjects.NewMoney(valueobjects.NewDecimal(-100.00), valueobjects.RUB)
		tx, err := model.NewTransaction(
			valueobjects.GenerateTransactionID(),
			testOrgID,
			money,
			vo.TransactionTypeDebit,
			idempotencyKey,
			mockClock,
		)

		assert.ErrorIs(t, err, model.ErrInvalidTransactionAmount)
		assert.Nil(t, tx)
	})

	t.Run("zero amount returns error", func(t *testing.T) {
		idempotencyKey, _ := valueobjects.NewIdempotencyKey("key-123")
		tx, err := model.NewTransaction(
			valueobjects.GenerateTransactionID(),
			testOrgID,
			fixtures.NewTestMoneyRUB(0),
			vo.TransactionTypeDebit,
			idempotencyKey,
			mockClock,
		)

		assert.ErrorIs(t, err, model.ErrInvalidTransactionAmount)
		assert.Nil(t, tx)
	})

	t.Run("invalid transaction type returns error", func(t *testing.T) {
		idempotencyKey, _ := valueobjects.NewIdempotencyKey("key-123")
		tx, err := model.NewTransaction(
			valueobjects.GenerateTransactionID(),
			testOrgID,
			validAmount,
			"INVALID_TYPE",
			idempotencyKey,
			mockClock,
		)

		assert.ErrorIs(t, err, model.ErrInvalidTransactionType)
		assert.Nil(t, tx)
	})
}

// TestTransaction_StatusTransitions проверяет переходы между статусами
func TestTransaction_StatusTransitions(t *testing.T) {
	testTime := time.Date(2025, 10, 1, 12, 0, 0, 0, time.UTC)
	mockClock := mocks.NewFixedClock(testTime)
	validAmount := fixtures.NewTestMoneyRUB(10000)
	testOrgID := valueobjects.GenerateOrganizationID()
	idempotencyKey, _ := valueobjects.NewIdempotencyKey("key-123")

	tx, _ := model.NewTransaction(
		valueobjects.GenerateTransactionID(),
		testOrgID,
		validAmount,
		vo.TransactionTypeDebit,
		idempotencyKey,
		mockClock,
	)

	t.Run("complete pending transaction", func(t *testing.T) {
		err := tx.Complete(mockClock)
		assert.NoError(t, err)
		assert.Equal(t, vo.StatusCompleted, tx.Status().Value())

		// Проверка события
		events := tx.PopEvents()
		assert.Len(t, events, 2)
		completedEvent, ok := events[1].(event.TransactionCompleted)
		assert.True(t, ok)
		assert.Equal(t, valueobjects.NewDateTime(testTime), completedEvent.CompletedAt)
	})

	t.Run("complete already completed transaction", func(t *testing.T) {
		// Уже завершена в предыдущем тесте
		err := tx.Complete(mockClock)
		assert.ErrorIs(t, err, model.ErrInvalidStatusTransition)
	})

	t.Run("fail pending transaction", func(t *testing.T) {
		// Создаем новую транзакцию
		tx, _ := model.NewTransaction(
			valueobjects.GenerateTransactionID(),
			testOrgID,
			validAmount,
			vo.TransactionTypeDebit,
			idempotencyKey,
			mockClock,
		)

		// Очистка списка событий
		_ = tx.PopEvents()

		err := tx.Fail("test reason", mockClock)
		assert.NoError(t, err)
		assert.Equal(t, vo.StatusFailed, tx.Status().Value())

		// Проверка события
		events := tx.PopEvents()
		assert.Len(t, events, 1)
		failedEvent, ok := events[0].(event.TransactionFailed)
		assert.True(t, ok)
		assert.Equal(t, "test reason", failedEvent.Reason)
	})

	t.Run("invalid status transitions", func(t *testing.T) {
		// Создаем завершенную транзакцию
		tx, _ := model.NewTransaction(
			valueobjects.GenerateTransactionID(),
			testOrgID,
			validAmount,
			vo.TransactionTypeDebit,
			idempotencyKey,
			mockClock,
		)
		tx.Complete(mockClock)

		// // Попытка перевести завершенную в статус pending
		// tx.status = vo.NewTransactionStatus(vo.StatusPending)
		assert.False(t, tx.CanBeCompleted())
		assert.False(t, tx.Status().CanTransitionTo(vo.StatusPending))
	})
}

// TestTransaction_Compensation проверяет механизм компенсации
func TestTransaction_Compensation(t *testing.T) {
	testTime := time.Date(2025, 10, 1, 12, 0, 0, 0, time.UTC)
	mockClock := mocks.NewFixedClock(testTime)

	validAmount := fixtures.NewTestMoneyRUB(10000)
	testOrgID := valueobjects.GenerateOrganizationID()
	idempotencyKey, _ := valueobjects.NewIdempotencyKey("key-123")

	t.Run("compensate completed transaction", func(t *testing.T) {
		tx, _ := model.NewTransaction(
			valueobjects.GenerateTransactionID(),
			testOrgID,
			validAmount,
			vo.TransactionTypeDebit,
			idempotencyKey,
			mockClock,
		)
		tx.Complete(mockClock)

		compensationID := valueobjects.GenerateTransactionID()
		compensation, err := tx.Compensate(compensationID, mockClock)

		assert.NoError(t, err)
		assert.NotNil(t, compensation)
		assert.Equal(t, vo.StatusCompensated, tx.Status().Value())
		assert.Equal(t, vo.TransactionTypeCredit, compensation.TransactionType())
		assert.Equal(t, validAmount, compensation.Amount())
		assert.False(t, compensation.IsCompleted())
		assert.True(t, compensation.CanBeCompleted())

		// Проверка событий
		txEvents := tx.PopEvents()
		compEvents := compensation.PopEvents()

		assert.Len(t, txEvents, 3)
		assert.Len(t, compEvents, 2)

		compEvent, ok := txEvents[2].(event.TransactionCompensated)
		assert.True(t, ok)
		assert.Equal(t, tx.ID(), compEvent.OriginalTransactionID)
		assert.Equal(t, compensationID, compEvent.CompensationID)
	})

	t.Run("cannot compensate non-completed transaction", func(t *testing.T) {
		tx, _ := model.NewTransaction(
			valueobjects.GenerateTransactionID(),
			testOrgID,
			validAmount,
			vo.TransactionTypeDebit,
			idempotencyKey,
			mockClock,
		)

		compensationID := valueobjects.GenerateTransactionID()
		compensation, err := tx.Compensate(compensationID, mockClock)

		assert.Error(t, err)
		assert.Nil(t, compensation)
		assert.Equal(t, vo.StatusPending, tx.Status().Value())
	})
}

// TestTransaction_Idempotency проверяет работу с идемпотентными ключами
func TestTransaction_Idempotency(t *testing.T) {
	testTime := time.Date(2025, 10, 1, 12, 0, 0, 0, time.UTC)
	mockClock := mocks.NewFixedClock(testTime)

	validAmount := fixtures.NewTestMoneyRUB(10000)
	testOrgID := valueobjects.GenerateOrganizationID()

	t.Run("valid idempotency key", func(t *testing.T) {
		key1, _ := valueobjects.NewIdempotencyKey("key-123")
		key2, _ := valueobjects.NewIdempotencyKey("key-123")

		tx, _ := model.NewTransaction(
			valueobjects.GenerateTransactionID(),
			testOrgID,
			validAmount,
			vo.TransactionTypeDebit,
			key1,
			mockClock,
		)

		err := tx.ValidateIdempotency(key2)
		assert.NoError(t, err)
	})

	t.Run("mismatched idempotency key", func(t *testing.T) {
		key1 := valueobjects.GenerateIdempotencyKey()
		key2 := valueobjects.GenerateIdempotencyKey()

		tx, _ := model.NewTransaction(
			valueobjects.GenerateTransactionID(),
			testOrgID,
			validAmount,
			vo.TransactionTypeDebit,
			key1,
			mockClock,
		)

		err := tx.ValidateIdempotency(key2)
		assert.ErrorIs(t, err, model.ErrIdempotencyKeyMismatch)
	})

	t.Run("invalid idempotency key format", func(t *testing.T) {
		_, err := valueobjects.NewIdempotencyKey("")
		assert.Error(t, err)

		_, err = valueobjects.NewIdempotencyKey(strings.Repeat("a", 65))
		assert.Error(t, err)

		_, err = valueobjects.NewIdempotencyKey("invalid@key")
		assert.Error(t, err)
	})
}

// TestTransaction_DomainEvents проверяет корректность генерации событий
func TestTransaction_DomainEvents(t *testing.T) {
	testTime := time.Date(2025, 10, 1, 12, 0, 0, 0, time.UTC)
	mockClock := mocks.NewFixedClock(testTime)

	validAmount := fixtures.NewTestMoneyRUB(10000)
	testOrgID := valueobjects.GenerateOrganizationID()
	idempotencyKey, _ := valueobjects.NewIdempotencyKey("key-123")

	t.Run("events are properly buffered and cleared", func(t *testing.T) {
		tx, _ := model.NewTransaction(
			valueobjects.GenerateTransactionID(),
			testOrgID,
			validAmount,
			vo.TransactionTypeDebit,
			idempotencyKey,
			mockClock,
		)

		// После создания должно быть 1 событие
		assert.Len(t, tx.PopEvents(), 1)
		assert.Len(t, tx.PopEvents(), 0) // Буфер должен быть очищен

		// После завершения должно быть 1 событие
		tx.Complete(mockClock)
		assert.Len(t, tx.PopEvents(), 1)
		assert.Len(t, tx.PopEvents(), 0)

		// После компенсации должно быть 1 событие
		compensationID := valueobjects.GenerateTransactionID()
		compensateTx, err := tx.Compensate(compensationID, mockClock)
		assert.NoError(t, err)
		assert.Len(t, tx.PopEvents(), 1)
		assert.Len(t, compensateTx.PopEvents(), 2)
	})

	t.Run("events contain correct data", func(t *testing.T) {
		txID := valueobjects.GenerateTransactionID()
		tx, _ := model.NewTransaction(
			txID,
			testOrgID,
			validAmount,
			vo.TransactionTypeDebit,
			idempotencyKey,
			mockClock,
		)

		events := tx.PopEvents()
		createdEvent := events[0].(event.TransactionCreated)

		assert.Equal(t, txID, createdEvent.TransactionID)
		assert.Equal(t, testOrgID, createdEvent.OrganizationID)
		assert.Equal(t, validAmount, createdEvent.Amount)
		assert.Equal(t, vo.TransactionTypeDebit, createdEvent.TransactionType)
		assert.Equal(t, valueobjects.NewDateTime(testTime), createdEvent.Timestamp)
	})
}

// TestTransactionStatus_Transitions проверяет логику переходов статусов
func TestTransactionStatus_Transitions(t *testing.T) {
	t.Run("valid transitions", func(t *testing.T) {
		pending := vo.NewTransactionStatus(vo.StatusPending)
		assert.True(t, pending.CanTransitionTo(vo.StatusCompleted))
		assert.True(t, pending.CanTransitionTo(vo.StatusFailed))
		assert.False(t, pending.CanTransitionTo(vo.StatusCompensated))

		completed := vo.NewTransactionStatus(vo.StatusCompleted)
		assert.False(t, completed.CanTransitionTo(vo.StatusPending))
		assert.False(t, completed.CanTransitionTo(vo.StatusFailed))
		assert.True(t, completed.CanTransitionTo(vo.StatusCompensated))

		failed := vo.NewTransactionStatus(vo.StatusFailed)
		assert.False(t, failed.CanTransitionTo(vo.StatusPending))
		assert.False(t, failed.CanTransitionTo(vo.StatusCompleted))
		assert.False(t, failed.CanTransitionTo(vo.StatusCompensated))

		compensated := vo.NewTransactionStatus(vo.StatusCompensated)
		assert.False(t, compensated.CanTransitionTo(vo.StatusPending))
		assert.False(t, compensated.CanTransitionTo(vo.StatusCompleted))
		assert.False(t, compensated.CanTransitionTo(vo.StatusFailed))
	})

	t.Run("invalid status handling", func(t *testing.T) {
		invalid := vo.NewTransactionStatus("INVALID_STATUS")
		assert.Equal(t, vo.StatusPending, invalid.Value())
	})
}
