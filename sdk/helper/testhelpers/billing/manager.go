// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"context"
	"sync"

	"github.com/hashicorp/vault/sdk/logical"
	"go.uber.org/atomic"
)

var _ logical.ConsumptionBillingManager = (*MockConsumptionBillingManager)(nil)

// MockConsumptionBillingManager is a test fake for logical.ConsumptionBillingManager.
// It records all WriteBillingData calls so tests can assert on billing counts, units,
// and per-mount attribution without wiring up the full Vault billing stack.
type MockConsumptionBillingManager struct {
	mu          sync.Mutex
	totalCount  *atomic.Uint64
	totalUnits  *atomic.Float64
	attribution map[string]logical.MountAttribution
}

// NewMockConsumptionBillingManager creates an initialized MockConsumptionBillingManager.
func NewMockConsumptionBillingManager() *MockConsumptionBillingManager {
	return &MockConsumptionBillingManager{
		totalCount:  atomic.NewUint64(0),
		totalUnits:  atomic.NewFloat64(0),
		attribution: make(map[string]logical.MountAttribution),
	}
}

// WriteBillingData accumulates billing data from a single plugin write. If data
// contains a "units" float64 key the float path is taken; otherwise "count" uint64 is used.
func (f *MockConsumptionBillingManager) WriteBillingData(_ context.Context, _ string, data map[string]interface{}) error {
	mountAccessor, _ := data["mountAccessor"].(string)
	mountPath, _ := data["mountPath"].(string)
	mountType, _ := data["mountType"].(string)
	mountRunningVersion, _ := data["mountRunningVersion"].(string)
	backendAwareUUID, _ := data["backendAwareUUID"].(string)

	var newCount interface{}

	if units, ok := data["units"].(float64); ok {
		f.totalUnits.Add(units)

		f.mu.Lock()
		defer f.mu.Unlock()
		prevUnits, _ := f.attribution[mountAccessor].Count.(float64)
		newCount = prevUnits + units
	} else {
		count, _ := data["count"].(uint64)
		f.totalCount.Add(count)

		f.mu.Lock()
		defer f.mu.Unlock()
		prevCount, _ := f.attribution[mountAccessor].Count.(uint64)
		newCount = prevCount + count
	}

	f.attribution[mountAccessor] = logical.MountAttribution{
		MountPath:           mountPath,
		MountAccessor:       mountAccessor,
		MountType:           mountType,
		MountRunningVersion: mountRunningVersion,
		BackendAwareUUID:    backendAwareUUID,
		Count:               newCount,
	}

	return nil
}

// TotalCount returns the cumulative count across all WriteBillingData calls with integer counts.
func (f *MockConsumptionBillingManager) TotalCount() uint64 {
	return f.totalCount.Load()
}

// TotalUnits returns the cumulative units across all WriteBillingData calls with float units.
func (f *MockConsumptionBillingManager) TotalUnits() float64 {
	return f.totalUnits.Load()
}

// GetAttribution returns the MountAttribution for the given mount accessor and whether it exists.
func (f *MockConsumptionBillingManager) GetAttribution(mountAccessor string) (logical.MountAttribution, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()

	attr, ok := f.attribution[mountAccessor]
	return attr, ok
}

// GetParentNamespaceID satisfies the logical.ConsumptionBillingManager interface.
// It always returns an empty string because namespace hierarchy is not relevant in test fakes.
func (f *MockConsumptionBillingManager) GetParentNamespaceID(_ string) string {
	return ""
}
