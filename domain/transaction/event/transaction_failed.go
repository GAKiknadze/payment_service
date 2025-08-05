package event

import "github.com/GAKiknadze/payment_service/domain/common/valueobjects"

type TransactionFailed struct {
	TransactionID valueobjects.TransactionID
	FailedAt      valueobjects.DateTime
	Reason        string
}
