package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/jackc/pgio"
)

type Int2 struct {
	Int    int16
	Status Status
}

func (dst *Int2) Set(src interface{}) error {
	if src == nil {
		*dst = Int2{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case int8:
		*dst = Int2{Int: int16(value), Status: Present}
	case uint8:
		*dst = Int2{Int: int16(value), Status: Present}
	case int16:
		*dst = Int2{Int: int16(value), Status: Present}
	case uint16:
		if value > math.MaxInt16 {
			return fmt.Errorf("%d is greater than maximum value for Int2", value)
		}
		*dst = Int2{Int: int16(value), Status: Present}
	case int32:
		if value < math.MinInt16 {
			return fmt.Errorf("%d is greater than maximum value for Int2", value)
		}
		if value > math.MaxInt16 {
			return fmt.Errorf("%d is greater than maximum value for Int2", value)
		}
		*dst = Int2{Int: int16(value), Status: Present}
	case uint32:
		if value > math.MaxInt16 {
			return fmt.Errorf("%d is greater than maximum value for Int2", value)
		}
		*dst = Int2{Int: int16(value), Status: Present}
	case int64:
		if value < math.MinInt16 {
			return fmt.Errorf("%d is greater than maximum value for Int2", value)
		}
		if value > math.MaxInt16 {
			return fmt.Errorf("%d is greater than maximum value for Int2", value)
		}
		*dst = Int2{Int: int16(value), Status: Present}
	case uint64:
		if value > math.MaxInt16 {
			return fmt.Errorf("%d is greater than maximum value for Int2", value)
		}
		*dst = Int2{Int: int16(value), Status: Present}
	case int:
		if value < math.MinInt16 {
			return fmt.Errorf("%d is greater than maximum value for Int2", value)
		}
		if value > math.MaxInt16 {
			return fmt.Errorf("%d is greater than maximum value for Int2", value)
		}
		*dst = Int2{Int: int16(value), Status: Present}
	case uint:
		if value > math.MaxInt16 {
			return fmt.Errorf("%d is greater than maximum value for Int2", value)
		}
		*dst = Int2{Int: int16(value), Status: Present}
	case string:
		num, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return err
		}
		*dst = Int2{Int: int16(num), Status: Present}
	case float32:
		if value > math.MaxInt16 {
			return fmt.Errorf("%f is greater than maximum value for Int2", value)
		}
		*dst = Int2{Int: int16(value), Status: Present}
	case float64:
		if value > math.MaxInt16 {
			return fmt.Errorf("%f is greater than maximum value for Int2", value)
		}
		*dst = Int2{Int: int16(value), Status: Present}
	case *int8:
		if value == nil {
			*dst = Int2{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint8:
		if value == nil {
			*dst = Int2{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int16:
		if value == nil {
			*dst = Int2{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint16:
		if value == nil {
			*dst = Int2{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int32:
		if value == nil {
			*dst = Int2{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint32:
		if value == nil {
			*dst = Int2{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int64:
		if value == nil {
			*dst = Int2{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint64:
		if value == nil {
			*dst = Int2{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int:
		if value == nil {
			*dst = Int2{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint:
		if value == nil {
			*dst = Int2{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *string:
		if value == nil {
			*dst = Int2{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *float32:
		if value == nil {
			*dst = Int2{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *float64:
		if value == nil {
			*dst = Int2{Status: Null}
		} else {
			return dst.Set(*value)
		}
	default:
		if originalSrc, ok := underlyingNumberType(src); ok {
			return dst.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to Int2", value)
	}

	return nil
}

func (dst Int2) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Int
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Int2) AssignTo(dst interface{}) error {
	return int64AssignTo(int64(src.Int), src.Status, dst)
}

func (dst *Int2) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Int2{Status: Null}
		return nil
	}

	n, err := strconv.ParseInt(string(src), 10, 16)
	if err != nil {
		return err
	}

	*dst = Int2{Int: int16(n), Status: Present}
	return nil
}

func (dst *Int2) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Int2{Status: Null}
		return nil
	}

	if len(src) != 2 {
		return fmt.Errorf("invalid length for int2: %v", len(src))
	}

	n := int16(binary.BigEndian.Uint16(src))
	*dst = Int2{Int: n, Status: Present}
	return nil
}

func (src Int2) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return append(buf, strconv.FormatInt(int64(src.Int), 10)...), nil
}

func (src Int2) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return pgio.AppendInt16(buf, src.Int), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Int2) Scan(src interface{}) error {
	if src == nil {
		*dst = Int2{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case int64:
		if src < math.MinInt16 {
			return fmt.Errorf("%d is greater than maximum value for Int2", src)
		}
		if src > math.MaxInt16 {
			return fmt.Errorf("%d is greater than maximum value for Int2", src)
		}
		*dst = Int2{Int: int16(src), Status: Present}
		return nil
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		return dst.DecodeText(nil, srcCopy)
	}

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src Int2) Value() (driver.Value, error) {
	switch src.Status {
	case Present:
		return int64(src.Int), nil
	case Null:
		return nil, nil
	default:
		return nil, errUndefined
	}
}

func (src Int2) MarshalJSON() ([]byte, error) {
	switch src.Status {
	case Present:
		return []byte(strconv.FormatInt(int64(src.Int), 10)), nil
	case Null:
		return []byte("null"), nil
	case Undefined:
		return nil, errUndefined
	}

	return nil, errBadStatus
}

func (dst *Int2) UnmarshalJSON(b []byte) error {
	var n *int16
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}

	if n == nil {
		*dst = Int2{Status: Null}
	} else {
		*dst = Int2{Int: *n, Status: Present}
	}

	return nil
}
