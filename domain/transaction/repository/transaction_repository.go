package repository

import (
	"context"
	"time"

	"github.com/GAKiknadze/payment_service/domain/common/valueobjects"
	"github.com/GAKiknadze/payment_service/domain/transaction/model"
)

// TransactionRepository - интерфейс доступа к данным транзакций
type TransactionRepository interface {
	// Save сохраняет транзакцию с управлением версиями
	Save(ctx context.Context, transaction *model.Transaction) error

	// FindByID находит транзакцию по идентификатору
	FindByID(ctx context.Context, id valueobjects.TransactionID) (*model.Transaction, error)

	// FindByOrganization возвращает транзакции организации в указанном периоде
	FindByOrganization(
		ctx context.Context,
		orgID valueobjects.OrganizationID,
		start, end valueobjects.DateTime,
	) ([]*model.Transaction, error)

	// FindPending возвращает зависшие транзакции для мониторинга
	FindPending(ctx context.Context, maxAge time.Duration) ([]*model.Transaction, error)

	// FindByidempotencyKey проверяет существование транзакции с таким ключом
	FindByidempotencyKey(
		ctx context.Context,
		key valueobjects.IdempotencyKey,
	) (*model.Transaction, error)
}
