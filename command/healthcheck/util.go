package healthcheck

import (
	"fmt"
	"time"
)

var (
	oneDay   = 24 * time.Hour
	oneWeek  = 7 * oneDay
	oneMonth = 30 * oneDay
	oneYear  = 365 * oneDay
)

var suffixDurationMap = map[string]time.Duration{
	"y":  oneYear,
	"mo": oneMonth,
	"w":  oneWeek,
	"d":  oneDay,
}
var orderedSuffixes = []string{"y", "mo", "w", "d"}

func FormatDuration(d time.Duration) string {
	var result string
	for _, suffix := range orderedSuffixes {
		unit := suffixDurationMap[suffix]
		if d > unit {
			quantity := int64(d / unit)
			result = fmt.Sprintf("%v%v%v", quantity, suffix, result)
			d = d - (time.Duration(quantity) * unit)
		}
	}

	if d > 0 {
		result = d.String() + result
	}

	return result
}
