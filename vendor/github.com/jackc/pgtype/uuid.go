package pgtype

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
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
	case interface{ Get() interface{} }:
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	case fmt.Stringer:
		value2 := value.String()
		return dst.Set(value2)
	case [16]byte:
		*dst = UUID{Bytes: value, Status: Present}
	case []byte:
		if value != nil {
			if len(value) != 16 {
				return fmt.Errorf("[]byte must be 16 bytes to convert to UUID: %d", len(value))
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
	case *string:
		if value == nil {
			*dst = UUID{Status: Null}
		} else {
			return dst.Set(*value)
		}
	default:
		if originalSrc, ok := underlyingUUIDType(src); ok {
			return dst.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to UUID", value)
	}

	return nil
}

func (dst UUID) Get() interface{} {
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
			*v = string(encodeUUID(src.Bytes))
			return nil
		default:
			if nextDst, retry := GetAssignToDstType(v); retry {
				return src.AssignTo(nextDst)
			}
		}
	case Null:
		return NullAssignTo(dst)
	}

	return fmt.Errorf("cannot assign %v into %T", src, dst)
}

// parseUUID converts a string UUID in standard form to a byte array.
func parseUUID(src string) (dst [16]byte, err error) {
	var uuidBuf [32]byte
	srcBuf := uuidBuf[:]

	switch len(src) {
	case 36:
		copy(srcBuf[0:8], src[:8])
		copy(srcBuf[8:12], src[9:13])
		copy(srcBuf[12:16], src[14:18])
		copy(srcBuf[16:20], src[19:23])
		copy(srcBuf[20:], src[24:])
	case 32:
		// dashes already stripped, assume valid
		copy(srcBuf, src)

	default:
		// assume invalid.
		return dst, fmt.Errorf("cannot parse UUID %v", src)
	}

	_, err = hex.Decode(dst[:], srcBuf)
	if err != nil {
		return dst, err
	}
	return dst, err
}

// encodeUUID converts a uuid byte array to UUID standard string form.
func encodeUUID(src [16]byte) (dst []byte) {
	var buf [36]byte
	dst = buf[:]

	hex.Encode(dst, src[:4])
	buf[8] = '-'
	hex.Encode(dst[9:13], src[4:6])
	buf[13] = '-'
	hex.Encode(dst[14:18], src[6:8])
	buf[18] = '-'
	hex.Encode(dst[19:23], src[8:10])
	buf[23] = '-'
	hex.Encode(dst[24:], src[10:])

	return
}

func (dst *UUID) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = UUID{Status: Null}
		return nil
	}

	if len(src) != 36 {
		return fmt.Errorf("invalid length for UUID: %v", len(src))
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
		return fmt.Errorf("invalid length for UUID: %v", len(src))
	}

	*dst = UUID{Status: Present}
	copy(dst.Bytes[:], src)
	return nil
}

func (src UUID) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	return append(buf, encodeUUID(src.Bytes)...), nil
}

func (src UUID) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
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

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src UUID) Value() (driver.Value, error) {
	return EncodeValueText(src)
}

func (src UUID) MarshalJSON() ([]byte, error) {
	switch src.Status {
	case Present:
		var buff bytes.Buffer
		buff.WriteByte('"')
		buff.Write(encodeUUID(src.Bytes))
		buff.WriteByte('"')
		return buff.Bytes(), nil
	case Null:
		return []byte("null"), nil
	case Undefined:
		return nil, errUndefined
	}
	return nil, errBadStatus
}

func (dst *UUID) UnmarshalJSON(src []byte) error {
	if bytes.Equal(src, []byte("null")) {
		return dst.Set(nil)
	}
	if len(src) != 38 {
		return fmt.Errorf("invalid length for UUID: %v", len(src))
	}
	return dst.Set(string(src[1 : len(src)-1]))
}
