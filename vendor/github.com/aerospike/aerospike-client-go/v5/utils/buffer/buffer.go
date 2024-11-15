// Copyright 2014-2021 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package buffer

import (
	"encoding/binary"
	"fmt"
	"math"
)

const (
	// SizeOfInt32 defines the size of int32
	SizeOfInt32 = uintptr(4)
	// SizeOfInt64 defines the size of int64
	SizeOfInt64 = uintptr(8)

	uint64sz = int(8)
	uint32sz = int(4)
	uint16sz = int(2)

	float32sz = int(4)
	float64sz = int(8)
)

// SizeOfInt defines the size of native int
var SizeOfInt uintptr

// Arch64Bits defines if the system is 64 bits
var Arch64Bits bool

// Arch32Bits defines if the system is 32 bits
var Arch32Bits bool

func init() {
	if 0 == ^uint(0xffffffff) {
		SizeOfInt = 4
	} else {
		SizeOfInt = 8
	}
	Arch64Bits = (SizeOfInt == SizeOfInt64)
	Arch32Bits = (SizeOfInt == SizeOfInt32)
}

// BytesToHexString converts a byte slice into a hex string
func BytesToHexString(buf []byte) string {
	hlist := make([]byte, 3*len(buf))

	for i := range buf {
		hex := fmt.Sprintf("%02x ", buf[i])
		idx := i * 3
		copy(hlist[idx:], hex)
	}
	return string(hlist)
}

// LittleBytesToInt32 converts a slice into int32; only maximum of 4 bytes will be used
func LittleBytesToInt32(buf []byte, offset int) int32 {
	l := len(buf[offset:])
	if l > uint32sz {
		l = uint32sz
	}
	r := int32(binary.LittleEndian.Uint32(buf[offset : offset+l]))
	return r
}

// BytesToInt64 converts a slice into int64; only maximum of 8 bytes will be used
func BytesToInt64(buf []byte, offset int) int64 {
	l := len(buf[offset:])
	if l > uint64sz {
		l = uint64sz
	}
	r := int64(binary.BigEndian.Uint64(buf[offset : offset+l]))
	return r
}

// VarBytesToInt64 will convert a 8, 4 or 2 byte slice into an int64
func VarBytesToInt64(buf []byte, offset int, len int) int64 {
	if len == 8 {
		return BytesToInt64(buf, offset)
	} else if len == 4 {
		return int64(BytesToInt32(buf, offset))
	} else if len == 2 {
		return int64(BytesToInt16(buf, offset))
	}

	val := int64(0)
	for i := 0; i < len; i++ {
		val <<= 8
		val |= int64(buf[offset+i] & 0xFF)
	}
	return val
}

// BytesToInt32 converts a slice into int32; only maximum of 4 bytes will be used
func BytesToInt32(buf []byte, offset int) int32 {
	return int32(binary.BigEndian.Uint32(buf[offset : offset+uint32sz]))
}

// BytesToUint32 converts a slice into uint32; only maximum of 4 bytes will be used
func BytesToUint32(buf []byte, offset int) uint32 {
	return binary.BigEndian.Uint32(buf[offset : offset+uint32sz])
}

// BytesToInt16 converts a slice of bytes to an int16
func BytesToInt16(buf []byte, offset int) int16 {
	return int16(binary.BigEndian.Uint16(buf[offset : offset+uint16sz]))
}

// BytesToUint16 converts a byte slice to a uint16
func BytesToUint16(buf []byte, offset int) uint16 {
	return binary.BigEndian.Uint16(buf[offset : offset+uint16sz])
}

// BytesToFloat32 converts a byte slice to a float32
func BytesToFloat32(buf []byte, offset int) float32 {
	bits := binary.BigEndian.Uint32(buf[offset : offset+float32sz])
	return math.Float32frombits(bits)
}

// BytesToFloat64 converts a byte slice to a float64
func BytesToFloat64(buf []byte, offset int) float64 {
	bits := binary.BigEndian.Uint64(buf[offset : offset+float64sz])
	return math.Float64frombits(bits)
}

// BytesToBool converts a byte slice to a bool
func BytesToBool(buf []byte, offset, length int) bool {
	if length <= 0 {
		return false
	}
	return buf[offset] != 0
}
