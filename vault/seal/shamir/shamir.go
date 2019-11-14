package shamir

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/seal"
)

// ShamirSeal implements the seal.Access interface for Shamir unseal
type ShamirSeal struct {
	logger log.Logger
	key    []byte
	aead   cipher.AEAD
}

// Ensure that we are implementing AutoSealAccess
var _ seal.Access = (*ShamirSeal)(nil)

// NewSeal creates a new ShamirSeal with the provided logger
func NewSeal(logger log.Logger) *ShamirSeal {
	seal := &ShamirSeal{
		logger: logger,
	}
	return seal
}

func (s *ShamirSeal) GetKey() []byte {
	return s.key
}

func (s *ShamirSeal) SetKey(key []byte) error {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aead, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return err
	}

	s.key = key
	s.aead = aead
	return nil
}

// Init is called during core.Initialize. No-op at the moment.
func (s *ShamirSeal) Init(_ context.Context) error {
	return nil
}

// Finalize is called during shutdown. This is a no-op since
// ShamirSeal doesn't require any cleanup.
func (s *ShamirSeal) Finalize(_ context.Context) error {
	return nil
}

// SealType returns the seal type for this particular seal implementation.
func (s *ShamirSeal) SealType() string {
	return seal.Shamir
}

// KeyID returns the last known key id.
func (s *ShamirSeal) KeyID() string {
	return ""
}

// Encrypt is used to encrypt the plaintext using the aead held by the seal.
func (s *ShamirSeal) Encrypt(_ context.Context, plaintext []byte) (*physical.EncryptedBlobInfo, error) {
	if plaintext == nil {
		return nil, fmt.Errorf("given plaintext for encryption is nil")
	}

	if s.aead == nil {
		return nil, errors.New("aead is not configured in the seal")
	}

	iv, err := uuid.GenerateRandomBytes(12)
	if err != nil {
		return nil, err
	}

	ciphertext := s.aead.Seal(nil, iv, plaintext, nil)

	return &physical.EncryptedBlobInfo{
		Ciphertext: append(iv, ciphertext...),
	}, nil
}

func (s *ShamirSeal) Decrypt(_ context.Context, in *physical.EncryptedBlobInfo) ([]byte, error) {
	if in == nil {
		return nil, fmt.Errorf("given plaintext for encryption is nil")
	}

	if s.aead == nil {
		return nil, errors.New("aead is not configured in the seal")
	}

	iv, ciphertext := in.Ciphertext[:12], in.Ciphertext[12:]

	plaintext, err := s.aead.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
