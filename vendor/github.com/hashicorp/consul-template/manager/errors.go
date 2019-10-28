package manager

import "fmt"

// ErrExitable is an interface that defines an integer ExitStatus() function.
type ErrExitable interface {
	ExitStatus() int
}

var _ error = new(ErrChildDied)
var _ ErrExitable = new(ErrChildDied)

// ErrChildDied is the error returned when the child process prematurely dies.
type ErrChildDied struct {
	code int
}

// NewErrChildDied creates a new error with the given exit code.
func NewErrChildDied(c int) *ErrChildDied {
	return &ErrChildDied{code: c}
}

// Error implements the error interface.
func (e *ErrChildDied) Error() string {
	return fmt.Sprintf("child process died with exit code %d", e.code)
}

// ExitStatus implements the ErrExitable interface.
func (e *ErrChildDied) ExitStatus() int {
	return e.code
}
