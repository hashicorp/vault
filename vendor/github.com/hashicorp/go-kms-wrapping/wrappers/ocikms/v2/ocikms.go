// Copyright © 2019, Oracle and/or its affiliates.
package ocikms

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"os"
	"sync/atomic"
	"time"

	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/oracle/oci-go-sdk/v60/common"
	"github.com/oracle/oci-go-sdk/v60/common/auth"
	"github.com/oracle/oci-go-sdk/v60/keymanagement"
)

const (
	// OCI KMS key ID to use for encryption and decryption
	EnvOciKmsWrapperKeyId   = "OCIKMS_WRAPPER_KEY_ID"
	EnvVaultOciKmsSealKeyId = "VAULT_OCIKMS_SEAL_KEY_ID"
	// OCI KMS crypto endpoint to use for encryption and decryption
	EnvOciKmsWrapperCryptoEndpoint   = "OCIKMS_WRAPPER_CRYPTO_ENDPOINT"
	EnvVaultOciKmsSealCryptoEndpoint = "VAULT_OCIKMS_CRYPTO_ENDPOINT"
	// OCI KMS management endpoint to manage keys
	EnvOciKmsWrapperManagementEndpoint   = "OCIKMS_WRAPPER_MANAGEMENT_ENDPOINT"
	EnvVaultOciKmsSealManagementEndpoint = "VAULT_OCIKMS_MANAGEMENT_ENDPOINT"
	// Maximum number of retries
	KMSMaximumNumberOfRetries = 5
)

type Wrapper struct {
	authTypeAPIKey bool   // true for user principal, false for instance principal, default is false
	keyId          string // OCI KMS keyId

	cryptoEndpoint     string // OCI KMS crypto endpoint
	managementEndpoint string // OCI KMS management endpoint

	cryptoClient     *keymanagement.KmsCryptoClient     // OCI KMS crypto client
	managementClient *keymanagement.KmsManagementClient // OCI KMS management client

	currentKeyId *atomic.Value // Current key version which is used for encryption/decryption
}

var _ wrapping.Wrapper = (*Wrapper)(nil)

// NewWrapper creates a new Wrapper seal with the provided logger
func NewWrapper() *Wrapper {
	k := &Wrapper{
		currentKeyId: new(atomic.Value),
	}
	k.currentKeyId.Store("")
	return k
}

func (k *Wrapper) SetConfig(_ context.Context, opt ...wrapping.Option) (*wrapping.WrapperConfig, error) {
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, err
	}

	// Check and set KeyId
	switch {
	case os.Getenv(EnvOciKmsWrapperKeyId) != "" && !opts.Options.WithDisallowEnvVars:
		k.keyId = os.Getenv(EnvOciKmsWrapperKeyId)
	case os.Getenv(EnvVaultOciKmsSealKeyId) != "" && !opts.Options.WithDisallowEnvVars:
		k.keyId = os.Getenv(EnvVaultOciKmsSealKeyId)
	case opts.WithKeyId != "":
		k.keyId = opts.WithKeyId
	default:
		return nil, fmt.Errorf("'%s' not found for OCI KMS seal configuration", KmsConfigKeyId)
	}
	// Check and set cryptoEndpoint
	switch {
	case os.Getenv(EnvOciKmsWrapperCryptoEndpoint) != "" && !opts.Options.WithDisallowEnvVars:
		k.cryptoEndpoint = os.Getenv(EnvOciKmsWrapperCryptoEndpoint)
	case os.Getenv(EnvVaultOciKmsSealCryptoEndpoint) != "" && !opts.Options.WithDisallowEnvVars:
		k.cryptoEndpoint = os.Getenv(EnvVaultOciKmsSealCryptoEndpoint)
	case opts.withCryptoEndpoint != "":
		k.cryptoEndpoint = opts.withCryptoEndpoint
	default:
		return nil, fmt.Errorf("'%s' not found for OCI KMS seal configuration", KmsConfigCryptoEndpoint)
	}

	// Check and set managementEndpoint
	switch {
	case os.Getenv(EnvOciKmsWrapperManagementEndpoint) != "" && !opts.Options.WithDisallowEnvVars:
		k.managementEndpoint = os.Getenv(EnvOciKmsWrapperManagementEndpoint)
	case os.Getenv(EnvVaultOciKmsSealManagementEndpoint) != "" && !opts.Options.WithDisallowEnvVars:
		k.managementEndpoint = os.Getenv(EnvVaultOciKmsSealManagementEndpoint)
	case opts.withManagementEndpoint != "":
		k.managementEndpoint = opts.withManagementEndpoint
	default:
		return nil, fmt.Errorf("'%s' not found for OCI KMS seal configuration", KmsConfigManagementEndpoint)
	}

	k.authTypeAPIKey = opts.withAuthTypeApiKey

	// Check and set OCI KMS crypto client
	if k.cryptoClient == nil {
		kmsCryptoClient, err := k.getOciKmsCryptoClient()
		if err != nil {
			return nil, fmt.Errorf("error initializing OCI KMS client: %w", err)
		}
		k.cryptoClient = kmsCryptoClient
	}

	// Check and set OCI KMS management client
	if k.managementClient == nil {
		kmsManagementClient, err := k.getOciKmsManagementClient()
		if err != nil {
			return nil, fmt.Errorf("error initializing OCI KMS client: %w", err)
		}
		k.managementClient = kmsManagementClient
	}

	// Calling Encrypt method with empty string just to validate keyId access and store current keyVersion
	encryptedBlobInfo, err := k.Encrypt(context.Background(), []byte(""), nil)
	if err != nil || encryptedBlobInfo == nil {
		return nil, fmt.Errorf("failed "+KmsConfigKeyId+" validation: %w", err)
	}

	// Map that holds non-sensitive configuration info
	wrapConfig := new(wrapping.WrapperConfig)
	wrapConfig.Metadata = make(map[string]string)
	wrapConfig.Metadata[KmsConfigKeyId] = k.keyId
	wrapConfig.Metadata[KmsConfigCryptoEndpoint] = k.cryptoEndpoint
	wrapConfig.Metadata[KmsConfigManagementEndpoint] = k.managementEndpoint
	if k.authTypeAPIKey {
		wrapConfig.Metadata["principal_type"] = "user"
	} else {
		wrapConfig.Metadata["principal_type"] = "instance"
	}

	return wrapConfig, nil
}

