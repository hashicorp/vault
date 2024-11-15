// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Nonce is an interface for generating and validating nonces. Different
// backend implementations have different performance and security
// characteristics.

package nonceutil

import (
	"time"
)

// NonceService is an interface for issuing and redeeming nonces, with
// a hook to periodically free resources when no redemptions have happened
// recently.
//
// A nonce is a unique token that can be given to a client, who can later
// "redeem" or use that token on a subsequent request to prove that the
// request has only been done once. No tracking of client->token is performed
// as part of this service. For an example use of nonces within a protocol,
// see IETF RFC 8555 Automatic Certificate Management Environment (ACME).
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
