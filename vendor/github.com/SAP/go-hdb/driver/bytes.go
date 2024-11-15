package driver

import (
	"database/sql/driver"
)

// NullBytes represents an []byte that may be null.
// NullBytes implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullBytes struct {
	Bytes []byte
	Valid bool // Valid is true if Bytes is not NULL
}

// Scan implements the Scanner interface.
func (n *NullBytes) Scan(value any) error {
	n.Bytes, n.Valid = value.([]byte)
	return nil
}

// Value implements the driver Valuer interface.
func (n NullBytes) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bytes, nil
}
