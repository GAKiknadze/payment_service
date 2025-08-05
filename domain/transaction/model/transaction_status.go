package model

// Состояния транзакции
const (
	StatusPending     = "PENDING"
	StatusCompleted   = "COMPLETED"
	StatusFailed      = "FAILED"
	StatusCompensated = "COMPENSATED"
)

// TransactionStatus - Value Object для управления статусами
type TransactionStatus struct {
	status string
}

// NewTransactionStatus создает валидный статус
func NewTransactionStatus(status string) TransactionStatus {
	// По умолчанию возвращаем PENDING для невалидных значений
	// (в реальности здесь должна быть строгая проверка)
	validStatuses := map[string]bool{
		StatusPending:     true,
		StatusCompleted:   true,
		StatusFailed:      true,
		StatusCompensated: true,
	}

	if !validStatuses[status] {
		status = StatusPending
	}

	return TransactionStatus{status: status}
}

// Value возвращает строковое представление статуса
func (s TransactionStatus) Value() string {
	return s.status
}

// CanTransitionTo проверяет допустимость перехода между статусами
func (s TransactionStatus) CanTransitionTo(newStatus string) bool {
	transitions := map[string]map[string]bool{
		StatusPending: {
			StatusCompleted:   true,
			StatusFailed:      true,
			StatusCompensated: false,
		},
		StatusCompleted: {
			StatusPending:     false,
			StatusFailed:      false,
			StatusCompensated: true,
		},
		StatusFailed: {
			StatusPending:     false,
			StatusCompleted:   false,
			StatusCompensated: false,
		},
		StatusCompensated: {
			StatusPending:   false,
			StatusCompleted: false,
			StatusFailed:    false,
		},
	}

	return transitions[s.status][newStatus]
}

// Equals сравнивает два статуса по значению
func (s TransactionStatus) Equals(other TransactionStatus) bool {
	return s.status == other.status
}
