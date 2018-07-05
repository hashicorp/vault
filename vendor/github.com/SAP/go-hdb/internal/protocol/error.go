/*
Copyright 2014 SAP SE

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package protocol

import (
	"fmt"

	"github.com/SAP/go-hdb/internal/bufio"
)

const (
	sqlStateSize = 5
	//bytes of fix length fields mod 8
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
	return fmt.Sprintf("errorCode %d, errorPosition %d, errorTextLength % d errorLevel %s, sqlState %s stmtNo %d errorText %s",
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
	numArg int
	idx    int
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
	return e.numArg
}

// SetIdx implements the driver.Error interface.
func (e *hdbErrors) SetIdx(idx int) {
	switch {
	case idx < 0:
		e.idx = 0
	case idx >= e.numArg:
		e.idx = e.numArg - 1
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
	if idx >= 0 && idx < e.numArg {
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

func (e *hdbErrors) kind() partKind {
	return pkError
}

func (e *hdbErrors) setNumArg(numArg int) {
	e.numArg = numArg
}

func (e *hdbErrors) read(rd *bufio.Reader) error {
	e.idx = 0 // init error index

	if e.errors == nil || e.numArg > cap(e.errors) {
		e.errors = make([]*hdbError, e.numArg)
	} else {
		e.errors = e.errors[:e.numArg]
	}

	for i := 0; i < e.numArg; i++ {
		_error := e.errors[i]
		if _error == nil {
			_error = new(hdbError)
			e.errors[i] = _error
		}

		_error.stmtNo = -1
		_error.errorCode = rd.ReadInt32()
		_error.errorPosition = rd.ReadInt32()
		_error.errorTextLength = rd.ReadInt32()
		_error.errorLevel = errorLevel(rd.ReadInt8())
		rd.ReadFull(_error.sqlState[:])

		// read error text as ASCII data as some errors return invalid CESU-8 characters
		// e.g: SQL HdbError 7 - feature not supported: invalid character encoding: <invaid CESU-8 characters>
		//	if e.errorText, err = rd.ReadCesu8(int(e.errorTextLength)); err != nil {
		//		return err
		//	}
		_error.errorText = make([]byte, int(_error.errorTextLength))
		rd.ReadFull(_error.errorText)

		if trace {
			outLogger.Printf("error %d: %s", i, _error)
		}

		if e.numArg == 1 {
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
			rd.Skip(1)
			break
		}

		pad := padBytes(int(fixLength + _error.errorTextLength))
		if pad != 0 {
			rd.Skip(pad)
		}
	}

	return rd.GetError()
}
