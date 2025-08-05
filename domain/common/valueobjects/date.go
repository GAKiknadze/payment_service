package valueobjects

import (
	"time"

	"github.com/GAKiknadze/payment_service/domain/common/interfaces"
)

// Date представляет календарную дату без времени
type Date struct {
	year  int
	month time.Month
	day   int
}

// NewDate создает объект даты
func NewDate(year int, month time.Month, day int) (Date, error) {
	// Проверка валидности даты
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	if t.Year() != year || t.Month() != month || t.Day() != day {
		return Date{}, ErrInvalidDate
	}
	return Date{year, month, day}, nil
}

// Today возвращает текущую дату
func Today(clock interfaces.Clock) Date {
	year, month, day := clock.Today()
	return Date{
		year:  year,
		month: month,
		day:   day,
	}
}

// Equals сравнивает даты
func (d Date) Equals(other Date) bool {
	return d.year == other.year && d.month == other.month && d.day == other.day
}
