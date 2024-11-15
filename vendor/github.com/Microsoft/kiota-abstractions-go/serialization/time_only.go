package serialization

import (
	"errors"
	"strings"
	"time"
)

// TimeOnly is represents the time part of a date time (time) value.
type TimeOnly struct {
	time time.Time
}

const timeOnlyFormat = "15:04:05.000000000"

var timeOnlyParsingFormats = map[int]string{
	0: "15:04:05", //Go doesn't seem to support optional parameters in time.Parse, which is sad
	1: "15:04:05.0",
	2: "15:04:05.00",
	3: "15:04:05.000",
	4: "15:04:05.0000",
	5: "15:04:05.00000",
	6: "15:04:05.000000",
	7: "15:04:05.0000000",
	8: "15:04:05.00000000",
	9: timeOnlyFormat,
}

// String returns the time only as a string following the RFC3339 standard.
func (t TimeOnly) String() string {
	return t.time.Format(timeOnlyFormat)
}

// ParseTimeOnly parses a string into a TimeOnly following the RFC3339 standard.
func ParseTimeOnly(s string) (*TimeOnly, error) {
	if len(strings.TrimSpace(s)) <= 0 {
		return nil, nil
	}
	splat := strings.Split(s, ".")
	parsingFormat := timeOnlyParsingFormats[0]
	if len(splat) > 1 {
		dotSectionLen := len(splat[1])
		if dotSectionLen >= len(timeOnlyParsingFormats) {
			return nil, errors.New("too many decimal places in time only string")
		}
		parsingFormat = timeOnlyParsingFormats[dotSectionLen]
	}
	timeValue, err := time.Parse(parsingFormat, s)
	if err != nil {
		return nil, err
	}
	return &TimeOnly{
		time: timeValue,
	}, nil
}

// NewTimeOnly creates a new TimeOnly from a time.Time.
func NewTimeOnly(t time.Time) *TimeOnly {
	return &TimeOnly{
		time: t,
	}
}
