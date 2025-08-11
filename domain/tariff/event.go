package tariff

import (
	"time"

	common "github.com/GAKiknadze/payment_service/domain/common/valueobject"
)

type EventTariffCreated struct {
	TariffID     common.TariffID
	Name         string
	Description  *string
	BillingCycle string
	Prices       []common.Price
	Quotas       []common.QuotaDefinition
	CreatedAt    time.Time
}

type EventTariffUpdated struct {
	TariffID             common.TariffID
	ChangedFields        []string
	UpdatedAt            time.Time
	RequiresNotification bool
	NewVersion           uint
}

type EventTariffArchived struct {
	TariffID                 common.TariffID
	ArchivedAt               time.Time
	Reason                   *string
	DeprecationDate          time.Time
	ActiveSubscriptionsCount uint
	NewVersion               uint
}

type EventPriceAdded struct {
	TariffID   common.TariffID
	Currency   string
	Amount     string
	IsDefault  bool
	AddedAt    time.Time
	NewVersion uint
}

type EventPriceRemoved struct {
	TariffID           common.TariffID
	Currency           string
	Price              common.Price
	WasDefault         bool
	NewDefaultCurrency string
	RemovedAt          time.Time
	NewVersion         uint
}

type EventQuotasUpdated struct {
	TariffID   common.TariffID
	OldQuotas  []common.QuotaDefinition
	NewQuotas  []common.QuotaDefinition
	UpdatedAt  time.Time
	NewVersion uint
}
