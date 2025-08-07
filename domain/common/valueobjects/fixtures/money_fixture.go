package fixtures

import "github.com/GAKiknadze/payment_service/domain/common/valueobjects"

func NewTestMoneyRUB(amount float64) valueobjects.Money {
	return valueobjects.MustMoney(valueobjects.NewDecimal(amount), valueobjects.RUB)
}
