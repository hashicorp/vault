// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials/providers"
)

// chainedCreds allows us to grab AliCloud credentials in the following order:
//
//  1. Credentials set via env variables if available
//  2. Credentials explicitly configured if available
//  3. Credentials inferred from instance metadata if available
//
// Since credentials are pulled every time a client is created, and clients are not
// long-lasting, this means that credentials could be changed or overridden at any
// point and they would be picked up without requiring a Vault restart.
func chainedCreds(key, secret string) (auth.Credential, error) {
	providerChain := []providers.Provider{
		providers.NewEnvCredentialProvider(),
		providers.NewConfigurationCredentialProvider(&providers.Configuration{
			AccessKeyID:     key,
			AccessKeySecret: secret,
		}),
		providers.NewInstanceMetadataProvider(),
	}
	return providers.NewChainProvider(providerChain).Retrieve()
}
