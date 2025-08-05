package valueobjects

import (
	"errors"
	"strings"

	"github.com/GAKiknadze/payment_service/internal/idgen"
)

var (
	ErrInvalidTransactionID = errors.New("invalid transaction ID format")
)

// TransactionID представляет идентификатор организации
type TransactionID string

// NewTransactionID создает новый идентификатор с валидацией
func NewTransactionID(id string) (TransactionID, error) {
	if !idgen.ValidatePrefixedID(id, "TR", 8) {
		return "", ErrInvalidOrganizationID
	}
	return TransactionID(strings.ToUpper(id)), nil
}

// GenerateTransactionID генерирует новый идентификатор
func GenerateTransactionID() TransactionID {
	return TransactionID(idgen.GeneratePrefixedID("TR", 8))
}

// String возвращает строковое представление
func (id TransactionID) String() string {
	return string(id)
}

// Equals сравнивает два идентификатора
func (id TransactionID) Equals(other TransactionID) bool {
	return string(id) == string(other)
}
