package valueobject

import "errors"

var ErrInvalidPrice = errors.New("invalid price")

type Price struct {
	id        string
	amount    MoneyAmount
	isDefault bool
}

// NewPrice - фабричный метод для создания объекта Price
func NewPrice(id string, amount MoneyAmount, isDefault bool) (Price, error) {
	// Проверка, что сумма валидна
	if !amount.IsValid() {
		return Price{}, ErrInvalidPrice
	}

	return Price{
		id:        id,
		amount:    amount,
		isDefault: isDefault,
	}, nil
}

// ID Уникальный идентификатор цены
func (p Price) ID() string {
	return p.id
}

// Amount возвращает сумму цены как MoneyAmount
func (p Price) Amount() MoneyAmount {
	return p.amount
}

// Currency возвращает валюту цены
func (p Price) Currency() Currency {
	return p.amount.Currency()
}

// IsDefault является ли валютой по умолчанию
func (p Price) IsDefault() bool {
	return p.isDefault
}

// Format возвращает отформатированное строковое представление цены
func (p Price) Format() string {
	return p.amount.Format()
}

// IsCompatibleWith проверяет совместимость валюты цены с валютой организации
func (p Price) IsCompatibleWith(organizationCurrency Currency) bool {
	return p.Currency().Code() == organizationCurrency.Code()
}
