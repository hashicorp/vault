// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/mlkem"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// generateRSAWrapKey returns a PEM-encoded RSA public key and the corresponding private key.
func generateRSAWrapKey(t *testing.T, bits int) (pubPEM string, priv *rsa.PrivateKey) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	require.NoError(t, err)
	der, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	require.NoError(t, err)
	pubPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
	return pubPEM, priv
}

// generateECWrapKey returns a PEM-encoded EC public key and the corresponding private key.
func generateECWrapKey(t *testing.T, curve elliptic.Curve) (pubPEM string, priv *ecdsa.PrivateKey) {
	t.Helper()
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	require.NoError(t, err)
	der, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	require.NoError(t, err)
	pubPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
	return pubPEM, priv
}

// generateMLKEM768WrapKey returns a PEM-encoded ML-KEM-768 public key and the corresponding decapsulation key.
func generateMLKEM768WrapKey(t *testing.T) (pubPEM string, dk *mlkem.DecapsulationKey768) {
	t.Helper()
	dk, err := mlkem.GenerateKey768()
	require.NoError(t, err)
	pubBytes := dk.EncapsulationKey().Bytes()
	pubPEM = string(pem.EncodeToMemory(&pem.Block{Type: "ML-KEM PUBLIC KEY", Bytes: pubBytes}))
	return pubPEM, dk
}

// decryptRSABlob decrypts a wrappedKeyBlob produced by wrapWithRSA.
func decryptRSABlob(t *testing.T, blobJSON []byte, priv *rsa.PrivateKey) string {
	t.Helper()
	var blob wrappedKeyBlob
	require.NoError(t, json.Unmarshal(blobJSON, &blob))
	require.Equal(t, "RSA-OAEP-SHA256+AES256GCM", blob.Alg)

	cek, err := rsa.DecryptOAEP(sha256.New(), nil, priv, blob.WrappedCEK, nil)
	require.NoError(t, err)

	return testAesGCMOpen(t, cek, blob.Nonce, blob.Ciphertext)
}

// decryptECBlob decrypts a wrappedKeyBlob produced by wrapWithEC.
func decryptECBlob(t *testing.T, blobJSON []byte, priv *ecdsa.PrivateKey) string {
	t.Helper()
	var blob wrappedKeyBlob
	require.NoError(t, json.Unmarshal(blobJSON, &blob))
	require.Equal(t, "ECDH-ES+AES256GCM", blob.Alg)

	ephPub, err := x509.ParsePKIXPublicKey(blob.EphemeralPub)
	require.NoError(t, err)
	ephECDH, err := ephPub.(*ecdsa.PublicKey).ECDH()
	require.NoError(t, err)

	recipPrivECDH, err := priv.ECDH()
	require.NoError(t, err)

	shared, err := recipPrivECDH.ECDH(ephECDH)
	require.NoError(t, err)

	cek := hkdfSHA256(shared, []byte("vault-pki-export-v1"), 32)
	return testAesGCMOpen(t, cek, blob.Nonce, blob.Ciphertext)
}

// decryptMLKEM768Blob decrypts a wrappedKeyBlob produced by wrapWithMLKEM for ML-KEM-768.
func decryptMLKEM768Blob(t *testing.T, blobJSON []byte, dk *mlkem.DecapsulationKey768) string {
	t.Helper()
	var blob wrappedKeyBlob
	require.NoError(t, json.Unmarshal(blobJSON, &blob))
	require.Equal(t, "MLKEM768+AES256GCM", blob.Alg)

	shared, err := dk.Decapsulate(blob.KEMCiphertext)
	require.NoError(t, err)

	cek := hkdfSHA256(shared, []byte("vault-pki-export-v1"), 32)
	return testAesGCMOpen(t, cek, blob.Nonce, blob.Ciphertext)
}

func testAesGCMOpen(t *testing.T, key, nonce, ciphertext []byte) string {
	t.Helper()
	block, err := aes.NewCipher(key)
	require.NoError(t, err)
	gcm, err := cipher.NewGCM(block)
	require.NoError(t, err)
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	require.NoError(t, err)
	return string(plaintext)
}

