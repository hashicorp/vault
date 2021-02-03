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
)

// BufferEx is a specialized buffer interface for aerospike client.
type BufferEx interface {
	WriteInt64(num int64) (int, error)
	WriteUint64(num uint64) (int, error)
	WriteInt32(num int32) (int, error)
	WriteUint32(num uint32) (int, error)
	WriteInt16(num int16) (int, error)
	WriteUint16(num uint16) (int, error)
	WriteFloat32(float float32) (int, error)
	WriteFloat64(float float64) (int, error)
	WriteByte(b byte) error
	WriteString(s string) (int, error)
	Write(b []byte) (int, error)
}

var _ BufferEx = &buffer{}

type buffer struct {
	dataBuffer []byte
	dataOffset int
}

func newBuffer(sz int) *buffer {
	return &buffer{
		dataBuffer: make([]byte, sz),
	}
}

// Int64ToBytes converts an int64 into slice of Bytes.
func (buf *buffer) WriteInt64(num int64) (int, error) {
	return buf.WriteUint64(uint64(num))
}

// Uint64ToBytes converts an uint64 into slice of Bytes.
func (buf *buffer) WriteUint64(num uint64) (int, error) {
	binary.BigEndian.PutUint64(buf.dataBuffer[buf.dataOffset:buf.dataOffset+8], num)
	buf.dataOffset += 8
	return 8, nil
}

// Int32ToBytes converts an int32 to a byte slice of size 4
func (buf *buffer) WriteInt32(num int32) (int, error) {
	return buf.WriteUint32(uint32(num))
}

// Uint32ToBytes converts an uint32 to a byte slice of size 4
func (buf *buffer) WriteUint32(num uint32) (int, error) {
	binary.BigEndian.PutUint32(buf.dataBuffer[buf.dataOffset:buf.dataOffset+4], num)
	buf.dataOffset += 4
	return 4, nil
}

// Uint32ToBytes converts an uint32 to a byte slice of size 4
func (buf *buffer) WriteUint32At(num uint32, index int) (int, error) {
	binary.BigEndian.PutUint32(buf.dataBuffer[index:index+4], num)
	return 4, nil
}

// Int16ToBytes converts an int16 to slice of bytes
func (buf *buffer) WriteInt16(num int16) (int, error) {
	return buf.WriteUint16(uint16(num))
}

// Int16ToBytes converts an int16 to slice of bytes
func (buf *buffer) WriteUint16(num uint16) (int, error) {
	binary.BigEndian.PutUint16(buf.dataBuffer[buf.dataOffset:buf.dataOffset+2], num)
	buf.dataOffset += 2
	return 2, nil
}

func (buf *buffer) WriteFloat32(float float32) (int, error) {
	bits := math.Float32bits(float)
	binary.BigEndian.PutUint32(buf.dataBuffer[buf.dataOffset:buf.dataOffset+4], bits)
	buf.dataOffset += 4
	return 4, nil
}

func (buf *buffer) WriteFloat64(float float64) (int, error) {
	bits := math.Float64bits(float)
	binary.BigEndian.PutUint64(buf.dataBuffer[buf.dataOffset:buf.dataOffset+8], bits)
	buf.dataOffset += 8
	return 8, nil
}

func (buf *buffer) WriteByte(b byte) error {
	buf.dataBuffer[buf.dataOffset] = b
	buf.dataOffset++
	return nil
}

func (buf *buffer) WriteString(s string) (int, error) {
	copy(buf.dataBuffer[buf.dataOffset:buf.dataOffset+len(s)], s)
	buf.dataOffset += len(s)
	return len(s), nil
}

func (buf *buffer) Write(b []byte) (int, error) {
	copy(buf.dataBuffer[buf.dataOffset:buf.dataOffset+len(b)], b)
	buf.dataOffset += len(b)
	return len(b), nil
}
