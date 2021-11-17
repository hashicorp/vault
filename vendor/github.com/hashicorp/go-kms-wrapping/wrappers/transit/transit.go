package transit

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"

	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
)

// Wrapper is a wrapper that leverages Vault's Transit secret
// engine
type Wrapper struct {
	logger       hclog.Logger
	client       transitClientEncryptor
	currentKeyID *atomic.Value
}

var _ wrapping.Wrapper = (*Wrapper)(nil)

// NewWrapper creates a new transit wrapper
func NewWrapper(opts *wrapping.WrapperOptions) *Wrapper {
	if opts == nil {
		opts = new(wrapping.WrapperOptions)
	}
	s := &Wrapper{
		logger:       opts.Logger,
		currentKeyID: new(atomic.Value),
	}
	s.currentKeyID.Store("")
	return s
}

// SetConfig processes the config info from the server config
func (s *Wrapper) SetConfig(config map[string]string) (map[string]string, error) {
	client, wrapperInfo, err := newTransitClient(s.logger, config)
	if err != nil {
		return nil, err
	}
	s.client = client

	// Send a value to test the wrapper and to set the current key id
	if _, err := s.Encrypt(context.Background(), []byte("a"), nil); err != nil {
		client.Close()
		return nil, err
	}

	return wrapperInfo, nil
}

// Init is called during core.Initialize
func (s *Wrapper) Init(_ context.Context) error {
	return nil
}

// Finalize is called during shutdown
func (s *Wrapper) Finalize(_ context.Context) error {
	s.client.Close()
	return nil
}

// Type returns the type for this particular Wrapper implementation
func (s *Wrapper) Type() string {
	return wrapping.Transit
}

// KeyID returns the last known key id
func (s *Wrapper) KeyID() string {
	return s.currentKeyID.Load().(string)
}

// HMACKeyID returns the last known HMAC key id
func (s *Wrapper) HMACKeyID() string {
	return ""
}

// Encrypt is used to encrypt using Vault's Transit engine
func (s *Wrapper) Encrypt(_ context.Context, plaintext, aad []byte) (blob *wrapping.EncryptedBlobInfo, err error) {
	ciphertext, err := s.client.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}

	splitKey := strings.Split(string(ciphertext), ":")
	if len(splitKey) != 3 {
		return nil, errors.New("invalid ciphertext returned")
	}
	keyID := splitKey[1]
	s.currentKeyID.Store(keyID)

	ret := &wrapping.EncryptedBlobInfo{
		Ciphertext: ciphertext,
		KeyInfo: &wrapping.KeyInfo{
			KeyID: keyID,
		},
	}
	return ret, nil
}

// Decrypt is used to decrypt the ciphertext
func (s *Wrapper) Decrypt(_ context.Context, in *wrapping.EncryptedBlobInfo, _ []byte) (pt []byte, err error) {
	plaintext, err := s.client.Decrypt(in.Ciphertext)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// GetClient returns the transit Wrapper's transitClientEncryptor
func (s *Wrapper) GetClient() transitClientEncryptor {
	return s.client
}
