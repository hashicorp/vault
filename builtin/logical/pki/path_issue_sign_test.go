// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"crypto/mldsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/stretchr/testify/require"
	"software.sslmate.com/src/go-pkcs12"
)

// TestPathIssueSign_KeyTypeAny is a regression test.  At one point, the signature bits
// were determined based on the keyType of the (CSR) key being signed, not by the (CA)
// signer key.  We also need to make sure we correctly set default bits for each key type.
func TestPathIssueSign_KeyTypeAny(t *testing.T) {
	// Try (Parent) Issuers of Several Different Types:
	// RSA2048; RSA3072; RSA4096; RSA8192; EC224; EC256; EC384; EC521; ED25519 (ignore keysize)
	// And Signature Bits of:
	// 224, 256, 384, 512 (+ Unset)
	// Try (Child) Issued Certificates of Each of Those Same (Key) Type (with a CSR)
	keyTypeOptions := []struct {
		name               string
		keyType            string
		keySize            int
		parameterSet       string
		defaultSigBits     int
		successfulBitSizes []int
		generatedCaBundle  string
		generatedLeafCsr   string
	}{
		{"rsa-2048", "rsa", 2048, "", 256, []int{0, 256, 384, 512}, rsa2048CaAndKey, rsa2048leafCsr},
		{"rsa-3072", "rsa", 3072, "", 256, []int{0, 256, 384, 512}, rsa3072CaAndKey, rsa3072leafCsr},
		{"rsa-4096", "rsa", 4096, "", 256, []int{0, 256, 384, 512}, rsa4096CaAndKey, rsa4096leafCsr},
		{"rsa-8192", "rsa", 8192, "", 256, []int{0, 256, 384, 512}, rsa8192CaAndKey, rsa8192leafCsr},
		{"ec-224", "ec", 224, "", 256, []int{0, 256}, ec224CaAndKey, ec224leafCsr},
		{"ec-256", "ec", 256, "", 256, []int{0, 256}, ec256CaAndKey, ec256leafCsr},
		{"ec-384", "ec", 384, "", 384, []int{0, 384}, ec384CaAndKey, ec384leafCsr},
		{"ec-521", "ec", 521, "", 512, []int{0, 512}, ec512CaAndKey, ec512leafCsr},
		{"ed25519", "ed25519", 0, "", 512, []int{0, 512}, ed25519CaAndKey, ed25519leafCsr},
		{"ml-dsa-44", "ml-dsa", 0, "44", mldsa.MLDSA44SignatureSize * 8, []int{0}, mldsa44CaAndKey, mldsa44LeafCsr},
		{"ml-dsa-65", "ml-dsa", 0, "65", mldsa.MLDSA65SignatureSize * 8, []int{0}, mldsa65CaAndKey, mldsa65LeafCsr},
		{"ml-dsa-87", "ml-dsa", 0, "87", mldsa.MLDSA87SignatureSize * 8, []int{0}, mldsa87CaAndKey, mldsa87LeafCsr},
	}
	signatureBitOptions := []int{0, 224, 256, 384, 512, 1117}

	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	for _, parentKeyType := range keyTypeOptions {
		resp, err := CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
			"pem_bundle": parentKeyType.generatedCaBundle,
		})
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected ca info")
		}
		if resp.IsError() {
			t.Fatal("expected successful import", resp.Error())
		}
		issuers := resp.Data["imported_issuers"].([]string)
		issuerId := issuers[0]
		resp, err = CBWrite(b, s, "issuer/"+issuerId, map[string]interface{}{
			"issuer_name": parentKeyType.name,
		})
		if err != nil {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected ca info updated")
		}
		if resp.IsError() {
			t.Fatal("expected successful update", resp.Error())
		}

		for _, signatureBitOption := range signatureBitOptions {
			roleName := "issuer_" + parentKeyType.name + "_sigBits_" + strconv.Itoa(signatureBitOption)

			// Now Make a Role with Key-Type Any, Set Signature-Bits:
			roleResp, err := CBWrite(b, s, "roles/"+roleName, map[string]interface{}{
				"allowed_domains":    "foobar.com",
				"allow_bare_domains": true,
				"max_ttl":            "2h",
				"key_type":           "any",
				"issuer_ref":         parentKeyType.name,
				"signature_bits":     signatureBitOption,
			})
			if err != nil { // We should never fail to create a role because of misconfigured Signature bits if the
				// configuration might be valid for a different issuer, since issuers are updated separately.
				t.Fatal(fmt.Errorf("test failed creating role %v: %v", roleName, err))
			}
			signingFailureExpected := false
			defaultSizeOverride := false
			if slices.Contains(parentKeyType.successfulBitSizes, signatureBitOption) {
				require.Equal(t, len(roleResp.Warnings), 0)
			} else {
				if parentKeyType.keyType == "rsa" {
					signingFailureExpected = true
					require.Contains(t, roleResp.Warnings, fmt.Sprintf("The Issuing Certificate %v for this role has a key algorithm, %v, incompatible with the set role signature bits, %d", parentKeyType.name, parentKeyType.keyType, signatureBitOption))
				} else {
					defaultSizeOverride = true
				}
			}

			// For Each Key-Type, Generate a CSR and try to have it signed by the Role
			for _, childKeyTypeOption := range keyTypeOptions {

				resp, err = CBWrite(b, s, "sign/"+roleName, map[string]interface{}{
					"common_name": "foobar.com",
					"csr":         childKeyTypeOption.generatedLeafCsr,
				})

				if signingFailureExpected {
					require.Error(t, err, "expected signing failure", roleName, childKeyTypeOption.keyType, childKeyTypeOption.keySize)
				} else {
					if err != nil {
						t.Fatal(fmt.Errorf("test failed signing csr with keyType %v size %d with role %v: %v", childKeyTypeOption.keyType, childKeyTypeOption.keySize, roleName, err))
					}
					if resp == nil {
						t.Fatal(fmt.Errorf("test role %v didn't give cert response to attemp to sign CSR with %v keyType", roleName, childKeyTypeOption))
					}

					rawCert := resp.Data["certificate"].(string)
					trimmedRawCert, _ := strings.CutPrefix(rawCert, "-----BEGIN CERTIFICATE-----\n")
					moreTrimmedCert, _ := strings.CutSuffix(trimmedRawCert, "-----END CERTIFICATE-----")
					cleanCert := strings.ReplaceAll(moreTrimmedCert, "\n", "")
					certBytes, err := base64.StdEncoding.DecodeString(cleanCert)
					if err != nil {
						require.NoError(t, err, "failed to decode certificate")
					}
					cert, err := x509.ParseCertificate(certBytes)
					if err != nil {
						require.NoError(t, err, "failed to parse certificate")
					}

					// Signature is Truncated, So We have to Look At Type To Get Digest Length
					resultingLength := signingBits(cert.SignatureAlgorithm)
					if signatureBitOption == 0 || defaultSizeOverride {
						require.Equal(t, parentKeyType.defaultSigBits, resultingLength, "signature length was not default size", "parent key type", parentKeyType.name, "signature bits set", signatureBitOption, "default signature bits", parentKeyType.defaultSigBits, "childKeyType", childKeyTypeOption.name)
					} else {
						require.Equal(t, signatureBitOption, resultingLength, "signature length was not what was set", "role", roleName, "parent key type", parentKeyType.name, "signature bits set", signatureBitOption, "childKeyType", childKeyTypeOption.name)
					}
				}
			}

		}
	}
}

func signingBits(alg x509.SignatureAlgorithm) int {
	switch alg {
	case x509.SHA256WithRSA, x509.SHA256WithRSAPSS, x509.ECDSAWithSHA256:
		return 256
	case x509.SHA384WithRSA, x509.SHA384WithRSAPSS, x509.ECDSAWithSHA384:
		return 384
	case x509.SHA512WithRSA, x509.SHA512WithRSAPSS, x509.ECDSAWithSHA512, x509.PureEd25519:
		return 512
	case x509.MLDSA44:
		return mldsa.MLDSA44SignatureSize * 8
	case x509.MLDSA65:
		return mldsa.MLDSA65SignatureSize * 8
	case x509.MLDSA87:
		return mldsa.MLDSA87SignatureSize * 8
	default:
		return 0
	}
}

