package valueobjects

import (
	"errors"
	"regexp"
	"strings"

	"github.com/GAKiknadze/payment_service/internal/idgen"
)

var (
	ErrInvalidIdempotencyKey = errors.New("invalid idempotency key format")
	keyPattern               = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`)
)

// IdempotencyKey представляет ключ для обеспечения идемпотентности
type IdempotencyKey string

// NewIdempotencyKey создает новый ключ с валидацией
func NewIdempotencyKey(key string) (IdempotencyKey, error) {
	if !keyPattern.MatchString(strings.ToLower(key)) {
		return "", ErrInvalidIdempotencyKey
	}
	return IdempotencyKey(strings.ToLower(key)), nil
}

// GenerateIdempotencyKey генерирует новый UUID-ключ
func GenerateIdempotencyKey() IdempotencyKey {
	return IdempotencyKey(idgen.GenerateUUID())
}

// String возвращает строковое представление
func (k IdempotencyKey) String() string {
	return string(k)
}

// Equals сравнивает два ключа
func (k IdempotencyKey) Equals(other IdempotencyKey) bool {
	return strings.EqualFold(string(k), string(other))
}

// Validate проверяет формат ключа
func (k IdempotencyKey) Validate() error {
	if !keyPattern.MatchString(string(k)) {
		return ErrInvalidIdempotencyKey
	}
	return nil
}
