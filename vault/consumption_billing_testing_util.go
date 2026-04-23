// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package vault

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
