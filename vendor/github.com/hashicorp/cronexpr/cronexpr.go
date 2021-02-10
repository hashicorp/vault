/*!
 * Copyright 2013 Raymond Hill
 *
 * Modifications 2020 - HashiCorp
 *
 * Project: github.com/gorhill/cronexpr
 * File: cronexpr.go
 * Version: 1.0
 * License: pick the one which suits you :
 *   GPL v3 see <https://www.gnu.org/licenses/gpl.html>
 *   APL v2 see <http://www.apache.org/licenses/LICENSE-2.0>
 *
 */

// Package cronexpr parses cron time expressions.
package cronexpr

/******************************************************************************/

import (
	"fmt"
	"sort"
	"time"
)

/******************************************************************************/

// A Expression represents a specific cron time expression as defined at
// <https://github.com/gorhill/cronexpr#implementation>
type Expression struct {
	expression             string
	secondList             []int
	minuteList             []int
	hourList               []int
	daysOfMonth            map[int]bool
	workdaysOfMonth        map[int]bool
	lastDayOfMonth         bool
	lastWorkdayOfMonth     bool
	daysOfMonthRestricted  bool
	actualDaysOfMonthList  []int
	monthList              []int
	daysOfWeek             map[int]bool
	specificWeekDaysOfWeek map[int]bool
	lastWeekDaysOfWeek     map[int]bool
	daysOfWeekRestricted   bool
	yearList               []int
}

/******************************************************************************/

// MustParse returns a new Expression pointer. It expects a well-formed cron
// expression. If a malformed cron expression is supplied, it will `panic`.
// See <https://github.com/gorhill/cronexpr#implementation> for documentation
// about what is a well-formed cron expression from this library's point of
// view.
func MustParse(cronLine string) *Expression {
	expr, err := Parse(cronLine)
	if err != nil {
		panic(err)
	}
	return expr
}

/******************************************************************************/

// Parse returns a new Expression pointer. An error is returned if a malformed
// cron expression is supplied.
// See <https://github.com/gorhill/cronexpr#implementation> for documentation
// about what is a well-formed cron expression from this library's point of
// view.
func Parse(cronLine string) (*Expression, error) {

	// Maybe one of the built-in aliases is being used
	cron := cronNormalizer.Replace(cronLine)

	indices := fieldFinder.FindAllStringIndex(cron, -1)
	fieldCount := len(indices)
	if fieldCount < 5 {
		return nil, fmt.Errorf("missing field(s)")
	}
	// ignore fields beyond 7th
	if fieldCount > 7 {
		fieldCount = 7
	}

	var expr = Expression{}
	var field = 0
	var err error

	// second field (optional)
	if fieldCount == 7 {
		err = expr.secondFieldHandler(cron[indices[field][0]:indices[field][1]])
		if err != nil {
			return nil, err
		}
		field += 1
	} else {
		expr.secondList = []int{0}
	}

	// minute field
	err = expr.minuteFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field += 1

	// hour field
	err = expr.hourFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field += 1

	// day of month field
	err = expr.domFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field += 1

	// month field
	err = expr.monthFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field += 1

	// day of week field
	err = expr.dowFieldHandler(cron[indices[field][0]:indices[field][1]])
	if err != nil {
		return nil, err
	}
	field += 1

	// year field
	if field < fieldCount {
		err = expr.yearFieldHandler(cron[indices[field][0]:indices[field][1]])
		if err != nil {
			return nil, err
		}
	} else {
		expr.yearList = yearDescriptor.defaultList
	}

	return &expr, nil
}

/******************************************************************************/

