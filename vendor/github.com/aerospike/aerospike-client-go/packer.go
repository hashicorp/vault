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
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"time"

	ParticleType "github.com/aerospike/aerospike-client-go/internal/particle_type"
	. "github.com/aerospike/aerospike-client-go/types"
	Buffer "github.com/aerospike/aerospike-client-go/utils/buffer"
)

var packObjectReflect func(BufferEx, interface{}, bool) (int, error)

func packIfcList(cmd BufferEx, list []interface{}) (int, error) {
	size := 0
	n, err := packArrayBegin(cmd, len(list))
	if err != nil {
		return n, err
	}
	size += n

	for i := range list {
		n, err := packObject(cmd, list[i], false)
		if err != nil {
			return 0, err
		}
		size += n
	}

	return size, err
}

// PackList packs any slice that implement the ListIter interface
func PackList(cmd BufferEx, list ListIter) (int, error) {
	return packList(cmd, list)
}

func packList(cmd BufferEx, list ListIter) (int, error) {
	size := 0
	n, err := packArrayBegin(cmd, list.Len())
	if err != nil {
		return n, err
	}
	size += n

	n, err = list.PackList(cmd)
	return size + n, err
}

func packValueArray(cmd BufferEx, list ValueArray) (int, error) {
	size := 0
	n, err := packArrayBegin(cmd, len(list))
	if err != nil {
		return n, err
	}
	size += n

	for i := range list {
		n, err := list[i].pack(cmd)
		if err != nil {
			return 0, err
		}
		size += n
	}

	return size, err
}

func packArrayBegin(cmd BufferEx, size int) (int, error) {
	if size < 16 {
		return packAByte(cmd, 0x90|byte(size))
	} else if size <= math.MaxUint16 {
		return packShort(cmd, 0xdc, int16(size))
	} else {
		return packInt(cmd, 0xdd, int32(size))
	}
}

func packIfcMap(cmd BufferEx, theMap map[interface{}]interface{}) (int, error) {
	size := 0
	n, err := packMapBegin(cmd, len(theMap))
	if err != nil {
		return n, err
	}
	size += n

	for k, v := range theMap {
		n, err := packObject(cmd, k, true)
		if err != nil {
			return 0, err
		}
		size += n
		n, err = packObject(cmd, v, false)
		if err != nil {
			return 0, err
		}
		size += n
	}

	return size, err
}

// PackJson packs json data
func PackJson(cmd BufferEx, theMap map[string]interface{}) (int, error) {
	return packJsonMap(cmd, theMap)
}

func packJsonMap(cmd BufferEx, theMap map[string]interface{}) (int, error) {
	size := 0
	n, err := packMapBegin(cmd, len(theMap))
	if err != nil {
		return n, err
	}
	size += n

	for k, v := range theMap {
		n, err := packString(cmd, k)
		if err != nil {
			return 0, err
		}
		size += n
		n, err = packObject(cmd, v, false)
		if err != nil {
			return 0, err
		}
		size += n
	}

	return size, err
}

// PackMap packs any map that implements the MapIter interface
func PackMap(cmd BufferEx, theMap MapIter) (int, error) {
	return packMap(cmd, theMap)
}

func packMap(cmd BufferEx, theMap MapIter) (int, error) {
	size := 0
	n, err := packMapBegin(cmd, theMap.Len())
	if err != nil {
		return n, err
	}
	size += n

	n, err = theMap.PackMap(cmd)
	return size + n, err
}

func packMapBegin(cmd BufferEx, size int) (int, error) {
	if size < 16 {
		return packAByte(cmd, 0x80|byte(size))
	} else if size <= math.MaxUint16 {
		return packShort(cmd, 0xde, int16(size))
	} else {
		return packInt(cmd, 0xdf, int32(size))
	}
}

// PackBytes backs a byte array
func PackBytes(cmd BufferEx, b []byte) (int, error) {
	return packBytes(cmd, b)
}

func packBytes(cmd BufferEx, b []byte) (int, error) {
	size := 0
	n, err := packByteArrayBegin(cmd, len(b)+1)
	if err != nil {
		return n, err
	}
	size += n

	n, err = packAByte(cmd, ParticleType.BLOB)
	if err != nil {
		return size + n, err
	}
	size += n

	n, err = packByteArray(cmd, b)
	if err != nil {
		return size + n, err
	}
	size += n

	return size, nil
}

func packByteArrayBegin(cmd BufferEx, length int) (int, error) {
	// Use string header codes for byte arrays.
	return packStringBegin(cmd, length)
}

