package memd

import (
	"errors"
)

// AppendULEB128_32 appends a 32-bit number encoded as ULEB128 to a byte slice
func AppendULEB128_32(b []byte, v uint32) []byte {
	for {
		c := uint8(v & 0x7f)
		v >>= 7
		if v != 0 {
			c |= 0x80
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}
	return b
}

// DecodeULEB128_32 decodes a ULEB128 encoded number into a uint32
func DecodeULEB128_32(b []byte) (uint32, int, error) {
	if len(b) == 0 {
		return 0, 0, errors.New("no data provided")
	}
	var u uint64
	var n int
	for i := 0; ; i++ {
		if i >= len(b) {
			return 0, 0, errors.New("encoded number is longer than provided data")
		}
		if i*7 > 32 {
			// oversize and then break to get caught below
			u = 0xffffffffffffffff
			break
		}

		u |= uint64(b[i]&0x7f) << (i * 7)

		if b[i]&0x80 == 0 {
			n = i + 1
			break
		}
	}

	if u > 0xffffffff {
		return 0, 0, errors.New("encoded data is longer than 32 bits")
	}

	return uint32(u), n, nil
}
