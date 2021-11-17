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
