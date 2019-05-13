package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"time"

	"github.com/jackc/pgx/pgio"
	"github.com/pkg/errors"
)

const pgTimestamptzHourFormat = "2006-01-02 15:04:05.999999999Z07"
const pgTimestamptzMinuteFormat = "2006-01-02 15:04:05.999999999Z07:00"
const pgTimestamptzSecondFormat = "2006-01-02 15:04:05.999999999Z07:00:00"
const microsecFromUnixEpochToY2K = 946684800 * 1000000

const (
	negativeInfinityMicrosecondOffset = -9223372036854775808
	infinityMicrosecondOffset         = 9223372036854775807
)

type Timestamptz struct {
	Time             time.Time
	Status           Status
	InfinityModifier InfinityModifier
}

func (dst *Timestamptz) Set(src interface{}) error {
	if src == nil {
		*dst = Timestamptz{Status: Null}
		return nil
	}

	switch value := src.(type) {
	case time.Time:
		*dst = Timestamptz{Time: value, Status: Present}
	default:
		if originalSrc, ok := underlyingTimeType(src); ok {
			return dst.Set(originalSrc)
		}
		return errors.Errorf("cannot convert %v to Timestamptz", value)
	}

	return nil
}

func (dst *Timestamptz) Get() interface{} {
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

func (src *Timestamptz) AssignTo(dst interface{}) error {
	switch src.Status {
	case Present:
		switch v := dst.(type) {
		case *time.Time:
			if src.InfinityModifier != None {
				return errors.Errorf("cannot assign %v to %T", src, dst)
			}
			*v = src.Time
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

func (dst *Timestamptz) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Timestamptz{Status: Null}
		return nil
	}

	sbuf := string(src)
	switch sbuf {
	case "infinity":
		*dst = Timestamptz{Status: Present, InfinityModifier: Infinity}
	case "-infinity":
		*dst = Timestamptz{Status: Present, InfinityModifier: -Infinity}
	default:
		var format string
		if sbuf[len(sbuf)-9] == '-' || sbuf[len(sbuf)-9] == '+' {
			format = pgTimestamptzSecondFormat
		} else if sbuf[len(sbuf)-6] == '-' || sbuf[len(sbuf)-6] == '+' {
			format = pgTimestamptzMinuteFormat
		} else {
			format = pgTimestamptzHourFormat
		}

		tim, err := time.Parse(format, sbuf)
		if err != nil {
			return err
		}

		*dst = Timestamptz{Time: tim, Status: Present}
	}

	return nil
}

func (dst *Timestamptz) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Timestamptz{Status: Null}
		return nil
	}

	if len(src) != 8 {
		return errors.Errorf("invalid length for timestamptz: %v", len(src))
	}

	microsecSinceY2K := int64(binary.BigEndian.Uint64(src))

	switch microsecSinceY2K {
	case infinityMicrosecondOffset:
		*dst = Timestamptz{Status: Present, InfinityModifier: Infinity}
	case negativeInfinityMicrosecondOffset:
		*dst = Timestamptz{Status: Present, InfinityModifier: -Infinity}
	default:
		microsecSinceUnixEpoch := microsecFromUnixEpochToY2K + microsecSinceY2K
		tim := time.Unix(microsecSinceUnixEpoch/1000000, (microsecSinceUnixEpoch%1000000)*1000)
		*dst = Timestamptz{Time: tim, Status: Present}
	}

	return nil
}

func (src *Timestamptz) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	var s string

	switch src.InfinityModifier {
	case None:
		s = src.Time.UTC().Format(pgTimestamptzSecondFormat)
	case Infinity:
		s = "infinity"
	case NegativeInfinity:
		s = "-infinity"
	}

	return append(buf, s...), nil
}

func (src *Timestamptz) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
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
func (dst *Timestamptz) Scan(src interface{}) error {
	if src == nil {
		*dst = Timestamptz{Status: Null}
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
		*dst = Timestamptz{Time: src, Status: Present}
		return nil
	}

	return errors.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *Timestamptz) Value() (driver.Value, error) {
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
