package limiter

import (
	"context"
	"fmt"
	"time"
)

// ErrStopped is the error returned when the store is stopped. All implementers
// should return this error for stoppable stores.
var ErrStopped = fmt.Errorf("store is stopped")

// Store is an interface for limiter storage backends.
//
// Keys should be hash, sanitized, or otherwise scrubbed of identifiable
// information they will be given to the store in plaintext. If you're rate
// limiting by IP address, for example, the IP address would be stored in the
// storage system in plaintext. This may be undesirable in certain situations,
// like when the store is a public database. In those cases, you should hash or
// HMAC the key before passing giving it to the store. If you want to encrypt
// the value, you must use homomorphic encryption to ensure the value always
// encrypts to the same ciphertext.
type Store interface {
	// Take takes a token from the given key if available, returning:
	//
	// - the configured limit size
	// - the number of remaining tokens in the interval
	// - the server time when new tokens will be available
	// - whether the take was successful
	// - any errors that occurred while performing the take - these should be
	//   backend errors (e.g. connection failures); Take() should never return an
	//   error for an bucket.
	//
	// If "ok" is false, the take was unsuccessful and the caller should NOT
	// service the request.
	//
	// See the note about keys on the interface documentation.
	Take(ctx context.Context, key string) (tokens, remaining, reset uint64, ok bool, err error)

	// Get gets the current limit and remaining tokens for the provided key. It
	// does not change any of the values.
	Get(ctx context.Context, key string) (tokens, remaining uint64, err error)

	// Set configures the limit at the provided key. If a limit already exists, it
	// is overwritten. This also sets the number of tokens in the bucket to the
	// limit.
	Set(ctx context.Context, key string, tokens uint64, interval time.Duration) error

	// Burst adds more tokens to the key's current bucket until the next interval
	// tick. This will allow the current bucket tick to exceed the maximum number
	// maximum ticks until the next interval.
	Burst(ctx context.Context, key string, tokens uint64) error

	// Close terminates the store and cleans up any data structures or connections
	// that may remain open. After a store is stopped, Take() should always return
	// zero values.
	Close(ctx context.Context) error
}
