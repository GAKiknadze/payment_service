# Common (общие сущности)

Этот домен содержит общие сущности и объекты, используемые в нескольких доменах системы.

## Value Objects

### MoneyAmount

*Сумма денежных средств с указанием валюты.*

**Содержит:**
- `amount` Числовое значение суммы (неотрицательное)
- `currency` Код валюты (например, RUB, USD, EUR)

**Используется в:**
- Billing Domain (Payment)
- Organization Domain (Balance)
- Subscription Domain (цены тарифов)
- Tariff Domain (Price)

### PaymentMethod

*Платежный метод для проведения финансовых операций.*

**Содержит:**
- `id` Уникальный идентификатор платежного метода
- `type` Тип платежного метода (Card, BankTransfer, PayPal)
- `token` Токенизированные данные платежного метода
- `displayData` Данные для отображения (например, последние 4 цифры карты)
- `isDefault` Является ли методом по умолчанию
- `isValid` Валидность метода (не просрочен, активен)

**Используется в:**
- Organization Domain (PaymentMethods collection)
- Billing Domain (для проведения платежей)
- Subscription Domain (при создании подписки)

### BillingCycle

*Периодичность списания средств для подписок.*

**Содержит:**
- `cycleType` Тип периодичности (`Hourly`, `Monthly`, `OneTime`)
- `isRecurring` Является ли цикл повторяющимся
- `displayName` Отображаемое название (например, "Ежемесячно")

**Методы:**
- `NewBillingCycle (cycleType BillingCycleType) (BillingCycle, error)` Фабричный метод для создания BillingCycle
- `CalculateNextBillingDate(currentDate time.Time) (time.Time, error)` Метод для расчета следующей даты списания

**Возможные ошибки:**
- `ErrUnsupportedBillingCycleType` Тип периодичности не поддерживается

**Используется в:**
- Tariff Domain (описание тарифа)
- Subscription Domain (управление подписками)
- Billing Domain (расписание списаний)

### AutoTopUpSettings

*Настройки автоматического пополнения баланса.*

**Содержит:**
- `isEnabled` Активировано ли автопополнение
- `threshold` Порог срабатывания (рассчитывается как 10% от месячной суммы подписок)
- `topUpAmount` Сумма пополнения
- `paymentMethodId` Идентификатор платежного метода для автопополнения
- `isValid` Валидность настроек (если isEnabled, то остальные параметры должны быть валидны)

**Используется в:**
- Organization Domain (настройки организации)
- Billing Domain (автопополнение баланса)

### QuotaDefinition

*Определение квоты для тарифа.*

**Содержит:**
- `resourceType` Тип ресурса (например, tokens, ssl_certificates, api_calls)
- `limit` Максимальный лимит ресурса
- `unit` Единица измерения (например, count, mb, requests)
- `isRecurring` Периодичность сброса (true для периодических тарифов)
- `resetPeriod` Период сброса в днях (для периодических квот)

**Используется в:**
- Tariff Domain (лимиты тарифа)
- Subscription Domain (квоты активной подписки)
- Quota Domain (определение квот)

### QuotaUsage

*Текущее использование квоты.*

**Содержит:**
- `resourceType` Тип ресурса
- `used` Использовано единиц ресурса
- `limit` Максимальный лимит
- `organizationId` Идентификатор организации
- `subscriptionId` Идентификатор подписки
- `periodStart` Начало текущего периода использования
- `periodEnd` Конец текущего периода использования
- `resetDate` Дата следующего сброса квоты
- `status` Статус использования (`Normal`, `Warning`, `Exceeded`)

**Используется в:**
- Subscription Domain (отслеживание использования)
- Quota Domain (проверка и обновление квот)
- Organization Domain (просмотр общего использования)

### Currency

*Информация о поддерживаемой валюте.*

**Содержит:**
- `code` Код валюты (например, RUB, USD, EUR)
- `symbol` Символ валюты (например, ₽, $, €)
- `name` Название валюты (например, "Российский рубль")
- `decimalPlaces` Количество десятичных знаков
- `isSupported` Поддерживается ли валюта в системе

**Используется в:**
- Organization Domain (валюта организации)
- Tariff Domain (цены в разных валютах)
- Billing Domain (обработка платежей в разных валютах)

### Price

*Цена услуги или тарифа в определенной валюте.*

**Содержит:**
- `id` Уникальный идентификатор цены
- `amount` Сумма
- `isDefault` Является ли валютой по умолчанию

**Используется в:**
- Tariff Domain (цены тарифа)
- Subscription Domain (расчет стоимости)
- Billing Domain (обработка платежей)