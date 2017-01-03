package jwt

import (
	"time"

	"github.com/SermoDigital/jose/crypto"
)

// JWT represents a JWT per RFC 7519.
// It's described as an interface instead of a physical structure
// because both JWS and JWEs can be JWTs. So, in order to use either,
// import one of those two packages and use their "NewJWT" (and other)
// functions.
type JWT interface {
	// Claims returns the set of Claims.
	Claims() Claims

	// Validate returns an error describing any issues found while
	// validating the JWT. For info on the fn parameter, see the
	// comment on ValidateFunc.
	Validate(key interface{}, method crypto.SigningMethod, v ...*Validator) error

	// Serialize serializes the JWT into its on-the-wire
	// representation.
	Serialize(key interface{}) ([]byte, error)
}

// ValidateFunc is a function that provides access to the JWT
// and allows for custom validation. Keep in mind that the Verify
// methods in the JWS/JWE sibling packages call ValidateFunc *after*
// validating the JWS/JWE, but *before* any validation per the JWT
// RFC. Therefore, the ValidateFunc can be used to short-circuit
// verification, but cannot be used to circumvent the RFC.
// Custom JWT implementations are free to abuse this, but it is
// not recommended.
type ValidateFunc func(Claims) error

// Validator represents some of the validation options.
type Validator struct {
	Expected Claims        // If non-nil, these are required to match.
	EXP      time.Duration // EXPLeeway
	NBF      time.Duration // NBFLeeway
	Fn       ValidateFunc  // See ValidateFunc for more information.

	_ struct{} // Require explicitly-named struct fields.
}

// Validate validates the JWT based on the expected claims in v.
// Note: it only validates the registered claims per
// https://tools.ietf.org/html/rfc7519#section-4.1
//
// Custom claims should be validated using v's Fn member.
func (v *Validator) Validate(j JWT) error {
	if iss, ok := v.Expected.Issuer(); ok &&
		j.Claims().Get("iss") != iss {
		return ErrInvalidISSClaim
	}
	if sub, ok := v.Expected.Subject(); ok &&
		j.Claims().Get("sub") != sub {
		return ErrInvalidSUBClaim
	}
	if iat, ok := v.Expected.IssuedAt(); ok {
		if t, ok := j.Claims().GetTime("iat"); !t.Equal(iat) || !ok {
			return ErrInvalidIATClaim
		}
	}
	if jti, ok := v.Expected.JWTID(); ok &&
		j.Claims().Get("jti") != jti {
		return ErrInvalidJTIClaim
	}

	if aud, ok := v.Expected.Audience(); ok {
		aud2, ok := j.Claims().Audience()
		if !ok || !ValidAudience(aud, aud2) {
			return ErrInvalidAUDClaim
		}
	}

	if v.Fn != nil {
		return v.Fn(j.Claims())
	}
	return nil
}

// SetClaim sets the claim with the given val.
func (v *Validator) SetClaim(claim string, val interface{}) {
	v.expect()
	v.Expected.Set(claim, val)
}

// SetIssuer sets the "iss" claim per
// https://tools.ietf.org/html/rfc7519#section-4.1.1
func (v *Validator) SetIssuer(iss string) {
	v.expect()
	v.Expected.Set("iss", iss)
}

// SetSubject sets the "sub" claim per
// https://tools.ietf.org/html/rfc7519#section-4.1.2
func (v *Validator) SetSubject(sub string) {
	v.expect()
	v.Expected.Set("sub", sub)
}

// SetAudience sets the "aud" claim per
// https://tools.ietf.org/html/rfc7519#section-4.1.3
func (v *Validator) SetAudience(aud string) {
	v.expect()
	v.Expected.Set("aud", aud)
}

// SetExpiration sets the "exp" claim per
// https://tools.ietf.org/html/rfc7519#section-4.1.4
func (v *Validator) SetExpiration(exp time.Time) {
	v.expect()
	v.Expected.Set("exp", exp)
}

// SetNotBefore sets the "nbf" claim per
// https://tools.ietf.org/html/rfc7519#section-4.1.5
func (v *Validator) SetNotBefore(nbf time.Time) {
	v.expect()
	v.Expected.Set("nbf", nbf)
}

// SetIssuedAt sets the "iat" claim per
// https://tools.ietf.org/html/rfc7519#section-4.1.6
func (v *Validator) SetIssuedAt(iat time.Time) {
	v.expect()
	v.Expected.Set("iat", iat)
}

// SetJWTID sets the "jti" claim per
// https://tools.ietf.org/html/rfc7519#section-4.1.7
func (v *Validator) SetJWTID(jti string) {
	v.expect()
	v.Expected.Set("jti", jti)
}

func (v *Validator) expect() {
	if v.Expected == nil {
		v.Expected = make(Claims)
	}
}
