// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
)

// These constants contain the accepted env vars; the Vault one is for backwards compat
const (
	EnvAwsKmsWrapperKeyId   = "AWSKMS_WRAPPER_KEY_ID"
	EnvVaultAwsKmsSealKeyId = "VAULT_AWSKMS_SEAL_KEY_ID"
)

const (
	// AwsKmsEncrypt is used to directly encrypt the data with KMS
	AwsKmsEncrypt = iota
	// AwsKmsEnvelopeAesGcmEncrypt is when a data encryption key is generated and
	// the data is encrypted with AES-GCM and the key is encrypted with KMS
	AwsKmsEnvelopeAesGcmEncrypt
)

// Wrapper represents credentials and Key information for the KMS Key used to
// encryption and decryption
type Wrapper struct {
	accessKey            string
	secretKey            string
	sessionToken         string
	region               string
	keyId                string
	endpoint             string
	sharedCredsFilename  string
	sharedCredsProfile   string
	roleArn              string
	roleSessionName      string
	webIdentityTokenFile string
	keyNotRequired       bool

	currentKeyId *atomic.Value

	client kmsiface.KMSAPI

	logger hclog.Logger
}

// Ensure that we are implementing Wrapper
var _ wrapping.Wrapper = (*Wrapper)(nil)

// NewWrapper creates a new AwsKms wrapper with the provided options
func NewWrapper() *Wrapper {
	k := &Wrapper{
		currentKeyId: new(atomic.Value),
	}
	k.currentKeyId.Store("")
	return k
}

// SetConfig sets the fields on the Wrapper object based on
// values from the config parameter.
//
// Order of precedence AWS values:
// * Environment variable
// * Passed in config map
// * Instance metadata role (access key and secret key)
// * Default values
func (k *Wrapper) SetConfig(_ context.Context, opt ...wrapping.Option) (*wrapping.WrapperConfig, error) {
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, err
	}

	k.keyNotRequired = opts.withKeyNotRequired
	k.logger = opts.withLogger

	// Check and set KeyId
	switch {
	case os.Getenv(EnvAwsKmsWrapperKeyId) != "" && !opts.withDisallowEnvVars:
		k.keyId = os.Getenv(EnvAwsKmsWrapperKeyId)
	case os.Getenv(EnvVaultAwsKmsSealKeyId) != "" && !opts.withDisallowEnvVars:
		k.keyId = os.Getenv(EnvVaultAwsKmsSealKeyId)
	case opts.WithKeyId != "":
		k.keyId = opts.WithKeyId
	case k.keyNotRequired:
		// key not required to set config
	default:
		return nil, fmt.Errorf("key id not found in env or config for aws kms wrapper configuration")
	}

	k.currentKeyId.Store(k.keyId)

	// Please see GetRegion for an explanation of the order in which region is parsed.
	k.region, err = awsutil.GetRegion(opts.withRegion)
	if err != nil {
		return nil, err
	}

	// Check and set AWS access key, secret key, and session token
	k.accessKey = opts.withAccessKey
	k.secretKey = opts.withSecretKey
	k.sessionToken = opts.withSessionToken
	k.sharedCredsFilename = opts.withSharedCredsFilename
	k.sharedCredsProfile = opts.withSharedCredsProfile
	k.webIdentityTokenFile = opts.withWebIdentityTokenFile
	k.roleSessionName = opts.withRoleSessionName
	k.roleArn = opts.withRoleArn

	if !opts.withDisallowEnvVars {
		k.endpoint = os.Getenv("AWS_KMS_ENDPOINT")
	}
	if k.endpoint == "" {
		k.endpoint = opts.withEndpoint
	}

	// Check and set k.client
	if k.client == nil {
		client, err := k.GetAwsKmsClient()
		if err != nil {
			return nil, fmt.Errorf("error initializing AWS KMS wrapping client: %w", err)
		}

		if !k.keyNotRequired {
			// Test the client connection using provided key ID
			keyInfo, err := client.DescribeKey(&kms.DescribeKeyInput{
				KeyId: aws.String(k.keyId),
			})
			if err != nil {
				return nil, fmt.Errorf("error fetching AWS KMS wrapping key information: %w", err)
			}
			if keyInfo == nil || keyInfo.KeyMetadata == nil || keyInfo.KeyMetadata.KeyId == nil {
				return nil, errors.New("no key information returned")
			}
			k.currentKeyId.Store(aws.StringValue(keyInfo.KeyMetadata.KeyId))
		}

		k.client = client
	}

	// Map that holds non-sensitive configuration info
	wrapConfig := new(wrapping.WrapperConfig)
	wrapConfig.Metadata = make(map[string]string)
	wrapConfig.Metadata["region"] = k.region
	wrapConfig.Metadata["key_id"] = k.keyId
	if k.endpoint != "" {
		wrapConfig.Metadata["endpoint"] = k.endpoint
	}

	return wrapConfig, nil
}

