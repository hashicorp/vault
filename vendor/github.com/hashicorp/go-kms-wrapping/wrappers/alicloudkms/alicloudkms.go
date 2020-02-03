package alicloudkms

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials/providers"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	wrapping "github.com/hashicorp/go-kms-wrapping"
)

// These constants contain the accepted env vars; the Vault one is for backwards compat
const (
	EnvAliCloudKMSWrapperKeyID   = "ALICLOUDKMS_WRAPPER_KEY_ID"
	EnvVaultAliCloudKMSSealKeyID = "VAULT_ALICLOUDKMS_SEAL_KEY_ID"
)

// Wrapper is a Wrapper that uses AliCloud's KMS
type Wrapper struct {
	client       kmsClient
	domain       string
	keyID        string
	currentKeyID *atomic.Value
}

// Ensure that we are implementing Wrapper
var _ wrapping.Wrapper = (*Wrapper)(nil)

// NewWrapper creates a new AliCloud Wrapper
func NewWrapper(opts *wrapping.WrapperOptions) *Wrapper {
	if opts == nil {
		opts = new(wrapping.WrapperOptions)
	}
	k := &Wrapper{
		currentKeyID: new(atomic.Value),
	}
	k.currentKeyID.Store("")
	return k
}

// SetConfig sets the fields on the AliCloudKMSWrapper object based on
// values from the config parameter.
//
// Order of precedence AliCloud values:
// * Environment variable
// * Value from Vault configuration file
// * Instance metadata role (access key and secret key)
func (k *Wrapper) SetConfig(config map[string]string) (map[string]string, error) {
	if config == nil {
		config = map[string]string{}
	}

	// Check and set KeyID
	switch {
	case os.Getenv(EnvAliCloudKMSWrapperKeyID) != "":
		k.keyID = os.Getenv(EnvAliCloudKMSWrapperKeyID)
	case os.Getenv(EnvVaultAliCloudKMSSealKeyID) != "":
		k.keyID = os.Getenv(EnvVaultAliCloudKMSSealKeyID)
	case config["kms_key_id"] != "":
		k.keyID = config["kms_key_id"]
	default:
		return nil, fmt.Errorf("'kms_key_id' not found for AliCloud KMS wrapper configuration")
	}

	region := ""
	if k.client == nil {
		// Check and set region.
		region = os.Getenv("ALICLOUD_REGION")
		if region == "" {
			ok := false
			if region, ok = config["region"]; !ok {
				region = "cn-beijing"
			}
		}

		// A domain isn't required, but it can be used to override the endpoint
		// returned by the region. An example value for a domain would be:
		// "kms.us-east-1.aliyuncs.com".
		k.domain = os.Getenv("ALICLOUD_DOMAIN")
		if k.domain == "" {
			k.domain = config["domain"]
		}

		// Build the optional, configuration-based piece of the credential chain.
		credConfig := &providers.Configuration{}

		if accessKey, ok := config["access_key"]; ok {
			credConfig.AccessKeyID = accessKey
		}

		if secretKey, ok := config["secret_key"]; ok {
			credConfig.AccessKeySecret = secretKey
		} else {
			if accessSecret, ok := config["access_secret"]; ok {
				credConfig.AccessKeySecret = accessSecret
			}
		}

		credentialChain := []providers.Provider{
			providers.NewEnvCredentialProvider(),
			providers.NewConfigurationCredentialProvider(credConfig),
			providers.NewInstanceMetadataProvider(),
		}
		credProvider := providers.NewChainProvider(credentialChain)

		creds, err := credProvider.Retrieve()
		if err != nil {
			return nil, err
		}
		clientConfig := sdk.NewConfig()
		clientConfig.Scheme = "https"
		client, err := kms.NewClientWithOptions(region, clientConfig, creds)
		if err != nil {
			return nil, err
		}
		k.client = client
	}

	// Test the client connection using provided key ID
	input := kms.CreateDescribeKeyRequest()
	input.KeyId = k.keyID
	input.Domain = k.domain

	keyInfo, err := k.client.DescribeKey(input)
	if err != nil {
		return nil, fmt.Errorf("error fetching AliCloud KMS key information: %w", err)
	}
	if keyInfo == nil || keyInfo.KeyMetadata.KeyId == "" {
		return nil, errors.New("no key information returned")
	}

	// Store the current key id. If using a key alias, this will point to the actual
	// unique key that that was used for this encrypt operation.
	k.currentKeyID.Store(keyInfo.KeyMetadata.KeyId)

	// Map that holds non-sensitive configuration info
	wrapperInfo := make(map[string]string)
	wrapperInfo["region"] = region
	wrapperInfo["kms_key_id"] = k.keyID
	if k.domain != "" {
		wrapperInfo["domain"] = k.domain
	}

	return wrapperInfo, nil
}

