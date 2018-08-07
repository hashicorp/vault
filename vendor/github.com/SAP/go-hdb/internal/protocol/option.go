/*
Copyright 2014 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/bufio"
)

type booleanType bool

func (t booleanType) String() string {
	return fmt.Sprintf("%t", t)
}

type intType int32

func (t intType) String() string {
	return fmt.Sprintf("%d", t)
}

type bigintType int64

func (t bigintType) String() string {
	return fmt.Sprintf("%d", t)
}

type doubleType float64

func (t doubleType) String() string {
	return fmt.Sprintf("%g", t)
}

type stringType []byte

type binaryStringType []byte

func (t binaryStringType) String() string {
	return fmt.Sprintf("%v", []byte(t))
}

//multi line options (number of lines in part header argumentCount)
type multiLineOptions []plainOptions

func (o multiLineOptions) size() int {
	size := 0
	for _, m := range o {
		size += m.size()
	}
	return size
}

//pointer: append multiLineOptions itself
func (o *multiLineOptions) read(rd *bufio.Reader, lineCnt int) {
	for i := 0; i < lineCnt; i++ {
		m := plainOptions{}
		cnt := rd.ReadInt16()
		m.read(rd, int(cnt))
		*o = append(*o, m)
	}
}

func (o multiLineOptions) write(wr *bufio.Writer) {
	for _, m := range o {
		wr.WriteInt16(int16(len(m)))
		m.write(wr)
	}
}

type plainOptions map[int8]interface{}

func (o plainOptions) size() int {
	size := 2 * len(o) //option + type
	for _, v := range o {
		switch v := v.(type) {
		default:
			outLogger.Fatalf("type %T not implemented", v)
		case booleanType:
			size++
		case intType:
			size += 4
		case bigintType:
			size += 8
		case doubleType:
			size += 8
		case stringType:
			size += (2 + len(v)) //length int16 + string length
		case binaryStringType:
			size += (2 + len(v)) //length int16 + string length
		}
	}
	return size
}

func (o plainOptions) read(rd *bufio.Reader, cnt int) {

	for i := 0; i < cnt; i++ {

		k := rd.ReadInt8()
		tc := rd.ReadB()

		switch TypeCode(tc) {

		default:
			outLogger.Fatalf("type code %s not implemented", TypeCode(tc))

		case tcBoolean:
			o[k] = booleanType(rd.ReadBool())

		case tcInteger:
			o[k] = intType(rd.ReadInt32())

		case tcBigint:
			o[k] = bigintType(rd.ReadInt64())

		case tcDouble:
			o[k] = doubleType(rd.ReadFloat64())

		case tcString:
			size := rd.ReadInt16()
			v := make([]byte, size)
			rd.ReadFull(v)
			o[k] = stringType(v)

		case tcBstring:
			size := rd.ReadInt16()
			v := make([]byte, size)
			rd.ReadFull(v)
			o[k] = binaryStringType(v)

		}
	}
}

func (o plainOptions) write(wr *bufio.Writer) {

	for k, v := range o {

		wr.WriteInt8(k)

		switch v := v.(type) {

		default:
			outLogger.Fatalf("type %T not implemented", v)

		case booleanType:
			wr.WriteInt8(int8(tcBoolean))
			wr.WriteBool(bool(v))

		case intType:
			wr.WriteInt8(int8(tcInteger))
			wr.WriteInt32(int32(v))

		case bigintType:
			wr.WriteInt8(int8(tcBigint))
			wr.WriteInt64(int64(v))

		case doubleType:
			wr.WriteInt8(int8(tcDouble))
			wr.WriteFloat64(float64(v))

		case stringType:
			wr.WriteInt8(int8(tcString))
			wr.WriteInt16(int16(len(v)))
			wr.Write(v)

		case binaryStringType:
			wr.WriteInt8(int8(tcBstring))
			wr.WriteInt16(int16(len(v)))
			wr.Write(v)
		}
	}
}
