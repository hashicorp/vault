package encoding

import (
	"time"

	"github.com/SAP/go-hdb/driver/internal/protocol/julian"
)

// Longdate.
func convertLongdateToTime(longdate int64) time.Time {
	const dayfactor = 10000000 * 24 * 60 * 60
	longdate--
	d := (longdate % dayfactor) * 100
	t := convertDaydateToTime((longdate / dayfactor) + 1)
	return t.Add(time.Duration(d))
}

// nanosecond: HDB - 7 digits precision (not 9 digits).
func convertTimeToLongdate(t time.Time) int64 {
	return (((((((convertTimeToDayDate(t)-1)*24)+int64(t.Hour()))*60)+int64(t.Minute()))*60)+int64(t.Second()))*1e7 + int64(t.Nanosecond()/1e2) + 1
}

// Seconddate.
func convertSeconddateToTime(seconddate int64) time.Time {
	const dayfactor = 24 * 60 * 60
	seconddate--
	d := (seconddate % dayfactor) * 1e9
	t := convertDaydateToTime((seconddate / dayfactor) + 1)
	return t.Add(time.Duration(d))
}
func convertTimeToSeconddate(t time.Time) int64 {
	return (((((convertTimeToDayDate(t)-1)*24)+int64(t.Hour()))*60)+int64(t.Minute()))*60 + int64(t.Second()) + 1
}

const julianHdb = 1721423 // 1 January 0001 00:00:00 (1721424) - 1

// Daydate.
func convertDaydateToTime(daydate int64) time.Time {
	return julian.DayToTime(int(daydate) + julianHdb)
}
func convertTimeToDayDate(t time.Time) int64 {
	return int64(julian.TimeToDay(t) - julianHdb)
}

// Secondtime.
func convertSecondtimeToTime(secondtime int) time.Time {
	return time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(int64(secondtime-1) * 1e9))
}
func convertTimeToSecondtime(t time.Time) int {
	return (t.Hour()*60+t.Minute())*60 + t.Second() + 1
}
