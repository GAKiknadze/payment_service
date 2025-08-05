package event

import "github.com/GAKiknadze/payment_service/domain/common/valueobjects"

type TransactionCompensated struct {
	OriginalTransactionID valueobjects.TransactionID
	CompensationID        valueobjects.TransactionID
}
