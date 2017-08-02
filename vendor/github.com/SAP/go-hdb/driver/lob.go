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

package driver

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io"

	p "github.com/SAP/go-hdb/internal/protocol"
)

// A Lob is the driver representation of a database large object field.
// A Lob object uses an io.Reader object as source for writing content to a database lob field.
// A Lob object uses an io.Writer object as destination for reading content from a database lob field.
// A Lob can be created by contructor method NewLob with io.Reader and io.Writer as parameters or
// created by new, setting io.Reader and io.Writer by SetReader and SetWriter methods.
type Lob struct {
	rd         io.Reader
	wr         io.Writer
	writeDescr *p.LobWriteDescr
}

// NewLob creates a new Lob instance with the io.Reader and io.Writer given as parameters.
func NewLob(rd io.Reader, wr io.Writer) *Lob {
	return &Lob{rd: rd, wr: wr}
}

// SetReader sets the io.Reader source for a lob field to be written to database.
func (l *Lob) SetReader(rd io.Reader) {
	l.rd = rd
}

// SetWriter sets the io.Writer destination for a lob field to be read from database.
func (l *Lob) SetWriter(wr io.Writer) {
	l.wr = wr
}

// Scan implements the database/sql/Scanner interface.
func (l *Lob) Scan(src interface{}) error {

	if l.wr == nil {
		return errors.New("lob error: initial writer")
	}

	ptr, ok := src.(int64)
	if !ok {
		return fmt.Errorf("lob: invalid pointer type %T", src)
	}

	descr := p.PointerToLobReadDescr(ptr)
	if err := descr.SetWriter(l.wr); err != nil {
		return err
	}
	return nil
}

// Value implements the database/sql/Valuer interface.
func (l *Lob) Value() (driver.Value, error) {
	if l.rd == nil {
		return nil, errors.New("lob error: initial reader")
	}
	if l.writeDescr == nil {
		l.writeDescr = new(p.LobWriteDescr)
	}
	l.writeDescr.SetReader(l.rd)
	return p.LobWriteDescrToPointer(l.writeDescr), nil
}
