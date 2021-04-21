package oidc

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/cap/oidc/internal/base62"
)

// ChallengeMethod represents PKCE code challenge methods as defined by RFC
// 7636.
type ChallengeMethod string

const (
	// PKCE code challenge methods as defined by RFC 7636.
	//
	// See: https://tools.ietf.org/html/rfc7636#page-9
	S256 ChallengeMethod = "S256" // SHA-256
)

// CodeVerifier represents an OAuth PKCE code verifier.
//
// See: https://tools.ietf.org/html/rfc7636#section-4.1
type CodeVerifier interface {

	// Verifier returns the code verifier (see:
	// https://tools.ietf.org/html/rfc7636#section-4.1)
	Verifier() string

	// Challenge returns the code verifier's code challenge (see:
	// https://tools.ietf.org/html/rfc7636#section-4.2)
	Challenge() string

	// Method returns the code verifier's challenge method (see
	// https://tools.ietf.org/html/rfc7636#section-4.2)
	Method() ChallengeMethod

	// Copy returns a copy of the verifier
	Copy() CodeVerifier
}

// S256Verifier represents an OAuth PKCE code verifier that uses the S256
// challenge method.  It implements the CodeVerifier interface.
type S256Verifier struct {
	verifier  string
	challenge string
	method    ChallengeMethod
}

// min len of 43 chars per https://tools.ietf.org/html/rfc7636#section-4.1
const verifierLen = 43

// NewCodeVerifier creates a new CodeVerifier (*S256Verifier).
//
// See: https://tools.ietf.org/html/rfc7636#section-4.1
func NewCodeVerifier() (*S256Verifier, error) {
	const op = "NewCodeVerifier"
	data, err := base62.Random(verifierLen)
	if err != nil {
		return nil, fmt.Errorf("%s: unable to create verifier data %w", op, err)
	}
	v := &S256Verifier{
		verifier: data, // no need to encode it, since bas62.Random uses a limited set of characters.
		method:   S256,
	}
	if v.challenge, err = CreateCodeChallenge(v); err != nil {
		return nil, fmt.Errorf("%s: unable to create code challenge: %w", op, err)
	}
	return v, nil
}

func (v *S256Verifier) Verifier() string        { return v.verifier }  // Verifier implements the CodeVerifier.Verifier() interface function.
func (v *S256Verifier) Challenge() string       { return v.challenge } // Challenge implements the CodeVerifier.Challenge() interface function.
func (v *S256Verifier) Method() ChallengeMethod { return v.method }    // Method implements the CodeVerifier.Method() interface function.

// Copy returns a copy of the verifier.
func (v *S256Verifier) Copy() CodeVerifier {
	return &S256Verifier{
		verifier:  v.verifier,
		challenge: v.challenge,
		method:    v.method,
	}
}

// CreateCodeChallenge creates a code challenge from the verifier. Supported
// ChallengeMethods: S256
//
// See: https://tools.ietf.org/html/rfc7636#section-4.2
func CreateCodeChallenge(v CodeVerifier) (string, error) {
	// currently, we only support S256
	if v.Method() != S256 {
		return "", fmt.Errorf("CreateCodeChallenge: %s is invalid: %w", v.Method(), ErrUnsupportedChallengeMethod)
	}
	h := sha256.New()
	_, _ = h.Write([]byte(v.Verifier())) // hash documents that Write will never return an Error
	sum := h.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(sum), nil
}
