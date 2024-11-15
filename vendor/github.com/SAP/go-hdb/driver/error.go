package driver

import (
	p "github.com/SAP/go-hdb/driver/internal/protocol"
)

// HDB error levels.
const (
	HdbWarning    = 0
	HdbError      = 1
	HdbFatalError = 2
)

// DBError represents a single error returned by the database server.
type DBError interface {
	Error() string   // Implements the golang error interface.
	StmtNo() int     // Returns the statement number of the error in multi statement contexts (e.g. bulk insert).
	Code() int       // Code return the database error code.
	Position() int   // Position returns the start position of erroneous sql statements sent to the database server.
	Level() int      // Level return one of the database server predefined error levels.
	Text() string    // Text return the error description sent from database server.
	IsWarning() bool // IsWarning returns true if the HDB error level equals 0.
	IsError() bool   // IsError returns true if the HDB error level equals 1.
	IsFatal() bool   // IsFatal returns true if the HDB error level equals 2.
}

// Error represents errors (an error collection) send by the database server.
type Error interface {
	Error() string   // Implements the golang error interface.
	NumError() int   // NumError returns the number of errors.
	Unwrap() []error // Unwrap implements the standard error Unwrap function for errors wrapping multiple errors.
	SetIdx(idx int)  // SetIdx sets the error index in case number of errors are greater 1 in the range of 0 <= index < NumError().
	DBError          // DBError functions for error in case of single error, for error set by SetIdx in case of error collection.
}

var (
	_ DBError = (*p.HdbError)(nil)
	_ Error   = (*p.HdbErrors)(nil)
)
