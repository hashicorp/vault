package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"math"
	"strconv"

	"github.com/jackc/pgx/pgio"
	"github.com/pkg/errors"
)

type Int8 struct {
	Int    int64
	Status Status
}

func (dst *Int8) Set(src interface{}) error {
	if src == nil {
		*dst = Int8{Status: Null}
		return nil
	}

	switch value := src.(type) {
	case int8:
		*dst = Int8{Int: int64(value), Status: Present}
	case uint8:
		*dst = Int8{Int: int64(value), Status: Present}
	case int16:
		*dst = Int8{Int: int64(value), Status: Present}
	case uint16:
		*dst = Int8{Int: int64(value), Status: Present}
	case int32:
		*dst = Int8{Int: int64(value), Status: Present}
	case uint32:
		*dst = Int8{Int: int64(value), Status: Present}
	case int64:
		*dst = Int8{Int: int64(value), Status: Present}
	case uint64:
		if value > math.MaxInt64 {
			return errors.Errorf("%d is greater than maximum value for Int8", value)
		}
		*dst = Int8{Int: int64(value), Status: Present}
	case int:
		if int64(value) < math.MinInt64 {
			return errors.Errorf("%d is greater than maximum value for Int8", value)
		}
		if int64(value) > math.MaxInt64 {
			return errors.Errorf("%d is greater than maximum value for Int8", value)
		}
		*dst = Int8{Int: int64(value), Status: Present}
	case uint:
		if uint64(value) > math.MaxInt64 {
			return errors.Errorf("%d is greater than maximum value for Int8", value)
		}
		*dst = Int8{Int: int64(value), Status: Present}
	case string:
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*dst = Int8{Int: num, Status: Present}
	default:
		if originalSrc, ok := underlyingNumberType(src); ok {
			return dst.Set(originalSrc)
		}
		return errors.Errorf("cannot convert %v to Int8", value)
	}

	return nil
}

func (dst *Int8) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Int
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Int8) AssignTo(dst interface{}) error {
	return int64AssignTo(int64(src.Int), src.Status, dst)
}

func (dst *Int8) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Int8{Status: Null}
		return nil
	}

	n, err := strconv.ParseInt(string(src), 10, 64)
	if err != nil {
		return err
	}

	*dst = Int8{Int: n, Status: Present}
	return nil
}

func (dst *Int8) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Int8{Status: Null}
		return nil
	}

	if len(src) != 8 {
		return errors.Errorf("invalid length for int8: %v", len(src))
	}

	n := int64(binary.BigEndian.Uint64(src))

	*dst = Int8{Int: n, Status: Present}
	return nil
}

func (src *Int8) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return append(buf, strconv.FormatInt(src.Int, 10)...), nil
}

func (src *Int8) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return pgio.AppendInt64(buf, src.Int), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Int8) Scan(src interface{}) error {
	if src == nil {
		*dst = Int8{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case int64:
		*dst = Int8{Int: src, Status: Present}
		return nil
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
func (src *Int8) Value() (driver.Value, error) {
	switch src.Status {
	case Present:
		return int64(src.Int), nil
	case Null:
		return nil, nil
	default:
		return nil, errUndefined
	}
}

func (src *Int8) MarshalJSON() ([]byte, error) {
	switch src.Status {
	case Present:
		return []byte(strconv.FormatInt(src.Int, 10)), nil
	case Null:
		return []byte("null"), nil
	case Undefined:
		return nil, errUndefined
	}

	return nil, errBadStatus
}

func (dst *Int8) UnmarshalJSON(b []byte) error {
	var n int64
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}

	*dst = Int8{Int: n, Status: Present}

	return nil
}
