package wrapping

import (
	"context"

	"github.com/hashicorp/go-kms-wrapping/internal/xor"
)

// TestWrapper is a wrapper that can be used for tests
type TestWrapper struct {
	wrapperType string
	secret      []byte
	keyID       string
}

var _ Wrapper = (*TestWrapper)(nil)

// NewTestWrapper constructs a test wrapper
func NewTestWrapper(secret []byte) *TestWrapper {
	return &TestWrapper{
		wrapperType: Test,
		secret:      secret,
		keyID:       "static-key",
	}
}

// Init initializes the test wrapper
func (t *TestWrapper) Init(_ context.Context) error {
	return nil
}

// Finalize finalizes the test wrapper
func (t *TestWrapper) Finalize(_ context.Context) error {
	return nil
}

// Type returns the type of the test wrapper
func (t *TestWrapper) Type() string {
	return t.wrapperType
}

// KeyID returns the configured key ID
func (t *TestWrapper) KeyID() string {
	return t.keyID
}

// HMACKeyID returns the configured HMAC key ID
func (t *TestWrapper) HMACKeyID() string {
	return ""
}

// SetKeyID allows setting the test wrapper's key ID
func (t *TestWrapper) SetKeyID(k string) {
	t.keyID = k
}

// Encrypt allows encrypting via the test wrapper
func (t *TestWrapper) Encrypt(_ context.Context, plaintext, _ []byte) (*EncryptedBlobInfo, error) {
	ct, err := t.obscureBytes(plaintext)
	if err != nil {
		return nil, err
	}

	return &EncryptedBlobInfo{
		Ciphertext: ct,
		KeyInfo: &KeyInfo{
			KeyID: t.KeyID(),
		},
	}, nil
}

// Decrypt allows decrypting via the test wrapper
func (t *TestWrapper) Decrypt(_ context.Context, dwi *EncryptedBlobInfo, _ []byte) ([]byte, error) {
	return t.obscureBytes(dwi.Ciphertext)
}

// obscureBytes is a helper to simulate "encryption/decryption"
// on protected values.
func (t *TestWrapper) obscureBytes(in []byte) ([]byte, error) {
	out := make([]byte, len(in))

	if len(t.secret) != 0 {
		// make sure they are the same length
		localSecret := make([]byte, len(in))
		copy(localSecret, t.secret)

		var err error

		out, err = xor.XORBytes(in, localSecret)
		if err != nil {
			return nil, err
		}

	} else {
		// if there is no secret, simply reverse the string
		for i := 0; i < len(in); i++ {
			out[i] = in[len(in)-1-i]
		}
	}

	return out, nil
}
