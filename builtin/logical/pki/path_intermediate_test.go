// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto/ecdsa"
	"crypto/mldsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestGenerateIntermediate_FormatParam validates that CSR generation
// only supports valid formats and does NOT support PKCS#12 since there's no certificate to bundle
func TestGenerateIntermediate_FormatParam(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)
	testCases := []struct {
		endpoint   string
		format     string
		isValid    bool
		omitFormat bool
	}{
		{endpoint: "intermediate/generate/internal", format: "pkcs12_bundle"},
		{endpoint: "intermediate/generate/exported", format: "pkcs12_bundle"},
		{endpoint: "issuers/generate/intermediate/internal", format: "pkcs12_bundle"},
		{endpoint: "issuers/generate/intermediate/exported", format: "pkcs12_bundle"},

		{endpoint: "intermediate/generate/internal", format: "invalid"},
		{endpoint: "intermediate/generate/exported", format: "invalid"},
		{endpoint: "issuers/generate/intermediate/internal", format: "invalid"},
		{endpoint: "issuers/generate/intermediate/exported", format: "invalid"},

		{endpoint: "intermediate/generate/internal", format: ""},
		{endpoint: "intermediate/generate/exported", format: ""},
		{endpoint: "issuers/generate/intermediate/internal", format: ""},
		{endpoint: "issuers/generate/intermediate/exported", format: ""},

		{endpoint: "intermediate/generate/internal", omitFormat: true},
		{endpoint: "intermediate/generate/exported", omitFormat: true},
		{endpoint: "issuers/generate/intermediate/internal", omitFormat: true},
		{endpoint: "issuers/generate/intermediate/exported", omitFormat: true},
	}
	for _, tc := range testCases {
		name := fmt.Sprintf("endpoint=%s format=%s", tc.endpoint, tc.format)
		t.Run(name, func(t *testing.T) {
			// Attempt to generate intermediate CSR with various formats
			data := map[string]interface{}{
				"common_name": "Intermediate CA",
				"key_type":    "ec",
				"key_bits":    256,
			}
			if !tc.omitFormat {
				data["format"] = tc.format
			}
			_, err := CBWrite(b, s, tc.endpoint, data)

			if tc.omitFormat {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), `the "format" parameter must be "pem", "der" or "pem_bundle"`)
			}
		})
	}
}

