package tariff_test

import (
	"testing"
	"time"

	"github.com/GAKiknadze/payment_service/domain/common/valueobject"
	"github.com/GAKiknadze/payment_service/domain/tariff"
	"github.com/shopspring/decimal"
)

// Вспомогательные функции для тестов
func createTestCurrency(currencyType valueobject.CurrencyType) valueobject.Currency {
	currency, _ := valueobject.NewCurrency(currencyType)
	return currency
}

func createTestPrice(id string, currencyType valueobject.CurrencyType, amount float64) valueobject.Price {
	currency := createTestCurrency(currencyType)
	moneyAmount, _ := valueobject.NewMoneyAmount(
		decimal.NewFromFloat(amount),
		currency,
	)
	price, _ := valueobject.NewPrice(id, moneyAmount, false)
	return price
}

func createTestQuota(resourceType string, limit float64) valueobject.QuotaDefinition {
	limitDecimal := decimal.NewFromFloat(limit)
	quota, _ := valueobject.NewQuotaDefinition(
		resourceType,
		limitDecimal,
		"count",
		true,
		30*24*time.Hour,
	)
	return quota
}

func createTestBillingCycle(cycleType valueobject.BillingCycleType) valueobject.BillingCycle {
	billingCycle, _ := valueobject.NewBillingCycle(cycleType)
	return billingCycle
}

func TestNewTariff_ValidParameters(t *testing.T) {
	// Given - валидные параметры для создания тарифа
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	// When - создаем тариф
	tar, err := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// Then - проверяем, что ошибка отсутствует и данные корректны
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if tar.ID().String() != id.String() {
		t.Errorf("Expected ID %s, got %s", id, tar.ID())
	}

	if tar.Name() != name {
		t.Errorf("Expected name %s, got %s", name, tar.Name())
	}

	if *tar.Description() != description {
		t.Errorf("Expected description %s, got %s", description, *tar.Description())
	}

	if tar.Status() != tariff.TariffStatusActive {
		t.Errorf("Expected status Active, got %s", tar.Status())
	}

	if tar.BillingCycle().Type() != billingCycle.Type() {
		t.Errorf("Expected billing cycle %s, got %s", billingCycle.Type(), tar.BillingCycle().Type())
	}

	if tar.IsExtendable() != isExtendable {
		t.Errorf("Expected isExtendable %v, got %v", isExtendable, tar.IsExtendable())
	}

	if tar.Prices()[0].Currency().Code() != "RUB" {
		t.Errorf("Expected currency RUB, got %s", tar.Prices()[0].Currency().Code())
	}

	if !tar.Quotas()[0].Limit().Equal(decimal.NewFromFloat(1000)) {
		t.Errorf("Expected quota limit 1000, got %s", tar.Quotas()[0].Limit())
	}

	if tar.Version() != 1 {
		t.Errorf("Expected version 1, got %d", tar.Version())
	}

	if len(tar.PopEvents()) != 1 {
		t.Error("Expected 1 event to be recorded")
	}
}

func TestNewTariff_EmptyID(t *testing.T) {
	// Given - пустой ID
	id, _ := valueobject.NewTariffID("")
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	// When - пытаемся создать тариф
	_, err := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for empty tariff ID, got nil")
	}

	if err.Error() != "tariff ID cannot be empty" {
		t.Errorf("Expected 'tariff ID cannot be empty', got %v", err)
	}
}

func TestNewTariff_EmptyName(t *testing.T) {
	// Given - пустое имя
	id := valueobject.GenerateTariffID()
	name := ""
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	// When - пытаемся создать тариф
	_, err := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for empty tariff name, got nil")
	}

	if err.Error() != "tariff name cannot be empty" {
		t.Errorf("Expected 'tariff name cannot be empty', got %v", err)
	}
}

func TestNewTariff_InvalidBillingCycle(t *testing.T) {
	// Given - невалидный billing cycle
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	// Создаем невалидный billing cycle
	billingCycle := valueobject.BillingCycle{}
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	// When - пытаемся создать тариф
	_, err := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for invalid billing cycle, got nil")
	}

	if err != tariff.ErrInvalidBillingCycle {
		t.Errorf("Expected ErrInvalidBillingCycle, got %v", err)
	}
}

