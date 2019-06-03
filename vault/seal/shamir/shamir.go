package shamir

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
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

// SetConfig sets the fields on the ShamirSeal object based on
// values from the config parameter.
func (s *ShamirSeal) SetConfig(config map[string]string) (map[string]string, error) {
	// Map that holds non-sensitive configuration info
	sealInfo := make(map[string]string)

	if config == nil || config["key"] == "" {
		return sealInfo, nil
	}

	keyB64 := config["key"]
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return sealInfo, err
	}

	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return sealInfo, err
	}

	aead, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return sealInfo, err
	}

	s.aead = aead

	return sealInfo, nil
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