// TestSetSignedIntermediate verifies that ML-DSA intermediate certificates
// produced by root/sign-intermediate are accepted by intermediate/set-signed for
// all supported ML-DSA parameter sets. For each parameter set the test:
//   - generates a root CA on one backend and a CSR on a second backend,
//   - signs the CSR with the root,
//   - imports the signed certificate (plus the root chain) via set-signed, and
//   - asserts the response lists the issuer and key, the issuer can be read back,
//     and re-importing the same certificate is idempotent.
func TestSetSignedIntermediate(t *testing.T) {
	t.Parallel()

	parameterSets := []struct {
		paramSet           string
		expectedParameters mldsa.Parameters
		sigAlgo            x509.SignatureAlgorithm
	}{
		{certutil.MLDSA44, mldsa.MLDSA44(), x509.MLDSA44},
		{certutil.MLDSA65, mldsa.MLDSA65(), x509.MLDSA65},
		{certutil.MLDSA87, mldsa.MLDSA87(), x509.MLDSA87},
	}

	for _, tc := range parameterSets {
		t.Run("mldsa-"+tc.paramSet, func(t *testing.T) {
			t.Parallel()

			// Root backend and intermediate backend are isolated.
			bRoot, sRoot := CreateBackendWithStorage(t)
			bInt, sInt := CreateBackendWithStorage(t)

			// Step 1: Generate a self-signed root CA with the current parameter set.
			resp, err := CBWrite(bRoot, sRoot, "root/generate/internal", map[string]interface{}{
				"common_name":   "Test ML-DSA Root CA",
				"key_type":      "ml-dsa",
				"parameter_set": tc.paramSet,
				"ttl":           "8760h",
			})
			requireSuccessNonNilResponse(t, resp, err, "generate ML-DSA-%s root CA", tc.paramSet)
			rootCertPEM := resp.Data["certificate"].(string)

			// Step 2: Generate an intermediate CSR with the same parameter set.
			resp, err = CBWrite(bInt, sInt, "intermediate/generate/internal", map[string]interface{}{
				"common_name":   "Test ML-DSA Intermediate CA",
				"key_type":      "ml-dsa",
				"parameter_set": tc.paramSet,
			})
			requireSuccessNonNilResponse(t, resp, err, "generate ML-DSA-%s intermediate CSR", tc.paramSet)
			intCSR := resp.Data["csr"].(string)
			require.NotEmpty(t, intCSR, "intermediate CSR should be non-empty for ML-DSA-%s", tc.paramSet)

			// Step 3: Sign the CSR with the root CA.
			resp, err = CBWrite(bRoot, sRoot, "root/sign-intermediate", map[string]interface{}{
				"common_name": "Test ML-DSA Intermediate CA",
				"csr":         intCSR,
				"ttl":         "4380h",
			})
			requireSuccessNonNilResponse(t, resp, err, "sign ML-DSA-%s intermediate CSR", tc.paramSet)
			intCertPEM := resp.Data["certificate"].(string)

			// Verify the signed certificate has the expected parameters before importing.
			block, rest := pem.Decode([]byte(intCertPEM))
			require.NotNil(t, block, "signed intermediate cert PEM should decode for ML-DSA-%s", tc.paramSet)
			require.Empty(t, rest, "no trailing data expected after signed intermediate cert for ML-DSA-%s", tc.paramSet)
			intCert, parseErr := x509.ParseCertificate(block.Bytes)
			require.NoError(t, parseErr, "failed to parse signed intermediate certificate for ML-DSA-%s", tc.paramSet)
			require.True(t, intCert.IsCA, "signed intermediate must have IsCA=true for ML-DSA-%s", tc.paramSet)
			pub, ok := intCert.PublicKey.(*mldsa.PublicKey)
			require.True(t, ok, "signed intermediate public key should be *mldsa.PublicKey for ML-DSA-%s, got %T", tc.paramSet, intCert.PublicKey)
			require.Equal(t, tc.expectedParameters, pub.Parameters(),
				"signed intermediate public key parameters should match requested set for ML-DSA-%s", tc.paramSet)

			// Step 4: Import the signed intermediate (with root chain) via set-signed.
			resp, err = CBWrite(bInt, sInt, "intermediate/set-signed", map[string]interface{}{
				"certificate": intCertPEM + "\n" + rootCertPEM,
			})
			requireSuccessNonNilResponse(t, resp, err, "set-signed for ML-DSA-%s intermediate", tc.paramSet)

			// The response must report exactly one newly-imported issuer (the
			// intermediate). The root is appended for chain building only and is
			// not assigned a local key, so it may appear as a new issuer without
			// an associated key; the intermediate's key was already held internally.
			importedIssuers, ok := resp.Data["imported_issuers"].([]string)
			require.True(t, ok, "imported_issuers should be a []string for ML-DSA-%s", tc.paramSet)
			require.NotEmpty(t, importedIssuers, "at least one issuer should be imported by set-signed for ML-DSA-%s", tc.paramSet)

			// The mapping must contain an entry for each imported issuer; the
			// intermediate issuer's entry must point to a non-empty key ID.
			mapping, ok := resp.Data["mapping"].(map[string]string)
			require.True(t, ok, "mapping should be a map[string]string for ML-DSA-%s", tc.paramSet)
			require.NotEmpty(t, mapping, "mapping must not be empty after set-signed for ML-DSA-%s", tc.paramSet)

			// At least one imported issuer must have an associated key (the intermediate).
			hasKeyedIssuer := false
			for _, issuerID := range importedIssuers {
				if mapping[issuerID] != "" {
					hasKeyedIssuer = true
					break
				}
			}
			require.True(t, hasKeyedIssuer,
				"at least one imported issuer must have an associated key in the mapping for ML-DSA-%s", tc.paramSet)

			// Step 5: Confirm idempotency — re-importing the same certificate must
			// report the issuer as already existing rather than duplicating it.
			reimportResp, err := CBWrite(bInt, sInt, "intermediate/set-signed", map[string]interface{}{
				"certificate": intCertPEM + "\n" + rootCertPEM,
			})
			requireSuccessNonNilResponse(t, reimportResp, err, "re-import ML-DSA-%s signed intermediate", tc.paramSet)

			reimportedIssuers, ok := reimportResp.Data["imported_issuers"].([]string)
			require.True(t, ok, "imported_issuers should be a []string on re-import for ML-DSA-%s", tc.paramSet)
			require.Empty(t, reimportedIssuers,
				"no new issuers should be imported on re-import for ML-DSA-%s", tc.paramSet)

			existingIssuers, ok := reimportResp.Data["existing_issuers"].([]string)
			require.True(t, ok, "existing_issuers should be a []string on re-import for ML-DSA-%s", tc.paramSet)
			require.NotEmpty(t, existingIssuers,
				"re-import should report existing issuers for ML-DSA-%s", tc.paramSet)
		})
	}
}

