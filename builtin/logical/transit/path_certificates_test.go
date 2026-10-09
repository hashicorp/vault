// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	cryptoRand "crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/pki"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

var templateCsr = `
-----BEGIN CERTIFICATE REQUEST-----
MIICRTCCAS0CAQAwADCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAM49
McW7u3ILuAJfSFLUtGOMGBytHmMFcjTiX+5JcajFj0Uszb+HQ7eIsJJNXhVc/7fg
Z01DZvcCqb9ChEWE3xi4GEkPMXay7p7G1ooSLnQp6Z0lL5CuIFfMVOTvjfhTwRaJ
l9v2mMlm80BeiAUBqeoyGVrIh5fKASxaE0jrhjAxhGzqrXdDnL8A4na6ArprV4iS
aEAziODd2WmplSKgUwEaFdeG1t1bJf3o5ZQRCnKNtQcAk8UmgtvFEO8ohGMln/Fj
O7u7s6iRhOGf1g1NCAP5pGqxNx3bjz5f/CUcTSIGAReEomg41QTIhD9muCTL8qnm
6lS87wkGTv7qbeIGB7sCAwEAAaAAMA0GCSqGSIb3DQEBCwUAA4IBAQAfjE+jNqIk
4V1tL3g5XPjxr2+QcwddPf8opmbAzgt0+TiIHcDGBAxsXyi7sC9E5AFfFp7W07Zv
r5+v4i529K9q0BgGtHFswoEnhd4dC8Ye53HtSoEtXkBpZMDrtbS7eZa9WccT6zNx
4taTkpptZVrmvPj+jLLFkpKJJ3d+Gbrp6hiORPadT+igLKkqvTeocnhOdAtt427M
RXTVgN14pV3tqO+5MXzNw5tGNPcwWARWwPH9eCRxLwLUuxE4Qu73pUeEFjDEfGkN
iBnlTsTXBOMqSGryEkmRaZslWDvblvYeObYw+uc3kCbJ7jRy9soVwkbb5FueF/yC
O1aQIm23HrrG
-----END CERTIFICATE REQUEST-----
`

func TestTransit_Certs_CreateCsr(t *testing.T) {
	for _, keyType := range []string{
		"rsa-2048", "rsa-3072", "rsa-4096",
		"ecdsa-p256", "ecdsa-p384", "ecdsa-p521",
		"ed25519", "aes256-gcm96",
	} {
		b, s := createBackendWithStorage(t)

		resp, err := b.HandleRequest(context.Background(), &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "keys/test-key",
			Storage:   s,
			Data:      map[string]interface{}{"type": keyType},
		})
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("resp: %#v\nerr: %v", resp, err)
		}

		testTransit_CreateCsr(t, b, s, "test-key", keyType, templateCsr, 0)
	}
}

func testTransit_CreateCsr(t *testing.T, b logical.Backend, s logical.Storage, keyName, keyType, pemTemplateCsr string, keyVersion int) {
	t.Helper()

	csrSignReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("keys/%s/csr", keyName),
		Storage:   s,
		Data: map[string]interface{}{
			"csr": pemTemplateCsr,
		},
	}

	if keyVersion != 0 {
		csrSignReq.Data["version"] = keyVersion
	}

	resp, err := b.HandleRequest(context.Background(), csrSignReq)

	switch keyType {
	case "rsa-2048", "rsa-3072", "rsa-4096", "ecdsa-p256", "ecdsa-p384", "ecdsa-p521", "ed25519":
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("failed to sign CSR, err:%v resp:%#v", err, resp)
		}

		signedCsrBytes, ok := resp.Data["csr"]
		if !ok {
			t.Fatal("expected response data to hold a 'csr' field")
		}

		signedCsr, err := parseCsr(signedCsrBytes.(string))
		if err != nil {
			t.Fatalf("failed to parse returned csr, err:%v", err)
		}

		templateCsr, err := parseCsr(pemTemplateCsr)
		if err != nil {
			t.Fatalf("failed to parse returned template csr, err:%v", err)
		}

		// NOTE: Check other fields?
		if !reflect.DeepEqual(signedCsr.Subject, templateCsr.Subject) {
			t.Fatalf("subjects should have matched, err:%v", err)
		}

	default:
		if err == nil || (resp != nil && !resp.IsError()) {
			t.Fatalf("should have failed to sign CSR, provided key type does not support signing")
		}
	}
}

