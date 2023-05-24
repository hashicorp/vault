// Nonce is a class for generating and validating nonces loosely based off
// the design of Let's Encrypt's Boulder nonce service here:
//
//   https://github.com/letsencrypt/boulder/blob/main/nonce/nonce.go

package nonce

import (
	"time"
)

// NonceService is an interface for issuing and redeeming nonces, with
// a hook to periodically free resources when no redemptions have happened
// recently.
//
// Notably, nonces are not guaranteed to be stored or persisted; nonces
// from one startup will not necessarily be valid from another.
type NonceService interface {
	// Before using a nonce service, it must be initialized. Failure to
	// initialize might result in panics or other unexpected results.
	Initialize() error

	// Get a nonce; returns three values:
	//
	// 1. The nonce itself, a base64-url-no-padding opaque value.
	// 2. A time at which the nonce will expire, based on the validity
	//    period specified at construction. By default, the service issues
	//    short-lived nonces.
	// 3. An error if one occurred during generation of the nonce.
	Get() (string, time.Time, error)

	// Redeem the given nonce, returning whether or not it was accepted. A
	// nonce given twice will be rejected if the service is a strict nonce
	// service, but potentially accepted if the nonce service is loose
	// (i.e., temporal revocation only).
	Redeem(string) bool

	// A hook to tidy the memory usage of the underlying implementation; is
	// implementation dependent. Some implementations may not return status
	// information.
	Tidy() *NonceStatus

	// If true, this is a strict only-once redemption service implementation,
	// else a nonce could be accepted more than once within some safety
	// window.
	IsStrict() bool

	// Whether or not this service is usable across nodes.
	IsCrossNode() bool
}

func NewNonceService() NonceService {
	// By default, we create an encrypted nonce service that is strict but not
	// cross node, using a default window of 90 seconds (equal to the default
	// context request timeout window).
	return NewNonceServiceWithValidity(90 * time.Second)
}

func NewNonceServiceWithValidity(validity time.Duration) NonceService {
	return newEncryptedNonceService(validity)
}

// Status information about the number of nonces in this service, perhaps
// local to this node. Presumably, the delta roughly correlates to present
// memory usage.
type NonceStatus struct {
	Issued      uint64
	Outstanding uint64
	Message     string
}
