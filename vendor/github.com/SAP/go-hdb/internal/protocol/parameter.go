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
	"database/sql/driver"
	"fmt"
	"io"

	"github.com/SAP/go-hdb/internal/bufio"
)

type parameterOptions int8

const (
	poMandatory parameterOptions = 0x01
	poOptional  parameterOptions = 0x02
	poDefault   parameterOptions = 0x04
)

var parameterOptionsText = map[parameterOptions]string{
	poMandatory: "mandatory",
	poOptional:  "optional",
	poDefault:   "default",
}

func (k parameterOptions) String() string {
	t := make([]string, 0, len(parameterOptionsText))

	for option, text := range parameterOptionsText {
		if (k & option) != 0 {
			t = append(t, text)
		}
	}
	return fmt.Sprintf("%v", t)
}

type parameterMode int8

const (
	pmIn    parameterMode = 0x01
	pmInout parameterMode = 0x02
	pmOut   parameterMode = 0x04
)

var parameterModeText = map[parameterMode]string{
	pmIn:    "in",
	pmInout: "inout",
	pmOut:   "out",
}

func (k parameterMode) String() string {
	t := make([]string, 0, len(parameterModeText))

	for mode, text := range parameterModeText {
		if (k & mode) != 0 {
			t = append(t, text)
		}
	}
	return fmt.Sprintf("%v", t)
}

// ParameterFieldSet contains database field metadata for parameters.
type ParameterFieldSet struct {
	fields        []*ParameterField
	_inputFields  []*ParameterField
	_outputFields []*ParameterField
	names         fieldNames
}

func newParameterFieldSet(size int) *ParameterFieldSet {
	return &ParameterFieldSet{
		fields:        make([]*ParameterField, size),
		_inputFields:  make([]*ParameterField, 0, size),
		_outputFields: make([]*ParameterField, 0, size),
		names:         newFieldNames(),
	}
}

// String implements the Stringer interface.
func (f *ParameterFieldSet) String() string {
	a := make([]string, len(f.fields))
	for i, f := range f.fields {
		a[i] = f.String()
	}
	return fmt.Sprintf("%v", a)
}

func (f *ParameterFieldSet) read(rd *bufio.Reader) {
	for i := 0; i < len(f.fields); i++ {
		field := newParameterField(f.names)
		field.read(rd)
		f.fields[i] = field
		if field.In() {
			f._inputFields = append(f._inputFields, field)
		}
		if field.Out() {
			f._outputFields = append(f._outputFields, field)
		}
	}

	pos := uint32(0)
	for _, offset := range f.names.sortOffsets() {
		if diff := int(offset - pos); diff > 0 {
			rd.Skip(diff)
		}
		b, size := readShortUtf8(rd)
		f.names.setName(offset, string(b))
		pos += uint32(1 + size)
	}
}

func (f *ParameterFieldSet) inputFields() []*ParameterField {
	return f._inputFields
}

func (f *ParameterFieldSet) outputFields() []*ParameterField {
	return f._outputFields
}

// NumInputField returns the number of input fields in a database statement.
func (f *ParameterFieldSet) NumInputField() int {
	return len(f._inputFields)
}

// NumOutputField returns the number of output fields of a query or stored procedure.
func (f *ParameterFieldSet) NumOutputField() int {
	return len(f._outputFields)
}

// Field returns the field at index idx.
func (f *ParameterFieldSet) Field(idx int) *ParameterField {
	return f.fields[idx]
}

// OutputField returns the output field at index idx.
func (f *ParameterFieldSet) OutputField(idx int) *ParameterField {
	return f._outputFields[idx]
}

// ParameterField contains database field attributes for parameters.
type ParameterField struct {
	fieldNames       fieldNames
	parameterOptions parameterOptions
	tc               TypeCode
	mode             parameterMode
	fraction         int16
	length           int16
	offset           uint32
	chunkReader      lobChunkReader
	lobLocatorID     locatorID
}

func newParameterField(fieldNames fieldNames) *ParameterField {
	return &ParameterField{fieldNames: fieldNames}
}

// String implements the Stringer interface.
func (f *ParameterField) String() string {
	return fmt.Sprintf("parameterOptions %s typeCode %s mode %s fraction %d length %d name %s",
		f.parameterOptions,
		f.tc,
		f.mode,
		f.fraction,
		f.length,
		f.Name(),
	)
}

// TypeCode returns the type code of the field.
func (f *ParameterField) TypeCode() TypeCode {
	return f.tc
}

// TypeLength returns the type length of the field.
// see https://golang.org/pkg/database/sql/driver/#RowsColumnTypeLength
func (f *ParameterField) TypeLength() (int64, bool) {
	if f.tc.isVariableLength() {
		return int64(f.length), true
	}
	return 0, false
}

