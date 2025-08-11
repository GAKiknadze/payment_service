package tariff

import (
	"errors"
	"fmt"
	"time"

	common "github.com/GAKiknadze/payment_service/domain/common/valueobject"
)

type TariffStatus string

const (
	TariffStatusActive   TariffStatus = "Active"
	TariffStatusArchived TariffStatus = "Archived"
)

type Tariff struct {
	id           common.TariffID
	name         string
	description  *string
	status       TariffStatus
	billingCycle common.BillingCycle
	isExtendable bool
	createdAt    time.Time
	updatedAt    time.Time
	archivedAt   time.Time
	prices       []common.Price
	quotas       []common.QuotaDefinition
	version      uint
	events       []interface{}
}

// NewTariff создает новый активный тариф
func NewTariff(
	id common.TariffID,
	name string,
	description *string,
	billingCycle common.BillingCycle,
	isExtendable bool,
	prices []common.Price,
	quotas []common.QuotaDefinition,
) (*Tariff, error) {
	// Валидация обязательных параметров
	if id.String() == "" {
		return nil, errors.New("tariff ID cannot be empty")
	}

	if name == "" {
		return nil, errors.New("tariff name cannot be empty")
	}

	// Проверка валидности billing cycle
	if !isValidBillingCycle(billingCycle.Type()) {
		return nil, ErrInvalidBillingCycle
	}

	// Проверка наличия цен для периодических тарифов
	if (billingCycle.Type() == common.BillingCycleHourly ||
		billingCycle.Type() == common.BillingCycleMonthly) && len(prices) == 0 {
		return nil, ErrMissingPrices
	}

	// Проверка квот
	if err := validateQuotas(quotas, billingCycle); err != nil {
		return nil, err
	}

	// Создаем тариф
	now := time.Now()
	tariff := &Tariff{
		id:           id,
		name:         name,
		description:  description,
		status:       TariffStatusActive,
		billingCycle: billingCycle,
		isExtendable: isExtendable,
		createdAt:    now,
		updatedAt:    now,
		prices:       prices,
		quotas:       quotas,
		version:      1,
	}

	// Генерируем событие создания
	tariff.recordEvent(EventTariffCreated{
		TariffID:     id,
		Name:         name,
		BillingCycle: string(billingCycle.Type()),
		CreatedAt:    now,
		Prices:       prices,
		Quotas:       quotas,
	})

	return tariff, nil
}

// UpdateNameAndDescription обновляет название и описание тарифа
func (t *Tariff) UpdateNameAndDescription(name string, description *string) error {
	if t.status == TariffStatusArchived {
		return ErrArchivedTariff
	}

	if name == "" {
		return errors.New("tariff name cannot be empty")
	}

	if t.name == name && (description == nil && t.description == nil ||
		description != nil && t.description != nil && *description == *t.description) {
		// Нет изменений
		return nil
	}

	// Сохраняем старые значения для события
	oldName := t.name
	oldDescription := t.description

	// Обновляем данные
	t.name = name
	t.description = description
	t.updatedAt = time.Now()
	t.version++

	// Генерируем событие обновления
	t.recordEvent(EventTariffUpdated{
		TariffID:             t.id,
		ChangedFields:        getChangedFields(oldName, t.name, oldDescription, t.description),
		UpdatedAt:            t.updatedAt,
		RequiresNotification: false,
		NewVersion:           t.version,
	})

	return nil
}

// AddPrice добавляет цену в новой валюте
func (t *Tariff) AddPrice(price common.Price, isDefault bool) error {
	if t.status == TariffStatusArchived {
		return ErrArchivedTariff
	}

	// Проверяем, что валюта еще не существует
	for _, p := range t.prices {
		if p.Currency().Code() == price.Currency().Code() {
			return ErrCurrencyAlreadyExists
		}
	}

	// Создаем копию цены с обновленным флагом isDefault
	newPrice, err := common.NewPrice(price.ID(), price.Amount(), isDefault)
	if err != nil {
		return err
	}

	// Добавляем цену
	t.prices = append(t.prices, newPrice)
	t.updatedAt = time.Now()
	t.version++

	// Если это первая цена, делаем ее дефолтной
	if len(t.prices) == 1 {
		t.prices[0], _ = common.NewPrice(t.prices[0].ID(), t.prices[0].Amount(), true)
	}

	// Генерируем событие добавления цены
	t.recordEvent(EventPriceAdded{
		TariffID:   t.id,
		Currency:   price.Currency().Code(),
		Amount:     price.Amount().Amount().String(),
		IsDefault:  isDefault,
		AddedAt:    t.updatedAt,
		NewVersion: t.version,
	})

	return nil
}

