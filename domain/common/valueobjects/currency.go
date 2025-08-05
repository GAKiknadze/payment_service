package valueobjects

import "strings"

// Currency представляет валюту (Value Object)
type Currency string

// Константы валют
const (
	RUB Currency = "RUB"
)

// IsValid проверяет, является ли валюта поддерживаемой
func (c Currency) IsValid() bool {
	switch strings.ToUpper(string(c)) {
	case "RUB":
		return true
	default:
		return false
	}
}

// String возвращает строковое представление валюты
func (c Currency) String() string {
	return string(c)
}
