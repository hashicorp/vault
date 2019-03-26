package gcpckms

import (
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/useragent"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault/seal"
	context "golang.org/x/net/context"
	"google.golang.org/api/option"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

const (
	// General GCP values, follows TF naming conventions
	EnvGCPCKMSSealCredsPath = "GOOGLE_CREDENTIALS"
	EnvGCPCKMSSealProject   = "GOOGLE_PROJECT"
	EnvGCPCKMSSealLocation  = "GOOGLE_REGION"

	// CKMS-specific values
	EnvGCPCKMSSealKeyRing   = "VAULT_GCPCKMS_SEAL_KEY_RING"
	EnvGCPCKMSSealCryptoKey = "VAULT_GCPCKMS_SEAL_CRYPTO_KEY"
)

// GCPKMSMechanism is the method used to encrypt/decrypt in the autoseal
type GCPKMSMechanism uint32

const (
	// GCPKMSEncrypt is used to directly encrypt the data with KMS
	GCPKMSEncrypt = iota
	// GCPKMSEnvelopeAESGCMEncrypt is when a data encryption key is generatated and
	// the data is encrypted with AESGCM and the key is encrypted with KMS
	GCPKMSEnvelopeAESGCMEncrypt
)

type GCPCKMSSeal struct {
	// Values specific to IAM
	credsPath string // Path to the creds file generated during service account creation

	// Values specific to Cloud KMS service
	project    string
	location   string
	keyRing    string
	cryptoKey  string
	parentName string // Parent path built from the above values

	currentKeyID *atomic.Value

	client *cloudkms.KeyManagementClient
	logger log.Logger
}

var _ seal.Access = (*GCPCKMSSeal)(nil)

func NewSeal(logger log.Logger) *GCPCKMSSeal {
	s := &GCPCKMSSeal{
		logger:       logger,
		currentKeyID: new(atomic.Value),
	}
	s.currentKeyID.Store("")
	return s
}

// SetConfig sets the fields on the GCPCKMSSeal object based on values from the
// config parameter. Environment variables take precedence over values provided
// in the Vault configuration file (i.e. values in the `seal "gcpckms"` stanza).
//
// Order of precedence for GCP credentials file:
// * GOOGLE_CREDENTIALS environment variable
// * `credentials` value from Value configuration file
// * GOOGLE_APPLICATION_CREDENTIALS (https://developers.google.com/identity/protocols/application-default-credentials)
func (s *GCPCKMSSeal) SetConfig(config map[string]string) (map[string]string, error) {
	if config == nil {
		config = map[string]string{}
	}

	// Do not return an error in this case. Let client initialization in
	// getClient() attempt to sort out where to get default credentials internally
	// within the SDK (e.g. checking for GOOGLE_APPLICATION_CREDENTIALS), and let
	// it error out there if none is found. This is here to establish precedence on
	// non-default input methods.
	switch {
	case os.Getenv(EnvGCPCKMSSealCredsPath) != "":
		s.credsPath = os.Getenv(EnvGCPCKMSSealCredsPath)
	case config["credentials"] != "":
		s.credsPath = config["credentials"]
	}

	switch {
	case os.Getenv(EnvGCPCKMSSealProject) != "":
		s.project = os.Getenv(EnvGCPCKMSSealProject)
	case config["project"] != "":
		s.project = config["project"]
	default:
		return nil, errors.New("'project' not found for GCP CKMS seal configuration")
	}

	switch {
	case os.Getenv(EnvGCPCKMSSealLocation) != "":
		s.location = os.Getenv(EnvGCPCKMSSealLocation)
	case config["region"] != "":
		s.location = config["region"]
	default:
		return nil, errors.New("'region' not found for GCP CKMS seal configuration")
	}

	switch {
	case os.Getenv(EnvGCPCKMSSealKeyRing) != "":
		s.keyRing = os.Getenv(EnvGCPCKMSSealKeyRing)
	case config["key_ring"] != "":
		s.keyRing = config["key_ring"]
	default:
		return nil, errors.New("'key_ring' not found for GCP CKMS seal configuration")
	}

	switch {
	case os.Getenv(EnvGCPCKMSSealCryptoKey) != "":
		s.cryptoKey = os.Getenv(EnvGCPCKMSSealCryptoKey)
	case config["crypto_key"] != "":
		s.cryptoKey = config["crypto_key"]
	default:
		return nil, errors.New("'crypto_key' not found for GCP CKMS seal configuration")
	}

	// Set the parent name for encrypt/decrypt requests
	s.parentName = fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", s.project, s.location, s.keyRing, s.cryptoKey)

	// Set and check s.client
	if s.client == nil {
		kmsClient, err := s.getClient()
		if err != nil {
			return nil, errwrap.Wrapf("error initializing GCP CKMS seal client: {{err}}", err)
		}
		s.client = kmsClient

		// Make sure user has permissions to encrypt (also checks if key exists)
		ctx := context.Background()
		if _, err := s.Encrypt(ctx, []byte("vault-gcpckms-test")); err != nil {
			return nil, errwrap.Wrapf("failed to encrypt with GCP CKMS - ensure the "+
				"key exists and the service account has at least "+
				"roles/cloudkms.cryptoKeyEncrypterDecrypter permission: {{err}}", err)
		}
	}

	// Map that holds non-sensitive configuration info to return
	sealInfo := make(map[string]string)
	sealInfo["project"] = s.project
	sealInfo["region"] = s.location
	sealInfo["key_ring"] = s.keyRing
	sealInfo["crypto_key"] = s.cryptoKey

	return sealInfo, nil
}

// Init is called during core.Initialize. No-op at the moment.
func (s *GCPCKMSSeal) Init(_ context.Context) error {
	return nil
}

// Finalize is called during shutdown. This is a no-op since
// GCPKMSSeal doesn't require any cleanup.
func (s *GCPCKMSSeal) Finalize(_ context.Context) error {
	return nil
}

// SealType returns the seal type for this particular seal implementation.
func (s *GCPCKMSSeal) SealType() string {
	return seal.GCPCKMS
}

// KeyID returns the last known key id.
func (s *GCPCKMSSeal) KeyID() string {
	return s.currentKeyID.Load().(string)
}

// Encrypt is used to encrypt the master key using the the AWS CMK.
// This returns the ciphertext, and/or any errors from this
// call. This should be called after s.client has been instantiated.
func (s *GCPCKMSSeal) Encrypt(ctx context.Context, plaintext []byte) (blob *physical.EncryptedBlobInfo, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "encrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", "gcpckms", "encrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "encrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", "gcpckms", "encrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "encrypt"}, 1)
	metrics.IncrCounter([]string{"seal", "gcpckms", "encrypt"}, 1)

	if plaintext == nil {
		return nil, errors.New("given plaintext for encryption is nil")
	}

	env, err := seal.NewEnvelope().Encrypt(plaintext)
	if err != nil {
		return nil, errwrap.Wrapf("error wrapping data: {{err}}", err)
	}

	resp, err := s.client.Encrypt(ctx, &kmspb.EncryptRequest{
		Name:      s.parentName,
		Plaintext: env.Key,
	})
	if err != nil {
		return nil, err
	}

	// Store current key id value
	s.currentKeyID.Store(resp.Name)

	ret := &physical.EncryptedBlobInfo{
		Ciphertext: env.Ciphertext,
		IV:         env.IV,
		KeyInfo: &physical.SealKeyInfo{
			Mechanism: GCPKMSEnvelopeAESGCMEncrypt,
			// Even though we do not use the key id during decryption, store it
			// to know exactly what version was used in encryption in case we
			// want to rewrap older entries
			KeyID:      resp.Name,
			WrappedKey: resp.Ciphertext,
		},
	}

	return ret, nil
}

