package tariff

import (
	"time"

	common "github.com/GAKiknadze/payment_service/domain/common/valueobject"
)

type Tariff struct {
	id           common.TariffID
	name         string
	description  *string
	status       string
	billingCycle common.BillingCycle
	isExtendable bool
	createdAt    time.Time
	updatedAt    time.Time
	archivedAt   time.Time
	prices       []common.Price
	version      uint
	events       []interface{}
}

// PopEvents извлекает и сбрасывает буфер доменных событий
func (o *Tariff) PopEvents() []interface{} {
	events := o.events
	o.events = nil
	return events
}
