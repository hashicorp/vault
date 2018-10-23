// Package backoff is a utility for repeatedly retrying functions with support from a variety of backoff algorithms.
package backoff

/*
Backoff is an interface for any type that can implement a backoff algorithm
and maintain its current state.
*/
type Backoff interface {
	// Compute and return the next backoff delay.
	Next() bool
	// Retry a function until an error or backoff delay condition is met.
	Retry(func() error) error
	// Reset the backoff delay to its initial value.
	Reset()
}
