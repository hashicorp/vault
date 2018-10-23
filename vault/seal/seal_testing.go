package seal

import (
	"context"

	"github.com/hashicorp/vault/physical"
)

type TestSeal struct{}

var _ Access = (*TestSeal)(nil)

func (s *TestSeal) Init(_ context.Context) error {
	return nil
}

func (t *TestSeal) Finalize(_ context.Context) error {
	return nil
}

func (t *TestSeal) SealType() string {
	return Test
}

func (t *TestSeal) KeyID() string {
	return "static-key"
}

func (t *TestSeal) Encrypt(_ context.Context, plaintext []byte) (*physical.EncryptedBlobInfo, error) {
	return &physical.EncryptedBlobInfo{
		Ciphertext: ReverseBytes(plaintext),
	}, nil
}

func (t *TestSeal) Decrypt(_ context.Context, dwi *physical.EncryptedBlobInfo) ([]byte, error) {
	return ReverseBytes(dwi.Ciphertext), nil
}

// reverseBytes is a helper to simulate "encryption/decryption"
// on protected values.
func ReverseBytes(in []byte) []byte {
	out := make([]byte, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = in[len(in)-1-i]
	}
	return out
}
