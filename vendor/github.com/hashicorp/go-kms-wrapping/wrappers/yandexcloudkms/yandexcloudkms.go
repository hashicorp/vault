package yandexcloudkms

import (
	"context"
	"fmt"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/kms/v1"
	"github.com/yandex-cloud/go-sdk/iamkey"
	"os"
	"sync/atomic"

	wrapping "github.com/hashicorp/go-kms-wrapping"
	ycsdk "github.com/yandex-cloud/go-sdk"
)

const (
	// Accepted env vars
	EnvYandexCloudOAuthToken            = "YANDEXCLOUD_OAUTH_TOKEN"
	EnvYandexCloudServiceAccountKeyFile = "YANDEXCLOUD_SERVICE_ACCOUNT_KEY_FILE"
	EnvYandexCloudKMSKeyID              = "YANDEXCLOUD_KMS_KEY_ID"

	// Accepted config parameters
	CfgYandexCloudOAuthToken            = "oauth_token"
	CfgYandexCloudServiceAccountKeyFile = "service_account_key_file"
	CfgYandexCloudKMSKeyID              = "kms_key_id"
)

// Wrapper represents credentials and key information for the KMS Key used to
// encryption and decryption
type Wrapper struct {
	client           kms.SymmetricCryptoServiceClient
	keyID            string
	currentVersionID *atomic.Value
}

// Ensure that we are implementing Wrapper
var _ wrapping.Wrapper = (*Wrapper)(nil)

// NewWrapper creates a new Yandex.Cloud wrapper
func NewWrapper(opts *wrapping.WrapperOptions) *Wrapper {
	if opts == nil {
		opts = new(wrapping.WrapperOptions)
	}
	k := &Wrapper{
		currentVersionID: new(atomic.Value),
	}
	k.currentVersionID.Store("")
	return k
}

// SetConfig sets the fields on the Wrapper object based on
// values from the config parameter.
//
// Order of precedence Yandex.Cloud values:
// * Environment variable
// * Value from Vault configuration file
// * Compute Instance metadata
func (k *Wrapper) SetConfig(config map[string]string) (map[string]string, error) {
	if config == nil {
		config = map[string]string{}
	}

	// Check and set KeyID
	k.keyID = coalesce(os.Getenv(EnvYandexCloudKMSKeyID), config[CfgYandexCloudKMSKeyID])
	if k.keyID == "" {
		return nil, fmt.Errorf(
			"neither '%s' environment variable nor '%s' config parameter is set",
			EnvYandexCloudKMSKeyID, CfgYandexCloudKMSKeyID,
		)
	}

	// Check and set k.client
	if k.client == nil {
		client, err := getYandexCloudKMSClient(
			coalesce(os.Getenv(EnvYandexCloudOAuthToken), config[CfgYandexCloudOAuthToken]),
			coalesce(os.Getenv(EnvYandexCloudServiceAccountKeyFile), config[CfgYandexCloudServiceAccountKeyFile]),
		)
		if err != nil {
			return nil, fmt.Errorf("error initializing Yandex.Cloud KMS client: %w", err)
		}

		if err := k.setClient(client); err != nil {
			return nil, fmt.Errorf("error setting Yandex.Cloud KMS client: %w", err)
		}
	}

	// Map that holds non-sensitive configuration info
	wrappingInfo := make(map[string]string)
	wrappingInfo["kms_key_id"] = k.keyID

	return wrappingInfo, nil
}

func (k *Wrapper) setClient(client kms.SymmetricCryptoServiceClient) error {
	k.client = client

	// Make sure all the required permissions are granted (also checks if key exists)
	_, err := k.Encrypt(context.Background(), []byte("go-kms-wrapping-test"), nil)
	if err != nil {
		return fmt.Errorf(
			"failed to encrypt with Yandex.Cloud KMS key - ensure the key exists and permission to encrypt the key is granted: %w", err)
	}

	return nil
}

// Init is called during core.Initialize. No-op at the moment.
func (k *Wrapper) Init(_ context.Context) error {
	return nil
}

// Finalize is called during shutdown. This is a no-op since
// Wrapper doesn't require any cleanup.
func (k *Wrapper) Finalize(_ context.Context) error {
	return nil
}

// Type returns the wrapping type for this particular Wrapper implementation
func (k *Wrapper) Type() string {
	return wrapping.YandexCloudKMS
}