func packObject(cmd BufferEx, obj interface{}, mapKey bool) (int, error) {
	switch v := obj.(type) {
	case Value:
		return v.pack(cmd)
	case []Value:
		return ValueArray(v).pack(cmd)
	case string:
		return packString(cmd, v)
	case []byte:
		return packBytes(cmd, obj.([]byte))
	case int8:
		return packAInt(cmd, int(v))
	case uint8:
		return packAInt(cmd, int(v))
	case int16:
		return packAInt(cmd, int(v))
	case uint16:
		return packAInt(cmd, int(v))
	case int32:
		return packAInt(cmd, int(v))
	case uint32:
		return packAInt(cmd, int(v))
	case int:
		if Buffer.Arch32Bits {
			return packAInt(cmd, v)
		}
		return packAInt64(cmd, int64(v))
	case uint:
		if Buffer.Arch32Bits {
			return packAInt(cmd, int(v))
		}
		return packAUInt64(cmd, uint64(v))
	case int64:
		return packAInt64(cmd, v)
	case uint64:
		return packAUInt64(cmd, v)
	case time.Time:
		return packAInt64(cmd, v.UnixNano())
	case nil:
		return packNil(cmd)
	case bool:
		return packBool(cmd, v)
	case float32:
		return packFloat32(cmd, v)
	case float64:
		return packFloat64(cmd, v)
	case struct{}:
		if mapKey {
			return 0, NewAerospikeError(SERIALIZE_ERROR, fmt.Sprintf("Maps, Slices, and bounded arrays other than Bounded Byte Arrays are not supported as Map keys. Value: %#v", v))
		}
		return packIfcMap(cmd, map[interface{}]interface{}{})
	case []interface{}:
		if mapKey {
			return 0, NewAerospikeError(SERIALIZE_ERROR, fmt.Sprintf("Maps, Slices, and bounded arrays other than Bounded Byte Arrays are not supported as Map keys. Value: %#v", v))
		}
		return packIfcList(cmd, v)
	case map[interface{}]interface{}:
		if mapKey {
			return 0, NewAerospikeError(SERIALIZE_ERROR, fmt.Sprintf("Maps, Slices, and bounded arrays other than Bounded Byte Arrays are not supported as Map keys. Value: %#v", v))
		}
		return packIfcMap(cmd, v)
	case ListIter:
		if mapKey {
			return 0, NewAerospikeError(SERIALIZE_ERROR, fmt.Sprintf("Maps, Slices, and bounded arrays other than Bounded Byte Arrays are not supported as Map keys. Value: %#v", v))
		}
		return packList(cmd, obj.(ListIter))
	case MapIter:
		if mapKey {
			return 0, NewAerospikeError(SERIALIZE_ERROR, fmt.Sprintf("Maps, Slices, and bounded arrays other than Bounded Byte Arrays are not supported as Map keys. Value: %#v", v))
		}
		return packMap(cmd, obj.(MapIter))
	}

	// try to see if the object is convertible to a concrete value.
	// This will be faster and much more memory efficient than reflection.
	if v := tryConcreteValue(obj); v != nil {
		return v.pack(cmd)
	}

	if packObjectReflect != nil {
		return packObjectReflect(cmd, obj, mapKey)
	}

	return 0, NewAerospikeError(SERIALIZE_ERROR, fmt.Sprintf("Type `%v (%s)` not supported to pack. ", obj, reflect.TypeOf(obj).String()))
}

func packAUInt64(cmd BufferEx, val uint64) (int, error) {
	return packUInt64(cmd, val)
}

func packAInt64(cmd BufferEx, val int64) (int, error) {
	if val >= 0 {
		if val < 128 {
			return packAByte(cmd, byte(val))
		}

		if val <= math.MaxUint8 {
			return packByte(cmd, 0xcc, byte(val))
		}

		if val <= math.MaxUint16 {
			return packShort(cmd, 0xcd, int16(val))
		}

		if val <= math.MaxUint32 {
			return packInt(cmd, 0xce, int32(val))
		}
		return packInt64(cmd, 0xd3, val)
	}

	if val >= -32 {
		return packAByte(cmd, 0xe0|(byte(val)+32))
	}

	if val >= math.MinInt8 {
		return packByte(cmd, 0xd0, byte(val))
	}

	if val >= math.MinInt16 {
		return packShort(cmd, 0xd1, int16(val))
	}

	if val >= math.MinInt32 {
		return packInt(cmd, 0xd2, int32(val))
	}
	return packInt64(cmd, 0xd3, val)
}

// PackInt64 packs an int64
func PackInt64(cmd BufferEx, val int64) (int, error) {
	return packAInt64(cmd, val)
}

func packAInt(cmd BufferEx, val int) (int, error) {
	return packAInt64(cmd, int64(val))
}