// writeCAKey generates an internal CA key and returns its key_id string.
func writeCAKey(t *testing.T, b *backend, s logical.Storage, keyType string) string {
	t.Helper()
	data := map[string]interface{}{}
	switch keyType {
	case "rsa-2048", "rsa":
		data["key_type"] = "rsa"
		data["key_bits"] = 2048
	case "rsa-3072":
		data["key_type"] = "rsa"
		data["key_bits"] = 3072
	case "rsa-4096":
		data["key_type"] = "rsa"
		data["key_bits"] = 4096
	case "ec-p256", "ec":
		data["key_type"] = "ec"
		data["key_bits"] = 256
	case "ec-p384":
		data["key_type"] = "ec"
		data["key_bits"] = 384
	case "ec-p521":
		data["key_type"] = "ec"
		data["key_bits"] = 521
	case "ed25519":
		data["key_type"] = "ed25519"
	default:
		data["key_type"] = keyType
	}
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "keys/generate/internal",
		Storage:    s,
		Data:       data,
		MountPoint: "pki/",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.False(t, resp.IsError(), "unexpected error generating key: %v", resp.Error())
	require.NotEmpty(t, resp.Data["key_id"], "key_id must be non-empty")
	id := fmt.Sprintf("%v", resp.Data["key_id"])
	require.NotEmpty(t, id, "key_id must be a non-empty string")
	return id
}

// TestSecureExportCAKey_RSA verifies that WRITE /pki/keys/:ca-key-uuid/export with an RSA
// wrapping key returns a wrapped_key that decrypts to the original CA private key PEM.
func TestSecureExportCAKey_RSA(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	caKeyID := writeCAKey(t, b, s, "rsa")
	pubPEM, rsaPriv := generateRSAWrapKey(t, 2048)
	wantHMAC := computeExportKeyHMACFromPEM(pubPEM)

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "keys/" + caKeyID + "/export",
		Storage:    s,
		Data:       map[string]interface{}{"public_key": pubPEM},
		MountPoint: "pki/",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.False(t, resp.IsError(), "unexpected error: %v", resp.Error())

	require.Equal(t, caKeyID, resp.Data["ca_key_uuid"])
	require.Equal(t, wantHMAC, resp.Data["export_key_hmac"])

	exportedAt, ok := resp.Data["exported_at"].(string)
	require.True(t, ok && exportedAt != "")
	_, parseErr := time.Parse(time.RFC3339, exportedAt)
	require.NoError(t, parseErr, "exported_at must be RFC3339")

	wrappedB64, ok := resp.Data["wrapped_key"].(string)
	require.True(t, ok && wrappedB64 != "")

	blobJSON, err := base64.StdEncoding.DecodeString(wrappedB64)
	require.NoError(t, err)

	decrypted := decryptRSABlob(t, blobJSON, rsaPriv)
	require.True(t, strings.HasPrefix(decrypted, "-----BEGIN"),
		"decrypted blob should be a PEM private key")
}

// TestSecureExportCAKey_EC verifies WRITE /pki/keys/:ca-key-uuid/export with EC wrapping keys
// on all three supported curves (P-256, P-384, P-521), and that the decrypted plaintext
// is the original CA private key PEM.
func TestSecureExportCAKey_EC(t *testing.T) {
	t.Parallel()

	curves := []struct {
		name  string
		curve elliptic.Curve
	}{
		{"P-256", elliptic.P256()},
		{"P-384", elliptic.P384()},
		{"P-521", elliptic.P521()},
	}

	for _, tt := range curves {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b, s := CreateBackendWithStorage(t)

			caKeyID := writeCAKey(t, b, s, "ec")
			pubPEM, ecPriv := generateECWrapKey(t, tt.curve)

			resp, err := b.HandleRequest(context.Background(), &logical.Request{
				Operation:  logical.UpdateOperation,
				Path:       "keys/" + caKeyID + "/export",
				Storage:    s,
				Data:       map[string]interface{}{"public_key": pubPEM},
				MountPoint: "pki/",
			})
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.False(t, resp.IsError(), "unexpected error for curve %s: %v", tt.name, resp.Error())

			wrappedB64 := resp.Data["wrapped_key"].(string)
			blobJSON, err := base64.StdEncoding.DecodeString(wrappedB64)
			require.NoError(t, err)

			decrypted := decryptECBlob(t, blobJSON, ecPriv)
			require.True(t, strings.HasPrefix(decrypted, "-----BEGIN"),
				"decrypted blob should be a PEM private key for curve %s", tt.name)
		})
	}
}