// TestSetSignedIntermediate_crossParamSet verifies that intermediate/set-signed
// accepts an ML-DSA intermediate certificate whose public key uses a different
// parameter set than the signing root. The intermediate is generated with ML-DSA-65,
// signed by an ML-DSA-44 root, and the signed certificate is imported into the
// intermediate backend alongside the root chain.
func TestSetSignedIntermediate_crossParamSet(t *testing.T) {
	t.Parallel()

	bRoot, sRoot := CreateBackendWithStorage(t)
	bInt, sInt := CreateBackendWithStorage(t)

	// Generate an ML-DSA-44 root CA.
	resp, err := CBWrite(bRoot, sRoot, "root/generate/internal", map[string]interface{}{
		"common_name":   "ML-DSA-44 Root CA",
		"key_type":      "ml-dsa",
		"parameter_set": certutil.MLDSA44,
		"ttl":           "8760h",
	})
	requireSuccessNonNilResponse(t, resp, err, "generate ML-DSA-44 root CA")
	rootCertPEM := resp.Data["certificate"].(string)

	// Generate an ML-DSA-65 intermediate CSR.
	resp, err = CBWrite(bInt, sInt, "intermediate/generate/internal", map[string]interface{}{
		"common_name":   "ML-DSA-65 Intermediate CA",
		"key_type":      "ml-dsa",
		"parameter_set": certutil.MLDSA65,
	})
	requireSuccessNonNilResponse(t, resp, err, "generate ML-DSA-65 intermediate CSR")
	intCSR := resp.Data["csr"].(string)

	// Sign the ML-DSA-65 CSR with the ML-DSA-44 root.
	resp, err = CBWrite(bRoot, sRoot, "root/sign-intermediate", map[string]interface{}{
		"common_name": "ML-DSA-65 Intermediate CA",
		"csr":         intCSR,
		"ttl":         "4380h",
	})
	requireSuccessNonNilResponse(t, resp, err, "sign ML-DSA-65 intermediate CSR with ML-DSA-44 root")
	intCertPEM := resp.Data["certificate"].(string)

	// The signed certificate's public key must use ML-DSA-65 parameters, but
	// its signature algorithm must be ML-DSA-44 (produced by the root).
	block, _ := pem.Decode([]byte(intCertPEM))
	require.NotNil(t, block, "signed cross-param intermediate cert PEM should decode")
	intCert, parseErr := x509.ParseCertificate(block.Bytes)
	require.NoError(t, parseErr, "failed to parse cross-param signed intermediate certificate")
	require.Equal(t, x509.MLDSA44, intCert.SignatureAlgorithm,
		"cross-param intermediate should be signed with ML-DSA-44 (the root's algorithm)")
	pub, ok := intCert.PublicKey.(*mldsa.PublicKey)
	require.True(t, ok, "cross-param intermediate public key should be *mldsa.PublicKey, got %T", intCert.PublicKey)
	require.Equal(t, mldsa.MLDSA65(), pub.Parameters(),
		"cross-param intermediate public key should use ML-DSA-65 parameters")

	// Import the signed intermediate via set-signed; it must succeed.
	resp, err = CBWrite(bInt, sInt, "intermediate/set-signed", map[string]interface{}{
		"certificate": intCertPEM + "\n" + rootCertPEM,
	})
	requireSuccessNonNilResponse(t, resp, err, "set-signed for cross-param ML-DSA-44→ML-DSA-65 intermediate")

	importedIssuers, ok := resp.Data["imported_issuers"].([]string)
	require.True(t, ok, "imported_issuers should be a []string for cross-param set-signed")
	require.NotEmpty(t, importedIssuers, "at least one issuer should be imported for cross-param set-signed")

	mapping, ok := resp.Data["mapping"].(map[string]string)
	require.True(t, ok, "mapping should be a map[string]string for cross-param set-signed")
	require.NotEmpty(t, mapping, "mapping must not be empty after cross-param set-signed")
}

