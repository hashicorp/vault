// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package logical

import "context"

// billing.ConsumptionBillingManager is an implementation of this interface that the backend can use to write billing data.
type ConsumptionBillingManager interface {
	WriteBillingData(ctx context.Context, pluginType string, data map[string]interface{}) error
}

// ================================
// Creates a null consumption billing manager that does nothing
var _ ConsumptionBillingManager = (*nullConsumptionBillingManager)(nil)

func NewNullConsumptionBillingManager() ConsumptionBillingManager {
	return &nullConsumptionBillingManager{}
}

type nullConsumptionBillingManager struct{}

func (n *nullConsumptionBillingManager) WriteBillingData(ctx context.Context, pluginType string, data map[string]interface{}) error {
	return nil
}
