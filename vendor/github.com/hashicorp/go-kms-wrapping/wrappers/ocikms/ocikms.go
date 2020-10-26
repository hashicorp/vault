// Copyright © 2019, Oracle and/or its affiliates.
package ocikms

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/common/auth"
	"github.com/oracle/oci-go-sdk/keymanagement"
)

const (
	// OCI KMS key ID to use for encryption and decryption
	EnvOCIKMSWrapperKeyID   = "OCIKMS_WRAPPER_KEY_ID"
	EnvVaultOCIKMSSealKeyID = "VAULT_OCIKMS_SEAL_KEY_ID"
	// OCI KMS crypto endpoint to use for encryption and decryption
	EnvOCIKMSWrapperCryptoEndpoint   = "OCIKMS_WRAPPER_CRYPTO_ENDPOINT"
	EnvVaultOCIKMSSealCryptoEndpoint = "VAULT_OCIKMS_CRYPTO_ENDPOINT"
	// OCI KMS management endpoint to manage keys
	EnvOCIKMSWrapperManagementEndpoint   = "OCIKMS_WRAPPER_MANAGEMENT_ENDPOINT"
	EnvVaultOCIKMSSealManagementEndpoint = "VAULT_OCIKMS_MANAGEMENT_ENDPOINT"
	// Maximum number of retries
	KMSMaximumNumberOfRetries = 5
	// keyID config
	KMSConfigKeyID = "key_id"
	// cryptoEndpoint config
	KMSConfigCryptoEndpoint = "crypto_endpoint"
	// managementEndpoint config
	KMSConfigManagementEndpoint = "management_endpoint"
	// authTypeAPIKey config
	KMSConfigAuthTypeAPIKey = "auth_type_api_key"
)

type Wrapper struct {
	authTypeAPIKey bool   // true for user principal, false for instance principal, default is false
	keyID          string // OCI KMS keyID

	cryptoEndpoint     string // OCI KMS crypto endpoint
	managementEndpoint string // OCI KMS management endpoint

	cryptoClient     *keymanagement.KmsCryptoClient     // OCI KMS crypto client
	managementClient *keymanagement.KmsManagementClient // OCI KMS management client

	currentKeyID *atomic.Value // Current key version which is used for encryption/decryption
}

var _ wrapping.Wrapper = (*Wrapper)(nil)

// NewWrapper creates a new Wrapper seal with the provided logger
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

func (k *Wrapper) SetConfig(config map[string]string) (map[string]string, error) {
	if config == nil {
		config = map[string]string{}
	}

	// Check and set KeyID
	switch {
	case os.Getenv(EnvOCIKMSWrapperKeyID) != "":
		k.keyID = os.Getenv(EnvOCIKMSWrapperKeyID)
	case os.Getenv(EnvVaultOCIKMSSealKeyID) != "":
		k.keyID = os.Getenv(EnvVaultOCIKMSSealKeyID)
	case config[KMSConfigKeyID] != "":
		k.keyID = config[KMSConfigKeyID]
	default:
		return nil, fmt.Errorf("'%s' not found for OCI KMS seal configuration", KMSConfigKeyID)
	}
	// Check and set cryptoEndpoint
	switch {
	case os.Getenv(EnvOCIKMSWrapperCryptoEndpoint) != "":
		k.cryptoEndpoint = os.Getenv(EnvOCIKMSWrapperCryptoEndpoint)
	case os.Getenv(EnvVaultOCIKMSSealCryptoEndpoint) != "":
		k.cryptoEndpoint = os.Getenv(EnvVaultOCIKMSSealCryptoEndpoint)
	case config[KMSConfigCryptoEndpoint] != "":
		k.cryptoEndpoint = config[KMSConfigCryptoEndpoint]
	default:
		return nil, fmt.Errorf("'%s' not found for OCI KMS seal configuration", KMSConfigCryptoEndpoint)
	}

	// Check and set managementEndpoint
	switch {
	case os.Getenv(EnvOCIKMSWrapperManagementEndpoint) != "":
		k.managementEndpoint = os.Getenv(EnvOCIKMSWrapperManagementEndpoint)
	case os.Getenv(EnvVaultOCIKMSSealManagementEndpoint) != "":
		k.managementEndpoint = os.Getenv(EnvVaultOCIKMSSealManagementEndpoint)
	case config[KMSConfigManagementEndpoint] != "":
		k.managementEndpoint = config[KMSConfigManagementEndpoint]
	default:
		return nil, fmt.Errorf("'%s' not found for OCI KMS seal configuration", KMSConfigManagementEndpoint)
	}

	// Check and set authTypeAPIKey
	var err error
	k.authTypeAPIKey = false
	authTypeAPIKeyStr := config[KMSConfigAuthTypeAPIKey]
	if authTypeAPIKeyStr != "" {
		k.authTypeAPIKey, err = strconv.ParseBool(authTypeAPIKeyStr)
		if err != nil {
			return nil, fmt.Errorf("failed parsing "+KMSConfigAuthTypeAPIKey+" parameter: %w", err)
		}
	}

	// Check and set OCI KMS crypto client
	if k.cryptoClient == nil {
		kmsCryptoClient, err := k.getOCIKMSCryptoClient()
		if err != nil {
			return nil, fmt.Errorf("error initializing OCI KMS client: %w", err)
		}
		k.cryptoClient = kmsCryptoClient
	}

	// Check and set OCI KMS management client
	if k.managementClient == nil {
		kmsManagementClient, err := k.getOCIKMSManagementClient()
		if err != nil {
			return nil, fmt.Errorf("error initializing OCI KMS client: %w", err)
		}
		k.managementClient = kmsManagementClient
	}

	// Calling Encrypt method with empty string just to validate keyId access and store current keyVersion
	encryptedBlobInfo, err := k.Encrypt(context.Background(), []byte(""), nil)
	if err != nil || encryptedBlobInfo == nil {
		return nil, fmt.Errorf("failed "+KMSConfigKeyID+" validation: %w", err)
	}

	// Map that holds non-sensitive configuration info
	wrapperInfo := make(map[string]string)
	wrapperInfo[KMSConfigKeyID] = k.keyID
	wrapperInfo[KMSConfigCryptoEndpoint] = k.cryptoEndpoint
	wrapperInfo[KMSConfigManagementEndpoint] = k.managementEndpoint
	if k.authTypeAPIKey {
		wrapperInfo["principal_type"] = "user"
	} else {
		wrapperInfo["principal_type"] = "instance"
	}

	return wrapperInfo, nil
}

