// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package seal

func ValidateMultiSealGenerationInfo(_ bool, _, _ *SealGenerationInfo, _ bool) error {
	return nil
}
