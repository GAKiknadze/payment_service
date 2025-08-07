package valueobjects

import "errors"

var (
	// ErrInvalidMoney ошибка невалидной денежной суммы
	ErrInvalidMoney = errors.New("invalid money value")
	// ErrNegativeAmount ошибка отрицательной суммы
	ErrNegativeAmount = errors.New("negative amount not allowed")
	// ErrDivisionByZero ошибка деления на ноль
	ErrDivisionByZero = errors.New("division by zero")
	// ErrInvalidDate
	ErrInvalidDate = errors.New("ivalid date")
	// ErrInvalidBillingPeriod ошибка отрицательного периода
	ErrInvalidBillingPeriod = errors.New("billing period must be positive")
)
