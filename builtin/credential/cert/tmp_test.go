// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cert

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestIsTSS2Key(t *testing.T) {
	tests := []struct {
		name     string
		keyData  []byte
		expected bool
	}{
		{
			name: "Valid TSS2 key",
			keyData: []byte(`-----BEGIN TSS2 PRIVATE KEY-----
VGhpcyBpcyBhIGZha2UgVFNTMiBwcml2YXRlIGtleSBkYXRhIGZvciB0ZXN0aW5n
-----END TSS2 PRIVATE KEY-----`),
			expected: true,
		},
		{
			name: "Standard RSA key",
			keyData: []byte(`-----BEGIN RSA PRIVATE KEY-----
VGhpcyBpcyBhIGZha2UgUlNBIHByaXZhdGUga2V5IGRhdGEgZm9yIHRlc3Rpbmc=
-----END RSA PRIVATE KEY-----`),
			expected: false,
		},
		{
			name: "Standard EC key",
			keyData: []byte(`-----BEGIN EC PRIVATE KEY-----
VGhpcyBpcyBhIGZha2UgRUMgcHJpdmF0ZSBrZXkgZGF0YSBmb3IgdGVzdGluZw==
-----END EC PRIVATE KEY-----`),
			expected: false,
		},
		{
			name:     "Invalid PEM",
			keyData:  []byte("not a pem file"),
			expected: false,
		},
		{
			name:     "Empty data",
			keyData:  []byte(""),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTSS2Key(tt.keyData)
			if result != tt.expected {
				t.Errorf("isTSS2Key() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestTPMSupport(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("TPM support requires Linux")
	}

	// Check if running in CI environment
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("TPM tests skipped in CI environment (no TPM hardware available)")
	}

	// Check if TPM device exists
	if _, err := os.Stat("/dev/tpmrm0"); os.IsNotExist(err) {
		t.Skip("TPM hardware not available (/dev/tpmrm0 not found)")
	}

	t.Run("TPM supported on Linux", func(t *testing.T) {
		if !tmpSupported() {
			t.Error("tmpSupported() should return true on Linux")
		}
	})

	t.Run("Get TPM device path", func(t *testing.T) {
		// Test default path
		path := getTPMDevicePath("")
		if path != "/dev/tpmrm0" {
			t.Errorf("getTPMDevicePath(\"\") = %s, expected /dev/tpmrm0", path)
		}

		// Test custom path
		customPath := "/dev/tpm1"
		path = getTPMDevicePath(customPath)
		if path != customPath {
			t.Errorf("getTPMDevicePath(%s) = %s, expected %s", customPath, path, customPath)
		}
	})
}

func TestTPMSupportNonLinux(t *testing.T) {
	if runtime.GOOS == "linux" {
		t.Skip("This test is for non-Linux platforms")
	}

	t.Run("TPM not supported on non-Linux", func(t *testing.T) {
		if tmpSupported() {
			t.Error("tmpSupported() should return false on non-Linux platforms")
		}
	})

	t.Run("TPM device path empty on non-Linux", func(t *testing.T) {
		path := getTPMDevicePath("test")
		if path != "" {
			t.Errorf("getTPMDevicePath() should return empty string on non-Linux, got %s", path)
		}
	})
}

func TestLoadStandardPrivateKey(t *testing.T) {
	tests := []struct {
		name        string
		keyData     []byte
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Invalid PEM",
			keyData:     []byte("not a pem file"),
			expectError: true,
			errorMsg:    "failed to parse PEM block",
		},
		{
			name: "Unsupported key type",
			keyData: []byte(`-----BEGIN UNKNOWN KEY-----
VGhpcyBpcyBhIGZha2UgdW5rbm93biBrZXkgdHlwZSBmb3IgdGVzdGluZw==
-----END UNKNOWN KEY-----`),
			expectError: true,
			errorMsg:    "unsupported private key type",
		},
		{
			name:        "Empty data",
			keyData:     []byte(""),
			expectError: true,
			errorMsg:    "failed to parse PEM block",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadStandardPrivateKey(tt.keyData)
			if tt.expectError {
				if err == nil {
					t.Error("Expected an error, but got none")
					return
				}
				if tt.errorMsg != "" && !strings.HasPrefix(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to start with %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}
