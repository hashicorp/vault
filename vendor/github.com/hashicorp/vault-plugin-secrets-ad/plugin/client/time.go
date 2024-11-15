// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"strconv"
	"time"
)

const (
	nanoSecondsPerSecond = 1000000000
	nanosInTick          = 100
	ticksPerSecond       = nanoSecondsPerSecond / nanosInTick
)

// ParseTicks parses dates represented as Active Directory LargeInts into times.
// Not all time fields are represented this way,
// so be sure to test that your particular time returns expected results.
// Some time fields represented as LargeInts include accountExpires, lastLogon, lastLogonTimestamp, and pwdLastSet.
// More: https://social.technet.microsoft.com/wiki/contents/articles/31135.active-directory-large-integer-attributes.aspx
func ParseTicks(ticks string) (time.Time, error) {
	i, err := strconv.ParseInt(ticks, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return TicksToTime(i), nil
}

// TicksToTime converts an ActiveDirectory time in ticks to a time.
// This algorithm is summarized as:
//
//	Many dates are saved in Active Directory as Large Integer values.
//	These attributes represent dates as the number of 100-nanosecond intervals since 12:00 AM January 1, 1601.
//	100-nanosecond intervals, equal to 0.0000001 seconds, are also called ticks.
//	Dates in Active Directory are always saved in Coordinated Universal Time, or UTC.
//	More: https://social.technet.microsoft.com/wiki/contents/articles/31135.active-directory-large-integer-attributes.aspx
//
// If we directly follow the above algorithm we encounter time.Duration limits of 290 years and int overflow issues.
// Thus below, we carefully sidestep those.
func TicksToTime(ticks int64) time.Time {
	origin := time.Date(1601, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
	secondsSinceOrigin := ticks / ticksPerSecond
	remainingNanoseconds := ticks % ticksPerSecond * 100
	return time.Unix(origin+secondsSinceOrigin, remainingNanoseconds).UTC()
}
