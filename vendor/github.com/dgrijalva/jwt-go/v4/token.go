package jwt

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

// TimeFunc provides the current time when parsing token to validate "exp" claim (expiration time).
// You can override it to use another time value.  This is useful for testing or if your
// server uses a different time zone than your tokens.
var TimeFunc = time.Now

// Token represents JWT Token.  Different fields will be used depending on whether you're
// creating or parsing/verifying a token.
type Token struct {
	Raw       string                 // The raw token.  Populated when you Parse a token
	Method    SigningMethod          // The signing method used or to be used
	Header    map[string]interface{} // The first segment of the token
	Claims    Claims                 // The second segment of the token
	Signature string                 // The third segment of the token.  Populated when you Parse a token
	Valid     bool                   // Is the token valid?  Populated when you Parse/Verify a token
}

// New creates a new Token.  Takes a signing method. Uses the default claims type, MapClaims.
func New(method SigningMethod) *Token {
	return NewWithClaims(method, MapClaims{})
}

// NewWithClaims creats a new token with a specified signing method and claims type
func NewWithClaims(method SigningMethod, claims Claims) *Token {
	return &Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": method.Alg(),
		},
		Claims: claims,
		Method: method,
	}
}

// SignedString returns the complete, signed token
func (t *Token) SignedString(key interface{}, opts ...SigningOption) (string, error) {
	var sig, sstr string
	var err error
	if sstr, err = t.SigningString(opts...); err != nil {
		return "", err
	}
	if sig, err = t.Method.Sign(sstr, key); err != nil {
		return "", err
	}
	return strings.Join([]string{sstr, sig}, "."), nil
}

// SigningString generates the signing string.  This is the
// most expensive part of the whole deal.  Unless you
// need this for something special, just go straight for
// the SignedString.
func (t *Token) SigningString(opts ...SigningOption) (string, error) {
	// Process options
	var cfg = new(signingOptions)
	for _, opt := range opts {
		opt(cfg)
	}
	// Setup default marshaller
	if cfg.marshaller == nil {
		cfg.marshaller = t.defaultMarshaller
	}

	// Encode the two parts, then combine
	inputParts := []interface{}{t.Header, t.Claims}
	parts := make([]string, 2)
	for i, v := range inputParts {
		ctx := CodingContext{FieldDescriptor(i), t.Header}
		jsonValue, err := cfg.marshaller(ctx, v)
		if err != nil {
			return "", err
		}
		parts[i] = EncodeSegment(jsonValue)
	}
	return strings.Join(parts, "."), nil
}

func (t *Token) defaultMarshaller(ctx CodingContext, v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Parse then validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
// If everything is kosher, err will be nil
// Claims type will be the default, MapClaims
func Parse(tokenString string, keyFunc Keyfunc, options ...ParserOption) (*Token, error) {
	return NewParser(options...).Parse(tokenString, keyFunc)
}

// ParseWithClaims is Parse, but with a specified Claims type
func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc, options ...ParserOption) (*Token, error) {
	return NewParser(options...).ParseWithClaims(tokenString, claims, keyFunc)
}

// EncodeSegment is used internally for JWT specific base64url encoding with padding stripped
func EncodeSegment(seg []byte) string {
	return base64.RawURLEncoding.EncodeToString(seg)
}

// DecodeSegment is used internally for JWT specific base64url encoding with padding stripped
func DecodeSegment(seg string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(seg)
}
