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

const (
	resultsetIDSize = 8
)

type columnOptions int8

const (
	coMandatory columnOptions = 0x01
	coOptional  columnOptions = 0x02
)

var columnOptionsText = map[columnOptions]string{
	coMandatory: "mandatory",
	coOptional:  "optional",
}

func (k columnOptions) String() string {
	t := make([]string, 0, len(columnOptionsText))

	for option, text := range columnOptionsText {
		if (k & option) != 0 {
			t = append(t, text)
		}
	}
	return fmt.Sprintf("%v", t)
}

//resultset id
type resultsetID struct {
	id *uint64
}

func (id *resultsetID) kind() partKind {
	return pkResultsetID
}

func (id *resultsetID) size() (int, error) {
	return resultsetIDSize, nil
}

func (id *resultsetID) numArg() int {
	return 1
}

func (id *resultsetID) setNumArg(int) {
	//ignore - always 1
}

func (id *resultsetID) read(rd *bufio.Reader) error {

	_id, err := rd.ReadUint64()
	if err != nil {
		return err
	}
	*id.id = _id

	if trace {
		outLogger.Printf("resultset id: %d", *id.id)
	}

	return nil
}

func (id *resultsetID) write(wr *bufio.Writer) error {

	if err := wr.WriteUint64(*id.id); err != nil {
		return err
	}
	if trace {
		outLogger.Printf("resultset id: %d", *id.id)
	}

	return nil
}

const (
	resultTableName = iota // used as index: start with 0
	resultSchemaName
	resultColumnName
	resultColumnDisplayName
	maxResultNames
)

type resultField struct {
	columnOptions           columnOptions
	tc                      typeCode
	fraction                int16
	length                  int16
	tablenameOffset         uint32
	schemanameOffset        uint32
	columnnameOffset        uint32
	columnDisplaynameOffset uint32
}

func newResultField() *resultField {
	return &resultField{}
}

func (f *resultField) String() string {
	return fmt.Sprintf("columnsOptions %s typeCode %s fraction %d length %d tablenameOffset %d schemanameOffset %d columnnameOffset %d columnDisplaynameOffset %d",
		f.columnOptions,
		f.tc,
		f.fraction,
		f.length,
		f.tablenameOffset,
		f.schemanameOffset,
		f.columnnameOffset,
		f.columnDisplaynameOffset,
	)
}

// Field interface
func (f *resultField) typeCode() typeCode {
	return f.tc
}

func (f *resultField) in() bool {
	return false
}

func (f *resultField) out() bool {
	return true
}

func (f *resultField) name(names map[uint32]string) string {
	return names[f.columnnameOffset]
}

func (f *resultField) nameOffsets() []uint32 {
	return []uint32{f.tablenameOffset, f.schemanameOffset, f.columnnameOffset, f.columnDisplaynameOffset}
}

//

func (f *resultField) read(rd *bufio.Reader) error {
	var err error

	if co, err := rd.ReadInt8(); err == nil {
		f.columnOptions = columnOptions(co)
	} else {
		return err
	}
	if tc, err := rd.ReadInt8(); err == nil {
		f.tc = typeCode(tc)
	} else {
		return err
	}
	if f.fraction, err = rd.ReadInt16(); err != nil {
		return err
	}
	if f.length, err = rd.ReadInt16(); err != nil {
		return err
	}

	if err := rd.Skip(2); err != nil { //filler
		return err
	}

	if f.tablenameOffset, err = rd.ReadUint32(); err != nil {
		return err
	}
	if f.schemanameOffset, err = rd.ReadUint32(); err != nil {
		return err
	}
	if f.columnnameOffset, err = rd.ReadUint32(); err != nil {
		return err
	}
	if f.columnDisplaynameOffset, err = rd.ReadUint32(); err != nil {
		return err
	}

	return nil
}

//resultset metadata
type resultMetadata struct {
	fieldSet *FieldSet
	numArg   int
}

func (r *resultMetadata) String() string {
	return fmt.Sprintf("result metadata: %s", r.fieldSet.fields)
}

func (r *resultMetadata) kind() partKind {
	return pkResultMetadata
}

func (r *resultMetadata) setNumArg(numArg int) {
	r.numArg = numArg
}

func (r *resultMetadata) read(rd *bufio.Reader) error {

	for i := 0; i < r.numArg; i++ {
		field := newResultField()
		if err := field.read(rd); err != nil {
			return err
		}
		r.fieldSet.fields[i] = field
	}

	pos := uint32(0)
	for _, offset := range r.fieldSet.nameOffsets() {
		if diff := int(offset - pos); diff > 0 {
			rd.Skip(diff)
		}

		b, size, err := readShortUtf8(rd)
		if err != nil {
			return err
		}

		r.fieldSet.names[offset] = string(b)

		pos += uint32(1 + size)
	}

	if trace {
		outLogger.Printf("read %s", r)
	}

	return nil
}

//resultset
type resultset struct {
	numArg      int
	fieldSet    *FieldSet
	fieldValues *FieldValues
}

func (r *resultset) String() string {
	return fmt.Sprintf("resultset: %s", r.fieldValues)
}

func (r *resultset) kind() partKind {
	return pkResultset
}

func (r *resultset) setNumArg(numArg int) {
	r.numArg = numArg
}

func (r *resultset) read(rd *bufio.Reader) error {
	if err := r.fieldValues.read(r.numArg, r.fieldSet, rd); err != nil {
		return err
	}
	if trace {
		outLogger.Printf("read %s", r)
	}
	return nil
}
