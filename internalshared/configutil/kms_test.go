// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package configutil

import (
	"os"
	"reflect"
	"testing"
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
			"Priority greater than 1 without matching var names",
			&KMS{
				Type:     "awskms",
				Priority: 2,
				Name:     "awskms",
			},
			map[string]string{"AWS_REGION": "test_region", "AWS_ACCESS_KEY_ID": "test_access_key", "AWS_SECRET_ACCESS_KEY": "test_secret_key", "VAULT_AWSKMS_SEAL_KEY_ID": "test_key_id"},
			map[string]string{},
		},
		{
			"Priority greater than 1 with matching var names",
			&KMS{
				Type:     "awskms",
				Priority: 2,
				Name:     "awskms",
			},
			map[string]string{"AWS_REGION_awskms": "test_region", "AWS_ACCESS_KEY_ID_awskms": "test_access_key", "AWS_SECRET_ACCESS_KEY_awskms": "test_secret_key", "VAULT_AWSKMS_SEAL_KEY_ID_awskms": "test_key_id"},
			map[string]string{"region": "test_region", "access_key": "test_access_key", "secret_key": "test_secret_key", "kms_key_id": "test_key_id"},
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
				if err := os.Setenv(envName, envVal); err != nil {
					t.Errorf("error setting environment vars for test: %s", err)
				}
			}

			if got := getEnvConfig(tt.kms); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEnvConfig() = %v, want %v", got, tt.want)
			}

			for env := range tt.envVars {
				if err := os.Unsetenv(env); err != nil {
					t.Errorf("error unsetting environment vars for test: %s", err)
				}
			}
		})
	}
}