// TestPathIssueSign_PKCS12Format validates PKCS12 output for /issue and /sign endpoints:
// - /issue: PKCS12 contains private key, issued cert, and CA chain
// - /sign: PKCS12 trust store (no private key), signed cert, and CA chain
// This test checks encoder parameter handling and basic structure. Interoperability and encryption
// are covered in TestPKCS12OpenSSLValidation and TestPKCS12AndJKSJavaValidation.
func TestPathIssueSign_PKCS12Format(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// Import a CA certificate for test setup
	resp, err := CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
		"pem_bundle": ec256CaAndKey,
	})
	requireSuccessNonNilResponse(t, resp, err)

	// Get root cert data for tests with intermediate issuer
	issuers := resp.Data["imported_issuers"].([]string)
	issuerID := issuers[0]
	resp, err = CBRead(b, s, "issuer/"+issuerID)
	requireSuccessNonNilResponse(t, resp, err)
	rootData := resp.Data
	rootCert := rootData["certificate"].(string)

	// Create a role that allows issuing/signing certificates
	resp, err = CBWrite(b, s, "roles/test-role", map[string]interface{}{
		"allow_any_name": true,
		"max_ttl":        "2h",
		"key_type":       "ec",
		"key_bits":       "256",
	})
	requireSuccessNonNilResponse(t, resp, err)

	// Assert invalid format fails
	_, err = CBWrite(b, s, "issue/test-role", map[string]interface{}{
		"format":      "invalid",
		"common_name": "test.example.com",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), `the "format" parameter must be "pem", "der", "pem_bundle", "pkcs12_bundle" or "jks_bundle"`)

	// Parse PEM bundle and CSR for validation
	expectedRoot, err := certutil.ParsePEMBundle(ec256CaAndKey)
	require.NoError(t, err)
	block, _ := pem.Decode([]byte(ec256leafCsr))
	require.NotNil(t, block)
	expectedCSR, err := x509.ParseCertificateRequest(block.Bytes)
	require.NoError(t, err)

	buildData := func(password string, omitPassword bool, removeRoot bool, encoder string) map[string]interface{} {
		data := map[string]interface{}{
			"format":      "pkcs12_bundle",
			"common_name": "test.example.com",
		}
		if !omitPassword {
			data["pkcs12_password"] = password
		}
		if removeRoot {
			data["remove_roots_from_chain"] = true
		}
		if encoder != "" {
			data["pkcs12_encoder"] = encoder
		}
		return data
	}

	t.Run("issue", func(t *testing.T) {
		testCases := []struct {
			name         string
			encoder      string
			omitPassword bool
			password     string
			removeRoot   bool
			shouldError  bool
		}{
			{name: "custom password", password: "123-secure-password"},
			{name: "default password", password: pkcs12.DefaultPassword, omitPassword: true},
			{name: "empty password", password: ""},
			{name: "without CA chain", password: "123-secure-password", removeRoot: true},
			{name: "with modern2026 encoder", password: "123-secure-password", encoder: "modern2026"},
			{name: "with modern2023 encoder", password: "123-secure-password", encoder: "modern2023"},
			{name: "with invalid encoder", encoder: "modern2020", shouldError: true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				data := buildData(tc.password, tc.omitPassword, tc.removeRoot, tc.encoder)
				resp, err := CBWrite(b, s, "issue/test-role", data)
				pkcs12Bytes := verifyAndDecodePKCS12(t, "issue/test-role", resp, err, tc.shouldError)

				if !tc.shouldError {
					_, cert, caCerts := requireDecodesPKCS12Chain(t, pkcs12Bytes, tc.password)

					require.False(t, cert.IsCA, "leaf should not be CA")
					if tc.removeRoot {
						require.Len(t, caCerts, 0, "should have no CA cert when remove_roots_from_chain is true")
					} else {
						require.Len(t, caCerts, 1, "should have single cert in CA chain")
						require.True(t, caCerts[0].IsCA, "cert should be CA")
						require.Equal(t, expectedRoot.Certificate.Subject.CommonName, caCerts[0].Subject.CommonName)
					}
				}
			})
		}
	})

	t.Run("sign", func(t *testing.T) {
		testCases := []struct {
			name        string
			encoder     string
			endpoint    string
			removeRoot  bool
			shouldError bool
		}{
			{name: "sign", endpoint: "sign"},
			{name: "sign verbatim", endpoint: "sign-verbatim"},
			{name: "sign without CA chain", endpoint: "sign", removeRoot: true},
			{name: "with modern2026 encoder", endpoint: "sign", encoder: "modern2026"},
			{name: "with modern2023 encoder", endpoint: "sign", encoder: "modern2023"},
			{name: "with invalid encoder", endpoint: "sign", encoder: "modern2020", shouldError: true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				password := "123-secure-password"
				data := buildData(password, false, tc.removeRoot, tc.encoder)
				data["csr"] = ec256leafCsr
				path := tc.endpoint + "/test-role"
				resp, err := CBWrite(b, s, path, data)
				pkcs12Bytes := verifyAndDecodePKCS12(t, path, resp, err, tc.shouldError)

				if !tc.shouldError {
					certs := requireDecodesPKCS12TrustStore(t, pkcs12Bytes, password)
					require.False(t, certs[0].IsCA, "first cert should be leaf (not CA)")
					require.Equal(t, expectedCSR.RawSubjectPublicKeyInfo, certs[0].RawSubjectPublicKeyInfo, "first cert should contain CSR public key")
					require.Equal(t, expectedRoot.Certificate.Subject, certs[0].Issuer, "first should be issued by the expected CA")

					if tc.removeRoot {
						require.Len(t, certs, 1, "bundle should only contain leaf cert when remove_roots_from_chain is true")
					} else {
						require.Len(t, certs, 2, "bundle should contain leaf + CA")
						require.True(t, certs[1].IsCA, "last cert should be CA")
						require.Equal(t, expectedRoot.Certificate.Subject, certs[1].Subject, "last cert should match the expected CA certificate")
					}
				}
			})
		}
	})

	t.Run("with intermediate issuer", func(t *testing.T) {
		// Setup an intermediate, signed by the root.
		b_int, s_int := CreateBackendWithStorage(t)
		resp, err = CBWrite(b_int, s_int, "intermediate/generate/exported", map[string]interface{}{
			"common_name": "intermediate myvault.com",
			"key_type":    "ec",
		})
		requireSuccessNonNilResponse(t, resp, err)
		intermediateData := resp.Data
		resp, err = CBWrite(b, s, "root/sign-intermediate", map[string]interface{}{
			"csr":    intermediateData["csr"],
			"format": "pem",
		})
		requireSuccessNonNilResponse(t, resp, err)
		intermediateSignedData := resp.Data
		intermediateCert := intermediateSignedData["certificate"].(string)
		resp, err = CBWrite(b_int, s_int, "intermediate/set-signed", map[string]interface{}{
			"certificate": intermediateCert + "\n" + rootCert + "\n",
		})
		requireSuccessNonNilResponse(t, resp, err)
		// Setup role for signing certs
		resp, err = CBWrite(b_int, s_int, "roles/test-role", map[string]interface{}{
			"allow_any_name": true,
			"key_type":       "ec",
			"key_bits":       "256",
			"max_ttl":        "2h",
		})
		requireSuccessNonNilResponse(t, resp, err)

		testCases := []struct {
			name       string
			endpoint   string
			removeRoot bool
		}{
			{name: "sign", endpoint: "sign"},
			{name: "sign without root", endpoint: "sign", removeRoot: true},
			{name: "issue", endpoint: "issue"},
			{name: "issue without root", endpoint: "issue", removeRoot: true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				password := "123-secure-password"
				data := buildData(password, false, tc.removeRoot, "")
				if tc.endpoint == "sign" {
					data["csr"] = ec256leafCsr
				}
				path := tc.endpoint + "/test-role"
				resp, err := CBWrite(b_int, s_int, path, data)
				pkcs12Bytes := verifyAndDecodePKCS12(t, path, resp, err, false)

				if tc.endpoint == "issue" {
					_, cert, caCerts := requireDecodesPKCS12Chain(t, pkcs12Bytes, password)
					require.False(t, cert.IsCA, "leaf should not be CA")

					if tc.removeRoot {
						require.Len(t, caCerts, 1, "should have one CA cert when remove_roots_from_chain is true")
						require.True(t, caCerts[0].IsCA, "cert should be CA")
					} else {
						require.Len(t, caCerts, 2, "should have two certs in CA chain")
						require.True(t, caCerts[0].IsCA, "cert should be CA")
						require.True(t, caCerts[1].IsCA, "cert should be CA")
						require.Equal(t, expectedRoot.Certificate.Subject.CommonName, caCerts[1].Subject.CommonName)

					}
				}

				if tc.endpoint == "sign" {
					certs := requireDecodesPKCS12TrustStore(t, pkcs12Bytes, password)
					// Verify first cert is the leaf
					require.False(t, certs[0].IsCA, "first cert should be leaf (not CA)")
					require.Equal(t, expectedCSR.RawSubjectPublicKeyInfo, certs[0].RawSubjectPublicKeyInfo, "first cert should contain CSR public key")

					if tc.removeRoot {
						require.Len(t, certs, 2, "should have leaf + intermediate (no root)")
						require.True(t, certs[1].IsCA, "second cert should be CA")
						require.Equal(t, certs[1].Subject, certs[0].Issuer, "leaf should be issued by intermediate")
					} else {
						require.Len(t, certs, 3, "should have leaf + intermediate + root")
						require.True(t, certs[1].IsCA, "second cert should be CA")
						require.Equal(t, certs[1].Subject, certs[0].Issuer, "leaf should be issued by intermediate")
						require.True(t, certs[2].IsCA, "third cert should be CA")
						require.Equal(t, certs[2].Subject, certs[1].Issuer, "intermediate should be issued by root")
						requireSignedBy(t, certs[0], certs[1])
						requireSignedBy(t, certs[1], certs[2])
						// Verify root is self-signed
						requireSignedBy(t, certs[2], certs[2])
					}
				}
			})
		}
	})
}

