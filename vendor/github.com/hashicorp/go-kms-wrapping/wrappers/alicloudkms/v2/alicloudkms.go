// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
)

// These constants contain the accepted env vars; the Vault one is for backwards compat
const (
	EnvAliCloudKmsWrapperKeyId   = "ALICLOUDKMS_WRAPPER_KEY_ID"
	EnvVaultAliCloudKmsSealKeyId = "VAULT_ALICLOUDKMS_SEAL_KEY_ID"
)

// Wrapper is a Wrapper that uses AliCloud's KMS
type Wrapper struct {
	client       kmsClient
	domain       string
	keyId        string
	currentKeyId *atomic.Value
}

// Ensure that we are implementing Wrapper
var _ wrapping.Wrapper = (*Wrapper)(nil)

// NewWrapper creates a new AliCloud Wrapper
func NewWrapper() *Wrapper {
	k := &Wrapper{
		currentKeyId: new(atomic.Value),
	}
	k.currentKeyId.Store("")
	return k
}

// SetConfig sets the fields on the AliCloudKMSWrapper object based on
// values from the config parameter.
//
// Order of precedence AliCloud values:
// * Environment variable
// * Value from Vault configuration file
// * Instance metadata role (access key and secret key)
func (k *Wrapper) SetConfig(_ context.Context, opt ...wrapping.Option) (*wrapping.WrapperConfig, error) {
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, err
	}

	// Check and set KeyId
	switch {
	case os.Getenv(EnvAliCloudKmsWrapperKeyId) != "" && !opts.Options.WithDisallowEnvVars:
		k.keyId = os.Getenv(EnvAliCloudKmsWrapperKeyId)
	case os.Getenv(EnvVaultAliCloudKmsSealKeyId) != "" && !opts.Options.WithDisallowEnvVars:
		k.keyId = os.Getenv(EnvVaultAliCloudKmsSealKeyId)
	case opts.WithKeyId != "":
		k.keyId = opts.WithKeyId
	default:
		return nil, fmt.Errorf("key id not found (env or config) for alicloud kms wrapper configuration")
	}

	region := ""
	if k.client == nil {
		// Check and set region.
		if !opts.Options.WithDisallowEnvVars {
			region = os.Getenv("ALICLOUD_REGION")
		}
		if region == "" {
			region = opts.withRegion
		}

		// A domain isn't required, but it can be used to override the endpoint
		// returned by the region. An example value for a domain would be:
		// "kms.us-east-1.aliyuncs.com".
		if !opts.Options.WithDisallowEnvVars {
			k.domain = os.Getenv("ALICLOUD_DOMAIN")
		}
		if k.domain == "" {
			k.domain = opts.withDomain
		}

		// Build the optional, configuration-based piece of the credential chain.
		credConfig := &providers.Configuration{
			AccessKeyID:     opts.withAccessKey,
			AccessKeySecret: opts.withSecretKey,
		}
		if credConfig.AccessKeySecret == "" {
			credConfig.AccessKeySecret = opts.withAccessSecret
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
	input.KeyId = k.keyId
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
	k.currentKeyId.Store(keyInfo.KeyMetadata.KeyId)

	// Map that holds non-sensitive configuration info
	wrapConfig := new(wrapping.WrapperConfig)
	wrapConfig.Metadata = make(map[string]string)
	wrapConfig.Metadata["region"] = region
	wrapConfig.Metadata["kms_key_id"] = k.keyId
	if k.domain != "" {
		wrapConfig.Metadata["domain"] = k.domain
	}

	return wrapConfig, nil
}

// Type returns the type for this particular wrapper implementation
func (k *Wrapper) Type(_ context.Context) (wrapping.WrapperType, error) {
	return wrapping.WrapperTypeAliCloudKms, nil
}

// KeyId returns the last known key id
func (k *Wrapper) KeyId(_ context.Context) (string, error) {
	return k.currentKeyId.Load().(string), nil
}

// Encrypt is used to encrypt the master key using the the AliCloud CMK.
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

	input := kms.CreateEncryptRequest()
	input.KeyId = k.keyId
	input.Plaintext = base64.StdEncoding.EncodeToString(env.Key)
	input.Domain = k.domain

	output, err := k.client.Encrypt(input)
	if err != nil {
		return nil, fmt.Errorf("error encrypting data: %w", err)
	}

	// Store the current key id.
	keyId := output.KeyId
	k.currentKeyId.Store(keyId)

	ret := &wrapping.BlobInfo{
		Ciphertext: env.Ciphertext,
		Iv:         env.Iv,
		KeyInfo: &wrapping.KeyInfo{
			KeyId:      keyId,
			WrappedKey: []byte(output.CiphertextBlob),
		},
	}

	return ret, nil
}

// Decrypt is used to decrypt the ciphertext. This should be called after Init.
func (k *Wrapper) Decrypt(_ context.Context, in *wrapping.BlobInfo, opt ...wrapping.Option) ([]byte, error) {
	if in == nil {
		return nil, fmt.Errorf("given input for decryption is nil")
	}

	// KeyId is not passed to this call because AliCloud handles this
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
		Iv:         in.Iv,
		Ciphertext: in.Ciphertext,
	}
	plaintext, err := wrapping.EnvelopeDecrypt(envInfo, opt...)
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
