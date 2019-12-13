// Copyright © 2019, Oracle and/or its affiliates.
package ocikms

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/seal"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/common/auth"
	"github.com/oracle/oci-go-sdk/keymanagement"
)

const (
	// OCI KMS key ID to use for encryption and decryption
	EnvOCIKMSSealKeyID = "VAULT_OCIKMS_SEAL_KEY_ID"
	// OCI KMS crypto endpoint to use for encryption and decryption
	EnvOCIKMSCryptoEndpoint = "VAULT_OCIKMS_CRYPTO_ENDPOINT"
	// OCI KMS management endpoint to manage keys
	EnvOCIKMSManagementEndpoint = "VAULT_OCIKMS_MANAGEMENT_ENDPOINT"
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

var (
	metricInit    = []string{"ocikms", "init"}
	metricEncrypt = []string{"ocikms", "encrypt"}
	metricDecrypt = []string{"ocikms", "decrypt"}

	metricInitFailed    = []string{"ocikms", "initFailed"}
	metricEncryptFailed = []string{"ocikms", "encryptFailed"}
	metricDecryptFailed = []string{"ocikms", "decryptFailed"}
)

// OCIKMSMechanism is the method used to encrypt/decrypt in auto unseal process
type OCIKMSMechanism uint32

type OCIKMSSeal struct {
	authTypeAPIKey bool   // true for user principal, false for instance principal, default is false
	keyID          string // OCI KMS keyID

	cryptoEndpoint     string // OCI KMS crypto endpoint
	managementEndpoint string // OCI KMS management endpoint

	cryptoClient     *keymanagement.KmsCryptoClient     // OCI KMS crypto client
	managementClient *keymanagement.KmsManagementClient // OCI KMS management client

	currentKeyID *atomic.Value // Current key version which is used for encryption/decryption

	logger log.Logger
}

var _ seal.Access = (*OCIKMSSeal)(nil)

// NewSeal creates a new OCIKMSSeal seal with the provided logger
func NewSeal(logger log.Logger) *OCIKMSSeal {
	k := &OCIKMSSeal{
		logger:       logger,
		currentKeyID: new(atomic.Value),
	}
	k.currentKeyID.Store("")
	return k
}

func (k *OCIKMSSeal) SetConfig(config map[string]string) (map[string]string, error) {
	defer metrics.MeasureSince(metricInit, time.Now())
	if config == nil {
		config = map[string]string{}
	}

	// Check and set KeyID
	switch {
	case os.Getenv(EnvOCIKMSSealKeyID) != "":
		k.keyID = os.Getenv(EnvOCIKMSSealKeyID)
	case config[KMSConfigKeyID] != "":
		k.keyID = config[KMSConfigKeyID]
	default:
		metrics.IncrCounter(metricInitFailed, 1)
		return nil, fmt.Errorf("'%s' not found for OCI KMS seal configuration", KMSConfigKeyID)
	}
	k.logger.Info("OCI KMS configuration", KMSConfigKeyID, k.keyID)

	// Check and set cryptoEndpoint
	switch {
	case os.Getenv(EnvOCIKMSCryptoEndpoint) != "":
		k.cryptoEndpoint = os.Getenv(EnvOCIKMSCryptoEndpoint)
	case config[KMSConfigCryptoEndpoint] != "":
		k.cryptoEndpoint = config[KMSConfigCryptoEndpoint]
	default:
		metrics.IncrCounter(metricInitFailed, 1)
		return nil, fmt.Errorf("'%s' not found for OCI KMS seal configuration", KMSConfigCryptoEndpoint)
	}
	k.logger.Info("OCI KMS configuration", KMSConfigCryptoEndpoint, k.cryptoEndpoint)

	// Check and set managementEndpoint
	switch {
	case os.Getenv(EnvOCIKMSManagementEndpoint) != "":
		k.managementEndpoint = os.Getenv(EnvOCIKMSManagementEndpoint)
	case config[KMSConfigManagementEndpoint] != "":
		k.managementEndpoint = config[KMSConfigManagementEndpoint]
	default:
		metrics.IncrCounter(metricInitFailed, 1)
		return nil, fmt.Errorf("'%s' not found for OCI KMS seal configuration", KMSConfigManagementEndpoint)
	}
	k.logger.Info("OCI KMS configuration", KMSConfigManagementEndpoint, k.managementEndpoint)

	// Check and set authTypeAPIKey
	var err error
	k.authTypeAPIKey = false
	authTypeAPIKeyStr := config[KMSConfigAuthTypeAPIKey]
	if authTypeAPIKeyStr != "" {
		k.authTypeAPIKey, err = strconv.ParseBool(authTypeAPIKeyStr)
		if err != nil {
			metrics.IncrCounter(metricInitFailed, 1)
			return nil, errwrap.Wrapf("failed parsing "+KMSConfigAuthTypeAPIKey+" parameter: {{err}}", err)
		}
	}
	if k.authTypeAPIKey {
		k.logger.Info("using OCI KMS with user principal")
	} else {
		k.logger.Info("using OCI KMS with instance principal")
	}

	// Check and set OCI KMS crypto client
	if k.cryptoClient == nil {
		kmsCryptoClient, err := k.getOCIKMSCryptoClient()
		if err != nil {
			metrics.IncrCounter(metricInitFailed, 1)
			return nil, errwrap.Wrapf("error initializing OCI KMS client: {{err}}", err)
		}
		k.cryptoClient = kmsCryptoClient
	}

	// Check and set OCI KMS management client
	if k.managementClient == nil {
		kmsManagementClient, err := k.getOCIKMSManagementClient()
		if err != nil {
			metrics.IncrCounter(metricInitFailed, 1)
			return nil, errwrap.Wrapf("error initializing OCI KMS client: {{err}}", err)
		}
		k.managementClient = kmsManagementClient
	}

	// Calling Encrypt method with empty string just to validate keyId access and store current keyVersion
	encryptedBlobInfo, err := k.Encrypt(context.Background(), []byte(""))
	if err != nil || encryptedBlobInfo == nil {
		metrics.IncrCounter(metricInitFailed, 1)
		return nil, errwrap.Wrapf("failed "+KMSConfigKeyID+" validation: {{err}}", err)
	}
	k.logger.Info("successfully validated " + KMSConfigKeyID)

	// Map that holds non-sensitive configuration info
	sealInfo := make(map[string]string)
	sealInfo[KMSConfigKeyID] = k.keyID
	sealInfo[KMSConfigCryptoEndpoint] = k.cryptoEndpoint
	sealInfo[KMSConfigManagementEndpoint] = k.managementEndpoint

	return sealInfo, nil
}

func (k *OCIKMSSeal) SealType() string {
	return seal.OCIKMS
}

func (k *OCIKMSSeal) KeyID() string {
	return k.currentKeyID.Load().(string)
}

func (k *OCIKMSSeal) Init(context.Context) error {
	return nil
}

func (k *OCIKMSSeal) Finalize(context.Context) error {
	return nil
}

func (k *OCIKMSSeal) Encrypt(ctx context.Context, plaintext []byte) (*physical.EncryptedBlobInfo, error) {
	defer metrics.MeasureSince(metricEncrypt, time.Now())
	if plaintext == nil {
		metrics.IncrCounter(metricEncryptFailed, 1)
		return nil, fmt.Errorf("given plaintext for encryption is nil")
	}

	env, err := seal.NewEnvelope().Encrypt(plaintext)
	if err != nil {
		metrics.IncrCounter(metricEncryptFailed, 1)
		return nil, errwrap.Wrapf("error wrapping data: {{err}}", err)
	}

	if k.cryptoClient == nil {
		metrics.IncrCounter(metricEncryptFailed, 1)
		return nil, fmt.Errorf("nil client")
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
		metrics.IncrCounter(metricEncryptFailed, 1)
		return nil, errwrap.Wrapf("error encrypting data: {{err}}", err)
	}

	// Note: It is potential a timing issue if the key gets rotated between this
	// getCurrentKeyVersion operation and above Encrypt operation
	keyVersion, err := k.getCurrentKeyVersion()
	if err != nil {
		metrics.IncrCounter(metricEncryptFailed, 1)
		return nil, errwrap.Wrapf("error getting current key version: {{err}}", err)
	}
	// Update key version
	k.currentKeyID.Store(keyVersion)

	ret := &physical.EncryptedBlobInfo{
		Ciphertext: env.Ciphertext,
		IV:         env.IV,
		KeyInfo: &physical.SealKeyInfo{
			// Storing current key version in case we want to re-wrap older entries
			KeyID:      keyVersion,
			WrappedKey: []byte(*output.Ciphertext),
		},
	}

	k.logger.Debug("successfully encrypted")
	return ret, nil
}

func (k *OCIKMSSeal) Decrypt(ctx context.Context, in *physical.EncryptedBlobInfo) ([]byte, error) {
	defer metrics.MeasureSince(metricDecrypt, time.Now())
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
		metrics.IncrCounter(metricDecryptFailed, 1)
		return nil, errwrap.Wrapf("error decrypting data: {{err}}", err)
	}
	envelopKey, err := base64.StdEncoding.DecodeString(*output.Plaintext)
	if err != nil {
		metrics.IncrCounter(metricDecryptFailed, 1)
		return nil, errwrap.Wrapf("error base64 decrypting data: {{err}}", err)
	}
	envInfo := &seal.EnvelopeInfo{
		Key:        envelopKey,
		IV:         in.IV,
		Ciphertext: in.Ciphertext,
	}

	plaintext, err := seal.NewEnvelope().Decrypt(envInfo)
	if err != nil {
		return nil, errwrap.Wrapf("error decrypting data: {{err}}", err)
	}

	return plaintext, nil
}

