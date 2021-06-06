package awskms

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/vault/sdk/helper/awsutil"
)

// These constants contain the accepted env vars; the Vault one is for backwards compat
const (
	EnvAWSKMSWrapperKeyID   = "AWSKMS_WRAPPER_KEY_ID"
	EnvVaultAWSKMSSealKeyID = "VAULT_AWSKMS_SEAL_KEY_ID"
)

const (
	// AWSKMSEncrypt is used to directly encrypt the data with KMS
	AWSKMSEncrypt = iota
	// AWSKMSEnvelopeAESGCMEncrypt is when a data encryption key is generated and
	// the data is encrypted with AESGCM and the key is encrypted with KMS
	AWSKMSEnvelopeAESGCMEncrypt
)

// Wrapper represents credentials and Key information for the KMS Key used to
// encryption and decryption
type Wrapper struct {
	accessKey    string
	secretKey    string
	sessionToken string
	region       string
	keyID        string
	endpoint     string

	currentKeyID *atomic.Value

	client kmsiface.KMSAPI

	logger hclog.Logger
}

// Ensure that we are implementing Wrapper
var _ wrapping.Wrapper = (*Wrapper)(nil)

// NewWrapper creates a new AWSKMS wrapper with the provided options
func NewWrapper(opts *wrapping.WrapperOptions) *Wrapper {
	if opts == nil {
		opts = new(wrapping.WrapperOptions)
	}
	k := &Wrapper{
		currentKeyID: new(atomic.Value),
		logger:       opts.Logger,
	}
	k.currentKeyID.Store("")
	return k
}

// SetConfig sets the fields on the Wrapper object based on
// values from the config parameter.
//
// Order of precedence AWS values:
// * Environment variable
// * Value from Vault configuration file
// * Instance metadata role (access key and secret key)
// * Default values
func (k *Wrapper) SetConfig(config map[string]string) (map[string]string, error) {
	if config == nil {
		config = map[string]string{}
	}

	// Check and set KeyID
	switch {
	case os.Getenv(EnvAWSKMSWrapperKeyID) != "":
		k.keyID = os.Getenv(EnvAWSKMSWrapperKeyID)
	case os.Getenv(EnvVaultAWSKMSSealKeyID) != "":
		k.keyID = os.Getenv(EnvVaultAWSKMSSealKeyID)
	case config["kms_key_id"] != "":
		k.keyID = config["kms_key_id"]
	default:
		return nil, fmt.Errorf("'kms_key_id' not found for AWS KMS wrapper configuration")
	}

	// Please see GetRegion for an explanation of the order in which region is parsed.
	var err error
	k.region, err = awsutil.GetRegion(config["region"])
	if err != nil {
		return nil, err
	}

	// Check and set AWS access key, secret key, and session token
	k.accessKey = config["access_key"]
	k.secretKey = config["secret_key"]
	k.sessionToken = config["session_token"]

	k.endpoint = os.Getenv("AWS_KMS_ENDPOINT")
	if k.endpoint == "" {
		if endpoint, ok := config["endpoint"]; ok {
			k.endpoint = endpoint
		}
	}

	// Check and set k.client
	if k.client == nil {
		client, err := k.GetAWSKMSClient()
		if err != nil {
			return nil, fmt.Errorf("error initializing AWS KMS wrapping client: %w", err)
		}

		// Test the client connection using provided key ID
		keyInfo, err := client.DescribeKey(&kms.DescribeKeyInput{
			KeyId: aws.String(k.keyID),
		})
		if err != nil {
			return nil, fmt.Errorf("error fetching AWS KMS wrapping key information: %w", err)
		}
		if keyInfo == nil || keyInfo.KeyMetadata == nil || keyInfo.KeyMetadata.KeyId == nil {
			return nil, errors.New("no key information returned")
		}
		k.currentKeyID.Store(aws.StringValue(keyInfo.KeyMetadata.KeyId))

		k.client = client
	}

	// Map that holds non-sensitive configuration info
	wrappingInfo := make(map[string]string)
	wrappingInfo["region"] = k.region
	wrappingInfo["kms_key_id"] = k.keyID
	if k.endpoint != "" {
		wrappingInfo["endpoint"] = k.endpoint
	}

	return wrappingInfo, nil
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
	return wrapping.AWSKMS
}

