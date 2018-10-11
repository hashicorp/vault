package pgtype

import (
	"database/sql/driver"
	"encoding/hex"

	"github.com/pkg/errors"
)

type Bytea struct {
	Bytes  []byte
	Status Status
}

func (dst *Bytea) Set(src interface{}) error {
	if src == nil {
		*dst = Bytea{Status: Null}
		return nil
	}

	switch value := src.(type) {
	case []byte:
		if value != nil {
			*dst = Bytea{Bytes: value, Status: Present}
		} else {
			*dst = Bytea{Status: Null}
		}
	default:
		if originalSrc, ok := underlyingBytesType(src); ok {
			return dst.Set(originalSrc)
		}
		return errors.Errorf("cannot convert %v to Bytea", value)
	}

	return nil
}

func (dst *Bytea) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Bytes
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Bytea) AssignTo(dst interface{}) error {
	switch src.Status {
	case Present:
		switch v := dst.(type) {
		case *[]byte:
			buf := make([]byte, len(src.Bytes))
			copy(buf, src.Bytes)
			*v = buf
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

// DecodeText only supports the hex format. This has been the default since
// PostgreSQL 9.0.
func (dst *Bytea) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Bytea{Status: Null}
		return nil
	}

	if len(src) < 2 || src[0] != '\\' || src[1] != 'x' {
		return errors.Errorf("invalid hex format")
	}

	buf := make([]byte, (len(src)-2)/2)
	_, err := hex.Decode(buf, src[2:])
	if err != nil {
		return err
	}

	*dst = Bytea{Bytes: buf, Status: Present}
	return nil
}

func (dst *Bytea) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Bytea{Status: Null}
		return nil
	}

	*dst = Bytea{Bytes: src, Status: Present}
	return nil
}

func (src *Bytea) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = append(buf, `\x`...)
	buf = append(buf, hex.EncodeToString(src.Bytes)...)
	return buf, nil
}

func (src *Bytea) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return append(buf, src.Bytes...), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Bytea) Scan(src interface{}) error {
	if src == nil {
		*dst = Bytea{Status: Null}
		return nil
	}

	switch src := src.(type) {
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		buf := make([]byte, len(src))
		copy(buf, src)
		*dst = Bytea{Bytes: buf, Status: Present}
		return nil
	}

	return errors.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *Bytea) Value() (driver.Value, error) {
	switch src.Status {
	case Present:
		return src.Bytes, nil
	case Null:
		return nil, nil
	default:
		return nil, errUndefined
	}
}
