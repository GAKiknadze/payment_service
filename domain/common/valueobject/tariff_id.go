package valueobject

import (
	"errors"

	"github.com/GAKiknadze/payment_service/internal/idgen/generic"
)

type tariffConfig struct{}

func (tariffConfig) Config() generic.IdConfig {
	return generic.IdConfig{
		Prefix: "TAR",
		Err:    ErrInvalidTariffID,
	}
}

var ErrInvalidTariffID = errors.New("invalid tariff ID format")

type TariffID = generic.ID[tariffConfig]

func NewTariffID(id string) (TariffID, error) {
	return generic.NewID[tariffConfig](id)
}

func GenerateTariffID() TariffID {
	return generic.GenerateID[tariffConfig]()
}
