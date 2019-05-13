package pgtype

import (
	"database/sql/driver"
	"strconv"

	"github.com/pkg/errors"
)

type Bool struct {
	Bool   bool
	Status Status
}

func (dst *Bool) Set(src interface{}) error {
	if src == nil {
		*dst = Bool{Status: Null}
		return nil
	}

	switch value := src.(type) {
	case bool:
		*dst = Bool{Bool: value, Status: Present}
	case string:
		bb, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		*dst = Bool{Bool: bb, Status: Present}
	default:
		if originalSrc, ok := underlyingBoolType(src); ok {
			return dst.Set(originalSrc)
		}
		return errors.Errorf("cannot convert %v to Bool", value)
	}

	return nil
}

func (dst *Bool) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Bool
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Bool) AssignTo(dst interface{}) error {
	switch src.Status {
	case Present:
		switch v := dst.(type) {
		case *bool:
			*v = src.Bool
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

func (dst *Bool) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Bool{Status: Null}
		return nil
	}

	if len(src) != 1 {
		return errors.Errorf("invalid length for bool: %v", len(src))
	}

	*dst = Bool{Bool: src[0] == 't', Status: Present}
	return nil
}

func (dst *Bool) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Bool{Status: Null}
		return nil
	}

	if len(src) != 1 {
		return errors.Errorf("invalid length for bool: %v", len(src))
	}

	*dst = Bool{Bool: src[0] == 1, Status: Present}
	return nil
}

func (src *Bool) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	if src.Bool {
		buf = append(buf, 't')
	} else {
		buf = append(buf, 'f')
	}

	return buf, nil
}

func (src *Bool) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	if src.Bool {
		buf = append(buf, 1)
	} else {
		buf = append(buf, 0)
	}

	return buf, nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Bool) Scan(src interface{}) error {
	if src == nil {
		*dst = Bool{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case bool:
		*dst = Bool{Bool: src, Status: Present}
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
func (src *Bool) Value() (driver.Value, error) {
	switch src.Status {
	case Present:
		return src.Bool, nil
	case Null:
		return nil, nil
	default:
		return nil, errUndefined
	}
}