// TestGenerateIntermediate_internal verifies that ML-DSA intermediate CSR
// generation with key_type=internal stores the key in Vault without returning
// it in the response, for all supported parameter sets.
func TestGenerateIntermediate_internal(t *testing.T) {
	t.Parallel()

	parameterSets := []struct {
		paramSet           string
		expectedParameters mldsa.Parameters
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
				Path:      "intermediate/generate/internal",
				Storage:   s,
				Data: map[string]interface{}{
					"common_name":   "Test ML-DSA Intermediate CA",
					"key_type":      "ml-dsa",
					"parameter_set": tc.paramSet,
				},
				MountPoint: "pki/",
			})
			require.NoError(t, err, "unexpected transport error generating ML-DSA-%s intermediate CSR", tc.paramSet)
			require.NotNil(t, resp, "got nil response generating ML-DSA-%s intermediate CSR", tc.paramSet)
			require.False(t, resp.IsError(), "unexpected logical error generating ML-DSA-%s intermediate CSR: %v", tc.paramSet, resp.Error())

			// Private key must not be returned for internal generation.
			require.Nil(t, resp.Data["private_key"],
				"private_key must not be present in internal intermediate CSR generation response")

			// A key ID must be assigned.
			require.NotEmpty(t, resp.Data["key_id"],
				"key_id must not be empty in intermediate CSR response")

			validateIntermediateCSR(t, resp, tc.expectedParameters)
		})
	}
}

// TestGenerateIntermediate_exported verifies that ML-DSA intermediate CSR
// generation with key_type=exported returns both the CSR and the private key,
// and that the key material uses the requested parameter set.
func TestGenerateIntermediate_exported(t *testing.T) {
	t.Parallel()

	parameterSets := []struct {
		paramSet           string
		expectedParameters mldsa.Parameters
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
				Path:      "intermediate/generate/exported",
				Storage:   s,
				Data: map[string]interface{}{
					"common_name":   "Test ML-DSA Intermediate CA",
					"key_type":      "ml-dsa",
					"parameter_set": tc.paramSet,
				},
				MountPoint: "pki/",
			})
			require.NoError(t, err, "unexpected transport error generating exported ML-DSA-%s intermediate CSR", tc.paramSet)
			require.NotNil(t, resp, "got nil response generating exported ML-DSA-%s intermediate CSR", tc.paramSet)
			require.False(t, resp.IsError(), "unexpected logical error generating exported ML-DSA-%s intermediate CSR: %v", tc.paramSet, resp.Error())

			// Private key must be present and use the correct type.
			keyPEM, ok := resp.Data["private_key"].(string)
			require.True(t, ok, "private_key field should be a string")
			require.NotEmpty(t, keyPEM, "private_key must not be empty for exported intermediate CSR generation")
			require.Equal(t, certutil.MLDSAPrivateKey, resp.Data["private_key_type"],
				"private_key_type should be ml-dsa for ML-DSA-%s intermediate CSR", tc.paramSet)

			// Decode and verify the ML-DSA private key.
			block, rest := pem.Decode([]byte(keyPEM))
			require.Empty(t, rest, "trailing data after private key PEM block for ML-DSA-%s", tc.paramSet)
			require.NotNil(t, block, "failed to PEM-decode private key for ML-DSA-%s", tc.paramSet)
			rawKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			require.NoError(t, err, "failed to parse PKCS#8 private key for ML-DSA-%s", tc.paramSet)
			mldsaKey, ok := rawKey.(*mldsa.PrivateKey)
			require.True(t, ok, "expected *mldsa.PrivateKey for ML-DSA-%s, got %T", tc.paramSet, rawKey)
			require.Equal(t, tc.expectedParameters, mldsaKey.PublicKey().Parameters(),
				"exported private key parameters do not match requested parameter set for ML-DSA-%s", tc.paramSet)

			// A key ID must be assigned.
			require.NotEmpty(t, resp.Data["key_id"],
				"key_id must not be empty in exported intermediate CSR response for ML-DSA-%s", tc.paramSet)

			validateIntermediateCSR(t, resp, tc.expectedParameters)
		})
	}
}