func (k *Wrapper) Type(_ context.Context) (wrapping.WrapperType, error) {
	return wrapping.WrapperTypeOciKms, nil
}

func (k *Wrapper) KeyId(_ context.Context) (string, error) {
	return k.currentKeyId.Load().(string), nil
}

func (k *Wrapper) Encrypt(ctx context.Context, plaintext []byte, opt ...wrapping.Option) (*wrapping.BlobInfo, error) {
	if plaintext == nil {
		return nil, errors.New("given plaintext for encryption is nil")
	}

	env, err := wrapping.EnvelopeEncrypt(plaintext, opt...)
	if err != nil {
		return nil, fmt.Errorf("error wrapping data: %w", err)
	}

	if k.cryptoClient == nil {
		return nil, errors.New("nil client")
	}

	// OCI KMS required base64 encrypted plain text before sending to the service
	encodedKey := base64.StdEncoding.EncodeToString(env.Key)

	// Build Encrypt Request
	requestMetadata := k.getRequestMetadata()
	encryptedDataDetails := keymanagement.EncryptDataDetails{
		KeyId:     &k.keyId,
		Plaintext: &encodedKey,
	}

	input := keymanagement.EncryptRequest{
		EncryptDataDetails: encryptedDataDetails,
		RequestMetadata:    requestMetadata,
	}
	output, err := k.cryptoClient.Encrypt(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error encrypting data: %w", err)
	}

	// Note: It is potential a timing issue if the key gets rotated between this
	// getCurrentKeyVersion operation and above Encrypt operation
	keyVersion, err := k.getCurrentKeyVersion()
	if err != nil {
		return nil, fmt.Errorf("error getting current key version: %w", err)
	}
	// Update key version
	k.currentKeyId.Store(keyVersion)

	ret := &wrapping.BlobInfo{
		Ciphertext: env.Ciphertext,
		Iv:         env.Iv,
		KeyInfo: &wrapping.KeyInfo{
			// Storing current key version in case we want to re-wrap older entries
			KeyId:      keyVersion,
			WrappedKey: []byte(*output.Ciphertext),
		},
	}

	return ret, nil
}

