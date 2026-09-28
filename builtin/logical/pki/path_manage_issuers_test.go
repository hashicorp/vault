// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build enterprise

package pki

import (
	"crypto/mldsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// mldsaParamSets is the full set of supported ML-DSA parameter sets paired with
// the expected mldsa.Parameters value for certificate/key assertions.
var mldsaParamSets = []struct {
	paramSet     string
	expectedAlgo mldsa.Parameters
}{
	{certutil.MLDSA44, mldsa.MLDSA44()},
	{certutil.MLDSA65, mldsa.MLDSA65()},
	{certutil.MLDSA87, mldsa.MLDSA87()},
}

// TestIssuerGenerateRoot_internal verifies that the issuers/generate/root/internal
// endpoint generates a valid self-signed ML-DSA root CA for each supported parameter set,
// stores the key internally, and does not return private key material in the response.
func TestIssuerGenerateRoot_internal(t *testing.T) {
	t.Parallel()

	for _, tc := range mldsaParamSets {
		t.Run("mldsa-"+tc.paramSet, func(t *testing.T) {
			t.Parallel()
			b, s := CreateBackendWithStorage(t)

			resp, err := CBWrite(b, s, "issuers/generate/root/internal", map[string]interface{}{
				"common_name":   "Test ML-DSA Root CA",
				"key_type":      "ml-dsa",
				"parameter_set": tc.paramSet,
				"ttl":           "8760h",
			})
			requireSuccessNonNilResponse(t, resp, err, "generate ML-DSA-%s root via issuers/generate/root/internal", tc.paramSet)

			require.Nil(t, resp.Data["private_key"],
				"private_key must not be present in internal root generation response for ML-DSA-%s", tc.paramSet)
			require.NotEmpty(t, resp.Data["key_id"], "key_id must not be empty for ML-DSA-%s", tc.paramSet)
			require.NotEmpty(t, resp.Data["issuer_id"], "issuer_id must not be empty for ML-DSA-%s", tc.paramSet)

			requireMLDSARootCert(t, resp, tc.expectedAlgo, tc.paramSet)
		})
	}
}

// TestIssuerGenerateRoot_exported verifies that the issuers/generate/root/exported
// endpoint generates a valid self-signed ML-DSA root CA for each supported parameter set
// and returns the private key and its type in the response.
func TestIssuerGenerateRoot_exported(t *testing.T) {
	t.Parallel()

	for _, tc := range mldsaParamSets {
		t.Run("mldsa-"+tc.paramSet, func(t *testing.T) {
			t.Parallel()
			b, s := CreateBackendWithStorage(t)

			resp, err := CBWrite(b, s, "issuers/generate/root/exported", map[string]interface{}{
				"common_name":   "Test ML-DSA Root CA",
				"key_type":      "ml-dsa",
				"parameter_set": tc.paramSet,
				"ttl":           "8760h",
			})
			requireSuccessNonNilResponse(t, resp, err, "generate ML-DSA-%s root via issuers/generate/root/exported", tc.paramSet)

			require.Equal(t, certutil.ParameterSet(tc.paramSet), resp.Data["parameter_set"],
				"parameter_set should match requested value for ML-DSA-%s root", tc.paramSet)

			requireMLDSAExportedPrivateKey(t, resp, tc.expectedAlgo, tc.paramSet)
			requireMLDSARootCert(t, resp, tc.expectedAlgo, tc.paramSet)
		})
	}
}

// TestIssuerGenerateRoot_existing verifies that the issuers/generate/root/existing
// endpoint generates a valid self-signed ML-DSA root CA using a previously-imported key
// and does not return private key material in the response.
func TestIssuerGenerateRoot_existing(t *testing.T) {
	t.Parallel()

	for _, tc := range mldsaParamSets {
		t.Run("mldsa-"+tc.paramSet, func(t *testing.T) {
			t.Parallel()
			b, s := CreateBackendWithStorage(t)

			keyPEM := generateMLDSAKeyPEM(t, tc.paramSet)
			importMLDSAKey(t, b, s, "mldsa-"+tc.paramSet, keyPEM)

			resp, err := CBWrite(b, s, "issuers/generate/root/existing", map[string]interface{}{
				"common_name": "Test ML-DSA Root CA",
				"ttl":         "8760h",
			})
			requireSuccessNonNilResponse(t, resp, err, "generate root from existing ML-DSA-%s key", tc.paramSet)

			require.Nil(t, resp.Data["private_key"],
				"private_key must not be present in existing root generation response for ML-DSA-%s", tc.paramSet)

			requireMLDSARootCert(t, resp, tc.expectedAlgo, tc.paramSet)
		})
	}
}

// TestIssuerGenerateIntermediate_internal verifies that the
// issuers/generate/intermediate/internal endpoint generates a valid ML-DSA
// intermediate CSR without returning private key material in the response.
func TestIssuerGenerateIntermediate_internal(t *testing.T) {
	t.Parallel()

	for _, tc := range mldsaParamSets {
		t.Run("mldsa-"+tc.paramSet, func(t *testing.T) {
			t.Parallel()
			b, s := CreateBackendWithStorage(t)

			resp, err := CBWrite(b, s, "issuers/generate/intermediate/internal", map[string]interface{}{
				"common_name":   "Test ML-DSA Intermediate CA",
				"key_type":      "ml-dsa",
				"parameter_set": tc.paramSet,
			})
			requireSuccessNonNilResponse(t, resp, err, "generate ML-DSA-%s intermediate CSR via issuers/generate/intermediate/internal", tc.paramSet)

			require.Nil(t, resp.Data["private_key"],
				"private_key must not be present in internal intermediate CSR response for ML-DSA-%s", tc.paramSet)
			require.NotEmpty(t, resp.Data["key_id"],
				"key_id must not be empty in intermediate CSR response for ML-DSA-%s", tc.paramSet)

			requireMLDSAIntermediateCSR(t, resp, tc.expectedAlgo, tc.paramSet)
		})
	}
}

// TestIssuerGenerateIntermediate_exported verifies that the
// issuers/generate/intermediate/exported endpoint generates a valid ML-DSA
// intermediate CSR and returns both the CSR and the private key.
func TestIssuerGenerateIntermediate_exported(t *testing.T) {
	t.Parallel()

	for _, tc := range mldsaParamSets {
		t.Run("mldsa-"+tc.paramSet, func(t *testing.T) {
			t.Parallel()
			b, s := CreateBackendWithStorage(t)

			resp, err := CBWrite(b, s, "issuers/generate/intermediate/exported", map[string]interface{}{
				"common_name":   "Test ML-DSA Intermediate CA",
				"key_type":      "ml-dsa",
				"parameter_set": tc.paramSet,
			})
			requireSuccessNonNilResponse(t, resp, err, "generate exported ML-DSA-%s intermediate CSR via issuers/generate/intermediate/exported", tc.paramSet)

			require.NotEmpty(t, resp.Data["key_id"],
				"key_id must not be empty in exported intermediate CSR response for ML-DSA-%s", tc.paramSet)

			requireMLDSAExportedPrivateKey(t, resp, tc.expectedAlgo, tc.paramSet)
			requireMLDSAIntermediateCSR(t, resp, tc.expectedAlgo, tc.paramSet)
		})
	}
}

// TestIssuerGenerateIntermediate_existing verifies that the
// issuers/generate/intermediate/existing endpoint generates a valid ML-DSA
// intermediate CSR using a previously-imported key without returning private key material.
func TestIssuerGenerateIntermediate_existing(t *testing.T) {
	t.Parallel()

	for _, tc := range mldsaParamSets {
		t.Run("mldsa-"+tc.paramSet, func(t *testing.T) {
			t.Parallel()
			b, s := CreateBackendWithStorage(t)

			keyPEM := generateMLDSAKeyPEM(t, tc.paramSet)
			importMLDSAKey(t, b, s, "mldsa-"+tc.paramSet, keyPEM)

			resp, err := CBWrite(b, s, "issuers/generate/intermediate/existing", map[string]interface{}{
				"common_name": "Test ML-DSA Intermediate CA",
			})
			requireSuccessNonNilResponse(t, resp, err, "generate intermediate CSR from existing ML-DSA-%s key", tc.paramSet)

			require.Nil(t, resp.Data["private_key"],
				"private_key must not be present in existing intermediate CSR response for ML-DSA-%s", tc.paramSet)

			requireMLDSAIntermediateCSR(t, resp, tc.expectedAlgo, tc.paramSet)
		})
	}
}

// TestImportIssuer_bundle verifies that an ML-DSA CA certificate plus its private key
// can be imported via issuers/import/bundle and produces a usable issuer.
func TestImportIssuer_bundle(t *testing.T) {
	t.Parallel()

	for _, tc := range mldsaParamSets {
		t.Run("mldsa-"+tc.paramSet, func(t *testing.T) {
			t.Parallel()
			b, s := CreateBackendWithStorage(t)

			// Generate an ML-DSA root and capture the certificate and private key.
			genResp, err := CBWrite(b, s, "issuers/generate/root/exported", map[string]interface{}{
				"common_name":   "Test ML-DSA Import Root",
				"key_type":      "ml-dsa",
				"parameter_set": tc.paramSet,
				"ttl":           "8760h",
			})
			requireSuccessNonNilResponse(t, genResp, err, "generate ML-DSA-%s root for import test", tc.paramSet)

			certPEM := genResp.Data["certificate"].(string)
			keyPEM := genResp.Data["private_key"].(string)

			// Import the certificate+key into a fresh backend.
			bImport, sImport := CreateBackendWithStorage(t)
			importResp, err := CBWrite(bImport, sImport, "issuers/import/bundle", map[string]interface{}{
				"pem_bundle": keyPEM + "\n" + certPEM,
			})
			requireSuccessNonNilResponse(t, importResp, err, "import ML-DSA-%s bundle", tc.paramSet)

			importedIssuers := requireImportResponse(t, importResp, 1, 1, tc.paramSet)

			// Mapping must link the new issuer to its key.
			mapping, ok := importResp.Data["mapping"].(map[string]string)
			require.True(t, ok, "mapping should be a map[string]string for ML-DSA-%s", tc.paramSet)
			require.NotEmpty(t, mapping[importedIssuers[0]],
				"imported issuer should have an associated key in the mapping for ML-DSA-%s", tc.paramSet)

			// Importing the same bundle again must report the issuer as pre-existing.
			reimportResp, err := CBWrite(bImport, sImport, "issuers/import/bundle", map[string]interface{}{
				"pem_bundle": keyPEM + "\n" + certPEM,
			})
			requireSuccessNonNilResponse(t, reimportResp, err, "re-import ML-DSA-%s bundle", tc.paramSet)

			existingIssuers, ok := reimportResp.Data["existing_issuers"].([]string)
			require.True(t, ok, "existing_issuers should be a []string on re-import for ML-DSA-%s", tc.paramSet)
			require.Len(t, existingIssuers, 1,
				"expected exactly one existing issuer on re-import for ML-DSA-%s", tc.paramSet)
		})
	}
}

// TestImportIssuer_cert verifies that an ML-DSA CA certificate (without a private key)
// can be imported via issuers/import/cert and produces a usable issuer entry.
func TestImportIssuer_cert(t *testing.T) {
	t.Parallel()

	for _, tc := range mldsaParamSets {
		t.Run("mldsa-"+tc.paramSet, func(t *testing.T) {
			t.Parallel()
			b, s := CreateBackendWithStorage(t)

			// Generate an ML-DSA root to get a certificate to import.
			genResp, err := CBWrite(b, s, "issuers/generate/root/internal", map[string]interface{}{
				"common_name":   "Test ML-DSA Cert Import Root",
				"key_type":      "ml-dsa",
				"parameter_set": tc.paramSet,
				"ttl":           "8760h",
			})
			requireSuccessNonNilResponse(t, genResp, err, "generate ML-DSA-%s root for cert import test", tc.paramSet)

			certPEM := genResp.Data["certificate"].(string)

			// Import only the certificate into a fresh backend via issuers/import/cert.
			bImport, sImport := CreateBackendWithStorage(t)
			importResp, err := CBWrite(bImport, sImport, "issuers/import/cert", map[string]interface{}{
				"pem_bundle": certPEM,
			})
			requireSuccessNonNilResponse(t, importResp, err, "import ML-DSA-%s cert-only", tc.paramSet)

			// No keys should be imported when using the cert-only endpoint.
			requireImportResponse(t, importResp, 1, 0, tc.paramSet)
		})
	}
}

// TestImportIssuer_cert_rejectsPrivateKey verifies that submitting a private key to
// the issuers/import/cert endpoint is rejected with an informative error.
func TestImportIssuer_cert_rejectsPrivateKey(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// Generate an ML-DSA root to obtain a real certificate and private key.
	genResp, err := CBWrite(b, s, "issuers/generate/root/exported", map[string]interface{}{
		"common_name":   "Test ML-DSA Root",
		"key_type":      "ml-dsa",
		"parameter_set": certutil.MLDSA44,
		"ttl":           "8760h",
	})
	requireSuccessNonNilResponse(t, genResp, err, "generate ML-DSA-44 root")

	certPEM := genResp.Data["certificate"].(string)
	keyPEM := genResp.Data["private_key"].(string)

	// Sending a key bundle to the cert-only endpoint must be rejected.
	// CBWrite surfaces logical.ErrorResponse as a Go error (not resp.IsError()).
	bImport, sImport := CreateBackendWithStorage(t)
	_, err = CBWrite(bImport, sImport, "issuers/import/cert", map[string]interface{}{
		"pem_bundle": keyPEM + "\n" + certPEM,
	})
	require.Error(t, err, "expected an error when private key is submitted to issuers/import/cert")
	require.Contains(t, err.Error(), "/issuers/import/bundle",
		"error message should suggest the bundle endpoint")
}

// TestRevokeIssuer verifies that an ML-DSA issuer can be revoked via the
// issuer/{ref}/revoke endpoint, that the revoked field is true, and that a
// second revocation is idempotent.
func TestRevokeIssuer(t *testing.T) {
	t.Parallel()

	for _, tc := range mldsaParamSets {
		t.Run("mldsa-"+tc.paramSet, func(t *testing.T) {
			t.Parallel()
			b, s := CreateBackendWithStorage(t)

			issuerName := "test-root-" + tc.paramSet
			genResp, err := CBWrite(b, s, "issuers/generate/root/internal", map[string]interface{}{
				"common_name":   "Test ML-DSA Revoke Root",
				"key_type":      "ml-dsa",
				"parameter_set": tc.paramSet,
				"issuer_name":   issuerName,
				"ttl":           "8760h",
			})
			requireSuccessNonNilResponse(t, genResp, err, "generate ML-DSA-%s root for revocation test", tc.paramSet)

			requireIssuerRevoked(t, b, s, issuerName, tc.paramSet)

			// A second revocation must be idempotent.
			requireIssuerRevoked(t, b, s, issuerName, tc.paramSet)
		})
	}
}

// generateMLDSAKeyPEM generates an ML-DSA private key for the given parameter
// set and returns it as a PKCS#8 PEM string.
func generateMLDSAKeyPEM(t *testing.T, paramSet string) string {
	t.Helper()

	bundle, err := certutil.CreateKeyBundle("ml-dsa", 0, rand.Reader, certutil.ParameterSet(paramSet))
	require.NoError(t, err, "failed generating ml-dsa-%s key bundle", paramSet)
	keyPEM, err := bundle.ToPrivateKeyPemString()
	require.NoError(t, err, "failed converting ml-dsa-%s key to PEM", paramSet)
	return keyPEM
}

// importMLDSAKey imports a PEM-encoded ML-DSA private key into the backend under
// the given key name and asserts the import succeeds.
func importMLDSAKey(t *testing.T, b *backend, s logical.Storage, keyName, keyPEM string) {
	t.Helper()

	resp, err := CBWrite(b, s, "keys/import", map[string]interface{}{
		"key_name":   keyName,
		"pem_bundle": keyPEM,
	})
	requireSuccessNonNilResponse(t, resp, err, "import ML-DSA key %q", keyName)
}

// requireMLDSARootCert asserts that the response contains a parseable self-signed
// CA certificate whose public key is an ML-DSA key with the expected parameters.
func requireMLDSARootCert(t *testing.T, resp *logical.Response, expectedParams mldsa.Parameters, paramSetLabel string) {
	t.Helper()

	certPEM, ok := resp.Data["certificate"].(string)
	require.True(t, ok, "certificate field should be a string for ML-DSA-%s", paramSetLabel)
	require.NotEmpty(t, certPEM, "certificate must not be empty for ML-DSA-%s", paramSetLabel)

	cert := parseCert(t, certPEM)
	require.True(t, cert.IsCA, "generated root certificate must have IsCA=true for ML-DSA-%s", paramSetLabel)
	require.Equal(t, x509.MLDSA, cert.PublicKeyAlgorithm,
		"certificate public key algorithm should be MLDSA for ML-DSA-%s", paramSetLabel)

	pub, ok := cert.PublicKey.(*mldsa.PublicKey)
	require.True(t, ok, "expected *mldsa.PublicKey in certificate for ML-DSA-%s, got %T", paramSetLabel, cert.PublicKey)
	require.Equal(t, expectedParams, pub.Parameters(),
		"certificate public key parameters do not match expected ML-DSA parameter set for ML-DSA-%s", paramSetLabel)
}

// requireMLDSAExportedPrivateKey asserts that the response contains a valid
// exported ML-DSA private key with the expected parameters and correct type metadata.
func requireMLDSAExportedPrivateKey(t *testing.T, resp *logical.Response, expectedParams mldsa.Parameters, paramSetLabel string) {
	t.Helper()

	keyPEM, ok := resp.Data["private_key"].(string)
	require.True(t, ok, "private_key field should be a string for ML-DSA-%s", paramSetLabel)
	require.NotEmpty(t, keyPEM, "private_key must not be empty for exported generation of ML-DSA-%s", paramSetLabel)
	require.Equal(t, certutil.MLDSAPrivateKey, resp.Data["private_key_type"],
		"private_key_type should be ml-dsa for ML-DSA-%s", paramSetLabel)

	block, rest := pem.Decode([]byte(keyPEM))
	require.Empty(t, rest, "trailing data after private key PEM block for ML-DSA-%s", paramSetLabel)
	require.NotNil(t, block, "failed to PEM-decode private key for ML-DSA-%s", paramSetLabel)

	rawKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	require.NoError(t, err, "failed to parse PKCS#8 private key for ML-DSA-%s", paramSetLabel)

	mldsaKey, ok := rawKey.(*mldsa.PrivateKey)
	require.True(t, ok, "expected *mldsa.PrivateKey for ML-DSA-%s, got %T", paramSetLabel, rawKey)
	require.Equal(t, expectedParams, mldsaKey.PublicKey().Parameters(),
		"exported private key parameters do not match expected ML-DSA parameter set for ML-DSA-%s", paramSetLabel)
}

// requireMLDSAIntermediateCSR asserts that the response contains a well-formed CSR
// whose public key is an ML-DSA key with the expected parameters and a valid self-signature.
func requireMLDSAIntermediateCSR(t *testing.T, resp *logical.Response, expectedParams mldsa.Parameters, paramSetLabel string) {
	t.Helper()

	csrPEM, ok := resp.Data["csr"].(string)
	require.True(t, ok, "csr field should be a string for ML-DSA-%s", paramSetLabel)
	require.NotEmpty(t, csrPEM, "csr must not be empty for ML-DSA-%s", paramSetLabel)

	block, rest := pem.Decode([]byte(csrPEM))
	require.NotNil(t, block, "failed to PEM-decode CSR for ML-DSA-%s", paramSetLabel)
	require.Empty(t, rest, "trailing data after CSR PEM block for ML-DSA-%s", paramSetLabel)

	csr, err := x509.ParseCertificateRequest(block.Bytes)
	require.NoError(t, err, "failed to parse certificate request for ML-DSA-%s", paramSetLabel)
	require.NoError(t, csr.CheckSignature(), "CSR signature check failed for ML-DSA-%s", paramSetLabel)

	pub, ok := csr.PublicKey.(*mldsa.PublicKey)
	require.True(t, ok, "expected *mldsa.PublicKey in CSR for ML-DSA-%s, got %T", paramSetLabel, csr.PublicKey)
	require.Equal(t, expectedParams, pub.Parameters(),
		"CSR public key parameters do not match expected ML-DSA parameter set for ML-DSA-%s", paramSetLabel)
}

// requireImportResponse asserts that the import response lists exactly
// expectedIssuers newly-imported issuers and expectedKeys newly-imported keys, and
// returns the imported issuer IDs for further assertions by the caller.
func requireImportResponse(t *testing.T, resp *logical.Response, expectedIssuers, expectedKeys int, paramSetLabel string) []string {
	t.Helper()

	importedIssuers, ok := resp.Data["imported_issuers"].([]string)
	require.True(t, ok, "imported_issuers should be a []string for ML-DSA-%s", paramSetLabel)
	require.Len(t, importedIssuers, expectedIssuers,
		"expected %d new issuer(s) in import response for ML-DSA-%s", expectedIssuers, paramSetLabel)

	importedKeys, ok := resp.Data["imported_keys"].([]string)
	require.True(t, ok, "imported_keys should be a []string for ML-DSA-%s", paramSetLabel)
	require.Len(t, importedKeys, expectedKeys,
		"expected %d new key(s) in import response for ML-DSA-%s", expectedKeys, paramSetLabel)

	return importedIssuers
}

// requireIssuerRevoked revokes the named issuer via issuer/{name}/revoke and
// asserts that the response contains revoked=true.
func requireIssuerRevoked(t *testing.T, b *backend, s logical.Storage, issuerName, paramSetLabel string) {
	t.Helper()

	revokeResp, err := CBWrite(b, s, "issuer/"+issuerName+"/revoke", map[string]interface{}{})
	requireSuccessNonNilResponse(t, revokeResp, err, "revoke ML-DSA-%s issuer %q", paramSetLabel, issuerName)

	revoked, ok := revokeResp.Data["revoked"].(bool)
	require.True(t, ok, "revoked field should be a bool in revocation response for ML-DSA-%s", paramSetLabel)
	require.True(t, revoked, "revoked field must be true after revoking ML-DSA-%s issuer", paramSetLabel)
}