// TestPathIssueSign_JKSFormat validates JKS (java keystore) output for /issue and /sign endpoints:
// - /issue: JKS contains private key, issued cert, and CA chain
// - /sign: JKS trust store (no private key), signed cert, and CA chain
// This test checks encoder parameter handling and basic structure.
// Interoperability is covered in TestPKCS12AndJKSJavaValidation.
func TestPathIssueSign_JKSFormat(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// Import a CA certificate for test setup
	resp, err := CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
		"pem_bundle": ec256CaAndKey,
	})
	requireSuccessNonNilResponse(t, resp, err)

	// Get root cert data for tests with intermediate issuer
	issuers := resp.Data["imported_issuers"].([]string)
	issuerID := issuers[0]
	resp, err = CBRead(b, s, "issuer/"+issuerID)
	requireSuccessNonNilResponse(t, resp, err)
	rootData := resp.Data
	rootCert := rootData["certificate"].(string)

	// Create a role that allows issuing/signing certificates
	resp, err = CBWrite(b, s, "roles/test-role", map[string]interface{}{
		"allow_any_name": true,
		"max_ttl":        "2h",
		"key_type":       "ec",
		"key_bits":       "256",
	})
	requireSuccessNonNilResponse(t, resp, err)

	// Parse PEM bundle and CSR for validation
	expectedRoot, err := certutil.ParsePEMBundle(ec256CaAndKey)
	require.NoError(t, err)
	block, _ := pem.Decode([]byte(ec256leafCsr))
	require.NotNil(t, block)
	expectedCSR, err := x509.ParseCertificateRequest(block.Bytes)
	require.NoError(t, err)

	buildData := func(password string, omitPassword bool, removeRoot bool, alias string) map[string]interface{} {
		data := map[string]interface{}{
			"format":      "jks_bundle",
			"common_name": "test.example.com",
		}
		if !omitPassword {
			data["jks_password"] = password
		}
		if removeRoot {
			data["remove_roots_from_chain"] = true
		}
		if alias != "" {
			data["jks_private_key_alias"] = alias
		}
		return data
	}

	t.Run("issue", func(t *testing.T) {
		testCases := []struct {
			name         string
			alias        string
			omitPassword bool
			password     string
			removeRoot   bool
		}{
			{name: "default password and alias", password: pkcs12.DefaultPassword, omitPassword: true},
			{name: "empty password", password: ""},
			{name: "without CA chain", password: "123-secure-password", removeRoot: true},
			{name: "with custom password and alias", password: "123-secure-password", alias: "myapp"},
			{name: "with custom numeric alias", password: "123-secure-password", alias: "5"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				data := buildData(tc.password, tc.omitPassword, tc.removeRoot, tc.alias)
				resp, err := CBWrite(b, s, "issue/test-role", data)
				jksBytes := verifyAndDecodeJKS(t, "issue/test-role", resp, err)

				expectedAlias := tc.alias
				if expectedAlias == "" {
					// Alias should be default if unset
					expectedAlias = "1"
				}

				_, cert, caCerts := requireDecodesJKSChain(t, jksBytes, tc.password, expectedAlias)
				require.False(t, cert.IsCA, "leaf should not be CA")
				if tc.removeRoot {
					require.Len(t, caCerts, 0, "should have no CA cert when remove_roots_from_chain is true")
				} else {
					require.Len(t, caCerts, 1, "should have single cert in CA chain")
					require.True(t, caCerts[0].IsCA, "cert should be CA")
					require.Equal(t, expectedRoot.Certificate.Subject.CommonName, caCerts[0].Subject.CommonName)
				}
			})
		}
	})

	t.Run("sign", func(t *testing.T) {
		testCases := []struct {
			name            string
			alias           string
			expectedAliases []string
			endpoint        string
			removeRoot      bool
		}{
			{name: "with defaults", endpoint: "sign", expectedAliases: []string{"1", "2"}},
			{name: "verbatim with defaults", endpoint: "sign-verbatim", expectedAliases: []string{"1", "2"}},
			{name: "without CA chain", endpoint: "sign", removeRoot: true, expectedAliases: []string{"1"}},
			// jks_private_key_alias should be ignored
			{name: "with numeric alias", endpoint: "sign", alias: "3", expectedAliases: []string{"1", "2"}},
			{name: "with non-numeric alias", endpoint: "sign", alias: "myapp", expectedAliases: []string{"1", "2"}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				password := "123-secure-password"
				data := buildData(password, false, tc.removeRoot, tc.alias)
				data["csr"] = ec256leafCsr
				path := tc.endpoint + "/test-role"
				resp, err := CBWrite(b, s, path, data)
				jksBytes := verifyAndDecodeJKS(t, path, resp, err)

				certs := requireDecodesJKSTrustStore(t, jksBytes, password, tc.expectedAliases)
				require.False(t, certs[0].IsCA, "first cert should be leaf (not CA)")
				require.Equal(t, expectedCSR.RawSubjectPublicKeyInfo, certs[0].RawSubjectPublicKeyInfo, "first cert should contain CSR public key")
				require.Equal(t, expectedRoot.Certificate.Subject, certs[0].Issuer, "first should be issued by the expected CA")

				if tc.removeRoot {
					require.Len(t, certs, 1, "bundle should only contain leaf cert when remove_roots_from_chain is true")
				} else {
					require.Len(t, certs, 2, "bundle should contain leaf + CA")
					require.True(t, certs[1].IsCA, "last cert should be CA")
					require.Equal(t, expectedRoot.Certificate.Subject, certs[1].Subject, "last cert should match the expected CA certificate")
				}
			})
		}
	})

	t.Run("with intermediate issuer", func(t *testing.T) {
		// Setup an intermediate, signed by the root.
		b_int, s_int := CreateBackendWithStorage(t)
		resp, err = CBWrite(b_int, s_int, "intermediate/generate/exported", map[string]interface{}{
			"common_name": "intermediate myvault.com",
			"key_type":    "ec",
		})
		requireSuccessNonNilResponse(t, resp, err)
		intermediateData := resp.Data
		resp, err = CBWrite(b, s, "root/sign-intermediate", map[string]interface{}{
			"csr":    intermediateData["csr"],
			"format": "pem",
		})
		requireSuccessNonNilResponse(t, resp, err)
		intermediateSignedData := resp.Data
		intermediateCert := intermediateSignedData["certificate"].(string)
		resp, err = CBWrite(b_int, s_int, "intermediate/set-signed", map[string]interface{}{
			"certificate": intermediateCert + "\n" + rootCert + "\n",
		})
		requireSuccessNonNilResponse(t, resp, err)
		// Setup role for signing certs
		resp, err = CBWrite(b_int, s_int, "roles/test-role", map[string]interface{}{
			"allow_any_name": true,
			"key_type":       "ec",
			"key_bits":       "256",
			"max_ttl":        "2h",
		})
		requireSuccessNonNilResponse(t, resp, err)

		testCases := []struct {
			name       string
			endpoint   string
			removeRoot bool
		}{
			{name: "sign", endpoint: "sign"},
			{name: "sign without root", endpoint: "sign", removeRoot: true},
			{name: "issue", endpoint: "issue"},
			{name: "issue without root", endpoint: "issue", removeRoot: true},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				password := "123-secure-password"
				data := buildData(password, false, tc.removeRoot, "")
				if tc.endpoint == "sign" {
					data["csr"] = ec256leafCsr
				}
				path := tc.endpoint + "/test-role"
				resp, err := CBWrite(b_int, s_int, path, data)
				jksBytes := verifyAndDecodeJKS(t, path, resp, err)

				if tc.endpoint == "issue" {
					_, cert, caCerts := requireDecodesJKSChain(t, jksBytes, password, "1")
					require.False(t, cert.IsCA, "leaf should not be CA")

					if tc.removeRoot {
						require.Len(t, caCerts, 1, "should have one CA cert when remove_roots_from_chain is true")
						require.True(t, caCerts[0].IsCA, "cert should be CA")
					} else {
						require.Len(t, caCerts, 2, "should have two certs in CA chain")
						require.True(t, caCerts[0].IsCA, "cert should be CA")
						require.True(t, caCerts[1].IsCA, "cert should be CA")
						require.Equal(t, expectedRoot.Certificate.Subject.CommonName, caCerts[1].Subject.CommonName)

					}
				}

				if tc.endpoint == "sign" {
					expectedAliases := []string{"1", "2", "3"}
					if tc.removeRoot {
						expectedAliases = []string{"1", "2"}
					}
					certs := requireDecodesJKSTrustStore(t, jksBytes, password, expectedAliases)
					// Verify first cert is the leaf
					require.False(t, certs[0].IsCA, "first cert should be leaf (not CA)")
					require.Equal(t, expectedCSR.RawSubjectPublicKeyInfo, certs[0].RawSubjectPublicKeyInfo, "first cert should contain CSR public key")

					if tc.removeRoot {
						require.Len(t, certs, 2, "should have leaf + intermediate (no root)")
						require.True(t, certs[1].IsCA, "second cert should be CA")
						require.Equal(t, certs[1].Subject, certs[0].Issuer, "leaf should be issued by intermediate")
					} else {
						require.Len(t, certs, 3, "should have leaf + intermediate + root")
						require.True(t, certs[1].IsCA, "second cert should be CA")
						require.Equal(t, certs[1].Subject, certs[0].Issuer, "leaf should be issued by intermediate")
						require.True(t, certs[2].IsCA, "third cert should be CA")
						require.Equal(t, certs[2].Subject, certs[1].Issuer, "intermediate should be issued by root")
						requireSignedBy(t, certs[0], certs[1])
						requireSignedBy(t, certs[1], certs[2])
						// Verify root is self-signed
						requireSignedBy(t, certs[2], certs[2])
					}
				}
			})
		}
	})
}

