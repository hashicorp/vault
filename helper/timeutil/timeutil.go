// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package timeutil

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func StartOfPreviousMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location()).AddDate(0, -1, 0)
}

func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// IsCurrentDay checks if :t: is in the current day, as defined by :compare:
// generally, pass in time.Now().UTC() as :compare:
func IsCurrentDay(t, compare time.Time) bool {
	thisDayStart := StartOfDay(compare)
	queryDayStart := StartOfDay(t)

	return queryDayStart.Equal(thisDayStart)
}

func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

func StartOfNextMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location()).AddDate(0, 1, 0)
}

// IsMonthStart checks if :t: is the start of the month
func IsMonthStart(t time.Time) bool {
	return t.Equal(StartOfMonth(t))
}

func EndOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	if month == time.December {
		return time.Date(year, time.December, 31, 23, 59, 59, 0, t.Location())
	} else {
		eom := time.Date(year, month+1, 1, 23, 59, 59, 0, t.Location())
		return eom.AddDate(0, 0, -1)
	}
}

// IsPreviousMonth checks if :t: is in the month directly before :toCompare:
func IsPreviousMonth(t, toCompare time.Time) bool {
	thisMonthStart := StartOfMonth(toCompare)
	previousMonthStart := StartOfMonth(thisMonthStart.AddDate(0, 0, -1))

	if t.Equal(previousMonthStart) {
		return true
	}
	return t.After(previousMonthStart) && t.Before(thisMonthStart)
}

// IsCurrentMonth checks if :t: is in the current month, as defined by :compare:
// generally, pass in time.Now().UTC() as :compare:
func IsCurrentMonth(t, compare time.Time) bool {
	thisMonthStart := StartOfMonth(compare)
	queryMonthStart := StartOfMonth(t)

	return queryMonthStart.Equal(thisMonthStart)
}

// GetMostRecentContiguousMonths finds the start time of the most
// recent set of continuous months.
//
// For example, if the most recent start time is Aug 15, then that range is just 1 month
// If the recent start times are Aug 1 and July 1 and June 15, then that range is
// three months and we return June 15.
//
// note: return slice will be nil if :startTimes: is nil
// :startTimes: must be sorted in decreasing order (see unit test for examples)
func GetMostRecentContiguousMonths(startTimes []time.Time) []time.Time {
	if len(startTimes) < 2 {
		// no processing needed if 0 or 1 months worth of logs
		return startTimes
	}

	out := []time.Time{startTimes[0]}
	if !IsMonthStart(out[0]) {
		// there is less than one contiguous month (most recent start time is after the start of this month)
		return out
	}

	i := 1
	for ; i < len(startTimes); i++ {
		if !IsMonthStart(startTimes[i]) || !IsPreviousMonth(startTimes[i], startTimes[i-1]) {
			break
		}

		out = append(out, startTimes[i])
	}

	// handle mid-month log starts
	if i < len(startTimes) {
		if IsPreviousMonth(StartOfMonth(startTimes[i]), startTimes[i-1]) {
			// the earliest part of the segment is mid-month, but still valid for this segment
			out = append(out, startTimes[i])
		}
	}

	return out
}

func InRange(t, start, end time.Time) bool {
	return (t.Equal(start) || t.After(start)) &&
		(t.Equal(end) || t.Before(end))
}

// ParseTimeFromPath returns a UTC time from a path of the form '<timestamp>/',
// where <timestamp> is a Unix timestamp
func ParseTimeFromPath(path string) (time.Time, error) {
	elems := strings.Split(path, "/")
	if len(elems) == 1 {
		// :path: is a directory that must have children
		return time.Time{}, errors.New("Invalid path provided")
	}

	unixSeconds, err := strconv.ParseInt(elems[0], 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("could not convert time from path segment %q. error: %w", elems[0], err)
	}

	return time.Unix(unixSeconds, 0).UTC(), nil
}

// Compute the N-month period before the given date.
// For example, if it is currently April 2020, then 12 months is April 2019 through March 2020.
func MonthsPreviousTo(months int, now time.Time) time.Time {
	firstOfMonth := StartOfMonth(now.UTC())
	return firstOfMonth.AddDate(0, -months, 0)
}

// Skip this test if too close to the end of a month!
func SkipAtEndOfMonth(t *testing.T) {
	t.Helper()

	thisMonth := StartOfMonth(time.Now().UTC())
	endOfMonth := EndOfMonth(thisMonth)
	if endOfMonth.Sub(time.Now()) < 10*time.Minute {
		t.Skip("too close to end of month")
	}
}

// This interface allows unit tests to substitute in a simulated Clock.
type Clock interface {
	Now() time.Time
	NewTicker(time.Duration) *time.Ticker
	NewTimer(time.Duration) *time.Timer
}

type DefaultClock struct{}

var _ Clock = (*DefaultClock)(nil)

func (_ DefaultClock) Now() time.Time {
	return time.Now()
}

func (_ DefaultClock) NewTicker(d time.Duration) *time.Ticker {
	return time.NewTicker(d)
}

func (_ DefaultClock) NewTimer(d time.Duration) *time.Timer {
	return time.NewTimer(d)
}

// NormalizeToYear returns date normalized to the latest date
// within one year of normal. Assumes the date argument is
// some date before normal.
func NormalizeToYear(date, normal time.Time) time.Time {
	for date.AddDate(1, 0, 0).Compare(normal) <= 0 {
		date = date.AddDate(1, 0, 0)
	}
	return date
}