// KeyID returns the last known key id
func (k *Wrapper) KeyID() string {
	return k.currentKeyID.Load().(string)
}

// HMACKeyID returns the last known HMAC key id
func (k *Wrapper) HMACKeyID() string {
	return ""
}

// Encrypt is used to encrypt the master key using the the AWS CMK.
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

	if k.client == nil {
		return nil, fmt.Errorf("nil client")
	}

	input := &kms.EncryptInput{
		KeyId:     aws.String(k.keyID),
		Plaintext: env.Key,
	}
	output, err := k.client.Encrypt(input)
	if err != nil {
		return nil, fmt.Errorf("error encrypting data: %w", err)
	}

	// Store the current key id
	//
	// When using a key alias, this will return the actual underlying key id
	// used for encryption.  This is helpful if you are looking to reencyrpt
	// your data when it is not using the latest key id. See these docs relating
	// to key rotation https://docs.aws.amazon.com/kms/latest/developerguide/rotate-keys.html
	keyID := aws.StringValue(output.KeyId)
	k.currentKeyID.Store(keyID)

	ret := &wrapping.EncryptedBlobInfo{
		Ciphertext: env.Ciphertext,
		IV:         env.IV,
		KeyInfo: &wrapping.KeyInfo{
			Mechanism: AWSKMSEnvelopeAESGCMEncrypt,
			// Even though we do not use the key id during decryption, store it
			// to know exactly the specific key used in encryption in case we
			// want to rewrap older entries
			KeyID:      keyID,
			WrappedKey: output.CiphertextBlob,
		},
	}

	return ret, nil
}

// Decrypt is used to decrypt the ciphertext. This should be called after Init.
func (k *Wrapper) Decrypt(_ context.Context, in *wrapping.EncryptedBlobInfo, aad []byte) (pt []byte, err error) {
	if in == nil {
		return nil, fmt.Errorf("given input for decryption is nil")
	}

	// Default to mechanism used before key info was stored
	if in.KeyInfo == nil {
		in.KeyInfo = &wrapping.KeyInfo{
			Mechanism: AWSKMSEncrypt,
		}
	}

	var plaintext []byte
	switch in.KeyInfo.Mechanism {
	case AWSKMSEncrypt:
		input := &kms.DecryptInput{
			CiphertextBlob: in.Ciphertext,
		}

		output, err := k.client.Decrypt(input)
		if err != nil {
			return nil, fmt.Errorf("error decrypting data: %w", err)
		}
		plaintext = output.Plaintext

	case AWSKMSEnvelopeAESGCMEncrypt:
		// KeyID is not passed to this call because AWS handles this
		// internally based on the metadata stored with the encrypted data
		input := &kms.DecryptInput{
			CiphertextBlob: in.KeyInfo.WrappedKey,
		}
		output, err := k.client.Decrypt(input)
		if err != nil {
			return nil, fmt.Errorf("error decrypting data encryption key: %w", err)
		}

		envInfo := &wrapping.EnvelopeInfo{
			Key:        output.Plaintext,
			IV:         in.IV,
			Ciphertext: in.Ciphertext,
		}
		plaintext, err = wrapping.NewEnvelope(nil).Decrypt(envInfo, aad)
		if err != nil {
			return nil, fmt.Errorf("error decrypting data: %w", err)
		}

	default:
		return nil, fmt.Errorf("invalid mechanism: %d", in.KeyInfo.Mechanism)
	}

	return plaintext, nil
}

// GetAWSKMSClient returns an instance of the KMS client.
func (k *Wrapper) GetAWSKMSClient() (*kms.KMS, error) {
	credsConfig := &awsutil.CredentialsConfig{}

	credsConfig.AccessKey = k.accessKey
	credsConfig.SecretKey = k.secretKey
	credsConfig.SessionToken = k.sessionToken
	credsConfig.Region = k.region
	credsConfig.Logger = k.logger

	credsConfig.HTTPClient = cleanhttp.DefaultClient()

	creds, err := credsConfig.GenerateCredentialChain()
	if err != nil {
		return nil, err
	}

	awsConfig := &aws.Config{
		Credentials: creds,
		Region:      aws.String(credsConfig.Region),
		HTTPClient:  cleanhttp.DefaultClient(),
	}

	if k.endpoint != "" {
		awsConfig.Endpoint = aws.String(k.endpoint)
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	client := kms.New(sess)

	return client, nil
}
