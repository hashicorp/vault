package pgtype

import (
	"database/sql/driver"
	"fmt"
)

type Ltree Text

func (dst *Ltree) Set(src interface{}) error {
	return (*Text)(dst).Set(src)
}

func (dst Ltree) Get() interface{} {
	return (Text)(dst).Get()
}

func (src *Ltree) AssignTo(dst interface{}) error {
	return (*Text)(src).AssignTo(dst)
}

func (src Ltree) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (Text)(src).EncodeText(ci, buf)
}

func (src Ltree) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case Null:
		return nil, nil
	case Undefined:
		return nil, errUndefined
	}
	buf = append(buf, 1)
	return append(buf, src.String...), nil
}

func (Ltree) PreferredResultFormat() int16 {
	return TextFormatCode
}

func (dst *Ltree) DecodeText(ci *ConnInfo, src []byte) error {
	return (*Text)(dst).DecodeText(ci, src)
}

func (dst *Ltree) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = Ltree{Status: Null}
		return nil
	}

	// Get Ltree version, only 1 is allowed
	version := src[0]
	if version != 1 {
		return fmt.Errorf("unsupported ltree version %d", version)
	}

	ltreeStr := string(src[1:])
	*dst = Ltree{String: ltreeStr, Status: Present}
	return nil
}

func (Ltree) PreferredParamFormat() int16 {
	return TextFormatCode
}

func (dst *Ltree) Scan(src interface{}) error {
	return (*Text)(dst).Scan(src)
}

func (src Ltree) Value() (driver.Value, error) {
	return (Text)(src).Value()
}
