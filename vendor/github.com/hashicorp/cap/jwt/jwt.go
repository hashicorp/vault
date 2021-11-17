package jwt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// DefaultLeewaySeconds defines the amount of leeway that's used by default
// for validating the "nbf" (Not Before) and "exp" (Expiration Time) claims.
const DefaultLeewaySeconds = 150

// Validator validates JSON Web Tokens (JWT) by providing signature
// verification and claims set validation.
type Validator struct {
	keySet KeySet
}

// NewValidator returns a Validator that uses the given KeySet to verify JWT signatures.
func NewValidator(keySet KeySet) (*Validator, error) {
	if keySet == nil {
		return nil, errors.New("keySet must not be nil")
	}

	return &Validator{
		keySet: keySet,
	}, nil
}

// Expected defines the expected claims values to assert when validating a JWT.
// For claims that involve validation of the JWT with respect to time, leeway
// fields are provided to account for potential clock skew.
type Expected struct {
	// The expected JWT "iss" (issuer) claim value. If empty, validation is skipped.
	Issuer string

	// The expected JWT "sub" (subject) claim value. If empty, validation is skipped.
	Subject string

	// The expected JWT "jti" (JWT ID) claim value. If empty, validation is skipped.
	ID string

	// The list of expected JWT "aud" (audience) claim values to match against.
	// The JWT claim will be considered valid if it matches any of the expected
	// audiences. If empty, validation is skipped.
	Audiences []string

	// SigningAlgorithms provides the list of expected JWS "alg" (algorithm) header
	// parameter values to match against. The JWS header parameter will be considered
	// valid if it matches any of the expected signing algorithms. The following
	// algorithms are supported: RS256, RS384, RS512, ES256, ES384, ES512, PS256,
	// PS384, PS512, EdDSA. If empty, defaults to RS256.
	SigningAlgorithms []Alg

	// NotBeforeLeeway provides the option to set an amount of leeway to use when
	// validating the "nbf" (Not Before) claim. If the duration is zero or not
	// provided, a default leeway of 150 seconds will be used. If the duration is
	// negative, no leeway will be used.
	NotBeforeLeeway time.Duration

	// ExpirationLeeway provides the option to set an amount of leeway to use when
	// validating the "exp" (Expiration Time) claim. If the duration is zero or not
	// provided, a default leeway of 150 seconds will be used. If the duration is
	// negative, no leeway will be used.
	ExpirationLeeway time.Duration

	// ClockSkewLeeway provides the option to set an amount of leeway to use when
	// validating the "nbf" (Not Before), "exp" (Expiration Time), and "iat" (Issued At)
	// claims. If the duration is zero or not provided, a default leeway of 60 seconds
	// will be used. If the duration is negative, no leeway will be used.
	ClockSkewLeeway time.Duration

	// Now provides the option to specify a func for determining what the current time is.
	// The func will be used to provide the current time when validating a JWT with respect to
	// the "nbf" (Not Before), "exp" (Expiration Time), and "iat" (Issued At) claims. If not
	// provided, defaults to returning time.Now().
	Now func() time.Time
}

