// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"database/sql/driver"
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
	rd io.Reader
	wr io.Writer
}

// NewLob creates a new Lob instance with the io.Reader and io.Writer given as parameters.
func NewLob(rd io.Reader, wr io.Writer) *Lob {
	return &Lob{rd: rd, wr: wr}
}

// Reader returns the io.Reader of the Lob.
func (l Lob) Reader() io.Reader {
	return l.rd
}

// SetReader sets the io.Reader source for a lob field to be written to database
// and return *Lob, to enable simple call chaining.
func (l *Lob) SetReader(rd io.Reader) *Lob {
	l.rd = rd
	return l
}

// Writer returns the io.Writer of the Lob.
func (l Lob) Writer() io.Writer {
	return l.wr
}

// SetWriter sets the io.Writer destination for a lob field to be read from database
// and return *Lob, to enable simple call chaining.
func (l *Lob) SetWriter(wr io.Writer) *Lob {
	l.wr = wr
	return l
}

// Scan implements the database/sql/Scanner interface.
func (l *Lob) Scan(src interface{}) error {
	if l.wr == nil {
		return fmt.Errorf("lob error: initial writer %[1]T %[1]v", l)
	}

	ws, ok := src.(p.WriterSetter)
	if !ok {
		return fmt.Errorf("lob: invalid scan type %T", src)
	}

	if err := ws.SetWriter(l.wr); err != nil {
		return err
	}
	return nil
}

// NullLob represents an Lob that may be null.
// NullLob implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullLob struct {
	Lob   *Lob
	Valid bool // Valid is true if Lob is not NULL
}

// Value implements the database/sql/Valuer interface.
func (l NullLob) Value() (driver.Value, error) {
	if !l.Valid {
		return nil, nil
	}
	return l.Lob, nil
}

// Scan implements the database/sql/Scanner interface.
func (l *NullLob) Scan(src interface{}) error {
	if src == nil {
		l.Valid = false
		return nil
	}
	if err := l.Lob.Scan(src); err != nil {
		return err
	}
	l.Valid = true
	return nil
}
