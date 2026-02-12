// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package vault

func (c *Core) ResetInMemoryTransitDataProtectionCallCounts() {
	c.consumptionBilling.DataProtectionCallCounts.Transit.Store(0)
}

func (c *Core) GetInMemoryTransitDataProtectionCallCounts() uint64 {
	return c.consumptionBilling.DataProtectionCallCounts.Transit.Load()
}

func (c *Core) ResetInMemoryTransformDataProtectionCallCounts() {
	c.consumptionBilling.DataProtectionCallCounts.Transform.Store(0)
}

func (c *Core) GetInMemoryTransformDataProtectionCallCounts() uint64 {
	return c.consumptionBilling.DataProtectionCallCounts.Transform.Load()
}

func (c *Core) SetInMemoryTransitDataProtectionCallCounts(count uint64) {
	c.consumptionBilling.DataProtectionCallCounts.Transit.Store(count)
}

func (c *Core) SetInMemoryTransformDataProtectionCallCounts(count uint64) {
	c.consumptionBilling.DataProtectionCallCounts.Transform.Store(count)
}
