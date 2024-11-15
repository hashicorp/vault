package pgtype

import (
	"database/sql/driver"
)

type Bit Varbit

func (dst *Bit) Set(src interface{}) error {
	return (*Varbit)(dst).Set(src)
}

func (dst Bit) Get() interface{} {
	return (Varbit)(dst).Get()
}

func (src *Bit) AssignTo(dst interface{}) error {
	return (*Varbit)(src).AssignTo(dst)
}

func (dst *Bit) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*Varbit)(dst).DecodeBinary(ci, src)
}

func (src Bit) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (Varbit)(src).EncodeBinary(ci, buf)
}

func (dst *Bit) DecodeText(ci *ConnInfo, src []byte) error {
	return (*Varbit)(dst).DecodeText(ci, src)
}

func (src Bit) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (Varbit)(src).EncodeText(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *Bit) Scan(src interface{}) error {
	return (*Varbit)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src Bit) Value() (driver.Value, error) {
	return (Varbit)(src).Value()
}
