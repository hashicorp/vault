package pgtype

import (
	"database/sql/driver"

	"github.com/pkg/errors"
)

type JSONB JSON

func (dst *JSONB) Set(src interface{}) error {
	return (*JSON)(dst).Set(src)
}

func (dst *JSONB) Get() interface{} {
	return (*JSON)(dst).Get()
}

func (src *JSONB) AssignTo(dst interface{}) error {
	return (*JSON)(src).AssignTo(dst)
}

func (dst *JSONB) DecodeText(ci *ConnInfo, src []byte) error {
	return (*JSON)(dst).DecodeText(ci, src)
}

func (dst *JSONB) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = JSONB{Status: Null}
		return nil
	}

	if len(src) == 0 {
		return errors.Errorf("jsonb too short")
	}

	if src[0] != 1 {
		return errors.Errorf("unknown jsonb version number %d", src[0])
	}

	*dst = JSONB{Bytes: src[1:], Status: Present}
	return nil

}

func (src *JSONB) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (*JSON)(src).EncodeText(ci, buf)
}

func (src *JSONB) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}

	buf = append(buf, 1)
	return append(buf, src.Bytes...), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *JSONB) Scan(src interface{}) error {
	return (*JSON)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src *JSONB) Value() (driver.Value, error) {
	return (*JSON)(src).Value()
}
