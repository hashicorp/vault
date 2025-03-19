// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package configutil

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-kms-wrapping/wrappers/ocikms/v2"
	"github.com/stretchr/testify/require"
)

func Test_getEnvConfig(t *testing.T) {
	tests := []struct {
		name    string
		kms     *KMS
		envVars map[string]string
		want    map[string]string
	}{
		{
			"AliCloud wrapper",
			&KMS{
				Type:     "alicloudkms",
				Priority: 1,
			},
			map[string]string{"ALICLOUD_REGION": "test_region", "ALICLOUD_DOMAIN": "test_domain", "ALICLOUD_ACCESS_KEY": "test_access_key", "ALICLOUD_SECRET_KEY": "test_secret_key", "VAULT_ALICLOUDKMS_SEAL_KEY_ID": "test_key_id"},
			map[string]string{"region": "test_region", "domain": "test_domain", "access_key": "test_access_key", "secret_key": "test_secret_key", "kms_key_id": "test_key_id"},
		},
		{
			"AWS KMS wrapper",
			&KMS{
				Type:     "awskms",
				Priority: 1,
			},
			map[string]string{"AWS_REGION": "test_region", "AWS_ACCESS_KEY_ID": "test_access_key", "AWS_SECRET_ACCESS_KEY": "test_secret_key", "VAULT_AWSKMS_SEAL_KEY_ID": "test_key_id"},
			map[string]string{"region": "test_region", "access_key": "test_access_key", "secret_key": "test_secret_key", "kms_key_id": "test_key_id"},
		},
		{
			"Azure KeyVault wrapper",
			&KMS{
				Type:     "azurekeyvault",
				Priority: 1,
			},
			map[string]string{"AZURE_TENANT_ID": "test_tenant_id", "AZURE_CLIENT_ID": "test_client_id", "AZURE_CLIENT_SECRET": "test_client_secret", "AZURE_ENVIRONMENT": "test_environment", "VAULT_AZUREKEYVAULT_VAULT_NAME": "test_vault_name", "VAULT_AZUREKEYVAULT_KEY_NAME": "test_key_name"},
			map[string]string{"tenant_id": "test_tenant_id", "client_id": "test_client_id", "client_secret": "test_client_secret", "environment": "test_environment", "vault_name": "test_vault_name", "key_name": "test_key_name"},
		},
		{
			"GCP CKMS wrapper",
			&KMS{
				Type:     "gcpckms",
				Priority: 1,
			},
			map[string]string{"GOOGLE_CREDENTIALS": "test_credentials", "GOOGLE_PROJECT": "test_project", "GOOGLE_REGION": "test_region", "VAULT_GCPCKMS_SEAL_KEY_RING": "test_key_ring", "VAULT_GCPCKMS_SEAL_CRYPTO_KEY": "test_crypto_key"},
			map[string]string{"credentials": "test_credentials", "project": "test_project", "region": "test_region", "key_ring": "test_key_ring", "crypto_key": "test_crypto_key"},
		},
		{
			"OCI KMS wrapper",
			&KMS{
				Type:     "ocikms",
				Priority: 1,
			},
			map[string]string{"VAULT_OCIKMS_SEAL_KEY_ID": "test_key_id", "VAULT_OCIKMS_CRYPTO_ENDPOINT": "test_crypto_endpoint", "VAULT_OCIKMS_MANAGEMENT_ENDPOINT": "test_management_endpoint"},
			map[string]string{"key_id": "test_key_id", "crypto_endpoint": "test_crypto_endpoint", "management_endpoint": "test_management_endpoint"},
		},
		{
			"Transit wrapper",
			&KMS{
				Type:     "transit",
				Priority: 1,
			},
			map[string]string{"VAULT_ADDR": "test_address", "VAULT_TOKEN": "test_token", "VAULT_TRANSIT_SEAL_KEY_NAME": "test_key_name", "VAULT_TRANSIT_SEAL_MOUNT_PATH": "test_mount_path"},
			map[string]string{"address": "test_address", "token": "test_token", "key_name": "test_key_name", "mount_path": "test_mount_path"},
		},
		{
			"Environment vars not set",
			&KMS{
				Type:     "awskms",
				Priority: 1,
			},
			map[string]string{},
			map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for envName, envVal := range tt.envVars {
				t.Setenv(envName, envVal)
			}

			if got := GetEnvConfigFunc(tt.kms); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEnvConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestParseKMSesURLConformance tests that all config attrs whose values can be
// URLs, IP addresses, or host:port addresses, when configured with an IPv6
// address, the normalized to be conformant with RFC-5942 §4
// See: https://rfc-editor.org/rfc/rfc5952.html
func TestParseKMSesURLConformance(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		config     string
		expected   map[string]string
		shouldFail bool
	}{
		"alicloudkms ipv4": {
			config: `
seal "alicloudkms" {
  region     = "us-east-1"
	domain     = "kms.us-east-1.aliyuncs.com"
  access_key = "0wNEpMMlzy7szvai"
  secret_key = "PupkTg8jdmau1cXxYacgE736PJj4cA"
  kms_key_id = "08c33a6f-4e0a-4a1b-a3fa-7ddfa1d4fb73"
}`,
			expected: map[string]string{
				"region":     "us-east-1",
				"domain":     "kms.us-east-1.aliyuncs.com",
				"access_key": "0wNEpMMlzy7szvai",
				"secret_key": "PupkTg8jdmau1cXxYacgE736PJj4cA",
				"kms_key_id": "08c33a6f-4e0a-4a1b-a3fa-7ddfa1d4fb73",
			},
		},
		"alicloudkms ipv6": {
			config: `
seal "alicloudkms" {
  region     = "us-east-1"
	domain     = "2001:db8:0:0:0:0:2:1"
  access_key = "0wNEpMMlzy7szvai"
  secret_key = "PupkTg8jdmau1cXxYacgE736PJj4cA"
  kms_key_id = "08c33a6f-4e0a-4a1b-a3fa-7ddfa1d4fb73"
}`,
			expected: map[string]string{
				"region":     "us-east-1",
				"domain":     "2001:db8::2:1",
				"access_key": "0wNEpMMlzy7szvai",
				"secret_key": "PupkTg8jdmau1cXxYacgE736PJj4cA",
				"kms_key_id": "08c33a6f-4e0a-4a1b-a3fa-7ddfa1d4fb73",
			},
		},
		"awskms ipv4": {
			config: `
seal "awskms" {
  region     = "us-east-1"
  access_key = "AKIAIOSFODNN7EXAMPLE"
  secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  kms_key_id = "19ec80b0-dfdd-4d97-8164-c6examplekey"
	endpoint   = "https://vpce-0e1bb1852241f8cc6-pzi0do8n.kms.us-east-1.vpce.amazonaws.com"
}`,
			expected: map[string]string{
				"region":     "us-east-1",
				"access_key": "AKIAIOSFODNN7EXAMPLE",
				"secret_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				"kms_key_id": "19ec80b0-dfdd-4d97-8164-c6examplekey",
				"endpoint":   "https://vpce-0e1bb1852241f8cc6-pzi0do8n.kms.us-east-1.vpce.amazonaws.com",
			},
		},
		"awskms ipv6": {
			config: `
seal "awskms" {
  region     = "us-east-1"
  access_key = "AKIAIOSFODNN7EXAMPLE"
  secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  kms_key_id = "19ec80b0-dfdd-4d97-8164-c6examplekey"
  endpoint   = "https://[2001:db8:0:0:0:0:2:1]:5984/my-aws-endpoint"
}`,
			expected: map[string]string{
				"region":     "us-east-1",
				"access_key": "AKIAIOSFODNN7EXAMPLE",
				"secret_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				"kms_key_id": "19ec80b0-dfdd-4d97-8164-c6examplekey",
				"endpoint":   "https://[2001:db8::2:1]:5984/my-aws-endpoint",
			},
		},
		"azurekeyvault ipv4": {
			config: `
seal "azurekeyvault" {
  tenant_id      = "46646709-b63e-4747-be42-516edeaf1e14"
  client_id      = "03dc33fc-16d9-4b77-8152-3ec568f8af6e"
  client_secret  = "DUJDS3..."
  vault_name     = "hc-vault"
  key_name       = "vault_key"
	resource       = "vault.azure.net"
}`,
			expected: map[string]string{
				"tenant_id":     "46646709-b63e-4747-be42-516edeaf1e14",
				"client_id":     "03dc33fc-16d9-4b77-8152-3ec568f8af6e",
				"client_secret": "DUJDS3...",
				"vault_name":    "hc-vault",
				"key_name":      "vault_key",
				"resource":      "vault.azure.net",
			},
		},
		"azurekeyvault ipv6": {
			config: `
seal "azurekeyvault" {
  tenant_id      = "46646709-b63e-4747-be42-516edeaf1e14"
  client_id      = "03dc33fc-16d9-4b77-8152-3ec568f8af6e"
  client_secret  = "DUJDS3..."
  vault_name     = "hc-vault"
  key_name       = "vault_key"
  resource       = "2001:db8:0:0:0:0:2:1",
}`,
			expected: map[string]string{
				"tenant_id":     "46646709-b63e-4747-be42-516edeaf1e14",
				"client_id":     "03dc33fc-16d9-4b77-8152-3ec568f8af6e",
				"client_secret": "DUJDS3...",
				"vault_name":    "hc-vault",
				"key_name":      "vault_key",
				"resource":      "2001:db8::2:1",
			},
		},
		"ocikms ipv4": {
			config: `
seal "ocikms" {
	key_id               = "ocid1.key.oc1.iad.afnxza26aag4s.abzwkljsbapzb2nrha5nt3s7s7p42ctcrcj72vn3kq5qx"
	crypto_endpoint      = "https://afnxza26aag4s-crypto.kms.us-ashburn-1.oraclecloud.com"
	management_endpoint  = "https://afnxza26aag4s-management.kms.us-ashburn-1.oraclecloud.com"
	auth_type_api_key    = "true"
}`,
			expected: map[string]string{
				"key_id":              "ocid1.key.oc1.iad.afnxza26aag4s.abzwkljsbapzb2nrha5nt3s7s7p42ctcrcj72vn3kq5qx",
				"crypto_endpoint":     "https://afnxza26aag4s-crypto.kms.us-ashburn-1.oraclecloud.com",
				"management_endpoint": "https://afnxza26aag4s-management.kms.us-ashburn-1.oraclecloud.com",
				"auth_type_api_key":   "true",
			},
		},
		"ocikms ipv6": {
			config: `
seal "ocikms" {
  key_id               = "https://[2001:db8:0:0:0:0:2:1]/abzwkljsbapzb2nrha5nt3s7s7p42ctcrcj72vn3kq5qx"
  crypto_endpoint      = "https://[2001:db8:0:0:0:0:2:1]/afnxza26aag4s-crypto"
  management_endpoint  = "https://[2001:db8:0:0:0:0:2:1]/afnxza26aag4s-management"
  auth_type_api_key    = "true"
}`,
			expected: map[string]string{
				"key_id":              "https://[2001:db8::2:1]/abzwkljsbapzb2nrha5nt3s7s7p42ctcrcj72vn3kq5qx",
				"crypto_endpoint":     "https://[2001:db8::2:1]/afnxza26aag4s-crypto",
				"management_endpoint": "https://[2001:db8::2:1]/afnxza26aag4s-management",
				"auth_type_api_key":   "true",
			},
		},
		"transit ipv4": {
			config: `
seal "transit" {
  address            = "https://vault:8200"
  token              = "s.Qf1s5zigZ4OX6akYjQXJC1jY"
  disable_renewal    = "false"
  key_name           = "transit_key_name"
  mount_path         = "transit/"
  namespace          = "ns1/"
}
`,
			expected: map[string]string{
				"address":         "https://vault:8200",
				"token":           "s.Qf1s5zigZ4OX6akYjQXJC1jY",
				"disable_renewal": "false",
				"key_name":        "transit_key_name",
				"mount_path":      "transit/",
				"namespace":       "ns1/",
			},
		},
		"transit ipv6": {
			config: `
seal "transit" {
  address            = "https://[2001:db8:0:0:0:0:2:1]:8200"
  token              = "s.Qf1s5zigZ4OX6akYjQXJC1jY"
  disable_renewal    = "false"
  key_name           = "transit_key_name"
  mount_path         = "transit/"
  namespace          = "ns1/"
}
`,
			expected: map[string]string{
				"address":         "https://[2001:db8::2:1]:8200",
				"token":           "s.Qf1s5zigZ4OX6akYjQXJC1jY",
				"disable_renewal": "false",
				"key_name":        "transit_key_name",
				"mount_path":      "transit/",
				"namespace":       "ns1/",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			kmses, err := ParseKMSes(tc.config)
			require.NoError(t, err)
			require.Len(t, kmses, 1)
			require.EqualValues(t, tc.expected, kmses[0].Config)
		})
	}
}

// TestMergeKMSEnvConfigAddrConformance tests that all env config whose values
// can be URLs, IP addresses, or host:port addresses, when configured with an
// an IPv6 address, the normalized to be conformant with RFC-5942 §4
// See: https://rfc-editor.org/rfc/rfc5952.html
func TestMergeKMSEnvConfigAddrConformance(t *testing.T) {
	for name, tc := range map[string]struct {
		sealType  string // default to name if none given
		kmsConfig map[string]string
		envVars   map[string]string
		expected  map[string]string
	}{
		"alicloudkms": {
			kmsConfig: map[string]string{
				"region":     "us-east-1",
				"domain":     "kms.us-east-1.aliyuncs.com",
				"access_key": "0wNEpMMlzy7szvai",
				"secret_key": "PupkTg8jdmau1cXxYacgE736PJj4cA",
				"kms_key_id": "08c33a6f-4e0a-4a1b-a3fa-7ddfa1d4fb73",
			},
			envVars: map[string]string{"ALICLOUD_DOMAIN": "2001:db8:0:0:0:0:2:1"},
			expected: map[string]string{
				"region":     "us-east-1",
				"domain":     "2001:db8::2:1",
				"access_key": "0wNEpMMlzy7szvai",
				"secret_key": "PupkTg8jdmau1cXxYacgE736PJj4cA",
				"kms_key_id": "08c33a6f-4e0a-4a1b-a3fa-7ddfa1d4fb73",
			},
		},
		"awskms": {
			kmsConfig: map[string]string{
				"region":     "us-east-1",
				"access_key": "AKIAIOSFODNN7EXAMPLE",
				"secret_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				"kms_key_id": "19ec80b0-dfdd-4d97-8164-c6examplekey",
				"endpoint":   "https://vpce-0e1bb1852241f8cc6-pzi0do8n.kms.us-east-1.vpce.amazonaws.com",
			},
			envVars: map[string]string{"AWS_KMS_ENDPOINT": "https://[2001:db8:0:0:0:0:2:1]:5984/my-aws-endpoint"},
			expected: map[string]string{
				"region":     "us-east-1",
				"access_key": "AKIAIOSFODNN7EXAMPLE",
				"secret_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				"kms_key_id": "19ec80b0-dfdd-4d97-8164-c6examplekey",
				"endpoint":   "https://[2001:db8::2:1]:5984/my-aws-endpoint",
			},
		},
		"azurekeyvault": {
			kmsConfig: map[string]string{
				"tenant_id":     "46646709-b63e-4747-be42-516edeaf1e14",
				"client_id":     "03dc33fc-16d9-4b77-8152-3ec568f8af6e",
				"client_secret": "DUJDS3...",
				"vault_name":    "hc-vault",
				"key_name":      "vault_key",
				"resource":      "vault.azure.net",
			},
			envVars: map[string]string{"AZURE_AD_RESOURCE": "2001:db8:0:0:0:0:2:1"},
			expected: map[string]string{
				"tenant_id":     "46646709-b63e-4747-be42-516edeaf1e14",
				"client_id":     "03dc33fc-16d9-4b77-8152-3ec568f8af6e",
				"client_secret": "DUJDS3...",
				"vault_name":    "hc-vault",
				"key_name":      "vault_key",
				"resource":      "2001:db8::2:1",
			},
		},
		"ocikms wrapper env vars": {
			sealType: "ocikms",
			kmsConfig: map[string]string{
				"key_id":              "ocid1.key.oc1.iad.afnxza26aag4s.abzwkljsbapzb2nrha5nt3s7s7p42ctcrcj72vn3kq5qx",
				"crypto_endpoint":     "https://afnxza26aag4s-crypto.kms.us-ashburn-1.oraclecloud.com",
				"management_endpoint": "https://afnxza26aag4s-management.kms.us-ashburn-1.oraclecloud.com",
				"auth_type_api_key":   "true",
			},
			envVars: map[string]string{
				ocikms.EnvOciKmsWrapperKeyId:              "https://[2001:db8:0:0:0:0:2:1]/abzwkljsbapzb2nrha5nt3s7s7p42ctcrcj72vn3kq5qx",
				ocikms.EnvOciKmsWrapperCryptoEndpoint:     "https://[2001:db8:0:0:0:0:2:1]/afnxza26aag4s-crypto",
				ocikms.EnvOciKmsWrapperManagementEndpoint: "https://[2001:db8:0:0:0:0:2:1]/afnxza26aag4s-management",
			},
			expected: map[string]string{
				"key_id":              "https://[2001:db8::2:1]/abzwkljsbapzb2nrha5nt3s7s7p42ctcrcj72vn3kq5qx",
				"crypto_endpoint":     "https://[2001:db8::2:1]/afnxza26aag4s-crypto",
				"management_endpoint": "https://[2001:db8::2:1]/afnxza26aag4s-management",
				"auth_type_api_key":   "true",
			},
		},
		"ocikms vault env vars": {
			sealType: "ocikms",
			kmsConfig: map[string]string{
				"key_id":              "ocid1.key.oc1.iad.afnxza26aag4s.abzwkljsbapzb2nrha5nt3s7s7p42ctcrcj72vn3kq5qx",
				"crypto_endpoint":     "https://afnxza26aag4s-crypto.kms.us-ashburn-1.oraclecloud.com",
				"management_endpoint": "https://afnxza26aag4s-management.kms.us-ashburn-1.oraclecloud.com",
				"auth_type_api_key":   "true",
			},
			envVars: map[string]string{
				ocikms.EnvVaultOciKmsSealKeyId:              "https://[2001:db8:0:0:0:0:2:1]/abzwkljsbapzb2nrha5nt3s7s7p42ctcrcj72vn3kq5qx",
				ocikms.EnvVaultOciKmsSealCryptoEndpoint:     "https://[2001:db8:0:0:0:0:2:1]/afnxza26aag4s-crypto",
				ocikms.EnvVaultOciKmsSealManagementEndpoint: "https://[2001:db8:0:0:0:0:2:1]/afnxza26aag4s-management",
			},
			expected: map[string]string{
				"key_id":              "https://[2001:db8::2:1]/abzwkljsbapzb2nrha5nt3s7s7p42ctcrcj72vn3kq5qx",
				"crypto_endpoint":     "https://[2001:db8::2:1]/afnxza26aag4s-crypto",
				"management_endpoint": "https://[2001:db8::2:1]/afnxza26aag4s-management",
				"auth_type_api_key":   "true",
			},
		},
		"transit addr not in config": {
			sealType: "transit",
			kmsConfig: map[string]string{
				"token":           "s.Qf1s5zigZ4OX6akYjQXJC1jY",
				"disable_renewal": "false",
				"key_name":        "transit_key_name",
				"mount_path":      "transit/",
				"namespace":       "ns1/",
			},
			envVars: map[string]string{"VAULT_ADDR": "https://[2001:db8:0:0:0:0:2:1]:8200"},
			expected: map[string]string{
				// NOTE: If our address has not been configured we'll fall back to VAULT_ADDR for transit.
				"address":         "https://[2001:db8::2:1]:8200",
				"token":           "s.Qf1s5zigZ4OX6akYjQXJC1jY",
				"disable_renewal": "false",
				"key_name":        "transit_key_name",
				"mount_path":      "transit/",
				"namespace":       "ns1/",
			},
		},
		"transit addr in config": {
			sealType: "transit",
			kmsConfig: map[string]string{
				"address":         "https://vault:8200",
				"token":           "s.Qf1s5zigZ4OX6akYjQXJC1jY",
				"disable_renewal": "false",
				"key_name":        "transit_key_name",
				"mount_path":      "transit/",
				"namespace":       "ns1/",
			},
			envVars: map[string]string{"VAULT_ADDR": "https://[2001:db8:0:0:0:0:2:1]:8200"},
			expected: map[string]string{
				// NOTE: If our address has been configured we don't consider VAULT_ADDR
				"address":         "https://vault:8200",
				"token":           "s.Qf1s5zigZ4OX6akYjQXJC1jY",
				"disable_renewal": "false",
				"key_name":        "transit_key_name",
				"mount_path":      "transit/",
				"namespace":       "ns1/",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			typ := name
			if tc.sealType != "" {
				typ = tc.sealType
			}
			kms := &KMS{
				Type:   typ,
				Config: tc.kmsConfig,
			}

			for envName, envVal := range tc.envVars {
				t.Setenv(envName, envVal)
			}

			require.NoError(t, mergeKMSEnvConfig(kms))
			require.EqualValues(t, tc.expected, kms.Config)
		})
	}
}
