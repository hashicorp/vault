// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"encoding/json"
	"maps"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/billing"
)

func (c *Core) ResetInMemoryTransitDataProtectionCallCounts() {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.SecretEngineCounts.Transit.MonthlyCount.Store(0)
	}
}

func (c *Core) GetInMemoryTransitDataProtectionCallCounts() uint64 {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		return cb.SecretEngineCounts.Transit.MonthlyCount.Load()
	}
	return 0
}

func (c *Core) ResetInMemoryTransformDataProtectionCallCounts() {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.SecretEngineCounts.Transform.MonthlyCount.Store(0)
	}
}

func (c *Core) GetInMemoryTransformDataProtectionCallCounts() uint64 {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		return cb.SecretEngineCounts.Transform.MonthlyCount.Load()
	}
	return 0
}

func (c *Core) SetInMemoryTransitDataProtectionCallCounts(count uint64) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.SecretEngineCounts.Transit.MonthlyCount.Store(count)
	}
}

func (c *Core) SetInMemoryTransformDataProtectionCallCounts(count uint64) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.SecretEngineCounts.Transform.MonthlyCount.Store(count)
	}
}

func (c *Core) SetInMemoryGcpKmsDataProtectionCallCounts(count uint64) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.SecretEngineCounts.GcpKms.MonthlyCount.Store(count)
	}
}

func (c *Core) GetInMemoryGcpKmsDataProtectionCallCounts() uint64 {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		return cb.SecretEngineCounts.GcpKms.MonthlyCount.Load()
	}
	return 0
}

func (c *Core) ResetInMemoryJwtSpiffeIdentityCounts() {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.SecretEngineCounts.Spiffe.MonthlyUnits.Store(0)
	}
}

func (c *Core) GetInMemoryJwtSpiffeIdentityCounts() float64 {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		return cb.SecretEngineCounts.Spiffe.MonthlyUnits.Load()
	}
	return 0
}

func (c *Core) SetInMemoryJwtSpiffeIdentityCounts(count float64) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.SecretEngineCounts.Spiffe.MonthlyUnits.Store(count)
	}
}

func (c *Core) GetInMemoryOidcCounts() float64 {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		return cb.SecretEngineCounts.Oidc.MonthlyUnits.Load()
	}
	return 0
}

// SetInMemoryOidcCounts sets the in-memory OIDC duration-adjusted units counter.
// The value should be pre-normalized (i.e. already converted via DurationAdjustedTokenCount).
func (c *Core) SetInMemoryOidcCounts(normalizedUnits float64) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.SecretEngineCounts.Oidc.MonthlyUnits.Store(normalizedUnits)
	}
}

// NewMockOSBackendFactory creates a mock OS backend factory for testing.
// The backend implements LIST operations for hosts and accounts to support
// billing enumeration testing.
func NewMockOSBackendFactory() logical.Factory {
	return func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
		b := &framework.Backend{
			BackendType: logical.TypeLogical,
			Paths: []*framework.Path{
				{
					Pattern: "hosts/?$",
					Operations: map[logical.Operation]framework.OperationHandler{
						logical.ListOperation: &framework.PathOperation{
							Callback: func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
								// List all hosts from storage
								hosts, err := req.Storage.List(ctx, "hosts/")
								if err != nil {
									return nil, err
								}
								return logical.ListResponse(hosts), nil
							},
						},
					},
				},
				{
					Pattern: "hosts/" + framework.GenericNameRegex("host") + "/accounts/?$",
					Fields: map[string]*framework.FieldSchema{
						"host": {
							Type:        framework.TypeString,
							Description: "Host name",
						},
					},
					Operations: map[logical.Operation]framework.OperationHandler{
						logical.ListOperation: &framework.PathOperation{
							Callback: func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
								// Get the host name from the path
								hostName := data.Get("host").(string)

								// Read the host entry from storage
								entry, err := req.Storage.Get(ctx, "hosts/"+hostName)
								if err != nil {
									return nil, err
								}
								if entry == nil {
									return logical.ListResponse([]string{}), nil
								}

								// Parse the JSON to extract account names
								var hostData map[string]interface{}
								if err := json.Unmarshal(entry.Value, &hostData); err != nil {
									return nil, err
								}

								accounts := []string{}
								if accountsMap, ok := hostData["accounts"].(map[string]interface{}); ok {
									for accountName := range accountsMap {
										accounts = append(accounts, accountName)
									}
								}

								return logical.ListResponse(accounts), nil
							},
						},
					},
				},
			},
		}
		if err := b.Setup(ctx, conf); err != nil {
			return nil, err
		}
		return b, nil
	}
}

