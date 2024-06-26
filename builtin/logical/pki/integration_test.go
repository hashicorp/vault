// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	vaulthttp "github.com/hashicorp/vault/http"
	vaultocsp "github.com/hashicorp/vault/sdk/helper/ocsp"
	"github.com/hashicorp/vault/sdk/helper/testhelpers/schema"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func TestIntegration_RotateRootUsesNext(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/rotate/internal",
		Storage:   s,
		Data: map[string]interface{}{
			"common_name": "test.com",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed rotate root")
	require.NotNil(t, resp, "got nil response from rotate root")
	require.False(t, resp.IsError(), "got an error from rotate root: %#v", resp)

	issuerId1 := resp.Data["issuer_id"].(issuing.IssuerID)
	issuerName1 := resp.Data["issuer_name"]

	require.NotEmpty(t, issuerId1, "issuer id was empty on initial rotate root command")
	require.Equal(t, "next", issuerName1, "expected an issuer name of next on initial rotate root command")

	// Call it again, we should get a new issuer id, but since next issuer_name is used we should get a blank value.
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/rotate/internal",
		Storage:   s,
		Data: map[string]interface{}{
			"common_name": "test.com",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed rotate root")
	require.NotNil(t, resp, "got nil response from rotate root")
	require.False(t, resp.IsError(), "got an error from rotate root: %#v", resp)

	issuerId2 := resp.Data["issuer_id"].(issuing.IssuerID)
	issuerName2 := resp.Data["issuer_name"]

	require.NotEmpty(t, issuerId2, "issuer id was empty on second rotate root command")
	require.NotEqual(t, issuerId1, issuerId2, "should have been different issuer ids")
	require.Empty(t, issuerName2, "expected a blank issuer name on the second rotate root command")

	// Call it again, making sure we can use our own name if desired.
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/rotate/internal",
		Storage:   s,
		Data: map[string]interface{}{
			"common_name": "test.com",
			"issuer_name": "next-cert",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed rotate root")
	require.NotNil(t, resp, "got nil response from rotate root")
	require.False(t, resp.IsError(), "got an error from rotate root: %#v", resp)

	issuerId3 := resp.Data["issuer_id"].(issuing.IssuerID)
	issuerName3 := resp.Data["issuer_name"]

	require.NotEmpty(t, issuerId3, "issuer id was empty on third rotate root command")
	require.NotEqual(t, issuerId3, issuerId1, "should have been different issuer id from initial")
	require.NotEqual(t, issuerId3, issuerId2, "should have been different issuer id from second call")
	require.Equal(t, "next-cert", issuerName3, "expected an issuer name that we specified on third rotate root command")
}

func TestIntegration_ReplaceRootNormal(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// generate roots
	genTestRootCa(t, b, s)
	issuerId2, _ := genTestRootCa(t, b, s)

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/replace",
		Storage:   s,
		Data: map[string]interface{}{
			"default": issuerId2.String(),
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed replacing root")
	require.NotNil(t, resp, "got nil response from replacing root")
	require.False(t, resp.IsError(), "got an error from replacing root: %#v", resp)

	replacedIssuer := resp.Data["default"]
	require.Equal(t, issuerId2, replacedIssuer, "expected return value to match issuer we set")

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.ReadOperation,
		Path:       "config/issuers",
		Storage:    s,
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed replacing root")
	require.NotNil(t, resp, "got nil response from replacing root")
	require.False(t, resp.IsError(), "got an error from replacing root: %#v", resp)

	defaultIssuer := resp.Data["default"]
	require.Equal(t, issuerId2, defaultIssuer, "expected default issuer to be updated")
}

func TestIntegration_ReplaceRootDefaultsToNext(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// generate roots
	genTestRootCa(t, b, s)
	issuerId2, _ := genTestRootCaWithIssuerName(t, b, s, "next")

	// Do not specify the default value to replace.
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "root/replace",
		Storage:    s,
		Data:       map[string]interface{}{},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed replacing root")
	require.NotNil(t, resp, "got nil response from replacing root")
	require.False(t, resp.IsError(), "got an error from replacing root: %#v", resp)

	replacedIssuer := resp.Data["default"]
	require.Equal(t, issuerId2, replacedIssuer, "expected return value to match the 'next' issuer we set")

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.ReadOperation,
		Path:       "config/issuers",
		Storage:    s,
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed replacing root")
	require.NotNil(t, resp, "got nil response from replacing root")
	require.False(t, resp.IsError(), "got an error from replacing root: %#v", resp)

	defaultIssuer := resp.Data["default"]
	require.Equal(t, issuerId2, defaultIssuer, "expected default issuer to be updated")
}

func TestIntegration_ReplaceRootBadIssuer(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// generate roots
	genTestRootCa(t, b, s)
	genTestRootCa(t, b, s)

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/replace",
		Storage:   s,
		Data: map[string]interface{}{
			"default": "a-bad-issuer-id",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed replacing root, should have been an error in the response.")
	require.NotNil(t, resp, "got nil response from replacing root")
	require.True(t, resp.IsError(), "did not get an error from replacing root: %#v", resp)

	// Make sure we trap replacing with default.
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/replace",
		Storage:   s,
		Data: map[string]interface{}{
			"default": "default",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed replacing root, should have been an error in the response.")
	require.NotNil(t, resp, "got nil response from replacing root")
	require.True(t, resp.IsError(), "did not get an error from replacing root: %#v", resp)

	// Make sure we trap replacing with blank string.
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/replace",
		Storage:   s,
		Data: map[string]interface{}{
			"default": "",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed replacing root, should have been an error in the response.")
	require.NotNil(t, resp, "got nil response from replacing root")
	require.True(t, resp.IsError(), "did not get an error from replacing root: %#v", resp)
}

func TestIntegration_SetSignedWithBackwardsPemBundles(t *testing.T) {
	t.Parallel()
	rootBackend, rootStorage := CreateBackendWithStorage(t)
	intBackend, intStorage := CreateBackendWithStorage(t)

	// generate root
	resp, err := rootBackend.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "issuers/generate/root/internal",
		Storage:   rootStorage,
		Data: map[string]interface{}{
			"common_name": "test.com",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed generating root ca")
	require.NotNil(t, resp, "got nil response from generating root ca")
	require.False(t, resp.IsError(), "got an error from generating root ca: %#v", resp)
	rootCert := resp.Data["certificate"].(string)

	schema.ValidateResponse(t, schema.GetResponseSchema(t, rootBackend.Route("issuers/generate/root/internal"), logical.UpdateOperation), resp, true)

	// generate intermediate
	resp, err = intBackend.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "issuers/generate/intermediate/internal",
		Storage:   intStorage,
		Data: map[string]interface{}{
			"common_name": "test.com",
		},
		MountPoint: "pki-int/",
	})
	require.NoError(t, err, "failed generating int ca")
	require.NotNil(t, resp, "got nil response from generating int ca")
	require.False(t, resp.IsError(), "got an error from generating int ca: %#v", resp)
	intCsr := resp.Data["csr"].(string)

	// sign csr
	resp, err = rootBackend.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/sign-intermediate",
		Storage:   rootStorage,
		Data: map[string]interface{}{
			"csr":    intCsr,
			"format": "pem_bundle",
		},
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed generating root ca")
	require.NotNil(t, resp, "got nil response from generating root ca")
	require.False(t, resp.IsError(), "got an error from generating root ca: %#v", resp)

	intCert := resp.Data["certificate"].(string)

	// Send in the chain backwards now and make sure we link intCert as default.
	resp, err = intBackend.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "intermediate/set-signed",
		Storage:   intStorage,
		Data: map[string]interface{}{
			"certificate": rootCert + "\n" + intCert + "\n",
		},
		MountPoint: "pki-int/",
	})
	require.NoError(t, err, "failed generating root ca")
	require.NotNil(t, resp, "got nil response from generating root ca")
	require.False(t, resp.IsError(), "got an error from generating root ca: %#v", resp)

	// setup role
	resp, err = intBackend.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/example",
		Storage:   intStorage,
		Data: map[string]interface{}{
			"allowed_domains":  "example.com",
			"allow_subdomains": "true",
			"max_ttl":          "1h",
		},
		MountPoint: "pki-int/",
	})
	require.NoError(t, err, "failed setting up role example")
	require.NotNil(t, resp, "got nil response from setting up role example: %#v", resp)

	schema.ValidateResponse(t, schema.GetResponseSchema(t, intBackend.Route("roles/example"), logical.UpdateOperation), resp, true)

	// Issue cert
	resp, err = intBackend.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "issue/example",
		Storage:   intStorage,
		Data: map[string]interface{}{
			"common_name": "test.example.com",
			"ttl":         "5m",
		},
		MountPoint: "pki-int/",
	})
	require.NoError(t, err, "failed issuing a leaf cert from int ca")
	require.NotNil(t, resp, "got nil response issuing a leaf cert from int ca")
	require.False(t, resp.IsError(), "got an error issuing a leaf cert from int ca: %#v", resp)

	schema.ValidateResponse(t, schema.GetResponseSchema(t, intBackend.Route("issue/example"), logical.UpdateOperation), resp, true)
}

