// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/stretchr/testify/require"
)

// TestAlignToBillingPeriodStart tests alignToBillingPeriodStart for correct alignment of the start time and end time
func TestAlignToBillingPeriodStart(t *testing.T) {
	// assume the billing start time is February 1, 2023
	billingStartTime := time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		givenTime time.Time
		expected  time.Time
		byEndTime bool
	}{
		"same as billing start time: given start time is February 1, 2023": {
			givenTime: time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC),
			expected:  billingStartTime,
		},
		"same as billing start time: given end time is February 1, 2023": {
			givenTime: time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC),
			expected:  billingStartTime.AddDate(-1, 0, 0),
			byEndTime: true,
		},
		"within the same billing period: given start time is February 3, 2023": {
			givenTime: time.Date(2023, time.June, 3, 0, 0, 0, 0, time.UTC),
			expected:  billingStartTime,
		},
		"within the same billing period: given end time is February 3, 2023": {
			givenTime: time.Date(2023, time.June, 3, 0, 0, 0, 0, time.UTC),
			expected:  billingStartTime,
			byEndTime: true,
		},
		"before the current billing period: given start time is December 20, 2022": {
			givenTime: time.Date(2022, time.December, 20, 0, 0, 0, 0, time.UTC),
			expected:  billingStartTime.AddDate(-1, 0, 0),
		},
		"before the current billing period: given end time is December 20, 2022": {
			givenTime: time.Date(2022, time.December, 20, 0, 0, 0, 0, time.UTC),
			expected:  billingStartTime.AddDate(-1, 0, 0),
			byEndTime: true,
		},
		"exactly one year before the billing start time: given start time is February 1, 2022": {
			givenTime: time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC),
			expected:  billingStartTime.AddDate(-1, 0, 0),
		},
		"exactly one year before the billing start time: given end time is February 1, 2022": {
			givenTime: time.Date(2022, time.February, 1, 0, 0, 0, 0, time.UTC),
			expected:  billingStartTime.AddDate(-2, 0, 0),
			byEndTime: true,
		},
		"several years in the past: given start time is March 10, 2015": {
			givenTime: time.Date(2015, time.March, 10, 0, 0, 0, 0, time.UTC),
			expected:  time.Date(2015, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
		"several years in the past: given end time is March 10, 2015": {
			givenTime: time.Date(2015, time.March, 10, 0, 0, 0, 0, time.UTC),
			expected:  time.Date(2015, time.February, 1, 0, 0, 0, 0, time.UTC),
			byEndTime: true,
		},
		"several years in the future: given start time is April 15, 2029": {
			givenTime: time.Date(2029, time.April, 15, 0, 0, 0, 0, time.UTC),
			expected:  billingStartTime,
		},
		"several years in the future: given end time is April 15, 2029": {
			givenTime: time.Date(2029, time.April, 15, 0, 0, 0, 0, time.UTC),
			expected:  billingStartTime,
			byEndTime: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := alignToBillingPeriodStart(billingStartTime, tt.givenTime, tt.byEndTime)
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
		errorMsg           string
	}{
		"No start or end time provided, should return times from current billing period start time to now": {
			givenStartTime:     time.Time{},
			givenEndTime:       time.Time{},
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: false,
		},
		"Start time provided to be 5 months after the current billing period start time within the current billing cycle": {
			givenStartTime: func() time.Time {
				fiveMonthsForward := currentBillingStartDate.AddDate(0, 5, 0)
				return time.Date(fiveMonthsForward.Year(), fiveMonthsForward.Month(), 23, 5, 0, 0, 0, time.UTC)
			}(),
			givenEndTime:       time.Time{},
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: true,
		},
		"Start time provided to be 7 months after the current time outside the current billing cycle": {
			givenStartTime: func() time.Time {
				sevenMonthsForward := currentTime.AddDate(0, 7, 0)
				return time.Date(sevenMonthsForward.Year(), sevenMonthsForward.Month(), 1, 2, 0, 0, 0, time.UTC)
			}(),
			givenEndTime:       time.Time{},
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: true,
		},
		"Start time provided to be 4 months after the current billing cycle end time outside the current billing cycle": {
			givenStartTime: func() time.Time {
				fourMonthsForward := currentBillingEndDate.AddDate(0, 4, 0)
				return time.Date(fourMonthsForward.Year(), fourMonthsForward.Month(), 3, 13, 0, 0, 0, time.UTC)
			}(),
			givenEndTime:       time.Time{},
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: true,
		},
		"Start time provided to be 2 years before the current time outside the current billing cycle": {
			givenStartTime: func() time.Time {
				twoYearsAgo := currentTime.AddDate(-2, 0, 0)
				return time.Date(twoYearsAgo.Year(), twoYearsAgo.Month(), 13, 21, 0, 0, 0, time.UTC)
			}(),
			givenEndTime:       time.Time{},
			expectedStartTime:  currentBillingStartDate.AddDate(-2, 0, 0),
			expectedEndTime:    currentBillingEndDate.AddDate(-2, 0, 0),
			expectAlignedTimes: true,
		},
		"Start time provided to be at exact current billing period start": {
			givenStartTime:     currentBillingStartDate,
			givenEndTime:       time.Time{},
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: false,
		},
		"Start time provided to be 2 years and 2 months before current billing period start": {
			givenStartTime:     currentBillingStartDate.AddDate(-2, -2, 0),
			givenEndTime:       time.Time{},
			expectedStartTime:  currentBillingStartDate.AddDate(-3, 0, 0),
			expectedEndTime:    currentBillingStartDate.AddDate(-2, 0, 0),
			expectAlignedTimes: true,
		},
		"Start time provided to be 2 years before current billing period start": {
			givenStartTime:     currentBillingStartDate.AddDate(-2, 0, 0),
			givenEndTime:       time.Time{},
			expectedStartTime:  currentBillingStartDate.AddDate(-2, 0, 0),
			expectedEndTime:    currentBillingStartDate.AddDate(-1, 0, 0),
			expectAlignedTimes: false,
		},
		"Start time provided to be 2 years and 2 months after the current billing period start": {
			givenStartTime:     currentBillingStartDate.AddDate(2, 2, 20),
			givenEndTime:       time.Time{},
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
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
		"End time provided to be 7 months after the current time outside the current billing cycle": {
			givenStartTime: time.Time{},
			givenEndTime: func() time.Time {
				sevenMonthsForward := currentTime.AddDate(0, 7, 0)
				return time.Date(sevenMonthsForward.Year(), sevenMonthsForward.Month(), 1, 2, 0, 0, 0, time.UTC)
			}(),
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: true,
		},
		"End time provided to be 4 months after the current billing cycle end time outside the current billing cycle": {
			givenStartTime: time.Time{},
			givenEndTime: func() time.Time {
				fourMonthsForward := currentBillingEndDate.AddDate(0, 4, 0)
				return time.Date(fourMonthsForward.Year(), fourMonthsForward.Month(), 3, 13, 0, 0, 0, time.UTC)
			}(),
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: true,
		},
		"End time provided to be 2 years before the current time outside the current billing cycle": {
			givenStartTime: time.Time{},
			givenEndTime: func() time.Time {
				twoYearsAgo := currentTime.AddDate(-2, 0, 0)
				return time.Date(twoYearsAgo.Year(), twoYearsAgo.Month(), 13, 21, 0, 0, 0, time.UTC)
			}(),
			expectedStartTime:  currentBillingStartDate.AddDate(-2, 0, 0),
			expectedEndTime:    currentBillingEndDate.AddDate(-2, 0, 0),
			expectAlignedTimes: true,
		},
		"End time provided to be at exact current billing period start": {
			givenStartTime:     time.Time{},
			givenEndTime:       currentBillingStartDate,
			expectedStartTime:  currentBillingStartDate.AddDate(-1, 0, 0),
			expectedEndTime:    currentBillingStartDate,
			expectAlignedTimes: false,
		},
		"End time provided to be 2 years and 2 months before current billing period start": {
			givenStartTime:     time.Time{},
			givenEndTime:       currentBillingStartDate.AddDate(-2, -2, 0),
			expectedStartTime:  currentBillingStartDate.AddDate(-3, 0, 0),
			expectedEndTime:    currentBillingStartDate.AddDate(-2, 0, 0),
			expectAlignedTimes: true,
		},
		"End time provided to be 2 years before current billing period start": {
			givenStartTime:     time.Time{},
			givenEndTime:       currentBillingStartDate.AddDate(-2, 0, 0),
			expectedStartTime:  currentBillingStartDate.AddDate(-3, 0, 0),
			expectedEndTime:    currentBillingStartDate.AddDate(-2, 0, 0),
			expectAlignedTimes: false,
		},
		"End time provided to be 2 years and 2 months after the current billing period start": {
			givenStartTime:     time.Time{},
			givenEndTime:       currentBillingStartDate.AddDate(2, 2, 20),
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: true,
		},
		"Start time is before the current billing cycle start date and end time is within the current billing cycle": {
			givenStartTime: func() time.Time {
				threeMonthsAgo := currentBillingStartDate.AddDate(0, -2, -3)
				return time.Date(threeMonthsAgo.Year(), threeMonthsAgo.Month(), 12, 11, 0, 0, 0, time.UTC)
			}(),
			givenEndTime: func() time.Time {
				oneMonthForward := currentBillingStartDate.AddDate(0, 1, 6)
				return time.Date(oneMonthForward.Year(), oneMonthForward.Month(), 13, 24, 0, 0, 0, time.UTC)
			}(),
			expectedStartTime:  currentBillingStartDate.AddDate(-1, 0, 0),
			expectedEndTime:    currentBillingEndDate.AddDate(-1, 0, 0),
			expectAlignedTimes: true,
		},
		"Start time is within the current billing cycle and end time is after the current time": {
			givenStartTime: func() time.Time {
				twoDaysForward := currentBillingStartDate.AddDate(0, 0, 2)
				return time.Date(twoDaysForward.Year(), twoDaysForward.Month(), 1, 14, 0, 0, 0, time.UTC)
			}(),
			givenEndTime: func() time.Time {
				oneMonthForward := currentTime.AddDate(0, 1, 10)
				return time.Date(oneMonthForward.Year(), oneMonthForward.Month(), 13, 24, 0, 0, 0, time.UTC)
			}(),
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: true,
		},
		"Both start time and end time are after the current time": {
			givenStartTime: func() time.Time {
				twoDaysForward := currentTime.AddDate(1, 2, 2)
				return time.Date(twoDaysForward.Year(), twoDaysForward.Month(), 3, 12, 0, 0, 0, time.UTC)
			}(),
			givenEndTime: func() time.Time {
				oneMonthForward := currentTime.AddDate(1, 2, 10)
				return time.Date(oneMonthForward.Year(), oneMonthForward.Month(), 13, 24, 0, 0, 0, time.UTC)
			}(),
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: true,
		},
		"Start and end time in the same billing period from 1 year ago": {
			givenStartTime:     currentBillingStartDate.AddDate(-1, 0, 0),
			givenEndTime:       currentBillingStartDate,
			expectedStartTime:  currentBillingStartDate.AddDate(-1, 0, 0),
			expectedEndTime:    currentBillingStartDate,
			expectAlignedTimes: false,
		},
		"Start time at exact current billing start, end time ignored and capped to today": {
			givenStartTime:     currentBillingStartDate,
			givenEndTime:       currentBillingEndDate,
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: true,
		},
		"End time in current billing period, should cap at today": {
			givenStartTime:     time.Time{},
			givenEndTime:       currentBillingStartDate.AddDate(0, 0, 23),
			expectedStartTime:  currentBillingStartDate,
			expectedEndTime:    currentTime,
			expectAlignedTimes: true,
		},
		"Start time after end time, expect error": {
			givenStartTime: currentBillingStartDate.AddDate(0, 2, 0),
			givenEndTime:   currentBillingStartDate.AddDate(0, 1, 0),
			errorMsg:       "start_time is later than end_time",
		},
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
			if tt.errorMsg != "" {
				require.Error(t, err)
				require.Equal(t, tt.errorMsg, err.Error())
			} else {
				require.NoError(t, err)
				currentTime = time.Now().UTC()

				require.Equal(t, tt.expectedStartTime, actualStart, "Expected start time did not match")
				// since end time might have the current value, we need to truncate the end times up to minutes when comparing
				// time.Now().UTC() is called at slightly different moments, resulting in tiny differences in the timestamp
				require.Equal(t, tt.expectedEndTime.Truncate(time.Minute), actualEnd.Truncate(time.Minute), "Expected end time did not match")
				require.Equal(t, tt.expectAlignedTimes, actualAligned, "Expected alignment flag did not match")
			}
		})
	}
}
