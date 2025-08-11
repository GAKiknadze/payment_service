package tariff

import "errors"

var (
	ErrTariffAlreadyArchived       = errors.New("tariff is already archived")
	ErrActiveSubscriptions         = errors.New("cannot modify tariff with active subscriptions")
	ErrInvalidBillingCycle         = errors.New("invalid billing cycle")
	ErrMissingPrices               = errors.New("missing prices for periodic billing cycle")
	ErrCurrencyAlreadyExists       = errors.New("price for this currency already exists")
	ErrArchivedTariff              = errors.New("tariff is archived")
	ErrLastPriceRemoval            = errors.New("cannot remove the last price")
	ErrIncompatibleQuotaDefinition = errors.New("quota definition is incompatible with billing cycle")
)
