package jwt

import (
	"crypto/subtle"
	"fmt"
	"time"
)

// DefaultValidationHelper is used by Claims.Valid if none is provided
var DefaultValidationHelper = &ValidationHelper{}

// ValidationHelper is built by the parser and passed
// to Claims.Value to carry parse/validation options
// This standalone type exists to allow implementations to do whatever custom
// behavior is required while still being able to call upon the standard behavior
// as necessary.
type ValidationHelper struct {
	nowFunc      func() time.Time // Override for time.Now. Mostly used for testing
	leeway       time.Duration    // Leeway to provide when validating time values
	audience     *string          // Expected audience value
	skipAudience bool             // Ignore aud check
	issuer       *string          // Expected issuer value. ignored if nil
}

// NewValidationHelper creates a validation helper from a list of parser options
// Not all parser options will impact validation
// If you already have a custom parser, you can use its ValidationHelper value
// instead of creating a new one
func NewValidationHelper(options ...ParserOption) *ValidationHelper {
	p := NewParser(options...)
	return p.ValidationHelper
}

func (h *ValidationHelper) now() time.Time {
	if h.nowFunc != nil {
		return h.nowFunc()
	}
	return TimeFunc()
}

// Before returns true if Now is before t
// Takes leeway into account
func (h *ValidationHelper) Before(t time.Time) bool {
	return h.now().Before(t.Add(-h.leeway))
}

// After returns true if Now is after t
// Takes leeway into account
func (h *ValidationHelper) After(t time.Time) bool {
	return h.now().After(t.Add(h.leeway))
}

// ValidateExpiresAt returns an error if the expiration time is invalid
// Takes leeway into account
func (h *ValidationHelper) ValidateExpiresAt(exp *Time) error {
	// 'exp' claim is not set. ignore.
	if exp == nil {
		return nil
	}

	// Expiration has passed
	if h.After(exp.Time) {
		delta := h.now().Sub(exp.Time)
		return &TokenExpiredError{At: h.now(), ExpiredBy: delta}
	}

	// Expiration has not passed
	return nil
}

// ValidateNotBefore returns an error if the nbf time has not been reached
// Takes leeway into account
func (h *ValidationHelper) ValidateNotBefore(nbf *Time) error {
	// 'nbf' claim is not set. ignore.
	if nbf == nil {
		return nil
	}

	// Nbf hasn't been reached
	if h.Before(nbf.Time) {
		delta := nbf.Time.Sub(h.now())
		return &TokenNotValidYetError{At: h.now(), EarlyBy: delta}
	}
	// Nbf has been reached. valid.
	return nil
}

// ValidateAudience verifies that aud contains the audience value provided
// by the WithAudience option.
// Per the spec (https://tools.ietf.org/html/rfc7519#section-4.1.3), if the aud
// claim is present,
func (h *ValidationHelper) ValidateAudience(aud ClaimStrings) error {
	// Skip flag
	if h.skipAudience {
		return nil
	}

	// If there's no audience claim, ignore
	if aud == nil || len(aud) == 0 {
		return nil
	}

	// If there is an audience claim, but no value provided, fail
	if h.audience == nil {
		return &InvalidAudienceError{Message: "audience value was expected but not provided"}
	}

	return h.ValidateAudienceAgainst(aud, *h.audience)
}

// ValidateAudienceAgainst checks that the compare value is included in the aud list
// It is used by ValidateAudience, but exposed as a helper for other implementations
func (h *ValidationHelper) ValidateAudienceAgainst(aud ClaimStrings, compare string) error {
	if aud == nil {
		return nil
	}

	// Compare provided value with aud claim.
	// This code avoids the early return to make this check more or less constant time.
	// I'm not certain that's actually required in this context.
	var match = false
	for _, audStr := range aud {
		if subtle.ConstantTimeCompare([]byte(audStr), []byte(compare)) == 1 {
			match = true
		}
	}
	if !match {
		return &InvalidAudienceError{Message: fmt.Sprintf("'%v' wasn't found in aud claim", compare)}
	}
	return nil

}

// ValidateIssuer checks the claim value against the value provided by WithIssuer
func (h *ValidationHelper) ValidateIssuer(iss string) error {
	// Always passes validation if issuer is not provided
	if h.issuer == nil {
		return nil
	}

	return h.ValidateIssuerAgainst(iss, *h.issuer)
}

// ValidateIssuerAgainst checks the claim value against the value provided, ignoring the WithIssuer value
func (h *ValidationHelper) ValidateIssuerAgainst(iss string, compare string) error {
	if subtle.ConstantTimeCompare([]byte(iss), []byte(compare)) == 1 {
		return nil
	}

	return &InvalidIssuerError{Message: "'iss' value doesn't match expectation"}
}
