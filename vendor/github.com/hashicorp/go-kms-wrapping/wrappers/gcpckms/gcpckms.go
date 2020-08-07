package gcpckms

import (
	"errors"
	"fmt"
	"os"
	"sync/atomic"

	cloudkms "cloud.google.com/go/kms/apiv1"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	context "golang.org/x/net/context"
	"google.golang.org/api/option"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

const (
	// General GCP values, follows TF naming conventions
	EnvGCPCKMSWrapperCredsPath = "GOOGLE_CREDENTIALS"
	EnvGCPCKMSWrapperProject   = "GOOGLE_PROJECT"
	EnvGCPCKMSWrapperLocation  = "GOOGLE_REGION"

	// CKMS-specific values
	EnvGCPCKMSWrapperKeyRing     = "GCPCKMS_WRAPPER_KEY_RING"
	EnvVaultGCPCKMSSealKeyRing   = "VAULT_GCPCKMS_SEAL_KEY_RING"
	EnvGCPCKMSWrapperCryptoKey   = "GCPCKMS_WRAPPER_CRYPTO_KEY"
	EnvVaultGCPCKMSSealCryptoKey = "VAULT_GCPCKMS_SEAL_CRYPTO_KEY"
)

const (
	// GCPKMSEncrypt is used to directly encrypt the data with KMS
	GCPKMSEncrypt = iota
	// GCPKMSEnvelopeAESGCMEncrypt is when a data encryption key is generatated and
	// the data is encrypted with AESGCM and the key is encrypted with KMS
	GCPKMSEnvelopeAESGCMEncrypt
)

type Wrapper struct {
	// Values specific to IAM
	credsPath string // Path to the creds file generated during service account creation

	// Values specific to Cloud KMS service
	project    string
	location   string
	keyRing    string
	cryptoKey  string
	parentName string // Parent path built from the above values

	userAgent string

	currentKeyID *atomic.Value

	client *cloudkms.KeyManagementClient
}

var _ wrapping.Wrapper = (*Wrapper)(nil)

func NewWrapper(opts *wrapping.WrapperOptions) *Wrapper {
	if opts == nil {
		opts = new(wrapping.WrapperOptions)
	}
	s := &Wrapper{
		currentKeyID: new(atomic.Value),
	}
	s.currentKeyID.Store("")
	return s
}

// SetConfig sets the fields on the Wrapper object based on values from the
// config parameter. Environment variables take precedence over values provided
// in the config struct.
//
// Order of precedence for GCP credentials file:
// * GOOGLE_CREDENTIALS environment variable
// * `credentials` value from Value configuration file
// * GOOGLE_APPLICATION_CREDENTIALS (https://developers.google.com/identity/protocols/application-default-credentials)
func (s *Wrapper) SetConfig(config map[string]string) (map[string]string, error) {
	if config == nil {
		config = map[string]string{}
	}

	s.userAgent = config["user_agent"]

	// Do not return an error in this case. Let client initialization in
	// getClient() attempt to sort out where to get default credentials internally
	// within the SDK (e.g. checking for GOOGLE_APPLICATION_CREDENTIALS), and let
	// it error out there if none is found. This is here to establish precedence on
	// non-default input methods.
	switch {
	case os.Getenv(EnvGCPCKMSWrapperCredsPath) != "":
		s.credsPath = os.Getenv(EnvGCPCKMSWrapperCredsPath)
	case config["credentials"] != "":
		s.credsPath = config["credentials"]
	}

	switch {
	case os.Getenv(EnvGCPCKMSWrapperProject) != "":
		s.project = os.Getenv(EnvGCPCKMSWrapperProject)
	case config["project"] != "":
		s.project = config["project"]
	default:
		return nil, errors.New("'project' not found for GCP CKMS wrapper configuration")
	}

	switch {
	case os.Getenv(EnvGCPCKMSWrapperLocation) != "":
		s.location = os.Getenv(EnvGCPCKMSWrapperLocation)
	case config["region"] != "":
		s.location = config["region"]
	default:
		return nil, errors.New("'region' not found for GCP CKMS wrapper configuration")
	}

	switch {
	case os.Getenv(EnvGCPCKMSWrapperKeyRing) != "":
		s.keyRing = os.Getenv(EnvGCPCKMSWrapperKeyRing)
	case os.Getenv(EnvVaultGCPCKMSSealKeyRing) != "":
		s.keyRing = os.Getenv(EnvVaultGCPCKMSSealKeyRing)
	case config["key_ring"] != "":
		s.keyRing = config["key_ring"]
	default:
		return nil, errors.New("'key_ring' not found for GCP CKMS wrapper configuration")
	}

	switch {
	case os.Getenv(EnvGCPCKMSWrapperCryptoKey) != "":
		s.cryptoKey = os.Getenv(EnvGCPCKMSWrapperCryptoKey)
	case os.Getenv(EnvVaultGCPCKMSSealCryptoKey) != "":
		s.cryptoKey = os.Getenv(EnvVaultGCPCKMSSealCryptoKey)
	case config["crypto_key"] != "":
		s.cryptoKey = config["crypto_key"]
	default:
		return nil, errors.New("'crypto_key' not found for GCP CKMS wrapper configuration")
	}

	// Set the parent name for encrypt/decrypt requests
	s.parentName = fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", s.project, s.location, s.keyRing, s.cryptoKey)

	// Set and check s.client
	if s.client == nil {
		kmsClient, err := s.getClient()
		if err != nil {
			return nil, fmt.Errorf("error initializing GCP CKMS wrapper client: %w", err)
		}
		s.client = kmsClient

		// Make sure user has permissions to encrypt (also checks if key exists)
		ctx := context.Background()
		if _, err := s.Encrypt(ctx, []byte("vault-gcpckms-test"), nil); err != nil {
			return nil, fmt.Errorf("failed to encrypt with GCP CKMS - ensure the "+
				"key exists and the service account has at least "+
				"roles/cloudkms.cryptoKeyEncrypterDecrypter permission: %w", err)
		}
	}

	// Map that holds non-sensitive configuration info to return
	wrapperInfo := make(map[string]string)
	wrapperInfo["project"] = s.project
	wrapperInfo["region"] = s.location
	wrapperInfo["key_ring"] = s.keyRing
	wrapperInfo["crypto_key"] = s.cryptoKey

	return wrapperInfo, nil
}

// Init is called during core.Initialize. No-op at the moment
func (s *Wrapper) Init(_ context.Context) error {
	return nil
}

// Finalize is called during shutdown. This is a no-op since
// Wrapper doesn't require any cleanup.
func (s *Wrapper) Finalize(_ context.Context) error {
	return nil
}

// Type returns the type for this particular wrapper implementation
func (s *Wrapper) Type() string {
	return wrapping.GCPCKMS
}

// KeyID returns the last known key id
func (s *Wrapper) KeyID() string {
	return s.currentKeyID.Load().(string)
}

// HMACKeyID returns the last known key id
func (s *Wrapper) HMACKeyID() string {
	return ""
}

// Encrypt is used to encrypt the master key using the the AWS CMK.
// This returns the ciphertext, and/or any errors from this
// call. This should be called after s.client has been instantiated.
func (s *Wrapper) Encrypt(ctx context.Context, plaintext, aad []byte) (blob *wrapping.EncryptedBlobInfo, err error) {
	if plaintext == nil {
		return nil, errors.New("given plaintext for encryption is nil")
	}

	env, err := wrapping.NewEnvelope(nil).Encrypt(plaintext, aad)
	if err != nil {
		return nil, fmt.Errorf("error wrapping data: %w", err)
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

	ret := &wrapping.EncryptedBlobInfo{
		Ciphertext: env.Ciphertext,
		IV:         env.IV,
		KeyInfo: &wrapping.KeyInfo{
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
func (s *Wrapper) Decrypt(ctx context.Context, in *wrapping.EncryptedBlobInfo, aad []byte) (pt []byte, err error) {
	if in.Ciphertext == nil {
		return nil, fmt.Errorf("given ciphertext for decryption is nil")
	}

	// Default to mechanism used before key info was stored
	if in.KeyInfo == nil {
		in.KeyInfo = &wrapping.KeyInfo{
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
			return nil, fmt.Errorf("failed to decrypt data: %w", err)
		}

		plaintext = resp.Plaintext

	case GCPKMSEnvelopeAESGCMEncrypt:
		resp, err := s.client.Decrypt(ctx, &kmspb.DecryptRequest{
			Name:       s.parentName,
			Ciphertext: in.KeyInfo.WrappedKey,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt envelope: %w", err)
		}

		envInfo := &wrapping.EnvelopeInfo{
			Key:        resp.Plaintext,
			IV:         in.IV,
			Ciphertext: in.Ciphertext,
		}
		plaintext, err = wrapping.NewEnvelope(nil).Decrypt(envInfo, aad)
		if err != nil {
			return nil, fmt.Errorf("error decrypting data with envelope: %w", err)
		}

	default:
		return nil, fmt.Errorf("invalid mechanism: %d", in.KeyInfo.Mechanism)
	}

	return plaintext, nil
}

func (s *Wrapper) getClient() (*cloudkms.KeyManagementClient, error) {
	client, err := cloudkms.NewKeyManagementClient(context.Background(),
		option.WithCredentialsFile(s.credsPath),
		option.WithUserAgent(s.userAgent),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create KMS client: %w", err)
	}

	return client, nil
}
