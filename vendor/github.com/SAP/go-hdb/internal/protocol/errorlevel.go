// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

// ErrorLevel send from database server.
type errorLevel int8

func (e errorLevel) String() string {
	switch e {
	case 0:
		return "Warning"
	case 1:
		return "Error"
	case 2:
		return "Fatal Error"
	default:
		return ""
	}
}

// HDB error level constants.
const (
	errorLevelWarning    errorLevel = 0
	errorLevelError      errorLevel = 1
	errorLevelFatalError errorLevel = 2
)
