// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestExportCrypto_ReversibleUnwrap validates that CA private keys wrapped with RSA, EC, and ML-KEM
// can be unwrapped using unwrapCAPrivateKey back to the original plaintext PEM.
// A real EC P-256 private key is used as the CA key plaintext so the round-trip
// exercises actual key material rather than an invalid DER payload.
func TestExportCrypto_ReversibleUnwrap(t *testing.T) {
	t.Parallel()

	// Generate a real EC P-256 private key to use as the CA key being wrapped.
	caPriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	caDER, err := x509.MarshalPKCS8PrivateKey(caPriv)
	require.NoError(t, err)
	caPrivKeyPEM := strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: caDER})))

	keyTypes := []string{
		"rsa-2048",
		"rsa-3072",
		"rsa-4096",
		"ec-p256",
		"ec-p384",
		"ec-p521",
		"ml-kem-768",
		"ml-kem-1024",
	}

	for _, kt := range keyTypes {
		kt := kt
		t.Run(kt, func(t *testing.T) {
			t.Parallel()

			privPEM, pubPEM, _, err := generateExportKeypair(kt)
			require.NoError(t, err)

			wrappedJSON, err := wrapCAPrivateKey(caPrivKeyPEM, pubPEM)
			require.NoError(t, err)

			unwrapped, err := unwrapCAPrivateKey(wrappedJSON, privPEM)
			require.NoError(t, err)
			require.Equal(t, caPrivKeyPEM, unwrapped)
		})
	}
}
