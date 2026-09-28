// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

//go:build enterprise

package pki

import (
	"context"
	"crypto/mldsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestGenerateRoot_internal(t *testing.T) {
	t.Parallel()

	parameterSets := []struct {
		paramSet      string
		wantMLDSAAlgo mldsa.Parameters
	}{
		{certutil.MLDSA44, mldsa.MLDSA44()},
		{certutil.MLDSA65, mldsa.MLDSA65()},
		{certutil.MLDSA87, mldsa.MLDSA87()},
	}

	for _, tc := range parameterSets {
		t.Run("ps"+tc.paramSet, func(t *testing.T) {
			t.Parallel()
			b, s := CreateBackendWithStorage(t)

			resp, err := b.HandleRequest(context.Background(), &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "root/generate/internal",
				Storage:   s,
				Data: map[string]interface{}{
					"common_name":   "Test ML-DSA Root CA",
					"key_type":      "ml-dsa",
					"parameter_set": tc.paramSet,
					"ttl":           "8760h",
				},
				MountPoint: "pki/",
			})
			require.NoError(t, err, "unexpected transport error generating ML-DSA-%s root", tc.paramSet)

			validateResponse(t, resp)

			// Private key must not appear for internal generation
			require.Nil(t, resp.Data["private_key"],
				"private_key must not be present in internal root generation response")
		})
	}
}

func TestGenerateRoot_exported(t *testing.T) {
	t.Parallel()

	parameterSets := []struct {
		paramSet      string
		wantMLDSAAlgo mldsa.Parameters
	}{
		{certutil.MLDSA44, mldsa.MLDSA44()},
		{certutil.MLDSA65, mldsa.MLDSA65()},
		{certutil.MLDSA87, mldsa.MLDSA87()},
	}

	for _, tc := range parameterSets {
		t.Run("mldsa-"+tc.paramSet, func(t *testing.T) {
			t.Parallel()
			b, s := CreateBackendWithStorage(t)

			resp, err := b.HandleRequest(context.Background(), &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "root/generate/exported",
				Storage:   s,
				Data: map[string]interface{}{
					"common_name":   "Test ML-DSA Root CA",
					"key_type":      "ml-dsa",
					"parameter_set": tc.paramSet,
					"ttl":           "8760h",
				},
				MountPoint: "pki/",
			})

			require.NoError(t, err, "unexpected transport error generating exported ML-DSA-%s root", tc.paramSet)

			validateResponse(t, resp)

			// Private key must be present for exported generation
			keyPEM, ok := resp.Data["private_key"].(string)
			require.True(t, ok, "private_key field should be a string")
			require.NotEmpty(t, keyPEM, "private_key must not be empty for exported generation")

			require.Equal(t, certutil.MLDSAPrivateKey, resp.Data["private_key_type"],
				"key_type should be ml-dsa for ML-DSA-%s root", tc.paramSet)
			require.Equal(t, certutil.ParameterSet(tc.paramSet), resp.Data["parameter_set"],
				"parameter_set should match requested value for ML-DSA-%s root", tc.paramSet)

			// Decode and verify the ML-DSA private key
			block, rest := pem.Decode([]byte(keyPEM))
			require.Empty(t, rest, "trailing data after private key PEM block")
			require.NotNil(t, block, "failed to PEM-decode private key")
			rawKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			require.NoError(t, err, "failed to parse PKCS#8 private key")
			mldsaKey, ok := rawKey.(*mldsa.PrivateKey)
			require.True(t, ok, "expected *mldsa.PrivateKey, got %T", rawKey)
			require.Equal(t, tc.wantMLDSAAlgo, mldsaKey.PublicKey().Parameters(),
				"exported private key parameters do not match requested parameter set")
		})
	}
}

func TestGenerateRoot_existing(t *testing.T) {
	t.Parallel()

	mldsa44Bundle, err := certutil.CreateKeyBundle("ml-dsa", 0, rand.Reader, "44")
	require.NoError(t, err, "failed generating an ml-dsa key bundle")
	mldsa44Pem, err := mldsa44Bundle.ToPrivateKeyPemString()
	require.NoError(t, err, "failed converting ml-dsa key to pem")

	mldsa65Bundle, err := certutil.CreateKeyBundle("ml-dsa", 0, rand.Reader, "65")
	require.NoError(t, err, "failed generating an ml-dsa key bundle")
	mldsa65Pem, err := mldsa65Bundle.ToPrivateKeyPemString()
	require.NoError(t, err, "failed converting ml-dsa key to pem")

	mldsa87Bundle, err := certutil.CreateKeyBundle("ml-dsa", 0, rand.Reader, "87")
	require.NoError(t, err, "failed generating an ml-dsa key bundle")
	mldsa87Pem, err := mldsa87Bundle.ToPrivateKeyPemString()
	require.NoError(t, err, "failed converting ml-dsa key to pem")

	testCases := map[string]struct {
		parameterSet string
		key          string
	}{
		"mldsa-44": {
			parameterSet: certutil.MLDSA44,
			key:          mldsa44Pem,
		},
		"mldsa-65": {
			parameterSet: certutil.MLDSA65,
			key:          mldsa65Pem,
		},
		"mldsa-87": {
			parameterSet: certutil.MLDSA87,
			key:          mldsa87Pem,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			b, s := CreateBackendWithStorage(t)

			resp, err := b.HandleRequest(context.Background(), &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "keys/import",
				Storage:   s,
				Data: map[string]interface{}{
					"key_name":   name,
					"pem_bundle": tc.key,
				},
				MountPoint: "pki/",
			})
			require.NoError(t, err)
			require.False(t, resp.IsError(), "got unexpected error response: %s", resp.Error())

			resp, err = b.HandleRequest(context.Background(), &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "root/generate/existing",
				Storage:   s,
				Data: map[string]interface{}{
					"common_name": "Test ML-DSA Root CA",
					"ttl":         "8760h",
				},
				MountPoint: "pki/",
			})
			require.NoError(t, err, "unexpected transport error generating exported ML-DSA-%s root", tc.parameterSet)

			validateResponse(t, resp)
		})
	}
}

func validateResponse(t *testing.T, resp *logical.Response) {
	require.NotNil(t, resp, "got nil response generating exported ML-DSA root")
	require.False(t, resp.IsError(), "unexpected logical error generating exported ML-DSA root: %v", resp.Error())

	// Key metadata
	require.NotEmpty(t, resp.Data["key_id"], "key_id must not be empty")
	require.NotEmpty(t, resp.Data["issuer_id"], "issuer_id must not be empty")

	// Certificate is parseable and uses ML-DSA
	certPEM, ok := resp.Data["certificate"].(string)
	require.True(t, ok, "certificate field should be a string")
	require.NotEmpty(t, certPEM, "certificate must not be empty")

	block, _ := pem.Decode([]byte(certPEM))
	require.NotNil(t, block, "failed to PEM-decode certificate")
	cert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err, "failed to parse generated certificate")
	require.Equal(t, x509.MLDSA, cert.PublicKeyAlgorithm,
		"certificate public key algorithm should be MLDSA")
	require.True(t, cert.IsCA, "generated root certificate must have IsCA=true")
}
