// Package julian provided julian time conversion functions.
package julian

import (
	"time"
)

const gregorianDay = 2299161                // Start date of Gregorian Calendar as Julian Day Number
var gregorianDate = DayToTime(gregorianDay) // Start date of Gregorian Calendar (1582-10-15)

// TimeToDay returns the Julian Date Number of time's date components.
// The algorithm is taken from https://en.wikipedia.org/wiki/Julian_day.
func TimeToDay(t time.Time) int {
	t = t.UTC()

	month := int(t.Month())

	a := (14 - month) / 12
	y := t.Year() + 4800 - a
	m := month + (12 * a) - 3

	if t.Before(gregorianDate) { // Julian Calendar
		return t.Day() + (153*m+2)/5 + 365*y + y/4 - 32083
	}
	// Gregorian Calendar
	return t.Day() + (153*m+2)/5 + 365*y + y/4 - y/100 + y/400 - 32045
}

// DayToTime returns the correcponding UTC date for a Julian Day Number.
// The algorithm is taken from https://en.wikipedia.org/wiki/Julian_day.
func DayToTime(jd int) time.Time {
	var f int

	if jd < gregorianDay {
		f = jd + 1401
	} else {
		f = jd + 1401 + (((4*jd+274277)/146097)*3)/4 - 38
	}

	e := 4*f + 3
	g := (e % 1461) / 4
	h := 5*g + 2
	day := (h%153)/5 + 1
	month := (h/153+2)%12 + 1
	year := (e / 1461) - 4716 + (12+2-month)/12

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
