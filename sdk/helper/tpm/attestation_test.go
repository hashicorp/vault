// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package tpm

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// generateTestCertificate creates a self-signed certificate for testing.
func generateTestCertificate(t *testing.T) []byte {
	t.Helper()

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		t.Fatalf("Failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "test-certificate",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	return certPEM
}

// TestAttestationResult_SaveToDirectory_LoadFromDirectory tests the round-trip
// of saving an AttestationResult to disk and loading it back.
func TestAttestationResult_SaveToDirectory_LoadFromDirectory(t *testing.T) {
	// Arrange: Create a test AttestationResult with all fields populated
	testCertPEM := generateTestCertificate(t)
	testCAChainPEM := generateTestCertificate(t)
	testAKBlob := []byte("test-ak-blob-data-12345")
	testAppBlob := []byte("test-app-blob-data-67890")
	testPublicKeySHA256 := "abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"

	original := &AttestationResult{
		CertPEM:    testCertPEM,
		CAChainPEM: testCAChainPEM,
		KeyMaterial: &TPMKeyMaterial{
			Type:            "tpm2",
			AKBlob:          testAKBlob,
			AppBlob:         testAppBlob,
			PublicKeySHA256: testPublicKeySHA256,
		},
		NotBefore: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		NotAfter:  time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Act: Save to directory
	err := original.SaveToDirectory(tmpDir, 0o022)
	if err != nil {
		t.Fatalf("SaveToDirectory failed: %v", err)
	}

	// Assert: Verify all expected files were created
	expectedFiles := []string{
		"client-key.json",
		"ak.blob",
		"app.blob",
		"ca_chain.pem",
		"client.crt",
	}

	for _, filename := range expectedFiles {
		path := filepath.Join(tmpDir, filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", filename)
		}
	}

	// Act: Load from directory
	loaded, err := LoadFromDirectory(tmpDir)
	if err != nil {
		t.Fatalf("LoadFromDirectory failed: %v", err)
	}

	// Assert: Verify loaded data matches original
	if string(loaded.CertPEM) != string(original.CertPEM) {
		t.Errorf("CertPEM mismatch:\ngot:  %s\nwant: %s", loaded.CertPEM, original.CertPEM)
	}

	if string(loaded.CAChainPEM) != string(original.CAChainPEM) {
		t.Errorf("CAChainPEM mismatch:\ngot:  %s\nwant: %s", loaded.CAChainPEM, original.CAChainPEM)
	}

	if loaded.KeyMaterial == nil {
		t.Fatal("KeyMaterial is nil after loading")
	}

	if loaded.KeyMaterial.Type != original.KeyMaterial.Type {
		t.Errorf("KeyMaterial.Type mismatch: got %q, want %q", loaded.KeyMaterial.Type, original.KeyMaterial.Type)
	}

	if string(loaded.KeyMaterial.AKBlob) != string(original.KeyMaterial.AKBlob) {
		t.Errorf("KeyMaterial.AKBlob mismatch: got %q, want %q", loaded.KeyMaterial.AKBlob, original.KeyMaterial.AKBlob)
	}

	if string(loaded.KeyMaterial.AppBlob) != string(original.KeyMaterial.AppBlob) {
		t.Errorf("KeyMaterial.AppBlob mismatch: got %q, want %q", loaded.KeyMaterial.AppBlob, original.KeyMaterial.AppBlob)
	}

	if loaded.KeyMaterial.PublicKeySHA256 != original.KeyMaterial.PublicKeySHA256 {
		t.Errorf("KeyMaterial.PublicKeySHA256 mismatch: got %q, want %q", loaded.KeyMaterial.PublicKeySHA256, original.KeyMaterial.PublicKeySHA256)
	}

	// Note: NotBefore and NotAfter are parsed from the certificate, so they will differ
	// from the original values. We just verify they're not zero.
	if loaded.NotBefore.IsZero() {
		t.Error("NotBefore should not be zero after loading")
	}
	if loaded.NotAfter.IsZero() {
		t.Error("NotAfter should not be zero after loading")
	}
}

// TestAttestationResult_SaveToDirectory_EmptyCAChain tests that saving works
// when CA chain is empty (which is valid).
func TestAttestationResult_SaveToDirectory_EmptyCAChain(t *testing.T) {
	// Arrange: Create result with empty CA chain
	testCertPEM := generateTestCertificate(t)

	result := &AttestationResult{
		CertPEM:    testCertPEM,
		CAChainPEM: nil, // Empty CA chain
		KeyMaterial: &TPMKeyMaterial{
			Type:            "tpm2",
			AKBlob:          []byte("test-ak"),
			AppBlob:         []byte("test-app"),
			PublicKeySHA256: "abc123",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour),
	}

	tmpDir := t.TempDir()

	// Act: Save to directory
	err := result.SaveToDirectory(tmpDir, 0o022)
	if err != nil {
		t.Fatalf("SaveToDirectory failed: %v", err)
	}

	// Assert: CA chain file should not exist
	caChainPath := filepath.Join(tmpDir, "ca_chain.pem")
	if _, err := os.Stat(caChainPath); !os.IsNotExist(err) {
		t.Errorf("ca_chain.pem should not exist when CAChainPEM is empty")
	}

	// Act: Load should still work
	loaded, err := LoadFromDirectory(tmpDir)
	if err != nil {
		t.Fatalf("LoadFromDirectory failed: %v", err)
	}

	// Assert: Loaded CA chain should be empty
	if len(loaded.CAChainPEM) != 0 {
		t.Errorf("Expected empty CAChainPEM, got %d bytes", len(loaded.CAChainPEM))
	}
}

// TestAttestationResult_SaveToDirectory_NilResult tests error handling for nil result.
func TestAttestationResult_SaveToDirectory_NilResult(t *testing.T) {
	// Arrange: nil result
	var result *AttestationResult

	tmpDir := t.TempDir()

	// Act & Assert: Should return error
	err := result.SaveToDirectory(tmpDir, 0o022)
	if err == nil {
		t.Fatal("Expected error for nil AttestationResult, got nil")
	}
	if err.Error() != "attestation result is nil" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// TestAttestationResult_SaveToDirectory_NonExistentDirectory tests error handling
// when destination directory doesn't exist.
func TestAttestationResult_SaveToDirectory_NonExistentDirectory(t *testing.T) {
	// Arrange: Valid result but non-existent directory
	result := &AttestationResult{
		CertPEM: generateTestCertificate(t),
		KeyMaterial: &TPMKeyMaterial{
			Type:            "tpm2",
			AKBlob:          []byte("test-ak"),
			AppBlob:         []byte("test-app"),
			PublicKeySHA256: "abc123",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour),
	}

	nonExistentDir := filepath.Join(t.TempDir(), "does-not-exist")

	// Act & Assert: Should return error
	err := result.SaveToDirectory(nonExistentDir, 0o022)
	if err == nil {
		t.Fatal("Expected error for non-existent directory, got nil")
	}
}

// TestAttestationResult_SaveToDirectory_FileAsDirectory tests error handling
// when destination path is a file, not a directory.
func TestAttestationResult_SaveToDirectory_FileAsDirectory(t *testing.T) {
	// Arrange: Create a file instead of directory
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "not-a-directory")
	if err := os.WriteFile(filePath, []byte("test"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result := &AttestationResult{
		CertPEM: generateTestCertificate(t),
		KeyMaterial: &TPMKeyMaterial{
			Type:            "tpm2",
			AKBlob:          []byte("test-ak"),
			AppBlob:         []byte("test-app"),
			PublicKeySHA256: "abc123",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour),
	}

	// Act & Assert: Should return error
	err := result.SaveToDirectory(filePath, 0o022)
	if err == nil {
		t.Fatal("Expected error when destination is a file, got nil")
	}
	if err.Error() != "destination path must be a directory" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// TestLoadFromDirectory_MissingCertificate tests error handling when certificate is missing.
func TestLoadFromDirectory_MissingCertificate(t *testing.T) {
	// Arrange: Directory with only key reference, no certificate
	tmpDir := t.TempDir()

	// Create minimal key reference file
	keyRefPath := filepath.Join(tmpDir, "client-key.json")
	keyRefJSON := []byte(`{
		"type": "tpm2",
		"version": 1,
		"public_key_sha256": "abc123"
	}`)
	if err := os.WriteFile(keyRefPath, keyRefJSON, 0o644); err != nil {
		t.Fatalf("Failed to create key reference: %v", err)
	}

	// Act & Assert: Should return error about missing certificate
	_, err := LoadFromDirectory(tmpDir)
	if err == nil {
		t.Fatal("Expected error for missing certificate, got nil")
	}
}

// TestLoadFromDirectory_MissingKeyReference tests error handling when key reference is missing.
func TestLoadFromDirectory_MissingKeyReference(t *testing.T) {
	// Arrange: Directory with only certificate, no key reference
	tmpDir := t.TempDir()

	certPath := filepath.Join(tmpDir, "client.crt")
	certPEM := generateTestCertificate(t)
	if err := os.WriteFile(certPath, certPEM, 0o644); err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	// Act & Assert: Should return error about missing key reference
	_, err := LoadFromDirectory(tmpDir)
	require.ErrorContains(t, err, "failed to read key reference")
}

// TestAttestationResult_SaveToDirectory_FilePermissions tests that file permissions
// are set correctly with umask applied.
func TestAttestationResult_SaveToDirectory_FilePermissions(t *testing.T) {
	// Arrange: Create result
	result := &AttestationResult{
		CertPEM:    generateTestCertificate(t),
		CAChainPEM: generateTestCertificate(t),
		KeyMaterial: &TPMKeyMaterial{
			Type:            "tpm2",
			AKBlob:          []byte("test-ak"),
			AppBlob:         []byte("test-app"),
			PublicKeySHA256: "abc123",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour),
	}

	tmpDir := t.TempDir()

	// Act: Save with umask 0o022
	err := result.SaveToDirectory(tmpDir, 0o022)
	if err != nil {
		t.Fatalf("SaveToDirectory failed: %v", err)
	}

	// Assert: Check private files have restrictive permissions
	privateFiles := []string{"client-key.json", "ak.blob", "app.blob"}
	for _, filename := range privateFiles {
		path := filepath.Join(tmpDir, filename)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("Failed to stat %s: %v", filename, err)
			continue
		}

		// With umask 0o022, 0o600 becomes 0o600 (no change for owner-only)
		expectedPerm := os.FileMode(0o600)
		if info.Mode().Perm() != expectedPerm {
			t.Errorf("File %s has permissions %o, expected %o", filename, info.Mode().Perm(), expectedPerm)
		}
	}

	// Assert: Check public files have appropriate permissions
	publicFiles := []string{"ca_chain.pem", "client.crt"}
	for _, filename := range publicFiles {
		path := filepath.Join(tmpDir, filename)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("Failed to stat %s: %v", filename, err)
			continue
		}

		// With umask 0o022, 0o644 becomes 0o644
		expectedPerm := os.FileMode(0o644)
		if info.Mode().Perm() != expectedPerm {
			t.Errorf("File %s has permissions %o, expected %o", filename, info.Mode().Perm(), expectedPerm)
		}
	}
}