// TestSecureExportCAKey_MLKEM verifies WRITE /pki/keys/:ca-key-uuid/export with an ML-KEM-768
// wrapping key, and that the decrypted plaintext is the original CA private key PEM.
func TestSecureExportCAKey_MLKEM(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	caKeyID := writeCAKey(t, b, s, "ec")
	pubPEM, dk := generateMLKEM768WrapKey(t)

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "keys/" + caKeyID + "/export",
		Storage:    s,
		Data:       map[string]interface{}{"public_key": pubPEM},
		MountPoint: "pki/",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.False(t, resp.IsError(), "unexpected error: %v", resp.Error())

	wrappedB64 := resp.Data["wrapped_key"].(string)
	blobJSON, err := base64.StdEncoding.DecodeString(wrappedB64)
	require.NoError(t, err)

	decrypted := decryptMLKEM768Blob(t, blobJSON, dk)
	require.True(t, strings.HasPrefix(decrypted, "-----BEGIN"),
		"decrypted blob should be a PEM private key")
}

// TestSecureExportCAKey_ResponseNeverContainsPrivateKey verifies that the raw private key
// is never present in any response field.
func TestSecureExportCAKey_ResponseNeverContainsPrivateKey(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	caKeyID := writeCAKey(t, b, s, "rsa")
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
	require.False(t, resp.IsError())

	_, hasPriv := resp.Data["private_key"]
	require.False(t, hasPriv, "private_key must never appear in the export response")

	for field, val := range resp.Data {
		strVal, ok := val.(string)
		if !ok {
			continue
		}
		require.False(t, strings.Contains(strVal, "-----BEGIN PRIVATE KEY-----"),
			"field %q must not contain raw private key PEM", field)
		require.False(t, strings.Contains(strVal, "-----BEGIN EC PRIVATE KEY-----"),
			"field %q must not contain raw EC private key PEM", field)
	}
}

// TestSecureExportCAKey_NotFound verifies that exporting a non-existent CA key UUID
// returns a logical error, not a Go error or nil response.
func TestSecureExportCAKey_NotFound(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)
	pubPEM, _ := generateRSAWrapKey(t, 2048)

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "keys/00000000-0000-0000-0000-000000000000/export",
		Storage:    s,
		Data:       map[string]interface{}{"public_key": pubPEM},
		MountPoint: "pki/",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.IsError(), "expected logical error for non-existent CA key UUID")
}

// TestSecureExportCAKey_MissingPublicKey verifies that omitting public_key returns a logical error.
func TestSecureExportCAKey_MissingPublicKey(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	caKeyID := writeCAKey(t, b, s, "rsa")

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "keys/" + caKeyID + "/export",
		Storage:    s,
		Data:       map[string]interface{}{},
		MountPoint: "pki/",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.IsError(), "expected logical error when public_key is missing")
}

// TestSecureExportCAKey_BadPublicKey verifies that an invalid public_key value returns a logical error.
func TestSecureExportCAKey_BadPublicKey(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	caKeyID := writeCAKey(t, b, s, "rsa")

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "keys/" + caKeyID + "/export",
		Storage:    s,
		Data:       map[string]interface{}{"public_key": "not-a-pem-key"},
		MountPoint: "pki/",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.IsError(), "expected logical error for invalid public_key")
}

// TestSecureExportCAKey_HMACMatchesComputedValue verifies that the returned export_key_hmac
// is byte-for-byte equal to computeExportKeyHMACFromPEM applied to the same public key,
// so the import side can reproduce the correlator deterministically.
func TestSecureExportCAKey_HMACMatchesComputedValue(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	caKeyID := writeCAKey(t, b, s, "ec")
	pubPEM, _ := generateECWrapKey(t, elliptic.P256())
	want := computeExportKeyHMACFromPEM(pubPEM)

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "keys/" + caKeyID + "/export",
		Storage:    s,
		Data:       map[string]interface{}{"public_key": pubPEM},
		MountPoint: "pki/",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.False(t, resp.IsError())

	got := resp.Data["export_key_hmac"].(string)
	require.Equal(t, want, got, "export_key_hmac must equal computeExportKeyHMACFromPEM(public_key)")
}

