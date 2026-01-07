// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"fmt"
	"sync"
	"time"
)

const (
	BillingSubPath       = "billing/"
	ReplicatedPrefix     = "replicated/"
	RoleHWMCountsHWM     = "maxRoleCounts/"
	LocalPrefix          = "local/"
	BillingWriteInterval = 10 * time.Minute
)

var BillingMonthStorageFormat = "%s/%d/%02d/%s"

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
	year := now.Year()
	month := now.Month()
	return fmt.Sprintf(localPrefix, month, year, billingMetric)
}
