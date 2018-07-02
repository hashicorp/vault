package jwt

import "errors"

var (
	// ErrTokenIsExpired is return when time.Now().Unix() is after
	// the token's "exp" claim.
	ErrTokenIsExpired = errors.New("token is expired")

	// ErrTokenNotYetValid is return when time.Now().Unix() is before
	// the token's "nbf" claim.
	ErrTokenNotYetValid = errors.New("token is not yet valid")

	// ErrInvalidISSClaim means the "iss" claim is invalid.
	ErrInvalidISSClaim = errors.New("claim \"iss\" is invalid")

	// ErrInvalidSUBClaim means the "sub" claim is invalid.
	ErrInvalidSUBClaim = errors.New("claim \"sub\" is invalid")

	// ErrInvalidIATClaim means the "iat" claim is invalid.
	ErrInvalidIATClaim = errors.New("claim \"iat\" is invalid")

	// ErrInvalidJTIClaim means the "jti" claim is invalid.
	ErrInvalidJTIClaim = errors.New("claim \"jti\" is invalid")

	// ErrInvalidAUDClaim means the "aud" claim is invalid.
	ErrInvalidAUDClaim = errors.New("claim \"aud\" is invalid")
)
