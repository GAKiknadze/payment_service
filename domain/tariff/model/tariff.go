package model

import (
	"errors"
	"time"

	"github.com/GAKiknadze/payment_service/domain/common/interfaces"
	"github.com/GAKiknadze/payment_service/domain/common/utils"
	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
	"github.com/GAKiknadze/payment_service/domain/tariff/event"
)

var (
	ErrInvalidTariffName    = errors.New("tariff name cannot be empty")
	ErrInvalidPrice         = errors.New("price must be positive")
	ErrPeriodAlreadyExpired = errors.New("period has already expired")
)

// Tariff - агрегатный корень тарифной системы
type Tariff struct {
	id      valueobjects.TariffID
	name    string
	price   valueobjects.Money
	period  valueobjects.TimeRange
	version int
	events  []interface{} // Буфер доменных событий
}

// NewTariff создает новый тариф с валидацией инвариантов
func NewTariff(
	id valueobjects.TariffID,
	name string,
	price valueobjects.Money,
	period valueobjects.TimeRange,
	clock interfaces.Clock,
) (*Tariff, error) {
	if name == "" {
		return nil, ErrInvalidTariffName
	}
	if price.IsNegative() {
		return nil, ErrInvalidPrice
	}
	if period.Duration() <= 0 {
		return nil, errors.New("period must have positive duration")
	}
	if period.End().Before(utils.ToDateTime(clock.Now())) {
		return nil, ErrPeriodAlreadyExpired
	}

	return &Tariff{
		id:     id,
		name:   name,
		price:  price,
		period: period,
	}, nil
}

// UpdatePrice изменяет цену с применением бизнес-правил
func (t *Tariff) UpdatePrice(
	newPrice valueobjects.Money,
	effectiveDate valueobjects.DateTime,
	clock interfaces.Clock,
) error {
	if newPrice.IsNegative() {
		return ErrInvalidPrice
	}
	if !t.period.Contains(effectiveDate) {
		return errors.New("effective date must be within tariff period")
	}
	if effectiveDate.Before(utils.ToDateTime(clock.Now())) {
		return errors.New("cannot set price for past dates")
	}

	// Генерация доменного события
	t.events = append(t.events, event.TariffPriceUpdated{
		TariffID:      t.id,
		OldPrice:      t.price,
		NewPrice:      newPrice,
		EffectiveDate: effectiveDate,
	})

	t.price = newPrice
	t.version++
	return nil
}

// CalculateCost вычисляет стоимость за указанный период
func (t *Tariff) CalculateCost(duration time.Duration) (valueobjects.Money, error) {
	if duration < 0 {
		return valueobjects.Money{}, errors.New("duration cannot be negative")
	}

	// Бизнес-логика расчета (например, пропорционально времени)
	hours := duration.Hours()
	hourlyRate, _ := t.price.Divide(valueobjects.NewDecimal(730)) // 730 часов в месяце
	cost, _ := hourlyRate.Multiply(valueobjects.NewDecimal(hours))
	return cost, nil
}

// PopEvents извлекает и сбрасывает буфер доменных событий
func (t *Tariff) PopEvents() []interface{} {
	events := t.events
	t.events = nil
	return events
}
