package tencentcloudkms

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"sync/atomic"

	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
)

// These constants are TencentCloud accepted env vars
const (
	PROVIDER_SECRET_ID      = "TENCENTCLOUD_SECRET_ID"
	PROVIDER_SECRET_KEY     = "TENCENTCLOUD_SECRET_KEY"
	PROVIDER_SECURITY_TOKEN = "TENCENTCLOUD_SECURITY_TOKEN"
	PROVIDER_REGION         = "TENCENTCLOUD_REGION"
	PROVIDER_KMS_KEY_ID     = "TENCENTCLOUD_KMS_KEY_ID"
)

// Wrapper is a wrapper that uses TencentCloud KMS
type Wrapper struct {
	accessKey    string
	secretKey    string
	sessionToken string
	region       string

	keyID        string
	currentKeyID *atomic.Value

	client kmsClient
}

// Ensure that we are implementing Wrapper
var _ wrapping.Wrapper = (*Wrapper)(nil)

// NewWrapper returns a new TencentCloud wrapper
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

// SetConfig sets the fields on the wrapper object based on TencentCloud config parameter
//
// Order of precedence values:
// * Environment variable
// * Instance metadata role
func (k *Wrapper) SetConfig(config map[string]string) (map[string]string, error) {
	if config == nil {
		config = map[string]string{}
	}

	switch {
	case os.Getenv(PROVIDER_KMS_KEY_ID) != "":
		k.keyID = os.Getenv(PROVIDER_KMS_KEY_ID)
	case config["kms_key_id"] != "":
		k.keyID = config["kms_key_id"]
	default:
		return nil, fmt.Errorf("'kms_key_id' not found for TencentCloud KMS wrapper configuration")
	}

	switch {
	case os.Getenv(PROVIDER_REGION) != "":
		k.region = os.Getenv(PROVIDER_REGION)
	case config["region"] != "":
		k.region = config["region"]
	default:
		k.region = "ap-guangzhou"
	}

	switch {
	case os.Getenv(PROVIDER_SECRET_ID) != "":
		k.accessKey = os.Getenv(PROVIDER_SECRET_ID)
	case config["access_key"] != "":
		k.accessKey = config["access_key"]
	default:
		return nil, fmt.Errorf("'access_key' not found for TencentCloud KMS wrapper configuration")
	}

	switch {
	case os.Getenv(PROVIDER_SECRET_KEY) != "":
		k.secretKey = os.Getenv(PROVIDER_SECRET_KEY)
	case config["secret_key"] != "":
		k.secretKey = config["secret_key"]
	default:
		return nil, fmt.Errorf("'secret_key' not found for TencentCloud KMS wrapper configuration")
	}

	switch {
	case os.Getenv(PROVIDER_SECURITY_TOKEN) != "":
		k.sessionToken = os.Getenv(PROVIDER_SECURITY_TOKEN)
	case config["session_token"] != "":
		k.sessionToken = config["session_token"]
	default:
		k.sessionToken = ""
	}

	if k.client == nil {
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.ReqMethod = "POST"
		cpf.HttpProfile.ReqTimeout = 300
		cpf.Language = "en-US"

		credential := common.NewTokenCredential(k.accessKey, k.secretKey, k.sessionToken)
		client, err := kms.NewClient(credential, k.region, cpf)
		if err != nil {
			return nil, fmt.Errorf("error initializing TencentCloud KMS client: %w", err)
		}

		input := kms.NewDescribeKeyRequest()
		input.KeyId = &k.keyID
		keyInfo, err := client.DescribeKey(input)
		if err != nil {
			return nil, fmt.Errorf("error fetching TencentCloud KMS information: %w", err)
		}

		if keyInfo.Response.KeyMetadata == nil || keyInfo.Response.KeyMetadata.KeyId == nil {
			return nil, fmt.Errorf("no key information return")
		}

		k.currentKeyID.Store(*keyInfo.Response.KeyMetadata.KeyId)
		k.client = client
	}

	wrappingInfo := make(map[string]string)
	wrappingInfo["region"] = k.region
	wrappingInfo["kms_key_id"] = k.keyID

	return wrappingInfo, nil
}

// Init is called during core.Initialize. No-op at the moment.
func (k *Wrapper) Init(_ context.Context) error {
	return nil
}

// Finalize is called during shutdown. No-op at the moment.
func (k *Wrapper) Finalize(_ context.Context) error {
	return nil
}

// Type returns the type for this particular wrapper implementation
func (k *Wrapper) Type() string {
	return wrapping.TencentCloudKMS
}

// KeyID returns the last known key id
func (k *Wrapper) KeyID() string {
	return k.currentKeyID.Load().(string)
}

// HMACKeyID returns nothing, it's here to satisfy the interface
func (k *Wrapper) HMACKeyID() string {
	return ""
}

// Encrypt is used to encrypt the master key using the the TencentCloud KMS.
// This returns the ciphertext, and/or any errors from this call.
// This should be called after the KMS client has been instantiated.
func (k *Wrapper) Encrypt(_ context.Context, plaintext, aad []byte) (blob *wrapping.EncryptedBlobInfo, err error) {
	if plaintext == nil {
		return nil, fmt.Errorf("given plaintext for encryption is nil")
	}

	env, err := wrapping.NewEnvelope(nil).Encrypt(plaintext, aad)
	if err != nil {
		return nil, fmt.Errorf("error wrapping data: %w", err)
	}

	input := kms.NewEncryptRequest()
	input.KeyId = &k.keyID
	input.Plaintext = common.StringPtr(base64.StdEncoding.EncodeToString(env.Key))

	output, err := k.client.Encrypt(input)
	if err != nil {
		return nil, fmt.Errorf("error encrypting data: %w", err)
	}

	keyID := *output.Response.KeyId
	k.currentKeyID.Store(keyID)

	ret := &wrapping.EncryptedBlobInfo{
		Ciphertext: env.Ciphertext,
		IV:         env.IV,
		KeyInfo: &wrapping.KeyInfo{
			KeyID:      keyID,
			WrappedKey: []byte(*output.Response.CiphertextBlob),
		},
	}

	return ret, nil
}

// Decrypt is used to decrypt the ciphertext using the the TencentCloud KMS.
// This should be called after the KMS client has been instantiated.
func (k *Wrapper) Decrypt(_ context.Context, in *wrapping.EncryptedBlobInfo, aad []byte) (pt []byte, err error) {
	if in == nil {
		return nil, fmt.Errorf("given input for decryption is nil")
	}

	input := kms.NewDecryptRequest()
	input.CiphertextBlob = common.StringPtr(string(in.KeyInfo.WrappedKey))

	output, err := k.client.Decrypt(input)
	if err != nil {
		return nil, fmt.Errorf("error decrypting data encryption key: %w", err)
	}

	keyBytes, err := base64.StdEncoding.DecodeString(*output.Response.Plaintext)
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
