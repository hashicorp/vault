package oidc

import "encoding/json"

// RefreshToken is an oauth refresh_token.
// See https://tools.ietf.org/html/rfc6749#section-1.5.
type RefreshToken string

// RedactedRefreshToken is the redacted string or json for an oauth refresh_token.
const RedactedRefreshToken = "[REDACTED: refresh_token]"

// String will redact the token.
func (t RefreshToken) String() string {
	return RedactedRefreshToken
}

// MarshalJSON will redact the token.
func (t RefreshToken) MarshalJSON() ([]byte, error) {
	return json.Marshal(RedactedRefreshToken)
}
