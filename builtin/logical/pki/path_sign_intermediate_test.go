// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"crypto/mldsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"software.sslmate.com/src/go-pkcs12"
)

// TestSignIntermediate_MaxPathLengthValidation verifies that when signing an
// intermediate CA, the requested max_path_length is validated against the
// signing CA's BasicConstraints pathLenConstraint per RFC 5280 4.2.1.9.
//
// If the signing CA has a pathLenConstraint of N, any intermediate it signs
// must have a pathLenConstraint strictly less than N.
//
// When max_path_length is not provided at all (omitted from the request), the
// validation block is skipped entirely and the request always succeeds.
//
// When max_path_length is explicitly set to -1 (meaning "no constraint on the
// intermediate"), and the signing CA already has a pathLenConstraint, the
// request is rejected because an unconstrained intermediate would violate the
// CA's constraint.
//
// When the generated certificate has MaxPathLen == 0 (pathLenConstraint=0),
// Vault adds a warning that the certificate cannot be used to issue
// intermediate CAs.
func TestSignIntermediate_MaxPathLengthValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		caMaxPathLength  int  // max_path_length set on the root CA (-1 means not set / unconstrained)
		intMaxPathLength int  // max_path_length requested when signing the intermediate (-1 = "no constraint")
		omitParam        bool // if true, do not include max_path_length in the sign request at all
		expectError      bool
		errorContains    string
		// Expected values on the issued certificate (only checked when expectError==false).
		expectedMaxPathLen int // expected cert.MaxPathLen; -1 means "no constraint"
	}{
		// --- Explicit positive values: validation is active ---
		{
			name:               "ca_path_len_2_int_path_len_1_allowed",
			caMaxPathLength:    2,
			intMaxPathLength:   1,
			expectError:        false,
			expectedMaxPathLen: 1,
		},
		{
			name:             "ca_path_len_2_int_path_len_2_rejected",
			caMaxPathLength:  2,
			intMaxPathLength: 2,
			expectError:      true,
			errorContains:    "requested max_path_length 2 is not allowed",
		},
		{
			name:             "ca_path_len_2_int_path_len_3_rejected",
			caMaxPathLength:  2,
			intMaxPathLength: 3,
			expectError:      true,
			errorContains:    "requested max_path_length 3 is not allowed",
		},
		{
			// pathLenConstraint=0 on the issued cert; Go represents this with
			// MaxPathLen==0 and MaxPathLenZero==true.
			name:               "ca_path_len_1_int_path_len_0_allowed",
			caMaxPathLength:    1,
			intMaxPathLength:   0,
			expectError:        false,
			expectedMaxPathLen: 0,
		},
		{
			name:             "ca_path_len_1_int_path_len_1_rejected",
			caMaxPathLength:  1,
			intMaxPathLength: 1,
			expectError:      true,
			errorContains:    "requested max_path_length 1 is not allowed",
		},
		{
			// pathLenConstraint=0 on the CA means it may not issue further CAs
			// with any pathLenConstraint (including 0).
			name:             "ca_path_len_0_int_path_len_0_rejected",
			caMaxPathLength:  0,
			intMaxPathLength: 0,
			expectError:      true,
			errorContains:    "requested max_path_length 0 is not allowed",
		},
		{
			// CA has no pathLenConstraint; any explicit positive value is allowed.
			name:               "ca_no_path_len_constraint_int_path_len_5_allowed",
			caMaxPathLength:    -1, // no constraint on CA
			intMaxPathLength:   5,
			expectError:        false,
			expectedMaxPathLen: 5,
		},

		// --- Explicit -1 ("no constraint on intermediate"): rejected only when CA is constrained ---
		{
			// CA has no constraint; explicit -1 is allowed and results in no
			// pathLenConstraint on the issued certificate.
			name:               "ca_no_constraint_int_explicit_neg1_allowed",
			caMaxPathLength:    -1,
			intMaxPathLength:   -1,
			expectError:        false,
			expectedMaxPathLen: -1,
		},
		{
			// CA has a constraint; explicit -1 means "no constraint on the
			// intermediate", which violates the CA's pathLenConstraint.
			name:             "ca_path_len_5_int_explicit_neg1_rejected",
			caMaxPathLength:  5,
			intMaxPathLength: -1,
			expectError:      true,
			errorContains:    "requested max_path_length -1 is not allowed",
		},
		{
			// CA has pathLenConstraint=0; explicit -1 is also rejected.
			name:             "ca_path_len_0_int_explicit_neg1_rejected",
			caMaxPathLength:  0,
			intMaxPathLength: -1,
			expectError:      true,
			errorContains:    "requested max_path_length -1 is not allowed",
		},

		// --- Parameter omitted entirely: validation is skipped ---
		// When omitted, Vault auto-derives the pathLenConstraint as (CA pathLen - 1).
		// When the CA has no constraint, the issued cert also has no constraint (-1).
		{
			// CA has no constraint; omitting the parameter results in no
			// pathLenConstraint on the issued certificate.
			name:               "ca_no_constraint_int_param_omitted_allowed",
			caMaxPathLength:    -1,
			omitParam:          true,
			expectError:        false,
			expectedMaxPathLen: -1,
		},
		{
			// CA has pathLenConstraint=5; omitting the parameter causes Vault to
			// auto-derive pathLenConstraint = 5 - 1 = 4 on the issued certificate.
			name:               "ca_path_len_5_int_param_omitted_allowed",
			caMaxPathLength:    5,
			omitParam:          true,
			expectError:        false,
			expectedMaxPathLen: 4,
		},
		{
			// CA has pathLenConstraint=1; omitting the parameter causes Vault to
			// auto-derive pathLenConstraint = 1 - 1 = 0 on the issued certificate.
			name:               "ca_path_len_1_int_param_omitted_allowed",
			caMaxPathLength:    1,
			omitParam:          true,
			expectError:        false,
			expectedMaxPathLen: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			b, s := CreateBackendWithStorage(t)

			// Generate root CA with the specified max_path_length.
			rootParams := map[string]interface{}{
				"common_name": "root.example.com",
				"ttl":         "87600h",
			}
			if tc.caMaxPathLength >= 0 {
				rootParams["max_path_length"] = tc.caMaxPathLength
			}

			resp, err := CBWrite(b, s, "root/generate/internal", rootParams)
			requireSuccessNonNilResponse(t, resp, err, "failed to generate root CA")

			// Generate an intermediate CSR.
			resp, err = CBWrite(b, s, "intermediate/generate/internal", map[string]interface{}{
				"common_name": "int.example.com",
			})
			requireSuccessNonNilResponse(t, resp, err, "failed to generate intermediate CSR")
			csr := resp.Data["csr"].(string)

			// Build the sign-intermediate request.
			signParams := map[string]interface{}{
				"common_name": "int.example.com",
				"csr":         csr,
				"ttl":         "43800h",
			}
			if !tc.omitParam {
				signParams["max_path_length"] = tc.intMaxPathLength
			}

			resp, err = CBWrite(b, s, "root/sign-intermediate", signParams)

			if tc.expectError {
				require.Error(t, err, "expected sign-intermediate to fail but it succeeded")
				if tc.errorContains != "" {
					require.Contains(t, err.Error(), tc.errorContains,
						"error message did not contain expected text")
				}
			} else {
				require.NoError(t, err, "expected sign-intermediate to succeed but it failed")
				require.NotNil(t, resp, "expected non-nil response from sign-intermediate")
				require.False(t, resp.IsError(), "sign-intermediate returned error: %v", resp.Error())

				cert := parseCert(t, resp.Data["certificate"].(string))
				require.True(t, cert.BasicConstraintsValid, "expected BasicConstraints to be set")
				require.True(t, cert.IsCA, "expected IsCA to be true")

				require.Equal(t, tc.expectedMaxPathLen == 0, cert.MaxPathLenZero,
					"issued certificate has unexpected MaxPathLenZero")
				require.Equal(t, tc.expectedMaxPathLen, cert.MaxPathLen,
					"issued certificate has unexpected MaxPathLen")
				if tc.expectedMaxPathLen == 0 {
					requireMaxPathLengthZeroWarning(t, resp.Warnings)
				}
			}
		})
	}
}