func (k *Wrapper) Type() string {
	return wrapping.OCIKMS
}

func (k *Wrapper) KeyID() string {
	return k.currentKeyID.Load().(string)
}

func (k *Wrapper) HMACKeyID() string {
	return ""
}

func (k *Wrapper) Init(context.Context) error {
	return nil
}

func (k *Wrapper) Finalize(context.Context) error {
	return nil
}

func (k *Wrapper) Encrypt(ctx context.Context, plaintext, aad []byte) (*wrapping.EncryptedBlobInfo, error) {
	if plaintext == nil {
		return nil, errors.New("given plaintext for encryption is nil")
	}

	env, err := wrapping.NewEnvelope(nil).Encrypt(plaintext, aad)
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
		KeyId:     &k.keyID,
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
	k.currentKeyID.Store(keyVersion)

	ret := &wrapping.EncryptedBlobInfo{
		Ciphertext: env.Ciphertext,
		IV:         env.IV,
		KeyInfo: &wrapping.KeyInfo{
			// Storing current key version in case we want to re-wrap older entries
			KeyID:      keyVersion,
			WrappedKey: []byte(*output.Ciphertext),
		},
	}

	return ret, nil
}

func (k *Wrapper) Decrypt(ctx context.Context, in *wrapping.EncryptedBlobInfo, aad []byte) ([]byte, error) {
	if in == nil {
		return nil, fmt.Errorf("given input for decryption is nil")
	}

	requestMetadata := k.getRequestMetadata()
	cipherTextBlob := string(in.KeyInfo.WrappedKey)
	decryptedDataDetails := keymanagement.DecryptDataDetails{
		KeyId:      &k.keyID,
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
		IV:         in.IV,
		Ciphertext: in.Ciphertext,
	}

	plaintext, err := wrapping.NewEnvelope(nil).Decrypt(envInfo, aad)
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
func (k *Wrapper) getOCIKMSCryptoClient() (*keymanagement.KmsCryptoClient, error) {
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
func (k *Wrapper) getOCIKMSManagementClient() (*keymanagement.KmsManagementClient, error) {
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
		KeyId:           &k.keyID,
		RequestMetadata: requestMetadata,
	}
	getKeyResponse, err := k.managementClient.GetKey(context.Background(), getKeyInput)
	if err != nil || getKeyResponse.CurrentKeyVersion == nil {
		return "", fmt.Errorf("failed getting current key version: %w", err)
	}

	return *getKeyResponse.CurrentKeyVersion, nil
}
