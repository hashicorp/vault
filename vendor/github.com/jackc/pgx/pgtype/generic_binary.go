package pgtype

import (
	"database/sql/driver"
)

// GenericBinary is a placeholder for binary format values that no other type exists
// to handle.
type GenericBinary Bytea

func (dst *GenericBinary) Set(src interface{}) error {
	return (*Bytea)(dst).Set(src)
}

func (dst *GenericBinary) Get() interface{} {
	return (*Bytea)(dst).Get()
}

func (src *GenericBinary) AssignTo(dst interface{}) error {
	return (*Bytea)(src).AssignTo(dst)
}

func (dst *GenericBinary) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*Bytea)(dst).DecodeBinary(ci, src)
}

func (src *GenericBinary) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Bytea)(src).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *GenericBinary) Scan(src interface{}) error {
	return (*Bytea)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *GenericBinary) Value() (driver.Value, error) {
	return (*Bytea)(src).Value()
}