// TestSignIntermediate_PKCS12AndJKSFormat validates PKCS12 and JKS output for intermediate signing
// and asserts encoder (PK12) or alias (JKS) parameter handling and trust store structure.
// PKCS12/JKS archive should be a trust store (no private key) containing the signed intermediate and CA chain.
func TestSignIntermediate_PKCS12AndJKSFormat(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// Generate a root CA
	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "Root CA",
		"issuer_name": "root-ca",
		"ttl":         "87600h",
		"key_type":    "ec",
		"key_bits":    256,
	})
	requireSuccessNonNilResponse(t, resp, err)

	// Generate an intermediate CSR
	resp, err = CBWrite(b, s, "intermediate/generate/internal", map[string]interface{}{
		"common_name": "Intermediate CA",
		"key_type":    "ec",
		"key_bits":    256,
	})
	requireSuccessNonNilResponse(t, resp, err)
	intermediateCsr := resp.Data["csr"].(string)

	for _, path := range []string{"root/sign-intermediate", "issuer/root-ca/sign-intermediate"} {

		t.Run("invalid format", func(t *testing.T) {
			_, err := CBWrite(b, s, path, map[string]interface{}{
				"format": "invalid",
				"csr":    intermediateCsr,
				"ttl":    "43800h",
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), `the "format" parameter must be "pem", "der", "pem_bundle", "pkcs12_bundle" or "jks_bundle"`)
		})

		testCasesPKCS12 := []struct {
			name         string
			encoder      string
			omitPassword bool
			password     string
			shouldError  bool
		}{
			{name: "custom password", password: "intermediate-password"},
			{name: "default password", password: pkcs12.DefaultPassword, omitPassword: true},
			{name: "empty password", password: ""},
			{name: "with modern2026 encoder", password: "intermediate-password", encoder: "modern2026"},
			{name: "with modern2023 encoder", password: "intermediate-password", encoder: "modern2023"},
			{name: "with invalid encoder", encoder: "modern2020", shouldError: true},
		}

		for _, tc := range testCasesPKCS12 {
			name := fmt.Sprintf("PKCS12 endpoint=%q %s", path, tc.name)
			t.Run(name, func(t *testing.T) {
				data := map[string]interface{}{
					"format": "pkcs12_bundle",
					"csr":    intermediateCsr,
					"ttl":    "43800h",
				}
				if !tc.omitPassword {
					data["pkcs12_password"] = tc.password
				}
				if tc.encoder != "" {
					data["pkcs12_encoder"] = tc.encoder
				}

				resp, err := CBWrite(b, s, path, data)
				pkcs12Bytes := verifyAndDecodePKCS12(t, path, resp, err, tc.shouldError)
				if !tc.shouldError {
					certs := requireDecodesPKCS12TrustStore(t, pkcs12Bytes, tc.password)
					// First cert should be the intermediate
					require.Len(t, certs, 2, "should contain intermediate + root CA")
					require.True(t, certs[0].IsCA, "first cert should be a CA")
					require.Equal(t, "Intermediate CA", certs[0].Subject.CommonName)
					require.True(t, certs[1].IsCA, "second cert should be CA")
					require.Equal(t, "Root CA", certs[1].Subject.CommonName)
				}
			})
		}

		testCasesJKS := []struct {
			name            string
			alias           string
			password        string
			expectedAliases []string
		}{
			{name: "default password and alias", expectedAliases: []string{"1", "2"}},
			{name: "custom password and alias", alias: "myapp", expectedAliases: []string{"1", "2"}, password: "my-very-secure-password"},
			{name: "custom numeric alias", alias: "4", expectedAliases: []string{"1", "2"}},
		}

		for _, tc := range testCasesJKS {
			name := fmt.Sprintf("JKS endpoint=%q %s", path, tc.name)
			t.Run(name, func(t *testing.T) {
				decryptPw := pkcs12.DefaultPassword
				data := map[string]interface{}{
					"format": "jks_bundle",
					"csr":    intermediateCsr,
					"ttl":    "43800h",
				}
				if tc.password != "" {
					decryptPw = tc.password
					data["jks_password"] = tc.password
				}
				if tc.alias != "" {
					data["jks_private_key_alias"] = tc.alias
				}

				resp, err := CBWrite(b, s, path, data)
				pkcs12Bytes := verifyAndDecodeJKS(t, path, resp, err)
				certs := requireDecodesJKSTrustStore(t, pkcs12Bytes, decryptPw, tc.expectedAliases)
				// First cert should be the intermediate
				require.Len(t, certs, 2, "should contain intermediate + root CA")
				require.True(t, certs[0].IsCA, "first cert should be a CA")
				require.Equal(t, "Intermediate CA", certs[0].Subject.CommonName)
				require.True(t, certs[1].IsCA, "second cert should be CA")
				require.Equal(t, "Root CA", certs[1].Subject.CommonName)
			})
		}
	}
}