func TestNewTariff_MissingPricesForPeriodic(t *testing.T) {
	// Given - периодический тариф без цен
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	// When - пытаемся создать тариф
	_, err := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for missing prices, got nil")
	}

	if err != tariff.ErrMissingPrices {
		t.Errorf("Expected ErrMissingPrices, got %v", err)
	}
}

func TestNewTariff_ValidOneTimeTariffWithoutPrices(t *testing.T) {
	// Given - разовый тариф без цен
	id := valueobject.GenerateTariffID()
	name := "SSL Certificate"
	description := "One-time SSL certificate"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleOneTime)
	isExtendable := true
	prices := []valueobject.Price{}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("ssl_certificates", 1),
	}

	// When - создаем тариф
	tariff, err := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// Then - проверяем, что ошибка отсутствует
	if err != nil {
		t.Fatalf("Expected no error for OneTime tariff without prices, got: %v", err)
	}

	if tariff.BillingCycle().Type() != valueobject.BillingCycleOneTime {
		t.Error("Expected OneTime billing cycle")
	}
}

func TestUpdateNameAndDescription_ValidUpdate(t *testing.T) {
	// Given - активный тариф
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	tar, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	newName := "Premium Plan"
	newDescription := "Updated description"

	// When - обновляем название и описание
	err := tar.UpdateNameAndDescription(newName, &newDescription)

	// Then - проверяем, что ошибка отсутствует
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if tar.Name() != newName {
		t.Errorf("Expected name %s, got %s", newName, tar.Name())
	}

	if *tar.Description() != newDescription {
		t.Errorf("Expected description %s, got %s", newDescription, *tar.Description())
	}

	if tar.Version() != 2 {
		t.Errorf("Expected version 2, got %d", tar.Version())
	}

	events := tar.PopEvents()
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	event, ok := events[0].(tariff.EventTariffUpdated)
	if !ok {
		t.Fatal("Expected EventTariffUpdated event")
	}

	if len(event.ChangedFields) != 2 {
		t.Errorf("Expected 2 changed fields, got %d", len(event.ChangedFields))
	}
}

func TestUpdateNameAndDescription_ArchivedTariff(t *testing.T) {
	// Given - архивный тариф
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	tar, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// Архивируем тариф
	reason := "Test reason"
	tar.Archive(&reason)

	newName := "Premium Plan"
	newDescription := "Updated description"

	// When - пытаемся обновить архивный тариф
	err := tar.UpdateNameAndDescription(newName, &newDescription)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for updating archived tariff, got nil")
	}

	if err != tariff.ErrArchivedTariff {
		t.Errorf("Expected ErrArchivedTariff, got %v", err)
	}
}

func TestUpdateNameAndDescription_NoChanges(t *testing.T) {
	// Given - активный тариф
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	tariff, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// When - пытаемся обновить на те же значения
	err := tariff.UpdateNameAndDescription(name, &description)

	// Then - проверяем, что ошибка отсутствует и нет изменений
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if tariff.Version() != 1 {
		t.Errorf("Expected version to remain 1, got %d", tariff.Version())
	}

	if len(tariff.PopEvents()) != 1 { // 1 событие - создание тарифа
		t.Error("Expected no additional events to be recorded")
	}
}