// RSA keys (and certificates to a lesser degree) are very slow to generate which is
// inconvenient for testing; so these have been pre-generated.  They are valid for
// 100 years; so shouldn't expire.  They were generated with the following code (and
// tuning the mount for a similarly long ttl):
//
//	resp, err := CBWrite(b, s, "root/generate/exported", map[string]interface{}{
//		"common_name": "foobar.com",
//		"issuer_name": parentKeyType.name,
//		"key_name":    parentKeyType.name,
//		"key_type":    parentKeyType.keyType,
//		"key_bits":    parentKeyType.keySize,
//		"ttl":         "1000000h",
//	})
const rsa2048CaAndKey = `-----BEGIN CERTIFICATE-----
MIIDNDCCAhygAwIBAgIURfqw7VetXbOAIRomrgKZvUPwvCUwDQYJKoZIhvcNAQEL
BQAwFTETMBEGA1UEAxMKZm9vYmFyLmNvbTAgFw0yNTExMTQyMDI3MjVaGA8yMTM5
MTIxNDEyMjc1NVowFTETMBEGA1UEAxMKZm9vYmFyLmNvbTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBALpBu3/YskXGW2WElPbXfUwibv9+vqMPlAud3iTh
an/RObRBMWK0heow4J6YY63z7PtF10XFxlwiXlv3HOrfNttvtuq4Ma1aC1JJDjZp
OAqniyoojtjOgH77TMDWEM1OIuX6vDwlgfoFySUjuSsrr1eK8yltFvTErqZJZ7RB
/RoHmzvEXeQdYBLV9gY+ku8kqDs7IL6lViKxnqQd33f5+2w4WvpOgVCZSTTyDEMZ
5kpLr/W9UaVx3Lp9WtWaEv2it36Bw+M3zjUxxidjoeHrQDzT1PCPA9mOVWT6s9uI
2WYfmbGAv79gWY3ziX/0XAED3txt/AKcYahBOy1AdrFa0zcCAwEAAaN6MHgwDgYD
VR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFNFSnZqM/NdZ
i+vmKgm7I858m1YgMB8GA1UdIwQYMBaAFNFSnZqM/NdZi+vmKgm7I858m1YgMBUG
A1UdEQQOMAyCCmZvb2Jhci5jb20wDQYJKoZIhvcNAQELBQADggEBAEJU2FLj94US
FpYek0ShL/M38YW0XM6NvAkA03YN1mzIfiGu122XWJwaG3BHGu/FOe3DNxQR65Cr
wT1ZHwzTuV6l54gUVbidGa3f5iscOvt+0PH+XXllVDfJymzUt+i67bC78bK3QEXs
bTIihjo1FH9XgTKTO5aFUgHiciBZqg5Df5REzoZM6s6HA5HWgjbyDTJICausUi+v
jQE1w2h8MrhMwmT+XopUWi2WAqqyTtenfETdA3xQciAZx9OU/4oy8uxjcH3q71/u
7DlylPrvCwOmO/hGUDR0pCZprwZ6fQ/OVYV9p+mtP5NIaRRJyg9bsJ3SyH1kzTGs
US17QT3YFNs=
-----END CERTIFICATE-----
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAukG7f9iyRcZbZYSU9td9TCJu/36+ow+UC53eJOFqf9E5tEEx
YrSF6jDgnphjrfPs+0XXRcXGXCJeW/cc6t8222+26rgxrVoLUkkONmk4CqeLKiiO
2M6AfvtMwNYQzU4i5fq8PCWB+gXJJSO5KyuvV4rzKW0W9MSupklntEH9GgebO8Rd
5B1gEtX2Bj6S7ySoOzsgvqVWIrGepB3fd/n7bDha+k6BUJlJNPIMQxnmSkuv9b1R
pXHcun1a1ZoS/aK3foHD4zfONTHGJ2Oh4etAPNPU8I8D2Y5VZPqz24jZZh+ZsYC/
v2BZjfOJf/RcAQPe3G38ApxhqEE7LUB2sVrTNwIDAQABAoIBACybsZxc+dVcPGeD
6Wl1Er05QfxPDrle8cYWeS28DxWttnRFaN6K/cepDSLuvHDdCtTjVTuQsoE+efrs
pDBcZXcIunZcxwkNl8iNVqoRaSqkFeBy9kNWsc+3wBovKrcBD7qk4pBFK2wGFrae
Z6q/O69rx/ET/3t/35RT4FJ7u3KQFtuvMujQqYyJvsmZ+GMqqqOS2Xfa9SJS+UDU
cUBNhseLtS4NeA/6gJaBzyFpBd/PTHxQUA/oct/xJr3mTD6whJQZ3vtDdx+GtwSe
lEvZdz12sCcMofbx74PWMefsiMRSFa6nRyrQdRNwjAqNj15bZtFOhlivtuiAu06i
dqkZ7gECgYEA7PMEwnbh9U2mwL2R3dlStUMZEsYWYUwh5QWlgu1pJ90VjPKgrZfJ
lu6l21fl3U23qlc+F19c6tSqHRf4euUugjXIp0h8E+SGc+MhadMgqGPLskUrPuqN
Wr4WV7MAlQqrfx7cFrU9PatugMiNRjVqp54cFU1otFrSoUPByWPFopECgYEAyTtX
Py2AArGPDli2EcWn+HpkYt5j8oip5J12wwU8a+BA3k+pB7G5PqdVwS0Xu5lUaJgG
do0vQZ02jqc8suSTeEZmVnQSU7Wx/+FDo6nWSS+0JVDjwLMuR80o5fmV01kwWjwN
/hPt4A/VrbjP1Bm11jh94NlgTrOaYUtj+Q2jbUcCgYAiNqjuR2ozGGZGmFjSlsm5
gJnDOzUKEYsnXZxbflpbtjGha3tF9Y/XKlhqhpObU9h8USKXD18ETXbOwqJPZH5F
sOxrMy0vViUP4LD3bdPeXKKR+CjZadbFToM9YIxp+ONwdI1E/iB8oh9PmyXDCH2A
/HSDouzGdgLJ5FW79ZsY8QKBgFZqlV0cPQzrE3QlxIp9R1T9un564pEU/2Cd/pJh
fUEWXMUbkIstV1AArGL46mg1wHnqT1w55UFYMkWwq/BnGK1eDjSyQ+yO6pHoOxPd
q5hiVApyYlwuloFfKWEZfa31bz5Q6/FgvZarNigUZavAHsaQG/6jWyhxGKsPpS8f
HD+hAoGAd5E+A9bHwdC1nIB48Ir7uv92jzqD6RTwuxYJteH9WAOSDiVZCGwRmJJT
ZkFDFj3Dr9T/45bJ8DUKv+CH5urX08jGRqCE6fAkaKKpJxcPsZxYtYEqdyoD3Vfl
2aql9ov9ut6BgP3yZhNGtJnW7HqYV2BV8xXPuPeEmSaaTAaioX4=
-----END RSA PRIVATE KEY-----`

