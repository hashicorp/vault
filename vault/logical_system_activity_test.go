// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestAlignToBillingPeriod tests alignToBillingPeriod for correct alignment of the start time
func TestAlignToBillingPeriod(t *testing.T) {
	// assume the billing start time is February 1, 2023
	billingStartTime := time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		givenTime time.Time
		expected  time.Time
	}{
		"same as billing start time: given start time is February 1, 2023": {
			givenTime: time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC),
			expected:  billingStartTime,
		},
		"within the same billing period: given start time is February 3, 2023": {
			givenTime: time.Date(2023, time.June, 3, 0, 0, 0, 0, time.UTC),
			expected:  billingStartTime,
		},
		"before the current billing period: given start time is December 20, 2022": {
			givenTime: time.Date(2022, time.December, 20, 0, 0, 0, 0, time.UTC),
			expected:  time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
		"exactly one year before the billing start time: given start time is February 1, 2022": {
			givenTime: time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC),
			expected:  time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
		"several years in the past: given start time is March 10, 2015": {
			givenTime: time.Date(2015, time.March, 10, 0, 0, 0, 0, time.UTC),
			expected:  time.Date(2015, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
		"several years in the future: given start time is April 15, 2029": {
			givenTime: time.Date(2029, time.April, 15, 0, 0, 0, 0, time.UTC),
			expected:  billingStartTime,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := alignToBillingPeriodStart(billingStartTime, tt.givenTime)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