func TestAddPrice_ValidAdd(t *testing.T) {
	// Given - активный тариф с одной валютой
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	tar, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// Добавляем цену в новой валюте
	newPrice := createTestPrice("price_2", valueobject.CurrencyKZT, 15.75)

	// When - добавляем цену
	err := tar.AddPrice(newPrice, true)

	// Then - проверяем, что ошибка отсутствует
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(tar.Prices()) != 2 {
		t.Errorf("Expected 2 prices, got %d", len(tar.Prices()))
	}

	// Проверяем, что новая цена добавлена
	price, found := tar.GetPriceByCurrency("USD")
	if !found {
		t.Error("Expected new price to be found")
	}

	if !price.Amount().Amount().Equal(decimal.NewFromFloat(15.75)) {
		t.Errorf("Expected amount 15.75, got %s", price.Amount().Amount())
	}

	if !price.IsDefault() {
		t.Error("Expected new price to be default")
	}

	if tar.Version() != 2 {
		t.Errorf("Expected version 2, got %d", tar.Version())
	}

	events := tar.PopEvents()
	if len(events) != 2 { // 1 событие - создание тарифа, 1 - добавление цены
		t.Fatalf("Expected 2 events, got %d", len(events))
	}

	event, ok := events[1].(tariff.EventPriceAdded)
	if !ok {
		t.Fatal("Expected EventPriceAdded event")
	}

	if event.Currency != "USD" {
		t.Errorf("Expected currency USD, got %s", event.Currency)
	}
}

func TestAddPrice_CurrencyAlreadyExists(t *testing.T) {
	// Given - активный тариф с валютой RUB
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	tar, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// Пытаемся добавить цену в уже существующей валюте
	existingPrice := createTestPrice("price_2", valueobject.CurrencyRUB, 100.50)

	// When - добавляем цену
	err := tar.AddPrice(existingPrice, false)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for existing currency, got nil")
	}

	if err != tariff.ErrCurrencyAlreadyExists {
		t.Errorf("Expected ErrCurrencyAlreadyExists, got %v", err)
	}
}

func TestAddPrice_ArchivedTariff(t *testing.T) {
	// Given - архивный тариф
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	tar, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// Архивируем тариф
	reason := "Test reason"
	tar.Archive(&reason)

	// Добавляем цену в новой валюте
	newPrice := createTestPrice("price_2", valueobject.CurrencyKZT, 15.75)

	// When - пытаемся добавить цену
	err := tar.AddPrice(newPrice, true)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for archived tariff, got nil")
	}

	if err != tariff.ErrArchivedTariff {
		t.Errorf("Expected ErrArchivedTariff, got %v", err)
	}
}

func TestRemovePrice_ValidRemove(t *testing.T) {
	// Given - активный тариф с двумя валютами
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
		createTestPrice("price_2", valueobject.CurrencyKZT, 15.75),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	tar, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// When - удаляем валюту USD
	err := tar.RemovePrice("USD")

	// Then - проверяем, что ошибка отсутствует
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(tar.Prices()) != 1 {
		t.Errorf("Expected 1 price, got %d", len(tar.Prices()))
	}

	// Проверяем, что осталась только валюта RUB
	_, found := tar.GetPriceByCurrency("USD")
	if found {
		t.Error("Expected USD price to be removed")
	}

	_, found = tar.GetPriceByCurrency("RUB")
	if !found {
		t.Error("Expected RUB price to remain")
	}

	if tar.Version() != 2 {
		t.Errorf("Expected version 2, got %d", tar.Version())
	}

	events := tar.PopEvents()
	if len(events) != 2 { // 1 событие - создание тарифа, 1 - удаление цены
		t.Fatalf("Expected 2 events, got %d", len(events))
	}

	event, ok := events[1].(tariff.EventPriceRemoved)
	if !ok {
		t.Fatal("Expected EventPriceRemoved event")
	}

	if event.Currency != "USD" {
		t.Errorf("Expected currency USD, got %s", event.Currency)
	}
}

func TestRemovePrice_LastPrice(t *testing.T) {
	// Given - активный тариф с одной валютой
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	tar, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// When - пытаемся удалить единственную валюту
	err := tar.RemovePrice("RUB")

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for removing last price, got nil")
	}

	if err != tariff.ErrLastPriceRemoval {
		t.Errorf("Expected ErrLastPriceRemoval, got %v", err)
	}
}