func TestTransit_Certs_ImportCertChain(t *testing.T) {
	// Create Cluster
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": Factory,
			"pki":     pki.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cores := cluster.Cores
	client := cores[0].Client

	// Mount transit backend
	err := client.Sys().Mount("transit", &api.MountInput{
		Type: "transit",
	})
	require.NoError(t, err)

	// Mount PKI backend
	err = client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
	})
	require.NoError(t, err)

	for _, keyType := range []string{
		"rsa-2048", "rsa-3072", "rsa-4096",
		"ecdsa-p256", "ecdsa-p384", "ecdsa-p521",
		"ed25519",
	} {
		_, err := client.Logical().Write(fmt.Sprintf("transit/keys/%s", keyType), map[string]interface{}{
			"type": keyType,
		})
		require.NoError(t, err)

		testTransit_ImportCertChain(t, client, keyType, keyType, 0)
	}
}

func testTransit_ImportCertChain(t *testing.T, apiClient *api.Client, keyType, keyName string, keyVersion int) {
	t.Helper()
	id, err := uuid.GenerateUUID()
	require.NoError(t, err)
	issuerName := fmt.Sprintf("%s-issuer-%s", keyType, id)

	// Setup a new CSR
	privKey, err := cryptoutil.GenerateRSAKey(cryptoRand.Reader, 3072)
	require.NoError(t, err)

	var csrTemplate x509.CertificateRequest
	csrTemplate.Subject.CommonName = "example.com"
	reqCsrBytes, err := x509.CreateCertificateRequest(cryptoRand.Reader, &csrTemplate, privKey)
	require.NoError(t, err)

	pemTemplateCsr := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: reqCsrBytes,
	})
	t.Logf("csr: %v", string(pemTemplateCsr))

	reqData := map[string]interface{}{
		"csr": string(pemTemplateCsr),
	}

	if keyVersion != 0 {
		reqData["version"] = keyVersion
	}

	// Create CSR from template CSR fields and key in transit
	resp, err := apiClient.Logical().Write(fmt.Sprintf("transit/keys/%s/csr", keyName), reqData)
	require.NoError(t, err)
	require.NotNil(t, resp)
	pemCsr := resp.Data["csr"].(string)

	// Generate PKI root
	resp, err = apiClient.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"issuer_name": issuerName,
		"common_name": "PKI Root X1",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	rootCertPEM := resp.Data["certificate"].(string)
	pemBlock, _ := pem.Decode([]byte(rootCertPEM))
	require.NotNil(t, pemBlock)

	rootCert, err := x509.ParseCertificate(pemBlock.Bytes)
	require.NoError(t, err)

	// Create role to be used in the certificate issuing
	resp, err = apiClient.Logical().Write("pki/roles/example-dot-com", map[string]interface{}{
		"issuer_ref":                         issuerName,
		"allowed_domains":                    "example.com",
		"allow_bare_domains":                 true,
		"basic_constraints_valid_for_non_ca": true,
		"key_type":                           "any",
	})
	require.NoError(t, err)

	// Sign the CSR
	resp, err = apiClient.Logical().Write("pki/sign/example-dot-com", map[string]interface{}{
		"issuer_ref": issuerName,
		"csr":        pemCsr,
		"ttl":        "10m",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	leafCertPEM := resp.Data["certificate"].(string)
	pemBlock, _ = pem.Decode([]byte(leafCertPEM))
	require.NotNil(t, pemBlock)

	leafCert, err := x509.ParseCertificate(pemBlock.Bytes)
	require.NoError(t, err)

	require.NoError(t, leafCert.CheckSignatureFrom(rootCert))
	t.Logf("root: %v", rootCertPEM)
	t.Logf("leaf: %v", leafCertPEM)

	certificateChain := strings.Join([]string{leafCertPEM, rootCertPEM}, "\n")
	reqData = map[string]interface{}{
		"certificate_chain": certificateChain,
	}
	if keyVersion != 0 {
		reqData["version"] = keyVersion
	}

	// Import certificate chain to transit key version
	resp, err = apiClient.Logical().Write(fmt.Sprintf("transit/keys/%s/set-certificate", keyName), reqData)
	require.NoError(t, err)
	require.NotNil(t, resp)

	resp, err = apiClient.Logical().Read(fmt.Sprintf("transit/keys/%s", keyName))
	require.NoError(t, err)
	require.NotNil(t, resp)
	keys, ok := resp.Data["keys"].(map[string]interface{})
	if !ok {
		t.Fatalf("could not cast Keys value")
	}
	var expectedVer int
	if keyVersion == 0 {
		latestVerRaw, ok := resp.Data["latest_version"].(json.Number)
		require.True(t, ok)
		latestVer, err := latestVerRaw.Int64()
		require.NoError(t, err)
		expectedVer = int(latestVer)
	} else {
		expectedVer = keyVersion
	}
	keyData, ok := keys[strconv.Itoa(expectedVer)].(map[string]interface{})
	if !ok {
		t.Fatalf("could not cast key version %d from keys", expectedVer)
	}
	require.NotEmpty(t, keyData["certificate_chain"])
}

func TestTransit_Certs_ImportInvalidCertChain(t *testing.T) {
	// Create Cluster
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": Factory,
			"pki":     pki.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cores := cluster.Cores
	client := cores[0].Client

	// Mount transit backend
	err := client.Sys().Mount("transit", &api.MountInput{
		Type: "transit",
	})
	require.NoError(t, err)

	// Mount PKI backend
	err = client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
	})
	require.NoError(t, err)

	for _, keyType := range []string{
		"rsa-2048", "rsa-3072", "rsa-4096",
		"ecdsa-p256", "ecdsa-p384", "ecdsa-p521",
		"ed25519",
	} {
		_, err := client.Logical().Write(fmt.Sprintf("transit/keys/%s", keyType), map[string]interface{}{
			"type": keyType,
		})
		require.NoError(t, err)

		testTransit_ImportInvalidCertChain(t, client, keyType, keyType, "")
	}
}

