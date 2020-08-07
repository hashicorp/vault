package dependency

import "errors"

// ErrStopped is a special error that is returned when a dependency is
// prematurely stopped, usually due to a configuration reload or a process
// interrupt.
var ErrStopped = errors.New("dependency stopped")

// ErrContinue is a special error which says to continue (retry) on error.
var ErrContinue = errors.New("dependency continue")

var ErrLeaseExpired = errors.New("lease expired or is not renewable")
