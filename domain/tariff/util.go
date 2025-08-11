package tariff

import "github.com/GAKiknadze/payment_service/domain/common/valueobject"

func isValidBillingCycle(cycleType valueobject.BillingCycleType) bool {
	return cycleType == valueobject.BillingCycleHourly ||
		cycleType == valueobject.BillingCycleMonthly ||
		cycleType == valueobject.BillingCycleOneTime
}

func validateQuotas(quotas []valueobject.QuotaDefinition, billingCycle valueobject.BillingCycle) error {
	// Для OneTime тарифов квоты могут быть фиксированными
	if billingCycle.Type() == valueobject.BillingCycleOneTime {
		return nil
	}

	// Проверяем, что все квоты являются периодическими
	for _, quota := range quotas {
		if !quota.IsRecurring() {
			return ErrIncompatibleQuotaDefinition
		}
	}

	return nil
}

func getChangedFields(oldName, newName string, oldDesc, newDesc *string) []string {
	changes := []string{}

	if oldName != newName {
		changes = append(changes, "name")
	}

	if (oldDesc == nil && newDesc != nil) ||
		(oldDesc != nil && newDesc == nil) ||
		(oldDesc != nil && newDesc != nil && *oldDesc != *newDesc) {
		changes = append(changes, "description")
	}

	return changes
}