func TestIntegration_CSRGeneration(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)
	testCases := []struct {
		keyType               string
		usePss                bool
		keyBits               int
		sigBits               int
		expectedPublicKeyType crypto.PublicKey
		expectedSignature     x509.SignatureAlgorithm
	}{
		{"rsa", false, 2048, 0, &rsa.PublicKey{}, x509.SHA256WithRSA},
		{"rsa", false, 2048, 384, &rsa.PublicKey{}, x509.SHA384WithRSA},
		// Add back once https://github.com/golang/go/issues/45990 is fixed.
		// {"rsa", true, 2048, 0, &rsa.PublicKey{}, x509.SHA256WithRSAPSS},
		// {"rsa", true, 2048, 512, &rsa.PublicKey{}, x509.SHA512WithRSAPSS},
		{"ec", false, 224, 0, &ecdsa.PublicKey{}, x509.ECDSAWithSHA256},
		{"ec", false, 256, 0, &ecdsa.PublicKey{}, x509.ECDSAWithSHA256},
		{"ec", false, 384, 0, &ecdsa.PublicKey{}, x509.ECDSAWithSHA384},
		{"ec", false, 521, 0, &ecdsa.PublicKey{}, x509.ECDSAWithSHA512},
		{"ec", false, 521, 224, &ecdsa.PublicKey{}, x509.ECDSAWithSHA512}, // We ignore signature_bits for ec
		{"ed25519", false, 0, 0, ed25519.PublicKey{}, x509.PureEd25519},   // We ignore both fields for ed25519
	}
	for _, tc := range testCases {
		keyTypeName := tc.keyType
		if tc.usePss {
			keyTypeName = tc.keyType + "-pss"
		}
		testName := fmt.Sprintf("%s-%d-%d", keyTypeName, tc.keyBits, tc.sigBits)
		t.Run(testName, func(t *testing.T) {
			resp, err := CBWrite(b, s, "intermediate/generate/internal", map[string]interface{}{
				"common_name":    "myint.com",
				"key_type":       tc.keyType,
				"key_bits":       tc.keyBits,
				"signature_bits": tc.sigBits,
				"use_pss":        tc.usePss,
			})
			requireSuccessNonNilResponse(t, resp, err)
			requireFieldsSetInResp(t, resp, "csr")

			csrString := resp.Data["csr"].(string)
			pemBlock, _ := pem.Decode([]byte(csrString))
			require.NotNil(t, pemBlock, "failed to parse returned csr pem block")
			csr, err := x509.ParseCertificateRequest(pemBlock.Bytes)
			require.NoError(t, err, "failed parsing certificate request")

			require.Equal(t, tc.expectedSignature, csr.SignatureAlgorithm,
				"Expected %s, got %s", tc.expectedSignature.String(), csr.SignatureAlgorithm.String())
			require.IsType(t, tc.expectedPublicKeyType, csr.PublicKey)
		})
	}
}