func TestArchive_ValidArchive(t *testing.T) {
	// Given - активный тариф
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	tar, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	reason := "Test reason"

	// When - архивируем тариф
	err := tar.Archive(&reason)

	// Then - проверяем, что ошибка отсутствует
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !tar.IsArchived() {
		t.Error("Expected tariff to be archived")
	}

	if tar.Version() != 2 {
		t.Errorf("Expected version 2, got %d", tar.Version())
	}

	events := tar.PopEvents()
	if len(events) != 2 { // 1 событие - создание тарифа, 1 - архивация
		t.Fatalf("Expected 2 events, got %d", len(events))
	}

	event, ok := events[1].(tariff.EventTariffArchived)
	if !ok {
		t.Fatal("Expected EventTariffArchived event")
	}

	if event.Reason == nil || *event.Reason != reason {
		t.Errorf("Expected reason %s, got %v", reason, event.Reason)
	}
}

func TestArchive_AlreadyArchived(t *testing.T) {
	// Given - архивный тариф
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	tar, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// Архивируем тариф
	reason := "Test reason"
	tar.Archive(&reason)

	// When - пытаемся архивировать уже архивный тариф
	err := tar.Archive(&reason)

	// Then - проверяем, что получена ожидаемая ошибка
	if err == nil {
		t.Fatal("Expected error for already archived tariff, got nil")
	}

	if err != tariff.ErrTariffAlreadyArchived {
		t.Errorf("Expected ErrTariffAlreadyArchived, got %v", err)
	}
}

func TestIsActive(t *testing.T) {
	// Given - активный тариф
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	tariff, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// When & Then - проверяем, что тариф активен
	if !tariff.IsActive() {
		t.Error("Expected tariff to be active")
	}

	// Given - архивный тариф
	reason := "Test reason"
	tariff.Archive(&reason)

	// When & Then - проверяем, что тариф не активен
	if tariff.IsActive() {
		t.Error("Expected tariff not to be active")
	}
}

func TestIsArchived(t *testing.T) {
	// Given - активный тариф
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	tariff, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// When & Then - проверяем, что тариф не архивный
	if tariff.IsArchived() {
		t.Error("Expected tariff not to be archived")
	}

	// Given - архивный тариф
	reason := "Test reason"
	tariff.Archive(&reason)

	// When & Then - проверяем, что тариф архивный
	if !tariff.IsArchived() {
		t.Error("Expected tariff to be archived")
	}
}

func TestGetPriceByCurrency(t *testing.T) {
	// Given - тариф с двумя валютами
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
		createTestPrice("price_2", valueobject.CurrencyKZT, 15.75),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	tariff, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// When & Then - проверяем получение цены по существующей валюте
	price, found := tariff.GetPriceByCurrency("RUB")
	if !found {
		t.Error("Expected RUB price to be found")
	}

	if !price.Amount().Amount().Equal(decimal.NewFromFloat(100.50)) {
		t.Errorf("Expected amount 100.50, got %s", price.Amount().Amount())
	}

	// When & Then - проверяем получение цены по несуществующей валюте
	_, found = tariff.GetPriceByCurrency("EUR")
	if found {
		t.Error("Expected EUR price not to be found")
	}
}

func TestGetDefaultPrice(t *testing.T) {
	// Given - тариф с двумя валютами, одна из которых дефолтная
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false

	// Создаем цену RUB как дефолтную
	priceRUB := createTestPrice("price_1", valueobject.CurrencyRUB, 100.50)
	priceRUB, _ = valueobject.NewPrice(priceRUB.ID(), priceRUB.Amount(), true)

	prices := []valueobject.Price{
		priceRUB,
		createTestPrice("price_2", valueobject.CurrencyKZT, 15.75),
	}

	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	tar, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// When & Then - проверяем получение дефолтной цены
	defaultPrice, found := tar.GetDefaultPrice()
	if !found {
		t.Error("Expected default price to be found")
	}

	if defaultPrice.Currency().Code() != "RUB" {
		t.Errorf("Expected default currency RUB, got %s", defaultPrice.Currency().Code())
	}

	// Given - тариф без явно установленной дефолтной цены
	priceRUB, _ = valueobject.NewPrice(priceRUB.ID(), priceRUB.Amount(), false)
	prices = []valueobject.Price{
		priceRUB,
		createTestPrice("price_2", valueobject.CurrencyKZT, 15.75),
	}

	tar, _ = tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// When & Then - проверяем, что возвращается первая цена как дефолтная
	defaultPrice, found = tar.GetDefaultPrice()
	if !found {
		t.Error("Expected default price to be found")
	}

	if defaultPrice.Currency().Code() != "RUB" {
		t.Errorf("Expected default currency RUB, got %s", defaultPrice.Currency().Code())
	}
}