const rsa3072CaAndKey = `-----BEGIN CERTIFICATE-----
MIIENDCCApygAwIBAgIUTQvhmoa/czFL4Y00RVeUxWWa7c4wDQYJKoZIhvcNAQEL
BQAwFTETMBEGA1UEAxMKZm9vYmFyLmNvbTAgFw0yNTExMTQyMjI1MDBaGA8yMTM5
MTIxNDE0MjUzMFowFTETMBEGA1UEAxMKZm9vYmFyLmNvbTCCAaIwDQYJKoZIhvcN
AQEBBQADggGPADCCAYoCggGBAKcaQt+AGRbBTkd5zwEeCeNBf72A74kBhs3R8BFK
0nBbbtjySZH05hDd5RNofuvW+zlr8m6sWBvyEZ3adNKhmGEiJyZtBqWMePI5p1yj
NyBZRd0M24uAPDtJYRnlKWkqBu9DIcKyHSGnqlqsQgcbRehs4Vn533WL5AhYd+Vx
EQF8fZjgZfMfJhhXE3Tnib3bS038njmxVMCMCT8+6V4tEGRZzK61FT0NqpDuXxYa
Zm57uEem8MY3TAtENtNrbNyyblQRdWZVmwpmI2oOJ9EMwTUU6Y8RdmuhNpsz4OO+
EWaFwGzg994H/kHda6vpyryxWDrAmwOFg5BxGiz+5KpUgtwzzq17rhEfVZhrdTMT
+NYFX52u3kFqqXXaTjLWGah96mNSC9FuORHGjwYVf7e5rcP2pQBHGfNmXfEVwnBX
tt6iJjp0sSr1ifrvdmcGBVR2mp4fYeqc8BW5EoUibRw+BPiGepikqz2U48rDTAPe
vQCsSDYU422qOA3FmmN8CvZ8IwIDAQABo3oweDAOBgNVHQ8BAf8EBAMCAQYwDwYD
VR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUANpvM+sId1IkzF1eSdBfyl19HFAwHwYD
VR0jBBgwFoAUANpvM+sId1IkzF1eSdBfyl19HFAwFQYDVR0RBA4wDIIKZm9vYmFy
LmNvbTANBgkqhkiG9w0BAQsFAAOCAYEAgUlvaCDvVqj6654dtFuMzoHcr6j8/GVb
xWloujxOmnVq9oCy1Pkhiv3g3iseBEU4ty19qbpTqhThJgbOYWagwz+gHQ9RPgaB
fFd4HzmjfQXhwW0+jJvLRGl5Qn3ydNasGPcxdhtEMVQtXUMYzX02JW8swqeSkk3D
xNnCAucxlcsK/NhmdFoNCebNMGegEkDOTJjVNbbPiuBSEzuySGbzvl+DGmUlshw4
O0UIyxE3+gq3BYT7HK1gQxyf6j3w/pxgO6yYoOWxMERruwTB/6omPXXB80rkQlsS
DzAIWqYOTGSKVLis+D80wL971eMd9fET7riHfXQVjLSvvVfM1If9jzrxs8Hc52TO
OZPgq3VWQLf0I/MhPGeoKfMVhqrX31TJzwJeI+UR3fsqzAZKTrnKeZ4m7WYntgHq
yV/3BMXxCWYljYv5k1mumbzk2uw17wC/KHB2ecTSIAAjMFKP5jTVf1atVy+aLyhO
SLBYd9gGLrRiV2zL61/g5HnYy+TlMoHe
-----END CERTIFICATE-----
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

const rsa4096CaAndKey = `-----BEGIN CERTIFICATE-----
MIIFNDCCAxygAwIBAgIUArsz+X8+6yiAov3a80qzEa7pF1EwDQYJKoZIhvcNAQEL
BQAwFTETMBEGA1UEAxMKZm9vYmFyLmNvbTAgFw0yNTExMTQyMjQyMDJaGA8yMTM5
MTIxNDE0NDIzMVowFTETMBEGA1UEAxMKZm9vYmFyLmNvbTCCAiIwDQYJKoZIhvcN
AQEBBQADggIPADCCAgoCggIBAM/ltvEBy4eEYhrzjhQd4k7xeRz86p3pdA3xcrba
B7jmxDsqw+UxQExEM6CFnTjJsPABSq3rPDiAOVwLM9NP5H356LZItA2ERaTvgCPg
vYt78bEiAF9Tis676NIUNh0u62ytxJnpD43jl6zSJugjIWEirNFuIiwMWUcXVkGK
Mt9QW7PtkNl7+VtVsR9VnIzBx/C/M4HzIRDIVa8m60vUbADA3u/PnI2MzIjh2ptA
dXFSMWxBK99o0E3DpzvQjjUh2B1d/UTupN6BwShWe+mVH466VuZKVtrN30Efm8xe
Tv4oZ6smbykQ/0bdVi68YiTBSfWzY1ROZ67SX9rA1hTNle6+7iXE51cJkhcGE8vi
QfcQXAe+2uttkxwrnl2qgkfL65FsWceYpL834Qb8UM3jlBI+EttDnansMP7M/+Y9
gu5GLK8+RwXIBqEb+GefsRHIeyUPEUZsHUUanrNXquYjMMSCeEfRe2941GYKVkPI
CKoo7MtTlJCymNpSuqCcBUcAcTje+j0JJyNYaFT2ZXC03Rf8bcrK1wryb7urOnfx
Fg1cGTJOasDFr8CH4G5KvwLcLwZfq8A/yp4jfyRbsQESDi8U7a8RFFXBrXkCLxgN
lO0YxKIJ0wDHpR0IYoyeYxVtVOP774ZGXd1s5VJSf0lke53+Ehzhjlf4P00xaNmb
IPT3AgMBAAGjejB4MA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MB0G
A1UdDgQWBBT3rlpZOrZu6SsSanJIDCDuaK36WzAfBgNVHSMEGDAWgBT3rlpZOrZu
6SsSanJIDCDuaK36WzAVBgNVHREEDjAMggpmb29iYXIuY29tMA0GCSqGSIb3DQEB
CwUAA4ICAQB0DGuWEBWDX5Kg7nayiKE7QvmiS7fDJstaXcnRbiPs1OitrCRKdDp+
QjtGzLQt/9DYLG+yPsxxLsRJxmTixq1qAFyG2Wka28F1QBAXfUiPBCE9cbHYN2Mk
M+Z9CMeuuZkXmY3wUG8uYOeNkXlrvshiCs6PZRMS6M7C8GL+4VeYld3kVuFqGeJT
p5kFwHTDy5kdqwrW0uLqL5VADPsxCvD1AeAio52J9zEMQoNSpOzb9kk4R3JM0Mse
eXIF6xol1F9IB12/XqjODe1IEFEZv+FVSM1JSyHpaO+vlQHw+OF7N7P7Dw/2YD0l
O1rqeLoLUIKNUuzzOTtoZjqNZ1irQHuN2+TwDJt72dGJDRh7XiRVxpZGRsH1DeH6
nig+VL2iCTPpv68Ge//eLRshRp/4dj0L7nleWw1ve830IOR7Sw4fpoILPFG7NECN
jLEU/hDb+8tsFJQF05m0X43SfxLq3ghexq/iIs09HvQ9AHsX1/Zz+mgyyvtI8e9/
wevNIplA5z3lVNjDUr91S2xZrGHk8zyAkQ925evvIeyiTOaW/zckMi4xyTVHkl0w
oFgquUU6saGjYcT9322ewvdoxsbYSlJR2Kzx0V96duSfY7DfxFgFoOh7oqzM/IJm
dUpagZCAn3I72DHt0uogRa5y3SpJwTlzCQhH1iH3bUFo6y0c/l3jIQ==
-----END CERTIFICATE-----
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

