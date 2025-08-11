package tariff

import (
	"time"

	common "github.com/GAKiknadze/payment_service/domain/common/valueobject"
)

type EventTariffCreated struct {
	tariffID     common.TariffID
	name         string
	description  *string
	billingCycle string
	prices       []common.Price
	quotas       []common.QuotaDefinition
	createdAt    time.Time
}

type EventTariffUpdated struct {
	tariffID             common.TariffID
	oldVersion           uint
	newVersion           uint
	changedFields        []interface{}
	updatedAt            time.Time
	requiresNotification bool
}

type EventTariffArchived struct {
	tariffID                 common.TariffID
	archivedAt               time.Time
	reason                   *string
	deprecationDate          time.Time
	activeSubscriptionsCount uint
}

type EventPriceAdded struct {
	tariffID  common.TariffID
	currency  string
	amount    string
	isDefault bool
	addedAt   time.Time
}

type EventPriceRemoved struct {
	tariffID           common.TariffID
	currency           string
	wasDefault         bool
	newDefaultCurrency string
	removedAt          time.Time
}
