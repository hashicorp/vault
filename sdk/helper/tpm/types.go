// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package tpm

// TPMConfig identifies a TPM device. It is the common device-address type
// shared by Attest, BuildTLSConfig, and EndorsementKeys.
type TPMConfig struct {
	// DevicePath is an optional path to the TPM device. When empty the
	// platform default is used.
	DevicePath string
}
