// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import "os"

func init() {
	if signed := os.Getenv("VAULT_LICENSE_CI"); signed != "" {
		// nosemgrep: tools.semgrep.ci.os-setenv-in-tests -- runs in init(), no *testing.T is available here.
		os.Setenv(EnvVaultLicense, signed)
	}
}
