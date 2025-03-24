// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/stretchr/testify/require"
)

// TestAlignToBillingPeriodStart tests alignToBillingPeriodStart for correct alignment of the start time
func TestAlignToBillingPeriodStart(t *testing.T) {
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
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestGetBillingPeriodTimes(t *testing.T) {
	// establish the current billing period
	currentTime := time.Now().UTC()
	// assume the date of the current billing period is 7 months prior to current time, but with the day set to 1 and the rest to 0
	sevenMonthsAgo := currentTime.AddDate(0, -7, 0).UTC()
	currentBillingStartDate := time.Date(sevenMonthsAgo.Year(), sevenMonthsAgo.Month(), 1, 0, 0, 0, 0, time.UTC)
	currentBillingEndDate := currentBillingStartDate.AddDate(1, 0, 0).UTC()

	tests := map[string]struct {
		givenStartTime     time.Time
		givenEndTime       time.Time
		expectedStartTime  time.Time
		expectedEndTime    time.Time
		expectAlignedTimes bool
	}{
		"No start or end time provided, should return times from current billing period start time to now": {
			givenStartTime:     time.Time{},
			givenEndTime:       time.Time{},
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: false,
		},
		"Start time provided to be 2 months after the current billing period start time within the current billing cycle": {
			givenStartTime: func() time.Time {
				twoMonthsForward := currentBillingStartDate.AddDate(0, 2, 0)
				return time.Date(twoMonthsForward.Year(), twoMonthsForward.Month(), 10, 0, 0, 0, 0, time.UTC)
			}(),
			givenEndTime:       time.Time{},
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: true,
		},
		"Start time provided to be 5 months after the current time outside the current billing cycle": {
			givenStartTime: func() time.Time {
				fiveMonthsForward := currentTime.AddDate(0, 5, 0)
				return time.Date(fiveMonthsForward.Year(), fiveMonthsForward.Month(), 2, 2, 0, 0, 0, time.UTC)
			}(),
			givenEndTime:       time.Time{},
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: true,
		},
		"Start time provided to be 3 months after the current billing cycle end time outside the current billing cycle": {
			givenStartTime: func() time.Time {
				threeMonthsForward := currentBillingEndDate.AddDate(0, 3, 0)
				return time.Date(threeMonthsForward.Year(), threeMonthsForward.Month(), 1, 23, 0, 0, 0, time.UTC)
			}(),
			givenEndTime:       time.Time{},
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: true,
		},
		"Start time provided to be 1 year before the current time outside the current billing cycle": {
			givenStartTime: func() time.Time {
				oneYearAgo := currentTime.AddDate(-1, 0, 0)
				return time.Date(oneYearAgo.Year(), oneYearAgo.Month(), 1, 0, 0, 0, 0, time.UTC)
			}(),
			givenEndTime:       time.Time{},
			expectedStartTime:  currentBillingStartDate.AddDate(-1, 0, 0),
			expectedEndTime:    currentBillingEndDate.AddDate(-1, 0, 0),
			expectAlignedTimes: true,
		},
		"End time provided to be 5 months after the current billing period start time within the current billing cycle": {
			givenStartTime: time.Time{},
			givenEndTime: func() time.Time {
				fiveMonthsForward := currentBillingStartDate.AddDate(0, 5, 0)
				return time.Date(fiveMonthsForward.Year(), fiveMonthsForward.Month(), 23, 5, 0, 0, 0, time.UTC)
			}(),
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: true,
		},
		// "Start and end time in the same billing period": {
		// 	givenStartTime:     time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC),
		// 	givenEndTime:       time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
		// 	expectedStartTime:  currentBillingStartDate,
		// 	expectedEndTime:    currentBillingEndDate,
		// 	expectAlignedTimes: false,
		// },
		// "Start time at exact billing start, end time ignored": {
		// 	givenStartTime:     time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC),
		// 	givenEndTime:       time.Time{},
		// 	expectedStartTime:  currentBillingStartDate,
		// 	expectedEndTime:    currentBillingEndDate,
		// 	expectAlignedTimes: false,
		// },
		// "End time in current billing period, should cap at today": {
		// 	givenStartTime:     time.Time{},
		// 	givenEndTime:       time.Date(2023, time.May, 15, 0, 0, 0, 0, time.UTC),
		// 	expectedStartTime:  currentBillingStartDate,
		// 	expectedEndTime:    currentBillingEndDate,
		// 	expectAlignedTimes: true,
		// },
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			d := &framework.FieldData{
				Schema: map[string]*framework.FieldSchema{
					"start_time": {
						Type: framework.TypeTime,
					},
					"end_time": {
						Type: framework.TypeTime,
					},
				},
				Raw: map[string]any{
					"start_time": tt.givenStartTime.Format(time.RFC3339Nano),
					"end_time":   tt.givenEndTime.Format(time.RFC3339Nano),
				},
			}

			actualStart, actualEnd, actualAligned, err := getBillingPeriodTimes(d, currentBillingStartDate)
			require.NoError(t, err)
			currentTime = time.Now().UTC()

			require.Equal(t, tt.expectedStartTime, actualStart, "Expected start time did not match")
			// since end time might have the current value, we need to truncate the end times up to minutes when comparing
			// time.Now().UTC() is called at slightly different moments, resulting in tiny differences in the timestamp
			require.Equal(t, tt.expectedEndTime.Truncate(time.Minute), actualEnd.Truncate(time.Minute), "Expected end time did not match")
			require.Equal(t, tt.expectAlignedTimes, actualAligned, "Expected alignment flag did not match")
		})
	}
}
