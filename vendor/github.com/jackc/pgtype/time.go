package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgio"
)

// Time represents the PostgreSQL time type. The PostgreSQL time is a time of day without time zone.
//
// Time is represented as the number of microseconds since midnight in the same way that PostgreSQL does. Other time
// and date types in pgtype can use time.Time as the underlying representation. However, pgtype.Time type cannot due
// to needing to handle 24:00:00. time.Time converts that to 00:00:00 on the following day.
type Time struct {
	Microseconds int64 // Number of microseconds since midnight
	Status       Status
}

// Set converts src into a Time and stores in dst.
func (dst *Time) Set(src interface{}) error {
	if src == nil {
		*dst = Time{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case time.Time:
		usec := int64(value.Hour())*microsecondsPerHour +
			int64(value.Minute())*microsecondsPerMinute +
			int64(value.Second())*microsecondsPerSecond +
			int64(value.Nanosecond())/1000
		*dst = Time{Microseconds: usec, Status: Present}
	case *time.Time:
		if value == nil {
			*dst = Time{Status: Null}
		} else {
			return dst.Set(*value)
		}
	default:
		if originalSrc, ok := underlyingTimeType(src); ok {
			return dst.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to Time", value)
	}

	return nil
}

func (dst Time) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Microseconds
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Time) AssignTo(dst interface{}) error {
	switch src.Status {
	case Present:
		switch v := dst.(type) {
		case *time.Time:
			// 24:00:00 is max allowed time in PostgreSQL, but time.Time will normalize that to 00:00:00 the next day.
			var maxRepresentableByTime int64 = 24*60*60*1000000 - 1
			if src.Microseconds > maxRepresentableByTime {
				return fmt.Errorf("%d microseconds cannot be represented as time.Time", src.Microseconds)
			}

			usec := src.Microseconds
			hours := usec / microsecondsPerHour
			usec -= hours * microsecondsPerHour
			minutes := usec / microsecondsPerMinute
			usec -= minutes * microsecondsPerMinute
			seconds := usec / microsecondsPerSecond
			usec -= seconds * microsecondsPerSecond
			ns := usec * 1000
			*v = time.Date(2000, 1, 1, int(hours), int(minutes), int(seconds), int(ns), time.UTC)
			return nil
		default:
			if nextDst, retry := GetAssignToDstType(dst); retry {
				return src.AssignTo(nextDst)
			}
			return fmt.Errorf("unable to assign to %T", dst)
		}
	case Null:
		return NullAssignTo(dst)
	}

	return fmt.Errorf("cannot decode %#v into %T", src, dst)
}

// DecodeText decodes from src into dst.
func (dst *Time) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Time{Status: Null}
		return nil
	}

	s := string(src)

	if len(s) < 8 {
		return fmt.Errorf("cannot decode %v into Time", s)
	}

	hours, err := strconv.ParseInt(s[0:2], 10, 64)
	if err != nil {
		return fmt.Errorf("cannot decode %v into Time", s)
	}
	usec := hours * microsecondsPerHour

	minutes, err := strconv.ParseInt(s[3:5], 10, 64)
	if err != nil {
		return fmt.Errorf("cannot decode %v into Time", s)
	}
	usec += minutes * microsecondsPerMinute

	seconds, err := strconv.ParseInt(s[6:8], 10, 64)
	if err != nil {
		return fmt.Errorf("cannot decode %v into Time", s)
	}
	usec += seconds * microsecondsPerSecond

	if len(s) > 9 {
		fraction := s[9:]
		n, err := strconv.ParseInt(fraction, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot decode %v into Time", s)
		}

		for i := len(fraction); i < 6; i++ {
			n *= 10
		}

		usec += n
	}

	*dst = Time{Microseconds: usec, Status: Present}

	return nil
}

// DecodeBinary decodes from src into dst.
func (dst *Time) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Time{Status: Null}
		return nil
	}

	if len(src) != 8 {
		return fmt.Errorf("invalid length for time: %v", len(src))
	}

	usec := int64(binary.BigEndian.Uint64(src))
	*dst = Time{Microseconds: usec, Status: Present}

	return nil
}

// EncodeText writes the text encoding of src into w.
func (src Time) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	usec := src.Microseconds
	hours := usec / microsecondsPerHour
	usec -= hours * microsecondsPerHour
	minutes := usec / microsecondsPerMinute
	usec -= minutes * microsecondsPerMinute
	seconds := usec / microsecondsPerSecond
	usec -= seconds * microsecondsPerSecond

	s := fmt.Sprintf("%02d:%02d:%02d.%06d", hours, minutes, seconds, usec)

	return append(buf, s...), nil
}

// EncodeBinary writes the binary encoding of src into w. If src.Time is not in
// the UTC time zone it returns an error.
func (src Time) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return pgio.AppendInt64(buf, src.Microseconds), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Time) Scan(src interface{}) error {
	if src == nil {
		*dst = Time{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		return dst.DecodeText(nil, srcCopy)
	case time.Time:
		return dst.Set(src)
	}

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src Time) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
