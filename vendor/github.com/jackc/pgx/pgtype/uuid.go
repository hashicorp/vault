package pgtype

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

type UUID struct {
	Bytes  [16]byte
	Status Status
}

func (dst *UUID) Set(src interface{}) error {
	if src == nil {
		*dst = UUID{Status: Null}
		return nil
	}

	switch value := src.(type) {
	case [16]byte:
		*dst = UUID{Bytes: value, Status: Present}
	case []byte:
		if value != nil {
			if len(value) != 16 {
				return errors.Errorf("[]byte must be 16 bytes to convert to UUID: %d", len(value))
			}
			*dst = UUID{Status: Present}
			copy(dst.Bytes[:], value)
		} else {
			*dst = UUID{Status: Null}
		}
	case string:
		uuid, err := parseUUID(value)
		if err != nil {
			return err
		}
		*dst = UUID{Bytes: uuid, Status: Present}
	default:
		if originalSrc, ok := underlyingPtrType(src); ok {
			return dst.Set(originalSrc)
		}
		return errors.Errorf("cannot convert %v to UUID", value)
	}

	return nil
}

func (dst *UUID) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst.Bytes
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *UUID) AssignTo(dst interface{}) error {
	switch src.Status {
	case Present:
		switch v := dst.(type) {
		case *[16]byte:
			*v = src.Bytes
			return nil
		case *[]byte:
			*v = make([]byte, 16)
			copy(*v, src.Bytes[:])
			return nil
		case *string:
			*v = encodeUUID(src.Bytes)
			return nil
		default:
			if nextDst, retry := GetAssignToDstType(v); retry {
				return src.AssignTo(nextDst)
			}
		}
	case Null:
		return NullAssignTo(dst)
	}

	return errors.Errorf("cannot assign %v into %T", src, dst)
}

// parseUUID converts a string UUID in standard form to a byte array.
func parseUUID(src string) (dst [16]byte, err error) {
	src = src[0:8] + src[9:13] + src[14:18] + src[19:23] + src[24:]
	buf, err := hex.DecodeString(src)
	if err != nil {
		return dst, err
	}

	copy(dst[:], buf)
	return dst, err
}

// encodeUUID converts a uuid byte array to UUID standard string form.
func encodeUUID(src [16]byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", src[0:4], src[4:6], src[6:8], src[8:10], src[10:16])
}

func (dst *UUID) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = UUID{Status: Null}
		return nil
	}

	if len(src) != 36 {
		return errors.Errorf("invalid length for UUID: %v", len(src))
	}

	buf, err := parseUUID(string(src))
	if err != nil {
		return err
	}

	*dst = UUID{Bytes: buf, Status: Present}
	return nil
}

func (dst *UUID) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = UUID{Status: Null}
		return nil
	}

	if len(src) != 16 {
		return errors.Errorf("invalid length for UUID: %v", len(src))
	}

	*dst = UUID{Status: Present}
	copy(dst.Bytes[:], src)
	return nil
}

func (src *UUID) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return append(buf, encodeUUID(src.Bytes)...), nil
}

func (src *UUID) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return append(buf, src.Bytes[:]...), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *UUID) Scan(src interface{}) error {
	if src == nil {
		*dst = UUID{Status: Null}
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

	return errors.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *UUID) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
