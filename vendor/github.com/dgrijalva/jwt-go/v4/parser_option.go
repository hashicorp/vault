package jwt

import "time"

// ParserOption implements functional options for parser behavior
// see: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
type ParserOption func(*Parser)

// WithValidMethods returns the ParserOption for specifying valid signing methods
func WithValidMethods(valid []string) ParserOption {
	return func(p *Parser) {
		p.validMethods = valid
	}
}

// WithJSONNumber returns the ParserOption for using json.Number instead of float64 when parsing
// numeric values. Used most commonly with MapClaims, but it can be useful in some cases with
// structured claims types
func WithJSONNumber() ParserOption {
	return func(p *Parser) {
		p.useJSONNumber = true
	}
}

// WithoutClaimsValidation returns the ParserOption for disabling claims validation
// This does not disable signature validation. Use this if you want intend to implement
// claims validation via other means
func WithoutClaimsValidation() ParserOption {
	return func(p *Parser) {
		p.skipClaimsValidation = true
	}
}

// WithLeeway returns the ParserOption for specifying the leeway window.
func WithLeeway(d time.Duration) ParserOption {
	return func(p *Parser) {
		p.ValidationHelper.leeway = d
	}
}

// WithAudience returns the ParserOption for specifying an expected aud member value
func WithAudience(aud string) ParserOption {
	return func(p *Parser) {
		p.ValidationHelper.audience = &aud
	}
}

// WithoutAudienceValidation returns the ParserOption that specifies audience check should be skipped
func WithoutAudienceValidation() ParserOption {
	return func(p *Parser) {
		p.ValidationHelper.skipAudience = true
	}
}

// WithIssuer returns the ParserOption that specifies a value to compare against the iss claim
func WithIssuer(iss string) ParserOption {
	return func(p *Parser) {
		p.ValidationHelper.issuer = &iss
	}
}

// TokenUnmarshaller is the function signature required to supply custom JSON decoding logic.
// It is the same as json.Marshal with the addition of the FieldDescriptor.
// The field value will let your marshaller know which field is being processed.
// This is to facilitate things like compression, where you wouldn't want to compress
// the head.
type TokenUnmarshaller func(ctx CodingContext, data []byte, v interface{}) error

// WithUnmarshaller returns the ParserOption that replaces the specified decoder
func WithUnmarshaller(um TokenUnmarshaller) ParserOption {
	return func(p *Parser) {
		p.unmarshaller = um
	}
}