func testTransit_ImportInvalidCertChain(t *testing.T, apiClient *api.Client, keyType, keyName string, version string) {
	t.Helper()
	id, err := uuid.GenerateUUID()
	require.NoError(t, err)
	issuerName := fmt.Sprintf("%s-issuer-%s", keyType, id)

	// Generate PKI root
	resp, err := apiClient.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"issuer_name": issuerName,
		"common_name": "PKI Root X1",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	rootCertPEM := resp.Data["certificate"].(string)
	pemBlock, _ := pem.Decode([]byte(rootCertPEM))
	require.NotNil(t, pemBlock)

	rootCert, err := x509.ParseCertificate(pemBlock.Bytes)
	require.NoError(t, err)

	pkiKeyType := "rsa"
	pkiKeyBits := "0"
	if strings.HasPrefix(keyType, "rsa") {
		pkiKeyBits = keyType[4:]
	} else if strings.HasPrefix(keyType, "ecdsa") {
		pkiKeyType = "ec"
		pkiKeyBits = keyType[7:]
	} else if keyType == "ed25519" {
		pkiKeyType = "ed25519"
		pkiKeyBits = "0"
	}

	// Create role to be used in the certificate issuing
	resp, err = apiClient.Logical().Write("pki/roles/example-dot-com", map[string]interface{}{
		"issuer_ref":                         issuerName,
		"allowed_domains":                    "example.com",
		"allow_bare_domains":                 true,
		"basic_constraints_valid_for_non_ca": true,
		"key_type":                           pkiKeyType,
		"key_bits":                           pkiKeyBits,
	})
	require.NoError(t, err)

	// XXX -- Note subtle error: we issue a certificate with a new key,
	// not using a CSR from Transit.
	resp, err = apiClient.Logical().Write("pki/issue/example-dot-com", map[string]interface{}{
		"common_name": "example.com",
		"issuer_ref":  issuerName,
		"ttl":         "10m",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	leafCertPEM := resp.Data["certificate"].(string)
	pemBlock, _ = pem.Decode([]byte(leafCertPEM))
	require.NotNil(t, pemBlock)

	leafCert, err := x509.ParseCertificate(pemBlock.Bytes)
	require.NoError(t, err)

	require.NoError(t, leafCert.CheckSignatureFrom(rootCert))
	t.Logf("root: %v", rootCertPEM)
	t.Logf("leaf: %v", leafCertPEM)

	certificateChain := strings.Join([]string{leafCertPEM, rootCertPEM}, "\n")

	// Import certificate chain to transit key version
	resp, err = apiClient.Logical().Write(fmt.Sprintf("transit/keys/%s/set-certificate", keyName), map[string]interface{}{
		"certificate_chain": certificateChain,
		"version":           version,
	})
	require.Error(t, err)
}

// TestTransit_Certs_CreateCsr_PreservesNonExtensionAttributes checks that
// parseCsr strips extensionRequest from Attributes without touching other
// attributes like challengePassword, and that challengePassword survives
// end-to-end re-signing through Transit.
func TestTransit_Certs_CreateCsr_PreservesNonExtensionAttributes(t *testing.T) {
	t.Parallel()
	for _, keyType := range []string{"rsa-2048", "ecdsa-p256", "ed25519"} {
		keyType := keyType
		t.Run(keyType, func(t *testing.T) {
			t.Parallel()
			b, s := createBackendWithStorage(t)

			resp, err := b.HandleRequest(context.Background(), &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "keys/test-key",
				Storage:   s,
				Data:      map[string]interface{}{"type": keyType},
			})
			require.NoError(t, err)
			require.False(t, resp != nil && resp.IsError(), "key creation failed: %v", resp)

			testTransit_CreateCsr_PreservesNonExtensionAttributes(t, b, s, "test-key")
		})
	}
}

// testTransit_CreateCsr_PreservesNonExtensionAttributes builds a template CSR
// with a challengePassword attribute and a critical custom extension, then:
//  1. Calls parseCsr directly: extensionRequest is gone, challengePassword unchanged.
//  2. Re-signs via Transit: challengePassword survives, Critical flag intact, signature valid.
func testTransit_CreateCsr_PreservesNonExtensionAttributes(t *testing.T, b logical.Backend, s logical.Storage, keyName string) {
	t.Helper()

	// PKCS#9 challengePassword (OID 1.2.840.113549.1.9.7) — a non-extension attribute.
	oidChallengePassword := asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 7}
	challengePasswordValue := "vault-challenge-secret"

	// A critical extension so that extensionRequest appears in Attributes after parsing.
	customOID := asn1.ObjectIdentifier{1, 2, 3, 4, 5}
	customExtValue, err := asn1.Marshal("vault-test-extension-value")
	require.NoError(t, err)

	// Set Attributes directly so challengePassword ends up in the encoded CSR.
	// ParseCertificateRequest will surface it back alongside the extensionRequest
	// attribute that CreateCertificateRequest adds for the extensions.
	templateCSR := &x509.CertificateRequest{
		Subject: pkix.Name{CommonName: "test.example.com"},
		Attributes: []pkix.AttributeTypeAndValueSET{
			{
				Type: oidChallengePassword,
				Value: [][]pkix.AttributeTypeAndValue{
					{{Type: oidChallengePassword, Value: challengePasswordValue}},
				},
			},
		},
		ExtraExtensions: []pkix.Extension{
			{
				Id:       customOID,
				Critical: true,
				Value:    customExtValue,
			},
		},
	}

	throwawayKey, err := cryptoutil.GenerateRSAKey(cryptoRand.Reader, 2048)
	require.NoError(t, err)

	templateDER, err := x509.CreateCertificateRequest(cryptoRand.Reader, templateCSR, throwawayKey)
	require.NoError(t, err)

	pemTemplate := string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: templateDER,
	}))

	// Confirm both attributes are present in the template before we run parseCsr.
	verifyBlock, _ := pem.Decode([]byte(pemTemplate))
	require.NotNil(t, verifyBlock)
	templateParsed, err := x509.ParseCertificateRequest(verifyBlock.Bytes)
	require.NoError(t, err)

	hasChallengeInTemplate := false
	hasExtReqInTemplate := false
	for _, attr := range templateParsed.Attributes {
		if attr.Type.Equal(oidChallengePassword) {
			hasChallengeInTemplate = true
		}
		if attr.Type.Equal(certutil.OidExtensionRequest) {
			hasExtReqInTemplate = true
		}
	}
	require.True(t, hasChallengeInTemplate, "challengePassword must be present in template before re-signing")
	require.True(t, hasExtReqInTemplate, "extensionRequest must be present in template (pre-condition for the filter to matter)")

	// Call parseCsr directly: extensionRequest must be gone, challengePassword untouched.
	parsedIntermediate, err := parseCsr(pemTemplate)
	require.NoError(t, err)

	hasChallengeAfterParse := false
	for _, attr := range parsedIntermediate.Attributes {
		require.False(t, attr.Type.Equal(certutil.OidExtensionRequest),
			"OidExtensionRequest was not removed from Attributes by parseCsr")
		if attr.Type.Equal(oidChallengePassword) {
			hasChallengeAfterParse = true
			require.Len(t, attr.Value, 1, "challengePassword attribute should have exactly one value set")
			require.Len(t, attr.Value[0], 1, "challengePassword value set should have exactly one entry")
			require.Equal(t, challengePasswordValue, attr.Value[0][0].Value,
				"challengePassword value changed by parseCsr")
		}
	}
	require.True(t, hasChallengeAfterParse, "challengePassword attribute missing from parseCsr output")

	// Re-sign via Transit and check the output CSR.
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("keys/%s/csr", keyName),
		Storage:   s,
		Data:      map[string]interface{}{"csr": pemTemplate},
	})
	require.NoError(t, err)
	require.False(t, resp != nil && resp.IsError(), "transit /csr failed: %v", resp)

	pemOut, ok := resp.Data["csr"].(string)
	require.True(t, ok, "response missing 'csr' field")

	block, _ := pem.Decode([]byte(pemOut))
	require.NotNil(t, block, "failed to PEM-decode output CSR")
	outCSR, err := x509.ParseCertificateRequest(block.Bytes)
	require.NoError(t, err, "output CSR is structurally invalid")

	// challengePassword must still be present with its original value.
	// (OidExtensionRequest reappears here because CreateCertificateRequest
	// always adds it when encoding extensions; filter correctness is
	// already covered by the parseCsr check above.)
	var foundChallengePassword bool
	for _, attr := range outCSR.Attributes {
		if !attr.Type.Equal(oidChallengePassword) {
			continue
		}
		foundChallengePassword = true
		require.Len(t, attr.Value, 1, "challengePassword attribute should have exactly one value set")
		require.Len(t, attr.Value[0], 1, "challengePassword value set should have exactly one entry")
		require.Equal(t, challengePasswordValue, attr.Value[0][0].Value,
			"challengePassword value changed during re-signing")
	}
	require.True(t, foundChallengePassword, "challengePassword attribute missing from re-signed CSR")

	// Custom critical extension: value and Critical flag must survive re-signing.
	var foundCustomExt bool
	for _, ext := range outCSR.Extensions {
		if !ext.Id.Equal(customOID) {
			continue
		}
		foundCustomExt = true
		require.True(t, ext.Critical,
			"Critical flag was dropped from custom extension during re-signing")
		require.Equal(t, customExtValue, ext.Value,
			"extension value bytes changed during re-signing")
	}
	require.True(t, foundCustomExt, "custom extension %v not found in re-signed CSR", customOID)

	// CSR signature must be valid.
	require.NoError(t, outCSR.CheckSignature(),
		"re-signed CSR has an invalid signature")
}

