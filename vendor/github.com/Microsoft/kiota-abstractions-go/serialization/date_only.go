package serialization

import (
	"strings"
	"time"
)

// DateOnly is a struct that represents a date only from a date time (Time).
type DateOnly struct {
	time time.Time
}

const dateOnlyFormat = "2006-01-02"

// String returns the date only as a string following the RFC3339 standard.
func (t DateOnly) String() string {
	return t.time.Format(dateOnlyFormat)
}

// ParseDateOnly parses a string into a DateOnly following the RFC3339 standard.
func ParseDateOnly(s string) (*DateOnly, error) {
	if len(strings.TrimSpace(s)) <= 0 {
		return nil, nil
	}
	timeValue, err := time.Parse(dateOnlyFormat, s)
	if err != nil {
		return nil, err
	}
	return &DateOnly{
		time: timeValue,
	}, nil
}

// NewDateOnly creates a new DateOnly from a time.Time.
func NewDateOnly(t time.Time) *DateOnly {
	return &DateOnly{
		time: t,
	}
}