func (k *OCIKMSSeal) getConfigProvider() (common.ConfigurationProvider, error) {
	var cp common.ConfigurationProvider
	var err error
	if k.authTypeAPIKey {
		cp = common.DefaultConfigProvider()
	} else {
		cp, err = auth.InstancePrincipalConfigurationProvider()
		if err != nil {
			return nil, errwrap.Wrapf("failed creating InstancePrincipalConfigurationProvider: {{err}}", err)
		}
	}
	return cp, nil
}

// Build OCI KMS crypto client
func (k *OCIKMSSeal) getOCIKMSCryptoClient() (*keymanagement.KmsCryptoClient, error) {
	cp, err := k.getConfigProvider()
	if err != nil {
		return nil, errwrap.Wrapf("failed creating configuration provider: {{err}}", err)
	}

	// Build crypto client
	kmsCryptoClient, err := keymanagement.NewKmsCryptoClientWithConfigurationProvider(cp, k.cryptoEndpoint)
	if err != nil {
		return nil, errwrap.Wrapf("failed creating NewKmsCryptoClientWithConfigurationProvider: {{err}}", err)
	}

	return &kmsCryptoClient, nil
}

// Build OCI KMS management client
func (k *OCIKMSSeal) getOCIKMSManagementClient() (*keymanagement.KmsManagementClient, error) {
	cp, err := k.getConfigProvider()
	if err != nil {
		return nil, errwrap.Wrapf("failed creating configuration provider: {{err}}", err)
	}

	// Build crypto client
	kmsManagementClient, err := keymanagement.NewKmsManagementClientWithConfigurationProvider(cp, k.managementEndpoint)
	if err != nil {
		return nil, errwrap.Wrapf("failed creating NewKmsCryptoClientWithConfigurationProvider: {{err}}", err)
	}

	return &kmsManagementClient, nil
}

// Request metadata includes retry policy
func (k *OCIKMSSeal) getRequestMetadata() common.RequestMetadata {
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

func (k *OCIKMSSeal) getCurrentKeyVersion() (string, error) {
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
		return "", errwrap.Wrapf("failed getting current key version: {{err}}", err)
	}

	return *getKeyResponse.CurrentKeyVersion, nil
}
