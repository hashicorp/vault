package pgtype

import "database/sql/driver"

// Unknown represents the PostgreSQL unknown type. It is either a string literal
// or NULL. It is used when PostgreSQL does not know the type of a value. In
// general, this will only be used in pgx when selecting a null value without
// type information. e.g. SELECT NULL;
type Unknown struct {
	String string
	Status Status
}

func (dst *Unknown) Set(src interface{}) error {
	return (*Text)(dst).Set(src)
}

func (dst *Unknown) Get() interface{} {
	return (*Text)(dst).Get()
}

// AssignTo assigns from src to dst. Note that as Unknown is not a general number
// type AssignTo does not do automatic type conversion as other number types do.
func (src *Unknown) AssignTo(dst interface{}) error {
	return (*Text)(src).AssignTo(dst)
}

func (dst *Unknown) DecodeText(ci *ConnInfo, src []byte) error {
	return (*Text)(dst).DecodeText(ci, src)
}

func (dst *Unknown) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*Text)(dst).DecodeBinary(ci, src)
}

// Scan implements the database/sql Scanner interface.
func (dst *Unknown) Scan(src interface{}) error {
	return (*Text)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *Unknown) Value() (driver.Value, error) {
	return (*Text)(src).Value()
}