// Type returns the wrapping type for this particular Wrapper implementation
func (k *Wrapper) Type(_ context.Context) (wrapping.WrapperType, error) {
	return wrapping.WrapperTypeAwsKms, nil
}

// KeyId returns the last known key id
func (k *Wrapper) KeyId(_ context.Context) (string, error) {
	return k.currentKeyId.Load().(string), nil
}

// Encrypt is used to encrypt the master key using the the AWS CMK.
// This returns the ciphertext, and/or any errors from this
// call. This should be called after the KMS client has been instantiated.
func (k *Wrapper) Encrypt(_ context.Context, plaintext []byte, opt ...wrapping.Option) (*wrapping.BlobInfo, error) {
	if plaintext == nil {
		return nil, fmt.Errorf("given plaintext for encryption is nil")
	}

	env, err := wrapping.EnvelopeEncrypt(plaintext, opt...)
	if err != nil {
		return nil, fmt.Errorf("error wrapping data: %w", err)
	}

	if k.client == nil {
		return nil, fmt.Errorf("nil client")
	}

	input := &kms.EncryptInput{
		KeyId:     aws.String(k.keyId),
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
	keyId := aws.StringValue(output.KeyId)
	k.currentKeyId.Store(keyId)

	ret := &wrapping.BlobInfo{
		Ciphertext: env.Ciphertext,
		Iv:         env.Iv,
		KeyInfo: &wrapping.KeyInfo{
			Mechanism: AwsKmsEnvelopeAesGcmEncrypt,
			// Even though we do not use the key id during decryption, store it
			// to know exactly the specific key used in encryption in case we
			// want to rewrap older entries
			KeyId:      keyId,
			WrappedKey: output.CiphertextBlob,
		},
	}

	return ret, nil
}

// Decrypt is used to decrypt the ciphertext. This should be called after Init.
func (k *Wrapper) Decrypt(_ context.Context, in *wrapping.BlobInfo, opt ...wrapping.Option) ([]byte, error) {
	if in == nil {
		return nil, fmt.Errorf("given input for decryption is nil")
	}

	// Default to mechanism used before key info was stored
	if in.KeyInfo == nil {
		in.KeyInfo = &wrapping.KeyInfo{
			Mechanism: AwsKmsEncrypt,
		}
	}

	var plaintext []byte
	switch in.KeyInfo.Mechanism {
	case AwsKmsEncrypt:
		input := &kms.DecryptInput{
			CiphertextBlob: in.Ciphertext,
		}

		output, err := k.client.Decrypt(input)
		if err != nil {
			return nil, fmt.Errorf("error decrypting data: %w", err)
		}
		plaintext = output.Plaintext

	case AwsKmsEnvelopeAesGcmEncrypt:
		// KeyId is not passed to this call because AWS handles this
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
			Iv:         in.Iv,
			Ciphertext: in.Ciphertext,
		}
		plaintext, err = wrapping.EnvelopeDecrypt(envInfo, opt...)
		if err != nil {
			return nil, fmt.Errorf("error decrypting data: %w", err)
		}

	default:
		return nil, fmt.Errorf("invalid mechanism: %d", in.KeyInfo.Mechanism)
	}

	return plaintext, nil
}

// Client returns the AWS KMS client used by the wrapper.
func (k *Wrapper) Client() kmsiface.KMSAPI {
	return k.client
}

// GetAwsKmsClient returns an instance of the KMS client.
func (k *Wrapper) GetAwsKmsClient() (*kms.KMS, error) {
	credsConfig := &awsutil.CredentialsConfig{}

	credsConfig.AccessKey = k.accessKey
	credsConfig.SecretKey = k.secretKey
	credsConfig.SessionToken = k.sessionToken
	credsConfig.Filename = k.sharedCredsFilename
	credsConfig.Profile = k.sharedCredsProfile
	credsConfig.RoleARN = k.roleArn
	credsConfig.RoleSessionName = k.roleSessionName
	credsConfig.WebIdentityTokenFile = k.webIdentityTokenFile
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
