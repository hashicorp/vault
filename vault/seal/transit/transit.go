package transit

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/seal"
)

// Seal is a seal that leverages Vault's Transit secret
// engine
type Seal struct {
	logger       log.Logger
	client       transitClientEncryptor
	currentKeyID *atomic.Value
}

var _ seal.Access = (*Seal)(nil)

// NewSeal creates a new transit seal
func NewSeal(logger log.Logger) *Seal {
	s := &Seal{
		logger:       logger.ResetNamed("seal-transit"),
		currentKeyID: new(atomic.Value),
	}
	s.currentKeyID.Store("")
	return s
}

// SetConfig processes the config info from the server config
func (s *Seal) SetConfig(config map[string]string) (map[string]string, error) {
	client, sealInfo, err := newTransitClient(s.logger, config)
	if err != nil {
		return nil, err
	}
	s.client = client

	// Send a value to test the seal and to set the current key id
	if _, err := s.Encrypt(context.Background(), []byte("a")); err != nil {
		client.Close()
		return nil, err
	}

	return sealInfo, nil
}

// Init is called during core.Initialize
func (s *Seal) Init(_ context.Context) error {
	return nil
}

// Finalize is called during shutdown
func (s *Seal) Finalize(_ context.Context) error {
	s.client.Close()
	return nil
}

// SealType returns the seal type for this particular seal implementation.
func (s *Seal) SealType() string {
	return seal.Transit
}

// KeyID returns the last known key id.
func (s *Seal) KeyID() string {
	return s.currentKeyID.Load().(string)
}

// Encrypt is used to encrypt using Vaults Transit engine
func (s *Seal) Encrypt(_ context.Context, plaintext []byte) (blob *physical.EncryptedBlobInfo, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "encrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", "transit", "encrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "encrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", "transit", "encrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "encrypt"}, 1)
	metrics.IncrCounter([]string{"seal", "transit", "encrypt"}, 1)

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

	ret := &physical.EncryptedBlobInfo{
		Ciphertext: ciphertext,
		KeyInfo: &physical.SealKeyInfo{
			KeyID: keyID,
		},
	}
	return ret, nil
}

// Decrypt is used to decrypt the ciphertext
func (s *Seal) Decrypt(_ context.Context, in *physical.EncryptedBlobInfo) (pt []byte, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "decrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", "transit", "decrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "decrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", "transit", "decrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "decrypt"}, 1)
	metrics.IncrCounter([]string{"seal", "transit", "decrypt"}, 1)

	plaintext, err := s.client.Decrypt(in.Ciphertext)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
