package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/driver/internal/protocol/encoding"
)

// ErrorLevel send from database server.
type errorLevel int8

var errorLevelStrs = [...]string{"Warning", "Error", "FatalError"}

func (e errorLevel) String() string {
	if int(e) < 0 || int(e) >= len(errorLevelStrs) {
		return ""
	}
	return errorLevelStrs[e]
}

// HDB error level constants.
const (
	errorLevelWarning    errorLevel = 0
	errorLevelError      errorLevel = 1
	errorLevelFatalError errorLevel = 2
)

const (
	sqlStateSize = 5
	/*
		bytes of fix length fields mod 8
		  - errorCode = 4, errorPosition = 4, errortextLength = 4, errorLevel = 1, sqlState = 5 => 18 bytes
		  - 18 mod 8 = 2
	*/
	fixLength = 2
)

// HANA Database errors.
const (
	HdbErrAuthenticationFailed = 10
	HdbErrWhileParsingProtocol = 1033
)

type sqlState [sqlStateSize]byte

// HdbError represents a single error returned by the server.
type HdbError struct {
	errorCode       int32
	errorPosition   int32
	errorTextLength int32
	errorLevel      errorLevel
	sqlState        sqlState
	stmtNo          int
	errorText       []byte
}

func (e *HdbError) String() string {
	return fmt.Sprintf("errorCode %d errorPosition %d errorTextLength %d errorLevel %s sqlState %s stmtNo %d errorText %s",
		e.errorCode,
		e.errorPosition,
		e.errorTextLength,
		e.errorLevel,
		e.sqlState,
		e.stmtNo,
		e.errorText,
	)
}

func (e *HdbError) Error() string {
	if e.stmtNo != -1 {
		return fmt.Sprintf("SQL %s %d - %s (statement no: %d)", e.errorLevel, e.errorCode, e.errorText, e.stmtNo)
	}
	return fmt.Sprintf("SQL %s %d - %s", e.errorLevel, e.errorCode, e.errorText)
}

// StmtNo implements the driver.DBError interface.
func (e *HdbError) StmtNo() int { return e.stmtNo }

// Code implements the driver.DBError interface.
func (e *HdbError) Code() int { return int(e.errorCode) }

// Position implements the driver.DBError interface.
func (e *HdbError) Position() int { return int(e.errorPosition) }

// Level implements the driver.DBError interface.
func (e *HdbError) Level() int { return int(e.errorLevel) }

// Text implements the driver.DBError interface.
func (e *HdbError) Text() string { return string(e.errorText) }

// IsWarning implements the driver.DBError interface.
func (e *HdbError) IsWarning() bool { return e.errorLevel == errorLevelWarning }

// IsError implements the driver.DBError interface.
func (e *HdbError) IsError() bool { return e.errorLevel == errorLevelError }

// IsFatal implements the driver.DBError interface.
func (e *HdbError) IsFatal() bool { return e.errorLevel == errorLevelFatalError }

// HdbErrors represent the collection of errors return by the server.
type HdbErrors struct {
	onlyWarnings bool
	errs         []*HdbError
	*HdbError
}

func (e *HdbErrors) String() string {
	var b []byte
	for i, err := range e.errs {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.String()...)
	}
	return string(b)
}

func (e *HdbErrors) Error() string {
	var b []byte
	for i, err := range e.errs {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.Error()...)
	}
	return string(b)
}

// NumError implements the driver.Error interface.
// NumErrors returns the number of all errors, including warnings.
func (e *HdbErrors) NumError() int { return len(e.errs) }

func (e *HdbErrors) Unwrap() []error {
	errs := make([]error, 0, len(e.errs))
	for _, err := range e.errs {
		errs = append(errs, err)
	}
	return errs
}

// SetIdx implements the driver.Error interface.
func (e *HdbErrors) SetIdx(idx int) {
	if idx >= 0 && idx < len(e.errs) {
		e.HdbError = e.errs[idx]
	}
}

// setStmtNo sets the statement number of the error.
func (e *HdbErrors) setStmtNo(idx, no int) {
	if idx >= 0 && idx < len(e.errs) {
		e.errs[idx].stmtNo = no
	}
}

func (e *HdbErrors) decodeNumArg(dec *encoding.Decoder, numArg int) error {
	e.onlyWarnings = true
	e.errs = nil

	for range numArg {
		err := new(HdbError)
		e.errs = append(e.errs, err)

		// err.stmtNo = -1
		err.stmtNo = 0
		/*
			in case of an hdb error when inserting one record (e.g. duplicate)
			- hdb does not return a rowsAffected part
			- SetStmtNo is not called and
			- the default value (formerly -1) is kept
			--> initialize stmtNo with zero
		*/
		err.errorCode = dec.Int32()
		err.errorPosition = dec.Int32()
		err.errorTextLength = dec.Int32()
		err.errorLevel = errorLevel(dec.Int8())
		dec.Bytes(err.sqlState[:])

		// read error text as ASCII data as some errors return invalid CESU-8 characters
		// e.g: SQL HdbError 7 - feature not supported: invalid character encoding: <invaid CESU-8 characters>
		//	if e.errorText, err = rd.ReadCesu8(int(e.errorTextLength)); err != nil {
		//		return err
		//	}
		err.errorText = make([]byte, int(err.errorTextLength))
		dec.Bytes(err.errorText)

		if e.onlyWarnings && !err.IsWarning() {
			e.onlyWarnings = false
		}

		if numArg == 1 {
			// Error (protocol error?):
			// if only one error (numArg == 1): s.ph.bufferLength is one byte greater than data to be read
			// if more than one error: s.ph.bufferlength matches read bytes + padding
			//
			// Examples:
			// driver test TestHDBWarning
			//   --> 18 bytes fix error bytes + 103 bytes error text => 121 bytes (7 bytes padding needed)
			//   but s.ph.bufferLength = 122 (standard padding would only consume 6 bytes instead of 7)
			// driver test TestBulkInsertDuplicates
			//   --> returns 3 errors (number of total bytes matches s.ph.bufferLength)
			dec.Skip(1)
			break
		}

		pad := padBytes(int(fixLength + err.errorTextLength))
		if pad != 0 {
			dec.Skip(pad)
		}
	}
	if len(e.errs) > 0 {
		e.HdbError = e.errs[0] // set default to first error
	}

	return dec.Error()
}
