package gore

import (
	"errors"
)

var (
	// ErrNumberFormat is returned when formatting number fails.
	ErrNumberFormat = errors.New("number format error")
)

// FixInt represents a fixed size int64 number
type FixInt int64

// Bytes converts a FixInt to a byte array
func (x FixInt) Bytes() []byte {
	b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	y := int64(x)
	b[0] = byte((y >> 56) & 0xFF)
	b[1] = byte((y >> 48) & 0xFF)
	b[2] = byte((y >> 40) & 0xFF)
	b[3] = byte((y >> 32) & 0xFF)
	b[4] = byte((y >> 24) & 0xFF)
	b[5] = byte((y >> 16) & 0xFF)
	b[6] = byte((y >> 8) & 0xFF)
	b[7] = byte(y & 0xFF)
	return b
}

// ToFixInt converts a fixed size byte array to a int64
func ToFixInt(b []byte) (int64, error) {
	if len(b) != 8 {
		return 0, ErrNumberFormat
	}
	var x int64
	for i := range b {
		x = (x << 8) + int64(b[i]&0xFF)
	}
	return x, nil
}

// VarInt represents a base-128 int64 number
type VarInt int64

// Bytes converts a VarInt to a byte array
func (x VarInt) Bytes() []byte {
	b := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0} // 9.14 bytes are needed
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	pos := 0
	for {
		b[pos] = byte((ux & 0x7F) | 0x80)
		pos++
		if ux >= 128 {
			ux >>= 7
		} else {
			break
		}
	}
	b[pos-1] &= 0x7F
	return b[0:pos]
}

// ToVarInt converts a base-128 byte array to a int64
func ToVarInt(b []byte) (int64, error) {
	if len(b) < 1 {
		return 0, ErrNumberFormat
	}
	var ux uint64
	for i := range b {
		ux += uint64((b[i] & 0x7F)) << uint(7*i)
		if b[i] < 128 {
			break
		}
	}
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, nil
}
