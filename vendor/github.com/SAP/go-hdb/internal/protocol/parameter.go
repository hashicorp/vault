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

type parameterField struct {
	parameterOptions parameterOptions
	tc               typeCode
	mode             parameterMode
	fraction         int16
	length           int16
	nameOffset       uint32
}

func newParameterField() *parameterField {
	return &parameterField{}
}

func (f *parameterField) String() string {
	return fmt.Sprintf("parameterOptions %s typeCode %s mode %s fraction %d length %d nameOffset %d",
		f.parameterOptions,
		f.tc,
		f.mode,
		f.fraction,
		f.length,
		f.nameOffset,
	)
}

// field interface
func (f *parameterField) typeCode() typeCode {
	return f.tc
}

func (f *parameterField) in() bool {
	return f.mode == pmInout || f.mode == pmIn
}

func (f *parameterField) out() bool {
	return f.mode == pmInout || f.mode == pmOut
}

func (f *parameterField) name(names map[uint32]string) string {
	return names[f.nameOffset]
}

func (f *parameterField) nameOffsets() []uint32 {
	return []uint32{f.nameOffset}
}

//

func (f *parameterField) read(rd *bufio.Reader) error {
	var err error

	if po, err := rd.ReadInt8(); err == nil {
		f.parameterOptions = parameterOptions(po)
	} else {
		return err
	}
	if tc, err := rd.ReadInt8(); err == nil {
		f.tc = typeCode(tc)
	} else {
		return err
	}
	if mode, err := rd.ReadInt8(); err == nil {
		f.mode = parameterMode(mode)
	} else {
		return err
	}
	if err := rd.Skip(1); err != nil { //filler
		return err
	}
	if f.nameOffset, err = rd.ReadUint32(); err != nil {
		return err
	}
	if f.length, err = rd.ReadInt16(); err != nil {
		return err
	}
	if f.fraction, err = rd.ReadInt16(); err != nil {
		return err
	}
	if err := rd.Skip(4); err != nil { //filler
		return err
	}
	return nil
}

// parameter metadata
type parameterMetadata struct {
	fieldSet *FieldSet
	numArg   int
}

func (m *parameterMetadata) String() string {
	return fmt.Sprintf("parameter metadata: %s", m.fieldSet.fields)
}

func (m *parameterMetadata) kind() partKind {
	return pkParameterMetadata
}

func (m *parameterMetadata) setNumArg(numArg int) {
	m.numArg = numArg
}

func (m *parameterMetadata) read(rd *bufio.Reader) error {

	for i := 0; i < m.numArg; i++ {
		field := newParameterField()
		if err := field.read(rd); err != nil {
			return err
		}
		m.fieldSet.fields[i] = field
	}

	pos := uint32(0)
	for _, offset := range m.fieldSet.nameOffsets() {
		if diff := int(offset - pos); diff > 0 {
			rd.Skip(diff)
		}

		b, size, err := readShortUtf8(rd)
		if err != nil {
			return err
		}

		m.fieldSet.names[offset] = string(b)

		pos += uint32(1 + size)
	}

	if trace {
		outLogger.Printf("read %s", m)
	}

	return nil
}

// parameters
type parameters struct {
	fields []field //input fields
	args   []driver.Value
}

func newParameters(fieldSet *FieldSet, args []driver.Value) *parameters {
	m := &parameters{
		fields: make([]field, 0, len(fieldSet.fields)),
		args:   args,
	}
	for _, field := range fieldSet.fields {
		if field.in() {
			m.fields = append(m.fields, field)
		}
	}
	return m
}

func (m *parameters) kind() partKind {
	return pkParameters
}

func (m *parameters) size() (int, error) {

	size := len(m.args)
	cnt := len(m.fields)

	for i, arg := range m.args {

		if arg == nil { // null value
			continue
		}

		// mass insert
		field := m.fields[i%cnt]

		fieldSize, err := fieldSize(field.typeCode(), arg)
		if err != nil {
			return 0, err
		}

		size += fieldSize
	}

	return size, nil
}

func (m *parameters) numArg() int {
	cnt := len(m.fields)

	if cnt == 0 { // avoid divide-by-zero (e.g. prepare without parameters)
		return 0
	}

	return len(m.args) / cnt
}

func (m parameters) write(wr *bufio.Writer) error {

	cnt := len(m.fields)

	for i, arg := range m.args {

		//mass insert
		field := m.fields[i%cnt]

		if err := writeField(wr, field.typeCode(), arg); err != nil {
			return err
		}
	}

	if trace {
		outLogger.Printf("parameters: %s", m)
	}

	return nil
}

// output parameter
type outputParameters struct {
	numArg      int
	fieldSet    *FieldSet
	fieldValues *FieldValues
}

func (r *outputParameters) String() string {
	return fmt.Sprintf("output parameters: %v", r.fieldValues)
}

func (r *outputParameters) kind() partKind {
	return pkOutputParameters
}

func (r *outputParameters) setNumArg(numArg int) {
	r.numArg = numArg // should always be 1
}

func (r *outputParameters) read(rd *bufio.Reader) error {
	if err := r.fieldValues.read(r.numArg, r.fieldSet, rd); err != nil {
		return err
	}
	if trace {
		outLogger.Printf("read %s", r)
	}
	return nil
}
