package valueobjects

import (
	"errors"
	"strings"

	"github.com/GAKiknadze/payment_service/internal/idgen"
)

var (
	ErrInvalidOrganizationID = errors.New("invalid organization ID format")
)

// OrganizationID представляет идентификатор организации
type OrganizationID string

// NewOrganizationID создает новый идентификатор с валидацией
func NewOrganizationID(id string) (OrganizationID, error) {
	if !idgen.ValidatePrefixedID(id, "ORG", 8) {
		return "", ErrInvalidOrganizationID
	}
	return OrganizationID(strings.ToUpper(id)), nil
}

// GenerateOrganizationID генерирует новый идентификатор
func GenerateOrganizationID() OrganizationID {
	return OrganizationID(idgen.GeneratePrefixedID("ORG", 8))
}

// String возвращает строковое представление
func (id OrganizationID) String() string {
	return string(id)
}

// Equals сравнивает два идентификатора
func (id OrganizationID) Equals(other OrganizationID) bool {
	return string(id) == string(other)
}
