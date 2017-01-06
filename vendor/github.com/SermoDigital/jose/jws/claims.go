package jws

import (
	"encoding/json"
	"time"

	"github.com/SermoDigital/jose"
	"github.com/SermoDigital/jose/jwt"
)

// Claims represents a set of JOSE Claims.
type Claims jwt.Claims

// Get retrieves the value corresponding with key from the Claims.
func (c Claims) Get(key string) interface{} {
	return jwt.Claims(c).Get(key)
}

// Set sets Claims[key] = val. It'll overwrite without warning.
func (c Claims) Set(key string, val interface{}) {
	jwt.Claims(c).Set(key, val)
}

// Del removes the value that corresponds with key from the Claims.
func (c Claims) Del(key string) {
	jwt.Claims(c).Del(key)
}

// Has returns true if a value for the given key exists inside the Claims.
func (c Claims) Has(key string) bool {
	return jwt.Claims(c).Has(key)
}

// MarshalJSON implements json.Marshaler for Claims.
func (c Claims) MarshalJSON() ([]byte, error) {
	return jwt.Claims(c).MarshalJSON()
}

// Base64 implements the Encoder interface.
func (c Claims) Base64() ([]byte, error) {
	return jwt.Claims(c).Base64()
}

// UnmarshalJSON implements json.Unmarshaler for Claims.
func (c *Claims) UnmarshalJSON(b []byte) error {
	if b == nil {
		return nil
	}

	b, err := jose.DecodeEscaped(b)
	if err != nil {
		return err
	}

	// Since json.Unmarshal calls UnmarshalJSON,
	// calling json.Unmarshal on *p would be infinitely recursive
	// A temp variable is needed because &map[string]interface{}(*p) is
	// invalid Go.

	tmp := map[string]interface{}(*c)
	if err = json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*c = Claims(tmp)
	return nil
}

// Issuer retrieves claim "iss" per its type in
// https://tools.ietf.org/html/rfc7519#section-4.1.1
func (c Claims) Issuer() (string, bool) {
	return jwt.Claims(c).Issuer()
}

// Subject retrieves claim "sub" per its type in
// https://tools.ietf.org/html/rfc7519#section-4.1.2
func (c Claims) Subject() (string, bool) {
	return jwt.Claims(c).Subject()
}

// Audience retrieves claim "aud" per its type in
// https://tools.ietf.org/html/rfc7519#section-4.1.3
func (c Claims) Audience() ([]string, bool) {
	return jwt.Claims(c).Audience()
}

// Expiration retrieves claim "exp" per its type in
// https://tools.ietf.org/html/rfc7519#section-4.1.4
func (c Claims) Expiration() (time.Time, bool) {
	return jwt.Claims(c).Expiration()
}

// NotBefore retrieves claim "nbf" per its type in
// https://tools.ietf.org/html/rfc7519#section-4.1.5
func (c Claims) NotBefore() (time.Time, bool) {
	return jwt.Claims(c).NotBefore()
}

// IssuedAt retrieves claim "iat" per its type in
// https://tools.ietf.org/html/rfc7519#section-4.1.6
func (c Claims) IssuedAt() (time.Time, bool) {
	return jwt.Claims(c).IssuedAt()
}

// JWTID retrieves claim "jti" per its type in
// https://tools.ietf.org/html/rfc7519#section-4.1.7
func (c Claims) JWTID() (string, bool) {
	return jwt.Claims(c).JWTID()
}

// RemoveIssuer deletes claim "iss" from c.
func (c Claims) RemoveIssuer() {
	jwt.Claims(c).RemoveIssuer()
}

// RemoveSubject deletes claim "sub" from c.
func (c Claims) RemoveSubject() {
	jwt.Claims(c).RemoveIssuer()
}

// RemoveAudience deletes claim "aud" from c.
func (c Claims) RemoveAudience() {
	jwt.Claims(c).Audience()
}

// RemoveExpiration deletes claim "exp" from c.
func (c Claims) RemoveExpiration() {
	jwt.Claims(c).RemoveExpiration()
}

// RemoveNotBefore deletes claim "nbf" from c.
func (c Claims) RemoveNotBefore() {
	jwt.Claims(c).NotBefore()
}

// RemoveIssuedAt deletes claim "iat" from c.
func (c Claims) RemoveIssuedAt() {
	jwt.Claims(c).IssuedAt()
}

// RemoveJWTID deletes claim "jti" from c.
func (c Claims) RemoveJWTID() {
	jwt.Claims(c).RemoveJWTID()
}

// SetIssuer sets claim "iss" per its type in
// https://tools.ietf.org/html/rfc7519#section-4.1.1
func (c Claims) SetIssuer(issuer string) {
	jwt.Claims(c).SetIssuer(issuer)
}

// SetSubject sets claim "iss" per its type in
// https://tools.ietf.org/html/rfc7519#section-4.1.2
func (c Claims) SetSubject(subject string) {
	jwt.Claims(c).SetSubject(subject)
}

// SetAudience sets claim "aud" per its type in
// https://tools.ietf.org/html/rfc7519#section-4.1.3
func (c Claims) SetAudience(audience ...string) {
	jwt.Claims(c).SetAudience(audience...)
}

// SetExpiration sets claim "exp" per its type in
// https://tools.ietf.org/html/rfc7519#section-4.1.4
func (c Claims) SetExpiration(expiration time.Time) {
	jwt.Claims(c).SetExpiration(expiration)
}

// SetNotBefore sets claim "nbf" per its type in
// https://tools.ietf.org/html/rfc7519#section-4.1.5
func (c Claims) SetNotBefore(notBefore time.Time) {
	jwt.Claims(c).SetNotBefore(notBefore)
}

// SetIssuedAt sets claim "iat" per its type in
// https://tools.ietf.org/html/rfc7519#section-4.1.6
func (c Claims) SetIssuedAt(issuedAt time.Time) {
	jwt.Claims(c).SetIssuedAt(issuedAt)
}

// SetJWTID sets claim "jti" per its type in
// https://tools.ietf.org/html/rfc7519#section-4.1.7
func (c Claims) SetJWTID(uniqueID string) {
	jwt.Claims(c).SetJWTID(uniqueID)
}

var (
	_ json.Marshaler   = (Claims)(nil)
	_ json.Unmarshaler = (*Claims)(nil)
)
