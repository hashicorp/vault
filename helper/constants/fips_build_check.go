// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build (!fips && (fips_140_2 || fips_140_3)) || (fips && !fips_140_2 && !fips_140_3) || (fips_140_2 && fips_140_3)

package constants

import "C"

// This function is the equivalent of an external (CGo) function definition,
// without implementation in any imported or built library. This results in
// a linker err if the above build constraints are satisfied:
//
//	/home/cipherboy/GitHub/cipherboy/vault-enterprise/helper/constants/fips_build_check.go:10: undefined reference to `github.com/hashicorp/vault/helper/constants.VaultFIPSBuildRequiresVersionAgnosticTagAndOneVersionTag'
//
// This indicates that a build error has occurred due to mismatched tags.
//
// In particular, we use this to enforce the following restrictions on build
// tags:
//
//   - If a versioned fips_140_* tag is specified, the unversioned tag must
//     also be.
//   - If the unversioned tag is specified, a versioned tag must be.
//   - Both versioned flags cannot be specified at the same time.
//
// In the unlikely event that a FFI implementation for this function exists
// in the future, it should be renamed to a new function which does not
// exist.
//
// This approach was chosen above the other implementation in fips_cgo_check.go
// because this version does not break static analysis tools: most tools do not
// cross the CGo boundary and thus do not know that the below function is
// missing an implementation. However, in the other file, the function call is
// not marked as CGo (in part large because the lack of a cgo build tag
// prohibits us from using the same technique) and thus it must be a Go
// declaration, that is missing.
func VaultFIPSBuildRequiresVersionAgnosticTagAndOneVersionTag()

func init() {
	VaultFIPSBuildRequiresVersionAgnosticTagAndOneVersionTag()
}
