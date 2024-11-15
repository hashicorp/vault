package protocol

import (
	"errors"
	"fmt"
)

// DecodeError represents a decoding error.
type DecodeError struct {
	row       int
	fieldName string
	err       error
}

func (e *DecodeError) Unwrap() error { return e.err }

func (e *DecodeError) Error() string {
	return fmt.Sprintf("decode error: %s row: %d fieldname: %s", e.err, e.row, e.fieldName)
}

// DecodeErrors represents a list of decoding errors.
type DecodeErrors []*DecodeError

func (errs DecodeErrors) rowErrors(row int) error {
	var rowErrs []error
	for _, err := range errs {
		if err.row == row {
			rowErrs = append(rowErrs, err)
		}
	}
	switch len(rowErrs) {
	case 0:
		return nil
	case 1:
		return rowErrs[0]
	default:
		return errors.Join(rowErrs...)
	}
}

// RowErrors returns errors if they were assigned to a row, nil otherwise.
func (errs DecodeErrors) RowErrors(row int) error {
	if len(errs) == 0 {
		return nil
	}
	return errs.rowErrors(row)
}
