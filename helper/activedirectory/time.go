package activedirectory

import (
	"strconv"
	"time"
)

// ParseTime parses dates represented as Active Directory LargeInts into times.
// Not all time fields are represented this way,
// so be sure to test that your particular time returns expected results.
// Some time fields represented as LargeInts include accountExpires, lastLogon, lastLogonTimestamp, and pwdLastSet.
// More: https://social.technet.microsoft.com/wiki/contents/articles/31135.active-directory-large-integer-attributes.aspx
func ParseTime(ticks string) (time.Time, error) {
	i, err := strconv.ParseInt(ticks, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return ToTime(i), nil
}

// ToTime converts an ActiveDirectory time in ticks to a time.
// This algorithm is summarized as:
//
// 		Many dates are saved in Active Directory as Large Integer values.
// 		These attributes represent dates as the number of 100-nanosecond intervals since 12:00 AM January 1, 1601.
//		100-nanosecond intervals, equal to 0.0000001 seconds, are also called ticks.
//		Dates in Active Directory are always saved in Coordinated Universal Time, or UTC.
//		More: https://social.technet.microsoft.com/wiki/contents/articles/31135.active-directory-large-integer-attributes.aspx
//
// If we directly follow the above algorithm we encounter time.Duration limits of 290 years and int overflow issues.
// Thus below, we carefully sidestep those.
func ToTime(ticks int64) time.Time {

	// Go durations are limited to 290 years.
	// So, let's subtract ticks since 12:00 AM January 1, 1901.
	// Then we can use that as our starting point.
	ticksSince1901 := ticks - int64(94670208000000000)

	// These are ticks, so let's convert them to nanoseconds.
	since1901nanos := time.Duration(ticksSince1901 * 100)

	since := time.Date(1901, time.January, 1, 0, 0, 0, 0, time.UTC)
	return since.Add(since1901nanos)
}
