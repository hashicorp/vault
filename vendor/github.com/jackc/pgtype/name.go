package pgtype

import (
	"database/sql/driver"
)

// Name is a type used for PostgreSQL's special 63-byte
// name data type, used for identifiers like table names.
// The pg_class.relname column is a good example of where the
// name data type is used.
//
// Note that the underlying Go data type of pgx.Name is string,
// so there is no way to enforce the 63-byte length. Inputting
// a longer name into PostgreSQL will result in silent truncation
// to 63 bytes.
//
// Also, if you have custom-compiled PostgreSQL and set
// NAMEDATALEN to a different value, obviously that number of
// bytes applies, rather than the default 63.
type Name Text

func (dst *Name) Set(src interface{}) error {
	return (*Text)(dst).Set(src)
}

func (dst Name) Get() interface{} {
	return (Text)(dst).Get()
}

func (src *Name) AssignTo(dst interface{}) error {
	return (*Text)(src).AssignTo(dst)
}

func (dst *Name) DecodeText(ci *ConnInfo, src []byte) error {
	return (*Text)(dst).DecodeText(ci, src)
}

func (dst *Name) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*Text)(dst).DecodeBinary(ci, src)
}

func (src Name) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (Text)(src).EncodeText(ci, buf)
}

func (src Name) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (Text)(src).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *Name) Scan(src interface{}) error {
	return (*Text)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src Name) Value() (driver.Value, error) {
	return (Text)(src).Value()
}