// PackString packs a string
func PackString(cmd BufferEx, val string) (int, error) {
	return packString(cmd, val)
}

func packStringBegin(cmd BufferEx, size int) (int, error) {
	if size < 32 {
		return packAByte(cmd, 0xa0|byte(size))
	} else if size < 256 {
		return packByte(cmd, 0xd9, byte(size))
	} else if size < 65536 {
		return packShort(cmd, 0xda, int16(size))
	}
	return packInt(cmd, 0xdb, int32(size))
}

func packString(cmd BufferEx, val string) (int, error) {
	size := 0
	slen := len(val) + 1
	n, err := packStringBegin(cmd, slen)
	if err != nil {
		return n, err
	}
	size += n

	if cmd != nil {
		n, err = 1, cmd.WriteByte(byte(ParticleType.STRING))
		if err != nil {
			return size + n, err
		}
		size += n

		n, err = cmd.WriteString(val)
		if err != nil {
			return size + n, err
		}
		size += n
	} else {
		size += 1 + len(val)
	}

	return size, nil
}

func packRawString(cmd BufferEx, val string) (int, error) {
	size := 0
	slen := len(val)
	n, err := packStringBegin(cmd, slen)
	if err != nil {
		return n, err
	}
	size += n

	if cmd != nil {
		n, err = cmd.WriteString(val)
		if err != nil {
			return size + n, err
		}
		size += n
	} else {
		size += len(val)
	}

	return size, nil
}

func packGeoJson(cmd BufferEx, val string) (int, error) {
	size := 0
	slen := len(val) + 1
	n, err := packByteArrayBegin(cmd, slen)
	if err != nil {
		return n, err
	}
	size += n

	if cmd != nil {
		n, err = 1, cmd.WriteByte(byte(ParticleType.GEOJSON))
		if err != nil {
			return size + n, err
		}
		size += n

		n, err = cmd.WriteString(val)
		if err != nil {
			return size + n, err
		}
		size += n
	} else {
		size += 1 + len(val)
	}

	return size, nil
}

func packByteArray(cmd BufferEx, src []byte) (int, error) {
	if cmd != nil {
		return cmd.Write(src)
	}
	return len(src), nil
}

func packInt64(cmd BufferEx, valType int, val int64) (int, error) {
	if cmd != nil {
		size, err := 1, cmd.WriteByte(byte(valType))
		if err != nil {
			return size, err
		}

		n, err := cmd.WriteInt64(val)
		return size + n, err
	}
	return 1 + 8, nil
}

// PackUInt64 packs a uint64
func PackUInt64(cmd BufferEx, val uint64) (int, error) {
	return packUInt64(cmd, val)
}

func packUInt64(cmd BufferEx, val uint64) (int, error) {
	if cmd != nil {
		size, err := 1, cmd.WriteByte(byte(0xcf))
		if err != nil {
			return size, err
		}

		n, err := cmd.WriteInt64(int64(val))
		return size + n, err
	}
	return 1 + 8, nil
}

func packInt(cmd BufferEx, valType int, val int32) (int, error) {
	if cmd != nil {
		size, err := 1, cmd.WriteByte(byte(valType))
		if err != nil {
			return size, err
		}
		n, err := cmd.WriteInt32(val)
		return size + n, err
	}
	return 1 + 4, nil
}

func packShort(cmd BufferEx, valType int, val int16) (int, error) {
	if cmd != nil {
		size, err := 1, cmd.WriteByte(byte(valType))
		if err != nil {
			return size, err
		}

		n, err := cmd.WriteInt16(val)
		return size + n, err
	}
	return 1 + 2, nil
}

// This method is not compatible with MsgPack specs and is only used by aerospike client<->server
// for wire transfer only
func packShortRaw(cmd BufferEx, val int16) (int, error) {
	if cmd != nil {
		return cmd.WriteInt16(val)
	}
	return 2, nil
}

func packInfinity(cmd BufferEx) (int, error) {
	if cmd != nil {
		size := 0
		n, err := 1, cmd.WriteByte(byte(0xd4))
		if err != nil {
			return n, err
		}
		size += n

		n, err = 1, cmd.WriteByte(0xff)
		if err != nil {
			return size + n, err
		}
		size += n

		n, err = 1, cmd.WriteByte(0x01)
		if err != nil {
			return size + n, err
		}
		size += n

		return size, nil
	}
	return 3, nil
}

func packWildCard(cmd BufferEx) (int, error) {
	if cmd != nil {
		size := 0
		n, err := 1, cmd.WriteByte(byte(0xd4))
		if err != nil {
			return n, err
		}
		size += n

		n, err = 1, cmd.WriteByte(0xff)
		if err != nil {
			return size + n, err
		}
		size += n

		n, err = 1, cmd.WriteByte(0x00)
		if err != nil {
			return size + n, err
		}
		size += n

		return size, nil
	}
	return 3, nil
}

