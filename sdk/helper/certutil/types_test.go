// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package certutil

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
)

func TestGetPrivateKeyTypeFromPublicKey(t *testing.T) {
	rsaKey, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("error generating rsa key: %s", err)
	}

	ecdsaKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		t.Fatalf("error generating ecdsa key: %s", err)
	}

	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("error generating ed25519 key: %s", err)
	}

	testCases := map[string]struct {
		publicKey       crypto.PublicKey
		expectedKeyType PrivateKeyType
	}{
		"rsa": {
			publicKey:       rsaKey.Public(),
			expectedKeyType: RSAPrivateKey,
		},
		"ecdsa": {
			publicKey:       ecdsaKey.Public(),
			expectedKeyType: ECPrivateKey,
		},
		"ed25519": {
			publicKey:       publicKey,
			expectedKeyType: Ed25519PrivateKey,
		},
		"bad key type": {
			publicKey:       []byte{},
			expectedKeyType: UnknownPrivateKey,
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			keyType := GetPrivateKeyTypeFromPublicKey(tt.publicKey)

			if keyType != tt.expectedKeyType {
				t.Fatalf("key type mismatch: expected %s, got %s", tt.expectedKeyType, keyType)
			}
		})
	}
}

// TestCertExtKeyUsageIsPresent validates the expected behavior of the CertExtKeyUsageIsPresent
// function, also serves as documenting our expected behaviors in different edge case inputs.
func TestCertExtKeyUsageIsPresent(t *testing.T) {
	type args struct {
		desiredExtKeyUsage CertExtKeyUsage
		extKeyUsages       CertExtKeyUsage
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{desiredExtKeyUsage: ServerAuthExtKeyUsage, extKeyUsages: 0}, false},
		{"present-one-to-one", args{desiredExtKeyUsage: ServerAuthExtKeyUsage, extKeyUsages: ServerAuthExtKeyUsage}, true},
		{"present-one-to-many", args{desiredExtKeyUsage: ServerAuthExtKeyUsage, extKeyUsages: ServerAuthExtKeyUsage | ClientAuthExtKeyUsage}, true},
		{"missing-one-to-one", args{desiredExtKeyUsage: ServerAuthExtKeyUsage, extKeyUsages: EmailProtectionExtKeyUsage}, false},
		{"missing-to-many", args{desiredExtKeyUsage: ServerAuthExtKeyUsage, extKeyUsages: EmailProtectionExtKeyUsage | ClientAuthExtKeyUsage}, false},
		// Don't treat AnyExtUsage as a special value
		{"any-ext-usage-no-match", args{desiredExtKeyUsage: ServerAuthExtKeyUsage, extKeyUsages: AnyExtKeyUsage}, false},
		{"any-ext-usage-matches-itself", args{desiredExtKeyUsage: AnyExtKeyUsage, extKeyUsages: AnyExtKeyUsage}, true},
		// The desiredExtKeyUsage should be a single ExtKeyUsage not multiple
		{"bad-usage-many-to-many-missing", args{desiredExtKeyUsage: ServerAuthExtKeyUsage | ClientAuthExtKeyUsage, extKeyUsages: ServerAuthExtKeyUsage | EmailProtectionExtKeyUsage}, false},
		{"bad-usage-many-to-many-direct", args{desiredExtKeyUsage: ServerAuthExtKeyUsage | ClientAuthExtKeyUsage, extKeyUsages: ServerAuthExtKeyUsage | ClientAuthExtKeyUsage}, false},
		{"bad-usage-many-to-many-with-additional", args{desiredExtKeyUsage: ServerAuthExtKeyUsage | ClientAuthExtKeyUsage, extKeyUsages: ServerAuthExtKeyUsage | ClientAuthExtKeyUsage | EmailProtectionExtKeyUsage}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CertExtKeyUsageIsPresent(tt.args.extKeyUsages, tt.args.desiredExtKeyUsage); got != tt.want {
				t.Errorf("IsExtKeyUsagePresent() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetCertEKUString Validates we get the expected string value back for each of our defined EKUs
func TestGetCertEKUString(t *testing.T) {
	tests := []struct {
		name   string
		eku    CertExtKeyUsage
		want   string
		wantOk bool
	}{
		{"unknown", CertExtKeyUsage(0), "unknown EKU 0", false},
		{"any", AnyExtKeyUsage, "any", true},
		{"serverauth", ServerAuthExtKeyUsage, "serverAuth", true},
		{"clientauth", ClientAuthExtKeyUsage, "clientAuth", true},
		{"codeSigning", CodeSigningExtKeyUsage, "codeSigning", true},
		{"email", EmailProtectionExtKeyUsage, "emailProtection", true},
		{"ipsec-1", IpsecEndSystemExtKeyUsage, "ipsecEndSystem", true},
		{"ipsec-2", IpsecTunnelExtKeyUsage, "ipsecTunnel", true},
		{"ipsec-3", IpsecUserExtKeyUsage, "ipsecUser", true},
		{"timestamp", TimeStampingExtKeyUsage, "timeStamping", true},
		{"ocsp", OcspSigningExtKeyUsage, "ocspSigning", true},
		{"msgc", MicrosoftServerGatedCryptoExtKeyUsage, "microsoftServerGatedCrypto", true},
		{"nsgc", NetscapeServerGatedCryptoExtKeyUsage, "netscapeServerGatedCrypto", true},
		{"mscommercialcs", MicrosoftCommercialCodeSigningExtKeyUsage, "microsoftCommercialCodeSigning", true},
		{"mskernel", MicrosoftKernelCodeSigningExtKeyUsage, "microsoftKernelCodeSigning", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetCertEKUString(tt.eku)
			if got != tt.want {
				t.Errorf("GetCertEKUString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantOk {
				t.Errorf("GetCertEKUString() got1 = %v, want %v", got1, tt.wantOk)
			}
		})
	}
}

// TestGetCertEKUFromString validates we can convert string values to the EKU we want
func TestGetCertEKUFromString(t *testing.T) {
	tests := []struct {
		name   string
		eku    string
		want   CertExtKeyUsage
		wantOk bool
	}{
		{"empty", "", CertExtKeyUsage(0), false},
		{"unknown", "invalid-eku", CertExtKeyUsage(0), false},
		{"any", "any", AnyExtKeyUsage, true},
		{"any-with-spaces", "   any    ", AnyExtKeyUsage, true},
		{"any-with-caps", "AnY", AnyExtKeyUsage, true},
		{"any-with-eku-suffix", "AnYExtKeyUsage", AnyExtKeyUsage, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetCertEKUFromString(tt.eku)
			if got != tt.want {
				t.Errorf("GetCertEKUFromString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantOk {
				t.Errorf("GetCertEKUFromString() got1 = %v, want %v", got1, tt.wantOk)
			}
		})
	}
}
