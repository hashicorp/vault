package jwt

// Claims is the interface used to hold the claims values of a token
// For a type to be a Claims object, it must have a Valid method that determines
// if the token is invalid for any supported reason
// Claims are parsed and encoded using the standard library's encoding/json
// package. Claims are passed directly to that.
type Claims interface {
	// A nil validation helper should use the default helper
	Valid(*ValidationHelper) error
}

// StandardClaims is a structured version of Claims Section, as referenced at
// https://tools.ietf.org/html/rfc7519#section-4.1
// See examples for how to use this with your own claim types
type StandardClaims struct {
	Audience  ClaimStrings `json:"aud,omitempty"`
	ExpiresAt *Time        `json:"exp,omitempty"`
	ID        string       `json:"jti,omitempty"`
	IssuedAt  *Time        `json:"iat,omitempty"`
	Issuer    string       `json:"iss,omitempty"`
	NotBefore *Time        `json:"nbf,omitempty"`
	Subject   string       `json:"sub,omitempty"`
}

// Valid validates standard claims using ValidationHelper
// Validates time based claims "exp, nbf" (see: WithLeeway)
// Validates "aud" if present in claims. (see: WithAudience, WithoutAudienceValidation)
// Validates "iss" if option is provided (see: WithIssuer)
func (c StandardClaims) Valid(h *ValidationHelper) error {
	var vErr error

	if h == nil {
		h = DefaultValidationHelper
	}

	if err := h.ValidateExpiresAt(c.ExpiresAt); err != nil {
		vErr = wrapError(err, vErr)
	}

	if err := h.ValidateNotBefore(c.NotBefore); err != nil {
		vErr = wrapError(err, vErr)
	}

	if err := h.ValidateAudience(c.Audience); err != nil {
		vErr = wrapError(err, vErr)
	}

	if err := h.ValidateIssuer(c.Issuer); err != nil {
		vErr = wrapError(err, vErr)
	}

	return vErr
}

// VerifyAudience compares the aud claim against cmp.
func (c *StandardClaims) VerifyAudience(h *ValidationHelper, cmp string) error {
	return h.ValidateAudienceAgainst(c.Audience, cmp)
}

// VerifyIssuer compares the iss claim against cmp.
func (c *StandardClaims) VerifyIssuer(h *ValidationHelper, cmp string) error {
	return h.ValidateIssuerAgainst(c.Issuer, cmp)
}
