package pgtype

import (
	"database/sql/driver"
)

// BPChar is fixed-length, blank padded char type
// character(n), char(n)
type BPChar Text

// Set converts from src to dst.
func (dst *BPChar) Set(src interface{}) error {
	return (*Text)(dst).Set(src)
}

// Get returns underlying value
func (dst *BPChar) Get() interface{} {
	return (*Text)(dst).Get()
}

// AssignTo assigns from src to dst.
func (src *BPChar) AssignTo(dst interface{}) error {
	if src.Status == Present {
		switch v := dst.(type) {
		case *rune:
			runes := []rune(src.String)
			if len(runes) == 1 {
				*v = runes[0]
				return nil
			}
		}
	}
	return (*Text)(src).AssignTo(dst)
}

func (dst *BPChar) DecodeText(ci *ConnInfo, src []byte) error {
	return (*Text)(dst).DecodeText(ci, src)
}

func (dst *BPChar) DecodeBinary(ci *ConnInfo, src []byte) error {
	return (*Text)(dst).DecodeBinary(ci, src)
}

func (src *BPChar) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Text)(src).EncodeText(ci, buf)
}

func (src *BPChar) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*Text)(src).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *BPChar) Scan(src interface{}) error {
	return (*Text)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *BPChar) Value() (driver.Value, error) {
	return (*Text)(src).Value()
}

func (src *BPChar) MarshalJSON() ([]byte, error) {
	return (*Text)(src).MarshalJSON()
}

func (dst *BPChar) UnmarshalJSON(b []byte) error {
	return (*Text)(dst).UnmarshalJSON(b)
}
