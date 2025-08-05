package model

import (
	"errors"

	"github.com/GAKiknadze/payment_service/domain/common/interfaces"
	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
	"github.com/GAKiknadze/payment_service/domain/transaction/event"
	vo "github.com/GAKiknadze/payment_service/domain/transaction/valueobjects"
)

var (
	ErrInvalidTransactionAmount   = errors.New("transaction amount must be positive")
	ErrInvalidTransactionType     = errors.New("invalid transaction type")
	ErrInvalidStatusTransition    = errors.New("invalid status transition")
	ErrIdempotencyKeyMismatch     = errors.New("idempotency key mismatch")
	ErrTransactionAlreadyComplete = errors.New("transaction is already completed")
)

// Transaction - агрегатный корень финансовых операций
type Transaction struct {
	id              valueobjects.TransactionID
	organization    valueobjects.OrganizationID
	amount          valueobjects.Money
	transactionType vo.TransactionType
	status          vo.TransactionStatus
	idempotencyKey  valueobjects.IdempotencyKey
	createdAt       valueobjects.DateTime
	completedAt     *valueobjects.DateTime
	version         int
	events          []interface{} // Буфер доменных событий
}

// NewTransaction создает новую транзакцию с валидацией бизнес-правил
func NewTransaction(
	id valueobjects.TransactionID,
	orgID valueobjects.OrganizationID,
	amount valueobjects.Money,
	txType vo.TransactionType,
	idempotencyKey valueobjects.IdempotencyKey,
	clock interfaces.Clock,
) (*Transaction, error) {
	// Проверка бизнес-инвариантов
	if amount.IsZero() || amount.IsNegative() {
		return nil, ErrInvalidTransactionAmount
	}

	if txType != vo.TransactionTypeDebit && txType != vo.TransactionTypeCredit {
		return nil, ErrInvalidTransactionType
	}

	now := valueobjects.NewDateTime(clock.Now())

	transaction := &Transaction{
		id:              id,
		organization:    orgID,
		amount:          amount,
		transactionType: txType,
		status:          vo.NewTransactionStatus(vo.StatusPending),
		idempotencyKey:  idempotencyKey,
		createdAt:       now,
		version:         1,
	}

	// Генерация доменного события
	transaction.events = append(transaction.events, event.TransactionCreated{
		TransactionID:   id,
		OrganizationID:  orgID,
		Amount:          amount,
		TransactionType: txType,
		Timestamp:       now,
	})

	return transaction, nil
}

// Complete завершает транзакцию с проверкой бизнес-правил
func (t *Transaction) Complete(clock interfaces.Clock) error {
	if !t.status.CanTransitionTo(vo.StatusCompleted) {
		return ErrInvalidStatusTransition
	}

	now := valueobjects.NewDateTime(clock.Now())
	t.completedAt = &now
	t.status = vo.NewTransactionStatus(vo.StatusCompleted)
	t.version++

	// Генерация доменного события
	t.events = append(t.events, event.TransactionCompleted{
		TransactionID: t.id,
		CompletedAt:   now,
	})

	return nil
}

// Fail помечает транзакцию как неудачную
func (t *Transaction) Fail(reason string, clock interfaces.Clock) error {
	if !t.status.CanTransitionTo(vo.StatusFailed) {
		return ErrInvalidStatusTransition
	}

	now := valueobjects.NewDateTime(clock.Now())
	t.status = vo.NewTransactionStatus(vo.StatusFailed)
	t.version++

	// Генерация доменного события
	t.events = append(t.events, event.TransactionFailed{
		TransactionID: t.id,
		FailedAt:      now,
		Reason:        reason,
	})

	return nil
}

// Compensate создает компенсирующую транзакцию
func (t *Transaction) Compensate(
	id valueobjects.TransactionID,
	clock interfaces.Clock,
) (*Transaction, error) {
	if t.status.Value() != vo.StatusCompleted {
		return nil, errors.New("can only compensate completed transactions")
	}

	// Создаем обратную транзакцию
	compensation, err := NewTransaction(
		id,
		t.organization,
		t.amount.Negate(),
		reverseTransactionType(t.transactionType),
		t.idempotencyKey,
		clock,
	)
	if err != nil {
		return nil, err
	}

	// Помечаем исходную транзакцию как компенсированную
	t.status = vo.NewTransactionStatus(vo.StatusCompensated)
	t.version++

	// Генерация доменных событий
	t.events = append(t.events, event.TransactionCompensated{
		OriginalTransactionID: t.id,
		CompensationID:        id,
	})

	compensation.events = append(compensation.events, event.TransactionCreated{
		TransactionID:   id,
		OrganizationID:  t.organization,
		Amount:          t.amount.Negate(),
		TransactionType: reverseTransactionType(t.transactionType),
		Timestamp:       valueobjects.NewDateTime(clock.Now()),
	})

	return compensation, nil
}

// ValidateIdempotency проверяет соответствие идемпотентного ключа
func (t *Transaction) ValidateIdempotency(key valueobjects.IdempotencyKey) error {
	if !t.idempotencyKey.Equals(key) {
		return ErrIdempotencyKeyMismatch
	}
	return nil
}

// CanBeCompleted проверяет, может ли транзакция быть завершена
func (t *Transaction) CanBeCompleted() bool {
	return t.status.CanTransitionTo(vo.StatusCompleted)
}

// IsCompleted проверяет, завершена ли транзакция
func (t *Transaction) IsCompleted() bool {
	return t.status.Value() == vo.StatusCompleted
}

// ID возвращает идентификатор транзакции
func (t *Transaction) ID() valueobjects.TransactionID {
	return t.id
}

// OrganizationID возвращает идентификатор организации
func (t *Transaction) OrganizationID() valueobjects.OrganizationID {
	return t.organization
}

// Amount возвращает сумму транзакции
func (t *Transaction) Amount() valueobjects.Money {
	return t.amount
}

// Status возвращает текущий статус
func (t *Transaction) Status() vo.TransactionStatus {
	return t.status
}

// PopEvents извлекает и сбрасывает буфер доменных событий
func (t *Transaction) PopEvents() []interface{} {
	events := t.events
	t.events = nil
	return events
}

// Вспомогательная функция для определения обратного типа транзакции
func reverseTransactionType(txType vo.TransactionType) vo.TransactionType {
	if txType == vo.TransactionTypeDebit {
		return vo.TransactionTypeCredit
	}
	return vo.TransactionTypeDebit
}
