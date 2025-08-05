package event

import "github.com/GAKiknadze/payment_service/domain/common/valueobjects"

type TariffPriceUpdated struct {
	TariffID      valueobjects.TariffID
	OldPrice      valueobjects.Money
	NewPrice      valueobjects.Money
	EffectiveDate valueobjects.DateTime
}
