// Copyright 2013-2020 Aerospike, Inc.
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

	"github.com/aerospike/aerospike-client-go/pkg/ripemd160"
	. "github.com/aerospike/aerospike-client-go/types"
)

type keyWriter struct {
	buffer [8]byte
	offset int
	hash   ripemd160.Digest
}

// Int64ToBytes converts an int64 into slice of Bytes.
func (vb *keyWriter) WriteInt64(num int64) (int, error) {
	return vb.WriteUint64(uint64(num))
}

// Uint64ToBytes converts an uint64 into slice of Bytes.
func (vb *keyWriter) WriteUint64(num uint64) (int, error) {
	binary.BigEndian.PutUint64(vb.buffer[:8], num)
	vb.hash.Write(vb.buffer[:8])
	return 8, nil
}

// Int32ToBytes converts an int32 to a byte slice of size 4
func (vb *keyWriter) WriteInt32(num int32) (int, error) {
	return vb.WriteUint32(uint32(num))
}

// Uint32ToBytes converts an uint32 to a byte slice of size 4
func (vb *keyWriter) WriteUint32(num uint32) (int, error) {
	binary.BigEndian.PutUint32(vb.buffer[:4], num)
	vb.hash.Write(vb.buffer[:4])
	return 4, nil
}

// Int16ToBytes converts an int16 to slice of bytes
func (vb *keyWriter) WriteInt16(num int16) (int, error) {
	return vb.WriteUint16(uint16(num))
}

// UInt16ToBytes converts an iuint16 to slice of bytes
func (vb *keyWriter) WriteUint16(num uint16) (int, error) {
	binary.BigEndian.PutUint16(vb.buffer[:2], num)
	vb.hash.Write(vb.buffer[:2])
	return 2, nil
}

func (vb *keyWriter) WriteFloat32(float float32) (int, error) {
	bits := math.Float32bits(float)
	binary.BigEndian.PutUint32(vb.buffer[:4], bits)
	vb.hash.Write(vb.buffer[:4])
	return 4, nil
}

func (vb *keyWriter) WriteFloat64(float float64) (int, error) {
	bits := math.Float64bits(float)
	binary.BigEndian.PutUint64(vb.buffer[:8], bits)
	vb.hash.Write(vb.buffer[:8])
	return 8, nil
}

func (vb *keyWriter) WriteByte(b byte) error {
	_, err := vb.hash.Write([]byte{b})
	return err
}

func (vb *keyWriter) WriteString(s string) (int, error) {
	// To avoid allocating memory, write the strings in small chunks
	l := len(s)
	const size = 128
	b := [size]byte{}
	cnt := 0
	for i := 0; i < l; i++ {
		b[cnt] = s[i]
		cnt++

		if cnt == size {
			vb.Write(b[:])
			cnt = 0
		}
	}

	if cnt > 0 {
		vb.Write(b[:cnt])
	}

	return len(s), nil
}

func (vb *keyWriter) Write(b []byte) (int, error) {
	vb.hash.Write(b)
	return len(b), nil
}

func (vb *keyWriter) writeKey(val Value) error {
	switch v := val.(type) {
	case IntegerValue:
		vb.WriteInt64(int64(v))
		return nil
	case LongValue:
		vb.WriteInt64(int64(v))
		return nil
	case FloatValue:
		vb.WriteFloat64(float64(v))
		return nil
	case StringValue:
		vb.WriteString(string(v))
		return nil
	case ListValue:
		v.pack(vb)
		return nil
	case *ListValue:
		v.pack(vb)
		return nil
	case *ListerValue:
		v.pack(vb)
		return nil
	case ValueArray:
		v.pack(vb)
		return nil
	case *ValueArray:
		v.pack(vb)
		return nil
	case BytesValue:
		vb.Write(v)
		return nil
	}

	return NewAerospikeError(PARAMETER_ERROR, "Key Generation Error. Value not supported: "+val.String())
}
