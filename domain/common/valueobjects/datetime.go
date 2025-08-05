package valueobjects

import (
	"errors"
	"time"
)

var (
	ErrInvalidDateTime = errors.New("invalid date time")
	ErrInvalidTimeZone = errors.New("invalid time zone")
)

// DateTime представляет момент времени с временной зоной
// Value Object: неизменяемый, сравнивается по значению
type DateTime struct {
	time time.Time
}

// NewDateTime создает новый DateTime с валидацией
func NewDateTime(t time.Time) DateTime {
	// Нормализуем время к UTC для консистентности
	return DateTime{
		time: t.UTC(),
	}
}

// Now создает DateTime для текущего момента
func Now() DateTime {
	return NewDateTime(time.Now())
}

// Time возвращает внутреннее значение времени
func (dt DateTime) Time() time.Time {
	return dt.time
}

// String возвращает ISO 8601 представление
func (dt DateTime) String() string {
	return dt.time.Format(time.RFC3339)
}

// Equals сравнивает два DateTime
func (dt DateTime) Equals(other DateTime) bool {
	return dt.time.Equal(other.time)
}

// Before проверяет, предшествует ли один момент другому
func (dt DateTime) Before(other DateTime) bool {
	return dt.time.Before(other.time)
}

// After проверяет, следует ли один момент за другим
func (dt DateTime) After(other DateTime) bool {
	return dt.time.After(other.time)
}

// Add добавляет продолжительность
func (dt DateTime) Add(duration time.Duration) DateTime {
	return NewDateTime(dt.time.Add(duration))
}

// Sub вычитает продолжительность
func (dt DateTime) Sub(other DateTime) time.Duration {
	return dt.time.Sub(other.time)
}

// Format форматирует время в указанном формате
func (dt DateTime) Format(layout string) string {
	return dt.time.Format(layout)
}

// IsZero проверяет, является ли время нулевым
func (dt DateTime) IsZero() bool {
	return dt.time.IsZero()
}
