package repository

import (
	"context"

	"github.com/GAKiknadze/payment_service/domain/common/interfaces"
	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
	"github.com/GAKiknadze/payment_service/domain/tariff/model"
)

// TariffRepository определяет контракт доступа к данным
// Чистый доменный интерфейс без технических деталей
type TariffRepository interface {
	// FindCurrent возвращает активный тариф на текущий момент
	FindCurrent(ctx context.Context, clock interfaces.Clock) (*model.Tariff, error)

	// FindByPeriod ищет тарифы в указанном временном диапазоне
	FindByPeriod(
		ctx context.Context,
		start, end valueobjects.DateTime,
	) ([]*model.Tariff, error)

	// Save сохраняет агрегат с управлением версиями
	Save(ctx context.Context, tariff *model.Tariff) error

	// FindByID получает тариф по идентификатору
	FindByID(ctx context.Context, id valueobjects.TariffID) (*model.Tariff, error)
}
