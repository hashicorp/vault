// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"fmt"
	"sync"
	"time"
)

const (
	BillingSubPath          = "billing/"
	ReplicatedPrefix        = "replicated/"
	RoleHWMCountsHWM        = "maxRoleCounts/"
	KvHWMCountsHWM          = "maxKvCounts/"
	LocalPrefix             = "local/"
	ThirdPartyPluginsPrefix = "thirdPartyPluginCounts/"
	BillingWriteInterval    = 10 * time.Minute
)

var BillingMonthStorageFormat = "%s%d/%02d/%s" // e.g replicated/2026/01/maxKvCounts/

type ConsumptionBilling struct {
	// BillingStorageLock controls access to the billing storage paths
	BillingStorageLock sync.RWMutex

	BillingConfig BillingConfig
}

type BillingConfig struct {
	// For testing purposes. The cadence at which billing metrics are updated
	MetricsUpdateCadence time.Duration
}

func GetMonthlyBillingPath(localPrefix string, now time.Time, billingMetric string) string {
	// Normalize to avoid double slashes since our prefixes include trailing "/".
	// Example: localPrefix="replicated/", billingMetric="maxKvCounts/" =>
	// "replicated/2026/01/maxKvCounts/"
	year := now.Year()
	month := int(now.Month())
	return fmt.Sprintf(BillingMonthStorageFormat, localPrefix, year, month, billingMetric)
}
