package serialization

import (
	"time"

	cjl "github.com/cjlapao/common-go/duration"
)

// ISODuration represents an ISO 8601 duration
type ISODuration struct {
	duration cjl.Duration
}

// GetYears returns the number of years.
func (i ISODuration) GetYears() int {
	return i.duration.Years
}

// GetWeeks returns the number of weeks.
func (i ISODuration) GetWeeks() int {
	return i.duration.Weeks
}

// GetDays returns the number of days.
func (i ISODuration) GetDays() int {
	return i.duration.Days
}

// GetHours returns the number of hours.
func (i ISODuration) GetHours() int {
	return i.duration.Hours
}

// GetMinutes returns the number of minutes.
func (i ISODuration) GetMinutes() int {
	return i.duration.Minutes
}

// GetSeconds returns the number of seconds.
func (i ISODuration) GetSeconds() int {
	return i.duration.Seconds
}

// GetMilliSeconds returns the number of milliseconds.
func (i ISODuration) GetMilliSeconds() int {
	return i.duration.MilliSeconds
}

// SetYears sets the number of years.
func (i ISODuration) SetYears(years int) {
	i.duration.Years = years
}

// SetWeeks sets the number of weeks.
func (i ISODuration) SetWeeks(weeks int) {
	i.duration.Weeks = weeks
}

// SetDays sets the number of days.
func (i ISODuration) SetDays(days int) {
	i.duration.Days = days
}

// SetHours sets the number of hours.
func (i ISODuration) SetHours(hours int) {
	i.duration.Hours = hours
}

// SetMinutes sets the number of minutes.
func (i ISODuration) SetMinutes(minutes int) {
	i.duration.Minutes = minutes
}

// SetSeconds sets the number of seconds.
func (i ISODuration) SetSeconds(seconds int) {
	i.duration.Seconds = seconds
}

// SetMilliSeconds sets the number of milliseconds.
func (i ISODuration) SetMilliSeconds(milliSeconds int) {
	i.duration.MilliSeconds = milliSeconds
}

// ParseISODuration parses a string into an ISODuration following the ISO 8601 standard.
func ParseISODuration(s string) (*ISODuration, error) {
	d, err := cjl.FromString(s)
	if err != nil {
		return nil, err
	}
	return &ISODuration{
		duration: *d,
	}, nil
}

// NewISODuration creates a new ISODuration from primitive values.
func NewDuration(years int, weeks int, days int, hours int, minutes int, seconds int, milliSeconds int) *ISODuration {
	return &ISODuration{
		duration: cjl.Duration{
			Years:        years,
			Weeks:        weeks,
			Days:         days,
			Hours:        hours,
			Minutes:      minutes,
			Seconds:      seconds,
			MilliSeconds: milliSeconds,
		},
	}
}

// String returns the ISO 8601 representation of the duration.
func (i ISODuration) String() string {
	return i.duration.String()
}

// FromDuration returns an ISODuration from a time.Duration.
func FromDuration(d time.Duration) *ISODuration {
	return NewDuration(0, 0, 0, 0, 0, 0, int(d.Truncate(time.Millisecond).Milliseconds()))
}

// ToDuration returns the time.Duration representation of the ISODuration.
func (d ISODuration) ToDuration() (time.Duration, error) {
	return d.duration.ToDuration()
}
