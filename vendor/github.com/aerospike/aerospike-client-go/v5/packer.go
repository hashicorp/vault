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
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"time"

	ParticleType "github.com/aerospike/aerospike-client-go/v5/internal/particle_type"
	"github.com/aerospike/aerospike-client-go/v5/types"

	Buffer "github.com/aerospike/aerospike-client-go/v5/utils/buffer"
)

var packObjectReflect func(BufferEx, interface{}, bool) (int, Error)

func packIfcList(cmd BufferEx, list []interface{}) (int, Error) {
	size := 0
	n, err := packArrayBegin(cmd, len(list))
	if err != nil {
		return n, err
	}
	size += n

	for i := range list {
		n, err = packObject(cmd, list[i], false)
		if err != nil {
			return 0, err
		}
		size += n
	}

	return size, err
}

// PackList packs any slice that implement the ListIter interface
func PackList(cmd BufferEx, list ListIter) (int, Error) {
	return packList(cmd, list)
}

func packList(cmd BufferEx, list ListIter) (int, Error) {
	size := 0
	n, err := packArrayBegin(cmd, list.Len())
	if err != nil {
		return n, err
	}
	size += n

	n, nerr := list.PackList(cmd)
	if nerr != nil {
		return size + n, newErrorAndWrap(nerr, types.SERIALIZE_ERROR)
	}
	return size + n, nil
}

func packValueArray(cmd BufferEx, list ValueArray) (int, Error) {
	size := 0
	n, err := packArrayBegin(cmd, len(list))
	if err != nil {
		return n, err
	}
	size += n

	for i := range list {
		n, err = list[i].pack(cmd)
		if err != nil {
			return 0, err
		}
		size += n
	}

	return size, err
}

func packArrayBegin(cmd BufferEx, size int) (int, Error) {
	if size < 16 {
		return packAByte(cmd, 0x90|byte(size))
	} else if size <= math.MaxUint16 {
		return packShort(cmd, 0xdc, int16(size))
	} else {
		return packInt(cmd, 0xdd, int32(size))
	}
}

