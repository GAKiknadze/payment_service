package event

import "github.com/GAKiknadze/payment_service/domain/common/valueobjects"

type BalanceUpdated struct {
	OrganizationID  valueobjects.OrganizationID
	PreviousBalance valueobjects.Money
	NewBalance      valueobjects.Money
	ChangeAmount    valueobjects.Money
	ChangeType      string // "deposit", "billing", etc.
	Timestamp       valueobjects.DateTime
}