// TestTransit_Certs_CreateCsr_PreservesExtensions checks that re-signing a
// template CSR preserves extension Critical flags and does not duplicate SANs.
func TestTransit_Certs_CreateCsr_PreservesExtensions(t *testing.T) {
	for _, keyType := range []string{"rsa-2048", "ecdsa-p256", "ed25519"} {
		t.Run(keyType, func(t *testing.T) {
			b, s := createBackendWithStorage(t)

			resp, err := b.HandleRequest(context.Background(), &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "keys/test-key",
				Storage:   s,
				Data:      map[string]interface{}{"type": keyType},
			})
			require.NoError(t, err)
			require.False(t, resp != nil && resp.IsError(), "key creation failed: %v", resp)

			testTransit_CreateCsr_PreservesExtensions(t, b, s, "test-key")
		})
	}
}

func testTransit_CreateCsr_PreservesExtensions(t *testing.T, b logical.Backend, s logical.Storage, keyName string) {
	t.Helper()

	// Template has SANs and a critical extension to exercise both failure paths.
	customOID := asn1.ObjectIdentifier{1, 2, 3, 4, 5}
	customExtValue, err := asn1.Marshal("vault-test-extension-value")
	require.NoError(t, err)

	templateCSR := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   "test.example.com",
			Organization: []string{"Vault Test Org"},
		},
		DNSNames:    []string{"test.example.com", "alt.example.com"},
		IPAddresses: []net.IP{net.ParseIP("192.168.1.1")},
		ExtraExtensions: []pkix.Extension{
			{
				Id:       customOID,
				Critical: true,
				Value:    customExtValue,
			},
		},
	}

	// Throwaway key just to produce a valid PEM; Transit replaces it with its own.
	throwawayKey, err := cryptoutil.GenerateRSAKey(cryptoRand.Reader, 2048)
	require.NoError(t, err)

	templateDER, err := x509.CreateCertificateRequest(cryptoRand.Reader, templateCSR, throwawayKey)
	require.NoError(t, err)

	pemTemplate := string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: templateDER,
	}))

	// Re-sign via Transit.
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("keys/%s/csr", keyName),
		Storage:   s,
		Data:      map[string]interface{}{"csr": pemTemplate},
	})
	require.NoError(t, err)
	require.False(t, resp != nil && resp.IsError(), "transit /csr failed: %v", resp)

	pemOut, ok := resp.Data["csr"].(string)
	require.True(t, ok, "response missing 'csr' field")

	// A duplicate SAN block would cause "asn1: syntax error: sequence truncated" here.
	block, _ := pem.Decode([]byte(pemOut))
	require.NotNil(t, block, "failed to PEM-decode output CSR")
	outCSR, err := x509.ParseCertificateRequest(block.Bytes)
	require.NoError(t, err, "output CSR is structurally invalid (possible duplicate SAN / malformed ASN.1)")

	// Subject DN preserved.
	require.Equal(t, "test.example.com", outCSR.Subject.CommonName,
		"Subject CommonName not preserved")
	require.Equal(t, []string{"Vault Test Org"}, outCSR.Subject.Organization,
		"Subject Organization not preserved")

	// SANs preserved.
	require.ElementsMatch(t, []string{"test.example.com", "alt.example.com"}, outCSR.DNSNames,
		"DNS SANs not preserved")
	require.Len(t, outCSR.IPAddresses, 1, "expected exactly one IP SAN")
	require.True(t, outCSR.IPAddresses[0].Equal(net.ParseIP("192.168.1.1")),
		"IP SAN not preserved")

	// Count raw SAN entries — Go accepts duplicates silently, but openssl rejects them.
	sanCount := 0
	for _, ext := range outCSR.Extensions {
		if ext.Id.Equal(certutil.OidExtensionSubjectAltName) {
			sanCount++
		}
	}
	require.Equal(t, 1, sanCount,
		"expected exactly one SAN extension in output CSR, got %d", sanCount)

	// Custom extension: Critical flag and value bytes must survive re-signing.
	var foundCustomExt bool
	for _, ext := range outCSR.Extensions {
		if ext.Id.Equal(customOID) {
			foundCustomExt = true
			require.True(t, ext.Critical,
				"Critical flag was dropped during re-signing")
			require.Equal(t, customExtValue, ext.Value,
				"extension value bytes not preserved during re-signing")
			break
		}
	}
	require.True(t, foundCustomExt,
		"custom extension OID %v not found in output CSR", customOID)

	// CSR signature must be valid.
	require.NoError(t, outCSR.CheckSignature(),
		"re-signed CSR has an invalid signature")
}

