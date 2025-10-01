// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import "os"

func init() {
	if signed := os.Getenv("VAULT_LICENSE_CI"); signed != "" {
		os.Setenv(EnvVaultLicense, signed)
	}
}
