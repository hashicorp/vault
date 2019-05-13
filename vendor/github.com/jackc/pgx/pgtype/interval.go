package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/pgio"
	"github.com/pkg/errors"
)

const (
	microsecondsPerSecond = 1000000
	microsecondsPerMinute = 60 * microsecondsPerSecond
	microsecondsPerHour   = 60 * microsecondsPerMinute
)

type Interval struct {
	Microseconds int64
	Days         int32
	Months       int32
	Status       Status
}

func (dst *Interval) Set(src interface{}) error {
	if src == nil {
		*dst = Interval{Status: Null}
		return nil
	}

	switch value := src.(type) {
	case time.Duration:
		*dst = Interval{Microseconds: int64(value) / 1000, Status: Present}
	default:
		if originalSrc, ok := underlyingPtrType(src); ok {
			return dst.Set(originalSrc)
		}
		return errors.Errorf("cannot convert %v to Interval", value)
	}

	return nil
}

func (dst *Interval) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Interval) AssignTo(dst interface{}) error {
	switch src.Status {
	case Present:
		switch v := dst.(type) {
		case *time.Duration:
			if src.Days > 0 || src.Months > 0 {
				return errors.Errorf("interval with months or days cannot be decoded into %T", dst)
			}
			*v = time.Duration(src.Microseconds) * time.Microsecond
			return nil
		default:
			if nextDst, retry := GetAssignToDstType(dst); retry {
				return src.AssignTo(nextDst)
			}
		}
	case Null:
		return NullAssignTo(dst)
	}

	return errors.Errorf("cannot decode %#v into %T", src, dst)
}

func (dst *Interval) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Interval{Status: Null}
		return nil
	}

	var microseconds int64
	var days int32
	var months int32

	parts := strings.Split(string(src), " ")

	for i := 0; i < len(parts)-1; i += 2 {
		scalar, err := strconv.ParseInt(parts[i], 10, 64)
		if err != nil {
			return errors.Errorf("bad interval format")
		}

		switch parts[i+1] {
		case "year", "years":
			months += int32(scalar * 12)
		case "mon", "mons":
			months += int32(scalar)
		case "day", "days":
			days = int32(scalar)
		}
	}

	if len(parts)%2 == 1 {
		timeParts := strings.SplitN(parts[len(parts)-1], ":", 3)
		if len(timeParts) != 3 {
			return errors.Errorf("bad interval format")
		}

		var negative bool
		if timeParts[0][0] == '-' {
			negative = true
			timeParts[0] = timeParts[0][1:]
		}

		hours, err := strconv.ParseInt(timeParts[0], 10, 64)
		if err != nil {
			return errors.Errorf("bad interval hour format: %s", timeParts[0])
		}

		minutes, err := strconv.ParseInt(timeParts[1], 10, 64)
		if err != nil {
			return errors.Errorf("bad interval minute format: %s", timeParts[1])
		}

		secondParts := strings.SplitN(timeParts[2], ".", 2)

		seconds, err := strconv.ParseInt(secondParts[0], 10, 64)
		if err != nil {
			return errors.Errorf("bad interval second format: %s", secondParts[0])
		}

		var uSeconds int64
		if len(secondParts) == 2 {
			uSeconds, err = strconv.ParseInt(secondParts[1], 10, 64)
			if err != nil {
				return errors.Errorf("bad interval decimal format: %s", secondParts[1])
			}

			for i := 0; i < 6-len(secondParts[1]); i++ {
				uSeconds *= 10
			}
		}

		microseconds = hours * microsecondsPerHour
		microseconds += minutes * microsecondsPerMinute
		microseconds += seconds * microsecondsPerSecond
		microseconds += uSeconds

		if negative {
			microseconds = -microseconds
		}
	}

	*dst = Interval{Months: months, Days: days, Microseconds: microseconds, Status: Present}
	return nil
}

func (dst *Interval) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Interval{Status: Null}
		return nil
	}

	if len(src) != 16 {
		return errors.Errorf("Received an invalid size for a interval: %d", len(src))
	}

	microseconds := int64(binary.BigEndian.Uint64(src))
	days := int32(binary.BigEndian.Uint32(src[8:]))
	months := int32(binary.BigEndian.Uint32(src[12:]))

	*dst = Interval{Microseconds: microseconds, Days: days, Months: months, Status: Present}
	return nil
}

func (src *Interval) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	if src.Months != 0 {
		buf = append(buf, strconv.FormatInt(int64(src.Months), 10)...)
		buf = append(buf, " mon "...)
	}

	if src.Days != 0 {
		buf = append(buf, strconv.FormatInt(int64(src.Days), 10)...)
		buf = append(buf, " day "...)
	}

	absMicroseconds := src.Microseconds
	if absMicroseconds < 0 {
		absMicroseconds = -absMicroseconds
		buf = append(buf, '-')
	}

	hours := absMicroseconds / microsecondsPerHour
	minutes := (absMicroseconds % microsecondsPerHour) / microsecondsPerMinute
	seconds := (absMicroseconds % microsecondsPerMinute) / microsecondsPerSecond
	microseconds := absMicroseconds % microsecondsPerSecond

	timeStr := fmt.Sprintf("%02d:%02d:%02d.%06d", hours, minutes, seconds, microseconds)
	return append(buf, timeStr...), nil
}

// EncodeBinary encodes src into w.
func (src *Interval) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = pgio.AppendInt64(buf, src.Microseconds)
	buf = pgio.AppendInt32(buf, src.Days)
	return pgio.AppendInt32(buf, src.Months), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Interval) Scan(src interface{}) error {
	if src == nil {
		*dst = Interval{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		return dst.DecodeText(nil, srcCopy)
	}

	return errors.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *Interval) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
