package pgtype

import (
	"database/sql"
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgio"
)

type Date struct {
	Time             time.Time
	Status           Status
	InfinityModifier InfinityModifier
}

const (
	negativeInfinityDayOffset = -2147483648
	infinityDayOffset         = 2147483647
)

func (dst *Date) Set(src interface{}) error {
	if src == nil {
		*dst = Date{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	if value, ok := src.(interface{ Value() (driver.Value, error) }); ok {
		v, err := value.Value()
		if err != nil {
			return fmt.Errorf("cannot get value %v for Date: %v", value, err)
		}
		return dst.Set(v)
	}

	switch value := src.(type) {
	case time.Time:
		*dst = Date{Time: value, Status: Present}
	case *time.Time:
		if value == nil {
			*dst = Date{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case string:
		return dst.DecodeText(nil, []byte(value))
	case *string:
		if value == nil {
			*dst = Date{Status: Null}
		} else {
			return dst.Set(*value)
		}
	default:
		if originalSrc, ok := underlyingTimeType(src); ok {
			return dst.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to Date", value)
	}

	return nil
}

func (dst Date) Get() interface{} {
	switch dst.Status {
	case Present:
		if dst.InfinityModifier != None {
			return dst.InfinityModifier
		}
		return dst.Time
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Date) AssignTo(dst interface{}) error {
	if scanner, ok := dst.(sql.Scanner); ok {
		var err error
		switch src.Status {
		case Present:
			if src.InfinityModifier != None {
				err = scanner.Scan(src.InfinityModifier.String())
			} else {
				err = scanner.Scan(src.Time)
			}
		case Null:
			err = scanner.Scan(nil)
		}
		if err != nil {
			return fmt.Errorf("unable assign %v to %T: %s", src, dst, err)
		}
		return nil
	}

	switch src.Status {
	case Present:
		switch v := dst.(type) {
		case *time.Time:
			if src.InfinityModifier != None {
				return fmt.Errorf("cannot assign %v to %T", src, dst)
			}
			*v = src.Time
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

func (dst *Date) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Date{Status: Null}
		return nil
	}

	sbuf := string(src)
	switch sbuf {
	case "infinity":
		*dst = Date{Status: Present, InfinityModifier: Infinity}
	case "-infinity":
		*dst = Date{Status: Present, InfinityModifier: -Infinity}
	default:
		if strings.HasSuffix(sbuf, " BC") {
			t, err := time.ParseInLocation("2006-01-02", strings.TrimRight(sbuf, " BC"), time.UTC)
			t2 := time.Date(1-t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
			if err != nil {
				return err
			}
			*dst = Date{Time: t2, Status: Present}
			return nil
		}
		t, err := time.ParseInLocation("2006-01-02", sbuf, time.UTC)
		if err != nil {
			return err
		}

		*dst = Date{Time: t, Status: Present}
	}

	return nil
}

func (dst *Date) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Date{Status: Null}
		return nil
	}

	if len(src) != 4 {
		return fmt.Errorf("invalid length for date: %v", len(src))
	}

	dayOffset := int32(binary.BigEndian.Uint32(src))

	switch dayOffset {
	case infinityDayOffset:
		*dst = Date{Status: Present, InfinityModifier: Infinity}
	case negativeInfinityDayOffset:
		*dst = Date{Status: Present, InfinityModifier: -Infinity}
	default:
		t := time.Date(2000, 1, int(1+dayOffset), 0, 0, 0, 0, time.UTC)
		*dst = Date{Time: t, Status: Present}
	}

	return nil
}

func (src Date) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	var s string

	switch src.InfinityModifier {
	case None:
		s = src.Time.Format("2006-01-02")
	case Infinity:
		s = "infinity"
	case NegativeInfinity:
		s = "-infinity"
	}

	return append(buf, s...), nil
}

func (src Date) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	var daysSinceDateEpoch int32
	switch src.InfinityModifier {
	case None:
		tUnix := time.Date(src.Time.Year(), src.Time.Month(), src.Time.Day(), 0, 0, 0, 0, time.UTC).Unix()
		dateEpoch := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Unix()

		secSinceDateEpoch := tUnix - dateEpoch
		daysSinceDateEpoch = int32(secSinceDateEpoch / 86400)
	case Infinity:
		daysSinceDateEpoch = infinityDayOffset
	case NegativeInfinity:
		daysSinceDateEpoch = negativeInfinityDayOffset
	}

	return pgio.AppendInt32(buf, daysSinceDateEpoch), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Date) Scan(src interface{}) error {
	if src == nil {
		*dst = Date{Status: Null}
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
		*dst = Date{Time: src, Status: Present}
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src Date) Value() (driver.Value, error) {
	switch src.Status {
	case Present:
		if src.InfinityModifier != None {
			return src.InfinityModifier.String(), nil
		}
		return src.Time, nil
	case Null:
		return nil, nil
	default:
		return nil, errUndefined
	}
}

func (src Date) MarshalJSON() ([]byte, error) {
	switch src.Status {
	case Null:
		return []byte("null"), nil
	case Undefined:
		return nil, errUndefined
	}

	if src.Status != Present {
		return nil, errBadStatus
	}

	var s string

	switch src.InfinityModifier {
	case None:
		s = src.Time.Format("2006-01-02")
	case Infinity:
		s = "infinity"
	case NegativeInfinity:
		s = "-infinity"
	}

	return json.Marshal(s)
}

func (dst *Date) UnmarshalJSON(b []byte) error {
	var s *string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	if s == nil {
		*dst = Date{Status: Null}
		return nil
	}

	switch *s {
	case "infinity":
		*dst = Date{Status: Present, InfinityModifier: Infinity}
	case "-infinity":
		*dst = Date{Status: Present, InfinityModifier: -Infinity}
	default:
		t, err := time.ParseInLocation("2006-01-02", *s, time.UTC)
		if err != nil {
			return err
		}

		*dst = Date{Time: t, Status: Present}
	}

	return nil
}
