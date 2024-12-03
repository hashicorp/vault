// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cert

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"math/big"
	mathrand "math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ocsp"
)

var ocspPort int

var source InMemorySource

type testLogger struct{}

func (t *testLogger) Log(args ...any) {
	fmt.Printf("%v", args)
}

func TestMain(m *testing.M) {
	source = make(InMemorySource)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return
	}

	ocspPort = listener.Addr().(*net.TCPAddr).Port
	srv := &http.Server{
		Addr:    "localhost:0",
		Handler: NewResponder(&testLogger{}, source, nil),
	}
	go func() {
		srv.Serve(listener)
	}()
	defer srv.Shutdown(context.Background())
	m.Run()
}

func TestCert_RoleResolve(t *testing.T) {
	certTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "example.com",
		},
		DNSNames:    []string{"example.com"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(262980 * time.Hour),
	}

	tempDir, connState, err := generateTestCertAndConnState(t, certTemplate)
	if tempDir != "" {
		defer os.RemoveAll(tempDir)
	}
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile(filepath.Join(tempDir, "ca_cert.pem"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", allowed{dns: "example.com"}, false),
			testAccStepLoginWithName(t, connState, "web"),
			testAccStepResolveRoleWithName(t, connState, "web"),
			// Test with caching disabled
			testAccStepSetRoleCacheSize(t, -1),
			testAccStepLoginWithName(t, connState, "web"),
			testAccStepResolveRoleWithName(t, connState, "web"),
		},
	})
}

func testAccStepResolveRoleWithName(t *testing.T, connState tls.ConnectionState, certName string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.ResolveRoleOperation,
		Path:            "login",
		Unauthenticated: true,
		ConnState:       &connState,
		Check: func(resp *logical.Response) error {
			if resp.Data["role"] != certName {
				t.Fatalf("Role was not as expected. Expected %s, received %s", certName, resp.Data["role"])
			}
			return nil
		},
		Data: map[string]interface{}{
			"name": certName,
		},
	}
}

func TestCert_RoleResolveWithoutProvidingCertName(t *testing.T) {
	certTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "example.com",
		},
		DNSNames:    []string{"example.com"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(262980 * time.Hour),
	}

	tempDir, connState, err := generateTestCertAndConnState(t, certTemplate)
	if tempDir != "" {
		defer os.RemoveAll(tempDir)
	}
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile(filepath.Join(tempDir, "ca_cert.pem"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", allowed{dns: "example.com"}, false),
			testAccStepLoginWithName(t, connState, "web"),
			testAccStepResolveRoleWithEmptyDataMap(t, connState, "web"),
			testAccStepSetRoleCacheSize(t, -1),
			testAccStepLoginWithName(t, connState, "web"),
			testAccStepResolveRoleWithEmptyDataMap(t, connState, "web"),
		},
	})
}

func testAccStepSetRoleCacheSize(t *testing.T, size int) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			"role_cache_size": size,
		},
	}
}

func testAccStepResolveRoleWithEmptyDataMap(t *testing.T, connState tls.ConnectionState, certName string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.ResolveRoleOperation,
		Path:            "login",
		Unauthenticated: true,
		ConnState:       &connState,
		Check: func(resp *logical.Response) error {
			if resp.Data["role"] != certName {
				t.Fatalf("Role was not as expected. Expected %s, received %s", certName, resp.Data["role"])
			}
			return nil
		},
		Data: map[string]interface{}{},
	}
}

func testAccStepResolveRoleExpectRoleResolutionToFail(t *testing.T, connState tls.ConnectionState, certName string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.ResolveRoleOperation,
		Path:            "login",
		Unauthenticated: true,
		ConnState:       &connState,
		ErrorOk:         true,
		Check: func(resp *logical.Response) error {
			if resp == nil && !resp.IsError() {
				t.Fatalf("Response was not an error: resp:%#v", resp)
			}

			errString, ok := resp.Data["error"].(string)
			if !ok {
				t.Fatal("Error not part of response.")
			}

			if _, dataKeyExists := resp.Data["data"]; dataKeyExists {
				t.Fatal("metadata key 'data' existed in response without feature enabled")
			}

			if !strings.Contains(errString, certAuthFailMsg) {
				t.Fatalf("Error was not due to invalid role name. Error: %s", errString)
			}
			return nil
		},
		Data: map[string]interface{}{
			"name": certName,
		},
	}
}