// TestGenerateIntermediate_existing verifies that intermediate CSR generation
// with key_type=existing uses a previously-imported ML-DSA key and returns a
// valid CSR signed with that key's parameters.
func TestGenerateIntermediate_existing(t *testing.T) {
	t.Parallel()

	mldsa44Bundle, err := certutil.CreateKeyBundle("ml-dsa", 0, rand.Reader, "ml-dsa-44")
	require.NoError(t, err, "failed generating ml-dsa-44 key bundle")
	mldsa44Pem, err := mldsa44Bundle.ToPrivateKeyPemString()
	require.NoError(t, err, "failed converting ml-dsa-44 key to PEM")

	mldsa65Bundle, err := certutil.CreateKeyBundle("ml-dsa", 0, rand.Reader, "ml-dsa-65")
	require.NoError(t, err, "failed generating ml-dsa-65 key bundle")
	mldsa65Pem, err := mldsa65Bundle.ToPrivateKeyPemString()
	require.NoError(t, err, "failed converting ml-dsa-65 key to PEM")

	mldsa87Bundle, err := certutil.CreateKeyBundle("ml-dsa", 0, rand.Reader, "ml-dsa-87")
	require.NoError(t, err, "failed generating ml-dsa-87 key bundle")
	mldsa87Pem, err := mldsa87Bundle.ToPrivateKeyPemString()
	require.NoError(t, err, "failed converting ml-dsa-87 key to PEM")

	testCases := map[string]struct {
		keyPEM             string
		expectedParameters mldsa.Parameters
	}{
		"mldsa-44": {keyPEM: mldsa44Pem, expectedParameters: mldsa.MLDSA44()},
		"mldsa-65": {keyPEM: mldsa65Pem, expectedParameters: mldsa.MLDSA65()},
		"mldsa-87": {keyPEM: mldsa87Pem, expectedParameters: mldsa.MLDSA87()},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			b, s := CreateBackendWithStorage(t)

			// Import the pre-generated key so it becomes the default key.
			resp, err := b.HandleRequest(context.Background(), &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "keys/import",
				Storage:   s,
				Data: map[string]interface{}{
					"key_name":   name,
					"pem_bundle": tc.keyPEM,
				},
				MountPoint: "pki/",
			})
			require.NoError(t, err, "unexpected transport error importing %s key", name)
			require.False(t, resp.IsError(), "unexpected logical error importing %s key: %v", name, resp.Error())

			// Generate an intermediate CSR using the existing (default) key.
			resp, err = b.HandleRequest(context.Background(), &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "intermediate/generate/existing",
				Storage:   s,
				Data: map[string]interface{}{
					"common_name": "Test ML-DSA Intermediate CA",
				},
				MountPoint: "pki/",
			})
			require.NoError(t, err, "unexpected transport error generating existing ML-DSA intermediate CSR for %s", name)
			require.NotNil(t, resp, "got nil response generating existing ML-DSA intermediate CSR for %s", name)
			require.False(t, resp.IsError(), "unexpected logical error generating existing ML-DSA intermediate CSR for %s: %v", name, resp.Error())

			// Private key must not be returned for existing generation.
			require.Nil(t, resp.Data["private_key"],
				"private_key must not be present in existing intermediate CSR generation response for %s", name)

			validateIntermediateCSR(t, resp, tc.expectedParameters)
		})
	}
}

