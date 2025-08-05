package event

import "github.com/GAKiknadze/payment_service/domain/common/valueobjects"

type OrganizationTerminated struct {
	OrganizationID valueobjects.OrganizationID
	FinalBalance   valueobjects.Money
	Timestamp      valueobjects.DateTime
}
