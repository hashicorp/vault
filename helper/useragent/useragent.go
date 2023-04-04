// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package useragent

import (
	"fmt"
	"runtime"

	"github.com/hashicorp/vault/version"
)

var (
	// projectURL is the project URL.
	projectURL = "https://www.vaultproject.io/"

	// rt is the runtime - variable for tests.
	rt = runtime.Version()

	// versionFunc is the func that returns the current version. This is a
	// function to take into account the different build processes and distinguish
	// between enterprise and oss builds.
	versionFunc = func() string {
		return version.GetVersion().VersionNumber()
	}
)

// String returns the consistent user-agent string for Vault.
//
// e.g. Vault/0.10.4 (+https://www.vaultproject.io/; go1.10.1)
func String() string {
	return fmt.Sprintf("Vault/%s (+%s; %s)",
		versionFunc(), projectURL, rt)
}
