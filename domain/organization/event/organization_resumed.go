package event

import "github.com/GAKiknadze/payment_service/domain/common/valueobjects"

type OrganizationResumed struct {
	OrganizationID valueobjects.OrganizationID
	Timestamp      valueobjects.DateTime
}