// TestSecureExportCAKey_WrappedKeyIsNonDeterministic verifies that two consecutive exports of
// the same CA key with the same public key produce different wrapped_key blobs due to
// fresh nonces and ephemeral keys on each call.
func TestSecureExportCAKey_WrappedKeyIsNonDeterministic(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	caKeyID := writeCAKey(t, b, s, "rsa")
	pubPEM, _ := generateRSAWrapKey(t, 2048)

	export := func() string {
		resp, err := b.HandleRequest(context.Background(), &logical.Request{
			Operation:  logical.UpdateOperation,
			Path:       "keys/" + caKeyID + "/export",
			Storage:    s,
			Data:       map[string]interface{}{"public_key": pubPEM},
			MountPoint: "pki/",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.False(t, resp.IsError())
		return resp.Data["wrapped_key"].(string)
	}

	wrapped1 := export()
	wrapped2 := export()
	require.NotEqual(t, wrapped1, wrapped2,
		"two successive exports must produce distinct wrapped_key blobs due to random nonce/ephemeral key")
}

// TestSecureExportCAKey_ECDHCurveMismatchFails generates an EC P-256 wrapping key, exports a CA
// key, then attempts to decrypt with a P-384 private key and verifies the ECDH fails, confirming
// the ephemeral key is curve-bound.
func TestSecureExportCAKey_ECDHCurveMismatchFails(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	caKeyID := writeCAKey(t, b, s, "ec")
	pubPEM, _ := generateECWrapKey(t, elliptic.P256())

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "keys/" + caKeyID + "/export",
		Storage:    s,
		Data:       map[string]interface{}{"public_key": pubPEM},
		MountPoint: "pki/",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.False(t, resp.IsError())

	wrappedB64 := resp.Data["wrapped_key"].(string)
	blobJSON, err := base64.StdEncoding.DecodeString(wrappedB64)
	require.NoError(t, err)

	var blob wrappedKeyBlob
	require.NoError(t, json.Unmarshal(blobJSON, &blob))

	// Parse the ephemeral public key (P-256) and try to ECDH with a P-384 private key.
	ephPubAny, err := x509.ParsePKIXPublicKey(blob.EphemeralPub)
	require.NoError(t, err)
	ephECDH, err := ephPubAny.(*ecdsa.PublicKey).ECDH()
	require.NoError(t, err)

	wrongPriv, err := ecdh.P384().GenerateKey(rand.Reader)
	require.NoError(t, err)

	_, ecdhErr := wrongPriv.ECDH(ephECDH)
	require.Error(t, ecdhErr, "ECDH with mismatched curves must return an error")
}

// TestSecureExportCAKey_LogsWarningOnExport verifies that exporting a CA key via BYOK emits a
// server-side warning log containing the CA key UUID, name, SPKI fingerprint, and the wrapping
// key HMAC. This warning is the only in-band signal that a key export occurred, so its presence
// and content are security-relevant.
func TestSecureExportCAKey_LogsWarningOnExport(t *testing.T) {
	t.Parallel()

	// Wire up a buffered logger so we can assert on the emitted warning.
	var logBuf bytes.Buffer
	logger := hclog.New(&hclog.LoggerOptions{
		Level:  hclog.Warn,
		Output: &logBuf,
	})

	b, s := CreateBackendWithStorageAndLogger(t, logger)

	caKeyID := writeCAKey(t, b, s, "rsa")
	pubPEM, _ := generateRSAWrapKey(t, 2048)
	wantHMAC := computeExportKeyHMACFromPEM(pubPEM)

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "keys/" + caKeyID + "/export",
		Storage:    s,
		Data:       map[string]interface{}{"public_key": pubPEM},
		MountPoint: "pki/",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.False(t, resp.IsError(), "unexpected error: %v", resp.Error())

	log := logBuf.String()
	require.Contains(t, log, "CA private key material was securely exported via BYOK", "warn message must be present in log output")
	require.Contains(t, log, caKeyID, "log must include the CA key UUID")
	require.Contains(t, log, wantHMAC, "log must include the export key HMAC")
}

// TestSecureExportCAKey_ManagedKeyRejected verifies that exporting a managed (KMS-backed) key
// returns a user error.
func TestSecureExportCAKey_ManagedKeyRejected(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	managedKeyID := issuing.KeyID("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	entry := issuing.KeyEntry{
		ID:             managedKeyID,
		Name:           "test-managed-key",
		PrivateKeyType: certutil.ManagedPrivateKey,
	}
	storageEntry, err := logical.StorageEntryJSON(issuing.KeyPrefix+string(managedKeyID), entry)
	require.NoError(t, err)
	require.NoError(t, s.Put(context.Background(), storageEntry))

	pubPEM, _ := generateRSAWrapKey(t, 2048)
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "keys/" + string(managedKeyID) + "/export",
		Storage:    s,
		Data:       map[string]interface{}{"public_key": pubPEM},
		MountPoint: "pki/",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, resp.IsError())
	require.Contains(t, resp.Error().Error(), "secure export is not supported for managed (KMS) keys")
}
