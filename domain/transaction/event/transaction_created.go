package event

import (
	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
	vo "github.com/GAKiknadze/payment_service/domain/transaction/valueobjects"
)

type TransactionCreated struct {
	TransactionID   valueobjects.TransactionID
	OrganizationID  valueobjects.OrganizationID
	Amount          valueobjects.Money
	TransactionType vo.TransactionType
	Timestamp       valueobjects.DateTime
}