func TestHasPrices(t *testing.T) {
	// Given - тариф с ценами
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	tariffWithPrices, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// Given - тариф без цен
	id = valueobject.GenerateTariffID()
	name = "SSL Certificate"
	description = "One-time SSL certificate"
	billingCycle = createTestBillingCycle(valueobject.BillingCycleOneTime)
	isExtendable = true
	prices = []valueobject.Price{}
	quotas = []valueobject.QuotaDefinition{
		createTestQuota("ssl_certificates", 1),
	}

	tariffWithoutPrices, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// When & Then - проверяем наличие цен
	if !tariffWithPrices.HasPrices() {
		t.Error("Expected tariff with prices to return true")
	}

	if tariffWithoutPrices.HasPrices() {
		t.Error("Expected tariff without prices to return false")
	}
}

func TestGetQuotaDefinition(t *testing.T) {
	// Given - тариф с квотами
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
		createTestQuota("api_calls", 5000),
	}

	tariff, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// When & Then - проверяем получение квоты по существующему типу
	quota, found := tariff.GetQuotaDefinition("tokens")
	if !found {
		t.Error("Expected tokens quota to be found")
	}

	if quota.ResourceType() != "tokens" {
		t.Errorf("Expected resource type tokens, got %s", quota.ResourceType())
	}

	// When & Then - проверяем получение квоты по несуществующему типу
	_, found = tariff.GetQuotaDefinition("storage")
	if found {
		t.Error("Expected storage quota not to be found")
	}
}

func TestCanSupportSubscriptions(t *testing.T) {
	// Given - периодический тариф с ценами
	id := valueobject.GenerateTariffID()
	name := "Basic Plan"
	description := "Test description"
	billingCycle := createTestBillingCycle(valueobject.BillingCycleMonthly)
	isExtendable := false
	prices := []valueobject.Price{
		createTestPrice("price_1", valueobject.CurrencyRUB, 100.50),
	}
	quotas := []valueobject.QuotaDefinition{
		createTestQuota("tokens", 1000),
	}

	tariffWithPrices, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// Given - периодический тариф без цен
	id = valueobject.GenerateTariffID()
	name = "Premium Plan"
	description = "Test description"
	billingCycle = createTestBillingCycle(valueobject.BillingCycleHourly)
	isExtendable = false
	prices = []valueobject.Price{}
	quotas = []valueobject.QuotaDefinition{
		createTestQuota("tokens", 5000),
	}

	tariffWithoutPrices, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// Given - разовый тариф без цен
	id = valueobject.GenerateTariffID()
	name = "SSL Certificate"
	description = "One-time SSL certificate"
	billingCycle = createTestBillingCycle(valueobject.BillingCycleOneTime)
	isExtendable = true
	prices = []valueobject.Price{}
	quotas = []valueobject.QuotaDefinition{
		createTestQuota("ssl_certificates", 1),
	}

	oneTimeTariff, _ := tariff.NewTariff(
		id, name, &description, billingCycle, isExtendable, prices, quotas,
	)

	// When & Then - проверяем поддержку подписок
	if !tariffWithPrices.CanSupportSubscriptions() {
		t.Error("Expected tariff with prices to support subscriptions")
	}

	if tariffWithoutPrices.CanSupportSubscriptions() {
		t.Error("Expected tariff without prices not to support subscriptions")
	}

	if !oneTimeTariff.CanSupportSubscriptions() {
		t.Error("Expected OneTime tariff to support subscriptions")
	}
}
