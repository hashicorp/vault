package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"github.com/jackc/pgio"
)

type Float4 struct {
	Float  float32
	Status Status
}

func (dst *Float4) Set(src interface{}) error {
	if src == nil {
		*dst = Float4{Status: Null}
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
		*dst = Float4{Float: value, Status: Present}
	case float64:
		*dst = Float4{Float: float32(value), Status: Present}
	case int8:
		*dst = Float4{Float: float32(value), Status: Present}
	case uint8:
		*dst = Float4{Float: float32(value), Status: Present}
	case int16:
		*dst = Float4{Float: float32(value), Status: Present}
	case uint16:
		*dst = Float4{Float: float32(value), Status: Present}
	case int32:
		f32 := float32(value)
		if int32(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float32", value)
		}
	case uint32:
		f32 := float32(value)
		if uint32(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float32", value)
		}
	case int64:
		f32 := float32(value)
		if int64(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float32", value)
		}
	case uint64:
		f32 := float32(value)
		if uint64(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float32", value)
		}
	case int:
		f32 := float32(value)
		if int(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float32", value)
		}
	case uint:
		f32 := float32(value)
		if uint(f32) == value {
			*dst = Float4{Float: f32, Status: Present}
		} else {
			return fmt.Errorf("%v cannot be exactly represented as float32", value)
		}
	case string:
		num, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		*dst = Float4{Float: float32(num), Status: Present}
	case *float64:
		if value == nil {
			*dst = Float4{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *float32:
		if value == nil {
			*dst = Float4{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int8:
		if value == nil {
			*dst = Float4{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint8:
		if value == nil {
			*dst = Float4{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int16:
		if value == nil {
			*dst = Float4{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint16:
		if value == nil {
			*dst = Float4{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int32:
		if value == nil {
			*dst = Float4{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint32:
		if value == nil {
			*dst = Float4{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int64:
		if value == nil {
			*dst = Float4{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint64:
		if value == nil {
			*dst = Float4{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *int:
		if value == nil {
			*dst = Float4{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *uint:
		if value == nil {
			*dst = Float4{Status: Null}
		} else {
			return dst.Set(*value)
		}
	case *string:
		if value == nil {
			*dst = Float4{Status: Null}
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

func (dst Float4) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Float
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Float4) AssignTo(dst interface{}) error {
	return float64AssignTo(float64(src.Float), src.Status, dst)
}

func (dst *Float4) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Float4{Status: Null}
		return nil
	}

	n, err := strconv.ParseFloat(string(src), 32)
	if err != nil {
		return err
	}

	*dst = Float4{Float: float32(n), Status: Present}
	return nil
}

func (dst *Float4) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Float4{Status: Null}
		return nil
	}

	if len(src) != 4 {
		return fmt.Errorf("invalid length for float4: %v", len(src))
	}

	n := int32(binary.BigEndian.Uint32(src))

	*dst = Float4{Float: math.Float32frombits(uint32(n)), Status: Present}
	return nil
}

func (src Float4) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = append(buf, strconv.FormatFloat(float64(src.Float), 'f', -1, 32)...)
	return buf, nil
}

func (src Float4) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = pgio.AppendUint32(buf, math.Float32bits(src.Float))
	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Float4) Scan(src interface{}) error {
	if src == nil {
		*dst = Float4{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case float64:
		*dst = Float4{Float: float32(src), Status: Present}
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
func (src Float4) Value() (driver.Value, error) {
	switch src.Status {
	case Present:
		return float64(src.Float), nil
	case Null:
		return nil, nil
	default:
		return nil, errUndefined
	}
}
