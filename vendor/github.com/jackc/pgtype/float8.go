package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"github.com/jackc/pgio"
)

type Float8 struct {
	Float  float64
	Status Status
}

func (dst *Float8) Set(src interface{}) error {
	if src == nil {
		*dst = Float8{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case float32:
		*dst = Float8{Float: float64(value), Status: Present}
	case float64:
		*dst = Float8{Float: value, Status: Present}
	case int8:
		*dst = Float8{Float: float64(value), Status: Present}
	case uint8:
		*dst = Float8{Float: float64(value), Status: Present}
	case int16:
		*dst = Float8{Float: float64(value), Status: Present}
	case uint16:
		*dst = Float8{Float: float64(value), Status: Present}
	case int32:
		*dst = Float8{Float: float64(value), Status: Present}
	case uint32:
		*dst = Float8{Float: float64(value), Status: Present}
	case int64:
		f64 := float64(value)
		if int64(f64) == value {
			*dst = Float8{Float: f64, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float64", value)
		}
	case uint64:
		f64 := float64(value)
		if uint64(f64) == value {
			*dst = Float8{Float: f64, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float64", value)
		}
	case int:
		f64 := float64(value)
		if int(f64) == value {
			*dst = Float8{Float: f64, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float64", value)
		}
	case uint:
		f64 := float64(value)
		if uint(f64) == value {
			*dst = Float8{Float: f64, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float64", value)
		}
	case string:
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		*dst = Float8{Float: float64(num), Status: Present}
	case *float64:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *float32:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int8:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint8:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int16:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint16:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int32:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint32:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int64:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint64:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *string:
		if value == nil {
			*dst = Float8{Status: Null}
		} else {
			return dst.Set(*value)
		}
	default:
		if originalSrc, ok := underlyingNumberType(src); ok {
			return dst.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to Float8", value)
	}

	return nil
}

func (dst Float8) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Float
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Float8) AssignTo(dst interface{}) error {
	return float64AssignTo(src.Float, src.Status, dst)
}

func (dst *Float8) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Float8{Status: Null}
		return nil
	}

	n, err := strconv.ParseFloat(string(src), 64)
	if err != nil {
		return err
	}

	*dst = Float8{Float: n, Status: Present}
	return nil
}

func (dst *Float8) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Float8{Status: Null}
		return nil
	}

	if len(src) != 8 {
		return fmt.Errorf("invalid length for float8: %v", len(src))
	}

	n := int64(binary.BigEndian.Uint64(src))

	*dst = Float8{Float: math.Float64frombits(uint64(n)), Status: Present}
	return nil
}

func (src Float8) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = append(buf, strconv.FormatFloat(float64(src.Float), 'f', -1, 64)...)
	return buf, nil
}

func (src Float8) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = pgio.AppendUint64(buf, math.Float64bits(src.Float))
	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Float8) Scan(src interface{}) error {
	if src == nil {
		*dst = Float8{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case float64:
		*dst = Float8{Float: src, Status: Present}
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
func (src Float8) Value() (driver.Value, error) {
	switch src.Status {
	case Present:
		return src.Float, nil
	case Null:
		return nil, nil
	default:
		return nil, errUndefined
	}
}
