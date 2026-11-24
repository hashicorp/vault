// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package tpm

// KeyReference contains metadata for locating and using a TPM key.
// It is typically stored as JSON alongside the certificate and key blobs.
type KeyReference struct {
	// Type is the key type, currently only "tpm2" is supported
	Type string `json:"type"`

	// Version is the key reference format version
	Version int `json:"version"`

	// PublicKeySHA256 is the SHA256 fingerprint of the public key
	PublicKeySHA256 string `json:"public_key_sha256"`
}

// TPMKeyMaterial contains the in-memory TPM key data needed for cryptographic operations.
// This type is used when working with TPM keys without requiring disk persistence.
type TPMKeyMaterial struct {
	// Type is the key type, currently only "tpm2" is supported
	Type string

	// AKBlob is the attestation key blob data
	AKBlob []byte

	// AppBlob is the application key blob data
	AppBlob []byte

	// PublicKeySHA256 is the SHA256 fingerprint of the public key
	PublicKeySHA256 string
}
