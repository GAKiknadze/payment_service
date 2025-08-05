package utils

import (
	"time"

	"github.com/GAKiknadze/payment_service/domain/common/interfaces"
	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
)

// ToDateTime преобразует time.Time в доменный тип
func ToDateTime(t time.Time) valueobjects.DateTime {
	return valueobjects.NewDateTime(t)
}

// TodayDate создает Date объект из текущего времени
func TodayDate(clock interfaces.Clock) (valueobjects.Date, error) {
	year, month, day := clock.Today()
	return valueobjects.NewDate(year, month, day)
}