// Decrypt is used to decrypt the ciphertext.
func (s *GCPCKMSSeal) Decrypt(ctx context.Context, in *physical.EncryptedBlobInfo) (pt []byte, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "decrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", "gcpckms", "decrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "decrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", "gcpckms", "decrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "decrypt"}, 1)
	metrics.IncrCounter([]string{"seal", "gcpckms", "decrypt"}, 1)

	if in.Ciphertext == nil {
		return nil, fmt.Errorf("given ciphertext for decryption is nil")
	}

	// Default to mechanism used before key info was stored
	if in.KeyInfo == nil {
		in.KeyInfo = &physical.SealKeyInfo{
			Mechanism: GCPKMSEncrypt,
		}
	}

	var plaintext []byte
	switch in.KeyInfo.Mechanism {
	case GCPKMSEncrypt:
		resp, err := s.client.Decrypt(ctx, &kmspb.DecryptRequest{
			Name:       s.parentName,
			Ciphertext: in.Ciphertext,
		})
		if err != nil {
			return nil, errwrap.Wrapf("failed to decrypt data: {{err}}", err)
		}

		plaintext = resp.Plaintext

	case GCPKMSEnvelopeAESGCMEncrypt:
		resp, err := s.client.Decrypt(ctx, &kmspb.DecryptRequest{
			Name:       s.parentName,
			Ciphertext: in.KeyInfo.WrappedKey,
		})
		if err != nil {
			return nil, errwrap.Wrapf("failed to decrypt envelope: {{err}}", err)
		}

		envInfo := &seal.EnvelopeInfo{
			Key:        resp.Plaintext,
			IV:         in.IV,
			Ciphertext: in.Ciphertext,
		}
		plaintext, err = seal.NewEnvelope().Decrypt(envInfo)
		if err != nil {
			return nil, errwrap.Wrapf("error decrypting data with envelope: {{err}}", err)
		}

	default:
		return nil, fmt.Errorf("invalid mechanism: %d", in.KeyInfo.Mechanism)
	}

	return plaintext, nil
}

func (s *GCPCKMSSeal) getClient() (*cloudkms.KeyManagementClient, error) {
	client, err := cloudkms.NewKeyManagementClient(context.Background(),
		option.WithCredentialsFile(s.credsPath),
		option.WithUserAgent(useragent.String()),
	)
	if err != nil {
		return nil, errwrap.Wrapf("failed to create KMS client: {{err}}", err)
	}

	return client, nil
}