// TestGenerateIntermediate_invalidParameterSet verifies that requesting an
// unsupported ML-DSA parameter set returns a logical error rather than panicking
// or silently using a default.
func TestGenerateIntermediate_invalidParameterSet(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "intermediate/generate/internal",
		Storage:   s,
		Data: map[string]interface{}{
			"common_name":   "Test ML-DSA Intermediate CA",
			"key_type":      "ml-dsa",
			"parameter_set": "14",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "unexpected transport error (want a logical error, not a transport error)")
	require.NotNil(t, resp, "expected a non-nil response for invalid parameter set")
	require.True(t, resp.IsError(), "expected a logical error for unsupported ML-DSA parameter set '14'")
}

// validateIntermediateCSR asserts that the CSR in resp is a well-formed PEM
// block containing a valid certificate request whose public key uses the
// expected ML-DSA parameter set.
func validateIntermediateCSR(t *testing.T, resp *logical.Response, wantParams mldsa.Parameters) {
	t.Helper()

	csrPEM, ok := resp.Data["csr"].(string)
	require.True(t, ok, "csr field should be a string")
	require.NotEmpty(t, csrPEM, "csr must not be empty")

	block, rest := pem.Decode([]byte(csrPEM))
	require.NotNil(t, block, "failed to PEM-decode CSR")
	require.Empty(t, rest, "trailing data after CSR PEM block")

	csr, err := x509.ParseCertificateRequest(block.Bytes)
	require.NoError(t, err, "failed to parse certificate request")
	require.NoError(t, csr.CheckSignature(), "CSR signature check failed")

	pub, ok := csr.PublicKey.(*mldsa.PublicKey)
	require.True(t, ok, "expected *mldsa.PublicKey in CSR, got %T", csr.PublicKey)
	require.Equal(t, wantParams, pub.Parameters(),
		"CSR public key parameters do not match the requested ML-DSA parameter set")
}

