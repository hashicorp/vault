// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package vault

func (c *Core) ResetInMemoryDataProtectionCallCounts() {
	c.consumptionBilling.DataProtectionCallCounts.Transit.Store(0)
}

func (c *Core) GetInMemoryTransitDataProtectionCallCounts() uint64 {
	return c.consumptionBilling.DataProtectionCallCounts.Transit.Load()
}
