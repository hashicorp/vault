// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginutil

import (
	"time"
)

type IdentityTokenRequest struct {
	// Key is the named identity token key to sign the token with
	Key string
	// Audience identifies the recipient of the token
	Audience string
	// TTL is the duration that the token will be valid for
	TTL time.Duration
}

type IdentityTokenResponse struct {
	// Token is the plugin identity token
	Token IdentityToken
	// TTL is the capped duration that the token is valid for
	TTL time.Duration
}

type IdentityToken string

// String returns a redacted token string. Use the Token() method
// to obtain the non-redacted token contents.
func (t IdentityToken) String() string {
	return "ey***"
}

// Token returns the non-redacted token contents.
func (t IdentityToken) Token() string {
	return string(t)
}