const rsa8192CaAndKey = `-----BEGIN CERTIFICATE-----
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
lb/Tm0GFPS0=
-----END CERTIFICATE-----
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

const ec224CaAndKey = `-----BEGIN CERTIFICATE-----
MIIBlDCCAUOgAwIBAgIUDIdGUQAa8UMKCxuIdBDQqMSccOMwCgYIKoZIzj0EAwIw
FTETMBEGA1UEAxMKZm9vYmFyLmNvbTAgFw0yNTExMTcxNDUwMDlaGA8yMTM5MTIx
NzA2NTAzOVowFTETMBEGA1UEAxMKZm9vYmFyLmNvbTBOMBAGByqGSM49AgEGBSuB
BAAhAzoABK/6iwAVzsJYvitjJXJ78YJrRXWExJD/kCmQehI+UrGPAPzHfp81RZuS
c6rHveBdJ3nXlzZCUHvXo3oweDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUw
AwEB/zAdBgNVHQ4EFgQU3jejWn+SRK8KMDn89eXMfiRgiWYwHwYDVR0jBBgwFoAU
3jejWn+SRK8KMDn89eXMfiRgiWYwFQYDVR0RBA4wDIIKZm9vYmFyLmNvbTAKBggq
hkjOPQQDAgM/ADA8AhwTua19D2CHI5Riv5WmBKz3uVLuSG8ZwYSX4ufwAhxDNrLX
2bMvnb8hnlEPI7gVhgB9bq0RGfpYE3v/
-----END CERTIFICATE-----
-----BEGIN EC PRIVATE KEY-----
MGgCAQEEHPybY8Z4tkbQV7HnbJ4Eq6AQpcccLxeA5fXoxbigBwYFK4EEACGhPAM6
AASv+osAFc7CWL4rYyVye/GCa0V1hMSQ/5ApkHoSPlKxjwD8x36fNUWbknOqx73g
XSd515c2QlB71w==
-----END EC PRIVATE KEY-----`

const ec256CaAndKey = `-----BEGIN CERTIFICATE-----
MIIBqTCCAU6gAwIBAgIUaRWkG+4WeYb1bQASXwBDq1hu0eMwCgYIKoZIzj0EAwIw
FTETMBEGA1UEAxMKZm9vYmFyLmNvbTAgFw0yNTExMTcxNTA0NTlaGA8yMTM5MTIx
NzA3MDUyOVowFTETMBEGA1UEAxMKZm9vYmFyLmNvbTBZMBMGByqGSM49AgEGCCqG
SM49AwEHA0IABPQf8311uaA/7ROV2vyjUGcyaJcc5YKshMg2VjDWvG8RrpUtvXLF
/VXW11zngxT97V8dpA1Lj1B1aaKIzx0CTfqjejB4MA4GA1UdDwEB/wQEAwIBBjAP
BgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBTd+NrxlY9Q4LisDjWcB574b1pn4jAf
BgNVHSMEGDAWgBTd+NrxlY9Q4LisDjWcB574b1pn4jAVBgNVHREEDjAMggpmb29i
YXIuY29tMAoGCCqGSM49BAMCA0kAMEYCIQDooDJ1p45cwYAUIwNYfU3HO1l1exor
pmNzZu+H0+d2tAIhAPalt/8lIeReeeaDcg2m0bsEBKpjm6VYCe4wWucQ7htW
-----END CERTIFICATE-----
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIIDoSaLJmKY4+qeeMMk7IjkvXoz9NOndpYa0NKlSGdJYoAoGCCqGSM49
AwEHoUQDQgAE9B/zfXW5oD/tE5Xa/KNQZzJolxzlgqyEyDZWMNa8bxGulS29csX9
VdbXXOeDFP3tXx2kDUuPUHVpoojPHQJN+g==
-----END EC PRIVATE KEY-----`

const ec384CaAndKey = `-----BEGIN CERTIFICATE-----
MIIB5jCCAWugAwIBAgIUKLfWtF+8L5myt+fvO9zmpRn6dxowCgYIKoZIzj0EAwMw
FTETMBEGA1UEAxMKZm9vYmFyLmNvbTAgFw0yNTExMTcxNzM1MzBaGA8yMTM5MTIx
NzA5MzYwMFowFTETMBEGA1UEAxMKZm9vYmFyLmNvbTB2MBAGByqGSM49AgEGBSuB
BAAiA2IABDd28Nf40cXJREW4BtC8Ig0dzingZ2jtU1pXS2edHrSQDclwBa8UYZe6
kgpkLNKeEYHhUlYPj98kxKl9E4ekjCn+CxJ5Zx8HPX3iPl01ORlrkI1ZcugNFSjP
aQ0kxB9JV6N6MHgwDgYDVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYD
VR0OBBYEFE61WPmDUVk+ClC/VZkarqvXKOPsMB8GA1UdIwQYMBaAFE61WPmDUVk+
ClC/VZkarqvXKOPsMBUGA1UdEQQOMAyCCmZvb2Jhci5jb20wCgYIKoZIzj0EAwMD
aQAwZgIxAKju3Se4UrkpUV59tcINouv/l24An3cHyeP2EPk8anRFyAe3P146lPqt
n1+B0Uzj+QIxALQMejPd/Mpps+CLMN13+fOigXLy6jsUXXTt3bUBvrfAC5udGOSC
crZno5EVVS04Wg==
-----END CERTIFICATE-----
-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDBUvahM6eA7z19A3p8Cny55ML83K0W+eIOF7nHLrlwaAlq6toNi3DoZ
SpcLHUOXS/CgBwYFK4EEACKhZANiAAQ3dvDX+NHFyURFuAbQvCINHc4p4Gdo7VNa
V0tnnR60kA3JcAWvFGGXupIKZCzSnhGB4VJWD4/fJMSpfROHpIwp/gsSeWcfBz19
4j5dNTkZa5CNWXLoDRUoz2kNJMQfSVc=
-----END EC PRIVATE KEY-----`

const ec512CaAndKey = `-----BEGIN CERTIFICATE-----
MIICLjCCAZGgAwIBAgIUeLTj4f4z04rgaEJSxBwXjdXxMlUwCgYIKoZIzj0EAwQw
FTETMBEGA1UEAxMKZm9vYmFyLmNvbTAgFw0yNTExMTcxNzM2NDRaGA8yMTM5MTIx
NzA5MzcxNFowFTETMBEGA1UEAxMKZm9vYmFyLmNvbTCBmzAQBgcqhkjOPQIBBgUr
gQQAIwOBhgAEAUr6WhqfUw43B0f6oWKkY3q83mYV2EE0zDhGlg5c1eIDDctD6Dby
jdWPwVjb9yyZ+1en2jveMgJqieN1vA0+Ov5WAb4f4JzIMsaCbEJl68riIh0zx+dg
O5hbisHR7+FWTg8HEQC0/EgSRlgcdBKQnv9tlozDxXHUBX1kQNsMdQoLidtco3ow
eDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUP794
iA+PSFLUTNYbzrB7o9ErmCswHwYDVR0jBBgwFoAUP794iA+PSFLUTNYbzrB7o9Er
mCswFQYDVR0RBA4wDIIKZm9vYmFyLmNvbTAKBggqhkjOPQQDBAOBigAwgYYCQSSi
X/IeqFFKhhyUEg757qkRnT1czdF6HmjaUrpCo5yblkVzHLTO8gxM4e8wDmxL0EFr
+Dxl/3OitC3Ro9D9BghzAkF2t3z3L353jnap1JEAX9gtXQVXoQFGaJ0gZ00iyx2W
pqa9U6pR8qMtDf9uGx01gINsqRruwDkprRXYBYcTU9NoWQ==
-----END CERTIFICATE-----
-----BEGIN EC PRIVATE KEY-----
MIHcAgEBBEIAvjYG3qoX52+P4HWKXNJlLVPwWysnmh/ABpZ+rA+VIyPujIYNfLem
QVaN9QMLKoaHudVh8IRFPOPGxlwOV4tKNeagBwYFK4EEACOhgYkDgYYABAFK+loa
n1MONwdH+qFipGN6vN5mFdhBNMw4RpYOXNXiAw3LQ+g28o3Vj8FY2/csmftXp9o7
3jICaonjdbwNPjr+VgG+H+CcyDLGgmxCZevK4iIdM8fnYDuYW4rB0e/hVk4PBxEA
tPxIEkZYHHQSkJ7/bZaMw8Vx1AV9ZEDbDHUKC4nbXA==
-----END EC PRIVATE KEY-----`

const ed25519CaAndKey = `-----BEGIN CERTIFICATE-----
MIIBaDCCARqgAwIBAgIUKsd3rEygnaf5fiWZUwyh2WlRQxAwBQYDK2VwMBUxEzAR
BgNVBAMTCmZvb2Jhci5jb20wIBcNMjUxMTE3MTc0MDE0WhgPMjEzOTEyMTcwOTQw
NDRaMBUxEzARBgNVBAMTCmZvb2Jhci5jb20wKjAFBgMrZXADIQBpSFBQV9cYAZQT
S0xslNFjdbpgH5rC5yEgUYdgVFWHoqN6MHgwDgYDVR0PAQH/BAQDAgEGMA8GA1Ud
EwEB/wQFMAMBAf8wHQYDVR0OBBYEFHNSUSvGUPM7WKBVhRvojQBnoU35MB8GA1Ud
IwQYMBaAFHNSUSvGUPM7WKBVhRvojQBnoU35MBUGA1UdEQQOMAyCCmZvb2Jhci5j
b20wBQYDK2VwA0EAuNsyYppYWemshxpBinChoxdahOUx8G6uozUSY0VN8HTHcVRP
KhhRZO3a2AzTKjpC/GxD3lrZbPjGpb8LxvEpCA==
-----END CERTIFICATE-----
-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIDALK7hTzLj3lAbHDGIBVGMRI7UnrXnUZaX/wQJABEWj
-----END PRIVATE KEY-----`

const mldsa44CaAndKey = `-----BEGIN CERTIFICATE-----
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
AAAAABIfKzw=
-----END CERTIFICATE-----
-----BEGIN PRIVATE KEY-----
MDQCAQAwCwYJYIZIAWUDBAMRBCKAIDZvtfIFECLbnz60XlOZIPWN1vdiLUMxxww5
s9XXrmUE
-----END PRIVATE KEY-----`

const mldsa65CaAndKey = `-----BEGIN CERTIFICATE-----
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
Jw==
-----END CERTIFICATE-----
-----BEGIN PRIVATE KEY-----
MDQCAQAwCwYJYIZIAWUDBAMSBCKAIEwDRQ81ZvR9MdhGX7eGmEY3quBZGeefV6jY
aGGDmHg7
-----END PRIVATE KEY-----`

const mldsa87CaAndKey = `-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----
-----BEGIN PRIVATE KEY-----
MDQCAQAwCwYJYIZIAWUDBAMTBCKAIGfeqFMmBHMKKSdULQmt3GrPjFwKkEDeoRfq
MCjV7UMr
-----END PRIVATE KEY-----`

// Generating CSRs can take a great deal of time, particularly when
// it requires generating long RSA keys.  Therefore, the following
// CSR were generated with this code, but have been cached here for
// testing.
//
//	goodCr := &x509.CertificateRequest{}
//	var csrKey any
//	var csrPem string
//	for _, childKeyTypeOption := range keyTypeOptions {
//		switch childKeyTypeOption.keyType {
//		case "rsa":
//			csrKey, err = rsa.GenerateKey(rand.Reader, childKeyTypeOption.keySize)
//		case "ec":
//			switch childKeyTypeOption.keySize {
//			case 224:
//				csrKey, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
//			case 256:
//				csrKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
//			case 384:
//				csrKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
//			case 521:
//				csrKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
//			}
//		case "ed25519":
//			_, csrKey, err = ed25519.GenerateKey(rand.Reader)
//		}
//		require.NoError(t, err, "failed generated key for CSR")
//		csr, err := x509.CreateCertificateRequest(rand.Reader, goodCr, csrKey)
//		require.NoError(t, err, "failed generating csr")
//
//		csrPem = strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{
//			Type:  "CERTIFICATE REQUEST",
//			Bytes: csr,
//		})))
const rsa2048leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIICRTCCAS0CAQAwADCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANKQ
24KKkOQ+JBhO+s322v4nbioOGVXV7h9YSwM8EIQnZg6N9QxfLTtnx78rfR73yI34
2v+eucxhWmixoFs2/JEye7BiJ5INV9j3eUXQfzmmM9OuLm7RELD+vXcwGKS3so9x
wz57vhttaNM73VylNQCJvoftbhLcKHdkDzqb8fTcDhP+m1v7VKGXtMFF/4wgiJDp
kV3EZX2LHeZ7GcHweIW0m+JgCJBmHH6+mCnhIEiZXemzpJ6MtVLCt/hXdBjUzK/e
LrGkKG9QJEO7HOZKc/KAtVfrROL3PJFIrqJWaD8UWqmuTF5JQOLGqG038QRE+Zcw
0woON+Ssq/Ncl5X5wXsCAwEAAaAAMA0GCSqGSIb3DQEBCwUAA4IBAQAUfcsyYePl
8LdE6V0HVEYLVCfN4n+MVAjWBUjAWJZdimfxTGqMaAUBcvuhyKI6kYRcIXTgptQm
DaXK0BX6EjEbindmA1BpOEWgZyb+heE2OV5VtMKyX/+fnJ2ZJ3VY7aAACsx91VNJ
9pMdGKrcKewP6E2pbS4yryx8NbbZRLR46LBFwfFIcoPfFgrSGmv6mjnNNzLXdF+w
nKJyPhjP5c3Aq54kkDwzu68JEZxapwL3Zm49sWoNb0gnkm/XtgkkRrCYamsUnqqq
LSjy4Ov7yvK31HAULfkyxXgVxNNdxAsyUGXT/pVc6F9dtqq8Ni0s/KmvKk+ghUi5
w49o1WDkKU1p
-----END CERTIFICATE REQUEST-----`