func testAccStepResolveRoleExpectRoleResolutionToFailWithData(t *testing.T, connState tls.ConnectionState, certName string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.ResolveRoleOperation,
		Path:            "login",
		Unauthenticated: true,
		ConnState:       &connState,
		ErrorOk:         true,
		Check: func(resp *logical.Response) error {
			if resp == nil && !resp.IsError() {
				t.Fatalf("Response was not an error: resp:%#v", resp)
			}

			errString, ok := resp.Data["error"].(string)
			if !ok {
				t.Fatal("Error not part of response.")
			}

			dataKeysRaw, dataKeyExists := resp.Data["data"]
			if !dataKeyExists {
				t.Fatal("metadata key 'data' did not exist in response feature enabled")
			}
			dataKeys, ok := dataKeysRaw.(map[string]string)
			if !ok {
				t.Fatalf("the 'data' field was not a map: %T", dataKeysRaw)
			}

			for _, key := range []string{"common_name", "serial_number", "authority_key_id", "subject_key_id"} {
				require.Contains(t, dataKeys, key, "response metadata key %s was missing in response: %v", key, resp)
			}

			if !strings.Contains(errString, certAuthFailMsg) {
				t.Fatalf("Error was not due to invalid role name. Error: %s", errString)
			}
			return nil
		},
		Data: map[string]interface{}{
			"name": certName,
		},
	}
}

func testAccStepResolveRoleOCSPFail(t *testing.T, connState tls.ConnectionState, certName string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.ResolveRoleOperation,
		Path:            "login",
		Unauthenticated: true,
		ConnState:       &connState,
		ErrorOk:         true,
		Check: func(resp *logical.Response) error {
			if resp == nil || !resp.IsError() {
				t.Fatalf("Response was not an error: resp:%#v", resp)
			}

			errString, ok := resp.Data["error"].(string)
			if !ok {
				t.Fatal("Error not part of response.")
			}

			if !strings.Contains(errString, certAuthFailMsg) {
				t.Fatalf("Error was not due to OCSP failure. Error: %s", errString)
			}
			return nil
		},
		Data: map[string]interface{}{
			"name": certName,
		},
	}
}

func TestCert_RoleResolve_RoleDoesNotExist(t *testing.T) {
	certTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "example.com",
		},
		DNSNames:    []string{"example.com"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(262980 * time.Hour),
	}

	tempDir, connState, err := generateTestCertAndConnState(t, certTemplate)
	if tempDir != "" {
		defer os.RemoveAll(tempDir)
	}
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile(filepath.Join(tempDir, "ca_cert.pem"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", allowed{dns: "example.com"}, false),
			testAccStepLoginWithName(t, connState, "web"),
			testAccStepResolveRoleExpectRoleResolutionToFail(t, connState, "notweb"),
		},
	})
}