// TestSignIntermediate_MLDSA exercises the full intermediate CA workflow for ML-DSA keys:
// an ML-DSA-44 root generates and signs an ML-DSA-65 intermediate CSR, the signed
// intermediate is imported into a separate PKI mount, and a leaf certificate is issued
// from the intermediate.
func TestSignIntermediate_MLDSA(t *testing.T) {
	t.Parallel()

	// Root CA backend (ML-DSA-44).
	bRoot, sRoot := CreateBackendWithStorage(t)
	// Intermediate CA backend (ML-DSA-65).
	bInt, sInt := CreateBackendWithStorage(t)

	// Step 1: Generate an ML-DSA-44 root CA.
	resp, err := CBWrite(bRoot, sRoot, "root/generate/internal", map[string]interface{}{
		"common_name":   "ML-DSA Root CA",
		"key_type":      "ml-dsa",
		"parameter_set": "ml-dsa-44",
		"ttl":           "87600h",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to generate ML-DSA-44 root CA")
	rootCertPEM := resp.Data["certificate"].(string)
	rootCert := parseCert(t, rootCertPEM)

	require.True(t, rootCert.IsCA, "root certificate should be a CA")
	require.True(t, rootCert.BasicConstraintsValid, "root certificate should have BasicConstraints")
	// A self-signed root is signed with its own ML-DSA-44 key.
	require.Equal(t, x509.MLDSA44, rootCert.SignatureAlgorithm,
		"root cert should use ML-DSA-44 signature algorithm")
	require.NotEmpty(t, rootCert.SubjectKeyId, "root certificate SubjectKeyId should be non-empty")

	rootPub, ok := rootCert.PublicKey.(*mldsa.PublicKey)
	require.True(t, ok, "root certificate public key should be *mldsa.PublicKey")
	require.Equal(t, mldsa.MLDSA44(), rootPub.Parameters(),
		"root certificate public key should use ML-DSA-44 parameters")

	// Step 2: Generate an ML-DSA-65 intermediate CSR on the intermediate backend.
	resp, err = CBWrite(bInt, sInt, "intermediate/generate/internal", map[string]interface{}{
		"common_name":   "ML-DSA Intermediate CA",
		"key_type":      "ml-dsa",
		"parameter_set": "ml-dsa-65",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to generate ML-DSA-65 intermediate CSR")
	intCSR := resp.Data["csr"].(string)
	require.NotEmpty(t, intCSR, "intermediate CSR should be non-empty")

	// Verify the CSR decodes as a valid PEM block.
	csrBlock, rest := pem.Decode([]byte(intCSR))
	require.NotNil(t, csrBlock, "intermediate CSR PEM should decode")
	require.Empty(t, rest, "intermediate CSR PEM should have no trailing data")

	// Step 3: Sign the intermediate CSR with the ML-DSA-44 root.
	resp, err = CBWrite(bRoot, sRoot, "root/sign-intermediate", map[string]interface{}{
		"common_name": "ML-DSA Intermediate CA",
		"csr":         intCSR,
		"ttl":         "43800h",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to sign ML-DSA intermediate CSR with ML-DSA-44 root")
	intCertPEM := resp.Data["certificate"].(string)
	intCert := parseCert(t, intCertPEM)

	require.True(t, intCert.IsCA, "intermediate certificate should be a CA")
	require.True(t, intCert.BasicConstraintsValid, "intermediate certificate should have BasicConstraints")
	require.Equal(t, x509.MLDSA44, intCert.SignatureAlgorithm,
		"intermediate cert should carry ML-DSA-44 signature algorithm (signed by the ML-DSA-44 root)")
	require.NotEmpty(t, intCert.SubjectKeyId, "intermediate certificate SubjectKeyId should be non-empty")

	intPub, ok := intCert.PublicKey.(*mldsa.PublicKey)
	require.True(t, ok, "intermediate certificate public key should be *mldsa.PublicKey")
	require.Equal(t, mldsa.MLDSA65(), intPub.Parameters(),
		"intermediate certificate public key should use ML-DSA-65 parameters")

	// The intermediate's issuer should match the root's subject.
	require.Equal(t, rootCert.RawSubject, intCert.RawIssuer,
		"intermediate certificate issuer should match root certificate subject")

	// Step 4: Import the signed intermediate (plus root chain) into the intermediate backend.
	resp, err = CBWrite(bInt, sInt, "intermediate/set-signed", map[string]interface{}{
		"certificate": intCertPEM + "\n" + rootCertPEM,
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to import signed ML-DSA intermediate certificate")

	// Step 5: Create a role on the intermediate backend and issue a leaf certificate.
	// The intermediate backend's default issuer is the ML-DSA-65 intermediate.
	_, err = CBWrite(bInt, sInt, "roles/mldsa-leaf-role", map[string]interface{}{
		"allow_any_name": true,
		"max_ttl":        "1h",
		"key_type":       "ml-dsa",
		"parameter_set":  "ml-dsa-65",
	})
	require.NoError(t, err, "failed to create ML-DSA leaf role")

	resp, err = CBWrite(bInt, sInt, "issue/mldsa-leaf-role", map[string]interface{}{
		"common_name": "leaf.example.com",
		"ttl":         "1h",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed to issue ML-DSA leaf certificate")

	leafCertPEM := resp.Data["certificate"].(string)
	leafCert := parseCert(t, leafCertPEM)

	require.False(t, leafCert.IsCA, "leaf certificate should not be a CA")
	// The leaf is signed by the ML-DSA-65 intermediate, so its SignatureAlgorithm must
	// be ML-DSA-65.
	require.Equal(t, x509.MLDSA65, leafCert.SignatureAlgorithm,
		"leaf certificate should use ML-DSA-65 signature algorithm (signed by ML-DSA-65 intermediate)")

	leafPub, ok := leafCert.PublicKey.(*mldsa.PublicKey)
	require.True(t, ok, "leaf certificate public key should be *mldsa.PublicKey")
	require.Equal(t, mldsa.MLDSA65(), leafPub.Parameters(),
		"leaf certificate public key should use ML-DSA-65 parameters (generated by the ML-DSA-65 role)")

	// Step 6: Validate the certificate chain: leaf → intermediate → root.
	rootPool := x509.NewCertPool()
	rootPool.AddCert(rootCert)

	intermediatePool := x509.NewCertPool()
	intermediatePool.AddCert(intCert)

	chains, err := leafCert.Verify(x509.VerifyOptions{
		Roots:         rootPool,
		Intermediates: intermediatePool,
		// Use the leaf's NotBefore so test TTLs do not cause time-based failures.
		CurrentTime: leafCert.NotBefore,
	})
	require.NoError(t, err, "leaf certificate chain verification failed")
	require.NotEmpty(t, chains, "expected at least one valid certificate chain")

	// The verified chain should be: leaf → intermediate → root (length 3).
	require.Len(t, chains[0], 3, "expected chain length of 3 (leaf, intermediate, root)")
	require.Equal(t, leafCert.SerialNumber, chains[0][0].SerialNumber,
		"first element in chain should be the leaf")
	require.Equal(t, intCert.SerialNumber, chains[0][1].SerialNumber,
		"second element in chain should be the intermediate")
	require.Equal(t, rootCert.SerialNumber, chains[0][2].SerialNumber,
		"third element in chain should be the root")
}
