// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

// evaluateConsumptionLimits is a no-op on OSS
func evaluateConsumptionLimits(totals map[string]float64, limits map[string]float64) []consumptionLimitState {
	return nil
}
