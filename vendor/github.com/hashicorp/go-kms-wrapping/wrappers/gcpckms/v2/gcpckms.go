// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpckms

import (
	"errors"
	"fmt"
	"os"
	"sync/atomic"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	context "golang.org/x/net/context"
	"google.golang.org/api/option"
)

const (
	// General GCP values, follows TF naming conventions
	EnvGcpCkmsWrapperCredsPath = "GOOGLE_CREDENTIALS"
	EnvGcpCkmsWrapperProject   = "GOOGLE_PROJECT"
	EnvGcpCkmsWrapperLocation  = "GOOGLE_REGION"

	// CKMS-specific values
	EnvGcpCkmsWrapperKeyRing     = "GCPCKMS_WRAPPER_KEY_RING"
	EnvVaultGcpCkmsSealKeyRing   = "VAULT_GCPCKMS_SEAL_KEY_RING"
	EnvGcpCkmsWrapperCryptoKey   = "GCPCKMS_WRAPPER_CRYPTO_KEY"
	EnvVaultGcpCkmsSealCryptoKey = "VAULT_GCPCKMS_SEAL_CRYPTO_KEY"
)

const (
	// GcpCkmsEncrypt is used to directly encrypt the data with KMS
	GcpCkmsEncrypt = iota
	// GcpCkmsEnvelopeAesGcmEncrypt is when a data encryption key is generatated and
	// the data is encrypted with AES-GCM and the key is encrypted with KMS
	GcpCkmsEnvelopeAesGcmEncrypt
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

	currentKeyId   *atomic.Value
	keyNotRequired bool

	client *cloudkms.KeyManagementClient
}

var keyPermissions = []string{
	"cloudkms.cryptoKeyVersions.useToSign",
	"cloudkms.cryptoKeyVersions.useToVerify",
	"cloudkms.cryptoKeyVersions.viewPublicKey",
	"cloudkms.cryptoKeyVersions.useToDecrypt",
	"cloudkms.cryptoKeyVersions.useToEncrypt",
}

var _ wrapping.Wrapper = (*Wrapper)(nil)

func NewWrapper() *Wrapper {
	s := &Wrapper{
		currentKeyId: new(atomic.Value),
	}
	s.currentKeyId.Store("")
	return s
}