func (k *Wrapper) Decrypt(ctx context.Context, in *wrapping.BlobInfo, opt ...wrapping.Option) ([]byte, error) {
	if in == nil {
		return nil, fmt.Errorf("given input for decryption is nil")
	}

	requestMetadata := k.getRequestMetadata()
	cipherTextBlob := string(in.KeyInfo.WrappedKey)
	decryptedDataDetails := keymanagement.DecryptDataDetails{
		KeyId:      &k.keyId,
		Ciphertext: &cipherTextBlob,
	}
	input := keymanagement.DecryptRequest{
		DecryptDataDetails: decryptedDataDetails,
		RequestMetadata:    requestMetadata,
	}
	output, err := k.cryptoClient.Decrypt(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %w", err)
	}
	envelopeKey, err := base64.StdEncoding.DecodeString(*output.Plaintext)
	if err != nil {
		return nil, fmt.Errorf("error base64 decrypting data: %w", err)
	}
	envInfo := &wrapping.EnvelopeInfo{
		Key:        envelopeKey,
		Iv:         in.Iv,
		Ciphertext: in.Ciphertext,
	}

	plaintext, err := wrapping.EnvelopeDecrypt(envInfo, opt...)
	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %w", err)
	}

	return plaintext, nil
}

func (k *Wrapper) getConfigProvider() (common.ConfigurationProvider, error) {
	var cp common.ConfigurationProvider
	var err error
	if k.authTypeAPIKey {
		cp = common.DefaultConfigProvider()
	} else {
		cp, err = auth.InstancePrincipalConfigurationProvider()
		if err != nil {
			return nil, fmt.Errorf("failed creating InstancePrincipalConfigurationProvider: %w", err)
		}
	}
	return cp, nil
}

// Build OCI KMS crypto client
func (k *Wrapper) getOciKmsCryptoClient() (*keymanagement.KmsCryptoClient, error) {
	cp, err := k.getConfigProvider()
	if err != nil {
		return nil, fmt.Errorf("failed creating configuration provider: %w", err)
	}

	// Build crypto client
	kmsCryptoClient, err := keymanagement.NewKmsCryptoClientWithConfigurationProvider(cp, k.cryptoEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed creating NewKmsCryptoClientWithConfigurationProvider: %w", err)
	}

	return &kmsCryptoClient, nil
}

// Build OCI KMS management client
func (k *Wrapper) getOciKmsManagementClient() (*keymanagement.KmsManagementClient, error) {
	cp, err := k.getConfigProvider()
	if err != nil {
		return nil, fmt.Errorf("failed creating configuration provider: %w", err)
	}

	// Build crypto client
	kmsManagementClient, err := keymanagement.NewKmsManagementClientWithConfigurationProvider(cp, k.managementEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed creating NewKmsCryptoClientWithConfigurationProvider: %w", err)
	}

	return &kmsManagementClient, nil
}

// Request metadata includes retry policy
func (k *Wrapper) getRequestMetadata() common.RequestMetadata {
	// Only retry for 5xx errors
	retryOn5xxFunc := func(r common.OCIOperationResponse) bool {
		return r.Error != nil && r.Response.HTTPResponse().StatusCode >= 500
	}
	return getRequestMetadataWithCustomizedRetryPolicy(retryOn5xxFunc)
}

func getRequestMetadataWithCustomizedRetryPolicy(fn func(r common.OCIOperationResponse) bool) common.RequestMetadata {
	return common.RequestMetadata{
		RetryPolicy: getExponentialBackoffRetryPolicy(uint(KMSMaximumNumberOfRetries), fn),
	}
}

func getExponentialBackoffRetryPolicy(n uint, fn func(r common.OCIOperationResponse) bool) *common.RetryPolicy {
	// The duration between each retry operation, you might want to wait longer each time the retry fails
	exponentialBackoff := func(r common.OCIOperationResponse) time.Duration {
		return time.Duration(math.Pow(float64(2), float64(r.AttemptNumber-1))) * time.Second
	}
	policy := common.NewRetryPolicy(n, fn, exponentialBackoff)
	return &policy
}

func (k *Wrapper) getCurrentKeyVersion() (string, error) {
	if k.managementClient == nil {
		return "", fmt.Errorf("managementClient has not yet initialized")
	}
	requestMetadata := k.getRequestMetadata()
	getKeyInput := keymanagement.GetKeyRequest{
		KeyId:           &k.keyId,
		RequestMetadata: requestMetadata,
	}
	getKeyResponse, err := k.managementClient.GetKey(context.Background(), getKeyInput)
	if err != nil || getKeyResponse.CurrentKeyVersion == nil {
		return "", fmt.Errorf("failed getting current key version: %w", err)
	}

	return *getKeyResponse.CurrentKeyVersion, nil
}
