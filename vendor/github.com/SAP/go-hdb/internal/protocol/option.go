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
func (o *multiLineOptions) read(rd *bufio.Reader, lineCnt int) error {

	for i := 0; i < lineCnt; i++ {

		m := plainOptions{}

		cnt, err := rd.ReadInt16()
		if err != nil {
			return err
		}
		if err := m.read(rd, int(cnt)); err != nil {
			return err
		}

		*o = append(*o, m)
	}
	return nil
}

func (o multiLineOptions) write(wr *bufio.Writer) error {
	for _, m := range o {

		if err := wr.WriteInt16(int16(len(m))); err != nil {
			return err
		}
		if err := m.write(wr); err != nil {
			return err
		}
	}
	return nil
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

func (o plainOptions) read(rd *bufio.Reader, cnt int) error {

	for i := 0; i < cnt; i++ {

		k, err := rd.ReadInt8()
		if err != nil {
			return err
		}

		tc, err := rd.ReadByte()
		if err != nil {
			return err
		}

		switch typeCode(tc) {

		default:
			outLogger.Fatalf("type code %s not implemented", typeCode(tc))

		case tcBoolean:
			if v, err := rd.ReadBool(); err == nil {
				o[k] = booleanType(v)
			} else {
				return err
			}

		case tcInt:
			if v, err := rd.ReadInt32(); err == nil {
				o[k] = intType(v)
			} else {
				return err
			}

		case tcBigint:
			if v, err := rd.ReadInt64(); err == nil {
				o[k] = bigintType(v)
			} else {
				return err
			}

		case tcDouble:
			if v, err := rd.ReadFloat64(); err == nil {
				o[k] = doubleType(v)
			} else {
				return err
			}

		case tcString:
			size, err := rd.ReadInt16()
			if err != nil {
				return err
			}
			v := make([]byte, size)
			if err := rd.ReadFull(v); err == nil {
				o[k] = stringType(v)
			} else {
				return err
			}

		case tcBstring:
			size, err := rd.ReadInt16()
			if err != nil {
				return err
			}
			v := make([]byte, size)
			if err := rd.ReadFull(v); err == nil {
				o[k] = binaryStringType(v)
			} else {
				return err
			}
		}
	}
	return nil
}

func (o plainOptions) write(wr *bufio.Writer) error {

	for k, v := range o {

		if err := wr.WriteInt8(k); err != nil {
			return err
		}

		switch v := v.(type) {

		default:
			outLogger.Fatalf("type %T not implemented", v)

		case booleanType:
			if err := wr.WriteInt8(int8(tcBoolean)); err != nil {
				return err
			}
			if err := wr.WriteBool(bool(v)); err != nil {
				return err
			}

		case intType:
			if err := wr.WriteInt8(int8(tcInt)); err != nil {
				return err
			}
			if err := wr.WriteInt32(int32(v)); err != nil {
				return err
			}

		case bigintType:
			if err := wr.WriteInt8(int8(tcBigint)); err != nil {
				return err
			}
			if err := wr.WriteInt64(int64(v)); err != nil {
				return err
			}

		case doubleType:
			if err := wr.WriteInt8(int8(tcDouble)); err != nil {
				return err
			}
			if err := wr.WriteFloat64(float64(v)); err != nil {
				return err
			}

		case stringType:
			if err := wr.WriteInt8(int8(tcString)); err != nil {
				return err
			}
			if err := wr.WriteInt16(int16(len(v))); err != nil {
				return err
			}
			if _, err := wr.Write(v); err != nil {
				return err
			}

		case binaryStringType:
			if err := wr.WriteInt8(int8(tcBstring)); err != nil {
				return err
			}
			if err := wr.WriteInt16(int16(len(v))); err != nil {
				return err
			}
			if _, err := wr.Write(v); err != nil {
				return err
			}
		}
	}
	return nil
}
