// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testcluster

const (
	// EnvVaultLicenseCI is the name of an environment variable that contains
	// a signed license string used for Vault Enterprise binary-based tests.
	// The binary will be run with the env var VAULT_LICENSE set to this value.
	EnvVaultLicenseCI = "VAULT_LICENSE_CI"

	// DefaultCAFile is the path to the CA file. This is a docker-specific
	// constant. TODO: needs to be moved to a more relevant place
	DefaultCAFile = "/vault/config/ca.pem"
)
