package pgtype

import (
	"database/sql/driver"
	"errors"
)

func DatabaseSQLValue(ci *ConnInfo, src Value) (interface{}, error) {
	if valuer, ok := src.(driver.Valuer); ok {
		return valuer.Value()
	}

	if textEncoder, ok := src.(TextEncoder); ok {
		buf, err := textEncoder.EncodeText(ci, nil)
		if err != nil {
			return nil, err
		}
		return string(buf), nil
	}

	if binaryEncoder, ok := src.(BinaryEncoder); ok {
		buf, err := binaryEncoder.EncodeBinary(ci, nil)
		if err != nil {
			return nil, err
		}
		return buf, nil
	}

	return nil, errors.New("cannot convert to database/sql compatible value")
}

func EncodeValueText(src TextEncoder) (interface{}, error) {
	var encBuf [36]byte
	buf, err := src.EncodeText(nil, encBuf[:0])
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}
	return string(buf), err
}
