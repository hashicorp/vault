package pgtype

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Text struct {
	String string
	Status Status
}

func (dst *Text) Set(src interface{}) error {
	if src == nil {
		*dst = Text{Status: Null}
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case string:
		*dst = Text{String: value, Status: Present}
	case *string:
		if value == nil {
			*dst = Text{Status: Null}
		} else {
			*dst = Text{String: *value, Status: Present}
		}
	case []byte:
		if value == nil {
			*dst = Text{Status: Null}
		} else {
			*dst = Text{String: string(value), Status: Present}
		}
	case fmt.Stringer:
		if value == fmt.Stringer(nil) {
			*dst = Text{Status: Null}
		} else {
			*dst = Text{String: value.String(), Status: Present}
		}
	default:
		// Cannot be part of the switch: If Value() returns nil on
		// non-string, we should still try to checks the underlying type
		// using reflection.
		//
		// For example the struct might implement driver.Valuer with
		// pointer receiver and fmt.Stringer with value receiver.
		if value, ok := src.(driver.Valuer); ok {
			if value == driver.Valuer(nil) {
				*dst = Text{Status: Null}
				return nil
			} else {
				v, err := value.Value()
				if err != nil {
					return fmt.Errorf("driver.Valuer Value() method failed: %w", err)
				}

				// Handles also v == nil case.
				if s, ok := v.(string); ok {
					*dst = Text{String: s, Status: Present}
					return nil
				}
			}
		}

		if originalSrc, ok := underlyingStringType(src); ok {
			return dst.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to Text", value)
	}

	return nil
}

func (dst Text) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.String
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Text) AssignTo(dst interface{}) error {
	switch src.Status {
	case Present:
		switch v := dst.(type) {
		case *string:
			*v = src.String
			return nil
		case *[]byte:
			*v = make([]byte, len(src.String))
			copy(*v, src.String)
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

func (Text) PreferredResultFormat() int16 {
	return TextFormatCode
}

func (dst *Text) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Text{Status: Null}
		return nil
	}

	*dst = Text{String: string(src), Status: Present}
	return nil
}

func (dst *Text) DecodeBinary(ci *ConnInfo, src []byte) error {
	return dst.DecodeText(ci, src)
}

func (Text) PreferredParamFormat() int16 {
	return TextFormatCode
}

func (src Text) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return append(buf, src.String...), nil
}

func (src Text) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return src.EncodeText(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *Text) Scan(src interface{}) error {
	if src == nil {
		*dst = Text{Status: Null}
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

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src Text) Value() (driver.Value, error) {
	switch src.Status {
	case Present:
		return src.String, nil
	case Null:
		return nil, nil
	default:
		return nil, errUndefined
	}
}

func (src Text) MarshalJSON() ([]byte, error) {
	switch src.Status {
	case Present:
		return json.Marshal(src.String)
	case Null:
		return []byte("null"), nil
	case Undefined:
		return nil, errUndefined
	}

	return nil, errBadStatus
}

func (dst *Text) UnmarshalJSON(b []byte) error {
	var s *string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	if s == nil {
		*dst = Text{Status: Null}
	} else {
		*dst = Text{String: *s, Status: Present}
	}

	return nil
}