const rsa3072leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIIDRTCCAa0CAQAwADCCAaIwDQYJKoZIhvcNAQEBBQADggGPADCCAYoCggGBALXT
IIejmww8yIh3/ly7d+/B5xowUhk9ij4sg5LAyxQtQa3ZWYNjipdGBG14G0nspG0e
Z1pevHDuhmODH7xi6shwo+ng8phTDYoW5/9aJn10lyHBsvuCic6kRuSZTqpS3NFk
8l2y4vHKYS2ZhcGlBujO4m/cbafkUN7v/a4L0h2fYq7EG/rGRud79OMnFaVhFqaR
macqukHqzhxo9Wj2y4iqdDiQpL0mm84ymKE27fSdfBDdkATKN3mUwfkkqBDM4bhM
ulX0saqjqxLB1IA2NBmt3jKO0HH36e7UlwDmBnEi2dpdMwvzLntCcAECUqSEDrZr
Ht+QZtBOinQ/1XWpYeCSV2XtO4dFfooM2/dSOFm3IK1nohVsIDyzFlRGAPHbT8xC
GbkxY53RWcM/sAKQsnx5kCdefrtQEX1DSOy1w1xEsHvJi1ptuCU0wLZNyq8f9nnr
RRoAFUJmCnjBm+KUOT4p6JVovwbr9I6/5DIq/iVqtXMVzx6NFxDCbv9ok8PXhwID
AQABoAAwDQYJKoZIhvcNAQELBQADggGBAC8raimxwX2lkOxS2Wo3vBYd+VOYD1z3
XXSfCdgqix0tscL09Sz1w0SOttKIoLpgZxMa7/p0fYtjpj1Y4SR2ee54kgRxdHHY
6VGPd6eVDfpYhb+6TEIBdZS2gIC150gOAA7vDrsIVXnLwjTOg6JwQbFgMvthVZ5l
ne+Y2MmSdO1EmKNwYxM3O80HDRIrcnU84OtLxxWGZRFGOaIwC8vFtArnHdHQDw/V
hsgFMov+00mcPUYrLHUy3kLZRm/PEwASgoy6CRhI5UM0LRARV8u+BUMX6B8pjri2
/eRm4euRg4t9jig+jU9oNM+coFKl8ZaHgwoDpmHiUjhGDAprZRwgYB2IhGqbG6TS
RmWkuzQD31uONdTrSIoVcDeVzREmzVDfERQ8KoQQ2nrGqlXgk2QUUmfHt6L+kpZ+
aG/m3iQLttdhotO9DMeZe0qkQ7bgwOnKdOhU9LdjeH1wBHqtKh2IoG+uZ8UOQp+f
uQeKKhf1NSEOyv81mrUvhdVD//6tspRzwg==
-----END CERTIFICATE REQUEST-----`

const rsa4096leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIIERTCCAi0CAQAwADCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAL93
Fn32HImbEe9aD6UqfLIM3VB5pM9LgehMbZUo1NBdkM5urpEe9fG6mP4Or2b3hvn/
HJmV0si4V+IU9DTfc3rmVRUm8weQ5HTyF5pqmbQWbiPpibV6t3BNCuWfLJEfhCW2
pPhyhyTriKIffuo994N3VRW+TXdGFsPvlPAt91miR/DH5GfNZAmDqlxJhgrT70dP
H1eE3PBwPhH8xaAWiccsF5/4R4idDn5DusoB4CWZXSx+FDQdMNTV7XSDOtgZax0R
RKo0nwipBJiIulgp0SsSu4zy8fN2G9EjKkROjiX1YWomqwiK7nrsjjuR5z63Qnwk
pzSGIQteTHSTFH8XoWQQa91pUUy22j3DKwcT9gJEptzt7NiDsbLdqR52dQRyfuYJ
UNDrbB0aC+k2rVBeHPEzvvm3rOJpS6dq+5IU/CCroi0KMt2+ZumMHDe2vGsuNvwS
NvWd4ibk/EZWtiOUTISuEVBOJvtaPGtNmd25CdC/7EN+PpM8vCkW8my0O/hN7frB
TOTR8KTc6LmgBPVXRMUyB6Gss2fI86m38ydZioLGtNCT2ZkAmZzp1w8YNsqdnGvy
kwDIgtjvFR1XtvGv2XO20hIDd7BjJAgDBDoRfU2gpA3RUpQZcyAK4JNexG//gQAV
1cwvInjPSeGsLEhsz9V07mt+OHKEyDJABUecnBWnAgMBAAGgADANBgkqhkiG9w0B
AQsFAAOCAgEAH7PJi7vGU7eghKMeM6Y4r5yYUcQC665AUsgKPax/pyX0E/blPZyL
RhTupB6y4G6He+l4EuoAUcVBx3k/uFHy4jsIElSo7HyqZi9UtUzOn448ZQbGg5fU
zC2qJp+KKEBYnpuAul2e0YDklklOk4oSsuXAb/gYEhiN5foqtn5hUme9AoXBRwyS
hGmJLwkZaJIFDzIlvG1NNfhNnzfTR4o1t/7YMdcp7sJ+Qa/IShc++Xr3vgE6dUmY
3aGQIlWeIkO6qxsuKcfKGt4wwkCRHPR981wT6xHEl3lgeAJg4qalFcCmxOXaWhz+
tppB/Hj1e6FYqfnCBkPe0jjLrenGDYeoSKQkXH4XK2I2VX+kGUPnYN5yaNP9N6hL
31dhrPHknB13mx3UQoGCD+tRglelHceLbgTZteiVC/PjKRN9ZJit98KalHOKQhOv
AFPtQXbBKn3jbIKwYxJE37e3HqQ69eBudFzo0bUwDvEgAsbZXmh6D5hKpEIMRDOd
usCQjur8/Mjwn/c+VI89wAxnBevlFvj6Sy3EDtnSUmnEJsQ8Qhrg2LwISkcodUqg
HUJkn7Yn06drnCTa4WkBBTJOisg39ZvQMKPnjMCxG80CiqgOOZztFH6YeDWfux6b
Tw/AULSSqwxWnEzEwr9CWpeZbgNbXrgGW/iQnCpEwH/X08PsYqSyYIs=
-----END CERTIFICATE REQUEST-----`

