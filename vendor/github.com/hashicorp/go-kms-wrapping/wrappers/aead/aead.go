package aead

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"errors"

	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/go-uuid"
)

// Wrapper implements the wrapping.Wrapper interface for Shamir
type Wrapper struct {
	keyBytes []byte
	aead     cipher.AEAD
}

// Ensure that we are implementing AutoSealAccess
var _ wrapping.Wrapper = (*Wrapper)(nil)

// NewWrapper creates a new Wrapper with the provided logger
func NewWrapper(opts *wrapping.WrapperOptions) *Wrapper {
	if opts == nil {
		opts = new(wrapping.WrapperOptions)
	}
	seal := new(Wrapper)
	return seal
}

func (s *Wrapper) GetKeyBytes() []byte {
	return s.keyBytes
}

func (s *Wrapper) SetAEAD(aead cipher.AEAD) {
	s.aead = aead
}

// SetAESGCMKeyBytes takes in a byte slice and constucts an AES-GCM AEAD from it
func (s *Wrapper) SetAESGCMKeyBytes(key []byte) error {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aead, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return err
	}

	s.keyBytes = key
	s.aead = aead
	return nil
}

// Init is a no-op at the moment
func (s *Wrapper) Init(_ context.Context) error {
	return nil
}

// Finalize is called during shutdown. This is a no-op since
// Wrapper doesn't require any cleanup.
func (s *Wrapper) Finalize(_ context.Context) error {
	return nil
}

// Type returns the seal type for this particular Wrapper implementation
func (s *Wrapper) Type() string {
	return wrapping.Shamir
}

// KeyID returns the last known key id
func (s *Wrapper) KeyID() string {
	return ""
}

// HMACKeyID returns the last known HMAC key id
func (s *Wrapper) HMACKeyID() string {
	return ""
}

// Encrypt is used to encrypt the plaintext using the aead held by the seal.
func (s *Wrapper) Encrypt(_ context.Context, plaintext, aad []byte) (*wrapping.EncryptedBlobInfo, error) {
	if plaintext == nil {
		return nil, errors.New("given plaintext for encryption is nil")
	}

	if s.aead == nil {
		return nil, errors.New("aead is not configured in the seal")
	}

	iv, err := uuid.GenerateRandomBytes(12)
	if err != nil {
		return nil, err
	}

	ciphertext := s.aead.Seal(nil, iv, plaintext, aad)

	return &wrapping.EncryptedBlobInfo{
		Ciphertext: append(iv, ciphertext...),
	}, nil
}

func (s *Wrapper) Decrypt(_ context.Context, in *wrapping.EncryptedBlobInfo, aad []byte) ([]byte, error) {
	if in == nil {
		return nil, errors.New("given plaintext for encryption is nil")
	}

	if s.aead == nil {
		return nil, errors.New("aead is not configured in the seal")
	}

	iv, ciphertext := in.Ciphertext[:12], in.Ciphertext[12:]

	plaintext, err := s.aead.Open(nil, iv, ciphertext, aad)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
