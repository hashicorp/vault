// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"context"
	"time"
)

// MountAttribution holds the metric count and mount metadata for a specific mount.
type MountAttribution struct {
	MountPath        string      `json:"mount_path"`         // User-facing mount path
	MountType        string      `json:"mount_type"`         // Plugin type (e.g., "kv", "aws", "pki")
	MountAccessor    string      `json:"mount_accessor"`     // Mount accessor (also used as map key)
	NamespaceID      string      `json:"namespace_id"`       // Namespace identifier
	NamespacePath    string      `json:"namespace_path"`     // User-facing namespace path
	Count            interface{} `json:"count"`              // Count of a specific metric under this mount (int or float64)
	BackendAwareUUID string      `json:"backend_aware_uuid"` // A stable identifier that is unique across clusters
}

// MetricTypeAttribution holds mount attribution data for a specific metric type (e.g., "kv", "aws_static").
type MetricTypeAttribution struct {
	Count       interface{}                 `json:"count"`        // Total count for this metric type (int or float64)
	Mounts      map[string]MountAttribution `json:"mounts"`       // Map from mount accessor to per-mount breakdown
	LastUpdated time.Time                   `json:"last_updated"` // Last time the count for this metric type was updated
}

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