func TestIntegration_AutoIssuer(t *testing.T) {
	t.Parallel()
	b, s := CreateBackendWithStorage(t)

	// Generate two roots. The first should become default under the existing
	// behavior; when we update the config and generate a second, it should
	// take over as default. Deleting the first and re-importing it will make
	// it default again, and then disabling the option and removing and
	// reimporting the second and creating a new root won't affect it again.
	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "Root X1",
		"issuer_name": "root-1",
		"key_type":    "ec",
	})

	requireSuccessNonNilResponse(t, resp, err)
	issuerIdOne := resp.Data["issuer_id"]
	require.NotEmpty(t, issuerIdOne)
	certOne := resp.Data["certificate"]
	require.NotEmpty(t, certOne)

	resp, err = CBRead(b, s, "config/issuers")
	requireSuccessNonNilResponse(t, resp, err)
	require.Equal(t, issuerIdOne, resp.Data["default"])

	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("config/issuers"), logical.ReadOperation), resp, true)

	// Enable the new config option.
	resp, err = CBWrite(b, s, "config/issuers", map[string]interface{}{
		"default":                       issuerIdOne,
		"default_follows_latest_issuer": true,
	})
	require.NoError(t, err)
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route("config/issuers"), logical.UpdateOperation), resp, true)

	// Now generate the second root; it should become default.
	resp, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "Root X2",
		"issuer_name": "root-2",
		"key_type":    "ec",
	})
	requireSuccessNonNilResponse(t, resp, err)
	issuerIdTwo := resp.Data["issuer_id"]
	require.NotEmpty(t, issuerIdTwo)
	certTwo := resp.Data["certificate"]
	require.NotEmpty(t, certTwo)

	resp, err = CBRead(b, s, "config/issuers")
	requireSuccessNonNilResponse(t, resp, err)
	require.Equal(t, issuerIdTwo, resp.Data["default"])

	// Deleting the first shouldn't affect the default issuer.
	_, err = CBDelete(b, s, "issuer/root-1")
	require.NoError(t, err)
	resp, err = CBRead(b, s, "config/issuers")
	requireSuccessNonNilResponse(t, resp, err)
	require.Equal(t, issuerIdTwo, resp.Data["default"])

	// But reimporting it should update it to the new issuer's value.
	resp, err = CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
		"pem_bundle": certOne,
	})
	requireSuccessNonNilResponse(t, resp, err)
	issuerIdOneReimported := issuing.IssuerID(resp.Data["imported_issuers"].([]string)[0])

	resp, err = CBRead(b, s, "config/issuers")
	requireSuccessNonNilResponse(t, resp, err)
	require.Equal(t, issuerIdOneReimported, resp.Data["default"])

	// Now update the config to disable this option again.
	_, err = CBWrite(b, s, "config/issuers", map[string]interface{}{
		"default":                       issuerIdOneReimported,
		"default_follows_latest_issuer": false,
	})
	require.NoError(t, err)

	// Generating a new root shouldn't update the default.
	resp, err = CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "Root X3",
		"issuer_name": "root-3",
		"key_type":    "ec",
	})
	requireSuccessNonNilResponse(t, resp, err)
	issuerIdThree := resp.Data["issuer_id"]
	require.NotEmpty(t, issuerIdThree)

	resp, err = CBRead(b, s, "config/issuers")
	requireSuccessNonNilResponse(t, resp, err)
	require.Equal(t, issuerIdOneReimported, resp.Data["default"])

	// Deleting and re-importing root 2 should also not affect it.
	_, err = CBDelete(b, s, "issuer/root-2")
	require.NoError(t, err)
	resp, err = CBRead(b, s, "config/issuers")
	requireSuccessNonNilResponse(t, resp, err)
	require.Equal(t, issuerIdOneReimported, resp.Data["default"])

	resp, err = CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
		"pem_bundle": certTwo,
	})
	requireSuccessNonNilResponse(t, resp, err)
	require.Equal(t, 1, len(resp.Data["imported_issuers"].([]string)))
	resp, err = CBRead(b, s, "config/issuers")
	requireSuccessNonNilResponse(t, resp, err)
	require.Equal(t, issuerIdOneReimported, resp.Data["default"])
}

