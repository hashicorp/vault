package shamir

import (
	"context"
	"fmt"
	"sync/atomic"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault/seal"
)

// AWSKMSSeal represents credentials and Key information for the KMS Key used to
// encryption and decryption
type ShamirSeal struct {
	currentKeyID *atomic.Value

	logger log.Logger
}

// Ensure that we are implementing AutoSealAccess
var _ seal.Access = (*ShamirSeal)(nil)

// NewSeal creates a new AWSKMS seal with the provided logger
func NewSeal(logger log.Logger) *ShamirSeal {
	k := &ShamirSeal{
		logger:       logger,
		currentKeyID: new(atomic.Value),
	}
	k.currentKeyID.Store("")
	return k
}

// SetConfig sets the fields on the AWSKMSSeal object based on
// values from the config parameter.
//
// Order of precedence AWS values:
// * Environment variable
// * Value from Vault configuration file
// * Instance metadata role (access key and secret key)
// * Default values
func (k *ShamirSeal) SetConfig(config map[string]string) (map[string]string, error) {
	// Map that holds non-sensitive configuration info
	sealInfo := make(map[string]string)

	return sealInfo, nil
}

// Init is called during core.Initialize. No-op at the moment.
func (k *ShamirSeal) Init(_ context.Context) error {
	return nil
}

// Finalize is called during shutdown. This is a no-op since
// AWSKMSSeal doesn't require any cleanup.
func (k *ShamirSeal) Finalize(_ context.Context) error {
	return nil
}

// SealType returns the seal type for this particular seal implementation.
func (k *ShamirSeal) SealType() string {
	return seal.AWSKMS
}

// KeyID returns the last known key id.
func (k *ShamirSeal) KeyID() string {
	return k.currentKeyID.Load().(string)
}

// Encrypt is used to encrypt the master key using the the AWS CMK.
// This returns the ciphertext, and/or any errors from this
// call. This should be called after the KMS client has been instantiated.
func (k *ShamirSeal) Encrypt(_ context.Context, plaintext []byte) (blob *physical.EncryptedBlobInfo, err error) {
	if plaintext == nil {
		return nil, fmt.Errorf("given plaintext for encryption is nil")
	}

	return
}

func (k *ShamirSeal) Decrypt(_ context.Context, e *physical.EncryptedBlobInfo) (pt []byte, err error) {
	if e == nil {
		return nil, fmt.Errorf("given plaintext for encryption is nil")
	}

	return
}
