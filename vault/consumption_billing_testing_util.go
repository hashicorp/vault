// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (c *Core) ResetInMemoryTransitDataProtectionCallCounts() {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.DataProtectionCallCounts.Transit.Store(0)
	}
}

func (c *Core) GetInMemoryTransitDataProtectionCallCounts() uint64 {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		return cb.DataProtectionCallCounts.Transit.Load()
	}
	return 0
}

func (c *Core) ResetInMemoryTransformDataProtectionCallCounts() {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.DataProtectionCallCounts.Transform.Store(0)
	}
}

func (c *Core) GetInMemoryTransformDataProtectionCallCounts() uint64 {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		return cb.DataProtectionCallCounts.Transform.Load()
	}
	return 0
}

func (c *Core) SetInMemoryTransitDataProtectionCallCounts(count uint64) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.DataProtectionCallCounts.Transit.Store(count)
	}
}

func (c *Core) SetInMemoryTransformDataProtectionCallCounts(count uint64) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.DataProtectionCallCounts.Transform.Store(count)
	}
}

func (c *Core) SetInMemoryGcpKmsDataProtectionCallCounts(count uint64) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.DataProtectionCallCounts.GcpKms.Store(count)
	}
}

func (c *Core) GetInMemoryGcpKmsDataProtectionCallCounts() uint64 {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		return cb.DataProtectionCallCounts.GcpKms.Load()
	}
	return 0
}

func (c *Core) ResetInMemoryJwtSpiffeIdentityCounts() {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.IdentityTokenUnits.SpiffeJwt.Store(0)
	}
}

func (c *Core) GetInMemoryJwtSpiffeIdentityCounts() float64 {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		return cb.IdentityTokenUnits.SpiffeJwt.Load()
	}
	return 0
}

func (c *Core) SetInMemoryJwtSpiffeIdentityCounts(count float64) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.IdentityTokenUnits.SpiffeJwt.Store(count)
	}
}

func (c *Core) GetInMemoryOidcCounts() float64 {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		return cb.IdentityTokenUnits.OidcTokenDuration.Load()
	}
	return 0
}

func (c *Core) SetInMemoryOidcCounts(tokenDuration float64) {
	c.consumptionBillingLock.RLock()
	cb := c.consumptionBilling
	c.consumptionBillingLock.RUnlock()

	if cb != nil {
		cb.IdentityTokenUnits.OidcTokenDuration.Store(tokenDuration)
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
