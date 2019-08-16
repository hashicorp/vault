// Copyright © 2019, Oracle and/or its affiliates.
package ocikms

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/seal"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/common/auth"
	"github.com/oracle/oci-go-sdk/keymanagement"
	"math"
	"os"
	"strconv"
	"time"
)

const (
	// OCI KMS key ID to use for encryption and decryption
	EnvOCIKMSSealKeyID = "VAULT_OCIKMS_SEAL_KEY_ID"
	// OCI KMS key ID to use for encryption and decryption
	EnvOCIKMSCryptoEndpoint = "VAULT_OCIKMS_CRYPTO_ENDPOINT"
	// Maximum number of retries
	KMSMaximumNumberOfRetries = 5
	// keyID config
	KMSConfigKeyID = "keyID"
	// cryptoEndpoint config
	KMSConfigCryptoEndpoint = "cryptoEndpoint"
	// authTypeAPIKey config
	KMSConfigAuthTypeAPIKey = "authTypeAPIKey"
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
	authTypeAPIKey bool                           // true for user principal, false for instance principal, default is false
	keyID          string                         // OCI KMS keyID
	cryptoEndpoint string                         // OCI KMS crypto endpoint
	cryptoClient   *keymanagement.KmsCryptoClient // OCI KMS crypto client
	logger         log.Logger
}

var _ seal.Access = (*OCIKMSSeal)(nil)

// NewSeal creates a new OCIKMSSeal seal with the provided logger
func NewSeal(logger log.Logger) *OCIKMSSeal {
	k := &OCIKMSSeal{
		logger: logger,
		keyID:  "",
	}
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

	// Check and set authTypeAPIKey
	var err error
	k.authTypeAPIKey = false
	authTypeAPIKeyStr := config[KMSConfigAuthTypeAPIKey]
	if authTypeAPIKeyStr != "" {
		k.authTypeAPIKey, err = strconv.ParseBool(authTypeAPIKeyStr)
		if err != nil {
			metrics.IncrCounter(metricInitFailed, 1)
			return nil, errwrap.Wrapf("failed parsing authTypeAPIKey parameter: {{err}}", err)
		}
	}
	if k.authTypeAPIKey {
		k.logger.Info("using OCI KMS with user principal")
	} else {
		k.logger.Info("using OCI KMS with instance principal")
	}

	// Check and set OCI KMS crypto client
	if k.cryptoClient == nil {
		kmsCryptoClient, err := k.getOCIKMSClient()
		if err != nil {
			metrics.IncrCounter(metricInitFailed, 1)
			return nil, errwrap.Wrapf("error initializing OCI KMS clients: {{err}}", err)
		}

		// KeyID validation by trying to generate a DEK
		keyLength := 32
		includePlaintextKey := true
		generateKeyDetails := keymanagement.GenerateKeyDetails{
			IncludePlaintextKey: &includePlaintextKey,
			KeyId:               &k.keyID,
			KeyShape: &keymanagement.KeyShape{
				Algorithm: "AES",
				Length:    &keyLength,
			},
		}
		requestMetadata := k.getRequestMetadata()
		dekInput := keymanagement.GenerateDataEncryptionKeyRequest{
			GenerateKeyDetails: generateKeyDetails,
			RequestMetadata:    requestMetadata,
		}
		generateDEKResponse, err := kmsCryptoClient.GenerateDataEncryptionKey(context.Background(), dekInput)
		if err != nil || generateDEKResponse.Ciphertext == nil {
			metrics.IncrCounter(metricInitFailed, 1)
			return nil, errwrap.Wrapf("failed keyID validation: {{err}}", err)
		}
		k.logger.Info("successfully validated keyID")

		// Store client
		k.cryptoClient = kmsCryptoClient
	}

	// Map that holds non-sensitive configuration info
	sealInfo := make(map[string]string)
	sealInfo[KMSConfigKeyID] = k.keyID
	sealInfo[KMSConfigCryptoEndpoint] = k.cryptoEndpoint

	return sealInfo, nil
}

// Build OCI KMS crypto client
func (k *OCIKMSSeal) getOCIKMSClient() (*keymanagement.KmsCryptoClient, error) {
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

	// Build crypto client
	kmsCryptoClient, err := keymanagement.NewKmsCryptoClientWithConfigurationProvider(cp, k.cryptoEndpoint)
	if err != nil {
		return nil, errwrap.Wrapf("failed creating NewKmsCryptoClientWithConfigurationProvider: {{err}}", err)
	}

	return &kmsCryptoClient, nil
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

func (k *OCIKMSSeal) SealType() string {
	return seal.OCIKMS
}

func (k *OCIKMSSeal) KeyID() string {
	return k.keyID
}

func (k *OCIKMSSeal) Init(context.Context) error {
	return nil
}

func (k *OCIKMSSeal) Finalize(context.Context) error {
	return nil
}

func (k *OCIKMSSeal) Encrypt(ctx context.Context, plaintext []byte) (*physical.EncryptedBlobInfo, error) {
	defer metrics.MeasureSince(metricEncrypt, time.Now())
	if plaintext == nil || len(plaintext) == 0 {
		metrics.IncrCounter(metricEncryptFailed, 1)
		return nil, fmt.Errorf("given plaintext for encryption is nil")
	}

	// KMS required base64 encrypted plain text before sending to the service
	encodedPlaintext := base64.StdEncoding.EncodeToString(plaintext)

	if k.cryptoClient == nil {
		metrics.IncrCounter(metricEncryptFailed, 1)
		return nil, fmt.Errorf("nil client")
	}

	// Build Encrypt Request
	requestMetadata := k.getRequestMetadata()
	encryptedDataDetails := keymanagement.EncryptDataDetails{
		KeyId:     &k.keyID,
		Plaintext: &encodedPlaintext,
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
	k.logger.Debug("successfully encrypted")

	ret := &physical.EncryptedBlobInfo{
		Ciphertext: []byte(*output.Ciphertext),
		KeyInfo: &physical.SealKeyInfo{
			KeyID: k.keyID,
		},
	}

	return ret, nil
}

func (k *OCIKMSSeal) Decrypt(ctx context.Context, in *physical.EncryptedBlobInfo) ([]byte, error) {
	defer metrics.MeasureSince(metricDecrypt, time.Now())
	if in == nil {
		return nil, fmt.Errorf("given input for decryption is nil")
	}

	requestMetadata := k.getRequestMetadata()

	cipherText := string(in.Ciphertext)
	decryptedDataDetails := keymanagement.DecryptDataDetails{
		KeyId:      &k.keyID,
		Ciphertext: &cipherText,
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
	plaintext, err := base64.StdEncoding.DecodeString(*output.Plaintext)
	if err != nil {
		metrics.IncrCounter(metricDecryptFailed, 1)
		return nil, errwrap.Wrapf("error base64 decrypting data: {{err}}", err)
	}
	k.logger.Debug("successfully decrypted")

	return plaintext, nil
}
