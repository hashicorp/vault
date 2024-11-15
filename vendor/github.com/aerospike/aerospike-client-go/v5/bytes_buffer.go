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
)

// BufferEx is a specialized buffer interface for aerospike client.
type BufferEx interface {
	WriteInt64(num int64) int
	WriteUint64(num uint64) int
	WriteInt32(num int32) int
	WriteUint32(num uint32) int
	WriteInt16(num int16) int
	WriteUint16(num uint16) int
	WriteFloat32(float float32) int
	WriteFloat64(float float64) int
	WriteBool(b bool) int
	WriteByte(b byte)
	WriteString(s string) (int, Error)
	Write(b []byte) (int, Error)
}

var _ BufferEx = &bufferEx{}

type bufferEx struct {
	dataBuffer []byte
	dataOffset int
}

func newBuffer(sz int) *bufferEx {
	return &bufferEx{
		dataBuffer: make([]byte, sz),
	}
}

// Bytes returns the content of the buffer
func (buf *bufferEx) Bytes() []byte {
	return buf.dataBuffer[:buf.dataOffset]
}

// Int64ToBytes converts an int64 into slice of Bytes.
func (buf *bufferEx) WriteInt64(num int64) int {
	return buf.WriteUint64(uint64(num))
}

// Uint64ToBytes converts an uint64 into slice of Bytes.
func (buf *bufferEx) WriteUint64(num uint64) int {
	binary.BigEndian.PutUint64(buf.dataBuffer[buf.dataOffset:buf.dataOffset+8], num)
	buf.dataOffset += 8
	return 8
}

// Int32ToBytes converts an int32 to a byte slice of size 4
func (buf *bufferEx) WriteInt32(num int32) int {
	return buf.WriteUint32(uint32(num))
}

// Uint32ToBytes converts an uint32 to a byte slice of size 4
func (buf *bufferEx) WriteUint32(num uint32) int {
	binary.BigEndian.PutUint32(buf.dataBuffer[buf.dataOffset:buf.dataOffset+4], num)
	buf.dataOffset += 4
	return 4
}

// Uint32ToBytes converts an uint32 to a byte slice of size 4
func (buf *bufferEx) WriteUint32At(num uint32, index int) int {
	binary.BigEndian.PutUint32(buf.dataBuffer[index:index+4], num)
	return 4
}

// Int16ToBytes converts an int16 to slice of bytes
func (buf *bufferEx) WriteInt16(num int16) int {
	return buf.WriteUint16(uint16(num))
}

func (buf *bufferEx) WriteInt16LittleEndian(num uint16) int {
	binary.LittleEndian.PutUint16(buf.dataBuffer[buf.dataOffset:buf.dataOffset+2], num)
	buf.dataOffset += 2
	return 2
}

// Int16ToBytes converts an int16 to slice of bytes
func (buf *bufferEx) WriteUint16(num uint16) int {
	binary.BigEndian.PutUint16(buf.dataBuffer[buf.dataOffset:buf.dataOffset+2], num)
	buf.dataOffset += 2
	return 2
}

func (buf *bufferEx) WriteFloat32(float float32) int {
	bits := math.Float32bits(float)
	binary.BigEndian.PutUint32(buf.dataBuffer[buf.dataOffset:buf.dataOffset+4], bits)
	buf.dataOffset += 4
	return 4
}

func (buf *bufferEx) WriteFloat64(float float64) int {
	bits := math.Float64bits(float)
	binary.BigEndian.PutUint64(buf.dataBuffer[buf.dataOffset:buf.dataOffset+8], bits)
	buf.dataOffset += 8
	return 8
}

func (buf *bufferEx) WriteByte(b byte) {
	buf.dataBuffer[buf.dataOffset] = b
	buf.dataOffset++
}

func (buf *bufferEx) WriteString(s string) (int, Error) {
	copy(buf.dataBuffer[buf.dataOffset:buf.dataOffset+len(s)], s)
	buf.dataOffset += len(s)
	return len(s), nil
}

func (buf *bufferEx) WriteBool(b bool) int {
	if b {
		buf.WriteByte(1)
	} else {
		buf.WriteByte(0)
	}
	return 1
}

func (buf *bufferEx) Write(b []byte) (int, Error) {
	copy(buf.dataBuffer[buf.dataOffset:buf.dataOffset+len(b)], b)
	buf.dataOffset += len(b)
	return len(b), nil
}