func TestCert_RoleResolveOCSP(t *testing.T) {
	certTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "example.com",
		},
		DNSNames:    []string{"example.com"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(262980 * time.Hour),
		OCSPServer:   []string{fmt.Sprintf("http://localhost:%d", ocspPort)},
	}
	tempDir, connState, err := generateTestCertAndConnState(t, certTemplate)
	if tempDir != "" {
		defer os.RemoveAll(tempDir)
	}
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile(filepath.Join(tempDir, "ca_cert.pem"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	issuer := parsePEM(ca)
	pkf, err := ioutil.ReadFile(filepath.Join(tempDir, "ca_key.pem"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	pk, err := certutil.ParsePEMBundle(string(pkf))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	tempDir, connState2, err := generateTestCertAndConnState(t, certTemplate)
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca2, err := ioutil.ReadFile(filepath.Join(tempDir, "ca_cert.pem"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	issuer2 := parsePEM(ca2)
	pkf2, err := ioutil.ReadFile(filepath.Join(tempDir, "ca_key.pem"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	pk2, err := certutil.ParsePEMBundle(string(pkf2))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	type caData struct {
		privateKey crypto.Signer
		caBytes    []byte
		caChain    []*x509.Certificate
		connState  tls.ConnectionState
	}

	ca1Data := caData{
		pk.PrivateKey,
		ca,
		issuer,
		connState,
	}
	ca2Data := caData{
		pk2.PrivateKey,
		ca2,
		issuer2,
		connState2,
	}

	cases := []struct {
		name        string
		failOpen    bool
		certStatus  int
		errExpected bool
		caData      caData
		ocspCaCerts string
	}{
		{name: "failFalseGoodCert", certStatus: ocsp.Good, caData: ca1Data},
		{name: "failFalseRevokedCert", certStatus: ocsp.Revoked, errExpected: true, caData: ca1Data},
		{name: "failFalseUnknownCert", certStatus: ocsp.Unknown, errExpected: true, caData: ca1Data},
		{name: "failTrueGoodCert", failOpen: true, certStatus: ocsp.Good, caData: ca1Data},
		{name: "failTrueRevokedCert", failOpen: true, certStatus: ocsp.Revoked, errExpected: true, caData: ca1Data},
		{name: "failTrueUnknownCert", failOpen: true, certStatus: ocsp.Unknown, caData: ca1Data},
		{name: "failFalseGoodCertExtraCas", certStatus: ocsp.Good, caData: ca2Data, ocspCaCerts: string(pkf2)},
		{name: "failFalseRevokedCertExtraCas", certStatus: ocsp.Revoked, errExpected: true, caData: ca2Data, ocspCaCerts: string(pkf2)},
		{name: "failFalseUnknownCertExtraCas", certStatus: ocsp.Unknown, errExpected: true, caData: ca2Data, ocspCaCerts: string(pkf2)},
		{name: "failTrueGoodCertExtraCas", failOpen: true, certStatus: ocsp.Good, caData: ca2Data, ocspCaCerts: string(pkf2)},
		{name: "failTrueRevokedCertExtraCas", failOpen: true, certStatus: ocsp.Revoked, errExpected: true, caData: ca2Data, ocspCaCerts: string(pkf2)},
		{name: "failTrueUnknownCertExtraCas", failOpen: true, certStatus: ocsp.Unknown, caData: ca2Data, ocspCaCerts: string(pkf2)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resp, err := ocsp.CreateResponse(c.caData.caChain[0], c.caData.caChain[0], ocsp.Response{
				Status:       c.certStatus,
				SerialNumber: certTemplate.SerialNumber,
				ProducedAt:   time.Now(),
				ThisUpdate:   time.Now(),
				NextUpdate:   time.Now().Add(time.Hour),
			}, c.caData.privateKey)
			if err != nil {
				t.Fatal(err)
			}
			source[certTemplate.SerialNumber.String()] = resp

			b := testFactory(t)
			b.(*backend).ocspClient.ClearCache()
			var resolveStep logicaltest.TestStep
			var loginStep logicaltest.TestStep
			if c.errExpected {
				loginStep = testAccStepLoginWithNameInvalid(t, c.caData.connState, "web")
				resolveStep = testAccStepResolveRoleOCSPFail(t, c.caData.connState, "web")
			} else {
				loginStep = testAccStepLoginWithName(t, c.caData.connState, "web")
				resolveStep = testAccStepResolveRoleWithName(t, c.caData.connState, "web")
			}
			logicaltest.Test(t, logicaltest.TestCase{
				CredentialBackend: b,
				Steps: []logicaltest.TestStep{
					testAccStepCertWithExtraParams(t, "web", c.caData.caBytes, "foo", allowed{dns: "example.com"}, false,
						map[string]interface{}{"ocsp_enabled": true, "ocsp_fail_open": c.failOpen, "ocsp_ca_certificates": c.ocspCaCerts}),
					testAccStepReadCertPolicy(t, "web", false, map[string]interface{}{"ocsp_enabled": true, "ocsp_fail_open": c.failOpen, "ocsp_ca_certificates": c.ocspCaCerts}),
					loginStep,
					resolveStep,
				},
			})
		})
	}
}

// TestCert_MetadataOnFailure verifies that we return the cert metadata
// in the response on failures if the configuration option is enabled.
func TestCert_MetadataOnFailure(t *testing.T) {
	certTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "example.com",
		},
		DNSNames:    []string{"example.com"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(262980 * time.Hour),
	}

	tempDir, connState, err := generateTestCertAndConnState(t, certTemplate)
	if tempDir != "" {
		defer os.RemoveAll(tempDir)
	}
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile(filepath.Join(tempDir, "ca_cert.pem"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testStepEnableMetadataFailures(),
			testAccStepCert(t, "web", ca, "foo", allowed{dns: "example.com"}, false),
			testAccStepLoginWithName(t, connState, "web"),
			testAccStepResolveRoleExpectRoleResolutionToFailWithData(t, connState, "notweb"),
		},
	})
}