// SetConfig sets the fields on the Wrapper object based on values from the
// config parameter.   Environment variables take precedence over values provided
// in the config struct.
//
// Order of precedence for GCP credentials file:
// * GOOGLE_CREDENTIALS environment variable
// * `credentials` value from Value configuration file
// * GOOGLE_APPLICATION_CREDENTIALS
// (https://developers.google.com/identity/protocols/application-default-credentials)
//
// Unless the WithKeyNotRequired(true) option is provided, as a result of
// successful configuration, the wrapper's KeyId will be set to the primary
// CryptoKeyVersion.
func (s *Wrapper) SetConfig(_ context.Context, opt ...wrapping.Option) (*wrapping.WrapperConfig, error) {
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, err
	}

	s.keyNotRequired = opts.withKeyNotRequired

	s.userAgent = opts.withUserAgent

	// Do not return an error in this case. Let client initialization in
	// createClient() attempt to sort out where to get default credentials internally
	// within the SDK (e.g. checking for GOOGLE_APPLICATION_CREDENTIALS), and let
	// it error out there if none is found. This is here to establish precedence on
	// non-default input methods.
	switch {
	case os.Getenv(EnvGcpCkmsWrapperCredsPath) != "" && !opts.Options.WithDisallowEnvVars:
		s.credsPath = os.Getenv(EnvGcpCkmsWrapperCredsPath)
	case opts.withCredentials != "":
		s.credsPath = opts.withCredentials
	}

	switch {
	case os.Getenv(EnvGcpCkmsWrapperProject) != "" && !opts.Options.WithDisallowEnvVars:
		s.project = os.Getenv(EnvGcpCkmsWrapperProject)
	case opts.withProject != "":
		s.project = opts.withProject
	default:
		return nil, errors.New("'project' not found for GCP CKMS wrapper configuration")
	}

	switch {
	case os.Getenv(EnvGcpCkmsWrapperLocation) != "" && !opts.Options.WithDisallowEnvVars:
		s.location = os.Getenv(EnvGcpCkmsWrapperLocation)
	case opts.withRegion != "":
		s.location = opts.withRegion
	default:
		return nil, errors.New("'region' not found for GCP CKMS wrapper configuration")
	}

	switch {
	case os.Getenv(EnvGcpCkmsWrapperKeyRing) != "" && !opts.Options.WithDisallowEnvVars:
		s.keyRing = os.Getenv(EnvGcpCkmsWrapperKeyRing)
	case os.Getenv(EnvVaultGcpCkmsSealKeyRing) != "" && !opts.Options.WithDisallowEnvVars:
		s.keyRing = os.Getenv(EnvVaultGcpCkmsSealKeyRing)
	case opts.withKeyRing != "":
		s.keyRing = opts.withKeyRing
	default:
		return nil, errors.New("'key_ring' not found for GCP CKMS wrapper configuration")
	}

	switch {
	case os.Getenv(EnvGcpCkmsWrapperCryptoKey) != "" && !opts.Options.WithDisallowEnvVars:
		s.cryptoKey = os.Getenv(EnvGcpCkmsWrapperCryptoKey)
	case os.Getenv(EnvVaultGcpCkmsSealCryptoKey) != "" && !opts.Options.WithDisallowEnvVars:
		s.cryptoKey = os.Getenv(EnvVaultGcpCkmsSealCryptoKey)
	case opts.withCryptoKey != "":
		s.cryptoKey = opts.withCryptoKey
	case s.keyNotRequired:
		// key not required to set config
	default:
		return nil, errors.New("'crypto_key' not found for GCP CKMS wrapper configuration")
	}

	// Set the parent name for encrypt/decrypt requests
	s.parentName = fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", s.project, s.location, s.keyRing, s.cryptoKey)

	// Set and check s.client
	if s.client == nil {
		kmsClient, err := s.createClient()
		if err != nil {
			return nil, fmt.Errorf("error initializing GCP CKMS wrapper client: %w", err)
		}
		s.client = kmsClient

		// Make sure user has permissions to encrypt or sign and check if key exists
		if !s.keyNotRequired {
			ctx := context.Background()
			k, err := s.client.GetCryptoKey(ctx, &kmspb.GetCryptoKeyRequest{Name: s.parentName})
			if err != nil {
				return nil, fmt.Errorf("error checking key existence: %s", err)
			}
			s.currentKeyId.Store(k.GetPrimary().GetName())

			permissions, err := s.client.ResourceIAM(s.parentName).TestPermissions(ctx, keyPermissions)
			if err != nil {
				return nil, err
			}

			if len(permissions) == 0 {
				return nil, errors.New("permissions check failed - ensure the service account has at least " +
					"roles/cloudkms.cryptoKeyEncrypterDecrypter permissions (for encryption keys) or " +
					"roles/cloudkms.signerVerifier permissions (for signing keys)")
			}
		}
	}

	// Map that holds non-sensitive configuration info to return
	wrapConfig := new(wrapping.WrapperConfig)
	wrapConfig.Metadata = make(map[string]string)
	wrapConfig.Metadata["project"] = s.project
	wrapConfig.Metadata["region"] = s.location
	wrapConfig.Metadata["key_ring"] = s.keyRing
	wrapConfig.Metadata["crypto_key"] = s.cryptoKey

	return wrapConfig, nil
}

// Type returns the type for this particular wrapper implementation
func (s *Wrapper) Type(_ context.Context) (wrapping.WrapperType, error) {
	return wrapping.WrapperTypeGcpCkms, nil
}

// KeyId returns the last known CryptoKeyVersion which is determined when the
// wrappers is configured (Unless the WithKeyNotRequired(true) option is
// provided during configuration) or after successful encryption operations.
func (s *Wrapper) KeyId(_ context.Context) (string, error) {
	return s.currentKeyId.Load().(string), nil
}