func packIfcMap(cmd BufferEx, theMap map[interface{}]interface{}) (int, Error) {
	size := 0
	n, err := packMapBegin(cmd, len(theMap))
	if err != nil {
		return n, err
	}
	size += n

	for k, v := range theMap {
		n, err = packObject(cmd, k, true)
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
func PackJson(cmd BufferEx, theMap map[string]interface{}) (int, Error) {
	return packJsonMap(cmd, theMap)
}

func packJsonMap(cmd BufferEx, theMap map[string]interface{}) (int, Error) {
	size := 0
	n, err := packMapBegin(cmd, len(theMap))
	if err != nil {
		return n, err
	}
	size += n

	for k, v := range theMap {
		n, err = packString(cmd, k)
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
func PackMap(cmd BufferEx, theMap MapIter) (int, Error) {
	return packMap(cmd, theMap)
}

func packMap(cmd BufferEx, theMap MapIter) (int, Error) {
	size := 0
	n, err := packMapBegin(cmd, theMap.Len())
	if err != nil {
		return n, err
	}
	size += n

	n, nerr := theMap.PackMap(cmd)
	if nerr != nil {
		return size + n, newErrorAndWrap(nerr, types.SERIALIZE_ERROR)
	}
	return size + n, nil
}

func packMapBegin(cmd BufferEx, size int) (int, Error) {
	if size < 16 {
		return packAByte(cmd, 0x80|byte(size))
	} else if size <= math.MaxUint16 {
		return packShort(cmd, 0xde, int16(size))
	} else {
		return packInt(cmd, 0xdf, int32(size))
	}
}

// PackBytes backs a byte array
func PackBytes(cmd BufferEx, b []byte) (int, Error) {
	return packBytes(cmd, b)
}

func packBytes(cmd BufferEx, b []byte) (int, Error) {
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

func packByteArrayBegin(cmd BufferEx, length int) (int, Error) {
	// Use string header codes for byte arrays.
	return packStringBegin(cmd, length)
}

func packObject(cmd BufferEx, obj interface{}, mapKey bool) (int, Error) {
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
			return 0, newError(types.SERIALIZE_ERROR, fmt.Sprintf("Maps, Slices, and bounded arrays other than Bounded Byte Arrays are not supported as Map keys. Value: %#v", v))
		}
		return packIfcMap(cmd, map[interface{}]interface{}{})
	case []interface{}:
		if mapKey {
			return 0, newError(types.SERIALIZE_ERROR, fmt.Sprintf("Maps, Slices, and bounded arrays other than Bounded Byte Arrays are not supported as Map keys. Value: %#v", v))
		}
		return packIfcList(cmd, v)
	case map[interface{}]interface{}:
		if mapKey {
			return 0, newError(types.SERIALIZE_ERROR, fmt.Sprintf("Maps, Slices, and bounded arrays other than Bounded Byte Arrays are not supported as Map keys. Value: %#v", v))
		}
		return packIfcMap(cmd, v)
	case ListIter:
		if mapKey {
			return 0, newError(types.SERIALIZE_ERROR, fmt.Sprintf("Maps, Slices, and bounded arrays other than Bounded Byte Arrays are not supported as Map keys. Value: %#v", v))
		}
		return packList(cmd, obj.(ListIter))
	case MapIter:
		if mapKey {
			return 0, newError(types.SERIALIZE_ERROR, fmt.Sprintf("Maps, Slices, and bounded arrays other than Bounded Byte Arrays are not supported as Map keys. Value: %#v", v))
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

	return 0, newError(types.SERIALIZE_ERROR, fmt.Sprintf("Type `%v (%s)` not supported to pack. ", obj, reflect.TypeOf(obj).String()))
}

func packAUInt64(cmd BufferEx, val uint64) (int, Error) {
	return packUInt64(cmd, val)
}

func packAInt64(cmd BufferEx, val int64) (int, Error) {
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
func PackInt64(cmd BufferEx, val int64) (int, Error) {
	return packAInt64(cmd, val)
}

func packAInt(cmd BufferEx, val int) (int, Error) {
	return packAInt64(cmd, int64(val))
}

// PackString packs a string
func PackString(cmd BufferEx, val string) (int, Error) {
	return packString(cmd, val)
}

func packStringBegin(cmd BufferEx, size int) (int, Error) {
	if size < 32 {
		return packAByte(cmd, 0xa0|byte(size))
	} else if size < 256 {
		return packByte(cmd, 0xd9, byte(size))
	} else if size < 65536 {
		return packShort(cmd, 0xda, int16(size))
	}
	return packInt(cmd, 0xdb, int32(size))
}

func packString(cmd BufferEx, val string) (int, Error) {
	size := 0
	slen := len(val) + 1
	n, err := packStringBegin(cmd, slen)
	if err != nil {
		return n, err
	}
	size += n

	if cmd != nil {
		cmd.WriteByte(byte(ParticleType.STRING))
		size++

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

func packRawString(cmd BufferEx, val string) (int, Error) {
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

func packGeoJson(cmd BufferEx, val string) (int, Error) {
	size := 0
	slen := len(val) + 1
	n, err := packByteArrayBegin(cmd, slen)
	if err != nil {
		return n, err
	}
	size += n

	if cmd != nil {
		cmd.WriteByte(byte(ParticleType.GEOJSON))
		size++

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

func packByteArray(cmd BufferEx, src []byte) (int, Error) {
	if cmd != nil {
		return cmd.Write(src)
	}
	return len(src), nil
}

func packInt64(cmd BufferEx, valType int, val int64) (int, Error) {
	if cmd != nil {
		cmd.WriteByte(byte(valType))
		cmd.WriteInt64(val)
	}
	return 1 + 8, nil
}

// PackUInt64 packs a uint64
func PackUInt64(cmd BufferEx, val uint64) (int, Error) {
	return packUInt64(cmd, val)
}

func packUInt64(cmd BufferEx, val uint64) (int, Error) {
	if cmd != nil {
		cmd.WriteByte(byte(0xcf))
		cmd.WriteInt64(int64(val))
	}
	return 1 + 8, nil
}

func packInt(cmd BufferEx, valType int, val int32) (int, Error) {
	if cmd != nil {
		cmd.WriteByte(byte(valType))
		cmd.WriteInt32(val)
	}
	return 1 + 4, nil
}

func packShort(cmd BufferEx, valType int, val int16) (int, Error) {
	if cmd != nil {
		cmd.WriteByte(byte(valType))
		cmd.WriteInt16(val)
	}
	return 1 + 2, nil
}

// This method is not compatible with MsgPack specs and is only used by aerospike client<->server
// for wire transfer only
func packShortRaw(cmd BufferEx, val int16) (int, Error) {
	if cmd != nil {
		cmd.WriteInt16(val)
	}
	return 2, nil
}

func packInfinity(cmd BufferEx) (int, Error) {
	if cmd != nil {
		cmd.WriteByte(byte(0xd4))
		cmd.WriteByte(0xff)
		cmd.WriteByte(0x01)
	}
	return 3, nil
}

func packWildCard(cmd BufferEx) (int, Error) {
	if cmd != nil {
		cmd.WriteByte(byte(0xd4))
		cmd.WriteByte(0xff)
		cmd.WriteByte(0x00)
	}
	return 3, nil
}

func packByte(cmd BufferEx, valType int, val byte) (int, Error) {
	if cmd != nil {
		cmd.WriteByte(byte(valType))
		cmd.WriteByte(val)
	}
	return 1 + 1, nil
}

// PackNil packs a nil value
func PackNil(cmd BufferEx) (int, Error) {
	return packNil(cmd)
}

func packNil(cmd BufferEx) (int, Error) {
	if cmd != nil {
		cmd.WriteByte(0xc0)
	}
	return 1, nil
}

// PackBool packs a bool value
func PackBool(cmd BufferEx, val bool) (int, Error) {
	return packBool(cmd, val)
}

func packBool(cmd BufferEx, val bool) (int, Error) {
	if cmd != nil {
		if val {
			cmd.WriteByte(0xc3)
		} else {
			cmd.WriteByte(0xc2)
		}
	}
	return 1, nil
}

// PackFloat32 packs float32 value
func PackFloat32(cmd BufferEx, val float32) (int, Error) {
	return packFloat32(cmd, val)
}

func packFloat32(cmd BufferEx, val float32) (int, Error) {
	if cmd != nil {
		cmd.WriteByte(0xca)
		cmd.WriteFloat32(val)
	}
	return 1 + 4, nil
}

// PackFloat64 packs float64 value
func PackFloat64(cmd BufferEx, val float64) (int, Error) {
	return packFloat64(cmd, val)
}

func packFloat64(cmd BufferEx, val float64) (int, Error) {
	if cmd != nil {
		cmd.WriteByte(0xcb)
		cmd.WriteFloat64(val)
	}
	return 1 + 8, nil
}

func packAByte(cmd BufferEx, val byte) (int, Error) {
	if cmd != nil {
		cmd.WriteByte(val)
	}
	return 1, nil
}

/***************************************************************************

	packer

***************************************************************************/

// packer implements a buffered packer
type packer struct {
	bytes.Buffer
	tempBuffer [8]byte
}

func newPacker() *packer {
	return &packer{}
}

// WriteInt64 writes an int64 to the buffer
func (vb *packer) WriteInt64(num int64) int {
	return vb.WriteUint64(uint64(num))
}

// WriteUint64 writes an uint64 to the buffer
func (vb *packer) WriteUint64(num uint64) int {
	binary.BigEndian.PutUint64(vb.tempBuffer[:8], num)
	n, _ := vb.Write(vb.tempBuffer[:8])
	return n
}

// WriteInt32 writes an int32 to the buffer
func (vb *packer) WriteInt32(num int32) int {
	return vb.WriteUint32(uint32(num))
}

// WriteUint32 writes an uint32 to the buffer
func (vb *packer) WriteUint32(num uint32) int {
	binary.BigEndian.PutUint32(vb.tempBuffer[:4], num)
	n, _ := vb.Write(vb.tempBuffer[:4])
	return n
}

// WriteInt16 writes an int16 to the buffer
func (vb *packer) WriteInt16(num int16) int {
	return vb.WriteUint16(uint16(num))
}

// WriteUint16 writes an uint16 to the buffer
func (vb *packer) WriteUint16(num uint16) int {
	binary.BigEndian.PutUint16(vb.tempBuffer[:2], num)
	n, _ := vb.Write(vb.tempBuffer[:2])
	return n
}

// WriteFloat32 writes an float32 to the buffer
func (vb *packer) WriteFloat32(float float32) int {
	bits := math.Float32bits(float)
	binary.BigEndian.PutUint32(vb.tempBuffer[:4], bits)
	n, _ := vb.Write(vb.tempBuffer[:4])
	return n
}

// WriteFloat64 writes an float64 to the buffer
func (vb *packer) WriteFloat64(float float64) int {
	bits := math.Float64bits(float)
	binary.BigEndian.PutUint64(vb.tempBuffer[:8], bits)
	n, _ := vb.Write(vb.tempBuffer[:8])
	return n
}

// WriteBool writes a bool to the buffer
func (vb *packer) WriteBool(b bool) int {
	if b {
		vb.WriteByte(1)
	} else {
		vb.WriteByte(0)
	}
	return 1
}

// WriteBytes writes a byte to the buffer
func (vb *packer) WriteByte(b byte) {
	vb.Write([]byte{b})
}

// Write writes a byte slice to the buffer
func (vb *packer) Write(b []byte) (int, Error) {
	n, err := vb.Buffer.Write(b)
	if err != nil {
		return n, newCommonError(err)
	}
	return n, nil
}

// WriteString writes a string to the buffer
func (vb *packer) WriteString(s string) (int, Error) {
	n, err := vb.Buffer.WriteString(s)
	if err != nil {
		return n, newCommonError(err)
	}
	return n, nil
}
