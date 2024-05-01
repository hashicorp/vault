// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginutil

import (
	"time"
)

const redactedTokenString = "ey***"

type IdentityTokenRequest struct {
	// Audience identifies the recipient of the token. The requested
	// value will be in the "aud" claim. Required.
	Audience string
	// TTL is the requested duration that the token will be valid for.
	// Optional with a default of 1hr.
	TTL time.Duration
}

type IdentityTokenResponse struct {
	// Token is the plugin identity token.
	Token IdentityToken
	// TTL is the duration that the token is valid for after truncation is applied.
	// The TTL may be truncated depending on the lifecycle of its signing key.
	TTL time.Duration
}

type IdentityToken string

// String returns a redacted token string. Use the Token() method
// to obtain the non-redacted token contents.
func (t IdentityToken) String() string {
	return redactedTokenString
}

// Token returns the non-redacted token contents.
func (t IdentityToken) Token() string {
	return string(t)
}
