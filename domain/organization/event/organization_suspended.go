package event

import "github.com/GAKiknadze/payment_service/domain/common/valueobjects"

type OrganizationSuspended struct {
	OrganizationID valueobjects.OrganizationID
	Timestamp      valueobjects.DateTime
}
