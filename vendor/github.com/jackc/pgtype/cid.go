package pgtype

import (
	"database/sql/driver"
)

// CID is PostgreSQL's Command Identifier type.
//
// When one does
//
// 	select cmin, cmax, * from some_table;
//
// it is the data type of the cmin and cmax hidden system columns.
//
// It is currently implemented as an unsigned four byte integer.
// Its definition can be found in src/include/c.h as CommandId
// in the PostgreSQL sources.
type CID pguint32

// Set converts from src to dst. Note that as CID is not a general
// number type Set does not do automatic type conversion as other number
// types do.
func (dst *CID) Set(src interface{}) error {
	return (*pguint32)(dst).Set(src)
}

func (dst CID) Get() interface{} {
	return (pguint32)(dst).Get()
}

// AssignTo assigns from src to dst. Note that as CID is not a general number
// type AssignTo does not do automatic type conversion as other number types do.
func (src *CID) AssignTo(dst interface{}) error {
	return (*pguint32)(src).AssignTo(dst)
}

func (dst *CID) DecodeText(ci *ConnInfo, src []byte) error {
	return (*pguint32)(dst).DecodeText(ci, src)
}

func (dst *CID) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*pguint32)(dst).DecodeBinary(ci, src)
}

func (src CID) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (pguint32)(src).EncodeText(ci, buf)
}

func (src CID) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (pguint32)(src).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *CID) Scan(src interface{}) error {
	return (*pguint32)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src CID) Value() (driver.Value, error) {
	return (pguint32)(src).Value()
}
