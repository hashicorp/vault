// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !fips

package constants

// IsFIPS returns true if Vault is operating in a FIPS-140-{2,3} mode.
func IsFIPS() bool {
	return false
}
