package model

import (
	"errors"

	"github.com/GAKiknadze/payment_service/domain/common/interfaces"
	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
	"github.com/GAKiknadze/payment_service/domain/organization/event"
	tariff "github.com/GAKiknadze/payment_service/domain/tariff/model"
)

var (
	ErrInvalidOrganizationName = errors.New("organization name cannot be empty")
	ErrInsufficientBalance     = errors.New("insufficient balance for operation")
	ErrOrganizationSuspended   = errors.New("organization is suspended due to zero balance")
	ErrBillingInProgress       = errors.New("billing operation is already in progress")
	ErrInvalidBillingPeriod    = errors.New("billing period must be positive")
)

// Organization - агрегатный корень управления организацией и её балансом
type Organization struct {
	id          valueobjects.OrganizationID
	name        string
	balance     valueobjects.Money
	tariff      *tariff.Tariff
	status      OrganizationStatus
	billingInfo BillingInfo
	version     int
	events      []interface{} // Буфер доменных событий
	isSuspended bool
}

// BillingInfo содержит информацию о последнем списании
type BillingInfo struct {
	LastBillingTime valueobjects.DateTime
	NextBillingTime valueobjects.DateTime
}

// OrganizationStatus представляет статус организации
type OrganizationStatus int

const (
	StatusActive OrganizationStatus = iota
	StatusSuspended
	StatusTerminated
)

// NewOrganization создает новую организацию с начальным балансом
func NewOrganization(
	id valueobjects.OrganizationID,
	name string,
	initialBalance valueobjects.Money,
	tariff *tariff.Tariff,
	clock interfaces.Clock,
) (*Organization, error) {
	if name == "" {
		return nil, ErrInvalidOrganizationName
	}
	if initialBalance.IsNegative() {
		return nil, errors.New("initial balance cannot be negative")
	}
	if tariff == nil {
		return nil, errors.New("tariff cannot be nil")
	}

	now := valueobjects.NewDateTime(clock.Now())

	return &Organization{
		id:      id,
		name:    name,
		balance: initialBalance,
		tariff:  tariff,
		status:  StatusActive,
		billingInfo: BillingInfo{
			LastBillingTime: now,
			NextBillingTime: calculateNextBillingTime(now, tariff),
		},
		isSuspended: false,
	}, nil
}

// Deposit пополняет баланс организации
func (o *Organization) Deposit(amount valueobjects.Money, clock interfaces.Clock) error {
	if !amount.IsPositive() {
		return errors.New("deposit amount must be positive")
	}
	if o.isSuspended {
		return ErrOrganizationSuspended
	}

	previousBalance := o.balance
	o.balance, _ = o.balance.Add(amount)

	// Генерация доменного события
	o.events = append(o.events, event.BalanceUpdated{
		OrganizationID:  o.id,
		PreviousBalance: previousBalance,
		NewBalance:      o.balance,
		ChangeAmount:    amount,
		ChangeType:      "deposit",
		Timestamp:       valueobjects.NewDateTime(clock.Now()),
	})

	o.version++
	return nil
}

// ProcessBilling выполняет списание по тарифу за указанный период
func (o *Organization) ProcessBilling(
	billingPeriod valueobjects.TimeRange,
	clock interfaces.Clock,
) error {
	if o.isSuspended {
		return ErrOrganizationSuspended
	}
	if billingPeriod.Duration() <= 0 {
		return ErrInvalidBillingPeriod
	}

	// Расчет стоимости
	cost, err := o.tariff.CalculateCost(billingPeriod.Duration())
	if err != nil {
		return err
	}

	// Проверка баланса
	if o.balance.LessThan(cost) {
		return ErrInsufficientBalance
	}

	// Выполняем списание
	previousBalance := o.balance
	o.balance, _ = o.balance.Subtract(cost)

	// Обновляем информацию о списаниях
	now := valueobjects.NewDateTime(clock.Now())
	o.billingInfo = BillingInfo{
		LastBillingTime: now,
		NextBillingTime: calculateNextBillingTime(now, o.tariff),
	}

	// Генерация событий
	o.events = append(o.events, event.BalanceUpdated{
		OrganizationID:  o.id,
		PreviousBalance: previousBalance,
		NewBalance:      o.balance,
		ChangeAmount:    cost,
		ChangeType:      "billing",
		Timestamp:       now,
	})

	o.events = append(o.events, event.BillingProcessed{
		OrganizationID: o.id,
		Amount:         cost,
		PeriodStart:    billingPeriod.Start(),
		PeriodEnd:      billingPeriod.End(),
		Timestamp:      now,
	})

	// Проверка на приостановку
	if o.balance.IsZero() {
		o.suspendOrganization(clock)
	}

	o.version++
	return nil
}

// CheckBalance проверяет возможность выполнения операции
func (o *Organization) CheckBalance(amount valueobjects.Money) bool {
	return !o.balance.LessThan(amount) && !o.isSuspended
}

// Resume возобновляет работу организации после пополнения баланса
func (o *Organization) Resume(clock interfaces.Clock) error {
	if !o.isSuspended {
		return nil // Уже активна
	}
	if o.balance.IsZero() {
		return errors.New("cannot resume with zero balance")
	}

	o.isSuspended = false
	o.status = StatusActive

	// Генерация события
	o.events = append(o.events, event.OrganizationResumed{
		OrganizationID: o.id,
		Timestamp:      valueobjects.NewDateTime(clock.Now()),
	})

	o.version++
	return nil
}

// Terminate окончательно закрывает организацию
func (o *Organization) Terminate(clock interfaces.Clock) error {
	if o.status == StatusTerminated {
		return errors.New("organization is already terminated")
	}

	o.status = StatusTerminated
	o.isSuspended = true

	// Генерация события
	o.events = append(o.events, event.OrganizationTerminated{
		OrganizationID: o.id,
		FinalBalance:   o.balance,
		Timestamp:      valueobjects.NewDateTime(clock.Now()),
	})

	o.version++
	return nil
}

// PopEvents извлекает и сбрасывает буфер доменных событий
func (o *Organization) PopEvents() []interface{} {
	events := o.events
	o.events = nil
	return events
}

// Внутренние методы
func (o *Organization) suspendOrganization(clock interfaces.Clock) {
	o.isSuspended = true
	o.status = StatusSuspended

	o.events = append(o.events, event.OrganizationSuspended{
		OrganizationID: o.id,
		Timestamp:      valueobjects.NewDateTime(clock.Now()),
	})
}

func calculateNextBillingTime(
	lastBillingTime valueobjects.DateTime,
	tariff *tariff.Tariff,
) valueobjects.DateTime {
	return valueobjects.NewDateTime(
		lastBillingTime.Time().AddDate(0, 1, 0),
	)
}
