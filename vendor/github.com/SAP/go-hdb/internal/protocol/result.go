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
	_id := rd.ReadUint64()
	*id.id = _id

	if trace {
		outLogger.Printf("resultset id: %d", *id.id)
	}

	return rd.GetError()
}

func (id *resultsetID) write(wr *bufio.Writer) error {
	wr.WriteUint64(*id.id)

	if trace {
		outLogger.Printf("resultset id: %d", *id.id)
	}

	return nil
}

// ResultFieldSet contains database field metadata for result fields.
type ResultFieldSet struct {
	fields []*ResultField
	names  fieldNames
}

func newResultFieldSet(size int) *ResultFieldSet {
	return &ResultFieldSet{
		fields: make([]*ResultField, size),
		names:  newFieldNames(),
	}
}

// String implements the Stringer interface.
func (f *ResultFieldSet) String() string {
	a := make([]string, len(f.fields))
	for i, f := range f.fields {
		a[i] = f.String()
	}
	return fmt.Sprintf("%v", a)
}

func (f *ResultFieldSet) read(rd *bufio.Reader) {
	for i := 0; i < len(f.fields); i++ {
		field := newResultField(f.names)
		field.read(rd)
		f.fields[i] = field
	}

	pos := uint32(0)
	for _, offset := range f.names.sortOffsets() {
		diff := int(offset - pos)
		if diff > 0 {
			rd.Skip(diff)
		}
		b, size := readShortUtf8(rd)
		f.names.setName(offset, string(b))
		pos += uint32(1 + size + diff)
	}
}

// NumField returns the number of fields of a query.
func (f *ResultFieldSet) NumField() int {
	return len(f.fields)
}

// Field returns the field at index idx.
func (f *ResultFieldSet) Field(idx int) *ResultField {
	return f.fields[idx]
}

const (
	tableName = iota
	schemaName
	columnName
	columnDisplayName
	maxNames
)

// ResultField contains database field attributes for result fields.
type ResultField struct {
	fieldNames    fieldNames
	columnOptions columnOptions
	tc            TypeCode
	fraction      int16
	length        int16
	offsets       [maxNames]uint32
}

func newResultField(fieldNames fieldNames) *ResultField {
	return &ResultField{fieldNames: fieldNames}
}

// String implements the Stringer interface.
func (f *ResultField) String() string {
	return fmt.Sprintf("columnsOptions %s typeCode %s fraction %d length %d tablename %s schemaname %s columnname %s columnDisplayname %s",
		f.columnOptions,
		f.tc,
		f.fraction,
		f.length,
		f.fieldNames.name(f.offsets[tableName]),
		f.fieldNames.name(f.offsets[schemaName]),
		f.fieldNames.name(f.offsets[columnName]),
		f.fieldNames.name(f.offsets[columnDisplayName]),
	)
}

// TypeCode returns the type code of the field.
func (f *ResultField) TypeCode() TypeCode {
	return f.tc
}

// TypeLength returns the type length of the field.
// see https://golang.org/pkg/database/sql/driver/#RowsColumnTypeLength
func (f *ResultField) TypeLength() (int64, bool) {
	if f.tc.isVariableLength() {
		return int64(f.length), true
	}
	return 0, false
}

// TypePrecisionScale returns the type precision and scale (decimal types) of the field.
// see https://golang.org/pkg/database/sql/driver/#RowsColumnTypePrecisionScale
func (f *ResultField) TypePrecisionScale() (int64, int64, bool) {
	if f.tc.isDecimalType() {
		return int64(f.length), int64(f.fraction), true
	}
	return 0, 0, false
}

// Nullable returns true if the field may be null, false otherwise.
// see https://golang.org/pkg/database/sql/driver/#RowsColumnTypeNullable
func (f *ResultField) Nullable() bool {
	return f.columnOptions == coOptional
}

// Name returns the result field name.
func (f *ResultField) Name() string {
	return f.fieldNames.name(f.offsets[columnDisplayName])
}

func (f *ResultField) read(rd *bufio.Reader) {
	f.columnOptions = columnOptions(rd.ReadInt8())
	f.tc = TypeCode(rd.ReadInt8())
	f.fraction = rd.ReadInt16()
	f.length = rd.ReadInt16()
	rd.Skip(2) //filler
	for i := 0; i < maxNames; i++ {
		offset := rd.ReadUint32()
		f.offsets[i] = offset
		f.fieldNames.addOffset(offset)
	}
}

//resultset metadata
type resultMetadata struct {
	resultFieldSet *ResultFieldSet
	numArg         int
}

func (r *resultMetadata) String() string {
	return fmt.Sprintf("result metadata: %s", r.resultFieldSet.fields)
}

func (r *resultMetadata) kind() partKind {
	return pkResultMetadata
}

func (r *resultMetadata) setNumArg(numArg int) {
	r.numArg = numArg
}

func (r *resultMetadata) read(rd *bufio.Reader) error {

	r.resultFieldSet.read(rd)

	if trace {
		outLogger.Printf("read %s", r)
	}

	return rd.GetError()
}

//resultset
type resultset struct {
	numArg         int
	s              *Session
	resultFieldSet *ResultFieldSet
	fieldValues    *FieldValues
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

	cols := len(r.resultFieldSet.fields)
	r.fieldValues.resize(r.numArg, cols)

	for i := 0; i < r.numArg; i++ {
		for j, field := range r.resultFieldSet.fields {
			var err error
			if r.fieldValues.values[i*cols+j], err = readField(r.s, rd, field.TypeCode()); err != nil {
				return err
			}
		}
	}

	if trace {
		outLogger.Printf("read %s", r)
	}
	return rd.GetError()
}
