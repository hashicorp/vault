package pgtype

import (
	"database/sql/driver"
)

// GenericText is a placeholder for text format values that no other type exists
// to handle.
type GenericText Text

func (dst *GenericText) Set(src interface{}) error {
	return (*Text)(dst).Set(src)
}

func (dst GenericText) Get() interface{} {
	return (Text)(dst).Get()
}

func (src *GenericText) AssignTo(dst interface{}) error {
	return (*Text)(src).AssignTo(dst)
}

func (dst *GenericText) DecodeText(ci *ConnInfo, src []byte) error {
	return (*Text)(dst).DecodeText(ci, src)
}

func (src GenericText) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (Text)(src).EncodeText(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *GenericText) Scan(src interface{}) error {
	return (*Text)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src GenericText) Value() (driver.Value, error) {
	return (Text)(src).Value()
}
