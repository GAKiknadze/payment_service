# Сервис оплаты: Окончательный список функций

## Оглавление
1. [OrganizationAppService](#organizationappservice)
2. [SubscriptionAppService](#subscriptionappservice)
3. [TariffAppService](#tariffappservice)
4. [BillingAppService](#billingappservice)
5. [QuotaAppService](#quotaappservice)

---

## OrganizationAppService
*Управление организациями с разделением прав доступа между пользователями и администраторами.*

### **Управление организациями**
- [**CreateOrganization**](./organization.md#createorganization)  
  Создание новой организации с фиксированной валютой (владелец может создать только одну организацию).

- [**GetOrganization**](./organization.md#getorganization)  
  Получение данных об организации (владелец видит только свою, администратор — любую).

- [**UpdateOrganization**](./organization.md#updateorganization)  
  Обновление данных организации (владелец — только название, администратор — все поля кроме валюты).

- [**DeleteOrganization**](./organization.md#deleteorganization)  
  Удаление организации (владелец — при нулевом балансе и отсутствии подписок, администратор — с принудительным удалением).

### **Управление балансом**
- [**TopUpBalance**](./organization.md#topupbalance)  
  Ручное пополнение баланса организации владельца.

- [**UpdateAutoTopUpSettings**](./organization.md#updateautotopupsettings)  
  Настройка параметров автопополнения (активация, сумма пополнения, платежный метод).

- [**ManagePaymentMethods**](./organization.md#managepaymentmethods)  
  Добавление и удаление платежных методов для организации.

- [**GetQuotaUsage**](./organization.md#getquotausage)  
  Просмотр текущего использования квот по всем ресурсам.

---

## SubscriptionAppService
*Управление подписками и тарифами с поддержкой апгрейда, даунгрейда и отмены.*

### **Управление подписками**
- [**GetActiveSubscription**](./subscription.md#getactivesubscription)  
  Получение списка активных подписок организации.

- [**CreateSubscription**](./subscription.md#createsubscription)  
  Подключение нового тарифа с предварительной оплатой.

- [**ConfirmSubscriptionPayment**](./subscription.md#confirmsubscriptionpayment)  
  Подтверждение оплаты и активация подписки после успешного платежа.

- [**CancelSubscription**](./subscription.md#cancelsubscription)  
  Отмена активной подписки с расчетом возврата средств по выбранной политике.

- [**UpgradeSubscription**](./subscription.md#upgradesubscription)  
  Смена тарифа на более дорогой с немедленным списанием разницы.

- [**DowngradeSubscription**](./subscription.md#downgradesubscription)  
  Смена тарифа на более дешевый с применением изменений с даты следующего списания.

- [**ExtendSubscription**](./subscription.md#extendsubscription)  
  Продление разовой подписки (например, SSL-сертификата).

---

## TariffAppService
*Управление тарифами (только для администратора).*

### **Управление тарифами**
- [**GetTariffsList**](./tariff.md#gettariffslist)  
  Получение списка активных тарифов (с опцией включения архивных).

- [**GetTariff**](./tariff.md#gettariff)  
  Получение детальной информации о тарифе.

- [**CreateTariff**](./tariff.md#createtariff)  
  Создание нового тарифа с ценами в разных валютах и лимитами ресурсов.

- [**UpdateTariff**](./tariff.md#updatetariff)  
  Обновление параметров тарифа (название, описание, цены, квоты).

- [**AddPriceToTariff**](./tariff.md#addpricetotariff)  
  Добавление цены в новой валюте к существующему тарифу.

- [**RemovePriceFromTariff**](./tariff.md#removepricefromtariff)  
  Удаление цены в определенной валюте из тарифа.

- [**ArchiveTariff**](./tariff.md#archivetariff)  
  Архивирование тарифа (запрет на новые подключения).

---

## BillingAppService
*Управление платежами и биллингом (администратор и системные процессы).*

### **Платежные операции**
- [**ProcessScheduledBilling**](./billing.md#processscheduledbilling)  
  Автоматическое списание средств по расписанию для всех активных подписок.

- [**TriggerAutoTopUp**](./billing.md#triggerautotopup)  
  Принудительный запуск автопополнения для тестирования.

- [**GetPaymentHistory**](./billing.md#getpaymenthistory)  
  Получение истории платежей с фильтрацией по организации, периоду и типу.

- [**GenerateBillingReport**](./billing.md#generatebillingreport)  
  Генерация финансовых отчетов за заданный период.

- [**ManualCharge**](./billing.md#manualcharge)  
  Ручное списание средств за разовое использование ресурса.

- [**RetryFailedPayment**](./billing.md#retryfailedpayment)  
  Повторная попытка обработки неудачного платежа.

- [**TopUpBalance**](./organization.md#topupbalance)  
  Ручное пополнение баланса (дублируется из OrganizationAppService для администратора).

---

## QuotaAppService
*Управление квотами и лимитами ресурсов (автоматизировано и доступно для диагностики).*

### **Контроль использования ресурсов**
- [**CheckQuotaUsage**](./quota.md#checkquotausage)  
  Проверка возможности использования ресурса с учетом текущих квот.

- [**IncrementUsage**](./quota.md#incrementusage)  
  Увеличение использования квоты при выполнении операции (вызывается системой автоматически).

- [**ResetQuotaUsage**](./quota.md)  
  Сброс использования квоты при начале нового биллингового периода.

- [**GetQuotaLimits**](./quota.md) 
  Получение лимитов квот для всех активных подписок организации.

- [**AdjustQuotaManually**](./quota.md)  
  Ручная корректировка использования квоты (только для администратора при ошибках системы).
