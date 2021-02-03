// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/protocol/encoding"
)

const (
	sqlStateSize = 5
	// bytes of fix length fields mod 8
	// - errorCode = 4, errorPosition = 4, errortextLength = 4, errorLevel = 1, sqlState = 5 => 18 bytes
	// - 18 mod 8 = 2
	fixLength = 2
)

type sqlState [sqlStateSize]byte

type hdbError struct {
	errorCode       int32
	errorPosition   int32
	errorTextLength int32
	errorLevel      errorLevel
	sqlState        sqlState
	stmtNo          int
	errorText       []byte
}

// String implements the Stringer interface.
func (e *hdbError) String() string {
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

// Error implements the Error interface.
func (e *hdbError) Error() string {
	if e.stmtNo != -1 {
		return fmt.Sprintf("SQL %s %d - %s (statement no: %d)", e.errorLevel, e.errorCode, e.errorText, e.stmtNo)
	}
	return fmt.Sprintf("SQL %s %d - %s", e.errorLevel, e.errorCode, e.errorText)
}

type hdbErrors struct {
	errors []*hdbError
	//numArg int
	idx int
}

// String implements the Stringer interface.
func (e *hdbErrors) String() string {
	return e.errors[e.idx].String()
}

// Error implements the golang error interface.
func (e *hdbErrors) Error() string {
	return e.errors[e.idx].Error()
}

// NumError implements the driver.Error interface.
func (e *hdbErrors) NumError() int {
	if e.errors == nil {
		return 0
	}
	return len(e.errors)
}

// SetIdx implements the driver.Error interface.
func (e *hdbErrors) SetIdx(idx int) {
	numError := e.NumError()
	switch {
	case idx < 0:
		e.idx = 0
	case idx >= numError:
		e.idx = numError - 1
	default:
		e.idx = idx
	}
}

// StmtNo implements the driver.Error interface.
func (e *hdbErrors) StmtNo() int {
	return e.errors[e.idx].stmtNo
}

// Code implements the driver.Error interface.
func (e *hdbErrors) Code() int {
	return int(e.errors[e.idx].errorCode)
}

// Position implements the driver.Error interface.
func (e *hdbErrors) Position() int {
	return int(e.errors[e.idx].errorPosition)
}

// Level implements the driver.Error interface.
func (e *hdbErrors) Level() int {
	return int(e.errors[e.idx].errorLevel)
}

// Text implements the driver.Error interface.
func (e *hdbErrors) Text() string {
	return string(e.errors[e.idx].errorText)
}

// IsWarning implements the driver.Error interface.
func (e *hdbErrors) IsWarning() bool {
	return e.errors[e.idx].errorLevel == errorLevelWarning
}

// IsError implements the driver.Error interface.
func (e *hdbErrors) IsError() bool {
	return e.errors[e.idx].errorLevel == errorLevelError
}

// IsFatal implements the driver.Error interface.
func (e *hdbErrors) IsFatal() bool {
	return e.errors[e.idx].errorLevel == errorLevelFatalError
}

func (e *hdbErrors) setStmtNo(idx, no int) {
	if idx >= 0 && idx < e.NumError() {
		e.errors[idx].stmtNo = no
	}
}

func (e *hdbErrors) isWarnings() bool {
	for _, _error := range e.errors {
		if _error.errorLevel != errorLevelWarning {
			return false
		}
	}
	return true
}

func (e *hdbErrors) reset(numArg int) {
	e.idx = 0 // init error index
	if e.errors == nil || numArg > cap(e.errors) {
		e.errors = make([]*hdbError, numArg)
	} else {
		e.errors = e.errors[:numArg]
	}
}

func (e *hdbErrors) decode(dec *encoding.Decoder, ph *partHeader) error {
	e.reset(ph.numArg())

	numArg := ph.numArg()
	for i := 0; i < numArg; i++ {
		_error := e.errors[i]
		if _error == nil {
			_error = new(hdbError)
			e.errors[i] = _error
		}

		_error.stmtNo = -1
		_error.errorCode = dec.Int32()
		_error.errorPosition = dec.Int32()
		_error.errorTextLength = dec.Int32()
		_error.errorLevel = errorLevel(dec.Int8())
		dec.Bytes(_error.sqlState[:])

		// read error text as ASCII data as some errors return invalid CESU-8 characters
		// e.g: SQL HdbError 7 - feature not supported: invalid character encoding: <invaid CESU-8 characters>
		//	if e.errorText, err = rd.ReadCesu8(int(e.errorTextLength)); err != nil {
		//		return err
		//	}
		_error.errorText = make([]byte, int(_error.errorTextLength))
		dec.Bytes(_error.errorText)

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

		pad := padBytes(int(fixLength + _error.errorTextLength))
		if pad != 0 {
			dec.Skip(pad)
		}
	}
	return dec.Error()
}
