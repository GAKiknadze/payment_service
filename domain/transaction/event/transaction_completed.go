package event

import "github.com/GAKiknadze/payment_service/domain/common/valueobjects"

type TransactionCompleted struct {
	TransactionID valueobjects.TransactionID
	CompletedAt   valueobjects.DateTime
}
