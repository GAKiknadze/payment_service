package model

import (
	"errors"
	"time"

	"github.com/GAKiknadze/payment_service/domain/common/interfaces"
	"github.com/GAKiknadze/payment_service/domain/common/utils"
	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
	"github.com/GAKiknadze/payment_service/domain/tariff/event"
	"github.com/shopspring/decimal"
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

	oldPrice := t.price

	t.price = newPrice

	// Генерация доменного события
	t.events = append(t.events, event.TariffPriceUpdated{
		TariffID:      t.id,
		OldPrice:      oldPrice,
		NewPrice:      newPrice,
		EffectiveDate: effectiveDate,
	})

	t.version++
	return nil
}

// CalculateCost вычисляет стоимость за указанный период
func (t *Tariff) CalculateCost(duration time.Duration) (valueobjects.Money, error) {
	// if duration < 0 {
	// 	return valueobjects.Money{}, valueobjects.ErrInvalidBillingPeriod
	// }

	// hours := duration.Hours()

	// if math.Abs(hours-730) < 0.001 {
	// 	return t.price, nil
	// }

	// numerator := t.price.Amount().Mul(decimal.NewFromFloat(hours))
	// result := numerator.Div(decimal.NewFromInt(730))
	// return valueobjects.MustMoney(result, t.Price().Currency()), nil

	if duration < 0 {
		return valueobjects.Money{}, valueobjects.ErrInvalidBillingPeriod
	}

	totalHours := t.period.Duration().Hours()

	if totalHours <= 0 {
		return valueobjects.Money{}, errors.New("tariff period must have positive duration")
	}

	ratio := duration.Hours() / totalHours

	// Стоимость = цена тарифа * доля периода
	cost, _ := t.price.Multiply(decimal.NewFromFloat(ratio))
	return cost, nil
}

// PopEvents извлекает и сбрасывает буфер доменных событий
func (t *Tariff) PopEvents() []interface{} {
	events := t.events
	t.events = nil
	return events
}

// ID возвращает идентификатор тарифа
func (t *Tariff) ID() valueobjects.TariffID {
	return t.id
}

// Name возвращает название тарифа
func (t *Tariff) Name() string {
	return t.name
}

// Price возвращает цена тарифа
func (t *Tariff) Price() valueobjects.Money {
	return t.price
}

// Period возвращает период действия тарифа
func (t *Tariff) Period() valueobjects.TimeRange {
	return t.period
}

// Version возвращает номер версии тарифа
func (t *Tariff) Version() int {
	return t.version
}