// Validate validates JWTs of the JWS compact serialization form.
//
// The given JWT is considered valid if:
//  1. Its signature is successfully verified.
//  2. Its claims set and header parameter values match what's given by Expected.
//  3. It's valid with respect to the current time. This means that the current
//     time must be within the times (inclusive) given by the "nbf" (Not Before)
//     and "exp" (Expiration Time) claims and after the time given by the "iat"
//     (Issued At) claim, with configurable leeway. See Expected.Now() for details
//     on how the current time is provided for validation.
func (v *Validator) Validate(ctx context.Context, token string, expected Expected) (map[string]interface{}, error) {
	// First, verify the signature to ensure subsequent validation is against verified claims
	allClaims, err := v.keySet.VerifySignature(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("error verifying token signature: %w", err)
	}

	// Validate the signing algorithm in the JWS header
	if err := validateSigningAlgorithm(token, expected.SigningAlgorithms); err != nil {
		return nil, fmt.Errorf("invalid algorithm (alg) header parameter: %w", err)
	}

	// Unmarshal all claims into the set of public JWT registered claims
	claims := jwt.Claims{}
	allClaimsJSON, err := json.Marshal(allClaims)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(allClaimsJSON, &claims); err != nil {
		return nil, err
	}

	// At least one of the "nbf" (Not Before), "exp" (Expiration Time), or "iat" (Issued At)
	// claims are required to be set.
	if claims.IssuedAt == nil {
		claims.IssuedAt = new(jwt.NumericDate)
	}
	if claims.Expiry == nil {
		claims.Expiry = new(jwt.NumericDate)
	}
	if claims.NotBefore == nil {
		claims.NotBefore = new(jwt.NumericDate)
	}
	if *claims.IssuedAt == 0 && *claims.Expiry == 0 && *claims.NotBefore == 0 {
		return nil, errors.New("no issued at (iat), not before (nbf), or expiration time (exp) claims in token")
	}

	// If "exp" (Expiration Time) is not set, then set it to the latest of
	// either the "iat" (Issued At) or "nbf" (Not Before) claims plus leeway.
	if *claims.Expiry == 0 {
		latestStart := *claims.IssuedAt
		if *claims.NotBefore > *claims.IssuedAt {
			latestStart = *claims.NotBefore
		}
		leeway := expected.ExpirationLeeway.Seconds()
		if expected.ExpirationLeeway.Seconds() < 0 {
			leeway = 0
		} else if expected.ExpirationLeeway.Seconds() == 0 {
			leeway = DefaultLeewaySeconds
		}
		*claims.Expiry = jwt.NumericDate(int64(latestStart) + int64(leeway))
	}

	// If "nbf" (Not Before) is not set, then set it to the "iat" (Issued At) if set.
	// Otherwise, set it to the "exp" (Expiration Time) minus leeway.
	if *claims.NotBefore == 0 {
		if *claims.IssuedAt != 0 {
			*claims.NotBefore = *claims.IssuedAt
		} else {
			leeway := expected.NotBeforeLeeway.Seconds()
			if expected.NotBeforeLeeway.Seconds() < 0 {
				leeway = 0
			} else if expected.NotBeforeLeeway.Seconds() == 0 {
				leeway = DefaultLeewaySeconds
			}
			*claims.NotBefore = jwt.NumericDate(int64(*claims.Expiry) - int64(leeway))
		}
	}

	// Set clock skew leeway to apply when validating all time-related claims
	cksLeeway := expected.ClockSkewLeeway
	if expected.ClockSkewLeeway.Seconds() < 0 {
		cksLeeway = 0
	} else if expected.ClockSkewLeeway.Seconds() == 0 {
		cksLeeway = jwt.DefaultLeeway
	}

	// Validate claims by asserting they're as expected
	if expected.Issuer != "" && expected.Issuer != claims.Issuer {
		return nil, fmt.Errorf("invalid issuer (iss) claim")
	}
	if expected.Subject != "" && expected.Subject != claims.Subject {
		return nil, fmt.Errorf("invalid subject (sub) claim")
	}
	if expected.ID != "" && expected.ID != claims.ID {
		return nil, fmt.Errorf("invalid ID (jti) claim")
	}
	if err := validateAudience(expected.Audiences, claims.Audience); err != nil {
		return nil, fmt.Errorf("invalid audience (aud) claim: %w", err)
	}

	// Validate that the token is not expired with respect to the current time
	now := time.Now()
	if expected.Now != nil {
		now = expected.Now()
	}
	if claims.NotBefore != nil && now.Add(cksLeeway).Before(claims.NotBefore.Time()) {
		return nil, errors.New("invalid not before (nbf) claim: token not yet valid")
	}
	if claims.Expiry != nil && now.Add(-cksLeeway).After(claims.Expiry.Time()) {
		return nil, errors.New("invalid expiration time (exp) claim: token is expired")
	}
	if claims.IssuedAt != nil && now.Add(cksLeeway).Before(claims.IssuedAt.Time()) {
		return nil, errors.New("invalid issued at (iat) claim: token issued in the future")
	}

	return allClaims, nil
}

// validateSigningAlgorithm checks whether the JWS "alg" (Algorithm) header
// parameter value for the given JWT matches any given in expectedAlgorithms.
// If expectedAlgorithms is empty, RS256 will be expected by default.
func validateSigningAlgorithm(token string, expectedAlgorithms []Alg) error {
	if err := SupportedSigningAlgorithm(expectedAlgorithms...); err != nil {
		return err
	}

	jws, err := jose.ParseSigned(token)
	if err != nil {
		return err
	}

	if len(jws.Signatures) == 0 {
		return fmt.Errorf("token must be signed")
	}
	if len(jws.Signatures) == 1 && len(jws.Signatures[0].Signature) == 0 {
		return fmt.Errorf("token must be signed")
	}
	if len(jws.Signatures) > 1 {
		return fmt.Errorf("token with multiple signatures not supported")
	}

	if len(expectedAlgorithms) == 0 {
		expectedAlgorithms = []Alg{RS256}
	}

	actual := Alg(jws.Signatures[0].Header.Algorithm)
	for _, expected := range expectedAlgorithms {
		if expected == actual {
			return nil
		}
	}

	return fmt.Errorf("token signed with unexpected algorithm")
}

// validateAudience returns an error if audClaim does not contain any audiences
// given by expectedAudiences. If expectedAudiences is empty, it skips validation
// and returns nil.
func validateAudience(expectedAudiences, audClaim []string) error {
	if len(expectedAudiences) == 0 {
		return nil
	}

	for _, v := range expectedAudiences {
		if contains(audClaim, v) {
			return nil
		}
	}

	return errors.New("audience claim does not match any expected audience")
}

func contains(sl []string, st string) bool {
	for _, s := range sl {
		if s == st {
			return true
		}
	}
	return false
}
