// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"github.com/hashicorp/go-kms-wrapping/wrappers/alicloudkms/v2"
	"github.com/hashicorp/go-kms-wrapping/wrappers/awskms/v2"
	"github.com/hashicorp/go-kms-wrapping/wrappers/azurekeyvault/v2"
	"github.com/hashicorp/go-kms-wrapping/wrappers/gcpckms/v2"
	"github.com/hashicorp/go-kms-wrapping/wrappers/ocikms/v2"
	"github.com/hashicorp/go-kms-wrapping/wrappers/transit/v2"
)

var (
	AliCloudKMSEnvVars = map[string]string{
		"ALICLOUD_REGION":                        "region",
		"ALICLOUD_DOMAIN":                        "domain",
		"ALICLOUD_ACCESS_KEY":                    "access_key",
		"ALICLOUD_SECRET_KEY":                    "secret_key",
		alicloudkms.EnvVaultAliCloudKmsSealKeyId: "kms_key_id",
		alicloudkms.EnvAliCloudKmsWrapperKeyId:   "kms_key_id",
	}

	AWSKMSEnvVars = map[string]string{
		"AWS_REGION":                   "region",
		"AWS_DEFAULT_REGION":           "region",
		"AWS_ACCESS_KEY_ID":            "access_key",
		"AWS_SESSION_TOKEN":            "session_token",
		"AWS_SECRET_ACCESS_KEY":        "secret_key",
		awskms.EnvVaultAwsKmsSealKeyId: "kms_key_id",
		awskms.EnvAwsKmsWrapperKeyId:   "kms_key_id",
		"AWS_KMS_ENDPOINT":             "endpoint",
	}

	AzureEnvVars = map[string]string{
		"AZURE_TENANT_ID":                              "tenant_id",
		"AZURE_CLIENT_ID":                              "client_id",
		"AZURE_CLIENT_SECRET":                          "client_secret",
		"AZURE_ENVIRONMENT":                            "environment",
		"AZURE_AD_RESOURCE":                            "resource",
		azurekeyvault.EnvAzureKeyVaultWrapperKeyName:   "key_name",
		azurekeyvault.EnvVaultAzureKeyVaultKeyName:     "key_name",
		azurekeyvault.EnvAzureKeyVaultWrapperVaultName: "vault_name",
		azurekeyvault.EnvVaultAzureKeyVaultVaultName:   "vault_name",
	}

	GCPCKMSEnvVars = map[string]string{
		gcpckms.EnvGcpCkmsWrapperCredsPath:   "credentials",
		"GOOGLE_APPLICATION_CREDENTIALS":     "credentials",
		gcpckms.EnvGcpCkmsWrapperProject:     "project",
		gcpckms.EnvGcpCkmsWrapperLocation:    "region",
		gcpckms.EnvVaultGcpCkmsSealCryptoKey: "crypto_key",
		gcpckms.EnvGcpCkmsWrapperCryptoKey:   "crypto_key",
		gcpckms.EnvGcpCkmsWrapperKeyRing:     "key_ring",
		gcpckms.EnvVaultGcpCkmsSealKeyRing:   "key_ring",
	}

	OCIKMSEnvVars = map[string]string{
		ocikms.EnvOciKmsWrapperCryptoEndpoint:       "crypto_endpoint",
		ocikms.EnvVaultOciKmsSealCryptoEndpoint:     "crypto_endpoint",
		ocikms.EnvOciKmsWrapperKeyId:                "key_id",
		ocikms.EnvVaultOciKmsSealKeyId:              "key_id",
		ocikms.EnvOciKmsWrapperManagementEndpoint:   "management_endpoint",
		ocikms.EnvVaultOciKmsSealManagementEndpoint: "management_endpoint",
	}

	TransitEnvVars = map[string]string{
		"VAULT_ADDR":                              "address",
		"VAULT_TOKEN":                             "token",
		"VAULT_NAMESPACE":                         "namespace",
		"VAULT_CACERT":                            "tls_ca_cert",
		"VAULT_CLIENT_CERT":                       "tls_client_cert",
		"VAULT_CLIENT_KEY":                        "tls_client_key",
		"VAULT_TLS_SERVER_NAME":                   "tls_server_name",
		"VAULT_SKIP_VERIFY":                       "tls_skip_verify",
		transit.EnvVaultTransitSealKeyName:        "key_name",
		transit.EnvTransitWrapperKeyName:          "key_name",
		transit.EnvTransitWrapperMountPath:        "mount_path",
		transit.EnvVaultTransitSealMountPath:      "mount_path",
		transit.EnvTransitWrapperDisableRenewal:   "disable_renewal",
		transit.EnvVaultTransitSealDisableRenewal: "disable_renewal",
	}

	// TransitPrioritizeConfigValues are the variables where file config takes precedence over env vars in transit seals
	TransitPrioritizeConfigValues = []string{
		"token",
		"address",
	}

	// TransitTLSConfigVars are the TLS config variables for transit seals
	// if one of them is set in file config, transit seals use the file config for all TLS values and ignore env vars
	// otherwise they use the env vars for TLS config
	TransitTLSConfigVars = []string{
		"tls_ca_cert",
		"tls_client_cert",
		"tls_client_key",
		"tls_server_name",
		"tls_skip_verify",
	}
)
