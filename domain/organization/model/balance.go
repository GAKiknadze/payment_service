package model

import (
	"errors"

	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
)

// Balance представляет баланс организации
// Value Object: неизменяемый, сравнивается по значению
type Balance struct {
	amount valueobjects.Money
}

// NewBalance создает валидный баланс
func NewBalance(initialAmount valueobjects.Money) (Balance, error) {
	if initialAmount.IsNegative() {
		return Balance{}, errors.New("balance cannot be negative")
	}
	return Balance{amount: initialAmount}, nil
}

// Add увеличивает баланс
func (b Balance) Add(amount valueobjects.Money) Balance {
	if amount.IsNegative() {
		panic("cannot add negative amount to balance")
	}
	newAmount, _ := b.amount.Add(amount)
	return Balance{amount: newAmount}
}

// Subtract уменьшает баланс
func (b Balance) Subtract(amount valueobjects.Money) (Balance, error) {
	if amount.IsNegative() {
		return Balance{}, errors.New("cannot subtract negative amount")
	}
	if b.amount.LessThan(amount) {
		return Balance{}, errors.New("insufficient balance")
	}
	newAmount, _ := b.amount.Subtract(amount)
	return Balance{amount: newAmount}, nil
}

// IsZero проверяет, является ли баланс нулевым
func (b Balance) IsZero() bool {
	return b.amount.IsZero()
}
