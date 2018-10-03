package pgtype

import (
	"database/sql/driver"
)

// XID is PostgreSQL's Transaction ID type.
//
// In later versions of PostgreSQL, it is the type used for the backend_xid
// and backend_xmin columns of the pg_stat_activity system view.
//
// Also, when one does
//
//  select xmin, xmax, * from some_table;
//
// it is the data type of the xmin and xmax hidden system columns.
//
// It is currently implemented as an unsigned four byte integer.
// Its definition can be found in src/include/postgres_ext.h as TransactionId
// in the PostgreSQL sources.
type XID pguint32

// Set converts from src to dst. Note that as XID is not a general
// number type Set does not do automatic type conversion as other number
// types do.
func (dst *XID) Set(src interface{}) error {
	return (*pguint32)(dst).Set(src)
}

func (dst *XID) Get() interface{} {
	return (*pguint32)(dst).Get()
}

// AssignTo assigns from src to dst. Note that as XID is not a general number
// type AssignTo does not do automatic type conversion as other number types do.
func (src *XID) AssignTo(dst interface{}) error {
	return (*pguint32)(src).AssignTo(dst)
}

func (dst *XID) DecodeText(ci *ConnInfo, src []byte) error {
	return (*pguint32)(dst).DecodeText(ci, src)
}

func (dst *XID) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*pguint32)(dst).DecodeBinary(ci, src)
}

func (src *XID) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*pguint32)(src).EncodeText(ci, buf)
}

func (src *XID) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*pguint32)(src).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *XID) Scan(src interface{}) error {
	return (*pguint32)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *XID) Value() (driver.Value, error) {
	return (*pguint32)(src).Value()
}