// TestLDAPAiaCrlUrls validates we can properly handle CRL urls that are ldap based.
func TestLDAPAiaCrlUrls(t *testing.T) {
	t.Parallel()

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		NumCores:    1,
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	singleCore := cluster.Cores[0]
	vault.TestWaitActive(t, singleCore.Core)
	client := singleCore.Client

	mountPKIEndpoint(t, client, "pki")

	// Attempt multiple urls
	crls := []string{
		"ldap://ldap.example.com/cn=example%20CA,dc=example,dc=com?certificateRevocationList;binary",
		"ldap://ldap.example.com/cn=CA,dc=example,dc=com?authorityRevocationList;binary",
	}

	_, err := client.Logical().Write("pki/config/urls", map[string]interface{}{
		"crl_distribution_points": crls,
	})
	require.NoError(t, err)

	resp, err := client.Logical().Read("pki/config/urls")
	require.NoError(t, err, "failed reading config/urls")
	require.NotNil(t, resp, "resp was nil")
	require.NotNil(t, resp.Data, "data within response was nil")
	require.NotEmpty(t, resp.Data["crl_distribution_points"], "crl_distribution_points was nil within data")
	require.Len(t, resp.Data["crl_distribution_points"], len(crls))

	for _, crlVal := range crls {
		require.Contains(t, resp.Data["crl_distribution_points"], crlVal)
	}

	resp, err = client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "Root R1",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["issuer_id"])
	rootIssuerId := resp.Data["issuer_id"].(string)

	_, err = client.Logical().Write("pki/roles/example-root", map[string]interface{}{
		"allowed_domains":  "example.com",
		"allow_subdomains": "true",
		"max_ttl":          "1h",
		"key_type":         "ec",
		"issuer_ref":       rootIssuerId,
	})
	require.NoError(t, err)

	resp, err = client.Logical().Write("pki/issue/example-root", map[string]interface{}{
		"common_name": "test.example.com",
		"ttl":         "5m",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["certificate"])

	certPEM := resp.Data["certificate"].(string)
	certBlock, _ := pem.Decode([]byte(certPEM))
	require.NotNil(t, certBlock)
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	require.NoError(t, err)

	require.EqualValues(t, crls, cert.CRLDistributionPoints)
}

