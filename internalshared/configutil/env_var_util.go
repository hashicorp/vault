package configutil

var (
	AliCloudKMSEnvVars = map[string]string{
		"ALICLOUD_REGION":               "region",
		"ALICLOUD_DOMAIN":               "domain",
		"ALICLOUD_ACCESS_KEY":           "access_key",
		"ALICLOUD_SECRET_KEY":           "secret_key",
		"VAULT_ALICLOUDKMS_SEAL_KEY_ID": "kms_key_id",
	}

	AWSKMSEnvVars = map[string]string{
		"AWS_REGION":               "region",
		"AWS_DEFAULT_REGION":       "region",
		"AWS_ACCESS_KEY_ID":        "access_key",
		"AWS_SESSION_TOKEN":        "session_token",
		"AWS_SECRET_ACCESS_KEY":    "secret_key",
		"VAULT_AWSKMS_SEAL_KEY_ID": "kms_key_id",
		"AWS_KMS_ENDPOINT":         "endpoint",
	}

	AzureEnvVars = map[string]string{
		"AZURE_TENANT_ID":                "tenant_id",
		"AZURE_CLIENT_ID":                "client_id",
		"AZURE_CLIENT_SECRET":            "client_secret",
		"AZURE_ENVIRONMENT":              "environment",
		"VAULT_AZUREKEYVAULT_VAULT_NAME": "vault_name",
		"VAULT_AZUREKEYVAULT_KEY_NAME":   "key_name",
		"AZURE_AD_RESOURCE":              "resource",
	}

	GCPCKMSEnvVars = map[string]string{
		"GOOGLE_CREDENTIALS":             "credentials",
		"GOOGLE_APPLICATION_CREDENTIALS": "credentials",
		"GOOGLE_PROJECT":                 "project",
		"GOOGLE_REGION":                  "region",
		"VAULT_GCPCKMS_SEAL_KEY_RING":    "key_ring",
		"VAULT_GCPCKMS_SEAL_CRYPTO_KEY":  "crypto_key",
	}

	OCIKMSEnvVars = map[string]string{
		"VAULT_OCIKMS_SEAL_KEY_ID":         "key_id",
		"VAULT_OCIKMS_CRYPTO_ENDPOINT":     "crypto_endpoint",
		"VAULT_OCIKMS_MANAGEMENT_ENDPOINT": "management_endpoint",
	}

	TransitEnvVars = map[string]string{
		"VAULT_ADDR":                         "address",
		"VAULT_TOKEN":                        "token",
		"VAULT_TRANSIT_SEAL_KEY_NAME":        "key_name",
		"VAULT_TRANSIT_SEAL_MOUNT_PATH":      "mount_path",
		"VAULT_NAMESPACE":                    "namespace",
		"VAULT_TRANSIT_SEAL_DISABLE_RENEWAL": "disable_renewal",
		"VAULT_CACERT":                       "tls_ca_cert",
		"VAULT_CLIENT_CERT":                  "tls_client_cert",
		"VAULT_CLIENT_KEY":                   "tls_client_key",
		"VAULT_TLS_SERVER_NAME":              "tls_server_name",
		"VAULT_SKIP_VERIFY":                  "tls_skip_verify",
	}
)
