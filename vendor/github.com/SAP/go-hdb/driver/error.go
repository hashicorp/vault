// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package driver

// HDB error levels.
const (
	HdbWarning    = 0
	HdbError      = 1
	HdbFatalError = 2
)

// Error represents errors send by the database server.
type Error interface {
	Error() string   // Implements the golang error interface.
	NumError() int   // NumError returns the number of errors.
	SetIdx(idx int)  // Sets the error index in case number of errors are greater 1 in the range of 0 <= index < NumError().
	StmtNo() int     // Returns the statement number of the error in multi statement contexts (e.g. bulk insert).
	Code() int       // Code return the database error code.
	Position() int   // Position returns the start position of erroneous sql statements sent to the database server.
	Level() int      // Level return one of the database server predefined error levels.
	Text() string    // Text return the error description sent from database server.
	IsWarning() bool // IsWarning returns true if the HDB error level equals 0.
	IsError() bool   // IsError returns true if the HDB error level equals 1.
	IsFatal() bool   // IsFatal returns true if the HDB error level equals 2.
}
