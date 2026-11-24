// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

// Package tpm provides utilities for TPM-based attestation and certificate
// management in Vault. It implements a two-step attestation workflow that
// allows clients to obtain certificates backed by TPM hardware.
//
// The package provides:
//   - TPMConfig: Device-address type used by all TPM operations
//   - Attest(): Performs the two-step TPM attestation workflow to obtain a certificate
//   - BuildTLSConfig(): Creates a tls.Config using TPM-backed certificates
//   - EndorsementKeys(): Returns the endorsement keys for the identified TPM
//   - TPMSigner: A crypto.Signer implementation that uses TPM hardware
//   - KeyReference: Metadata for locating and using TPM keys using files on disk
//   - TPMKeyMaterial: The in-memory equivalent of KeyReference
//
// Example usage:
//
//	cfg := tpm.TPMConfig{DevicePath: "/dev/tpm0"}
//
//	// Fetch endorsement keys to enroll the TPM.
//	eks, err := tpm.EndorsementKeys(cfg)
//	if err != nil {
//	    return err
//	}
//
//	// Run the two-phase attestation workflow to obtain a certificate.
//	result, err := tpm.Attest(ctx, client, &tpm.AttestationConfig{
//	    MountPath: "auth/tpm",
//	    TPMConfig: cfg,
//	})
//	if err != nil {
//	    return err
//	}
//
//	// Build a TLS config that presents the TPM-backed certificate.
//	tlsConfig, err := tpm.BuildTLSConfig(result, caCert, caPath, cfg)
//	if err != nil {
//	    return err
//	}
//
//	client.SetTLSConfig(tlsConfig)
package tpm
