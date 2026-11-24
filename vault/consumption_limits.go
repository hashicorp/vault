// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

// consumptionLimitState is the result of comparing one metric's current-month
// total against its configured soft limit. Consumption limits are an
// Enterprise-only feature, so the state is only populated in Enterprise builds.
type consumptionLimitState struct {
	MetricName  string
	Total       float64
	Limit       float64
	Utilization float64 // Total / Limit
	Reached     bool    // Total >= Limit
}
