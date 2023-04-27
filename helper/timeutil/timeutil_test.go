// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timeutil

import (
	"reflect"
	"testing"
	"time"
)

func TestTimeutil_StartOfPreviousMonth(t *testing.T) {
	testCases := []struct {
		Input    time.Time
		Expected time.Time
	}{
		{
			Input:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Expected: time.Date(2019, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Input:    time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
			Expected: time.Date(2019, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Input:    time.Date(2020, 3, 31, 23, 59, 59, 999999999, time.UTC),
			Expected: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		result := StartOfPreviousMonth(tc.Input)
		if !result.Equal(tc.Expected) {
			t.Errorf("start of month before %v is %v, got %v", tc.Input, tc.Expected, result)
		}
	}
}

func TestTimeutil_StartOfMonth(t *testing.T) {
	testCases := []struct {
		Input    time.Time
		Expected time.Time
	}{
		{
			Input:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			Expected: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Input:    time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
			Expected: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Input:    time.Date(2020, 1, 1, 0, 0, 0, 1, time.UTC),
			Expected: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Input:    time.Date(2020, 1, 31, 23, 59, 59, 999999999, time.UTC),
			Expected: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Input:    time.Date(2020, 2, 28, 1, 2, 3, 4, time.UTC),
			Expected: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		result := StartOfMonth(tc.Input)
		if !result.Equal(tc.Expected) {
			t.Errorf("start of %v is %v, expected %v", tc.Input, result, tc.Expected)
		}
	}
}

func TestTimeutil_IsMonthStart(t *testing.T) {
	testCases := []struct {
		input    time.Time
		expected bool
	}{
		{
			input:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			input:    time.Date(2020, 1, 1, 0, 0, 0, 1, time.UTC),
			expected: false,
		},
		{
			input:    time.Date(2020, 4, 5, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		{
			input:    time.Date(2020, 1, 31, 23, 59, 59, 999999999, time.UTC),
			expected: false,
		},
	}

	for _, tc := range testCases {
		result := IsMonthStart(tc.input)
		if result != tc.expected {
			t.Errorf("is %v the start of the month? expected %t, got %t", tc.input, tc.expected, result)
		}
	}
}

func TestTimeutil_EndOfMonth(t *testing.T) {
	testCases := []struct {
		Input    time.Time
		Expected time.Time
	}{
		{
			// The current behavior does not use the nanoseconds
			// because we didn't want to clutter the result of end-of-month reporting.
			Input:    time.Date(2020, 1, 31, 23, 59, 59, 0, time.UTC),
			Expected: time.Date(2020, 1, 31, 23, 59, 59, 0, time.UTC),
		},
		{
			Input:    time.Date(2020, 1, 31, 23, 59, 59, 999999999, time.UTC),
			Expected: time.Date(2020, 1, 31, 23, 59, 59, 0, time.UTC),
		},
		{
			Input:    time.Date(2020, 1, 15, 1, 2, 3, 4, time.UTC),
			Expected: time.Date(2020, 1, 31, 23, 59, 59, 0, time.UTC),
		},
		{
			// Leap year
			Input:    time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
			Expected: time.Date(2020, 2, 29, 23, 59, 59, 0, time.UTC),
		},
		{
			// non-leap year
			Input:    time.Date(2100, 2, 1, 0, 0, 0, 0, time.UTC),
			Expected: time.Date(2100, 2, 28, 23, 59, 59, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		result := EndOfMonth(tc.Input)
		if !result.Equal(tc.Expected) {
			t.Errorf("end of %v is %v, expected %v", tc.Input, result, tc.Expected)
		}
	}
}

func TestTimeutil_IsPreviousMonth(t *testing.T) {
	testCases := []struct {
		tInput       time.Time
		compareInput time.Time
		expected     bool
	}{
		{
			tInput:       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			compareInput: time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC),
			expected:     false,
		},
		{
			tInput:       time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC),
			compareInput: time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC),
			expected:     true,
		},
		{
			// leap year (false)
			tInput:       time.Date(2019, 12, 29, 10, 10, 10, 0, time.UTC),
			compareInput: time.Date(2020, 2, 29, 10, 10, 10, 0, time.UTC),
			expected:     false,
		},
		{
			// leap year (true)
			tInput:       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			compareInput: time.Date(2020, 2, 29, 10, 10, 10, 0, time.UTC),
			expected:     true,
		},
		{
			tInput:       time.Date(2018, 5, 5, 5, 0, 0, 0, time.UTC),
			compareInput: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected:     false,
		},
		{
			// test normalization. want to make subtracting 1 month from 3/30/2020 doesn't yield 2/30/2020, normalized
			// to 3/1/2020
			tInput:       time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
			compareInput: time.Date(2020, 3, 30, 0, 0, 0, 0, time.UTC),
			expected:     true,
		},
	}

	for _, tc := range testCases {
		result := IsPreviousMonth(tc.tInput, tc.compareInput)
		if result != tc.expected {
			t.Errorf("%v in previous month to %v? expected %t, got %t", tc.tInput, tc.compareInput, tc.expected, result)
		}
	}
}

func TestTimeutil_IsCurrentMonth(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		input    time.Time
		expected bool
	}{
		{
			input:    now,
			expected: true,
		},
		{
			input:    StartOfMonth(now).AddDate(0, 0, -1),
			expected: false,
		},
		{
			input:    EndOfMonth(now).AddDate(0, 0, -1),
			expected: true,
		},
		{
			input:    StartOfMonth(now).AddDate(-1, 0, 0),
			expected: false,
		},
	}

	for _, tc := range testCases {
		result := IsCurrentMonth(tc.input, now)
		if result != tc.expected {
			t.Errorf("invalid result. expected %t for %v", tc.expected, tc.input)
		}
	}
}

func TestTimeUtil_ContiguousMonths(t *testing.T) {
	testCases := []struct {
		input    []time.Time
		expected []time.Time
	}{
		{
			input: []time.Time{
				time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 5, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 5, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input: []time.Time{
				time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input: []time.Time{
				time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			input:    []time.Time{},
			expected: []time.Time{},
		},
		{
			input:    nil,
			expected: nil,
		},
		{
			input: []time.Time{
				time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			expected: []time.Time{
				time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tc := range testCases {
		result := GetMostRecentContiguousMonths(tc.input)

		if !reflect.DeepEqual(tc.expected, result) {
			t.Errorf("invalid contiguous segment returned. expected %v, got %v", tc.expected, result)
		}
	}
}

func TestTimeUtil_ParseTimeFromPath(t *testing.T) {
	testCases := []struct {
		input       string
		expectedOut time.Time
		expectError bool
	}{
		{
			input:       "719020800/1",
			expectedOut: time.Unix(719020800, 0).UTC(),
			expectError: false,
		},
		{
			input:       "1601415205/3",
			expectedOut: time.Unix(1601415205, 0).UTC(),
			expectError: false,
		},
		{
			input:       "baddata/3",
			expectedOut: time.Time{},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		result, err := ParseTimeFromPath(tc.input)
		gotError := err != nil

		if result != tc.expectedOut {
			t.Errorf("bad timestamp on input %q. expected: %v got: %v", tc.input, tc.expectedOut, result)
		}
		if gotError != tc.expectError {
			t.Errorf("bad error status on input %q. expected error: %t, got error: %t", tc.input, tc.expectError, gotError)
		}
	}
}