// RemovePrice удаляет цену в указанной валюте
func (t *Tariff) RemovePrice(currencyCode string) error {
	if t.status == TariffStatusArchived {
		return ErrArchivedTariff
	}

	if len(t.prices) <= 1 {
		return ErrLastPriceRemoval
	}

	// Находим индекс цены с указанной валютой
	index := -1
	var removedPrice common.Price
	var wasDefault bool

	for i, p := range t.prices {
		if p.Currency().Code() == currencyCode {
			index = i
			removedPrice = p
			wasDefault = p.IsDefault()
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("price for currency %s not found", currencyCode)
	}

	// Удаляем цену
	t.prices = append(t.prices[:index], t.prices[index+1:]...)

	// Если удаляемая цена была дефолтной, устанавливаем новую дефолтную цену
	var newDefaultCurrency string
	if wasDefault && len(t.prices) > 0 {
		t.prices[0], _ = common.NewPrice(t.prices[0].ID(), t.prices[0].Amount(), true)
		newDefaultCurrency = t.prices[0].Currency().Code()
	}

	t.updatedAt = time.Now()
	t.version++

	// Генерируем событие удаления цены
	t.recordEvent(EventPriceRemoved{
		TariffID:           t.id,
		Currency:           currencyCode,
		Price:              removedPrice,
		WasDefault:         wasDefault,
		NewDefaultCurrency: newDefaultCurrency,
		RemovedAt:          t.updatedAt,
		NewVersion:         t.version,
	})

	return nil
}

// UpdateQuotas обновляет квоты тарифа
func (t *Tariff) UpdateQuotas(newQuotas []common.QuotaDefinition) error {
	if t.status == TariffStatusArchived {
		return ErrArchivedTariff
	}

	// Проверяем валидность новых квот
	if err := validateQuotas(newQuotas, t.billingCycle); err != nil {
		return err
	}

	// Сохраняем старые квоты для события
	oldQuotas := make([]common.QuotaDefinition, len(t.quotas))
	copy(oldQuotas, t.quotas)

	// Обновляем квоты
	t.quotas = newQuotas
	t.updatedAt = time.Now()
	t.version++

	// Генерируем событие обновления квот
	t.recordEvent(EventQuotasUpdated{
		TariffID:   t.id,
		OldQuotas:  oldQuotas,
		NewQuotas:  newQuotas,
		UpdatedAt:  t.updatedAt,
		NewVersion: t.version,
	})

	return nil
}

// Archive архивирует тариф
func (t *Tariff) Archive(reason *string) error {
	if t.status == TariffStatusArchived {
		return ErrTariffAlreadyArchived
	}

	archivedAt := time.Now()
	t.status = TariffStatusArchived
	t.archivedAt = archivedAt
	t.updatedAt = archivedAt
	t.version++

	// Генерируем событие архивации
	t.recordEvent(EventTariffArchived{
		TariffID:                 t.id,
		ArchivedAt:               archivedAt,
		Reason:                   reason,
		DeprecationDate:          archivedAt.AddDate(0, 1, 0), // Пример: уведомление за 1 месяц
		ActiveSubscriptionsCount: 0,
		NewVersion:               t.version, // Это значение будет обновлено позже
	})

	return nil
}

// IsActive проверяет, является ли тариф активным
func (t *Tariff) IsActive() bool {
	return t.status == TariffStatusActive
}

// IsArchived проверяет, является ли тариф архивным
func (t *Tariff) IsArchived() bool {
	return t.status == TariffStatusArchived
}

// GetPriceByCurrency возвращает цену для указанной валюты
func (t *Tariff) GetPriceByCurrency(currencyCode string) (common.Price, bool) {
	for _, price := range t.prices {
		if price.Currency().Code() == currencyCode {
			return price, true
		}
	}
	return common.Price{}, false
}

// GetDefaultPrice возвращает цену по умолчанию
func (t *Tariff) GetDefaultPrice() (common.Price, bool) {
	for _, price := range t.prices {
		if price.IsDefault() {
			return price, true
		}
	}

	// Если дефолтная цена не установлена, возвращаем первую доступную
	if len(t.prices) > 0 {
		return t.prices[0], true
	}

	return common.Price{}, false
}

// HasPrices проверяет наличие цен
func (t *Tariff) HasPrices() bool {
	return len(t.prices) > 0
}

// GetQuotaDefinition возвращает определение квоты для указанного типа ресурса
func (t *Tariff) GetQuotaDefinition(resourceType string) (common.QuotaDefinition, bool) {
	for _, quota := range t.quotas {
		if quota.ResourceType() == resourceType {
			return quota, true
		}
	}
	return common.QuotaDefinition{}, false
}

// CanSupportSubscriptions проверяет, может ли тариф поддерживать подписки
func (t *Tariff) CanSupportSubscriptions() bool {
	// Периодические тарифы должны иметь цены
	if (t.billingCycle.Type() == common.BillingCycleHourly ||
		t.billingCycle.Type() == common.BillingCycleMonthly) && len(t.prices) == 0 {
		return false
	}

	return true
}

func (t Tariff) ID() common.TariffID {
	return t.id
}

func (t Tariff) Status() TariffStatus {
	return t.status
}

func (t Tariff) IsExtendable() bool {
	return t.isExtendable
}

func (t Tariff) Name() string {
	return t.name
}

func (t Tariff) Description() *string {
	return t.description
}

func (t Tariff) Version() uint {
	return t.version
}

func (t Tariff) BillingCycle() common.BillingCycle {
	return t.billingCycle
}

func (t Tariff) Prices() []common.Price {
	return t.prices
}

func (t Tariff) Quotas() []common.QuotaDefinition {
	return t.quotas
}

// PopEvents извлекает и сбрасывает буфер доменных событий
func (t Tariff) PopEvents() []interface{} {
	events := t.events
	t.events = nil
	return events
}

// recordEvent добавляет событие в буфер
func (t Tariff) recordEvent(event interface{}) {
	t.events = append(t.events, event)
}