// TypePrecisionScale returns the type precision and scale (decimal types) of the field.
// see https://golang.org/pkg/database/sql/driver/#RowsColumnTypePrecisionScale
func (f *ParameterField) TypePrecisionScale() (int64, int64, bool) {
	if f.tc.isDecimalType() {
		return int64(f.length), int64(f.fraction), true
	}
	return 0, 0, false
}

// Nullable returns true if the field may be null, false otherwise.
// see https://golang.org/pkg/database/sql/driver/#RowsColumnTypeNullable
func (f *ParameterField) Nullable() bool {
	return f.parameterOptions == poOptional
}

// In returns true if the parameter field is an input field.
func (f *ParameterField) In() bool {
	return f.mode == pmInout || f.mode == pmIn
}

// Out returns true if the parameter field is an output field.
func (f *ParameterField) Out() bool {
	return f.mode == pmInout || f.mode == pmOut
}

// Name returns the parameter field name.
func (f *ParameterField) Name() string {
	return f.fieldNames.name(f.offset)
}

// SetLobReader sets the io.Reader if a Lob parameter field.
func (f *ParameterField) SetLobReader(rd io.Reader) error {
	f.chunkReader = newLobChunkReader(f.TypeCode().isCharBased(), rd)
	return nil
}

//

func (f *ParameterField) read(rd *bufio.Reader) {
	f.parameterOptions = parameterOptions(rd.ReadInt8())
	f.tc = TypeCode(rd.ReadInt8())
	f.mode = parameterMode(rd.ReadInt8())
	rd.Skip(1) //filler
	f.offset = rd.ReadUint32()
	f.fieldNames.addOffset(f.offset)
	f.length = rd.ReadInt16()
	f.fraction = rd.ReadInt16()
	rd.Skip(4) //filler
}

// parameter metadata
type parameterMetadata struct {
	prmFieldSet *ParameterFieldSet
	numArg      int
}

func (m *parameterMetadata) String() string {
	return fmt.Sprintf("parameter metadata: %s", m.prmFieldSet.fields)
}

func (m *parameterMetadata) kind() partKind {
	return pkParameterMetadata
}

func (m *parameterMetadata) setNumArg(numArg int) {
	m.numArg = numArg
}

func (m *parameterMetadata) read(rd *bufio.Reader) error {

	m.prmFieldSet.read(rd)

	if trace {
		outLogger.Printf("read %s", m)
	}

	return rd.GetError()
}

// input parameters
type inputParameters struct {
	inputFields []*ParameterField
	args        []driver.NamedValue
}

func newInputParameters(inputFields []*ParameterField, args []driver.NamedValue) *inputParameters {
	return &inputParameters{inputFields: inputFields, args: args}
}

func (p *inputParameters) String() string {
	return fmt.Sprintf("input parameters: %v", p.args)
}

func (p *inputParameters) kind() partKind {
	return pkParameters
}

func (p *inputParameters) size() (int, error) {

	size := len(p.args)
	cnt := len(p.inputFields)

	for i, arg := range p.args {

		if arg.Value == nil { // null value
			continue
		}

		// mass insert
		field := p.inputFields[i%cnt]

		fieldSize, err := fieldSize(field.TypeCode(), arg)
		if err != nil {
			return 0, err
		}

		size += fieldSize
	}

	return size, nil
}

func (p *inputParameters) numArg() int {
	cnt := len(p.inputFields)

	if cnt == 0 { // avoid divide-by-zero (e.g. prepare without parameters)
		return 0
	}

	return len(p.args) / cnt
}

func (p *inputParameters) write(wr *bufio.Writer) error {

	cnt := len(p.inputFields)

	for i, arg := range p.args {

		//mass insert
		field := p.inputFields[i%cnt]

		if err := writeField(wr, field.TypeCode(), arg); err != nil {
			return err
		}
	}

	if trace {
		outLogger.Printf("input parameters: %s", p)
	}

	return nil
}

// output parameter
type outputParameters struct {
	numArg       int
	s            *Session
	outputFields []*ParameterField
	fieldValues  *FieldValues
}

func (p *outputParameters) String() string {
	return fmt.Sprintf("output parameters: %v", p.fieldValues)
}

func (p *outputParameters) kind() partKind {
	return pkOutputParameters
}

func (p *outputParameters) setNumArg(numArg int) {
	p.numArg = numArg // should always be 1
}

func (p *outputParameters) read(rd *bufio.Reader) error {

	cols := len(p.outputFields)
	p.fieldValues.resize(p.numArg, cols)

	for i := 0; i < p.numArg; i++ {
		for j, field := range p.outputFields {
			var err error
			if p.fieldValues.values[i*cols+j], err = readField(p.s, rd, field.TypeCode()); err != nil {
				return err
			}
		}
	}

	if trace {
		outLogger.Printf("read %s", p)
	}
	return rd.GetError()
}
