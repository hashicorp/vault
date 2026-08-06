// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package pluginutil

// IsEnterpriseVersion checks for the "version" metadata identifier in a plugin's
// semantic version
func IsEnterpriseVersion(v string) bool {
	return false
}