const rsa8192leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
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
THW7sEerAulu
-----END CERTIFICATE REQUEST-----`

const ec224leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIGnMFcCAQAwADBOMBAGByqGSM49AgEGBSuBBAAhAzoABDPEiaV7GSzlRPZoY1IE
gbhxZHhP9vqReTNRfAZSCzIjrIAb46j8HvpfruQUCiXnc9F+RDE90LRuoAAwCgYI
KoZIzj0EAwIDQAAwPQIdAOCbImu4+9pUb8BgbHXa13eY3HaKsNdc8m3bizQCHGB0
xwWb7zUBGOrVVyqs/vFhBLq16OTzw3b/wS0=
-----END CERTIFICATE REQUEST-----`

const ec256leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIG7MGICAQAwADBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABNBo3srjDRR9pIKp
Xkg4A4cYLZIdVzY/atAevIWb8tNtIwPu14VOFs1C01YwxQS2TY0iSwgUyAo+BIws
01nGKc6gADAKBggqhkjOPQQDAgNJADBGAiEA76gHPJbAOX85DHwWsgzn+9GpV0t7
qX98i/SIM6z2yIACIQCZfb8cCcx3PL34RzOw4WGDero7QOUvQ9DQp+w036KJRg==
-----END CERTIFICATE REQUEST-----`

const ec384leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIH2MH8CAQAwADB2MBAGByqGSM49AgEGBSuBBAAiA2IABN0yemr+Ij23X7Y1SYbJ
Naf5y2gYlfi9lDG4JG9dcGJ+rc1XQKJ8zesfODjnf0uvJ+tM4ujhZsmz7kF4MKu0
pCdmfF8TfWdvE7ys0wuRlxOgngDsx8sQSbweoFEk1hrRfaAAMAoGCCqGSM49BAMD
A2cAMGQCMFQlHDMjwsKYpfxPmJQDeErzjJynPy+8KAWiJ3zMj39Otdrx0md+SyBG
ibjgw3uXTQIwMdYhMKTEFtF36nqtNTtPODc5MMySHRfxB9Sy0TnjdTyGvdvrHzSr
AWEW03wZUpbW
-----END CERTIFICATE REQUEST-----`

const ec512leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MIIBQzCBpQIBADAAMIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQBctfkWcBYsL0j
xDocFiSLCLujQoMokv+1wBc+J9oWfmYFatpqdd1OlS2A8UdaLc8HIZCIaeV6rUBy
7/LrqZI1Zx0BTZt7Yl3KimUgCLkrq11WLKlQdxuc59ejFUtQR4ci1sR63MGenPgp
/aWUETKoQ8O/Xvur3nkHtVoFD9lmmD2PEeCgADAKBggqhkjOPQQDBAOBjAAwgYgC
QgFpswcXJ+bpXDihwtExKPTTwRIVv7t0JMHgQjolIfSf6T20P9KuuTIYKuls87y9
DGQo5Ku/tj6PkZFwO6VRsbgn5QJCAI6feCOcmwttGXN6YzoajN4458oyD+/UuQ9/
RnoT+li4M/y1QETDjT0h7IIf5lHYulzI+rP7VCE9IU3PS3CeYrdx
-----END CERTIFICATE REQUEST-----`

const ed25519leafCsr = `-----BEGIN CERTIFICATE REQUEST-----
MH8wMwIBADAAMCowBQYDK2VwAyEAmX98EACbAIG9PNcwMG6Zo5rkLDrQ02n8ZEXs
A4AbUdKgADAFBgMrZXADQQCmxK2WBPw5S7Bs2xHmBgG3+yF5HxPjVrVloRWxY9/O
vaeNkZybMZSaATniE3IAG04xnddeDR3MZZwDuUamuZcA
-----END CERTIFICATE REQUEST-----`

const mldsa44LeafCsr = `-----BEGIN CERTIFICATE REQUEST-----
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
-----END CERTIFICATE REQUEST-----`

const mldsa65LeafCsr = `-----BEGIN CERTIFICATE REQUEST-----
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
-----END CERTIFICATE REQUEST-----`

const mldsa87LeafCsr = `-----BEGIN CERTIFICATE REQUEST-----
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
AAAAAAAAAAAAAAAAAAAAAAAABQsSFx4jKi4=
-----END CERTIFICATE REQUEST-----`
