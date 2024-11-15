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

package aerospike

import (
	"encoding/binary"
	"math"

	"github.com/aerospike/aerospike-client-go/v5/pkg/ripemd160"
	"github.com/aerospike/aerospike-client-go/v5/types"
)

type keyWriter struct {
	buffer [8]byte
	hash   ripemd160.Digest
}

// WriteInt64 writes a int64 to the key
func (vb *keyWriter) WriteInt64(num int64) int {
	return vb.WriteUint64(uint64(num))
}

// WriteUint64 writes a uint64 to the key
func (vb *keyWriter) WriteUint64(num uint64) int {
	binary.BigEndian.PutUint64(vb.buffer[:8], num)
	vb.hash.Write(vb.buffer[:8])
	return 8
}

// WriteInt32 writes a int32 to the key
func (vb *keyWriter) WriteInt32(num int32) int {
	return vb.WriteUint32(uint32(num))
}

// WriteUint32 writes a uint32 to the key
func (vb *keyWriter) WriteUint32(num uint32) int {
	binary.BigEndian.PutUint32(vb.buffer[:4], num)
	vb.hash.Write(vb.buffer[:4])
	return 4
}

// WriteInt16 writes a int16 to the key
func (vb *keyWriter) WriteInt16(num int16) int {
	return vb.WriteUint16(uint16(num))
}

// WriteUint16 writes a uint16 to the key
func (vb *keyWriter) WriteUint16(num uint16) int {
	binary.BigEndian.PutUint16(vb.buffer[:2], num)
	vb.hash.Write(vb.buffer[:2])
	return 2
}

// WriteFloat32 writes a float32 to the key
func (vb *keyWriter) WriteFloat32(float float32) int {
	bits := math.Float32bits(float)
	binary.BigEndian.PutUint32(vb.buffer[:4], bits)
	vb.hash.Write(vb.buffer[:4])
	return 4
}

// WriteFloat64 writes a float64 to the key
func (vb *keyWriter) WriteFloat64(float float64) int {
	bits := math.Float64bits(float)
	binary.BigEndian.PutUint64(vb.buffer[:8], bits)
	vb.hash.Write(vb.buffer[:8])
	return 8
}

// WriteBool writes a bool to the key
func (vb *keyWriter) WriteBool(b bool) int {
	if b {
		vb.hash.Write([]byte{1})
	} else {
		vb.hash.Write([]byte{0})
	}
	return 1
}

// WriteByte writes a byte to the key
func (vb *keyWriter) WriteByte(b byte) {
	vb.hash.Write([]byte{b})
}

// WriteString writes a string to the key
func (vb *keyWriter) WriteString(s string) (int, Error) {
	// To avoid allocating memory, write the strings in small chunks
	l := len(s)
	const size = 128
	b := [size]byte{}
	cnt := 0
	sz := 0
	for i := 0; i < l; i++ {
		b[cnt] = s[i]
		cnt++

		if cnt == size {
			n, err := vb.Write(b[:])
			if err != nil {
				return sz + n, err
			}
			sz += n
			cnt = 0
		}
	}

	if cnt > 0 {
		n, err := vb.Write(b[:cnt])
		if err != nil {
			return sz + n, err
		}
	}

	return len(s), nil
}

func (vb *keyWriter) Write(b []byte) (int, Error) {
	n, err := vb.hash.Write(b)
	if err != nil {
		return n, newCommonError(err)
	}
	return n, nil
}

func (vb *keyWriter) writeKey(val Value) Error {
	switch v := val.(type) {
	case IntegerValue:
		vb.WriteInt64(int64(v))
		return nil
	case LongValue:
		vb.WriteInt64(int64(v))
		return nil
	case StringValue:
		vb.WriteString(string(v))
		return nil
	case BytesValue:
		vb.Write(v)
		return nil
	}

	return newError(types.PARAMETER_ERROR, "Key Generation Error. Value not supported: "+val.String())
}
