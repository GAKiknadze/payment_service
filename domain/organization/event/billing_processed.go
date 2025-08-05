package event

import "github.com/GAKiknadze/payment_service/domain/common/valueobjects"

type BillingProcessed struct {
	OrganizationID valueobjects.OrganizationID
	Amount         valueobjects.Money
	PeriodStart    valueobjects.DateTime
	PeriodEnd      valueobjects.DateTime
	Timestamp      valueobjects.DateTime
}