// Encrypt is used to encrypt the master key using the the AWS CMK.
// This returns the ciphertext, and/or any errors from this
// call. This should be called after s.client has been instantiated.
// After a successful call, the wrapper's KeyId will be set to the key's id +
// it's version (example of the version appended at the very end of the key's id
// projects/<proj-id>/locations/<location-id>/keyRings/<keyring-id>/cryptoKeys/<key-id>/cryptoKeyVersions/<key-version-id>).
// Note: only the key's id (without it's version) is used when making GCP
// Encrypt/Decrypt calls.
func (s *Wrapper) Encrypt(ctx context.Context, plaintext []byte, opt ...wrapping.Option) (*wrapping.BlobInfo, error) {
	if plaintext == nil {
		return nil, errors.New("given plaintext for encryption is nil")
	}

	env, err := wrapping.EnvelopeEncrypt(plaintext, opt...)
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
	s.currentKeyId.Store(resp.Name)

	ret := &wrapping.BlobInfo{
		Ciphertext: env.Ciphertext,
		Iv:         env.Iv,
		KeyInfo: &wrapping.KeyInfo{
			Mechanism: GcpCkmsEnvelopeAesGcmEncrypt,
			// Even though we do not use the key id during decryption, store it
			// to know exactly what version was used in encryption in case we
			// want to rewrap older entries
			KeyId:      resp.Name,
			WrappedKey: resp.Ciphertext,
		},
	}

	return ret, nil
}

// Decrypt is used to decrypt the ciphertext.
func (s *Wrapper) Decrypt(ctx context.Context, in *wrapping.BlobInfo, opt ...wrapping.Option) ([]byte, error) {
	if in.Ciphertext == nil {
		return nil, fmt.Errorf("given ciphertext for decryption is nil")
	}

	// Default to mechanism used before key info was stored
	if in.KeyInfo == nil {
		in.KeyInfo = &wrapping.KeyInfo{
			Mechanism: GcpCkmsEncrypt,
		}
	}

	var plaintext []byte
	switch in.KeyInfo.Mechanism {
	case GcpCkmsEncrypt:
		resp, err := s.client.Decrypt(ctx, &kmspb.DecryptRequest{
			Name:       s.parentName,
			Ciphertext: in.Ciphertext,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt data: %w", err)
		}

		plaintext = resp.Plaintext

	case GcpCkmsEnvelopeAesGcmEncrypt:
		resp, err := s.client.Decrypt(ctx, &kmspb.DecryptRequest{
			Name:       s.parentName,
			Ciphertext: in.KeyInfo.WrappedKey,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt envelope: %w", err)
		}

		envInfo := &wrapping.EnvelopeInfo{
			Key:        resp.Plaintext,
			Iv:         in.Iv,
			Ciphertext: in.Ciphertext,
		}
		plaintext, err = wrapping.EnvelopeDecrypt(envInfo, opt...)
		if err != nil {
			return nil, fmt.Errorf("error decrypting data with envelope: %w", err)
		}

	default:
		return nil, fmt.Errorf("invalid mechanism: %d", in.KeyInfo.Mechanism)
	}

	return plaintext, nil
}

// Client returns the GCP KMS client used by the wrapper.
func (s *Wrapper) Client() *cloudkms.KeyManagementClient {
	return s.client
}

// createClient returns a configured GCP KMS client.
func (s *Wrapper) createClient() (*cloudkms.KeyManagementClient, error) {
	client, err := cloudkms.NewKeyManagementClient(context.Background(),
		option.WithCredentialsFile(s.credsPath),
		option.WithUserAgent(s.userAgent),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create KMS client: %w", err)
	}

	return client, nil
}

// KeyRingResourceName returns the relative resource name of the configured key ring.
func (s *Wrapper) KeyRingResourceName() string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", s.project, s.location, s.keyRing)
}

// LocationName returns the relative location name.
func (s *Wrapper) LocationName() string {
	return fmt.Sprintf("projects/%s/locations/%s", s.project, s.location)
}
