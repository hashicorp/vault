package pgtype

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"

	"github.com/jackc/pgio"
)

type Varbit struct {
	Bytes  []byte
	Len    int32 // Number of bits
	Status Status
}

func (dst *Varbit) Set(src interface{}) error {
	return fmt.Errorf("cannot convert %v to Varbit", src)
}

func (dst Varbit) Get() interface{} {
	switch dst.Status {
	case Present:
		return dst
	case Null:
		return nil
	default:
		return dst.Status
	}
}

func (src *Varbit) AssignTo(dst interface{}) error {
	return fmt.Errorf("cannot assign %v to %T", src, dst)
}

func (dst *Varbit) DecodeText(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Varbit{Status: Null}
		return nil
	}

	bitLen := len(src)
	byteLen := bitLen / 8
	if bitLen%8 > 0 {
		byteLen++
	}
	buf := make([]byte, byteLen)

	for i, b := range src {
		if b == '1' {
			byteIdx := i / 8
			bitIdx := uint(i % 8)
			buf[byteIdx] = buf[byteIdx] | (128 >> bitIdx)
		}
	}

	*dst = Varbit{Bytes: buf, Len: int32(bitLen), Status: Present}
	return nil
}

func (dst *Varbit) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Varbit{Status: Null}
		return nil
	}

	if len(src) < 4 {
		return fmt.Errorf("invalid length for varbit: %v", len(src))
	}

	bitLen := int32(binary.BigEndian.Uint32(src))
	rp := 4

	*dst = Varbit{Bytes: src[rp:], Len: bitLen, Status: Present}
	return nil
}

func (src Varbit) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	for i := int32(0); i < src.Len; i++ {
		byteIdx := i / 8
		bitMask := byte(128 >> byte(i%8))
		char := byte('0')
		if src.Bytes[byteIdx]&bitMask > 0 {
			char = '1'
		}
		buf = append(buf, char)
	}

	return buf, nil
}

func (src Varbit) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = pgio.AppendInt32(buf, src.Len)
	return append(buf, src.Bytes...), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Varbit) Scan(src interface{}) error {
	if src == nil {
		*dst = Varbit{Status: Null}
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
func (src Varbit) Value() (driver.Value, error) {
	return EncodeValueText(src)
}