func packByte(cmd BufferEx, valType int, val byte) (int, error) {
	if cmd != nil {
		size := 0
		n, err := 1, cmd.WriteByte(byte(valType))
		if err != nil {
			return n, err
		}
		size += n

		n, err = 1, cmd.WriteByte(val)
		if err != nil {
			return size + n, err
		}
		size += n

		return size, nil
	}
	return 1 + 1, nil
}

// Pack nil packs a nil value
func PackNil(cmd BufferEx) (int, error) {
	return packNil(cmd)
}

func packNil(cmd BufferEx) (int, error) {
	if cmd != nil {
		return 1, cmd.WriteByte(0xc0)
	}
	return 1, nil
}

// Pack bool packs a bool value
func PackBool(cmd BufferEx, val bool) (int, error) {
	return packBool(cmd, val)
}

func packBool(cmd BufferEx, val bool) (int, error) {
	if cmd != nil {
		if val {
			return 1, cmd.WriteByte(0xc3)
		}
		return 1, cmd.WriteByte(0xc2)
	}
	return 1, nil
}

// PackFloat32 packs float32 value
func PackFloat32(cmd BufferEx, val float32) (int, error) {
	return packFloat32(cmd, val)
}

func packFloat32(cmd BufferEx, val float32) (int, error) {
	if cmd != nil {
		size := 0
		n, err := 1, cmd.WriteByte(0xca)
		if err != nil {
			return n, err
		}
		size += n
		n, err = cmd.WriteFloat32(val)
		return size + n, err
	}
	return 1 + 4, nil
}

// PackFloat64 packs float64 value
func PackFloat64(cmd BufferEx, val float64) (int, error) {
	return packFloat64(cmd, val)
}

func packFloat64(cmd BufferEx, val float64) (int, error) {
	if cmd != nil {
		size := 0
		n, err := 1, cmd.WriteByte(0xcb)
		if err != nil {
			return n, err
		}
		size += n
		n, err = cmd.WriteFloat64(val)
		return size + n, err
	}
	return 1 + 8, nil
}

func packAByte(cmd BufferEx, val byte) (int, error) {
	if cmd != nil {
		return 1, cmd.WriteByte(val)
	}
	return 1, nil
}

// packer implements a buffered packer
type packer struct {
	bytes.Buffer
	tempBuffer [8]byte
}

func newPacker() *packer {
	return &packer{}
}

// Int64ToBytes converts an int64 into slice of Bytes.
func (vb *packer) WriteInt64(num int64) (int, error) {
	return vb.WriteUint64(uint64(num))
}

// Uint64ToBytes converts an uint64 into slice of Bytes.
func (vb *packer) WriteUint64(num uint64) (int, error) {
	binary.BigEndian.PutUint64(vb.tempBuffer[:8], num)
	vb.Write(vb.tempBuffer[:8])
	return 8, nil
}

// Int32ToBytes converts an int32 to a byte slice of size 4
func (vb *packer) WriteInt32(num int32) (int, error) {
	return vb.WriteUint32(uint32(num))
}

// Uint32ToBytes converts an uint32 to a byte slice of size 4
func (vb *packer) WriteUint32(num uint32) (int, error) {
	binary.BigEndian.PutUint32(vb.tempBuffer[:4], num)
	vb.Write(vb.tempBuffer[:4])
	return 4, nil
}

// Int16ToBytes converts an int16 to slice of bytes
func (vb *packer) WriteInt16(num int16) (int, error) {
	return vb.WriteUint16(uint16(num))
}

// UInt16ToBytes converts an iuint16 to slice of bytes
func (vb *packer) WriteUint16(num uint16) (int, error) {
	binary.BigEndian.PutUint16(vb.tempBuffer[:2], num)
	vb.Write(vb.tempBuffer[:2])
	return 2, nil
}

func (vb *packer) WriteFloat32(float float32) (int, error) {
	bits := math.Float32bits(float)
	binary.BigEndian.PutUint32(vb.tempBuffer[:4], bits)
	vb.Write(vb.tempBuffer[:4])
	return 4, nil
}

func (vb *packer) WriteFloat64(float float64) (int, error) {
	bits := math.Float64bits(float)
	binary.BigEndian.PutUint64(vb.tempBuffer[:8], bits)
	vb.Write(vb.tempBuffer[:8])
	return 8, nil
}

func (vb *packer) WriteByte(b byte) error {
	_, err := vb.Write([]byte{b})
	return err
}

func (vb *packer) WriteString(s string) (int, error) {
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
