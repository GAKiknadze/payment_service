package tariff

import (
	"time"

	common "github.com/GAKiknadze/payment_service/domain/common/valueobject"
)

type TariffFilter struct{}

type ITariffRepository interface {
	Create(tariff *Tariff) (common.TariffID, error)
	GetByID(tariffID common.TariffID) (*Tariff, error)
	GetActiveTariffs() ([]Tariff, error)
	GetTariffs(filter TariffFilter, page int, pageSize int) ([]Tariff, int, error)
	Update(tariff *Tariff) error
	Archive(tariffID common.TariffID, archivedAt time.Time) error
	HasActiveSubscriptions(tariffID common.TariffID) (bool, error)
	GetPriceByCurrency(tariffID common.TariffID, currency string) (*common.Price, error)
}
