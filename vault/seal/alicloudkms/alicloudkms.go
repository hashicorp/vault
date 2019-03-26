package alicloudkms

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials/providers"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault/seal"
)

const (
	EnvAliCloudKMSSealKeyID = "VAULT_ALICLOUDKMS_SEAL_KEY_ID"
)

type AliCloudKMSSeal struct {
	logger       log.Logger
	client       kmsClient
	domain       string
	keyID        string
	currentKeyID *atomic.Value
}

// Ensure that we are implementing AutoSealAccess
var _ seal.Access = (*AliCloudKMSSeal)(nil)

func NewSeal(logger log.Logger) *AliCloudKMSSeal {
	k := &AliCloudKMSSeal{
		logger:       logger,
		currentKeyID: new(atomic.Value),
	}
	k.currentKeyID.Store("")
	return k
}

// SetConfig sets the fields on the AliCloudKMSSeal object based on
// values from the config parameter.
//
// Order of precedence AliCloud values:
// * Environment variable
// * Value from Vault configuration file
// * Instance metadata role (access key and secret key)
func (k *AliCloudKMSSeal) SetConfig(config map[string]string) (map[string]string, error) {
	if config == nil {
		config = map[string]string{}
	}

	// Check and set KeyID
	switch {
	case os.Getenv(EnvAliCloudKMSSealKeyID) != "":
		k.keyID = os.Getenv(EnvAliCloudKMSSealKeyID)
	case config["kms_key_id"] != "":
		k.keyID = config["kms_key_id"]
	default:
		return nil, fmt.Errorf("'kms_key_id' not found for AliCloud KMS seal configuration")
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
		return nil, errwrap.Wrapf("error fetching AliCloud KMS sealkey information: {{err}}", err)
	}
	if keyInfo == nil || keyInfo.KeyMetadata.KeyId == "" {
		return nil, errors.New("no key information returned")
	}

	// Store the current key id. If using a key alias, this will point to the actual
	// unique key that that was used for this encrypt operation.
	k.currentKeyID.Store(keyInfo.KeyMetadata.KeyId)

	// Map that holds non-sensitive configuration info
	sealInfo := make(map[string]string)
	sealInfo["region"] = region
	sealInfo["kms_key_id"] = k.keyID
	if k.domain != "" {
		sealInfo["domain"] = k.domain
	}

	return sealInfo, nil
}

// Init is called during core.Initialize. No-op at the moment.
func (k *AliCloudKMSSeal) Init(_ context.Context) error {
	return nil
}

// Finalize is called during shutdown. This is a no-op since
// AliCloudKMSSeal doesn't require any cleanup.
func (k *AliCloudKMSSeal) Finalize(_ context.Context) error {
	return nil
}

// SealType returns the seal type for this particular seal implementation.
func (k *AliCloudKMSSeal) SealType() string {
	return seal.AliCloudKMS
}

// KeyID returns the last known key id.
func (k *AliCloudKMSSeal) KeyID() string {
	return k.currentKeyID.Load().(string)
}

// Encrypt is used to encrypt the master key using the the AliCloud CMK.
// This returns the ciphertext, and/or any errors from this
// call. This should be called after the KMS client has been instantiated.
func (k *AliCloudKMSSeal) Encrypt(_ context.Context, plaintext []byte) (blob *physical.EncryptedBlobInfo, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "encrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", "alicloudkms", "encrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "encrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", "alicloudkms", "encrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "encrypt"}, 1)
	metrics.IncrCounter([]string{"seal", "alicloudkms", "encrypt"}, 1)

	if plaintext == nil {
		return nil, fmt.Errorf("given plaintext for encryption is nil")
	}

	env, err := seal.NewEnvelope().Encrypt(plaintext)
	if err != nil {
		return nil, errwrap.Wrapf("error wrapping data: {{err}}", err)
	}

	input := kms.CreateEncryptRequest()
	input.KeyId = k.keyID
	input.Plaintext = base64.StdEncoding.EncodeToString(env.Key)
	input.Domain = k.domain

	output, err := k.client.Encrypt(input)
	if err != nil {
		return nil, errwrap.Wrapf("error encrypting data: {{err}}", err)
	}

	// Store the current key id.
	keyID := output.KeyId
	k.currentKeyID.Store(keyID)

	ret := &physical.EncryptedBlobInfo{
		Ciphertext: env.Ciphertext,
		IV:         env.IV,
		KeyInfo: &physical.SealKeyInfo{
			KeyID:      keyID,
			WrappedKey: []byte(output.CiphertextBlob),
		},
	}

	return ret, nil
}

// Decrypt is used to decrypt the ciphertext. This should be called after Init.
func (k *AliCloudKMSSeal) Decrypt(_ context.Context, in *physical.EncryptedBlobInfo) (pt []byte, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "decrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", "alicloudkms", "decrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "decrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", "alicloudkms", "decrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "decrypt"}, 1)
	metrics.IncrCounter([]string{"seal", "alicloudkms", "decrypt"}, 1)

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
		return nil, errwrap.Wrapf("error decrypting data encryption key: {{err}}", err)
	}

	keyBytes, err := base64.StdEncoding.DecodeString(output.Plaintext)
	if err != nil {
		return nil, err
	}

	envInfo := &seal.EnvelopeInfo{
		Key:        keyBytes,
		IV:         in.IV,
		Ciphertext: in.Ciphertext,
	}
	plaintext, err := seal.NewEnvelope().Decrypt(envInfo)
	if err != nil {
		return nil, errwrap.Wrapf("error decrypting data: {{err}}", err)
	}

	return plaintext, nil
}

type kmsClient interface {
	Decrypt(request *kms.DecryptRequest) (response *kms.DecryptResponse, err error)
	DescribeKey(request *kms.DescribeKeyRequest) (response *kms.DescribeKeyResponse, err error)
	Encrypt(request *kms.EncryptRequest) (response *kms.EncryptResponse, err error)
}