func TestIntegrationOCSPClientWithPKI(t *testing.T) {
	t.Parallel()

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": Factory,
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

	err := client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "32h",
		},
	})
	require.NoError(t, err)

	resp, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "Root R1",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["issuer_id"])
	rootIssuerId := resp.Data["issuer_id"].(string)

	// Set URLs pointing to the issuer.
	_, err = client.Logical().Write("pki/config/cluster", map[string]interface{}{
		"path":     client.Address() + "/v1/pki",
		"aia_path": client.Address() + "/v1/pki",
	})
	require.NoError(t, err)

	_, err = client.Logical().Write("pki/config/urls", map[string]interface{}{
		"enable_templating":       true,
		"crl_distribution_points": "{{cluster_aia_path}}/issuer/{{issuer_id}}/crl/der",
		"issuing_certificates":    "{{cluster_aia_path}}/issuer/{{issuer_id}}/der",
		"ocsp_servers":            "{{cluster_aia_path}}/ocsp",
	})
	require.NoError(t, err)

	// Build an intermediate CA
	resp, err = client.Logical().Write("pki/intermediate/generate/internal", map[string]interface{}{
		"common_name": "Int X1",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["csr"])
	intermediateCSR := resp.Data["csr"].(string)

	resp, err = client.Logical().Write("pki/root/sign-intermediate", map[string]interface{}{
		"csr": intermediateCSR,
		"ttl": "20h",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["certificate"])
	intermediateCert := resp.Data["certificate"]

	resp, err = client.Logical().Write("pki/intermediate/set-signed", map[string]interface{}{
		"certificate": intermediateCert,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["imported_issuers"])
	rawImportedIssuers := resp.Data["imported_issuers"].([]interface{})
	require.Equal(t, len(rawImportedIssuers), 1)
	importedIssuer := rawImportedIssuers[0].(string)
	require.NotEmpty(t, importedIssuer)

	// Set intermediate as default.
	_, err = client.Logical().Write("pki/config/issuers", map[string]interface{}{
		"default": importedIssuer,
	})
	require.NoError(t, err)

	// Setup roles for root, intermediate.
	_, err = client.Logical().Write("pki/roles/example-root", map[string]interface{}{
		"allowed_domains":  "example.com",
		"allow_subdomains": "true",
		"max_ttl":          "1h",
		"key_type":         "ec",
		"issuer_ref":       rootIssuerId,
	})
	require.NoError(t, err)

	_, err = client.Logical().Write("pki/roles/example-int", map[string]interface{}{
		"allowed_domains":  "example.com",
		"allow_subdomains": "true",
		"max_ttl":          "1h",
		"key_type":         "ec",
	})
	require.NoError(t, err)

	// Issue certs and validate them against OCSP.
	for _, path := range []string{"pki/issue/example-int", "pki/issue/example-root"} {
		t.Logf("Validating against path: %v", path)
		resp, err = client.Logical().Write(path, map[string]interface{}{
			"common_name": "test.example.com",
			"ttl":         "5m",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Data)
		require.NotEmpty(t, resp.Data["certificate"])
		require.NotEmpty(t, resp.Data["issuing_ca"])
		require.NotEmpty(t, resp.Data["serial_number"])

		certPEM := resp.Data["certificate"].(string)
		certBlock, _ := pem.Decode([]byte(certPEM))
		require.NotNil(t, certBlock)
		cert, err := x509.ParseCertificate(certBlock.Bytes)
		require.NoError(t, err)
		require.NotNil(t, cert)

		issuerPEM := resp.Data["issuing_ca"].(string)
		issuerBlock, _ := pem.Decode([]byte(issuerPEM))
		require.NotNil(t, issuerBlock)
		issuer, err := x509.ParseCertificate(issuerBlock.Bytes)
		require.NoError(t, err)
		require.NotNil(t, issuer)

		serialNumber := resp.Data["serial_number"].(string)

		testLogger := hclog.New(hclog.DefaultOptions)

		conf := &vaultocsp.VerifyConfig{
			OcspFailureMode: vaultocsp.FailOpenFalse,
			ExtraCas:        []*x509.Certificate{cluster.CACert},
		}
		ocspClient := vaultocsp.New(func() hclog.Logger {
			return testLogger
		}, 10)

		_, err = client.Logical().Write("pki/revoke", map[string]interface{}{
			"serial_number": serialNumber,
		})
		require.NoError(t, err)

		err = ocspClient.VerifyLeafCertificate(context.Background(), cert, issuer, conf)
		require.Error(t, err)
	}
}

func genTestRootCa(t *testing.T, b *backend, s logical.Storage) (issuing.IssuerID, issuing.KeyID) {
	return genTestRootCaWithIssuerName(t, b, s, "")
}

func genTestRootCaWithIssuerName(t *testing.T, b *backend, s logical.Storage, issuerName string) (issuing.IssuerID, issuing.KeyID) {
	data := map[string]interface{}{
		"common_name": "test.com",
	}
	if len(issuerName) > 0 {
		data["issuer_name"] = issuerName
	}
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "issuers/generate/root/internal",
		Storage:    s,
		Data:       data,
		MountPoint: "pki/",
	})
	require.NoError(t, err, "failed generating root ca")
	require.NotNil(t, resp, "got nil response from generating root ca")
	require.False(t, resp.IsError(), "got an error from generating root ca: %#v", resp)

	issuerId := resp.Data["issuer_id"].(issuing.IssuerID)
	keyId := resp.Data["key_id"].(issuing.KeyID)

	require.NotEmpty(t, issuerId, "returned issuer id was empty")
	require.NotEmpty(t, keyId, "returned key id was empty")

	return issuerId, keyId
}