// Next returns the closest time instant immediately following `fromTime` which
// matches the cron expression `expr`.
//
// The `time.Location` of the returned time instant is the same as that of
// `fromTime`.
//
// The zero value of time.Time is returned if no matching time instant exists
// or if a `fromTime` is itself a zero value.
func (expr *Expression) Next(fromTime time.Time) time.Time {
	// Special case
	if fromTime.IsZero() {
		return fromTime
	}

	loc := fromTime.Location()
	t := fromTime.Add(time.Second - time.Duration(fromTime.Nanosecond())*time.Nanosecond)

WRAP:

	// let's find the next date that satisfies condition
	v := t.Year()
	if i := sort.SearchInts(expr.yearList, v); i == len(expr.yearList) {
		return time.Time{}
	} else if v != expr.yearList[i] {
		t = time.Date(expr.yearList[i], time.Month(expr.monthList[0]), 1, 0, 0, 0, 0, loc)
	}

	v = int(t.Month())
	if i := sort.SearchInts(expr.monthList, v); i == len(expr.monthList) {
		// try again with a new year
		t = time.Date(t.Year()+1, time.Month(expr.monthList[0]), 1, 0, 0, 0, 0, loc)
		goto WRAP
	} else if v != expr.monthList[i] {
		t = time.Date(t.Year(), time.Month(expr.monthList[i]), 1, 0, 0, 0, 0, loc)
	}

	expr.actualDaysOfMonthList = expr.calculateActualDaysOfMonth(t.Year(), int(t.Month()))
	if len(expr.actualDaysOfMonthList) == 0 {
		t = time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, loc)
		goto WRAP
	}

	v = t.Day()
	if i := sort.SearchInts(expr.actualDaysOfMonthList, v); i == len(expr.actualDaysOfMonthList) {
		t = time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, loc)
		goto WRAP
	} else if v != expr.actualDaysOfMonthList[i] {
		t = time.Date(t.Year(), t.Month(), expr.actualDaysOfMonthList[i], 0, 0, 0, 0, loc)

		// in San Palo, before 2019, there may be no midnight (or multiple midnights)
		// due to DST
		if t.Hour() != 0 {
			if t.Hour() > 12 {
				t = t.Add(time.Duration(24-t.Hour()) * time.Hour)
			} else {
				t = t.Add(time.Duration(-t.Hour()) * time.Hour)
			}
		}
	}

	if timeZoneInDay(t) {
		goto SLOW_CLOCK
	}

	// Fast path where hours/minutes behave as expected trivially
	v = t.Hour()
	if i := sort.SearchInts(expr.hourList, v); i == len(expr.hourList) {
		t = time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, loc)
		goto WRAP
	} else if v != expr.hourList[i] {
		t = time.Date(t.Year(), t.Month(), t.Day(), expr.hourList[i], expr.minuteList[0], expr.secondList[0], 0, loc)
	}

	v = t.Minute()
	if i := sort.SearchInts(expr.minuteList, v); i == len(expr.minuteList) {
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour()+1, 0, 0, 0, loc)
		goto WRAP
	} else if v != expr.minuteList[i] {
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), expr.minuteList[i], expr.secondList[0], 0, loc)
	}

	v = t.Second()
	if i := sort.SearchInts(expr.secondList, v); i == len(expr.secondList) {
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()+1, 0, 0, loc)
		goto WRAP
	} else if v != expr.secondList[i] {
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), expr.secondList[i], 0, loc)
	}

	return t

SLOW_CLOCK:
	// daylight saving effect is here, where odd things happen:
	// An hour may have 60 minutes, 30 minutes or 90 minutes;
	// partial hours may "repeat"!
	for !sortContains(expr.hourList, t.Hour()) {
		hourBefore := t.Hour()
		t = t.Add(time.Hour)
		if hourBefore == t.Hour() {
			t = t.Add(time.Hour)
		}
		t = t.Truncate(time.Minute)
		if t.Minute() != 0 {
			t = t.Add(-1 * time.Minute * time.Duration(t.Minute()))
		}

		if t.Hour() == 0 {
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
			goto WRAP
		}
	}

	for !sortContains(expr.minuteList, t.Minute()) {
		hoursBefore := t.Hour()
		t = t.Truncate(time.Minute).Add(time.Minute)
		if hoursBefore != t.Hour() {
			goto WRAP
		}
	}

	v = t.Second()
	t = t.Truncate(time.Minute)
	if i := sort.SearchInts(expr.secondList, v); i == len(expr.secondList) {
		t = t.Add(time.Minute)
		goto WRAP
	} else {
		t = t.Add(time.Duration(expr.secondList[i]) * time.Second)
	}

	return t
}

/******************************************************************************/

// NextN returns a slice of `n` closest time instants immediately following
// `fromTime` which match the cron expression `expr`.
//
// The time instants in the returned slice are in chronological ascending order.
// The `time.Location` of the returned time instants is the same as that of
// `fromTime`.
//
// A slice with len between [0-`n`] is returned, that is, if not enough existing
// matching time instants exist, the number of returned entries will be less
// than `n`.
func (expr *Expression) NextN(fromTime time.Time, n uint) []time.Time {
	nextTimes := make([]time.Time, 0, n)
	if n > 0 {
		fromTime = expr.Next(fromTime)
		for {
			if fromTime.IsZero() {
				break
			}
			nextTimes = append(nextTimes, fromTime)
			n -= 1
			if n == 0 {
				break
			}
			fromTime = expr.Next(fromTime)
		}
	}
	return nextTimes
}
