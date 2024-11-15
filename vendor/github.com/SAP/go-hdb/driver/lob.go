package driver

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"

	p "github.com/SAP/go-hdb/driver/internal/protocol"
	"github.com/SAP/go-hdb/driver/internal/unsafe"
)

func scanLob(src any, wr io.Writer) error {
	scanner, ok := src.(p.LobScanner)
	if !ok {
		return fmt.Errorf("lob: invalid scan type %T", src)
	}
	if err := scanner.Scan(wr); err != nil {
		var dbErr Error
		if errors.As(err, &dbErr) && dbErr.Code() == p.HdbErrWhileParsingProtocol {
			return ErrNestedQuery
		}
		return err
	}
	return nil
}

// ScanLobBytes supports scanning Lob data into a byte slice.
// This enables using []byte based custom types for scanning Lobs instead of using a Lob object.
// For usage please refer to the example.
func ScanLobBytes(src any, b *[]byte) error {
	if b == nil {
		return fmt.Errorf("lob scan error: parameter b %T is nil", b)
	}
	wr := new(bytes.Buffer) // cannot pool as we use the underlaying buffer (*).
	if err := scanLob(src, wr); err != nil {
		return err
	}
	*b = wr.Bytes() // (*) use underlaying buffer.
	return nil
}

// ScanLobString supports scanning Lob data into a string.
// This enables using string based custom types for scanning Lobs instead of using a Lob object.
// For usage please refer to the example.
func ScanLobString(src any, s *string) error {
	if s == nil {
		return fmt.Errorf("lob scan error: parameter s %T is nil", s)
	}
	wr := new(bytes.Buffer) // cannot pool as we use the underlaying buffer (*).
	if err := scanLob(src, wr); err != nil {
		return err
	}
	*s = unsafe.ByteSlice2String(wr.Bytes()) // (*) use underlaying buffer.
	return nil
}

// ScanLobWriter supports scanning Lob data into a io.Writer object.
// This enables using io.Writer based custom types for scanning Lobs instead of using a Lob object.
// For usage please refer to the example.
func ScanLobWriter(src any, wr io.Writer) error {
	if wr == nil {
		return fmt.Errorf("lob scan error: parameter wr %T is nil", wr)
	}
	return scanLob(src, wr)
}

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
func (l *Lob) Scan(src any) error {
	if l.wr == nil {
		l.wr = new(bytes.Buffer)
	}
	return ScanLobWriter(src, l.wr)
}

// NullLob represents an Lob that may be null.
// NullLob implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullLob struct {
	Lob   *Lob
	Valid bool // Valid is true if Lob is not NULL
}

// Scan implements the database/sql/Scanner interface.
func (n *NullLob) Scan(value any) error {
	/*
		In contrast to the Null[T] Scan implementation we do not
		create a new lob instance in case of value == nil to
		enable reuse of n.Lob.

		func (n *Null[T]) Scan(value any) error {
			if value == nil {
				n.V, n.Valid = *new(T), false
				return nil
			}
			n.Valid = true
			return convertAssign(&n.V, value)
		}
	*/
	if value == nil {
		n.Valid = false
		return nil
	}
	if n.Lob == nil {
		n.Lob = new(Lob)
	}
	n.Valid = true
	return n.Lob.Scan(value)
}

// Value implements the database/sql/Valuer interface.
func (n NullLob) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Lob, nil
}
