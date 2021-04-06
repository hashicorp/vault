package oidc

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"

	"gopkg.in/square/go-jose.v2"
)

// IDToken is an oidc id_token.
// See https://openid.net/specs/openid-connect-core-1_0.html#IDToken.
type IDToken string

// RedactedIDToken is the redacted string or json for an oidc id_token.
const RedactedIDToken = "[REDACTED: id_token]"

// String will redact the token.
func (t IDToken) String() string {
	return RedactedIDToken
}

// MarshalJSON will redact the token.
func (t IDToken) MarshalJSON() ([]byte, error) {
	return json.Marshal(RedactedIDToken)
}

// Claims retrieves the IDToken claims.
func (t IDToken) Claims(claims interface{}) error {
	const op = "IDToken.Claims"
	if len(t) == 0 {
		return fmt.Errorf("%s: id_token is empty: %w", op, ErrInvalidParameter)
	}
	if claims == nil {
		return fmt.Errorf("%s: claims interface is nil: %w", op, ErrNilParameter)
	}
	return UnmarshalClaims(string(t), claims)
}

// VerifyAccessToken verifies the at_hash claim of the id_token against the hash
// of the access_token.
//
// It will return true when it can verify the access_token. It will return false
// when it's unable to verify the access_token.
//
// It will return an error whenever it's possible to verify the access_token and
// the verification fails.
//
// Note: while we support signing id_tokens with EdDSA, unfortunately the
// access_token hash cannot be verified without knowing the key's curve. See:
// https://bitbucket.org/openid/connect/issues/1125
//
// For more info about verifying access_tokens returned during an OIDC flow see:
// https://openid.net/specs/openid-connect-core-1_0.html#CodeIDToken
func (t IDToken) VerifyAccessToken(accessToken AccessToken) (bool, error) {
	const op = "VerifyAccessToken"
	canVerify, err := t.verifyHashClaim("at_hash", string(accessToken))
	if err != nil {
		return canVerify, fmt.Errorf("%s: %w", op, err)
	}
	return canVerify, nil
}

// VerifyAuthorizationCode verifies the c_hash claim of the id_token against the
// hash of the authorization code.
//
// It will return true when it can verify the authorization code. It will return
// false when it's unable to verify the authorization code.
//
// It will return an error whenever it's possible to verify the authorization
// code and the verification fails.
//
// Note: while we support signing id_tokens with EdDSA, unfortunately the
// authorization code hash cannot be verified without knowing the key's curve.
// See: https://bitbucket.org/openid/connect/issues/1125
//
// For more info about authorization code verification using the id_token's
// c_hash claim see:
// https://openid.net/specs/openid-connect-core-1_0.html#HybridIDToken
func (t IDToken) VerifyAuthorizationCode(code string) (bool, error) {
	const op = "VerifyAccessToken"
	canVerify, err := t.verifyHashClaim("c_hash", code)
	if err != nil {
		return canVerify, fmt.Errorf("%s: %w", op, err)
	}
	return canVerify, nil
}

func (t IDToken) verifyHashClaim(claimName string, token string) (bool, error) {
	const op = "verifyHashClaim"
	var claims map[string]interface{}
	if err := t.Claims(&claims); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	tokenHash, ok := claims[claimName].(string)
	if !ok {
		return false, nil
	}

	jws, err := jose.ParseSigned(string(t))
	if err != nil {
		return false, fmt.Errorf("%s: malformed jwt (%v): %w", op, err, ErrMalformedToken)
	}
	switch len(jws.Signatures) {
	case 0:
		return false, fmt.Errorf("%s: id_token not signed: %w", op, ErrTokenNotSigned)
	case 1:
	default:
		return false, fmt.Errorf("%s: multiple signatures on id_token not supported", op)
	}
	sig := jws.Signatures[0]
	if _, ok := supportedAlgorithms[Alg(sig.Header.Algorithm)]; !ok {
		return false, fmt.Errorf("%s: id_token signed with algorithm %q: %w", op, sig.Header.Algorithm, ErrUnsupportedAlg)
	}
	sigAlgorithm := Alg(sig.Header.Algorithm)

	var h hash.Hash
	switch sigAlgorithm {
	case RS256, ES256, PS256:
		h = sha256.New()
	case RS384, ES384, PS384:
		h = sha512.New384()
	case RS512, ES512, PS512:
		h = sha512.New()
	case EdDSA:
		return false, nil
	default:
		return false, fmt.Errorf("%s: unsupported signing algorithm %s: %w", op, sigAlgorithm, ErrUnsupportedAlg)
	}
	_, _ = h.Write([]byte(token)) // hash documents that Write will never return an error
	sum := h.Sum(nil)[:h.Size()/2]
	actual := base64.RawURLEncoding.EncodeToString(sum)
	if actual != tokenHash {
		switch claimName {
		case "at_hash":
			return false, fmt.Errorf("%s: %w", op, ErrInvalidAtHash)
		case "c_hash":
			return false, fmt.Errorf("%s: %w", op, ErrInvalidCodeHash)
		}
	}
	return true, nil
}
