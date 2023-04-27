package testcluster

const (
	// EnvVaultLicenseCI is the name of an environment variable that contains
	// a signed license string used for Vault Enterprise binary-based tests.
	// The binary will be run with the env var VAULT_LICENSE set to this value.
	EnvVaultLicenseCI = "VAULT_LICENSE_CI"
)
