package limiter

import "io"

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
	//
	// If "ok" is false, the take was unsuccessful and the caller should NOT
	// service the request.
	//
	// See the note about keys on the interface documentation.
	Take(key string) (limit, remaining, reset uint64, ok bool)

	// Close terminates the store and cleans up any data structures or connections
	// that may remain open. After a store is stopped, Take() should always return
	// zero values.
	io.Closer
}
