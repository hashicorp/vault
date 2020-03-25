// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// AwsKmsProviderOptions specifies options for configuring the AWS KMS provider.
type AwsKmsProviderOptions struct {
	AccessKeyID     string
	SecretAccessKey string
}

// AwsKmsProvider creates a new AwsKmsProviderOptions instance.
func AwsKmsProvider() *AwsKmsProviderOptions {
	return &AwsKmsProviderOptions{}
}

// SetAccessKeyID specifies the AWS access key ID.
func (akpo *AwsKmsProviderOptions) SetAccessKeyID(accessKeyID string) *AwsKmsProviderOptions {
	akpo.AccessKeyID = accessKeyID
	return akpo
}

// SetSecretAccessKey specifies the AWS secret access key.
func (akpo *AwsKmsProviderOptions) SetSecretAccessKey(secretAccessKey string) *AwsKmsProviderOptions {
	akpo.SecretAccessKey = secretAccessKey
	return akpo
}

// LocalKmsProviderOptions specifies options for configuring a local KMS provider.
type LocalKmsProviderOptions struct {
	MasterKey []byte
}

// LocalKmsProvider creates a new LocalKmsProviderOptions instance.
func LocalKmsProvider() *LocalKmsProviderOptions {
	return &LocalKmsProviderOptions{}
}

// SetMasterKey specifies the local master key.
func (lkpo *LocalKmsProviderOptions) SetMasterKey(key []byte) *LocalKmsProviderOptions {
	lkpo.MasterKey = key
	return lkpo
}
