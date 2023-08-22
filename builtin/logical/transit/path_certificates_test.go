// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transit

import (
	"context"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/pki"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func TestTransit_Certs_CreateCsr(t *testing.T) {
	// NOTE: Use an existing CSR or generate one here?
	templateCsr := `
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

	testTransit_CreateCsr(t, "rsa-2048", templateCsr)
	testTransit_CreateCsr(t, "rsa-3072", templateCsr)
	testTransit_CreateCsr(t, "rsa-4096", templateCsr)
	testTransit_CreateCsr(t, "ecdsa-p256", templateCsr)
	testTransit_CreateCsr(t, "ecdsa-p384", templateCsr)
	testTransit_CreateCsr(t, "ecdsa-p521", templateCsr)
	testTransit_CreateCsr(t, "ed25519", templateCsr)
	testTransit_CreateCsr(t, "aes256-gcm96", templateCsr)
}

func testTransit_CreateCsr(t *testing.T, keyType, pemTemplateCsr string) {
	var resp *logical.Response
	var err error
	b, s := createBackendWithStorage(t)

	// Create the policy
	policyReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "keys/test-key",
		Storage:   s,
		Data: map[string]interface{}{
			"type": keyType,
		},
	}
	resp, err = b.HandleRequest(context.Background(), policyReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("resp: %#v\nerr: %v", resp, err)
	}

	csrSignReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "keys/test-key/csr",
		Storage:   s,
		Data: map[string]interface{}{
			"csr": pemTemplateCsr,
		},
	}

	resp, err = b.HandleRequest(context.Background(), csrSignReq)

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

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)
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

	testTransit_ImportCertChain(t, client, "rsa-2048")
	testTransit_ImportCertChain(t, client, "rsa-3072")
	testTransit_ImportCertChain(t, client, "rsa-4096")
	testTransit_ImportCertChain(t, client, "ecdsa-p256")
	testTransit_ImportCertChain(t, client, "ecdsa-p384")
	testTransit_ImportCertChain(t, client, "ecdsa-p521")
	testTransit_ImportCertChain(t, client, "ed25519")
}

func testTransit_ImportCertChain(t *testing.T, apiClient *api.Client, keyType string) {
	keyName := fmt.Sprintf("%s", keyType)
	issuerName := fmt.Sprintf("%s-issuer", keyType)

	// Create transit key
	_, err := apiClient.Logical().Write(fmt.Sprintf("transit/keys/%s", keyName), map[string]interface{}{
		"type": keyType,
	})
	require.NoError(t, err)

	// Setup a new CSR
	privKey, err := rsa.GenerateKey(cryptoRand.Reader, 3072)
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

	// Create CSR from template CSR fields and key in transit
	resp, err := apiClient.Logical().Write(fmt.Sprintf("transit/keys/%s/csr", keyName), map[string]interface{}{
		"csr": string(pemTemplateCsr),
	})
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
	// Import certificate chain to transit key version
	resp, err = apiClient.Logical().Write(fmt.Sprintf("transit/keys/%s/set-certificate", keyName), map[string]interface{}{
		"certificate_chain": certificateChain,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	resp, err = apiClient.Logical().Read(fmt.Sprintf("transit/keys/%s", keyName))
	require.NoError(t, err)
	require.NotNil(t, resp)
	keys, ok := resp.Data["keys"].(map[string]interface{})
	if !ok {
		t.Fatalf("could not cast Keys value")
	}
	keyData, ok := keys["1"].(map[string]interface{})
	if !ok {
		t.Fatalf("could not cast key version 1 from keys")
	}
	_, present := keyData["certificate_chain"]
	if !present {
		t.Fatalf("certificate chain not present in key version 1")
	}
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

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)
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

	testTransit_ImportInvalidCertChain(t, client, "rsa-2048")
	testTransit_ImportInvalidCertChain(t, client, "rsa-3072")
	testTransit_ImportInvalidCertChain(t, client, "rsa-4096")
	testTransit_ImportInvalidCertChain(t, client, "ecdsa-p256")
	testTransit_ImportInvalidCertChain(t, client, "ecdsa-p384")
	testTransit_ImportInvalidCertChain(t, client, "ecdsa-p521")
	testTransit_ImportInvalidCertChain(t, client, "ed25519")
}

func testTransit_ImportInvalidCertChain(t *testing.T, apiClient *api.Client, keyType string) {
	keyName := fmt.Sprintf("%s", keyType)
	issuerName := fmt.Sprintf("%s-issuer", keyType)

	// Create transit key
	_, err := apiClient.Logical().Write(fmt.Sprintf("transit/keys/%s", keyName), map[string]interface{}{
		"type": keyType,
	})
	require.NoError(t, err)

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
	} else if strings.HasPrefix(keyType, "ecdas") {
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
	})
	require.Error(t, err)
}
