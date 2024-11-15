package ber

import (
	"errors"
	"fmt"
	"io"
)

func readLength(reader io.Reader) (length int, read int, err error) {
	// length byte
	b, err := readByte(reader)
	if err != nil {
		if Debug {
			fmt.Printf("error reading length byte: %v\n", err)
		}
		return 0, 0, unexpectedEOF(err)
	}
	read++

	switch {
	case b == 0xFF:
		// Invalid 0xFF (x.600, 8.1.3.5.c)
		return 0, read, errors.New("invalid length byte 0xff")

	case b == LengthLongFormBitmask:
		// Indefinite form, we have to decode packets until we encounter an EOC packet (x.600, 8.1.3.6)
		length = LengthIndefinite

	case b&LengthLongFormBitmask == 0:
		// Short definite form, extract the length from the bottom 7 bits (x.600, 8.1.3.4)
		length = int(b) & LengthValueBitmask

	case b&LengthLongFormBitmask != 0:
		// Long definite form, extract the number of length bytes to follow from the bottom 7 bits (x.600, 8.1.3.5.b)
		lengthBytes := int(b) & LengthValueBitmask
		// Protect against overflow
		// TODO: support big int length?
		if lengthBytes > 8 {
			return 0, read, errors.New("long-form length overflow")
		}

		// Accumulate into a 64-bit variable
		var length64 int64
		for i := 0; i < lengthBytes; i++ {
			b, err = readByte(reader)
			if err != nil {
				if Debug {
					fmt.Printf("error reading long-form length byte %d: %v\n", i, err)
				}
				return 0, read, unexpectedEOF(err)
			}
			read++

			// x.600, 8.1.3.5
			length64 <<= 8
			length64 |= int64(b)
		}

		// Cast to a platform-specific integer
		length = int(length64)
		// Ensure we didn't overflow
		if int64(length) != length64 {
			return 0, read, errors.New("long-form length overflow")
		}

	default:
		return 0, read, errors.New("invalid length byte")
	}

	return length, read, nil
}

func encodeLength(length int) []byte {
	lengthBytes := encodeUnsignedInteger(uint64(length))
	if length > 127 || len(lengthBytes) > 1 {
		longFormBytes := []byte{LengthLongFormBitmask | byte(len(lengthBytes))}
		longFormBytes = append(longFormBytes, lengthBytes...)
		lengthBytes = longFormBytes
	}
	return lengthBytes
}
