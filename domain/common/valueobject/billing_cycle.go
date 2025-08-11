package valueobject

import (
	"errors"
	"time"
)

// BillingCycleType - перечисление типов циклов
type BillingCycleType string

const (
	BillingCycleHourly  BillingCycleType = "Hourly"
	BillingCycleMonthly BillingCycleType = "Monthly"
	BillingCycleOneTime BillingCycleType = "OneTime"
)

var ErrUnsupportedBillingCycleType = errors.New("unsupported billing cycle type")

type BillingCycle struct {
	cycleType   BillingCycleType
	isRecurring bool
	displayName string
}

// NewBillingCycle - фабричный метод для создания BillingCycle
func NewBillingCycle(cycleType BillingCycleType) (BillingCycle, error) {
	switch cycleType {
	case BillingCycleHourly:
		return BillingCycle{
			cycleType:   BillingCycleHourly,
			isRecurring: true,
			displayName: string(BillingCycleHourly),
		}, nil
	case BillingCycleMonthly:
		return BillingCycle{
			cycleType:   BillingCycleMonthly,
			isRecurring: true,
			displayName: string(BillingCycleMonthly),
		}, nil
	case BillingCycleOneTime:
		return BillingCycle{
			cycleType:   BillingCycleOneTime,
			isRecurring: false,
			displayName: string(BillingCycleOneTime),
		}, nil
	default:
		return BillingCycle{}, ErrUnsupportedBillingCycleType
	}
}

// CalculateNextBillingDate - метод для расчета следующей даты списания
// Учитывает особенности календаря (разное количество дней в месяцах)
func (bc BillingCycle) CalculateNextBillingDate(currentDate time.Time) (time.Time, error) {
	if !bc.isRecurring {
		return time.Time{}, nil // Для OneTime нет следующего списания
	}

	switch bc.cycleType {
	case BillingCycleHourly:
		return currentDate.Add(time.Hour), nil

	case BillingCycleMonthly:
		// Для месячного цикла прибавляем 1 месяц с корректной обработкой дней
		year, month, _ := currentDate.Date()
		nextMonth := month + 1
		nextYear := year

		// Обработка перехода на следующий год
		if nextMonth > 12 {
			nextMonth = 1
			nextYear++
		}

		// Определяем последний день следующего месяца
		lastDayOfMonth := time.Date(nextYear, nextMonth+1, 0, 0, 0, 0, 0, currentDate.Location()).Day()

		// Берем день текущей даты, но не больше последнего дня следующего месяца
		day := currentDate.Day()
		if day > lastDayOfMonth {
			day = lastDayOfMonth
		}

		// Сохраняем время (часы, минуты, секунды, наносекунды)
		return time.Date(nextYear, nextMonth, day,
			currentDate.Hour(), currentDate.Minute(), currentDate.Second(),
			currentDate.Nanosecond(), currentDate.Location()), nil

	default:
		return time.Time{}, ErrUnsupportedBillingCycleType
	}
}