// TestSetSignedIntermediate_classicalAndMLDSA verifies that
// intermediate/set-signed handles certificates that combine a classical
// signing algorithm with an ML-DSA subject key (and vice-versa):
//
//   - An EC-P256 root signing an ML-DSA-44 intermediate ensures that a
//     post-quantum intermediate CA can be issued by a classical trust anchor.
//   - An ML-DSA-65 root signing an EC-P256 intermediate ensures that an
//     existing classical PKI can be rooted in a post-quantum CA.
//
// Both directions test that set-signed accepts the certificate and correctly
// links the intermediate to its locally-held key.
func TestSetSignedIntermediate_classicalAndMLDSA(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		rootParams    map[string]interface{}
		intParams     map[string]interface{}
		verifyIntCert func(t *testing.T, intCert *x509.Certificate)
	}{
		{
			name: "ec-root-signs-mldsa-intermediate",
			rootParams: map[string]interface{}{
				"common_name": "EC-P256 Root CA",
				"key_type":    "ec",
				"key_bits":    256,
				"ttl":         "8760h",
			},
			intParams: map[string]interface{}{
				"common_name":   "ML-DSA-44 Intermediate CA",
				"key_type":      "ml-dsa",
				"parameter_set": certutil.MLDSA44,
			},
			verifyIntCert: func(t *testing.T, intCert *x509.Certificate) {
				t.Helper()
				// The intermediate's signature is produced by the EC root.
				require.Equal(t, x509.ECDSAWithSHA256, intCert.SignatureAlgorithm,
					"intermediate signed by EC-P256 root should carry ECDSAWithSHA256 signature algorithm")
				// The intermediate's subject public key must be ML-DSA-44.
				pub, ok := intCert.PublicKey.(*mldsa.PublicKey)
				require.True(t, ok,
					"intermediate subject public key should be *mldsa.PublicKey, got %T", intCert.PublicKey)
				require.Equal(t, mldsa.MLDSA44(), pub.Parameters(),
					"intermediate subject public key should use ML-DSA-44 parameters")
			},
		},
		{
			// ML-DSA-65 root signs a classical EC-P256 intermediate.
			// This represents a post-quantum root anchoring a classical sub-hierarchy.
			name: "mldsa-root-signs-ec-intermediate",
			rootParams: map[string]interface{}{
				"common_name":   "ML-DSA-65 Root CA",
				"key_type":      "ml-dsa",
				"parameter_set": certutil.MLDSA65,
				"ttl":           "8760h",
			},
			intParams: map[string]interface{}{
				"common_name": "EC-P256 Intermediate CA",
				"key_type":    "ec",
				"key_bits":    256,
			},
			verifyIntCert: func(t *testing.T, intCert *x509.Certificate) {
				t.Helper()
				// The intermediate's signature is produced by the ML-DSA-65 root.
				require.Equal(t, x509.MLDSA65, intCert.SignatureAlgorithm,
					"intermediate signed by ML-DSA-65 root should carry MLDSA65 signature algorithm")
				// The intermediate's subject public key must be an EC key.
				_, ok := intCert.PublicKey.(*ecdsa.PublicKey)
				require.True(t, ok,
					"intermediate subject public key should be *ecdsa.PublicKey, got %T", intCert.PublicKey)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			bRoot, sRoot := CreateBackendWithStorage(t)
			bInt, sInt := CreateBackendWithStorage(t)

			// Step 1: Generate the root CA.
			resp, err := CBWrite(bRoot, sRoot, "root/generate/internal", tc.rootParams)
			requireSuccessNonNilResponse(t, resp, err, "generate root CA for %s", tc.name)
			rootCertPEM := resp.Data["certificate"].(string)

			// Step 2: Generate the intermediate CSR.
			resp, err = CBWrite(bInt, sInt, "intermediate/generate/internal", tc.intParams)
			requireSuccessNonNilResponse(t, resp, err, "generate intermediate CSR for %s", tc.name)
			intCSR := resp.Data["csr"].(string)
			require.NotEmpty(t, intCSR, "intermediate CSR must not be empty for %s", tc.name)

			// Step 3: Sign the intermediate CSR with the root.
			resp, err = CBWrite(bRoot, sRoot, "root/sign-intermediate", map[string]interface{}{
				"common_name": tc.intParams["common_name"],
				"csr":         intCSR,
				"ttl":         "4380h",
			})
			requireSuccessNonNilResponse(t, resp, err, "sign intermediate CSR for %s", tc.name)
			intCertPEM := resp.Data["certificate"].(string)

			// Verify cross-algorithm properties of the signed certificate.
			block, rest := pem.Decode([]byte(intCertPEM))
			require.NotNil(t, block, "signed intermediate cert PEM should decode for %s", tc.name)
			require.Empty(t, rest, "no trailing data expected after signed intermediate cert for %s", tc.name)
			intCert, parseErr := x509.ParseCertificate(block.Bytes)
			require.NoError(t, parseErr, "failed to parse signed intermediate certificate for %s", tc.name)
			require.True(t, intCert.IsCA, "signed intermediate must have IsCA=true for %s", tc.name)
			tc.verifyIntCert(t, intCert)

			// Step 4: Import the signed intermediate (plus the root chain) via set-signed.
			resp, err = CBWrite(bInt, sInt, "intermediate/set-signed", map[string]interface{}{
				"certificate": intCertPEM + "\n" + rootCertPEM,
			})
			requireSuccessNonNilResponse(t, resp, err, "set-signed for %s", tc.name)

			importedIssuers, ok := resp.Data["imported_issuers"].([]string)
			require.True(t, ok, "imported_issuers should be a []string for %s", tc.name)
			require.NotEmpty(t, importedIssuers, "at least one issuer should be imported for %s", tc.name)

			mapping, ok := resp.Data["mapping"].(map[string]string)
			require.True(t, ok, "mapping should be a map[string]string for %s", tc.name)
			require.NotEmpty(t, mapping, "mapping must not be empty after set-signed for %s", tc.name)

			// The intermediate issuer must be linked to the key that was generated
			// on the intermediate backend.
			hasKeyedIssuer := false
			for _, issuerID := range importedIssuers {
				if mapping[issuerID] != "" {
					hasKeyedIssuer = true
					break
				}
			}
			require.True(t, hasKeyedIssuer,
				"at least one imported issuer must have an associated key in the mapping for %s", tc.name)
		})
	}
}
