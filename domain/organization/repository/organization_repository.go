package repository

import (
	"context"

	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
	"github.com/GAKiknadze/payment_service/domain/organization/model"
)

// OrganizationRepository определяет контракт доступа к данным организаций
type OrganizationRepository interface {
	// FindByID получает организацию по идентификатору
	FindByID(ctx context.Context, id valueobjects.OrganizationID) (*model.Organization, error)

	// FindByStatus ищет организации по статусу
	FindByStatus(ctx context.Context, status model.OrganizationStatus) ([]*model.Organization, error)

	// FindForBilling возвращает организации, требующие списания
	FindForBilling(ctx context.Context, maxBillingTime valueobjects.DateTime) ([]*model.Organization, error)

	// Save сохраняет агрегат с управлением версиями
	Save(ctx context.Context, org *model.Organization) error

	// UpdateStatus обновляет статус организации (опционально)
	UpdateStatus(ctx context.Context, id valueobjects.OrganizationID, status model.OrganizationStatus) error
}
