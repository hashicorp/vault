// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

// Package tpmutil holds TPM helpers that both Vault and the Vault CLI use.
package tpmutil

import (
	"crypto/sha256"
	"encoding/hex"
)

// IDFromDER derives the canonical TPM ID from a DER-encoded public key.
// The ID is "sha256-" followed by the hex-encoded SHA-256 of the DER bytes.
// This is the single authoritative definition of the TPM ID algorithm; both
// the identity store enrollment path and the CLI command use this function.
func IDFromDER(der []byte) string {
	sum := sha256.Sum256(der)
	return "sha256-" + hex.EncodeToString(sum[:])
}
