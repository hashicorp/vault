package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgio"
)

const pgTimestampFormat = "2006-01-02 15:04:05.999999999"

// Timestamp represents the PostgreSQL timestamp type. The PostgreSQL
// timestamp does not have a time zone. This presents a problem when
// translating to and from time.Time which requires a time zone. It is highly
// recommended to use timestamptz whenever possible. Timestamp methods either
// convert to UTC or return an error on non-UTC times.
type Timestamp struct {
	Time             time.Time // Time must always be in UTC.
	Status           Status
	InfinityModifier InfinityModifier
}

// Set converts src into a Timestamp and stores in dst. If src is a
// time.Time in a non-UTC time zone, the time zone is discarded.
func (dst *Timestamp) Set(src interface{}) error {
	if src == nil {
		*dst = Timestamp{Status: Null}
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
		*dst = Timestamp{Time: time.Date(value.Year(), value.Month(), value.Day(), value.Hour(), value.Minute(), value.Second(), value.Nanosecond(), time.UTC), Status: Present}
	case *time.Time:
		if value == nil {
			*dst = Timestamp{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case string:
		return dst.DecodeText(nil, []byte(value))
	case *string:
		if value == nil {
			*dst = Timestamp{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case InfinityModifier:
		*dst = Timestamp{InfinityModifier: value, Status: Present}
	default:
		if originalSrc, ok := underlyingTimeType(src); ok {
			return dst.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to Timestamp", value)
	}

	return nil
}

func (dst Timestamp) Get() interface{} {
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

func (src *Timestamp) AssignTo(dst interface{}) error {
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

// DecodeText decodes from src into dst. The decoded time is considered to
// be in UTC.
func (dst *Timestamp) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Timestamp{Status: Null}
		return nil
	}

	sbuf := string(src)
	switch sbuf {
	case "infinity":
		*dst = Timestamp{Status: Present, InfinityModifier: Infinity}
	case "-infinity":
		*dst = Timestamp{Status: Present, InfinityModifier: -Infinity}
	default:
		if strings.HasSuffix(sbuf, " BC") {
			t, err := time.Parse(pgTimestampFormat, strings.TrimRight(sbuf, " BC"))
			t2 := time.Date(1-t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
			if err != nil {
				return err
			}
			*dst = Timestamp{Time: t2, Status: Present}
			return nil
		}
		tim, err := time.Parse(pgTimestampFormat, sbuf)
		if err != nil {
			return err
		}

		*dst = Timestamp{Time: tim, Status: Present}
	}

	return nil
}

// DecodeBinary decodes from src into dst. The decoded time is considered to
// be in UTC.
func (dst *Timestamp) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Timestamp{Status: Null}
		return nil
	}

	if len(src) != 8 {
		return fmt.Errorf("invalid length for timestamp: %v", len(src))
	}

	microsecSinceY2K := int64(binary.BigEndian.Uint64(src))

	switch microsecSinceY2K {
	case infinityMicrosecondOffset:
		*dst = Timestamp{Status: Present, InfinityModifier: Infinity}
	case negativeInfinityMicrosecondOffset:
		*dst = Timestamp{Status: Present, InfinityModifier: -Infinity}
	default:
		tim := time.Unix(
			microsecFromUnixEpochToY2K/1000000+microsecSinceY2K/1000000,
			(microsecFromUnixEpochToY2K%1000000*1000)+(microsecSinceY2K%1000000*1000),
		).UTC()
		*dst = Timestamp{Time: tim, Status: Present}
	}

	return nil
}

// EncodeText writes the text encoding of src into w. If src.Time is not in
// the UTC time zone it returns an error.
func (src Timestamp) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}
	if src.Time.Location() != time.UTC {
		return nil, fmt.Errorf("cannot encode non-UTC time into timestamp")
	}

	var s string

	switch src.InfinityModifier {
	case None:
		s = src.Time.Truncate(time.Microsecond).Format(pgTimestampFormat)
	case Infinity:
		s = "infinity"
	case NegativeInfinity:
		s = "-infinity"
	}

	return append(buf, s...), nil
}

// EncodeBinary writes the binary encoding of src into w. If src.Time is not in
// the UTC time zone it returns an error.
func (src Timestamp) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}
	if src.Time.Location() != time.UTC {
		return nil, fmt.Errorf("cannot encode non-UTC time into timestamp")
	}

	var microsecSinceY2K int64
	switch src.InfinityModifier {
	case None:
		microsecSinceUnixEpoch := src.Time.Unix()*1000000 + int64(src.Time.Nanosecond())/1000
		microsecSinceY2K = microsecSinceUnixEpoch - microsecFromUnixEpochToY2K
	case Infinity:
		microsecSinceY2K = infinityMicrosecondOffset
	case NegativeInfinity:
		microsecSinceY2K = negativeInfinityMicrosecondOffset
	}

	return pgio.AppendInt64(buf, microsecSinceY2K), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Timestamp) Scan(src interface{}) error {
	if src == nil {
		*dst = Timestamp{Status: Null}
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
		*dst = Timestamp{Time: src, Status: Present}
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src Timestamp) Value() (driver.Value, error) {
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