// CreateMockOSHost creates a mock OS host entry in storage with the specified accounts.
// This is a helper function for tests that need to populate OS backend storage.
func CreateMockOSHost(ctx context.Context, storage logical.Storage, hostName string, accountNames []string) error {
	// Build the accounts map structure
	accountsMap := make(map[string]interface{})
	for _, accountName := range accountNames {
		accountsMap[accountName] = map[string]string{"username": "testuser"}
	}

	// Create the host entry structure and marshal to JSON
	hostMap := map[string]interface{}{"accounts": accountsMap}
	value, err := json.Marshal(hostMap)
	if err != nil {
		return err
	}

	// Create a mock host entry with accounts
	hostEntry := &logical.StorageEntry{
		Key:   "hosts/" + hostName,
		Value: value,
	}
	return storage.Put(ctx, hostEntry)
}

// DeleteMockOSHost deletes a mock OS host entry from storage.
// This is a helper function for tests that need to clean up OS backend storage.
func DeleteMockOSHost(ctx context.Context, storage logical.Storage, hostName string) error {
	return storage.Delete(ctx, "hosts/"+hostName)
}

func (c *Core) ResetInMemoryExternalCaCertUnits() {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.SecretEngineCounts.ExternalCa.MonthlyUnits.Store(0)
	}
}

func (c *Core) GetInMemoryExternalCaCertUnits() float64 {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		return cb.SecretEngineCounts.ExternalCa.MonthlyUnits.Load()
	}
	return 0
}

func (c *Core) SetInMemoryExternalCaCertUnits(count float64) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.SecretEngineCounts.ExternalCa.MonthlyUnits.Store(count)
	}
}

func (c *Core) GetInMemoryTransitAttribution() map[string]logical.MountAttribution {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.SecretEngineCounts.Transit.MountAttributionLock.RLock()
		defer cb.SecretEngineCounts.Transit.MountAttributionLock.RUnlock()
		return maps.Clone(cb.SecretEngineCounts.Transit.MountAttribution)
	}
	return nil
}

func (c *Core) GetInMemoryTransformAttribution() map[string]logical.MountAttribution {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.SecretEngineCounts.Transform.MountAttributionLock.RLock()
		defer cb.SecretEngineCounts.Transform.MountAttributionLock.RUnlock()
		return maps.Clone(cb.SecretEngineCounts.Transform.MountAttribution)
	}
	return nil
}

func (c *Core) GetInMemoryGcpKmsAttribution() map[string]logical.MountAttribution {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.SecretEngineCounts.GcpKms.MountAttributionLock.RLock()
		defer cb.SecretEngineCounts.GcpKms.MountAttributionLock.RUnlock()
		return maps.Clone(cb.SecretEngineCounts.GcpKms.MountAttribution)
	}
	return nil
}

func (c *Core) GetInMemoryExternalCaAttribution() map[string]logical.MountAttribution {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.SecretEngineCounts.ExternalCa.MountAttributionLock.RLock()
		defer cb.SecretEngineCounts.ExternalCa.MountAttributionLock.RUnlock()
		return maps.Clone(cb.SecretEngineCounts.ExternalCa.MountAttribution)
	}
	return nil
}

func (c *Core) GetInMemorySpiffeAttribution() map[string]logical.MountAttribution {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.SecretEngineCounts.Spiffe.MountAttributionLock.RLock()
		defer cb.SecretEngineCounts.Spiffe.MountAttributionLock.RUnlock()
		return maps.Clone(cb.SecretEngineCounts.Spiffe.MountAttribution)
	}
	return nil
}

// GetConsumptionBillingManagerConcrete returns the underlying *billing.ConsumptionBilling,
// available to tests outside the vault package that need to inspect or manipulate raw state.
func (c *Core) GetCoreConsumptionBillingManager() *billing.ConsumptionBilling {
	return c.consumptionBilling
}

func (c *Core) GetInMemoryOidcAttribution() map[string]logical.MountAttribution {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.SecretEngineCounts.Oidc.MountAttributionLock.RLock()
		defer cb.SecretEngineCounts.Oidc.MountAttributionLock.RUnlock()
		return maps.Clone(cb.SecretEngineCounts.Oidc.MountAttribution)
	}
	return nil
}
