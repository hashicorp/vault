// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build (fips || fips_140_2 || fips_140_3) && !cgo

package constants

func init() {
	// See note in fips_build_check.go.
	//
	// This function call is missing a declaration, causing the build to
	// fail on improper tags (fips specified but cgo not specified). This
	// ensures Vault fails to build if a FIPS build is requested but CGo
	// support is not enabled.
	//
	// Note that this could confuse static analysis tools as this function
	// should not ever be defined. If this function is defined in the future,
	// the below reference should be renamed to a new name that is not
	// defined to ensure we get a build failure.
	VaultFIPSBuildTagMustEnableCGo()
}
