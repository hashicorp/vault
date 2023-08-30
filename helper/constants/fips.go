// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !fips

package constants

// IsFIPS returns true if Vault is operating in a FIPS-140-{2,3} mode.
func IsFIPS() bool {
	return false
}
