package pgtype

import (
	"database/sql/driver"
)

type Varchar Text

// Set converts from src to dst. Note that as Varchar is not a general
// number type Set does not do automatic type conversion as other number
// types do.
func (dst *Varchar) Set(src interface{}) error {
	return (*Text)(dst).Set(src)
}

func (dst *Varchar) Get() interface{} {
	return (*Text)(dst).Get()
}

// AssignTo assigns from src to dst. Note that as Varchar is not a general number
// type AssignTo does not do automatic type conversion as other number types do.
func (src *Varchar) AssignTo(dst interface{}) error {
	return (*Text)(src).AssignTo(dst)
}

func (dst *Varchar) DecodeText(ci *ConnInfo, src []byte) error {
	return (*Text)(dst).DecodeText(ci, src)
}

func (dst *Varchar) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*Text)(dst).DecodeBinary(ci, src)
}

func (src *Varchar) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Text)(src).EncodeText(ci, buf)
}

func (src *Varchar) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Text)(src).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *Varchar) Scan(src interface{}) error {
	return (*Text)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *Varchar) Value() (driver.Value, error) {
	return (*Text)(src).Value()
}

func (src *Varchar) MarshalJSON() ([]byte, error) {
	return (*Text)(src).MarshalJSON()
}

func (dst *Varchar) UnmarshalJSON(b []byte) error {
	return (*Text)(dst).UnmarshalJSON(b)
}
