package seal

import (
	"context"

	"github.com/hashicorp/vault/helper/xor"
	"github.com/hashicorp/vault/physical"
)

type TestSeal struct {
	Type   string
	secret []byte
}

var _ Access = (*TestSeal)(nil)

func NewTestSeal(secret []byte) *TestSeal {
	return &TestSeal{
		Type:   Test,
		secret: secret,
	}
}

func (s *TestSeal) Init(_ context.Context) error {
	return nil
}

func (t *TestSeal) Finalize(_ context.Context) error {
	return nil
}

func (t *TestSeal) SealType() string {
	return t.Type
}

func (t *TestSeal) KeyID() string {
	return "static-key"
}

func (t *TestSeal) Encrypt(_ context.Context, plaintext []byte) (*physical.EncryptedBlobInfo, error) {
	ct, err := t.obscureBytes(plaintext)
	if err != nil {
		return nil, err
	}

	return &physical.EncryptedBlobInfo{
		Ciphertext: ct,
		KeyInfo: &physical.SealKeyInfo{
			KeyID: t.KeyID(),
		},
	}, nil
}

func (t *TestSeal) Decrypt(_ context.Context, dwi *physical.EncryptedBlobInfo) ([]byte, error) {
	return t.obscureBytes(dwi.Ciphertext)
}

// obscureBytes is a helper to simulate "encryption/decryption"
// on protected values.
func (t *TestSeal) obscureBytes(in []byte) ([]byte, error) {
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
