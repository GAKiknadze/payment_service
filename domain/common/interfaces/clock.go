package interfaces

import "time"

// Clock абстрагирует работу со временем в домене
type Clock interface {
	// Возвращаем стандартный time.Time - это разрешено в DDD
	// так как это фундаментальный тип, а не инфраструктурная деталь
	Now() time.Time

	// Для календарных операций используем упрощенный интерфейс
	Today() (year int, month time.Month, day int)
}