// KeyID returns the last known key id
func (k *Wrapper) KeyID() string {
	return k.currentVersionID.Load().(string)
}

// HMACKeyID returns the last known HMAC key id
func (k *Wrapper) HMACKeyID() string {
	return ""
}

// Encrypt is used to encrypt the master key using Yandex.Cloud symmetric key.
// This returns the ciphertext, and/or any errors from this
// call. This should be called after the KMS client has been instantiated.
func (k *Wrapper) Encrypt(ctx context.Context, plaintext, aad []byte) (blob *wrapping.EncryptedBlobInfo, err error) {
	if plaintext == nil {
		return nil, fmt.Errorf("given plaintext for encryption is nil")
	}

	env, err := wrapping.NewEnvelope(nil).Encrypt(plaintext, aad)
	if err != nil {
		return nil, fmt.Errorf("error wrapping data: %w", err)
	}

	if k.client == nil {
		return nil, fmt.Errorf("nil client")
	}

	encryptResponse, err := k.client.Encrypt(
		ctx,
		&kms.SymmetricEncryptRequest{
			KeyId:     k.keyID,
			Plaintext: env.Key,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error encrypting data: %w", err)
	}

	// Store the current version id
	k.currentVersionID.Store(encryptResponse.VersionId)

	ret := &wrapping.EncryptedBlobInfo{
		Ciphertext: env.Ciphertext,
		IV:         env.IV,
		KeyInfo: &wrapping.KeyInfo{
			// Even though we do not use the version id during decryption, store it
			// to know exactly the specific key version used in encryption in case we
			// want to rewrap older entries
			KeyID:      encryptResponse.VersionId,
			WrappedKey: encryptResponse.Ciphertext,
		},
	}

	return ret, nil
}

// Decrypt is used to decrypt the ciphertext. This should be called after Init.
func (k *Wrapper) Decrypt(ctx context.Context, in *wrapping.EncryptedBlobInfo, aad []byte) (pt []byte, err error) {
	if in == nil {
		return nil, fmt.Errorf("given input for decryption is nil")
	}

	decryptResponse, err := k.client.Decrypt(
		ctx,
		&kms.SymmetricDecryptRequest{
			KeyId:      k.keyID,
			Ciphertext: in.KeyInfo.WrappedKey,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error decrypting data encryption key: %w", err)
	}

	envInfo := &wrapping.EnvelopeInfo{
		Key:        decryptResponse.Plaintext,
		IV:         in.IV,
		Ciphertext: in.Ciphertext,
	}
	plaintext, err := wrapping.NewEnvelope(nil).Decrypt(envInfo, aad)
	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %w", err)
	}

	return plaintext, nil
}

// GetYandexCloudKMSClient returns an instance of the KMS client.
func getYandexCloudKMSClient(oauthToken string, serviceAccountKeyFile string) (kms.SymmetricCryptoServiceClient, error) {
	credentials, err := getCredentials(oauthToken, serviceAccountKeyFile)
	if err != nil {
		return nil, fmt.Errorf("error getting credentials: %w", err)
	}

	sdk, err := ycsdk.Build(
		context.Background(),
		ycsdk.Config{Credentials: credentials},
	)
	if err != nil {
		return nil, fmt.Errorf("error building Yandex.Cloud SDK instance: %w", err)
	}

	return sdk.KMSCrypto().SymmetricCrypto(), nil
}

func getCredentials(oauthToken string, serviceAccountKeyFile string) (ycsdk.Credentials, error) {
	if oauthToken != "" && serviceAccountKeyFile != "" {
		return nil, fmt.Errorf("error configuring authentication: both OAuth token and service account key file are specified")
	}

	// Yandex account authentication (via Oauth token)
	if oauthToken != "" {
		return ycsdk.OAuthToken(oauthToken), nil
	}

	// Service account authentication (via authorized key)
	if serviceAccountKeyFile != "" {
		key, err := iamkey.ReadFromJSONFile(serviceAccountKeyFile)
		if err != nil {
			return nil, fmt.Errorf("error reading service account key file: %w", err)
		}
		return ycsdk.ServiceAccountKey(key)
	}

	// Compute Instance Service Account authentication
	return ycsdk.InstanceServiceAccount(), nil
}

func coalesce(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
