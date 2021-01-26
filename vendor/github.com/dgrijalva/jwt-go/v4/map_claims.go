package jwt

// MapClaims is the Claims type that uses the map[string]interface{} for JSON decoding
// This is the default Claims type if you don't supply one
type MapClaims map[string]interface{}

// VerifyAudience compares the aud claim against cmp.
func (m MapClaims) VerifyAudience(h *ValidationHelper, cmp string) error {
	if aud, err := ParseClaimStrings(m["aud"]); err == nil && aud != nil {
		return h.ValidateAudienceAgainst(aud, cmp)
	} else if err != nil {
		return &MalformedTokenError{Message: "couldn't parse 'aud' value"}
	}
	return nil
}

// VerifyIssuer compares the iss claim against cmp.
func (m MapClaims) VerifyIssuer(h *ValidationHelper, cmp string) error {
	iss, ok := m["iss"].(string)
	if !ok {
		return &InvalidIssuerError{Message: "'iss' expected but not present"}
	}
	return h.ValidateIssuerAgainst(iss, cmp)
}

// Valid validates standard claims using ValidationHelper
// Validates time based claims "exp, nbf" (see: WithLeeway)
// Validates "aud" if present in claims. (see: WithAudience, WithoutAudienceValidation)
// Validates "iss" if option is provided (see: WithIssuer)
func (m MapClaims) Valid(h *ValidationHelper) error {
	var vErr error

	if h == nil {
		h = DefaultValidationHelper
	}

	exp, err := m.LoadTimeValue("exp")
	if err != nil {
		return err
	}

	if err = h.ValidateExpiresAt(exp); err != nil {
		vErr = wrapError(err, vErr)
	}

	nbf, err := m.LoadTimeValue("nbf")
	if err != nil {
		return err
	}

	if err = h.ValidateNotBefore(nbf); err != nil {
		vErr = wrapError(err, vErr)
	}

	// Try to parse the 'aud' claim
	if aud, err := ParseClaimStrings(m["aud"]); err == nil && aud != nil {
		// If it's present and well formed, validate
		if err = h.ValidateAudience(aud); err != nil {
			vErr = wrapError(err, vErr)
		}
	} else if err != nil {
		// If it's present and not well formed, return an error
		return &MalformedTokenError{Message: "couldn't parse 'aud' value"}
	}

	iss, _ := m["iss"].(string)
	if err = h.ValidateIssuer(iss); err != nil {
		vErr = wrapError(err, vErr)
	}

	return vErr
}

// LoadTimeValue extracts a *Time value from a key in m
func (m MapClaims) LoadTimeValue(key string) (*Time, error) {
	value, ok := m[key]
	if !ok {
		// No value present in map
		return nil, nil
	}

	return ParseTime(value)
}
