package awskms

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/awsutil"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault/seal"
)

const (
	// EnvAWSKMSSealKeyID is the AWS KMS key ID to use for encryption and decryption
	EnvAWSKMSSealKeyID = "VAULT_AWSKMS_SEAL_KEY_ID"
)

// AWSKMSMechanism is the method used to encrypt/decrypt in the autoseal
type AWSKMSMechanism uint32

const (
	// AWSKMSEncrypt is used to directly encrypt the data with KMS
	AWSKMSEncrypt = iota
	// AWSKMSEnvelopeAESGCMEncrypt is when a data encryption key is generated and
	// the data is encrypted with AESGCM and the key is encrypted with KMS
	AWSKMSEnvelopeAESGCMEncrypt
)

// AWSKMSSeal represents credentials and Key information for the KMS Key used to
// encryption and decryption
type AWSKMSSeal struct {
	accessKey    string
	secretKey    string
	sessionToken string
	region       string
	keyID        string
	endpoint     string

	currentKeyID *atomic.Value

	client kmsiface.KMSAPI
	logger log.Logger
}

// Ensure that we are implementing AutoSealAccess
var _ seal.Access = (*AWSKMSSeal)(nil)

// NewSeal creates a new AWSKMS seal with the provided logger
func NewSeal(logger log.Logger) *AWSKMSSeal {
	k := &AWSKMSSeal{
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
func (k *AWSKMSSeal) SetConfig(config map[string]string) (map[string]string, error) {
	if config == nil {
		config = map[string]string{}
	}

	// Check and set KeyID
	switch {
	case os.Getenv(EnvAWSKMSSealKeyID) != "":
		k.keyID = os.Getenv(EnvAWSKMSSealKeyID)
	case config["kms_key_id"] != "":
		k.keyID = config["kms_key_id"]
	default:
		return nil, fmt.Errorf("'kms_key_id' not found for AWS KMS seal configuration")
	}

	// Check and set region
	region, regionOk := config["region"]
	switch {
	case os.Getenv("AWS_REGION") != "":
		k.region = os.Getenv("AWS_REGION")
	case os.Getenv("AWS_DEFAULT_REGION") != "":
		k.region = os.Getenv("AWS_DEFAULT_REGION")
	case regionOk && region != "":
		k.region = region
	default:
		k.region = "us-east-1"

		// If available, get the region from EC2 instance metadata
		sess, err := session.NewSession(nil)
		if err != nil {
			k.logger.Warn(fmt.Sprintf("unable to begin session: %s, defaulting region to %s", err, k.region))
			break
		}

		// This will hang for ~10 seconds if the agent isn't running on an EC2 instance
		region, err := ec2metadata.New(sess).Region()
		if err != nil {
			k.logger.Warn(fmt.Sprintf("unable to retrieve region from ec2 instance metadata: %s, defaulting region to %s", err, k.region))
			break
		}
		k.region = region
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
		client, err := k.getAWSKMSClient()
		if err != nil {
			return nil, errwrap.Wrapf("error initializing AWS KMS sealclient: {{err}}", err)
		}

		// Test the client connection using provided key ID
		keyInfo, err := client.DescribeKey(&kms.DescribeKeyInput{
			KeyId: aws.String(k.keyID),
		})
		if err != nil {
			return nil, errwrap.Wrapf("error fetching AWS KMS sealkey information: {{err}}", err)
		}
		if keyInfo == nil || keyInfo.KeyMetadata == nil || keyInfo.KeyMetadata.KeyId == nil {
			return nil, errors.New("no key information returned")
		}
		k.currentKeyID.Store(aws.StringValue(keyInfo.KeyMetadata.KeyId))

		k.client = client
	}

	// Map that holds non-sensitive configuration info
	sealInfo := make(map[string]string)
	sealInfo["region"] = k.region
	sealInfo["kms_key_id"] = k.keyID
	if k.endpoint != "" {
		sealInfo["endpoint"] = k.endpoint
	}

	return sealInfo, nil
}

// Init is called during core.Initialize. No-op at the moment.
func (k *AWSKMSSeal) Init(_ context.Context) error {
	return nil
}

// Finalize is called during shutdown. This is a no-op since
// AWSKMSSeal doesn't require any cleanup.
func (k *AWSKMSSeal) Finalize(_ context.Context) error {
	return nil
}

// SealType returns the seal type for this particular seal implementation.
func (k *AWSKMSSeal) SealType() string {
	return seal.AWSKMS
}

// KeyID returns the last known key id.
func (k *AWSKMSSeal) KeyID() string {
	return k.currentKeyID.Load().(string)
}

// Encrypt is used to encrypt the master key using the the AWS CMK.
// This returns the ciphertext, and/or any errors from this
// call. This should be called after the KMS client has been instantiated.
func (k *AWSKMSSeal) Encrypt(_ context.Context, plaintext []byte) (blob *physical.EncryptedBlobInfo, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "encrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", "awskms", "encrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "encrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", "awskms", "encrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "encrypt"}, 1)
	metrics.IncrCounter([]string{"seal", "awskms", "encrypt"}, 1)

	if plaintext == nil {
		return nil, fmt.Errorf("given plaintext for encryption is nil")
	}

	env, err := seal.NewEnvelope().Encrypt(plaintext)
	if err != nil {
		return nil, errwrap.Wrapf("error wrapping data: {{err}}", err)
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
		return nil, errwrap.Wrapf("error encrypting data: {{err}}", err)
	}

	// store the current key id
	keyID := aws.StringValue(output.KeyId)
	k.currentKeyID.Store(keyID)

	ret := &physical.EncryptedBlobInfo{
		Ciphertext: env.Ciphertext,
		IV:         env.IV,
		KeyInfo: &physical.SealKeyInfo{
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
func (k *AWSKMSSeal) Decrypt(_ context.Context, in *physical.EncryptedBlobInfo) (pt []byte, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "decrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", "awskms", "decrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "decrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", "awskms", "decrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "decrypt"}, 1)
	metrics.IncrCounter([]string{"seal", "awskms", "decrypt"}, 1)

	if in == nil {
		return nil, fmt.Errorf("given input for decryption is nil")
	}

	// Default to mechanism used before key info was stored
	if in.KeyInfo == nil {
		in.KeyInfo = &physical.SealKeyInfo{
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
			return nil, errwrap.Wrapf("error decrypting data: {{err}}", err)
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
			return nil, errwrap.Wrapf("error decrypting data encryption key: {{err}}", err)
		}

		envInfo := &seal.EnvelopeInfo{
			Key:        output.Plaintext,
			IV:         in.IV,
			Ciphertext: in.Ciphertext,
		}
		plaintext, err = seal.NewEnvelope().Decrypt(envInfo)
		if err != nil {
			return nil, errwrap.Wrapf("error decrypting data: {{err}}", err)
		}

	default:
		return nil, fmt.Errorf("invalid mechanism: %d", in.KeyInfo.Mechanism)
	}

	return plaintext, nil
}

// getAWSKMSClient returns an instance of the KMS client.
func (k *AWSKMSSeal) getAWSKMSClient() (*kms.KMS, error) {
	credsConfig := &awsutil.CredentialsConfig{}

	credsConfig.AccessKey = k.accessKey
	credsConfig.SecretKey = k.secretKey
	credsConfig.SessionToken = k.sessionToken
	credsConfig.Region = k.region

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
