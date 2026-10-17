// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// generateRootCA creates a root CA via root/generate/exported.
func generateRootCA(t *testing.T, b *backend, s logical.Storage) (issuerID, keyID, privateKeyPEM, certPEM string) {
	t.Helper()
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/generate/exported",
		Storage:   s,
		Data: map[string]interface{}{
			"common_name": "test-ca",
			"key_type":    "ec",
			"key_bits":    256,
			"ttl":         "1h",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.False(t, resp.IsError(), "unexpected error generating root CA: %v", resp.Error())

	issuerID = fmt.Sprintf("%v", resp.Data["issuer_id"])
	keyID = fmt.Sprintf("%v", resp.Data["key_id"])
	privateKeyPEM = resp.Data["private_key"].(string)
	certPEM = resp.Data["certificate"].(string)
	return issuerID, keyID, privateKeyPEM, certPEM
}

// exportCAKey wraps the CA key via WRITE /pki/keys/:id/export using a fresh RSA wrapping key.
func exportCAKey(t *testing.T, b *backend, s logical.Storage, caKeyID string) (wrappedKeyB64, exportKeyHMAC string) {
	t.Helper()
	pubPEM, _ := generateRSAWrapKey(t, 2048)

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "keys/" + caKeyID + "/export",
		Storage:    s,
		Data:       map[string]interface{}{"public_key": pubPEM},
		MountPoint: "pki/",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.False(t, resp.IsError(), "unexpected error exporting CA key: %v", resp.Error())

	return resp.Data["wrapped_key"].(string), resp.Data["export_key_hmac"].(string)
}

// readIssuer reads GET /pki/issuer/:id.
func readIssuer(t *testing.T, b *backend, s logical.Storage, issuerID string) *logical.Response {
	t.Helper()
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.ReadOperation,
		Path:       "issuer/" + issuerID,
		Storage:    s,
		MountPoint: "pki/",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.False(t, resp.IsError(), "unexpected error reading issuer: %v", resp.Error())
	return resp
}

// warningsContaining filters resp.Warnings to those containing substr.
func warningsContaining(resp *logical.Response, substr string) []string {
	var matched []string
	for _, w := range resp.Warnings {
		if strings.Contains(w, substr) {
			matched = append(matched, w)
		}
	}
	return matched
}

// extractTimestampFromWarning parses the RFC3339 timestamp embedded in an export warning of
// the form "private key exported on <RFC3339> with <HMAC>".
func extractTimestampFromWarning(t *testing.T, warning string) time.Time {
	t.Helper()
	// Warning format: "private key exported on <timestamp> with <hmac>"
	after, found := strings.CutPrefix(warning, "private key exported on ")
	require.True(t, found, "warning does not start with expected prefix: %q", warning)
	timestampStr, _, found := strings.Cut(after, " with ")
	require.True(t, found, "warning missing ' with ' separator: %q", warning)
	ts, err := time.Parse(time.RFC3339, timestampStr)
	require.NoError(t, err, "timestamp in warning is not valid RFC3339: %q", timestampStr)
	return ts
}

// TestCAKeyExportRecord_WarnOnIssuerRead verifies that after a CA key is exported via
// WRITE /pki/keys/:id/export, reading the associated issuer produces a warning containing
// a valid RFC3339 timestamp and the export key HMAC.
func TestCAKeyExportRecord_WarnOnIssuerRead(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	before := time.Now().UTC().Truncate(time.Second)
	issuerID, keyID, _, _ := generateRootCA(t, b, s)
	_, exportHMAC := exportCAKey(t, b, s, keyID)
	after := time.Now().UTC()

	resp := readIssuer(t, b, s, issuerID)

	warnings := warningsContaining(resp, "private key exported on")
	require.Len(t, warnings, 1, "expected exactly one export warning, got: %v", resp.Warnings)
	require.Contains(t, warnings[0], exportHMAC, "warning should contain the export key HMAC")

	ts := extractTimestampFromWarning(t, warnings[0])
	require.False(t, ts.Before(before), "warning timestamp %v should not be before export started at %v", ts, before)
	require.False(t, ts.After(after), "warning timestamp %v should not be after export completed at %v", ts, after)
}

// TestCAKeyExportRecord_NoWarnWhenNotExported verifies that reading an issuer whose
// key has never been exported via the BYOK path produces no export warnings.
func TestCAKeyExportRecord_NoWarnWhenNotExported(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	issuerID, _, _, _ := generateRootCA(t, b, s)

	resp := readIssuer(t, b, s, issuerID)

	warnings := warningsContaining(resp, "private key exported on")
	require.Empty(t, warnings, "expected no export warnings for an unexported key, got: %v", resp.Warnings)
}

// TestCAKeyExportRecord_MultipleExportsProduceMultipleWarnings verifies that exporting the
// same CA key twice results in two distinct warnings on the issuer read, each carrying the
// respective export key HMAC.
func TestCAKeyExportRecord_MultipleExportsProduceMultipleWarnings(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	issuerID, keyID, _, _ := generateRootCA(t, b, s)

	_, hmac1 := exportCAKey(t, b, s, keyID)
	_, hmac2 := exportCAKey(t, b, s, keyID)

	resp := readIssuer(t, b, s, issuerID)

	warnings := warningsContaining(resp, "private key exported on")
	require.Len(t, warnings, 2, "expected two export warnings, got: %v", resp.Warnings)

	allWarnings := strings.Join(warnings, " ")
	require.Contains(t, allWarnings, hmac1, "warnings should contain first export HMAC")
	require.Contains(t, allWarnings, hmac2, "warnings should contain second export HMAC")
}

// TestCAKeyExportRecord_WarningPersistsAfterDeleteAndReimport verifies that export warnings
// survive a full delete-and-reimport cycle of the key and issuer.  The export record is keyed
// on the SHA-256 fingerprint of the CA key's SubjectPublicKeyInfo, which is identical before
// and after reimport, so the warning must still appear on the new issuer UUID.
func TestCAKeyExportRecord_WarningPersistsAfterDeleteAndReimport(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	issuerID, keyID, privateKeyPEM, certPEM := generateRootCA(t, b, s)
	_, exportHMAC := exportCAKey(t, b, s, keyID)

	// Confirm the warning is present before deletion.
	resp := readIssuer(t, b, s, issuerID)
	require.NotEmpty(t, warningsContaining(resp, "private key exported on"),
		"expected export warning before delete")

	// Delete the issuer first (key deletion fails if it is still referenced).
	_, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.DeleteOperation,
		Path:       "issuer/" + issuerID,
		Storage:    s,
		MountPoint: "pki/",
	})
	require.NoError(t, err)

	// Delete the key.
	_, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.DeleteOperation,
		Path:       "key/" + keyID,
		Storage:    s,
		MountPoint: "pki/",
	})
	require.NoError(t, err)

	// Reimport the original key + certificate as a bundle so the issuer is
	// automatically linked to the reimported key.
	reimportResp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "issuers/import/bundle",
		Storage:   s,
		Data: map[string]interface{}{
			"pem_bundle": privateKeyPEM + "\n" + certPEM,
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err)
	require.NotNil(t, reimportResp)
	require.False(t, reimportResp.IsError(), "unexpected error reimporting bundle: %v", reimportResp.Error())

	importedIssuers, ok := reimportResp.Data["imported_issuers"].([]string)
	require.True(t, ok && len(importedIssuers) == 1,
		"expected exactly one imported issuer, got: %v", reimportResp.Data["imported_issuers"])
	newIssuerID := importedIssuers[0]

	// Read the reimported issuer and verify the warning is still present.
	resp = readIssuer(t, b, s, newIssuerID)

	warnings := warningsContaining(resp, "private key exported on")
	require.NotEmpty(t, warnings, "expected export warning to survive delete-and-reimport, got none")
	require.Contains(t, warnings[0], exportHMAC,
		"warning should still reference the original export key HMAC after reimport")
}

// TestCaPublicKeyFingerprint_MLKEM verifies that caPublicKeyFingerprint correctly parses
// ML-KEM-768 and ML-KEM-1024 private keys, extracting and hashing their public encapsulation keys
// securely (without exposing or hashing the raw private key material directly).
func TestCaPublicKeyFingerprint_MLKEM(t *testing.T) {
	t.Parallel()

	for _, keyType := range []string{"ml-kem-768", "ml-kem-1024"} {
		keyType := keyType
		t.Run(keyType, func(t *testing.T) {
			t.Parallel()

			privPEM, _, pubDER, err := generateExportKeypair(keyType)
			require.NoError(t, err)

			// Generate fingerprint from private key PEM
			gotFingerprint, err := caPublicKeyFingerprint(privPEM)
			require.NoError(t, err)

			// Manually compute expected fingerprint over public key DER
			sum := sha256.Sum256(pubDER)
			wantFingerprint := "sha256:" + hex.EncodeToString(sum[:])

			require.Equal(t, wantFingerprint, gotFingerprint, "fingerprint must be computed securely over the public key bytes")
		})
	}
}