// Init is called during core.Initialize. No-op at the moment.
func (k *Wrapper) Init(_ context.Context) error {
	return nil
}

// Finalize is called during shutdown. This is a no-op since
// AliCloudKMSWrapper doesn't require any cleanup.
func (k *Wrapper) Finalize(_ context.Context) error {
	return nil
}

// Type returns the type for this particular wrapper implementation
func (k *Wrapper) Type() string {
	return wrapping.AliCloudKMS
}

// KeyID returns the last known key id
func (k *Wrapper) KeyID() string {
	return k.currentKeyID.Load().(string)
}

// HMACKeyID returns nothing, it's here to satisfy the interface
func (k *Wrapper) HMACKeyID() string {
	return ""
}

// Encrypt is used to encrypt the master key using the the AliCloud CMK.
// This returns the ciphertext, and/or any errors from this
// call. This should be called after the KMS client has been instantiated.
func (k *Wrapper) Encrypt(_ context.Context, plaintext, aad []byte) (blob *wrapping.EncryptedBlobInfo, err error) {
	if plaintext == nil {
		return nil, fmt.Errorf("given plaintext for encryption is nil")
	}

	env, err := wrapping.NewEnvelope(nil).Encrypt(plaintext, aad)
	if err != nil {
		return nil, fmt.Errorf("error wrapping data: %w", err)
	}

	input := kms.CreateEncryptRequest()
	input.KeyId = k.keyID
	input.Plaintext = base64.StdEncoding.EncodeToString(env.Key)
	input.Domain = k.domain

	output, err := k.client.Encrypt(input)
	if err != nil {
		return nil, fmt.Errorf("error encrypting data: %w", err)
	}

	// Store the current key id.
	keyID := output.KeyId
	k.currentKeyID.Store(keyID)

	ret := &wrapping.EncryptedBlobInfo{
		Ciphertext: env.Ciphertext,
		IV:         env.IV,
		KeyInfo: &wrapping.KeyInfo{
			KeyID:      keyID,
			WrappedKey: []byte(output.CiphertextBlob),
		},
	}

	return ret, nil
}

// Decrypt is used to decrypt the ciphertext. This should be called after Init.
func (k *Wrapper) Decrypt(_ context.Context, in *wrapping.EncryptedBlobInfo, aad []byte) (pt []byte, err error) {
	if in == nil {
		return nil, fmt.Errorf("given input for decryption is nil")
	}

	// KeyID is not passed to this call because AliCloud handles this
	// internally based on the metadata stored with the encrypted data
	input := kms.CreateDecryptRequest()
	input.CiphertextBlob = string(in.KeyInfo.WrappedKey)
	input.Domain = k.domain

	output, err := k.client.Decrypt(input)
	if err != nil {
		return nil, fmt.Errorf("error decrypting data encryption key: %w", err)
	}

	keyBytes, err := base64.StdEncoding.DecodeString(output.Plaintext)
	if err != nil {
		return nil, err
	}

	envInfo := &wrapping.EnvelopeInfo{
		Key:        keyBytes,
		IV:         in.IV,
		Ciphertext: in.Ciphertext,
	}
	plaintext, err := wrapping.NewEnvelope(nil).Decrypt(envInfo, aad)
	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %w", err)
	}

	return plaintext, nil
}

type kmsClient interface {
	Decrypt(request *kms.DecryptRequest) (response *kms.DecryptResponse, err error)
	DescribeKey(request *kms.DescribeKeyRequest) (response *kms.DescribeKeyResponse, err error)
	Encrypt(request *kms.EncryptRequest) (response *kms.EncryptResponse, err error)
}
