package oidc

import "encoding/json"

//  AccessToken is an oauth access_token.
type AccessToken string

// RedactedAccessToken is the redacted string or json for an oauth access_token.
const RedactedAccessToken = "[REDACTED: access_token]"

// String will redact the token.
func (t AccessToken) String() string {
	return RedactedAccessToken
}

// MarshalJSON will redact the token.
func (t AccessToken) MarshalJSON() ([]byte, error) {
	return json.Marshal(RedactedAccessToken)
}
