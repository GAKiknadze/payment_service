package valueobjects

import (
	"errors"
	"strings"

	"github.com/GAKiknadze/payment_service/internal/idgen"
)

var (
	ErrInvalidTariffID = errors.New("invalid tariff ID format")
)

// TariffID представляет идентификатор тарифа
type TariffID string

// NewTariffID создает новый идентификатор с валидацией
func NewTariffID(id string) (TariffID, error) {
	if !idgen.ValidatePrefixedID(id, "TAR", 8) {
		return "", ErrInvalidTariffID
	}
	return TariffID(strings.ToUpper(id)), nil
}

// GenerateTariffID генерирует новый идентификатор
func GenerateTariffID() TariffID {
	return TariffID(idgen.GeneratePrefixedID("TAR", 8))
}

// String возвращает строковое представление
func (id TariffID) String() string {
	return string(id)
}

// Equals сравнивает два идентификатора
func (id TariffID) Equals(other TariffID) bool {
	return string(id) == string(other)
}
