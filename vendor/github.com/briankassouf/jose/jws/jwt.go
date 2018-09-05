package jws

import (
	"net/http"
	"time"

	"github.com/briankassouf/jose"
	"github.com/briankassouf/jose/crypto"
	"github.com/briankassouf/jose/jwt"
)

// NewJWT creates a new JWT with the given claims.
func NewJWT(claims Claims, method crypto.SigningMethod) jwt.JWT {
	j, ok := New(claims, method).(*jws)
	if !ok {
		panic("jws.NewJWT: runtime panic: New(...).(*jws) != true")
	}
	j.sb[0].protected.Set("typ", "JWT")
	j.isJWT = true
	return j
}

// Serialize helps implements jwt.JWT.
func (j *jws) Serialize(key interface{}) ([]byte, error) {
	if j.isJWT {
		return j.Compact(key)
	}
	return nil, ErrIsNotJWT
}

// Claims helps implements jwt.JWT.
func (j *jws) Claims() jwt.Claims {
	if j.isJWT {
		if c, ok := j.payload.v.(Claims); ok {
			return jwt.Claims(c)
		}
	}
	return nil
}

// ParseJWTFromRequest tries to find the JWT in an http.Request.
// This method will call ParseMultipartForm if there's no token in the header.
func ParseJWTFromRequest(req *http.Request) (jwt.JWT, error) {
	if b, ok := fromHeader(req); ok {
		return ParseJWT(b)
	}
	if b, ok := fromForm(req); ok {
		return ParseJWT(b)
	}
	return nil, ErrNoTokenInRequest
}

// ParseJWT parses a serialized jwt.JWT into a physical jwt.JWT.
// If its payload isn't a set of claims (or able to be coerced into
// a set of claims) it'll return an error stating the
// JWT isn't a JWT.
func ParseJWT(encoded []byte) (jwt.JWT, error) {
	t, err := parseCompact(encoded, true)
	if err != nil {
		return nil, err
	}
	c, ok := t.Payload().(map[string]interface{})
	if !ok {
		return nil, ErrIsNotJWT
	}
	t.SetPayload(Claims(c))
	return t, nil
}

// IsJWT returns true if the JWS is a JWT.
func (j *jws) IsJWT() bool {
	return j.isJWT
}

func (j *jws) Validate(key interface{}, m crypto.SigningMethod, v ...*jwt.Validator) error {
	if j.isJWT {
		if err := j.Verify(key, m); err != nil {
			return err
		}
		var v1 jwt.Validator
		if len(v) > 0 {
			v1 = *v[0]
		}
		c, ok := j.payload.v.(Claims)
		if ok {
			if err := v1.Validate(j); err != nil {
				return err
			}
			return jwt.Claims(c).Validate(jose.Now(), v1.EXP, v1.NBF)
		}
	}
	return ErrIsNotJWT
}

// Conv converts a func(Claims) error to type jwt.ValidateFunc.
func Conv(fn func(Claims) error) jwt.ValidateFunc {
	if fn == nil {
		return nil
	}
	return func(c jwt.Claims) error {
		return fn(Claims(c))
	}
}

// NewValidator returns a jwt.Validator.
func NewValidator(c Claims, exp, nbf time.Duration, fn func(Claims) error) *jwt.Validator {
	return &jwt.Validator{
		Expected: jwt.Claims(c),
		EXP:      exp,
		NBF:      nbf,
		Fn:       Conv(fn),
	}
}

var _ jwt.JWT = (*jws)(nil)