// TestTransit_Certs_CreateCsr_PreservesOtherNameSAN verifies that otherName/UPN
// entries inside the SAN extension survive re-signing unchanged.
func TestTransit_Certs_CreateCsr_PreservesOtherNameSAN(t *testing.T) {
	for _, keyType := range []string{"rsa-2048", "ecdsa-p256", "ed25519"} {
		t.Run(keyType, func(t *testing.T) {
			b, s := createBackendWithStorage(t)

			resp, err := b.HandleRequest(context.Background(), &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "keys/test-key",
				Storage:   s,
				Data:      map[string]interface{}{"type": keyType},
			})
			require.NoError(t, err)
			require.False(t, resp != nil && resp.IsError(), "key creation failed: %v", resp)

			testTransit_CreateCsr_PreservesOtherNameSAN(t, b, s, "test-key")
		})
	}
}

func testTransit_CreateCsr_PreservesOtherNameSAN(t *testing.T, b logical.Backend, s logical.Storage, keyName string) {
	t.Helper()

	// Hand-craft a SAN extension with a dNSName and an otherName/UPN (tag 0).
	// Go's stdlib has no API for otherName so we build the raw ASN.1 directly.
	// OtherName ::= SEQUENCE { type-id OID, value [0] EXPLICIT ANY }
	upnOID := asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 311, 20, 2, 3} // UPN OID
	upnValue := "user@corp.example.com"

	// UTF8String inside [0] EXPLICIT.
	upnUTF8, err := asn1.Marshal(asn1.RawValue{
		Tag:   asn1.TagUTF8String,
		Bytes: []byte(upnValue),
	})
	require.NoError(t, err)

	// Full OtherName SEQUENCE.
	otherNameSeq, err := asn1.Marshal(struct {
		OID   asn1.ObjectIdentifier
		Value asn1.RawValue `asn1:"tag:0,explicit"`
	}{
		OID: upnOID,
		Value: asn1.RawValue{
			FullBytes: upnUTF8,
		},
	})
	require.NoError(t, err)

	// Peel the outer SEQUENCE wrapper to get the raw content bytes.
	var otherNameParsed asn1.RawValue
	_, err = asn1.Unmarshal(otherNameSeq, &otherNameParsed)
	require.NoError(t, err)

	// GeneralName CHOICE, context tag 0 (otherName).
	otherNameRaw := asn1.RawValue{
		Class:      asn1.ClassContextSpecific,
		Tag:        0,
		IsCompound: true,
		Bytes:      otherNameParsed.Bytes, // content only, no outer tag+length
	}

	// dNSName GeneralName, context tag 2.
	dnsRaw := asn1.RawValue{
		Class: asn1.ClassContextSpecific,
		Tag:   2,
		Bytes: []byte("dns.corp.example.com"),
	}

	// SubjectAltName SEQUENCE wrapping both entries.
	sanSeq, err := asn1.Marshal(asn1.RawValue{
		Tag:        asn1.TagSequence,
		IsCompound: true,
		Bytes: func() []byte {
			dns, _ := asn1.Marshal(dnsRaw)
			other, _ := asn1.Marshal(otherNameRaw)
			return append(dns, other...)
		}(),
	})
	require.NoError(t, err)

	sanExt := pkix.Extension{
		Id:    certutil.OidExtensionSubjectAltName,
		Value: sanSeq,
	}

	// Use ExtraExtensions so the raw SAN blob is passed verbatim.
	templateCSR := &x509.CertificateRequest{
		Subject:         pkix.Name{CommonName: "test.example.com"},
		ExtraExtensions: []pkix.Extension{sanExt},
	}

	throwawayKey, err := cryptoutil.GenerateRSAKey(cryptoRand.Reader, 2048)
	require.NoError(t, err)

	templateDER, err := x509.CreateCertificateRequest(cryptoRand.Reader, templateCSR, throwawayKey)
	require.NoError(t, err)

	pemTemplate := string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: templateDER,
	}))

	// Re-sign via Transit.
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("keys/%s/csr", keyName),
		Storage:   s,
		Data:      map[string]interface{}{"csr": pemTemplate},
	})
	require.NoError(t, err)
	require.False(t, resp != nil && resp.IsError(), "transit /csr failed: %v", resp)

	pemOut, ok := resp.Data["csr"].(string)
	require.True(t, ok, "response missing 'csr' field")

	outBlock, _ := pem.Decode([]byte(pemOut))
	require.NotNil(t, outBlock, "failed to PEM-decode output CSR")
	outCSR, err := x509.ParseCertificateRequest(outBlock.Bytes)
	require.NoError(t, err, "output CSR is structurally invalid")

	// Exactly one SAN in the output — no duplicates.
	var outSANExt *pkix.Extension
	for i := range outCSR.Extensions {
		if outCSR.Extensions[i].Id.Equal(certutil.OidExtensionSubjectAltName) {
			outSANExt = &outCSR.Extensions[i]
		}
	}
	sanCount := 0
	for _, ext := range outCSR.Extensions {
		if ext.Id.Equal(certutil.OidExtensionSubjectAltName) {
			sanCount++
		}
	}
	require.Equal(t, 1, sanCount, "expected exactly one SAN extension in output CSR, got %d", sanCount)
	require.NotNil(t, outSANExt, "SAN extension missing from output CSR")

	// Raw SAN bytes must be identical — proves otherName was not rebuilt from struct fields.
	require.Equal(t, sanExt.Value, outSANExt.Value,
		"SAN extension bytes changed during re-signing: otherName/UPN may have been lost")

	require.Contains(t, outCSR.DNSNames, "dns.corp.example.com",
		"dNSName SAN not preserved")

	require.NoError(t, outCSR.CheckSignature(),
		"re-signed CSR has an invalid signature")
}
