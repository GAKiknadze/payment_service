package valueobjects

import (
	"errors"
	"time"
)

var (
	ErrInvalidTimeRange = errors.New("invalid time range: start must be before end")
)

// TimeRange представляет временной интервал
type TimeRange struct {
	start DateTime
	end   DateTime
}

// NewTimeRange создает новый временной интервал с валидацией
func NewTimeRange(start, end time.Time) (TimeRange, error) {
	startDT := NewDateTime(start)
	endDT := NewDateTime(end)

	if !endDT.After(startDT) {
		return TimeRange{}, ErrInvalidTimeRange
	}

	return TimeRange{
		start: startDT,
		end:   endDT,
	}, nil
}

// Start возвращает начало интервала
func (tr TimeRange) Start() DateTime {
	return tr.start
}

// End возвращает конец интервала
func (tr TimeRange) End() DateTime {
	return tr.end
}

// Duration возвращает продолжительность интервала
func (tr TimeRange) Duration() time.Duration {
	return tr.end.Time().Sub(tr.start.Time())
}

// Contains проверяет, содержит ли интервал указанный момент
func (tr TimeRange) Contains(moment DateTime) bool {
	return moment.Equals(tr.start) || moment.Equals(tr.end) || moment.After(tr.start) && moment.Before(tr.end)
}

// Equals сравнивает два временных интервала
func (tr TimeRange) Equals(other TimeRange) bool {
	return tr.start.Equals(other.start) && tr.end.Equals(other.end)
}
