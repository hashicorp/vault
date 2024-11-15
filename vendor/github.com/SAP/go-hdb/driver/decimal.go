package driver

import (
	"database/sql/driver"
	"fmt"
	"math/big"
)

// A Decimal is the driver representation of a database decimal field value as big.Rat.
type Decimal big.Rat

// Scan implements the database/sql/Scanner interface.
func (d *Decimal) Scan(src any) error {
	r, ok := src.(*big.Rat)
	if !ok {
		return fmt.Errorf("decimal: invalid data type %T", src)
	}
	(*big.Rat)(d).Set(r)
	return nil
}

// Value implements the database/sql/Valuer interface.
func (d Decimal) Value() (driver.Value, error) {
	return (*big.Rat)(&d), nil
}

// NullDecimal represents an Decimal that may be null.
// NullDecimal implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullDecimal struct {
	Decimal *Decimal
	Valid   bool // Valid is true if Decimal is not NULL
}

// Scan implements the Scanner interface.
func (n *NullDecimal) Scan(value any) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	r, ok := value.(*big.Rat)
	if !ok {
		return fmt.Errorf("decimal: invalid data type %T", value)
	}
	n.Valid = true
	if n.Decimal == nil {
		n.Decimal = &Decimal{}
	}
	(*big.Rat)(n.Decimal).Set(r)
	return nil
}

// Value implements the driver Valuer interface.
func (n NullDecimal) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	if n.Decimal == nil {
		return nil, fmt.Errorf("invalid decimal value %v", n.Decimal)
	}
	return (*big.Rat)(n.Decimal), nil
}
